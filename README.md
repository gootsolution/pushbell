![GitHub Tag](https://img.shields.io/github/v/tag/gootsolution/pushbell?style=flat)
![GitHub License](https://img.shields.io/github/license/gootsolution/pushbell)
![Go Report Card](https://goreportcard.com/badge/github.com/gootsolution/pushbell)

# pushbell

pushbell is a Go library for sending web push notifications with support
for the VAPID (Voluntary Application Server Identification) specification.

## Features

- Full implementation of the [encryption](https://datatracker.ietf.org/doc/html/rfc8291)
  and [Web Push](https://datatracker.ietf.org/doc/html/rfc8030) specification,
  including [VAPID](https://datatracker.ietf.org/doc/html/rfc8292).
- Support for multiple push services (Firefox, Chrome, etc.)
- Use [fasthttp](https://github.com/valyala/fasthttp) client
- Simple and intuitive API
- Error handling and retries on failures
- Compatibility with different Go versions

## Installation

```shell
go get -u github.com/gootsolution/pushbell
```

## Example

```go
package main
```

## Documentation

Detailed API documentation is available on GoDoc.

## License

This project is distributed under the MIT License.

## Support and Contributing

If you have any questions, issues, or suggestions, please create an issue in this repository. Pull requests with fixes
and improvements are also welcome.