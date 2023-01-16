package dyad_test

import (
	. "github.com/dogmatiq/dyad"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Clone()", func() {
	When("the source value is an interface", func() {
		It("copies the value encapsulated by the interface", func() {
			value := "<value>"

			src := any(&value)
			dst := Clone(src)

			Expect(dst).To(Equal(src))

			value = "<changed>"
			Expect(dst).ToNot(Equal(src))
		})

		It("handles nil values", func() {
			var src any
			dst := Clone(src)

			Expect(dst).To(BeNil())
		})
	})

	When("the source value is a pointer", func() {
		It("copies the pointed-to-value", func() {
			value := "<value>"

			src := &value
			dst := Clone(src)

			Expect(*dst).To(Equal(*src))

			value = "<changed>"
			Expect(*dst).ToNot(Equal(*src))
		})

		It("handles nil values", func() {
			var src *int
			dst := Clone(src)

			Expect(dst).To(BeNil())
		})
	})

	When("the source value is a slice", func() {
		It("copies the slice itself", func() {
			src := []int{1, 2, 3}
			dst := Clone(src)

			Expect(dst).To(Equal(src))

			src[0] = 123
			Expect(dst).ToNot(Equal(src))
		})

		It("copies the elements within the slice", func() {
			original := "<value>"

			src := []*string{&original}
			dst := Clone(src)

			Expect(dst).To(Equal(src))

			original = "<changed>"
			Expect(dst).ToNot(Equal(src))
		})

		It("handles nil values", func() {
			var src []int
			dst := Clone(src)

			Expect(dst).To(BeNil())
		})
	})

	When("the source value is a map", func() {
		It("copies the map itself", func() {
			src := map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			}
			dst := Clone(src)

			Expect(dst).To(Equal(src))

			delete(src, "one")
			Expect(dst).ToNot(Equal(src))
		})

		It("copies the keys within the map", func() {
			key := "<key>"

			src := map[*string]int{&key: 123}
			dst := Clone(src)

			Expect(dst).To(HaveLen(1))

			for k := range dst {
				Expect(*k).To(Equal(key))
				Expect(k).ToNot(BeIdenticalTo(&key))
			}
		})

		It("copies the elements within the map", func() {
			elem := "<value>"

			src := map[string]*string{"<key>": &elem}
			dst := Clone(src)

			Expect(dst).To(Equal(src))

			elem = "<changed>"
			Expect(dst).ToNot(Equal(src))
		})

		It("handles nil values", func() {
			var src map[string]int
			dst := Clone(src)

			Expect(dst).To(BeNil())
		})
	})

	When("the source value is a struct", func() {
		It("copies the struct itself", func() {
			type Source struct {
				Value string
			}

			src := Source{"<value>"}
			dst := Clone(src)

			Expect(dst).To(Equal(src))
		})

		It("copies unexported fields", func() {
			type Source struct {
				value string
			}

			src := Source{"<value>"}
			dst := Clone(src)

			Expect(dst).To(Equal(src))
		})

		It("copies embedded fields", func() {
			type Embedded struct {
				Value string
			}

			type Source struct {
				Embedded
			}

			src := Source{Embedded{"<value>"}}
			dst := Clone(src)

			Expect(dst).To(Equal(src))
		})

		It("copies the field values within the struct", func() {
			original := "<value>"

			type Source struct {
				Ptr *string
			}

			src := Source{&original}
			dst := Clone(src)

			Expect(dst).To(Equal(src))

			original = "<changed>"
			Expect(dst).ToNot(Equal(src))
		})

		It("copies interface values within the struct", func() {
			original := "<value>"

			type Source struct {
				Value any
			}

			src := Source{&original}
			dst := Clone(src)

			Expect(dst).To(Equal(src))

			original = "<changed>"
			Expect(dst).ToNot(Equal(src))
		})

		It("copies a struct within a nil-interface field", func() {
			type Source struct {
				Value any
			}

			src := Source{}
			dst := Clone(src)

			Expect(dst).To(Equal(src))
		})
	})

	When("the source value is a channel", func() {
		It("panics", func() {
			Expect(func() {
				src := make(chan int, 1)
				Clone(src)
			}).To(PanicWith(MatchError(
				"cannot clone type: chan int",
			)))
		})
	})

	When("the source value is a basic type", func() {
		It("returns the same value", func() {
			Expect(Clone(true)).To(BeTrue())

			Expect(Clone(uintptr(123))).To(Equal(uintptr(123)))

			Expect(Clone("<string>")).To(Equal("<string>"))

			Expect(Clone(int(123))).To(Equal(int(123)))
			Expect(Clone(int8(123))).To(Equal(int8(123)))
			Expect(Clone(int16(123))).To(Equal(int16(123)))
			Expect(Clone(int32(123))).To(Equal(int32(123)))
			Expect(Clone(int64(123))).To(Equal(int64(123)))

			Expect(Clone(uint(123))).To(Equal(uint(123)))
			Expect(Clone(uint8(123))).To(Equal(uint8(123)))
			Expect(Clone(uint16(123))).To(Equal(uint16(123)))
			Expect(Clone(uint32(123))).To(Equal(uint32(123)))
			Expect(Clone(uint64(123))).To(Equal(uint64(123)))

			Expect(Clone(float32(123.45))).To(Equal(float32(123.45)))
			Expect(Clone(float64(123.45))).To(Equal(float64(123.45)))

			Expect(Clone(complex64(123))).To(Equal(complex64(123)))
			Expect(Clone(complex128(123))).To(Equal(complex128(123)))
		})
	})
})