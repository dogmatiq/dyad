package dyad

// An Option changes the behavior of a clone operation.
//
// The signature of this function is not part of the public API and may change
// at any time without warning.
type Option func(*cloneOptions)

type cloneOptions struct {
	channelStrategy         ChannelStrategy
	unexportedFieldStrategy UnexportedFieldStrategy
}

// ChannelStrategy is an enumeration of strategies that can be used by Clone()
// when a channel is encountered.
type ChannelStrategy int

const (
	// PanicOnChannel causes Clone() to panic when it encounters a
	// channel.
	//
	// This is the default behavior.
	PanicOnChannel ChannelStrategy = iota

	// ShareChannels causes Clone() to share the same channel between the
	// original and cloned values.
	ShareChannels

	// IgnoreChannels causes Clone() to use a nil value when a channel is
	// encountered.
	IgnoreChannels
)

// WithChannelStrategy is an option that controls how Clone() behaves when it
// encounters a channel.
func WithChannelStrategy(s ChannelStrategy) Option {
	return func(opts *cloneOptions) {
		opts.channelStrategy = s
	}
}

// UnexportedFieldStrategy is an enumeration of strategies that can be used by
// Clone() when an unexported struct field is encountered.
type UnexportedFieldStrategy int

const (
	// PanicOnUnexportedField causes Clone() to panic when it encounters an
	// unexported struct field.
	//
	// This is the default behavior.
	PanicOnUnexportedField UnexportedFieldStrategy = iota

	// CloneUnexportedFields causes Clone() to clone unexported struct fields
	// just as it would any other value.
	CloneUnexportedFields

	// IgnoreUnexportedFields causes Clone() to skip cloning of unexported
	// struct fields, leaving them as their zero value.
	IgnoreUnexportedFields
)

// WithUnexportedFieldStrategy is an option that controls how Clone() behaves
// when it encounters an unexported field.
func WithUnexportedFieldStrategy(s UnexportedFieldStrategy) Option {
	return func(opts *cloneOptions) {
		opts.unexportedFieldStrategy = s
	}
}
