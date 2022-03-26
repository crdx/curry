# curry

**curry** is a currency conversion tool. It uses the exchangerates.io historical API.

Real-time rates are not available on the free plan, but historical ones are, so yesterday's rates are used for conversions.

## Installation

```bash
go install github.com/crdx/curry@latest
```

## Configuration

Place a valid exchangerates.io API key in `~/.config/curry/api_key`.

## CLI

```
Usage:
    curry raw
    curry list
    curry clear
    curry [options] <value>

Convert currencies.

Commands:
    raw      Print raw json
    clear    Clear the cache
    list     List available currencies

Options:
    -f, --from TYPE    From this currency [default: GBP]
    -t, --to TYPE      To this currency [default: GBP]
    -q, --quiet        Show only the value
    -C, --no-color     Disable colours
    -h, --help         Show help
```

## Examples

Convert 10 USD to GBP.

```bash
curry -f USD -t GBP 10
```

View a list of currencies.

```bash
curry list
```

## Commands

### Convert

Use `-f/--from` and `-t/--to` to specify source and target currency identifiers. Use `-q/--quiet` to output only the value.

### Raw

Show the raw JSON returned from the API.

### List

List available currency identifiers.

### Clear

Delete the current rates cache (see below).

## Cache

The first call will fetch yesterday's rates and cache them locally, using the cache for any further calls.

Cached rates are stored at `~/.cache/curry/<YYYY-mm-dd>.json`.

## Contributions

Open an [issue](https://github.com/crdx/curry/issues) or send a [pull request](https://github.com/crdx/curry/pulls).

## Licence

[MIT](LICENCE.md).
