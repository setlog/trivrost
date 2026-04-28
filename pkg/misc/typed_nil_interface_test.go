package misc_test

import (
	"testing"

	"github.com/setlog/trivrost/pkg/misc"
)

func TestIsNil(t *testing.T) {
	var value int32 = 42
	var pointer *int32 = &value
	var nilPointer *int32 = nil
	var typedInterfaceValue interface{} = value
	var typedInterfacePointer interface{} = pointer
	var typedNilInterface interface{} = nilPointer
	var untypedNilInterface interface{} = nil
	if misc.IsNil(value) {
		t.Fatal("IsNil() on value was true.")
	}
	if misc.IsNil(pointer) {
		t.Fatal("IsNil() on pointer was true.")
	}
	if !misc.IsNil(nilPointer) {
		t.Fatal("IsNil() on nil pointer was false.")
	}
	if misc.IsNil(typedInterfaceValue) {
		t.Fatal("IsNil() on typed interface value was true.")
	}
	if misc.IsNil(typedInterfacePointer) {
		t.Fatal("IsNil() on typed interface pointer was true.")
	}
	if !misc.IsNil(typedNilInterface) {
		t.Fatal("IsNil() on nil interface was false.")
	}
	if !misc.IsNil(untypedNilInterface) {
		t.Fatal("IsNil() on untyped nil interface was false.")
	}
}
