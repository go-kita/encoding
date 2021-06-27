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

var _marshalers = atomic.NewUnsafePointer(unsafe.Pointer(&map[string]Marshaler{}))

// RegisterMarshaler register a Marshaler with a specific type name. The registered
// Marshaler can be retrieve by the same type name later.
// It more that one Marshaler registered by the same type name, the later one wins.
// This function will ignore the case of the name.
func RegisterMarshaler(name string, marshaler Marshaler) Marshaler {
	if len(name) == 0 || marshaler == nil {
		return nil
	}
	name = strings.ToLower(name)
	for {
		o := (*map[string]Marshaler)(_marshalers.Load())
		om, exist := (*o)[name]
		var mp map[string]Marshaler
		if exist {
			mp = make(map[string]Marshaler, len(*o))
		} else {
			mp = make(map[string]Marshaler, len(*o)+1)
		}
		for n := range *o {
			mp[n] = (*o)[n]
		}
		mp[name] = marshaler
		if _marshalers.CAS(unsafe.Pointer(o), unsafe.Pointer(&mp)) {
			return om
		}
	}
}

// GetMarshaler retrieve the registered Marshaler with type name.
// If no marshaler is registered with the name, nil will be returned.
// This function will ignore the case of the name.
func GetMarshaler(name string) Marshaler {
	return (*(*map[string]Marshaler)(_marshalers.Load()))[name]
}

var _unmarshalers = atomic.NewUnsafePointer(unsafe.Pointer(&map[string]Unmarshaler{}))

// RegisterUnmarshaler register a Unmarshaler with a specific type name. The registered
// Unmarshaler can be retrieve by the same type name later.
// It more that one Unmarshaler registered by the same type name, the later one wins.
// This function will ignore the case of the name.
func RegisterUnmarshaler(name string, unmarshaler Unmarshaler) Unmarshaler {
	if len(name) == 0 || unmarshaler == nil {
		return nil
	}
	name = strings.ToLower(name)
	for {
		o := (*map[string]Unmarshaler)(_unmarshalers.Load())
		ou, exist := (*o)[name]
		var mp map[string]Unmarshaler
		if exist {
			mp = make(map[string]Unmarshaler, len(*o))
		} else {
			mp = make(map[string]Unmarshaler, len(*o)+1)
		}
		for n := range *o {
			mp[n] = (*o)[n]
		}
		mp[name] = unmarshaler
		if _unmarshalers.CAS(unsafe.Pointer(o), unsafe.Pointer(&mp)) {
			return ou
		}
	}
}

// GetUnmarshaler retrieve the registered Unmarshaler with type name.
// If no unmarshaler is registered with the name, nil will be returned.
// This function will ignore the case of the name.
func GetUnmarshaler(name string) Unmarshaler {
	return (*(*map[string]Unmarshaler)(_unmarshalers.Load()))[name]
}
