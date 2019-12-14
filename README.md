# Bindings for Monobank API

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
}
```
