package main

func (app *application) setRouters() {
	app.router.Use(app.requestLogger)
	app.router.Post("/value/", app.getValue)
	app.router.Post("/value", app.getValue)
	app.router.Post("/update/", app.updateValue)
	app.router.Post("/update", app.updateValue)
	app.router.NotFound(errorNotFound)
}
