package router

type Cors struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
}

func EmptyCors() *Cors {
	return &Cors{
		allowedOrigins: make([]string, 0),
		allowedMethods: make([]string, 0),
		allowedHeaders: make([]string, 0),
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
