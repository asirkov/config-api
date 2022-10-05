package log

type MultiLogAdapter struct {
	Logger

	loggers []Logger
}

func NewMultiLogAdapter(l ...Logger) Logger {
	return &MultiLogAdapter{
		loggers: l,
	}
}

func (a *MultiLogAdapter) Log(entry Entry) {
	if a.loggers != nil {
		for _, logger := range a.loggers {
			logger.Log(entry)
		}
	}
}
