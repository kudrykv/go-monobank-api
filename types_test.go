package mono_test

import (
	"testing"
	"time"

	mono "github.com/kudrykv/go-monobank-api"
)

func TestTime_Time(t *testing.T) {
	expectEquals(t, mono.Time(123).Time(), time.Unix(123, 0))
}

func TestTime_UnmarshalJSON(t *testing.T) {
	tt := mono.Time(0)
	expectNoError(t, tt.UnmarshalJSON([]byte("null")))
	expectDeepEquals(t, tt, mono.Time(0))

	expectNoError(t, tt.UnmarshalJSON([]byte("1554466347")))
	expectDeepEquals(t, tt, mono.Time(1554466347))

	expectError(t, tt.UnmarshalJSON([]byte("1554bubu47")), "strconv.ParseInt: parsing \"1554bubu47\": invalid syntax")
}
