# random-data - writes random data to stdout

`random-data` writes pseudo-random data to stdout for testing. The data can be written as raw bytes or base64 encodde.

## Install

```
go install github.com/ipfs/go-test/cli/random-data
```

## Usage

```sh
> random-bytes -help
NAME
  ./random-data - Write random data to stdout

USAGE
  ./random-data [options]

OPTIONS:
  -b64
        base-64 encode output
  -seed int
        random seed, 0 or unset for current time
  -size int
        number of bytes to generate
```

## Examples

```sh
random-data -size=64 -b64
vRujjyEvx8lYiELflaDINvkm5nfueWGCdzEOxhRtz7N2EQjoyrpoMdVVOrwAgNO0tVojDAgu0JpU4hKSsdVl8A==
```

Note: Specifying the same seed will produce the same results.
