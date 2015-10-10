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

func TestParentAndChildren(t *testing.T) {
	parent := NewContext()
	child1 := NewContext()
	child2 := NewContext()
	child3 := NewContext()

	child1.SetParent(parent)
	child2.SetParent(parent)
	child3.SetParent(child2)

	if child1.Parent != parent {
		t.Error("child1 should be a child of parent")
	}

	if len(parent.Children) != 2 {
		t.Error("parent should have two children")
	}

	if child3.Parent != child2 {
		t.Error("child3 should be a child of child2")
	}

	if len(child2.Children) != 1 {
		t.Error("child2 should have exactly one child")
	}
}
