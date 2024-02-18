package log

var logger = New() //nolint:gochecknoglobals

func Trace(message any) {
	logger.Trace(message)
}

func Debug(message any) {
	logger.Debug(message)
}

func Info(message any) {
	logger.Info(message)
}

func Warn(message any) {
	logger.Info(message)
}

func Error(message any) {
	logger.Error(message)
}

func SetLevel(level ConfigLevel) {
	logger.SetLevel(level)
}

func SetFormat(format Format) {
	logger.SetFormat(format)
}
