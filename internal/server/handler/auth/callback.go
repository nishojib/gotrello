package authhandler

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"nishojib/gotrello/internal/server/handler"
	authui "nishojib/gotrello/ui/html/auth"

	"github.com/alexedwards/scs/v2"
)

func AuthCallback(session *scs.SessionManager) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		accessToken := r.URL.Query().Get("access_token")

		if len(accessToken) == 0 {
			return handler.Render(w, r, authui.CallbackScript())
		}

		session.Put(r.Context(), "access_token", accessToken)

		redirectTo, ok := session.Get(r.Context(), "redirect_to").(string)
		if !ok {
			slog.Info("no redirect url set")
		}
		session.Remove(r.Context(), "redirect_to")

		redirectUrl, err := url.QueryUnescape(redirectTo)
		if err != nil {
			slog.Info("error unescaping redirect url")
		}

		http.Redirect(w, r, fmt.Sprintf("/%s", redirectUrl), http.StatusSeeOther)
		return nil
	})
}
