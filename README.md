# pushbell

![GitHub Release](https://img.shields.io/github/v/release/gootsolution/pushbell)
![GitHub License](https://img.shields.io/github/license/gootsolution/pushbell)
![Go Report Card](https://goreportcard.com/badge/github.com/gootsolution/pushbell)
[![Lint & Test](https://github.com/gootsolution/pushbell/actions/workflows/lint-and-test.yml/badge.svg)](https://github.com/gootsolution/pushbell/actions/workflows/lint-and-test.yml)

pushbell is a Go library for sending web push notifications with support
for the VAPID (Voluntary Application Server Identification) specification.

## Features

- Full implementation of the [encryption](https://datatracker.ietf.org/doc/html/rfc8291)
  and [Web Push](https://datatracker.ietf.org/doc/html/rfc8030) specification,
  including [VAPID](https://datatracker.ietf.org/doc/html/rfc8292).
- Support for multiple push services (Firefox, Chrome, etc.).
- Use [fasthttp](https://github.com/valyala/fasthttp) client.
- Simple and intuitive API.

## Installation

```shell
go get -u github.com/gootsolution/pushbell
```

## Example

```go
package main

import (
	"github.com/gootsolution/pushbell"
)

func main() {
	applicationServerPrivateKey := "QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ"
	applicationServerPublicKey := "BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY"

	opts := pushbell.NewOptions().ApplyKeys(applicationServerPublicKey, applicationServerPrivateKey)

	pb, err := pushbell.NewService(opts)
	if err != nil {
		panic(err)
	}

	subscriptionEndpoint := "https://fcm.googleapis.com/fcm/send/e2CN0r8ft38:APA91bES3NaBHe_GgsRp_3Ir7f18L38wA5XYRoqZCbjMPEWnkKa07uxheWE5MGZncsPOr0_34zLaFljVqmNqW76KhPSrjdy_pdInnHPEIYAZpdcIYk8oIfo1F_84uKMSqIDXRhngL76S"
	subscriptionAuth := "rm_owGF0xliyVXsrZk1LzQ"
	subscriptionP256DH := "BKm5pKbGwkTxu7dJuuLyTCBOCuCi1Fs01ukzjUL5SEX1-b-filqeYASY6gy_QpPHGErGqAyQDYAtprNWYdcsM3Y"

	message := []byte("{\"title\": \"My first message\"}")

	err = pb.Send(&pushbell.Push{
		Endpoint:  subscriptionEndpoint,
		Auth:      subscriptionAuth,
		P256DH:    subscriptionP256DH,
		Plaintext: message,
	})
	if err != nil {
		panic(err)
	}
}

```

**NOTE:** You can use [this](https://gootsolution.github.io/pushbell/) to play around and make tests without your
service workers.

## Documentation

Detailed API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/gootsolution/pushbell).

## License

This project is distributed under the MIT License.

## Support and Contributing

If you have any questions, issues, or suggestions, please create an issue in this repository. Pull requests with fixes
and improvements are also welcome.
