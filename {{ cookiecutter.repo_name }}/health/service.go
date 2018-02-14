package health

import (
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/gtg"
)

const DefaultHealthPath = "/__health"

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

func NewHealthService(appSystemCode string, appName string, appDescription string) *HealthService {
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

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	return gtg.Status{GoodToGo: true}
}

func (service *HealthService) GTG() gtg.Status {
	return gtg.FailFastParallelCheck(service.gtgChecks)()
}
