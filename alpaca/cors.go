package alpaca

type Cors struct {
	allowAllOrigins  []string
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
	allowCredentials bool
	maxAge           int
}

func newCors(c *Cors) *Cors {

	if c.allowAllOrigins == nil {
		c.allowAllOrigins = []string{"*"}
	}

	if c.allowMethods == nil {
		c.allowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}

	if c.allowHeaders == nil {
		c.allowHeaders = []string{"Content-Type", "Authorization"}
	}

	if c.exposeHeaders == nil {
		c.exposeHeaders = nil
	}

	if c.maxAge == 0 {
		c.maxAge = 24
	}

	return &Cors{
		allowAllOrigins:  c.allowAllOrigins,
		allowMethods:     c.allowMethods,
		allowHeaders:     c.allowHeaders,
		exposeHeaders:    c.exposeHeaders,
		allowCredentials: c.allowCredentials,
		maxAge:           c.maxAge * 3600,
	}
}
