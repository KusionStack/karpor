package config

type Controller struct {
	config *Config
}

func NewController(config *Config) *Controller {
	return &Controller{
		config: config,
	}
}

func (c *Controller) Get() *Config {
	return c.config
}
