# fin

Go library for financial calculations.

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![License Badge][license badge]][LICENSE]

## Overview

[fin][] provides the following financial calculations:

- Various financial ratios (e.g., ROIC, ROE, TIE)
- Internal Rate of Return (IRR) & Modified Internal Rate of Return (MIRR)
- Net Present Value (NPV)
- Payback Period & Discounted Payback Period
- Monte Carlo Simulation (MCS) — not fully implemented

## Installation

```bash
$ go get github.com/goinvest/fin
```


## Contributing

Contributions are welcome! To contribute please:

1. Fork the repository
2. Create a feature branch
3. Code
4. Submit a [pull request][]

### Testing

Prior to submitting a [pull request][], please run:

```bash
$ make check
```

To update and view the test coverage report:

```bash
$ make cover
```

## License

[fin][] is released under the MIT license. Please see the [LICENSE][] file for
more information.

[fin]: https://github.com/goinvest/fin
[godoc badge]: https://godoc.org/github.com/goinvest/fin?status.svg
[godoc link]: https://godoc.org/github.com/goinvest/fin
[LICENSE]: https://github.com/goinvest/fin/blob/master/LICENSE
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/goinvest/fin
[report card]: https://goreportcard.com/report/github.com/goinvest/fin
