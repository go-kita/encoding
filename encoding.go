package encoding

import (
	"context"
	"strings"
	"unsafe"

	"go.uber.org/atomic"
)

// Marshaler can encode a value of supported type into binary data.
type Marshaler interface {
	// Marshal encodes a value of supported type into binary data.
	// If the type of the value is not supported, an error will be returned.
	// Any error occurred during the encoding process would be returned.
	// When the returned error is not nil, the content of the result is not
	// guaranteed.
	//
	// If an encoding can be extract from the context, it would be used to
	// encode the final binary data.
	Marshal(ctx context.Context, v interface{}) ([]byte, error)
}

// Unmarshaler can decode value from the binary data to supported type.
type Unmarshaler interface {
	// Unmarshal decodes value from the binary data.
	// If the type of the target value is not supported, an error will be returned.
	// Any error occurred during the decoding process would be returned.
	// When the returned error is not nil, the state of the target value is not
	// guaranteed.
	//
	// If an encoding can be extract from the context, it would be used to
	// decode the input binary data before decoding to the value.
	//
	// Unmarshaler must copy the data if it wishes to retain the data
	// after returning.
	Unmarshal(ctx context.Context, data []byte, v interface{}) error
}

// MarshalerSupplier is a function which supplies Marshalers.
type MarshalerSupplier func() Marshaler

// UnmarshalerSupplier is a function which supplies Unmarshalers.
type UnmarshalerSupplier func() Unmarshaler

var _marshalerSuppliers = atomic.NewUnsafePointer(unsafe.Pointer(&map[string]MarshalerSupplier{}))

// RegisterMarshaler register a MarshalerSupplier with a specific type name.
// It more that one MarshalerSupplier registered by the same type name, the later one wins.
// This function will ignore the case of the name.
func RegisterMarshaler(name string, supplier MarshalerSupplier) MarshalerSupplier {
	if len(name) == 0 || supplier == nil {
		return nil
	}
	name = strings.ToLower(name)
	for {
		o := (*map[string]MarshalerSupplier)(_marshalerSuppliers.Load())
		om, exist := (*o)[name]
		var mp map[string]MarshalerSupplier
		if exist {
			mp = make(map[string]MarshalerSupplier, len(*o))
		} else {
			mp = make(map[string]MarshalerSupplier, len(*o)+1)
		}
		for n := range *o {
			mp[n] = (*o)[n]
		}
		mp[name] = supplier
		if _marshalerSuppliers.CAS(unsafe.Pointer(o), unsafe.Pointer(&mp)) {
			return om
		}
	}
}

// GetMarshaler retrieve a Marshaler by type name.
// If no MarshalerSupplier is registered with the name, nil will be returned.
// This function will ignore the case of the name.
func GetMarshaler(name string) Marshaler {
	name = strings.ToLower(name)
	if supplier := (*(*map[string]MarshalerSupplier)(_marshalerSuppliers.Load()))[name]; supplier != nil {
		return supplier()
	}
	return nil
}

var _unmarshalerSuppliers = atomic.NewUnsafePointer(unsafe.Pointer(&map[string]UnmarshalerSupplier{}))

// RegisterUnmarshaler register a UnmarshalerSupplier with a specific type name.
// It more that one UnmarshalerSupplier registered by the same type name, the later one wins.
// This function will ignore the case of the name.
func RegisterUnmarshaler(name string, supplier UnmarshalerSupplier) UnmarshalerSupplier {
	if len(name) == 0 || supplier == nil {
		return nil
	}
	name = strings.ToLower(name)
	for {
		o := (*map[string]UnmarshalerSupplier)(_unmarshalerSuppliers.Load())
		ou, exist := (*o)[name]
		var mp map[string]UnmarshalerSupplier
		if exist {
			mp = make(map[string]UnmarshalerSupplier, len(*o))
		} else {
			mp = make(map[string]UnmarshalerSupplier, len(*o)+1)
		}
		for n := range *o {
			mp[n] = (*o)[n]
		}
		mp[name] = supplier
		if _unmarshalerSuppliers.CAS(unsafe.Pointer(o), unsafe.Pointer(&mp)) {
			return ou
		}
	}
}

// GetUnmarshaler retrieve an Unmarshaler by type name.
// If no UnmarshalerSupplier is registered with the name, nil will be returned.
// This function will ignore the case of the name.
func GetUnmarshaler(name string) Unmarshaler {
	name = strings.ToLower(name)
	if supplier := (*(*map[string]UnmarshalerSupplier)(_unmarshalerSuppliers.Load()))[name]; supplier != nil {
		return supplier()
	}
	return nil
}
