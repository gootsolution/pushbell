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

import (
	"errors"
	"log"
	"time"

	"github.com/gootsolution/pushbell"
)

func main() {
	applicationServerPrivateKey := "QxfAyO5dMMrSvDT2_xHxW5aktYPWGE_hT42RKlHilpQ"
	applicationServerPublicKey := "BIRM67G3W1fva-ephDo220BbiaOOy-SBk2uzHsmlqMXp_OmkKxYW96cOK5EWnKdkLg2i7N4FYfuxIwm7JWThVSY"
	applicationServerSubject := "mailto:webpush@example.com"

	pb, err := pushbell.New(applicationServerPrivateKey, applicationServerPublicKey, applicationServerSubject)
	if err != nil {
		panic(err)
	}

	subscriptionEndpoint := "https://fcm.googleapis.com/fcm/send/e2CN0r8ft38:APA91bES3NaBHe_GgsRp_3Ir7f18L38wA5XYRoqZCbjMPEWnkKa07uxheWE5MGZncsPOr0_34zLaFljVqmNqW76KhPSrjdy_pdInnHPEIYAZpdcIYk8oIfo1F_84uKMSqIDXRhngL76S"
	subscriptionAuth := "rm_owGF0xliyVXsrZk1LzQ"
	subscriptionP256DH := "BKm5pKbGwkTxu7dJuuLyTCBOCuCi1Fs01ukzjUL5SEX1-b-filqeYASY6gy_QpPHGErGqAyQDYAtprNWYdcsM3Y"

	message := []byte("{\"title\": \"My first message\"}")

	if err = pb.Send(
		subscriptionEndpoint,
		subscriptionAuth,
		subscriptionP256DH,
		message,
		pushbell.UrgencyHigh,
		time.Hour,
	); err != nil {
		switch {
		case errors.Is(err, pushbell.ErrPushGone):
			log.Println(err)
		default:
			panic(err)
		}
	}
}
```

**NOTE:** You can use [this](https://gootsolution.github.io/pushbell/) to play around and make tests without your
service workers.

## Documentation

Detailed API documentation is available on GoDoc.

## License

This project is distributed under the MIT License.

## Support and Contributing

If you have any questions, issues, or suggestions, please create an issue in this repository. Pull requests with fixes
and improvements are also welcome.