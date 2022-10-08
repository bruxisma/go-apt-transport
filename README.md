# Overview

[![Go Reference][godoc-svg]][godoc-lnk]
[![Coverage][codecov-svg]][codecov-lnk]

This library provides a common set of operations for writing `apt` transport
methods in golang. It does not handle the exact logic that covers *every*
possible transport protocol, but does allow for an easy to use interface when
communicating with `apt` or `apt-get`.

In addition to providing a simple interface for acquiring resources in custom
transport methods, this API also exports all pieces necessary to recreate the
common interfaces. This is done to provide users with more granular support
over `Message` deserialization, `Method` acquisition, and introspection.

# Usage

To add this as a go module simply do

```console
$ go get occult.work/apt/transport@latest
```

**NOTE**: There is no `occult.work/apt` package.

# Testing

To run unit tests, simply run `go test ./...`. All mocking, testing data, etc.
is taken care of. For more "fun" output, users can use
[`gotestfmt`](https://github.com/haveyoudebuggedit/gotestfmt)

```console
$ go test -v ./... -json -cover ${PWD} 2>&1 | gotestfmt
```

[godoc-svg]: https://pkg.go.dev/badge/occult.work/apt/transport.svg
[godoc-lnk]: https://pkg.go.dev/occult.work/apt/transport

[codecov-svg]: https://codecov.io/gh/bruxisma/go-apt-transport/branch/main/graph/badge.svg
[codecov-lnk]: https://codecov.io/gh/bruxisma/go-apt-transport
