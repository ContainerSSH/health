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
	srv, err := health.New(logger)

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

	client, err := http.NewClient(http.ClientConfiguration{
		URL:            "http://localhost:23074",
		AllowRedirects: false,
		Timeout:        5 * time.Second,
	},
		logger,
	)

	if err != nil {
		t.Fatal(err)
	}

	response := ""

	srv.ChangeStatus(true)
	status, err := client.Get("", &response)

	if err != nil {
		t.Fatal(err)
	}

	if status != 200 {
		t.Fatalf("Unexpected status code: %d", status)
	}

	if response != "ok" {
		t.Fatalf("Unexpected response: %s", response)
	}

	srv.ChangeStatus(false)
	status, err = client.Get("", &response)

	if err != nil {
		t.Fatal(err)
	}

	if status != 503 {
		t.Fatalf("Unexpected status code: %d", status)
	}

	if response != "not ok" {
		t.Fatalf("Unexpected response: %s", response)
	}
}
