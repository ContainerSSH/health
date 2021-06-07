[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH Health Check Library</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/health?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/health)

This is a health check service returning "ok" if all required ContainerSSH services are running.

<p align="center"><strong>⚠⚠⚠ Warning: This is a developer documentation. ⚠⚠⚠</strong><br />The user documentation for ContainerSSH is located at <a href="https://containerssh.io">containerssh.io</a>.</p>

## Using this service 

This library uses ContainerSSH' own [HTTP](https://github.com/containerssh/http) implementation to create an HTTP server that returns "ok" when all services are up.

You can instantiate this service as described in the [service library](https://github.com/containerssh/service) as follows:

```go
srv, err := health.New(logger)

if err != nil {
	t.Fatal(err)
}
```