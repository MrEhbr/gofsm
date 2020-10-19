# gofsm

![CI](https://github.com/MrEhbr/gofsm/workflows/CI/badge.svg)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/MrEhbr/gofsm/blob/master/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/MrEhbr/gofsm.svg)](https://github.com/MrEhbr/gofsm/releases)
[![codecov](https://codecov.io/gh/MrEhbr/gofsm/branch/master/graph/badge.svg)](https://codecov.io/gh/MrEhbr/gofsm)
![Made by Alexey Burmistrov](https://img.shields.io/badge/made%20by-Alexey%20Burmistrov-blue.svg?style=flat)

gofsm is a command line tool that generates finite state machine for Go struct.

## Usage

```console
Usage: gofsm gen -p ./examples/transitions -s Order -f State -o order_fsm.go -t ./examples/transitions/transitions.json
   --package, -p      package where struct is located (default: default is current dir(.))
   --struct, -s       struct name
   --field, -f        state field of struct
   --output, -o       output file name (default: default srcdir/<struct>_fsm.go)
   --transitions, -t  path to file with transitions
   --noGenerate, -g   don't put //go:generate instruction to the generated code (default: false)
```

This will generate [finite state machine](./examples/transitions/order_fsm.go) for struct Order with transitions defined in [./examples/transitions/transitions.json](./examples/transitions/transitions.json) file

## Install

### Using go

```console
go get -u github.com/MrEhbr/gofsm/cmd/gofsm
```

### Download releases

<https://github.com/MrEhbr/gofsm/releases>

## License

© 2020 [Alexey Burmistrov]

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE`](LICENSE)). See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: Apache-2.0`
