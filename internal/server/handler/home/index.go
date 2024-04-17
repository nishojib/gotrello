package homehandler

import (
	"net/http"
	"nishojib/gotrello/internal/server/handler"
	homeui "nishojib/gotrello/ui/html/home"
)

func Index() http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		return handler.Render(w, r, homeui.Index())
	})
}
