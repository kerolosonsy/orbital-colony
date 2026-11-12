package main

// Contrast compares one topic between Gin and Fiber in one short sentence per framework.
type Contrast struct {
	Topic string
	Gin   string
	Fiber string
}

var Comparison = []Contrast{
	{Topic: "handler", Gin: "func(*gin.Context) with no return; the engine is an http.Handler.", Fiber: "func(*fiber.Ctx) error on fasthttp; the app is not an http.Handler."},
	{Topic: "errors", Gin: "NoRoute and NoMethod handlers write 404 and 405 after the middleware chain.", Fiber: "Returned errors, including router 404 and 405, reach one ErrorHandler."},
	{Topic: "testing", Gin: "engine.ServeHTTP with httptest.NewRecorder, no port.", Fiber: "app.Test(req, -1) runs a request in memory and returns an *http.Response."},
	{Topic: "lifetime", Gin: "The *gin.Context is pooled; never keep it after the handler returns.", Fiber: "Params and body bytes are reused buffers; strings.Clone before storing."},
}
