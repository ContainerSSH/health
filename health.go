package health

import (
	"github.com/containerssh/http"
	"github.com/containerssh/log"
	"github.com/containerssh/service"
)

// Config is the configuration for the HealthCheckService.
type Config struct {
	http.ServerConfiguration
}

// New creates a new HTTP health service on port 23074
func New(config Config, logger log.Logger) (HealthCheckService, error) {

	rh := &requestHandler{}

	handler := http.NewServerHandler(rh, logger)
	svc, err := http.NewServer(
		"health",
		config.ServerConfiguration,
		handler,
		logger,
		func(url string) {},
	)

	if err != nil {
		return nil, err
	}

	return &healthCheckService{
		Service:        svc,
		requestHandler: rh,
	}, nil
}

// HealthCheckService is an HTTP service that lets you change the status via ChangeStatus().
type HealthCheckService interface {
	service.Service
	ChangeStatus(ok bool)
}

type healthCheckService struct {
	service.Service
	requestHandler *requestHandler
}

func (h *healthCheckService) ChangeStatus(ok bool) {
	h.requestHandler.ok = ok
}

type requestHandler struct {
	ok bool
}

func (r requestHandler) OnRequest(request http.ServerRequest, response http.ServerResponse) error {
	if r.ok {
		response.SetBody("ok")
	} else {
		response.SetBody("not ok")
		response.SetStatus(503)
	}
	return nil
}
