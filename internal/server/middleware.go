package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"nishojib/gotrello/internal/data/models"
	"nishojib/gotrello/internal/server/handler"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

// func (s *Server) authenticate() MiddlewareFunc {
// return func(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Vary", "Authorization")

// 		authorizationHeader := r.Header.Get("Authorization")

// 		if authorizationHeader == "" {
// 			r = s.contextSetUser(r, data.AnonymousUser)
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		headerParts := strings.Split(authorizationHeader, " ")
// 		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
// 			s.invalidAuthenticationTokenResponse(w, r)
// 			return
// 		}

// 		token := headerParts[1]

// 		v := validator.New()

// 		if data.ValidateTokenPlainText(v, token); !v.Valid() {
// 			s.invalidAuthenticationTokenResponse(w, r)
// 			return
// 		}

// 		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, data.ErrRecordNotFound):
// 				s.invalidAuthenticationTokenResponse(w, r)
// 			default:
// 				s.serverErrorResponse(w, r, err)
// 			}
// 			return
// 		}

// 		r = app.contextSetUser(r, user)

// 		next.ServeHTTP(w, r)
// 	})

// }

// func requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := app.contextGetUser(r)

// 		if user.IsAnonymous() {
// 			app.authenticationRequiredResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
// 	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := app.contextGetUser(r)

// 		if !user.Activated {
// 			app.inactiveAccountResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})

// 	return app.requireAuthenticatedUser(fn)
// }

// func requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
// 	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := app.contextGetUser(r)

// 		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
// 		if err != nil {
// 			app.serverErrorResponse(w, r, err)
// 			return
// 		}

// 		if !permissions.Include(code) {
// 			app.notPermittedResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})

// 	return app.requireActivatedUser(fn)
// }

func requireAccessToken(
	session *scs.SessionManager,
	sbClient *supabase.Client,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "static") {
				next.ServeHTTP(w, r)
				return
			}

			accessToken, ok := session.Get(r.Context(), "access_token").(string)
			if !ok {
				slog.Error(fmt.Sprintf("malformed access token %v", accessToken))
				next.ServeHTTP(w, r)
				return
			}

			response, err := sbClient.Auth.User(r.Context(), accessToken)
			if err != nil {
				slog.Error("couldn't find user")
				next.ServeHTTP(w, r)
				return
			}

			userID, err := uuid.Parse(response.ID)
			if err != nil {
				slog.Error(fmt.Sprintf("couldn't parse uuid, %s", response.ID))
				next.ServeHTTP(w, r)
				return
			}

			user := models.AuthenticatedUser{
				ID:          userID,
				Email:       response.Email,
				AccessToken: accessToken,
				LoggedIn:    true,
			}

			ctx := context.WithValue(r.Context(), models.UserContextKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func requireAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "static") {
			next.ServeHTTP(w, r)
			return
		}

		user := handler.GetAuthenticatedUser(r)
		if !user.LoggedIn {
			path := r.URL.Path
			path = strings.Replace(path, "/", "", 1)
			path = url.QueryEscape(path)
			http.Redirect(w, r, fmt.Sprintf("/login?redirect_to=%s", path), http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
