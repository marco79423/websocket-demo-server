package config

type IConfig interface {
	GetName() string     // 名稱
	GetLogLevel() string // Log 層級

	GetAddress() string
}

func NewConfig() (IConfig, error) {
	return &config{}, nil
}

type config struct {
}

func (c *config) GetName() string {
	return "Websocket Demo Server"
}

func (c *config) GetLogLevel() string {
	return "debug"
}

func (c *config) GetAddress() string {
	return ":8080"
}
