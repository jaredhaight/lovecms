package handlers

import (
	"github.com/jaredhaight/lovecms/internal"
	"log/slog"
	"net/http"
)

type SetupHandler struct {
	logger *slog.Logger
	config *internal.Config
}

func NewSetupHandler(logger slog.Logger, config internal.Config) *SetupHandler {
	return &SetupHandler{
		logger: &logger,
		config: &config,
	}
}

func (sh *SetupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
