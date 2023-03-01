package dyad_test

import (
	"fmt"

	"github.com/dogmatiq/dyad"
)

func ExampleClone() {
	type Value struct {
		Data map[any]any
	}

	src := Value{
		Data: map[any]any{
			"key": "original value",
		},
	}

	dst := dyad.Clone(src)
	dst.Data["key"] = "altered value"

	fmt.Println(src.Data["key"])
	fmt.Println(dst.Data["key"])

	// Output:
	// original value
	// altered value
}
