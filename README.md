# Bindings for Monobank API

[![GoDoc](https://godoc.org/github.com/kudrykv/go-monobank-api?status.svg)](https://godoc.org/github.com/kudrykv/go-monobank-api)
[![CI](https://github.com/kudrykv/go-monobank-api/workflows/CI/badge.svg)](https://github.com/kudrykv/go-monobank-api/actions?query=workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/kudrykv/go-monobank-api)](https://goreportcard.com/report/github.com/kudrykv/go-monobank-api)
[![codecov](https://codecov.io/gh/kudrykv/go-monobank-api/branch/master/graph/badge.svg)](https://codecov.io/gh/kudrykv/go-monobank-api)

---

The library is the client to work with Monobank API.
One of its features is no dependencies on 3rd-party libraries.

## Usage

```go
package main

import (
  "context"
  "fmt"

  mono "github.com/kudrykv/go-monobank-api"
)

func main() {
  // public can be initialized with optional parameters:
  //   mono.NewPublic(
  //     mono.WithDomain("custom-domain"),
  //     mono.WithClient(&http.Client{}),
  //     mono.WithUnmarshaller(customUnmarshaller),
  //   )
  public := mono.NewPublic()

  currencies, err := public.Currency(context.Background())
  if err != nil {
    panic(err)
  }

  fmt.Println(currencies)

  // NewPersonal can be configured in the same way the public client.
  private := mono.NewPersonal("api-token")

  info, err := private.ClientInfo(context.Background())
  if err != nil {
    panic(err)
  }

  fmt.Println(info)
}
```

## Monobank API Documentation
https://api.monobank.ua/docs/

You can obtain your personal token [here](https://api.monobank.ua).

## Progress
- [x] Public API
- [x] Personal API
- [ ] Corporate API (to do)