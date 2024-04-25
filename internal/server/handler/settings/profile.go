package settingshandler

import (
	"net/http"
	"nishojib/gotrello/internal/server/handler"
	settingsui "nishojib/gotrello/ui/html/settings"
)

func Profile() http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		return handler.Render(w, r, settingsui.Profile())
	})
}
