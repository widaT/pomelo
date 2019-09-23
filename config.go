package pomelo

type Config struct {
	EnableGizp     bool
	Address        string
	ParseMultiForm bool
	ErrLog         string
	AccLog         string
	LogMaxSize     int64
	LogMaxFiles    int
}

type Option func(config *Config)

func NewConfig(opts ...Option) *Config {
	config := &Config{
		Address:        "0.0.0.0:8080",
		ParseMultiForm: true,
		LogMaxSize:     1 << 30,
		LogMaxFiles:    7,
	}
	for _, o := range opts {
		o(config)
	}
	return config
}

func EnableGzip(sw bool) Option {
	return func(config *Config) {
		config.EnableGizp = sw
	}
}

func Address(ip string) Option {
	return func(config *Config) {
		config.Address = ip
	}
}

func ParseMultiForm(b bool) Option {
	return func(config *Config) {
		config.ParseMultiForm = b
	}
}

func ELog(path string) Option {
	return func(config *Config) {
		config.ErrLog = path
	}
}

func ALog(path string) Option {
	return func(config *Config) {
		config.AccLog = path
	}
}

func LogMaxSize(size int64) Option {
	return func(config *Config) {
		config.LogMaxSize = size
	}
}

func LogMaxFiles(num int) Option {
	return func(config *Config) {
		config.LogMaxFiles = num
	}
}
