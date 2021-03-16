package resources

import (
	"github.com/go-chi/chi"
	_ "github.com/ra9dev/PROJECTNAME/api"
	"github.com/swaggo/http-swagger"
	"path/filepath"
)

type SwaggerResource struct {
	FilesPath string
}

func (sr SwaggerResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/*", httpSwagger.Handler(
		httpSwagger.URL(filepath.Join(sr.FilesPath, "swagger.json")),
	))
	return r
}
