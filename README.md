# Bindings for Monobank API

[![GoDoc](https://godoc.org/github.com/kudrykv/go-monobank-api?status.svg)](https://godoc.org/github.com/kudrykv/go-monobank-api)
[![CI](https://github.com/kudrykv/go-monobank-api/workflows/CI/badge.svg)](https://github.com/kudrykv/go-monobank-api/actions?query=workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/kudrykv/go-monobank-api)](https://goreportcard.com/report/github.com/kudrykv/go-monobank-api)
[![codecov](https://codecov.io/gh/kudrykv/go-monobank-api/branch/master/graph/badge.svg)](https://codecov.io/gh/kudrykv/go-monobank-api)

---

go-monobank-api is the library to interact with the Monobank API.
It provides clients for working with public and personal API.

One of its features is no dependencies on 3rd-party libraries.

The library allows to work with webhooks.
It is possible to either parse the webhook response using the helper method, or receive webhooks in channel.

## Usage

### Basic usage

```go
package main

import (
  "context"
  "fmt"

  mono "github.com/kudrykv/go-monobank-api"
)

func main() {
  public := mono.NewPublic()

  currencies, err := public.Currency(context.Background())
  if err != nil {
    panic(err)
  }

  fmt.Println(currencies)

  private := mono.NewPersonal("api-token")

  info, err := private.ClientInfo(context.Background())
  if err != nil {
    panic(err)
  }

  fmt.Println(info)
}
```

### Webhooks

Example app for using channels:
```go
package main

import (
  "context"
  "fmt"
  "net/http"

  mono "github.com/kudrykv/go-monobank-api"
)

func main() {
  personal := mono.NewPersonal("api-token")
  if err := personal.SetWebhook(context.Background(), "https://domain/webhook"); err != nil {
    panic(err)
  }

  webhookChan, handlerFunc := personal.ListenForWebhooks(context.Background())

  mux := http.NewServeMux()
  mux.HandleFunc("/webhook", handlerFunc)

  go func() {
    wh := <-webhookChan

    fmt.Println(wh)
  }()

  if err := http.ListenAndServe("", mux); err != nil {
    panic(err)
  }
}
```

Example app for using helper func:
```go
package main

import (
  "context"
  "fmt"
  "net/http"

  mono "github.com/kudrykv/go-monobank-api"
)

func main() {
  personal := mono.NewPersonal("api-token")
  if err := personal.SetWebhook(context.Background(), "https://domain/webhook"); err != nil {
    panic(err)
  }

  mux := http.NewServeMux()
  mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    webhook, err := personal.ParseWebhook(r.Context(), r.Body)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    
    fmt.Println(webhook)
    w.WriteHeader(http.StatusOK)
  })

  if err := http.ListenAndServe("", mux); err != nil {
    panic(err)
  }
}

```

## Support

Is something missing or works in unexpected way?
[Create for that an issue](https://github.com/kudrykv/go-monobank-api/issues/new).

## Monobank API Documentation
https://api.monobank.ua/docs/

You can obtain your personal token [here](https://api.monobank.ua).

## Progress
- [x] Public API
- [x] Personal API
- [ ] Corporate API (to do)