package health_test

import (
	"github.com/containerssh/health"
	"github.com/containerssh/http"
	"github.com/containerssh/log"
	"github.com/containerssh/service"
	"testing"
	"time"
)

func TestOk(t *testing.T) {
	logger := log.NewTestLogger(t)
	config := health.Config{
		Enable: true,
		ServerConfiguration: http.ServerConfiguration{
			Listen: "127.0.0.1:23074",
		},
		Client: http.ClientConfiguration{
			URL:            "http://127.0.0.1:23074",
			AllowRedirects: false,
			Timeout:        5 * time.Second,
		},
	}

	srv, err := health.New(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	l := service.NewLifecycle(srv)

	running := make(chan struct{})
	l.OnRunning(func(s service.Service, l service.Lifecycle) {
		running <- struct{}{}
	})

	go func() {
		_ = l.Run()
	}()

	<-running

	client, err := health.NewClient(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	if client.Run() {
		t.Fatal("Health check did not fail, even though status is false.")
	}

	srv.ChangeStatus(true)
	if !client.Run() {
		t.Fatal("Health check failed, even though status is true.")
	}

	srv.ChangeStatus(false)
	if client.Run() {
		t.Fatal("Health check did not fail, even though status is false.")
	}
}
