package route

import (
	http "be/neurade/v2/internal/http/controller"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type RouteConfig struct {
	App                     *chi.Mux
	UserController          *http.UserController
	LLMController           *http.LLMController
	CourseController        *http.CourseController
	AssignmentController    *http.AssignmentController
	PrController            *http.PrController
	GitHubWebhookController *http.GitHubWebhookController
}

func (c *RouteConfig) Setup() *chi.Mux {
	r := c.App

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.UserController.Register)
		r.Post("/login", c.UserController.Login)
	})

	r.Route("/llm", func(r chi.Router) {
		r.Post("/create", c.LLMController.Create)
		r.Get("/{id}", c.LLMController.GetById)
	})

	r.Route("/course", func(r chi.Router) {
		r.Post("/create", c.CourseController.Create)
		r.Get("/owner/{user_id}", c.CourseController.GetAllByOwner)

		r.Route("/{course_id}", func(r chi.Router) {
			r.Get("/", c.CourseController.GetByID)
			r.Route("/pr", func(r chi.Router) {
				r.Get("/", c.PrController.GetAllByCourse)
				r.Get("/{pr_id}", c.PrController.GetByID)
				r.Post("/create", c.PrController.Create)
			})
			r.Route("/assignment", func(r chi.Router) {
				r.Get("/", c.AssignmentController.GetAllByCourse)
				// r.Get("/{assignment_id}", c.AssignmentController.GetByID)
				r.Post("/create", c.AssignmentController.Create)
			})
		})
		r.Get("/assignment/{assignment_id}", c.AssignmentController.GetByID)
	})

	r.Route("/webhook", func(r chi.Router) {
		r.Post("/fetch-pull-requests", c.GitHubWebhookController.FetchPullRequests)
		r.Get("/course/{course_id}/pull-requests", c.GitHubWebhookController.GetPullRequestsByCourse)
	})

	return r
}
