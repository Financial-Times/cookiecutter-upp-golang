package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"sync"
	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
{% if cookiecutter.add_sample_http_endpoint == "yes" %}
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
{% endif %}
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

	log.SetLevel(log.InfoLevel)
	log.Infof("[Startup] {{ cookiecutter.service_name }} is starting ")

	app.Action = func() {
		log.Infof("System code: %s, App Name: %s, Port: %s", *appSystemCode, *appName, *port)

		go func() {
			serveEndpoints(*appSystemCode, *appName, *port{%- if cookiecutter.add_sample_http_endpoint == "yes" -%}, requestHandler{}{%- endif -%})
		}()

		// todo: insert app code here

		waitForSignal()
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}

func serveEndpoints(appSystemCode string, appName string, port string{%- if cookiecutter.add_sample_http_endpoint == "yes" -%}, requestHandler requestHandler{%- endif -%}) {
	healthService := HealthService{}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc(healthPath, http.HandlerFunc(healthService.Health(appSystemCode, appName, appDescription)))
	serveMux.HandleFunc(status.GTGPath, status.NewGoodToGoHandler(healthService.GTG))
	serveMux.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
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
			log.Infof("HTTP server closing with message: %v", err)
		}
		wg.Done()
	}()

	waitForSignal()
	log.Infof("[Shutdown] {{ cookiecutter.service_name }} is shutting down")

	if err := server.Close(); err != nil {
		log.Errorf("Unable to stop http server: %v", err)
	}

	wg.Wait()
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
