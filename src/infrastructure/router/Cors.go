package router

type Cors struct {
	allowedOrigins   []string
	allowedMethods   []string
	allowedHeaders   []string
	allowCredentials bool
}

func EmptyCors() *Cors {
	return &Cors{
		allowedOrigins:   make([]string, 0),
		allowedMethods:   make([]string, 0),
		allowedHeaders:   make([]string, 0),
		allowCredentials: false,
	}
}

func (c *Cors) AllowedOrigins(allowedOrigins ...string) *Cors {
	c.allowedOrigins = allowedOrigins
	return c
}

func (c *Cors) AllowedMethods(allowedMethods ...string) *Cors {
	c.allowedMethods = allowedMethods
	return c
}

func (c *Cors) AllowedHeaders(allowedHeaders ...string) *Cors {
	c.allowedHeaders = allowedHeaders
	return c
}

func (c *Cors) AllowCredentials() *Cors {
	c.allowCredentials = true
	return c
}

func (c *Cors) NotAllowCredentials() *Cors {
	c.allowCredentials = false
	return c
}
