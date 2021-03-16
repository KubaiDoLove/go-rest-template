package resources

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type VersionResponse struct {
	Version string `json:"version"`
}

type VersionResource struct {
	Version string
}

func (vr VersionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", vr.CurrentVersion)

	return r
}

func (vr VersionResource) CurrentVersion(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, VersionResponse{
		Version: vr.Version,
	})
}
