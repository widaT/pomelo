package pomelo

type Config struct {
	EnableGizp     bool
	Address        string
	ParseMultiForm bool //是否预先解析 MultipartForm

	LogPath string
}

type Option func(config *Config)

func NewConfig(opts ...Option) *Config {
	config := &Config{
		Address: "0.0.0.0:8080",
		LogPath: "logs/",
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

func LogPath(path string) Option {
	return func(config *Config) {
		config.LogPath = path
	}
}
