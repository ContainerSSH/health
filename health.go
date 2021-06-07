package health

import (
	"github.com/containerssh/http"
	"github.com/containerssh/log"
	"github.com/containerssh/service"
)

// New creates a new HTTP health service on port 23074
func New(logger log.Logger) (service.Service, error) {

	handler := http.NewServerHandler(&requestHandler{}, logger)
	return http.NewServer(
		"health",
		http.ServerConfiguration{Listen: "127.0.0.1:23074"},
		handler,
		logger,
		func(url string) {},
	)
}

type requestHandler struct {
	ok bool
}

func (r requestHandler) OnRequest(request http.ServerRequest, response http.ServerResponse) error {
	response.SetBody("ok")
	return nil
}
