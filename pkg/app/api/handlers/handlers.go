package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/oneils/ynab-helper/bot/platform/web"
	"go.mongodb.org/mongo-driver/mongo"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *mongo.Database) http.Handler {
	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown)

	// Register debug check endpoints.
	cg := checkGroup{
		build: build,
		db:    db,
	}
	app.Handle(http.MethodGet, "/v1/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/v1/liveness", cg.liveness)

	return app
}
