package misc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.False(t, misc.IsNil(value), "IsNil() on value was true.")
	assert.False(t, misc.IsNil(pointer), "IsNil() on pointer was true.")
	assert.True(t, misc.IsNil(nilPointer), "IsNil() on nil pointer was false.")
	assert.False(t, misc.IsNil(typedInterfaceValue), "IsNil() on typed interface value was true.")
	assert.False(t, misc.IsNil(typedInterfacePointer), "IsNil() on typed interface pointer was true.")
	assert.True(t, misc.IsNil(typedNilInterface), "IsNil() on nil interface was false.")
	assert.True(t, misc.IsNil(untypedNilInterface), "IsNil() on untyped nil interface was false.")
}
