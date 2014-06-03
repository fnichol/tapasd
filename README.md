# tapasd

A concurrent [Ruby Tapas][ruby_tapas] episode downloader, curiously written in Go. tapasd is focused on fetching the movie files for local mirroring and loading onto tablets for viewing later offline.

Features:

* Simple runtime dependencies--in other words, none
* Architected as a [twelve-factor app][12factor], therefore easy to deploy with [Docker][docker]
* Defaults to a process which rechecks the XML feed every 6 hours
* Maintains no state except for the downloaded content

[![baby-gopher](https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)

## Usage

### Vanilla

To perform a download of all episodes in the feed to the current directory with 4 workers which re-checks every 6 hours, you only need to provide 2 things: a username and password for the DPD site:

```sh
tapasd -user=sally@example.com -pass=secret
```

If you'd prefer to perform the check-and-download once, then use the `-oneshot` flag:

```sh
tapasd -user=sally@example.com -pass=secret -oneshot
```

For more details on the other options and modes, run `tapasd -h`:

```
Usage of tapasd:
  -concurrency=4: data directory for downloads
  -data="/Users/fnichol/Projects/go/src/github.com/fnichol/tapasd": data directory for downloads
  -interval=21600: number of seconds to sleep between retrying
  -oneshot=false: check and download once, then quit
  -pass="[required]": pass for RubyTapas account (required)
  -user="[required]": user for RubyTapas account (required)
```

### Docker

In our example, we'll use a [data container pattern](http://docs.docker.io/use/working_with_volumes/) to persist the downloaded episodes between re-launching of `tapasd` services. First, we'll create a named container (`tapasd_data`) with a volume of `/data`:

```sh
docker run -v /data --name tapasd_data busybox true
```

Next, we'll launch a `tapasd` service, mounting in the shared volume from our `tapasd_data` container. For good measure, we'll also give this container a name of `tapasd`:

```sh
docker run -d --volumes-from tapasd_data --name tapasd fnichol/tapasd -user="user@example.com" -pass="secret"
```

Finally, if you want to check on its progress, simply follow the log output from the `tapasd` container:

```sh
docker logs -f tapasd
```

For a bonus, you can launch an interactive container with the data mounted in `/data` with:

```sh
docker run --rm -t -i --volumes-from tapasd_data busybox sh
```

Killing off the `tapasd` service is easy:

```sh
docker kill tapasd
```

If you wanted to re-start it at a later date:

```sh
docker start tapasd
```

And to remove the `tapasd` container:

```sh
docker rm tapasd
```

Again, your data is persisted in the volume associated with the `tapasd_data` container. To free up disk space by removing the data, simply:

```sh
docker rm tapasd_data
```

## Installation

### Source

#### Clone

```sh
mkdir -p $GOPATH/src/github.com/fnichol
cd $GOPATH/src/github.com/fnichol
git clone https://github.com/fnichol/tapasd.git
```

#### Build

```sh
cd $GOPATH/src/github.com/fnichol/tapasd
./build
```

This will generate a binary called `./bin/tapasd`.

### Docker

A Docker trusted build exists at [fnichol/tapasd](https://index.docker.io/u/fnichol/tapasd/), and can be pulled down with:

```sh
docker pull fnichol/tapasd
```

## Development

* Source host at [GitHub][repo]
* Report issues/questions/feature requests on [GitHub Issues][issues]

Pull requests are very welcome! Make sure your patches are well tested. Ideally create a topic branch for every separate change you make. For example:

1. Fork the repo
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add support for feature'`)
4. Push to the branch (`git push origin my-new-feature`)

## Authors

Created and maintained by [Fletcher Nichol][fnichol] (<fnichol@nichol.ca>).

## License

MIT (see [License.txt][license])

[fnichol]:  https://github.com/fnichol
[repo]:     https://github.com/fnichol/tapasd
[issues]:   https://github.com/fnichol/tapasd/issues
[license]:  https://github.com/fnichol/tapasd/blob/master/License.txt

[12factor]:   http://12factor.net/
[docker]:     https://www.docker.io/
[ruby_tapas]: http://www.rubytapas.com/
