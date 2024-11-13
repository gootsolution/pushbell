package pushbell

import (
	"net/http"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/gootsolution/pushbell/internal/httpclient"
)

type Options struct {
	ApplicationServerPublicKey  string            // [RFC 8292] ECDH public key.
	ApplicationServerPrivateKey string            // [RFC 8292] ECDH private key.
	ApplicationServerSubject    string            // [RFC 8292] Either a "mailto:" (email) or a "https:" URI.
	HttpClient                  httpclient.Client // Custom client for request.
	CheckStatusCode             bool              // If true, return error based on status code.
	KeyRotationInterval         time.Duration     // If set, enable encryption keys rotation.
}

func NewOptions() *Options {
	return &Options{
		ApplicationServerSubject: "https://github.com/gootsolution/pushbell",
	}
}

func (o *Options) ApplyKeys(publicKey, privateKey string) *Options {
	o.ApplicationServerPublicKey = publicKey
	o.ApplicationServerPrivateKey = privateKey

	return o
}

func (o *Options) SetSubject(subject string) *Options {
	o.ApplicationServerSubject = subject

	return o
}

func (o *Options) SetCheckStatusCode(checkStatusCode bool) *Options {
	o.CheckStatusCode = checkStatusCode

	return o
}

func (o *Options) SetKeyRotation(interval time.Duration) *Options {
	o.KeyRotationInterval = interval

	return o
}

func (o *Options) SetFastHttpClient(client *fasthttp.Client) *Options {
	o.HttpClient = httpclient.FastHttp(client)

	return o
}

func (o *Options) SetStdHttpClient(client *http.Client) *Options {
	o.HttpClient = httpclient.StdHttp(client)

	return o
}
