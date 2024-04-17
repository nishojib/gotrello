package authhandler

import (
	"net/http"
	"nishojib/gotrello/internal/server/handler"

	"github.com/nedpals/supabase-go"
)

func LoginWithGoogle(sbClient *supabase.Client) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		response, err := sbClient.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
			Provider:   "google",
			RedirectTo: "http://localhost:4000/auth/callback",
		})
		if err != nil {
			return err
		}

		http.Redirect(w, r, response.URL, http.StatusSeeOther)
		return nil
	})
}
