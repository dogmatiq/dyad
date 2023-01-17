package dyad

import (
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

func clone[T any](src T, options []Option) (dst T, err error) {
	var ctx cloneContext

	for _, o := range options {
		o(&ctx.options)
	}

	srcV := reflect.ValueOf(&src).Elem()
	dstV := reflect.ValueOf(&dst).Elem()

	err = cloneInto(
		ctx.WithPath("%s", srcV.Type()),
		srcV,
		dstV,
	)

	return dst, err
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

var (
	timeType = typeOf[time.Time]()
)

func cloneInto(
	ctx cloneContext,
	src, dst reflect.Value,
) error {
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
		return cloneInterfaceInto(ctx, src, dst)
	case reflect.Ptr:
		return clonePtrInto(ctx, src, dst)
	case reflect.Slice:
		return cloneSliceInto(ctx, src, dst)
	case reflect.Map:
		return cloneMapInto(ctx, src, dst)
	case reflect.Struct:
		return cloneStructInto(ctx, src, dst)
	case reflect.Chan:
		return cloneChannelInto(ctx, src, dst)
	default:
		dst.Set(src)
		return nil
	}
}

func cloneInterfaceInto(
	ctx cloneContext,
	src, dst reflect.Value,
) error {
	if src.IsNil() {
		return nil
	}

	srcElem := src.Elem()
	dstElem := reflect.New(srcElem.Type()).Elem()

	if err := cloneInto(
		ctx.WithPath("(%s)", srcElem.Type()),
		srcElem,
		dstElem,
	); err != nil {
		return err
	}

	dst.Set(dstElem)

	return nil
}

func clonePtrInto(
	ctx cloneContext,
	src, dst reflect.Value,
) error {
	if src.IsNil() {
		return nil
	}

	srcElem := src.Elem()
	dstPtr := reflect.New(srcElem.Type())
	dstElem := dstPtr.Elem()

	if err := cloneInto(ctx, srcElem, dstElem); err != nil {
		return err
	}

	dst.Set(dstPtr)

	return nil
}

func cloneSliceInto(
	ctx cloneContext,
	src, dst reflect.Value,
) error {
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
			ctx.WithPath("[%d]", i),
			src.Index(i),
			dst.Index(i),
		); err != nil {
			return err
		}
	}

	return nil
}

func cloneMapInto(
	ctx cloneContext,
	src, dst reflect.Value,
) error {
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
		ctx := ctx.WithPath("[%#v]", srcKey.Interface())
		srcElem := src.MapIndex(srcKey)

		dstKey := reflect.New(keyType).Elem()
		if err := cloneInto(ctx, srcKey, dstKey); err != nil {
			return err
		}

		dstElem := reflect.New(elemType).Elem()
		if err := cloneInto(ctx, srcElem, dstElem); err != nil {
			return err
		}

		dst.SetMapIndex(dstKey, dstElem)
	}

	return nil
}

func cloneStructInto(
	ctx cloneContext,
	src, dst reflect.Value,
) error {
	size := src.NumField()
	srcType := src.Type()

	for i := 0; i < size; i++ {
		field := srcType.Field(i)
		srcField := src.Field(i)
		dstField := dst.Field(i)

		// If the field is unexported
		if field.PkgPath != "" {
			switch ctx.options.unexportedFieldStrategy {
			case CloneUnexportedFields:
				srcField = unsafereflect.MakeMutable(srcField)
				dstField = unsafereflect.MakeMutable(dstField)
			case IgnoreUnexportedFields:
				continue
			default:
				return ctx.Error(
					"struct cannot be cloned due to unexported field (%s.%s), try the dyad.WithUnexportedFieldStrategy() option",
					srcType,
					field.Name,
				)
			}
		}

		if err := cloneInto(
			ctx.WithPath(".%s", field.Name),
			srcField,
			dstField,
		); err != nil {
			return err
		}
	}

	return nil
}

func cloneChannelInto(
	ctx cloneContext,
	src, dst reflect.Value,
) error {
	switch ctx.options.channelStrategy {
	case ShareChannels:
		dst.Set(src)
	case IgnoreChannels:
	default:
		return ctx.Error("channels cannot be cloned, try the dyad.WithChannelStrategy() option")
	}

	return nil
}
