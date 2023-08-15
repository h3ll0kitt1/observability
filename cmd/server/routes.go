package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *application) setRouters() {

	app.router.Use(app.requestLogger)
	app.router.Use(app.gzipper)

	app.router.Route("/", func(r chi.Router) {
		app.router.Get("/", app.getAll)
		app.router.Get("/ping", app.ping)

		app.router.Route("/value", func(router chi.Router) {
			router.Post("/", app.getValue)
			router.Get("/counter/{name}", app.getCounter)
			router.Get("/gauge/{name}", app.getGauge)
			router.Get("/{other}/{name}", errorUnknown)
		})

		app.router.Post("/updates/", app.updateList)

		app.router.Route("/update", func(router chi.Router) {
			router.Post("/", app.updateValue)

			router.Route("/counter", func(router chi.Router) {
				router.Post("/", errorNoName)
				router.Post("/{name}/{value}", app.updateCounter)
			})

			router.Route("/gauge", func(router chi.Router) {
				router.Post("/", errorNoName)
				router.Post("/{name}/{value}", app.updateGauge)
			})
		})
	})

	app.router.NotFound(errorNotFound)
}
