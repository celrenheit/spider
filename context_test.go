package spider

import (
	"net/http"
	"testing"
)

func TestContextStore(t *testing.T) {
	ctx := NewContext()

	req := &http.Request{}
	ctx.Set("req", req)

	val := ctx.Get("req")
	switch val.(type) {
	case *http.Request:
	default:
		t.Error("Not of type string")
	}

	if val != req {
		t.Error("Not equal")
	}
}
