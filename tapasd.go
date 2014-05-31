package main

import (
	"encoding/xml"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

const feedUrl = "https://rubytapas.dpdcart.com/feed"

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
}

type Item struct {
	Title     string    `xml:"title"`
	Enclosure Enclosure `xml:"enclosure"`
}

func main() {
	defDataDir, _ := os.Getwd()
	user := flag.String("user", "[required]", "user for RubyTapas account (required)")
	pass := flag.String("pass", "[required]", "pass for RubyTapas account (required)")
	dataDir := flag.String("data", defDataDir, "data directory for downloads")
	concurrency := flag.Int("concurrency", 4, "data directory for downloads")
	daemon := flag.Bool("daemon", false, "daemon mode which checks feed periodically")
	interval := flag.Int("interval", 60*60*6, "number of seconds to sleep between reprocessing")
	flag.Parse()
	if *user == "[required]" || *pass == "[required]" {
		log.Fatal("-user and -pass flags are required")
	}

	for {
		items := make(chan Item)
		go generate(items, *user, *pass)
		process(items, *concurrency, *user, *pass, *dataDir)
		log.Println("Processing and downloading complete")
		if !*daemon {
			break
		}
		log.Printf("Sleeping for %d seconds\n", *interval)
		timer := time.NewTimer(time.Duration(*interval) * time.Second)
		<-timer.C
		log.Println("Waking to process and download")
	}
}

func generate(urls chan Item, user string, pass string) {
	defer close(urls)
	request, err := http.NewRequest("GET", feedUrl, nil)
	request.SetBasicAuth(user, pass)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("URL %s returned status: %s\n", feedUrl, response.Status)
	}

	decoder := xml.NewDecoder(response.Body)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "item" {
				var item Item
				decoder.DecodeElement(&item, &startElement)
				urls <- item
			}
		}
	}
}

func process(items chan Item, concurrency int, user string, pass string, dataDir string) {
	var wg sync.WaitGroup
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		worker := i
		go func() {
			defer wg.Done()
			for item := range items {
				download(worker, item, user, pass, dataDir)
			}
		}()
	}
	wg.Wait()
}

func download(worker int, item Item, user string, pass string, dataDir string) {
	url, err := url.Parse(item.Enclosure.Url)
	if err != nil {
		log.Fatal(err)
	}
	fname := path.Join(dataDir, "rubytapas-"+slugify(item.Title)+path.Ext(url.Path))

	existing, err := os.Stat(fname)
	if !os.IsNotExist(err) {
		if existing.Size() != item.Enclosure.Length {
			log.Printf("[%02d] Deleting and retrying %s, incorrect file size (expected: %d, actual: %d)", worker, fname, item.Enclosure.Length, existing.Size())
			os.Remove(fname)
		} else {
			log.Printf("[%02d] Skipping %s", worker, fname)
			return
		}
	}

	log.Printf("[%02d] Downloading: %s to %s\n", worker, url.String(), fname)
	request, err := http.NewRequest("GET", url.String(), nil)
	request.SetBasicAuth(user, pass)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Printf("[%02d] URL %s returned status: %s\n", worker, url, response.Status)
		return
	}

	tempfile := path.Join(path.Dir(fname), "."+path.Base(fname))
	if _, err := os.Stat(tempfile); !os.IsNotExist(err) {
		os.Remove(tempfile)
	}
	temp, err := os.Create(tempfile)
	if err != nil {
		log.Fatal(err)
	}
	defer temp.Close()
	defer os.Remove(tempfile)

	written, err := io.Copy(temp, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if written != item.Enclosure.Length {
		log.Printf("[%02d] Deleting %s, incorrect file size (expected: %d, actual: %d)", worker, fname, item.Enclosure.Length, written)
		os.Remove(fname)
		return
	}
	err = os.Rename(tempfile, fname)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[%02d] Download complete: %s\n", worker, fname)
}

func slugify(title string) string {
	invalidSlugPatterns := regexp.MustCompile(`[^a-z0-9 _-]`)
	whitespacePatterns := regexp.MustCompile(`\s+`)

	result := strings.ToLower(title)
	result = invalidSlugPatterns.ReplaceAllString(result, "")
	result = whitespacePatterns.ReplaceAllString(result, "-")
	return result
}
