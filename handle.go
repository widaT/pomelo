package pomelo

// for 404 page
func NotFound(c *Context) {
	c.NotFound("The requested URL " + c.Request.URL.Path + " was not found on this server")
}
