# tapasd

A concurrent [Ruby Tapas][ruby_tapas] episode downloader, curiously written in Go. tapasd is focused on fetching the movie files for local mirroring and loading onto tablets for viewing later offline.

Features:

* Simple runtime dependencies--in other words, none
* Architected as a [twelve-factor app][12factor], therefore easy to deploy with [Docker][docker]
* Defaults to a process which rechecks the XML feed every 6 hours
* Maintains no state except for the downloaded content

[![baby-gopher](https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)

## Usage

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
