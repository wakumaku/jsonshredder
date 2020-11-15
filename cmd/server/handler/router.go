package handler

import (
	"wakumaku/jsonshredder/internal/service"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

// Router returns a router with all the paths and handlers
func Router(shredderSvc *service.Shredder, forwardSvc *service.Forwarder, logger *zerolog.Logger) *mux.Router {
	// Builds router handler paths
	router := mux.NewRouter()

	// Transformations endpoint
	router.HandleFunc("/{transformation}", loggerHandler(ShredderForwarder(shredderSvc, nil), logger))

	// Transform and Forward endpoint
	router.HandleFunc("/{transformation}/{forwarder}", loggerHandler(ShredderForwarder(shredderSvc, forwardSvc), logger))

	return router
}
