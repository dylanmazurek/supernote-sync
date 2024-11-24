package supernotelocal

type options struct {
	host string
	port int
}

func defaultOptions() options {
	defaultOptions := options{}

	return defaultOptions
}

type Option func(*options)

func WithHost(host string) Option {
	return func(o *options) {
		o.host = host
	}
}

func WithPort(port int) Option {
	return func(o *options) {
		o.port = port
	}
}
