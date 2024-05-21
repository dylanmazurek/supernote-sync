package supernote

type options struct {
	username string
	password string
}

func defaultOptions() options {
	defaultOptions := options{}

	return defaultOptions
}

type Option func(*options)

func WithCredentials(username string, password string) Option {
	return func(o *options) {
		o.username = username
		o.password = password
	}
}
