package cachedmysql

import "time"

const (
	defaultExpiry = time.Hour * 24 * 7
)

type (
	Options struct {
		Expiry time.Duration
	}

	Option func(o *Options)
)

func newOptions(opts ...Option) Options {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}

	if o.Expiry <= 0 {
		o.Expiry = defaultExpiry
	}

	return o
}

func WithExpiry(expiry time.Duration) Option {
	return func(o *Options) {
		o.Expiry = expiry
	}
}
