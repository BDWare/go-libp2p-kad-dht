# go-libp2p-kad-dht

Fork of [libp2p/go-libp2p-kad-dht](https://github.com/libp2p/go-libp2p-kad-dht) used by BDWare projects.

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-BDWare-yellow.svg?style=flat-square)](https://bdware.org/))
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/bdware/go-libp2p-kad-dht)
[![Build Status](https://travis-ci.com/BDWare/go-libp2p-kad-dht.svg?branch=master)](https://travis-ci.com/BDWare/go-libp2p-kad-dht)

> A Kademlia DHT implementation on go-libp2p

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Contribute](#contribute)
- [License](#license)

## Install

```sh
go get github.com/bdware/go-libp2p-kad-dht
```

Add the following lines to your `go.mod`:
```
replace github.com/libp2p/go-libp2p-kad-dht => github.com/bdware/go-libp2p-kad-dht {{KAD_DHT_VERSION}}
replace github.com/libp2p/go-libp2p-kbucket => github.com/bdware/go-libp2p-kbucket {{KBUCKET_VERSION}}
``` 

## Usage

Go to https://godoc.org/github.com/bdware/go-libp2p-kad-dht.

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/libp2p/go-libp2p-kad-dht/issues).

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

[MIT](LICENSE)

Copyright for portions of this fork are held by [Protocol Labs Inc., 2016] as
part of the original [go-libp2p-kad-dht](https://github.com/libp2p/go-libp2p-kad-dht) project.

All other copyright for this fork are held by [The BDWare Authors, 2020].

All rights reserved.

---

The last gx published version of this module was: 4.4.34: QmXuNFLZc6Nb5akB4sZsxK3doShsFKT1sZFvxLXJvZQwAW
