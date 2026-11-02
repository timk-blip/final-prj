package api

import (
	"github.com/go-chi/chi/v5"
)

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", nextdateHandler)
	r.Post("/api/signin", signinHandler)

	r.Group(func(gr chi.Router) {
		gr.Use(auth)
		gr.Post("/api/task", addTaskHandler)
		gr.Get("/api/task", getTaskHandler)
		gr.Put("/api/task", putTaskHandler)
		gr.Delete("/api/task", deleteTaskHandler)
		gr.Get("/api/tasks", tasksHandler)
		gr.Post("/api/task/done", doneTaskHandler)
	})
}
