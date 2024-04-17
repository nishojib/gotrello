package authhandler

import (
	"net/http"
	"nishojib/gotrello/internal/server/handler"
	"nishojib/gotrello/internal/validator"
	authui "nishojib/gotrello/ui/html/auth"

	"github.com/alexedwards/scs/v2"
	"github.com/nedpals/supabase-go"
)

func Login(session *scs.SessionManager) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		path := r.URL.Query().Get("redirect_to")
		if len(path) == 0 {
			path = "/"
		}

		session.Put(r.Context(), "redirect_to", path)

		return handler.Render(w, r, authui.Login())
	})
}

func LoginCreate(sbClient *supabase.Client) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		email := r.FormValue("email")

		v := validator.New()
		v.Check(email != "", "email", "Email must be provided")
		v.Check(
			validator.Matches(email, validator.EmailRX),
			"email",
			"Email must be a valid",
		)
		if !v.Valid() {
			return handler.Render(
				w,
				r,
				authui.LoginForm(authui.LoginParams{Email: email}, authui.LoginErrors{
					Email: v.Errors["email"],
				}),
			)
		}

		if err := sbClient.Auth.SendMagicLink(r.Context(), email); err != nil {
			return handler.Render(
				w,
				r,
				authui.LoginForm(
					authui.LoginParams{Email: email},
					authui.LoginErrors{Email: err.Error()},
				),
			)
		}

		return handler.Render(w, r, authui.MagicLinkSuccess(email))
	})
}
