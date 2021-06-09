package health

import (
	"fmt"

	"github.com/containerssh/http"
	"github.com/containerssh/log"
	"github.com/containerssh/service"
)

// Config is the configuration for the Service.
type Config struct {
	Enable                   bool `json:"enable" yaml:"enable"`
	http.ServerConfiguration `json:",inline" yaml:",inline" default:"{\"listen\":\"0.0.0.0:7000\"}"`
	Client                   http.ClientConfiguration `json:"client" yaml:"client" default:"{\"url\":\"http://127.0.0.1:7000/\"}"`
}

func (c Config) Validate() error {
	if !c.Enable {
		return nil
	}
	if err := c.ServerConfiguration.Validate(); err != nil {
		return err
	}
	if err := c.Client.Validate(); err != nil {
		return fmt.Errorf("invalid client configuration (%w)", err)
	}
	return nil
}

// New creates a new HTTP health service on port 23074
func New(config Config, logger log.Logger) (Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	handler := &requestHandler{}
	svc, err := http.NewServer(
		"Health check endpoint",
		config.ServerConfiguration,
		http.NewServerHandlerNegotiate(handler, logger),
		logger,
		func(url string) {
			logger.Info(log.NewMessage(MServiceAvailable, "Health check endpoint available at %s", url))
		},
	)
	if err != nil {
		return nil, err
	}

	return &healthCheckService{
		Service:        svc,
		requestHandler: handler,
	}, nil
}

// NewClient creates a new health check client based on the supplied configuration. If the health check is not enabled
// no client is returned.
func NewClient(config Config, logger log.Logger) (Client, error) {
	if !config.Enable {
		return nil, nil
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid health check configuration (%w)", err)
	}

	httpClient, err := http.NewClient(config.Client, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create health check client (%w)", err)
	}

	return &healthCheckClient{
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// Service is an HTTP service that lets you change the status via ChangeStatus().
type Service interface {
	service.Service
	ChangeStatus(ok bool)
}

// Client is the client to run health checks.
type Client interface {
	// Run runs a HTTP query against the health check service.
	Run() bool
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

func (r requestHandler) OnRequest(_ http.ServerRequest, response http.ServerResponse) error {
	if r.ok {
		response.SetBody("ok")
	} else {
		response.SetBody("not ok")
		response.SetStatus(503)
	}
	return nil
}

type healthCheckClient struct {
	httpClient http.Client
	logger     log.Logger
}

func (h *healthCheckClient) Run() bool {
	responseBody := ""
	statusCode, err := h.httpClient.Get("", &responseBody)
	if err != nil {
		h.logger.Warning(log.Wrap(err, ERequestFailed, "Request to health check endpoint failed"))
	}
	return statusCode == 200
}
