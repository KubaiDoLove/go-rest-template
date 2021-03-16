package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ra9dev/PROJECTNAME/internal/ports/http/resources"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const (
	defaultAddr            = ":8080"
	readTimeout            = 5 * time.Second
	writeTimeout           = 30 * time.Second
	idleConnectionsTimeout = 3 * time.Second
	compressLevel          = 5
	cacheMaxAge            = 300
)

type Server struct {
	Address           string
	FilesDir          http.Dir
	CertFile, KeyFile *string
	IsTesting         bool

	idleConnectionsCh chan struct{}
	appCtx            context.Context
	version           string
}

func NewServer(ctx context.Context, options ...ServerOption) *Server {
	srv := &Server{
		Address: defaultAddr,

		idleConnectionsCh: make(chan struct{}),
		appCtx:            ctx,
	}

	for _, enhance := range options {
		enhance(srv)
	}

	return srv
}

func (srv *Server) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.NewCompressor(compressLevel).Handler)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(srv.IsTesting),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           cacheMaxAge,
	}))

	r.Mount("/version", resources.VersionResource{Version: srv.version}.Routes())
	if srv.IsTesting {
		filesRoute := "/files"
		r.Mount(filesRoute, resources.NewFilesResource(srv.FilesDir).Routes())
		r.Mount("/swagger", resources.SwaggerResource{FilesPath: filesRoute}.Routes())
	}

	return r
}

func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}

	return []string{}
}

func (srv *Server) Run() error {
	s := &http.Server{
		Addr:         srv.Address,
		Handler:      chi.ServerBaseContext(srv.appCtx, srv.setupRouter()),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	go srv.ListenCtxForGracefulTermination(s)
	log.Printf("Serving HTTP on \"%s\"", srv.Address)

	switch true {
	case srv.CertFile != nil && srv.KeyFile != nil:
		if err := s.ListenAndServeTLS(*srv.CertFile, *srv.KeyFile); err != nil {
			return err
		}
	default:
		if err := s.ListenAndServe(); err != nil {
			return err
		}
	}

	return nil
}

func (srv *Server) ListenCtxForGracefulTermination(s *http.Server) {
	<-srv.appCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), idleConnectionsTimeout)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] HTTP server Shutdown: %v", err)
	}

	log.Println("Processed idle connections successfully before termination")
	close(srv.idleConnectionsCh)
}

func (srv *Server) WaitForGracefulTermination() {
	<-srv.idleConnectionsCh
}
