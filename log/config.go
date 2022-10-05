package log

type Level string

const (
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
)

type Config struct {
	LogLevel            Level  `yaml:"logLevel"`
	ConsoleEnabled      bool   `yaml:"consoleEnabled"`
	LogFilesEnabled     bool   `yaml:"logFilesEnabled"`
	LogFilesDestination string `yaml:"logFilesDestination"`
	MaxSize             int    `yaml:"maxSize"`
	MaxAge              int    `yaml:"maxAge"`
	MaxBackups          int    `yaml:"maxBackups"`
}
