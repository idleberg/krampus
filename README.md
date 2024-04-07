# :japanese_ogre: krampus

> Command-line tool to kill one or more processes by their port number.

![License](https://img.shields.io/github/license/idleberg/krampus?style=for-the-badge)
![Version](https://img.shields.io/github/v/release/idleberg/krampus?sort=semver&style=for-the-badge)
[![Build](https://img.shields.io/github/actions/workflow/status/idleberg/node-dent/default.yml?style=for-the-badge)](https://github.com/idleberg/krampus/actions)

## Installation

### Homebrew

```sh
$ brew install idleberg/asahi/krampus --build-from-source
```

### Go

```sh
$ go install github.com/idleberg/krampus
```

:warning: Make sure that your Go binaries path (usually `$HOME/go/bin` or `%USERPROFILE%\go\bin` on Windows) is in your system's `PATH`.

## Usage

```sh
// Let's kill your web-server
$ krampus 80 443
```

## Benchmark

```
Benchmark 1: node-krampus 3333
  Time (mean ± σ):     869.7 ms ±   6.0 ms    [User: 775.7 ms, System: 158.6 ms]
  Range (min … max):   861.7 ms … 884.3 ms    10 runs
 
Benchmark 2: go-krampus 3333
  Time (mean ± σ):      61.5 ms ±   2.5 ms    [User: 37.8 ms, System: 19.5 ms]
  Range (min … max):    58.7 ms …  73.3 ms    43 runs
 
Summary
  go-krampus 3333 ran
   14.15 ± 0.58 times faster than node-krampus 3333
```

## Credit

Name and idea inspired by Mario Nebl's [krampus](https://www.npmjs.com/package/krampus) for NodeJS.

## License

This work is licensed under [The MIT License](LICENSE).
  
