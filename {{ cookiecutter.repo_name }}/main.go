package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"sync"

   "github.com/Financial-Times/{{ cookiecutter.repo_name }}/health"

	"github.com/jawher/mow.cli"
   api "github.com/Financial-Times/api-endpoint"
	log "github.com/sirupsen/logrus"
{% if cookiecutter.add_sample_http_endpoint == "yes" %}
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
{% endif %}
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	status "github.com/Financial-Times/service-status-go/httphandlers"
)

const appDescription = "{{ cookiecutter.project_short_description }}"

func main() {
	app := cli.App("{{ cookiecutter.service_name }}", appDescription)

	appSystemCode := app.String(cli.StringOpt{
		Name:   "app-system-code",
		Value:  "{{ cookiecutter.system_code }}",
		Desc:   "System Code of the application",
		EnvVar: "APP_SYSTEM_CODE",
	})

	appName := app.String(cli.StringOpt{
		Name:   "app-name",
		Value:  "{{ cookiecutter.app_name }}",
		Desc:   "Application name",
		EnvVar: "APP_NAME",
	})

	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})

   apiYml := app.String(cli.StringOpt{
		Name:   "api-yml",
		Value:  "./api.yml",
		Desc:   "Location of the OpenAPI YML file.",
		EnvVar: "API_YML",
	})

   log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Infof("[Startup] {{ cookiecutter.service_name }} is starting ")

	app.Action = func() {
		log.Infof("System code: %s, App Name: %s, Port: %s", *appSystemCode, *appName, *port)

		go func() {
			serveEndpoints(*appSystemCode, *appName, *port, apiYml{%- if cookiecutter.add_sample_http_endpoint == "yes" -%}, requestHandler{}{%- endif -%})
		}()

		// todo: insert app code here

		waitForSignal()
	}
	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Error("{{ cookiecutter.repo_name }} could not start!")
		return
	}
}

func serveEndpoints(appSystemCode string, appName string, port string, apiYml *string{%- if cookiecutter.add_sample_http_endpoint == "yes" -%}, requestHandler requestHandler{%- endif -%}) {
	healthService := health.NewHealthService(appSystemCode, appName, appDescription)

	serveMux := http.NewServeMux()

	serveMux.HandleFunc(health.DefaultHealthPath, http.HandlerFunc(fthealth.Handler(healthService.Health())))
	serveMux.HandleFunc(status.GTGPath, status.NewGoodToGoHandler(healthService.GTG))
	serveMux.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)

   if apiYml != nil {
		apiEndpoint, err := api.NewAPIEndpointForFile(*apiYml)
		if err != nil {
			log.WithError(err).WithField("file", apiYml).Warn("Failed to serve the API Endpoint for this service. Please validate the file exists, and that it fits the OpenAPI specification.")
		} else {
			serveMux.HandleFunc(api.DefaultPath, apiEndpoint.ServeHTTP)
		}
	}

{% if cookiecutter.add_sample_http_endpoint == "yes" %}
	servicesRouter := mux.NewRouter()
	servicesRouter.HandleFunc("/sample", requestHandler.sampleMessage).Methods("POST")
	//todo: add new handlers here

	var monitoringRouter http.Handler = servicesRouter
	monitoringRouter = httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), monitoringRouter)
	monitoringRouter = httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry, monitoringRouter)

	serveMux.Handle("/", monitoringRouter)
{% endif %}
	server := &http.Server{Addr: ":" + port, Handler: serveMux}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.WithError(err).Info("HTTP server closing with message")
		}
		wg.Done()
	}()

	waitForSignal()
	log.Infof("[Shutdown] {{ cookiecutter.service_name }} is shutting down")

	if err := server.Close(); err != nil {
		log.WithError(err).Error("Unable to stop http server")
	}

	wg.Wait()
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
