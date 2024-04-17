package authhandler

import (
	"fmt"
	"net/http"
	"nishojib/gotrello/internal/server/handler"

	"github.com/alexedwards/scs/v2"
	"github.com/nedpals/supabase-go"
)

func Logout(session *scs.SessionManager, sbClient *supabase.Client) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		accessToken, ok := session.Get(r.Context(), "access_token").(string)
		if !ok {
			return fmt.Errorf("malformed access token %v", accessToken)
		}

		if err := sbClient.Auth.SignOut(r.Context(), accessToken); err != nil {
			return err
		}

		session.Remove(r.Context(), "access_token")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	})
}
