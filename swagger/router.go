package swagger

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, filepath.Join(r.URL.Path, "/index.html"), http.StatusMovedPermanently)
	})
	r.Get("/*", httpSwagger.Handler())

	return r
}
