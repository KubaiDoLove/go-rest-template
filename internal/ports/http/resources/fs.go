package resources

import (
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

type FilesResource struct {
	fs http.Handler
}

func NewFilesResource(filesDir http.Dir) *FilesResource {
	return &FilesResource{
		fs: http.FileServer(filesDir),
	}
}

func (fr FilesResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/*", fr.ServeFile)

	return r
}

func (fr FilesResource) ServeFile(w http.ResponseWriter, r *http.Request) {
	route := chi.RouteContext(r.Context()).RoutePattern()
	pathPrefix := strings.TrimSuffix(route, "/*")

	http.StripPrefix(pathPrefix, fr.fs).ServeHTTP(w, r)
}
