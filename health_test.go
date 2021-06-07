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
	srv, err := health.New(
		health.Config{
			ServerConfiguration: http.ServerConfiguration{
				Listen: "127.0.0.1:23074",
			},
		},
		logger)

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
		URL:            "http://127.0.0.1:23074",
		AllowRedirects: false,
		Timeout:        5 * time.Second,
	},
		logger,
	)

	if err != nil {
		t.Fatal(err)
	}

	srv.ChangeStatus(true)
	checkStatusResponse(t, client, 200, "ok")

	srv.ChangeStatus(false)

	checkStatusResponse(t, client, 503, "not ok")
}

func checkStatusResponse(t *testing.T, client http.Client, expectedStatusCode int, expectedResponse string) {
	response := ""
	status, err := client.Get("", &response)

	if err != nil {
		t.Fatal(err)
	}

	if status != expectedStatusCode {
		t.Fatalf("Unexpected status code: %d", status)
	}

	if response != expectedResponse {
		t.Fatalf("Unexpected response: %s", response)
	}
}
