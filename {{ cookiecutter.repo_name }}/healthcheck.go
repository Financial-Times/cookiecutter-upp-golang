package main

import (
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/gtg"
)

const healthPath = "/__health"

type HealthService struct {
}

func (service *HealthService) Health(appSystemCode string, appName string, description string) func(w http.ResponseWriter, r *http.Request) {
	checks := []fthealth.Check{service.HealthCheck()}
	hc := fthealth.HealthCheck{
		SystemCode:  appSystemCode,
		Name:        appName,
		Description: appDescription,
		Checks:      checks,
	}
	return fthealth.Handler(hc)
}

func (service *HealthService) HealthCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Sample healthcheck has no impact",
		Name:             "Sample healthcheck",
		PanicGuide:       "https://dewey.ft.com/cookie-cutter-test.html",
		Severity:         1,
		TechnicalSummary: "Sample healthcheck has no technical details",
		Checker:          service.Checker,
	}
}

func (service *HealthService) Checker() (string, error) {
	return "Sample is healthy", nil
}

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	return gtg.Status{GoodToGo: true}
}

func (service *HealthService) GTG() gtg.Status {
	check := func() gtg.Status {
		return gtgCheck(service.Checker)
	}

	return gtg.FailFastParallelCheck([]gtg.StatusChecker{check})()
}
