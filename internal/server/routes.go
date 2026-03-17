package server

import (
	_ "embed"
	"net/http"

	"github.com/throskam/ki"
	"github.com/throskam/kix/auth"
	"github.com/throskam/kix/sess"
	"github.com/throskam/memo/internal/lib"
	"github.com/throskam/memo/internal/server/handlers"
	"github.com/throskam/memo/web"
)

func addRoutes(
	r ki.Router,
	store sess.SessionStore,
	ts *lib.TopicService,
	ps *lib.ProjectService,
	us *lib.UserService,
	jwks auth.JWKS,
	config Config,
) {
	authenticated := Authenticated()
	anonymous := Anonymous()
	csrf := CSRF()

	aboutController := handlers.NewAboutController()
	authController := handlers.NewAuthController(us, jwks, config.Mail.SMTP, config.App.BaseURL, config.Mail.From, config.Mail.APIKey)
	homeController := handlers.NewHomeController(ps, ts)
	profileController := handlers.NewProfileController(us)
	projectController := handlers.NewProjectController(ps, ts)
	topicController := handlers.NewTopicController(ts, ps)

	// Static files.
	static := web.Static()
	fs := http.FileServer(http.FS(static))

	r.Mount("/static", fs)

	r.Group(func(r ki.Router) {
		r.Use(Session(store))
		r.Use(Authenticate(us, jwks, store))

		r.Get("/{$}", handlers.NewRedirectHandler("/home"))

		r.Route("/auth", func(r ki.Router) {
			r.Get(
				"/{$}",
				authController.PageGet,
				ki.WithMiddleware(anonymous),
				ki.WithName("auth:page:get"),
			)

			r.Route("/passwordless", func(r ki.Router) {
				r.Post("/send", authController.PasswordlessSend, ki.WithName("auth:passwordless:send"))
				r.Get("/verify", authController.PasswordlessVerify, ki.WithName("auth:passwordless:verify"))
			})

			r.Post(
				"/logout",
				handlers.NewLogoutHandler(store),
				ki.WithMiddleware(authenticated),
				ki.WithName("auth:page:logout"),
			)
		})

		r.Route("/about", func(r ki.Router) {
			r.Get("/{$}", aboutController.PageGet, ki.WithName("about:page:get"))
		})

		// Private routes.
		r.Group(func(r ki.Router) {
			r.Use(authenticated)
			r.Use(func(h http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					h.ServeHTTP(w, r)
				})
			})
			r.Use(csrf)

			r.Route("/home", func(r ki.Router) {
				r.Get("/{$}", homeController.PageGet, ki.WithName("home:page:get"))

				r.Route("/project-list", func(r ki.Router) {
					r.Get("/{$}", homeController.ProjectListGet, ki.WithName("home:project-list:get"))
				})

				r.Route("/project-create", func(r ki.Router) {
					r.Post("/submit", homeController.ProjectCreateSubmit, ki.WithName("home:project-create:submit"))
				})

				r.Route("/project-item", func(r ki.Router) {
					r.Post("/delete", homeController.ProjectItemDelete, ki.WithName("home:project-item:delete"))
				})
			})

			r.Route("/project", func(r ki.Router) {
				r.Get("/{project}", projectController.PageGet, ki.WithName("project:page:get"))
			})

			r.Route("/topic", func(r ki.Router) {
				r.Get("/{topic}", topicController.PageGet, ki.WithName("topic:page:get"))

				r.Route("/overview", func(r ki.Router) {
					r.Get("/{$}", topicController.OverviewGet, ki.WithName("topic:overview:get"))
					r.Get("/edit", topicController.OverviewEdit, ki.WithName("topic:overview:edit"))

					r.Post("/collapse", topicController.OverviewCollapse, ki.WithName("topic:overview:collapse"))
					r.Post("/expand", topicController.OverviewExpand, ki.WithName("topic:overview:expand"))

					r.Post("/save", topicController.OverviewSave, ki.WithName("topic:overview:save"))
				})

				r.Route("/descendant-list", func(r ki.Router) {
					r.Get("/{$}", topicController.DescendantList, ki.WithName("topic:descendant-list:get"))

					r.Post("/move", topicController.DescendantListMove, ki.WithName("topic:descendant-list:move"))
				})

				r.Route("/descendant-item", func(r ki.Router) {
					r.Post("/collapse", topicController.DescendantCollapse, ki.WithName("topic:descendant-item:collapse"))
					r.Post("/expand", topicController.DescendantExpand, ki.WithName("topic:descendant-item:expand"))
					r.Post("/delete", topicController.DescendantDelete, ki.WithName("topic:descendant-item:delete"))
				})

				r.Route("/descendant-create", func(r ki.Router) {
					r.Post("/submit", topicController.DescendantCreateSubmit, ki.WithName("topic:descendant-create:submit"))
				})

				r.Route("/toolbar", func(r ki.Router) {
					r.Post("/enable-selection-mode", topicController.ToolbarEnableSelectionMode, ki.WithName("topic:toolbar:enable-selection-mode"))
					r.Post("/disable-selection-mode", topicController.ToolbarDisableSelectionMode, ki.WithName("topic:toolbar:disable-selection-mode"))

					r.Post("/expand-recursive", topicController.ToolbarExpandRecursive, ki.WithName("topic:toolbar:expand-recursive"))
					r.Post("/collapse-recursive", topicController.ToolbarCollapseRecursive, ki.WithName("topic:toolbar:collapse-recursive"))
				})
			})

			r.Route("/profile", func(r ki.Router) {
				r.Get("/{$}", profileController.PageGet, ki.WithName("profile:page:get"))

				r.Route("/info", func(r ki.Router) {
					r.Get("/{$}", profileController.InfoGet, ki.WithName("profile:info:get"))
					r.Get("/edit", profileController.InfoEdit, ki.WithName("profile:info:edit"))

					r.Post("/save", profileController.InfoSave, ki.WithName("profile:info:save"))
				})
			})
		})

		r.Mount("/", handlers.NewNotFoundHanlder())
	})
}
