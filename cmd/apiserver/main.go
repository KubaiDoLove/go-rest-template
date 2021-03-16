// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"fmt"
	"github.com/ra9dev/PROJECTNAME/internal/ports/configs"
	"github.com/ra9dev/PROJECTNAME/internal/ports/http"
	"github.com/ra9dev/PROJECTNAME/pkg/log"
	"github.com/ra9dev/PROJECTNAME/pkg/os"
	stdlog "log"

	"github.com/ra9dev/PROJECTNAME/internal/adapters/database"
	"github.com/ra9dev/PROJECTNAME/internal/adapters/database/drivers"
)

var version = "unknown"

func main() {
	fmt.Printf("PROJECTNAME %s\n", version)

	appCtx, cancelAppCtx := context.WithCancel(context.Background())
	defer cancelAppCtx()
	go os.CatchTermination(cancelAppCtx)

	config := configs.NewAppConfig()
	log.Setup(config.LogMode)

	ds := database.New(drivers.DataStoreConfig{
		URL:  config.DataStoreURL,
		Name: config.DataStoreName,
		DB:   config.DataStoreDB,
	})
	if err := ds.Connect(); err != nil {
		stdlog.Printf("[ERROR] cannot connect to datastore %s: %v", ds.Name(), err)
		return
	}
	stdlog.Printf("Connected to %s", ds.Name())
	defer ds.Close(appCtx)

	httpSrv := http.NewServer(
		appCtx,
		http.WithVersion(version),
		http.WithCustomAddress(config.ListenAddr),
		http.WithSSL(config.CertFile, config.KeyFile),
		http.WithFiles(config.FilesDir),
		http.WithTestingMode(config.IsTesting),
	)
	if err := httpSrv.Run(); err != nil {
		stdlog.Println(err)
		return
	}

	httpSrv.WaitForGracefulTermination()
	stdlog.Printf("[WARN] process terminated")
}
