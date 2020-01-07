package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type service struct {
	flags ValidatorFlags
}

func (s service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	errs := validateDeploymentConfig(s.flags.DeploymentConfigUrl, s.flags.SkipUrlCheck, s.flags.SkipJarChek)
	if len(errs) > 0 {
		// It would be more correct to return http.StatusOK (since we were able to perform the validation, albeit with a bad outcome)
		// and have the requestor parse some error message from the response body, but we decided not to do this because Prometheus'
		// monitoring cannot do much beyond checking for the response's status code.
		w.WriteHeader(http.StatusInternalServerError)

		s.writeHtmlErrorList(w, errs)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func (s service) writeHtmlErrorList(w io.Writer, errs []error) {
	io.WriteString(w, "<!doctype html>\n")
	io.WriteString(w, "<html lang=\"en\">\n")
	io.WriteString(w, "  <head><meta charset=\"UTF-8\"><title>Validator Response</title></head>\n")
	io.WriteString(w, "  <body>\n")
	io.WriteString(w, "    There were errors validating the deployment-config at "+htmlLink(s.flags.DeploymentConfigUrl)+":\n")
	io.WriteString(w, "    <ul>\n")
	for _, err := range errs {
		io.WriteString(w, "      <li>"+html.EscapeString(err.Error())+"</li>\n")
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
