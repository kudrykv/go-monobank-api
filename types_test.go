package mono_test

import (
	"testing"
	"time"

	mono "github.com/kudrykv/go-monobank-api"
)

func TestTime_Time(t *testing.T) {
	expectEquals(t, mono.Time(123).Time(), time.Unix(123, 0))
}
