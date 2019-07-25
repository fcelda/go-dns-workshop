# DNS Workshop with Go

## Requirements

Go version with modules support.

## Quick Start

Clone the repository and download module dependencies:

```sh
$ cd ~/devel
$ git clone https://github.com/fcelda/go-dns-workshop
$ cd go-dns-workshop
$ go mod download
```

Build and start the server:

```sh
$ go build ./cmd/server
$ ./server
```

Build and test the client:

```sh
$ go build ./cmd/client
$ ./client -name hello.test -type TXT
```

## Resources

Documentation:

- [miekg/dns documentation](https://godoc.org/github.com/miekg/dns)
- [Hello, and welcome to DNS!](https://powerdns.org/hello-dns/)

DNS protocol specification:

- [RFC 1034: DNS Concepts and Facialities](https://tools.ietf.org/html/rfc1034)
- [RFC 1035: DNS Implementation and Specification](https://tools.ietf.org/html/rfc1035)
- [RFC 8499: DNS Terminology](https://tools.ietf.org/html/rfc8499)
- [IANA: DNS Parameters](https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml)
