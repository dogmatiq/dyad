package dyad

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dogmatiq/dyad/internal/unsafereflect"
)

// Clone returns a deep copy of src.
func Clone[T any](src T, options ...Option) (dst T) {
	dst, err := clone(src, options)
	if err != nil {
		panic(err)
	}

	return dst
}

// An Option changes the behavior of a clone operation.
//
// The signature of this function is not part of the public API and may change
// at any time without warning.
type Option func(*cloneOptions)

type cloneOptions struct {
	channelStrategy         ChannelStrategy
	unexportedFieldStrategy UnexportedFieldStrategy
}

func clone[T any](src T, options []Option) (dst T, err error) {
	var opts cloneOptions

	for _, o := range options {
		o(&opts)
	}

	err = cloneInto(
		reflect.ValueOf(&src).Elem(),
		reflect.ValueOf(&dst).Elem(),
		opts,
	)

	return dst, err
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

var (
	timeType = typeOf[time.Time]()
)

func cloneInto(src, dst reflect.Value, opts cloneOptions) error {
	if !src.IsValid() {
		return nil
	}

	switch src.Type() {
	case timeType:
		dst.Set(src)
		return nil
	}

	switch src.Kind() {
	case reflect.Interface:
		return cloneInterfaceInto(src, dst, opts)
	case reflect.Ptr:
		return clonePtrInto(src, dst, opts)
	case reflect.Slice:
		return cloneSliceInto(src, dst, opts)
	case reflect.Map:
		return cloneMapInto(src, dst, opts)
	case reflect.Struct:
		return cloneStructInto(src, dst, opts)
	case reflect.Chan:
		return cloneChannelInto(src, dst, opts)
	default:
		dst.Set(src)
		return nil
	}
}

func cloneInterfaceInto(src, dst reflect.Value, opts cloneOptions) error {
	if src.IsNil() {
		return nil
	}

	srcElem := src.Elem()
	dstElem := reflect.New(srcElem.Type()).Elem()

	if err := cloneInto(srcElem, dstElem, opts); err != nil {
		return err
	}

	dst.Set(dstElem)

	return nil
}

func clonePtrInto(src, dst reflect.Value, opts cloneOptions) error {
	if src.IsNil() {
		return nil
	}

	srcElem := src.Elem()
	dstPtr := reflect.New(srcElem.Type())
	dstElem := dstPtr.Elem()

	if err := cloneInto(srcElem, dstElem, opts); err != nil {
		return err
	}

	dst.Set(dstPtr)

	return nil
}

func cloneSliceInto(src, dst reflect.Value, opts cloneOptions) error {
	if src.IsNil() {
		return nil
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
		if err := cloneInto(
			src.Index(i),
			dst.Index(i),
			opts,
		); err != nil {
			return err
		}
	}

	return nil
}

func cloneMapInto(src, dst reflect.Value, opts cloneOptions) error {
	if src.IsNil() {
		return nil
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
		if err := cloneInto(srcKey, dstKey, opts); err != nil {
			return err
		}

		dstElem := reflect.New(elemType).Elem()
		if err := cloneInto(srcElem, dstElem, opts); err != nil {
			return err
		}

		dst.SetMapIndex(dstKey, dstElem)
	}

	return nil
}

func cloneStructInto(src, dst reflect.Value, opts cloneOptions) error {
	size := src.NumField()
	srcType := src.Type()

	for i := 0; i < size; i++ {
		field := srcType.Field(i)
		srcField := src.Field(i)
		dstField := dst.Field(i)

		// If the field is unexported
		if field.PkgPath != "" {
			switch opts.unexportedFieldStrategy {
			case CloneUnexportedFields:
				srcField = unsafereflect.MakeMutable(srcField)
				dstField = unsafereflect.MakeMutable(dstField)
			case IgnoreUnexportedFields:
				continue
			default:
				return fmt.Errorf(
					"cannot clone %s.%s, try the dyad.WithUnexportedFieldStrategy() option",
					srcType,
					field.Name,
				)
			}
		}

		if err := cloneInto(srcField, dstField, opts); err != nil {
			return err
		}
	}

	return nil
}

func cloneChannelInto(src, dst reflect.Value, opts cloneOptions) error {
	switch opts.channelStrategy {
	case ShareChannels:
		dst.Set(src)
	case IgnoreChannels:
	default:
		return fmt.Errorf(
			"cannot clone %s, try the dyad.WithChannelStrategy() option",
			src.Type(),
		)
	}

	return nil
}
