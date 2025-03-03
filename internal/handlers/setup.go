package handlers

import (
	"github.com/jaredhaight/lovecms/internal/application"
	"log/slog"
	"net/http"
)

type SetupHandler struct {
	logger *slog.Logger
	config *application.Config
}

func NewSetupHandler(logger slog.Logger, config application.Config) *SetupHandler {
	return &SetupHandler{
		logger: &logger,
		config: &config,
	}
}

func (sh *SetupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
