package render

import (
	"errors"
	"net/http"
	"reflect"
)

// Binder interface for managing request payloads.
type Binder interface {
	Bind(r *http.Request) error
}

// Bind decodes a request body and executes the Binder interface of the destination.
func Bind(r *http.Request, v Binder) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("render: destination must be a non-nil pointer")
	}

	// Create a new temporary instance of the underlying type
	temp := reflect.New(val.Elem().Type()).Interface().(Binder)

	if err := Decode(r, temp); err != nil {
		return err
	}
	if err := temp.Bind(r); err != nil {
		return err
	}

	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(temp).Elem())
	return nil
}
