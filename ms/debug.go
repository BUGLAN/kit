package ms

type Debug struct {
}

func WithDebug() MicroServiceOption {
	return func(ms *MicroService) {
		ms.debug = true
	}
}
