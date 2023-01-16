package dyad

import (
	"reflect"

	"github.com/dogmatiq/dyad/internal/unsafereflect"
)

// Clone returns a deep copy of src.
func Clone[T any](src T, options ...Option) (dst T) {
	clone(
		reflect.ValueOf(src),
		reflect.ValueOf(&dst).Elem(),
	)

	return dst
}

// An Option changes the behavior of a clone operation.
type Option interface {
	future()
}

func clone(src, dst reflect.Value) {
	if !src.IsValid() {
		return
	}

	switch src.Kind() {
	case reflect.Interface:
		cloneInterface(src, dst)
	case reflect.Ptr:
		clonePtr(src, dst)
	case reflect.Slice:
		cloneSlice(src, dst)
	case reflect.Map:
		cloneMap(src, dst)
	case reflect.Struct:
		cloneStruct(src, dst)
	default:
		dst.Set(src)
	}
}

func cloneInterface(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}

	srcElem := src.Elem()
	dstElem := reflect.New(srcElem.Type()).Elem()

	clone(srcElem, dstElem)
	dst.Set(dstElem)
}

func clonePtr(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}

	srcElem := src.Elem()
	dstPtr := reflect.New(srcElem.Type())
	dstElem := dstPtr.Elem()

	clone(srcElem, dstElem)
	dst.Set(dstPtr)
}

func cloneSlice(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}

	size := src.Len()

	dst.Set(
		reflect.MakeSlice(
			src.Type(),
			size,
			src.Cap(),
		),
	)

	for i := 0; i < size; i++ {
		clone(
			src.Index(i),
			dst.Index(i),
		)
	}
}

func cloneMap(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}

	mapType := src.Type()
	keyType := mapType.Key()
	elemType := mapType.Elem()

	dst.Set(
		reflect.MakeMap(mapType),
	)

	for _, srcKey := range src.MapKeys() {
		srcElem := src.MapIndex(srcKey)
		dstKey := reflect.New(keyType).Elem()
		dstElem := reflect.New(elemType).Elem()

		clone(srcKey, dstKey)
		clone(srcElem, dstElem)

		dst.SetMapIndex(dstKey, dstElem)
	}
}

func cloneStruct(src, dst reflect.Value) {
	size := src.NumField()
	srcType := src.Type()

	for i := 0; i < size; i++ {
		field := srcType.Field(i)
		srcField := src.Field(i)
		dstField := dst.Field(i)

		// If the field is unexported
		if field.PkgPath != "" {
			srcField = unsafereflect.MakeMutable(srcField)
			dstField = unsafereflect.MakeMutable(dstField)
		}

		clone(srcField, dstField)
	}
}
