package main

import (
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/gtg"
	"time"
)

type HealthService struct {
	config       *HealthConfig
	healthChecks []fthealth.Check
	gtgChecks    []gtg.StatusChecker
}

type HealthConfig struct {
	appSystemCode  string
	appName        string
	appDescription string
}

func newHealthService(appSystemCode string, appName string, appDescription string) *HealthService {
	hc := &HealthService{
		config: &HealthConfig{
			appSystemCode:  appSystemCode,
			appName:        appName,
			appDescription: appDescription,
		},
	}
	hc.healthChecks = []fthealth.Check{hc.sampleCheck()}
	check := func() gtg.Status {
		return gtgCheck(hc.sampleChecker)
	}
	var gtgChecks []gtg.StatusChecker
	gtgChecks = append(hc.gtgChecks, check)
	hc.gtgChecks = gtgChecks
	return hc
}

func (service *HealthService) Health() fthealth.HC {
	return &fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{
			SystemCode:  service.config.appSystemCode,
			Name:        service.config.appName,
			Description: service.config.appDescription,
			Checks:      service.healthChecks,
		},
		Timeout: 10 * time.Second,
	}
}

func (service *HealthService) sampleCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Sample healthcheck has no impact",
		Name:             "Sample healthcheck",
		PanicGuide:       "https://dewey.ft.com/{{ cookiecutter.system_code }}.html",
		Severity:         1,
		TechnicalSummary: "Sample healthcheck has no technical details",
		Checker:          service.sampleChecker,
	}
}

func (service *HealthService) sampleChecker() (string, error) {
	return "Sample is healthy", nil
}

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	return gtg.Status{GoodToGo: true}
}

func (service *HealthService) GTG() gtg.Status {
	return gtg.FailFastParallelCheck(service.gtgChecks)()
}
