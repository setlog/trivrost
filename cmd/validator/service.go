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
	flags ValidatorFlags
}

func (s service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	configUrl := req.URL.Query().Get(configUrlParameterName)
	if configUrl == "" {
		if s.flags.DeploymentConfigUrl != "" {
			configUrl = s.flags.DeploymentConfigUrl
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Required query parameter '"+configUrlParameterName+"' is missing.\n")
			return
		}
	}
	log.Printf("Validating deployment-config at %s...\n", configUrl)
	reps := validateDeploymentConfig(configUrl, s.flags.SkipUrlCheck, s.flags.SkipJarChek)
	logReports(reps)
	log.Printf("Finished validation of deployment-config at %s...\n", configUrl)

	haveError := reps.HaveError()
	if haveError {
		// It would be more correct to return http.StatusOK (since we were able to perform the validation, albeit with a bad outcome)
		// and have the requester parse some error message from the response body, but we decided not to do this because Prometheus'
		// monitoring cannot do much beyond checking for the response's status code.
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	s.writeReportListAsHTML(w, reps)
}

func (s service) writeReportListAsHTML(w io.Writer, reps reports) {
	io.WriteString(w, "<!doctype html>\n")
	io.WriteString(w, "<html lang=\"en\">\n")
	io.WriteString(w, "  <head><meta charset=\"UTF-8\"><title>Validator Response</title></head>\n")
	io.WriteString(w, "  <body>\n")
	if reps.HaveError() {
		io.WriteString(w, "    There were errors validating the deployment-config at "+htmlLink(s.flags.DeploymentConfigUrl)+":\n")
	} else {
		io.WriteString(w, "    The deployment-config at "+htmlLink(s.flags.DeploymentConfigUrl)+" was validated successfully:\n")
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

func actAsService(flags ValidatorFlags) {
	s := &http.Server{
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 30,
		Addr:         ":" + strconv.Itoa(flags.Port),
	}
	mux := http.NewServeMux()
	mux.Handle("/validate", service{flags: flags})
	s.Handler = mux
	log.Fatal(s.ListenAndServe())
}
