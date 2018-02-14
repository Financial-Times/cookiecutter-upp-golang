package health

import fthealth "github.com/Financial-Times/go-fthealth/v1_1"

func (service *HealthService) sampleCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Sample healthcheck has no impact",
		Name:             "Sample healthcheck",
		PanicGuide:       "https://dewey.in.ft.com/view/system/{{ cookiecutter.system_code }}",
		Severity:         1,
		TechnicalSummary: "Sample healthcheck has no technical details",
		Checker:          service.sampleChecker,
	}
}

func (service *HealthService) sampleChecker() (string, error) {
	return "Sample is healthy", nil
}
