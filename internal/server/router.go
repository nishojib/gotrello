package server

import (
	"fmt"
	"net/http"
	authhandler "nishojib/gotrello/internal/server/handler/auth"
	homehandler "nishojib/gotrello/internal/server/handler/home"
	projecthandler "nishojib/gotrello/internal/server/handler/project"
	taskhandler "nishojib/gotrello/internal/server/handler/task"
	"nishojib/gotrello/ui"
	"time"

	"github.com/alexedwards/scs/bunstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/nedpals/supabase-go"
	"github.com/uptrace/bun"
)

// New returns a new http.Handler that routes requests to the correct handler.
func Routes(db *bun.DB, sbClient *supabase.Client, rps int, limiterEnabled bool) http.Handler {
	var err error
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	if sessionManager.Store, err = bunstore.New(db); err != nil {
		panic(fmt.Errorf("error creating session manager: %w", err))
	}

	router := chi.NewRouter()

	if limiterEnabled {
		router.Use(httprate.LimitByIP(rps, 1*time.Minute))
	}

	router.Use(
		middleware.Recoverer,
		sessionManager.LoadAndSave,
		requireAccessToken(sessionManager, sbClient),
	)

	router.Handle("/*", http.FileServer(http.FS(ui.Files)))

	router.Mount("/debug", middleware.Profiler())

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		notFound(w)
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		methodNotAllowed(w)
	})

	router.Get("/", homehandler.Index())

	router.Get("/login", authhandler.Login(sessionManager))
	router.Post("/login", authhandler.LoginCreate(sbClient))
	router.Get("/login/provider/google", authhandler.LoginWithGoogle(sbClient))
	router.Post("/logout", authhandler.Logout(sessionManager, sbClient))
	router.Get("/auth/callback", authhandler.AuthCallback(sessionManager))

	router.Group(func(authRouter chi.Router) {
		authRouter.Use(requireAuth)

		authRouter.Get("/projects", projecthandler.List(db))
		authRouter.Get("/projects/{projectID}", projecthandler.Show(db))

		authRouter.Post("/tasks", taskhandler.Create(db))
		authRouter.Delete("/tasks/{taskID}", taskhandler.Delete(db))
		authRouter.Post("/move-item", taskhandler.Move(db))

	})

	return router
}
