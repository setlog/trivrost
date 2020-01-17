package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const configUrlParameterName = "configurl"

type service struct {
	flags *ValidatorFlags
}

func (s service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	configUrl := getConfigUrl(req, s.flags)
	if configUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Required query parameter '"+configUrlParameterName+"' is missing.\n")
		return
	}
	log.Printf("Validating deployment-config at %s...\n", configUrl)
	reps := validateDeploymentConfig(configUrl, s.flags.SkipUrlCheck, s.flags.SkipJarChek)
	logReports(reps, false)
	log.Printf("Finished validation of deployment-config at %s...\n", configUrl)

	w.WriteHeader(http.StatusOK)
	writeReportListAsHTML(w, reps, configUrl)
}

func getConfigUrl(req *http.Request, flags *ValidatorFlags) string {
	configUrl := req.URL.Query().Get(configUrlParameterName)
	if configUrl == "" {
		configUrl = flags.DeploymentConfigUrl
	}
	return configUrl
}

func writeReportListAsHTML(w io.Writer, reps reports, deploymentConfigUrl string) {
	io.WriteString(w, "<!doctype html>\n")
	io.WriteString(w, "<html lang=\"en\">\n")
	io.WriteString(w, "  <head><meta charset=\"UTF-8\"><title>Validator Response</title></head>\n")
	io.WriteString(w, "  <body>\n")
	if reps.HaveError() {
		io.WriteString(w, "    There were errors validating the deployment-config at "+htmlLink(deploymentConfigUrl)+":\n")
	} else {
		io.WriteString(w, "    The deployment-config at "+htmlLink(deploymentConfigUrl)+" was validated successfully:\n")
	}
	io.WriteString(w, "    <ul>\n")
	for _, rep := range reps {
		if rep.isError {
			io.WriteString(w, "      <li><span style=\"color:#f21310;\">"+html.EscapeString(rep.message)+"</span></li>\n")
		} else {
			io.WriteString(w, "      <li>"+html.EscapeString(rep.message)+"</li>\n")
		}
	}
	io.WriteString(w, "    </ul>\n")
	io.WriteString(w, "  </body>\n")
	io.WriteString(w, "</html>\n")
}

func htmlLink(url string) string {
	return "<a href=\"" + url + "\">" + html.EscapeString(url) + "</a>"
}

func actAsService(flags *ValidatorFlags) {
	httpServer := &http.Server{
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 30,
		Addr:         ":" + strconv.Itoa(flags.Port),
	}
	mux := http.NewServeMux()
	s := service{flags: flags}
	mux.Handle("/validate", s)
	mux.Handle("/metrics", metrics{service: s})
	httpServer.Handler = mux
	log.Printf("Listening for /validate and /metrics HTTP requests on port %d.\n", flags.Port)
	log.Fatal(httpServer.ListenAndServe())
}
