package git

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var services = []struct {
	method  string
	pattern *regexp.Regexp
	impl    func(s *SmartHTTPHandler, w http.ResponseWriter, r *http.Request)
}{
	{http.MethodGet, regexp.MustCompile("^/info/refs$"), (*SmartHTTPHandler).InfoRefs},
}

var rpcServices = []struct {
	name string
	impl func(s *SmartHTTPHandler, w http.ResponseWriter, r *http.Request)
}{
	{"git-upload-pack", (*SmartHTTPHandler).UploadPack},
}

type SmartHTTPHandler struct{}

func (s *SmartHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, service := range services {
		if r.Method == service.method && service.pattern.MatchString(r.URL.Path) {
			service.impl(s, w, r)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *SmartHTTPHandler) InfoRefs(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service")
	switch serviceName {
	case "":
		http.Error(w, "empty service", http.StatusBadRequest)
	default:
		for _, service := range rpcServices {
			if service.name == serviceName {
				w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", serviceName))
				service.impl(s, w, r)
				return
			}
		}
		http.Error(w, fmt.Sprintf("unknown service: %s", serviceName), http.StatusBadRequest)
	}
}

func (s *SmartHTTPHandler) UploadPack(w http.ResponseWriter, r *http.Request) {
	log.Println("UploadPack")
}
