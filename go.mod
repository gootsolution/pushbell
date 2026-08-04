module github.com/gootsolution/pushbell

go 1.26

retract (
	v1.0.1 // Broken key decoder
	v1.0.0 // Broken encryption service
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/valyala/fasthttp v1.73.0
	golang.org/x/crypto v0.54.0
)

require (
	github.com/andybalholm/brotli v1.2.2 // indirect
	github.com/klauspost/compress v1.19.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
)
