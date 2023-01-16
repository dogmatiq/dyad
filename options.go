package dyad

// ChannelStrategy is an enumeration of strategies that can be used by Clone()
// when a channel is encountered.
type ChannelStrategy int

const (
	// PanicOnChannel causes Clone() to panic when it encounters a
	// channel.
	//
	// This is the default behavior.
	PanicOnChannel ChannelStrategy = iota

	// ShareChannel causes Clone() to share the same channel between the
	// original and cloned values.
	ShareChannel
)

// WithChannelStrategy is an option that controls how Clone() behaves when it
// encounters a channel.
func WithChannelStrategy(s ChannelStrategy) Option {
	return func(opts *cloneOptions) {
		opts.channelStrategy = s
	}
}
