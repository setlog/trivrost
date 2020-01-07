package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
)

type service struct {
	flags ValidatorFlags
}

func (s service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	errs := validateDeploymentConfig(s.flags.DeploymentConfigUrl, s.flags.SkipUrlCheck, s.flags.SkipJarChek)
	if len(errs) > 0 {
		// It would be more correct to return http.StatusOK (since we were able to perform the validation, albeit with a bad outcome)
		// and have the requestor parse some error message from the response body, but we decided not to do this because Prometheus'
		// monitoring cannot do much beyond checking for the response's status code.
		w.WriteHeader(http.StatusInternalServerError)

		writeHtmlErrorList(w, s.flags.DeploymentConfigUrl, errs)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func writeHtmlErrorList(w io.Writer, deploymentConfigUrl string, errs []error) {
	io.WriteString(w, "<html>\n")
	io.WriteString(w, "  <head><meta charset=\"UTF-8\"><title>Validator Response</title></head>\n")
	io.WriteString(w, "  <body>\n")
	io.WriteString(w, "    There were errors validating the deployment-config at "+htmlLink(deploymentConfigUrl)+":\n")
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
	handler := service{flags: flags}
	http.Handle("/validate", handler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(flags.Port), nil))
}
