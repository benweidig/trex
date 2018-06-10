# T-Rex, the Tree Explorer [![Build Status](https://travis-ci.org/benweidig/trex.svg?branch=master)](https://travis-ci.org/benweidig/trex)

CLI tool for visualizing JSON/YAML files.

## Caveats

This project isn't finished yet, and is not handling big files very well...

## Install

You can either build from source, use the .deb-files. Will be added to my [homebrew tap](https://github.com/benweidig/homebrew-tap) after the next release maybe.

## Usage
```
trex [-m/--monochrome] [<filepath>]
```

## Arguments

| Argument          | Default | Description                        |
| ----------------- | ------- | ---------------------------------- |
| -m / --monochrome | false   | Don't use ANSI colors              |

ANSI colors might be disabled automatically if the terminal doesn't seem to support it, but the detection is not perfect.

## Vendoring

The project has a custom version of [rivo/tview](https://github.com/rivo/tview) vendored, due to an open PR from me, will be removed when merged.

## License

MIT. See [LICENSE](LICENSE).
