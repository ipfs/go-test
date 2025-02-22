# random-files - create random filesystem hierarchies

`random-files` creates random filesystem hierarchies for testing

## Install

```
go install github.com/ipfs/go-test/cli/random-files
```

## Usage

```sh
> random-files -help
NAME
  random-files - Write a random filesystem hierarchy to each <path>

USAGE
  random-files [options] <path>...

OPTIONS:
  -depth int
        depth of the directory tree including the root directory (default 2)
  -dirs int
        number of subdirectories at each depth (default 5)
  -files int
        number of files at each depth (default 10)
  -filesize int
        bytes of random data in each file (default 4096)
  -q    do not print files and directories
  -random-dirs
        randomize number of subdirectories, from 1 to -dirs
  -random-files
        randomize number of files, from 1 to -files
  -random-size
        randomize file size, from 1 to -filesize (default true)
  -seed int
        random seed, 0 for current time
```

## Examples

```sh
> random-files -depth=2 -files=3 -seed=1701 foo
foo/rwd67uvnj9yz-
foo/7vovyvr9
foo/fjv0w0
foo/gyubi50rec5/
foo/gyubi50rec5/vr6x-ce4uupj
foo/gyubi50rec5/ob9ud0e8lt_2e
foo/gyubi50rec5/11gip6zea
foo/nzu5j29-sh-ku4/
foo/nzu5j29-sh-ku4/vcs1629n
foo/nzu5j29-sh-ku4/rky_i_qsxrp
foo/nzu5j29-sh-ku4/xr1usy5ic0
foo/w30dzrx2w4_d/
foo/w30dzrx2w4_d/7ued6
foo/w30dzrx2w4_d/r1d3j
foo/w30dzrx2w4_d/av7d09i-av
foo/s6ha-58/
foo/s6ha-58/nukjsxg7t
foo/s6ha-58/7of_84
foo/s6ha-58/h0jgq8mu1n7u
foo/tq_8/
foo/tq_8/sx-a2jgmz_mk6
foo/tq_8/9hzrksz8
foo/tq_8/8b5swu
```

It made:

```sh
> tree foo
foo
├── 7vovyvr9
├── fjv0w0
├── gyubi50rec5
│   ├── 11gip6zea
│   ├── ob9ud0e8lt_2e
│   └── vr6x-ce4uupj
├── nzu5j29-sh-ku4
│   ├── rky_i_qsxrp
│   ├── vcs1629n
│   └── xr1usy5ic0
├── rwd67uvnj9yz-
├── s6ha-58
│   ├── 7of_84
│   ├── h0jgq8mu1n7u
│   └── nukjsxg7t
├── tq_8
│   ├── 8b5swu
│   ├── 9hzrksz8
│   └── sx-a2jgmz_mk6
└── w30dzrx2w4_d
    ├── 7ued6
    ├── av7d09i-av
    └── r1d3j

6 directories, 18 files
```

Note: Specifying the same seed will produce the same results.


### Acknowledgments

Credit to [Juan Benet](https://github.com/jbenet) as the author of [`go-random-files`](https://github.com/jbenet/go-random-files) from which this code was derived.
