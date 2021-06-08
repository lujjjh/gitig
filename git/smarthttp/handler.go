package smarthttp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/lujjjh/gitig/git/packetline"
)

const (
	protocolHeader   = "Git-Protocol"
	protocolVersion2 = "version=2"
)

var services = []struct {
	method  string
	pattern *regexp.Regexp
	impl    func(s *Handler, w http.ResponseWriter, r *http.Request)
}{
	{http.MethodGet, regexp.MustCompile("^/info/refs$"), (*Handler).InfoRefs},
	{http.MethodPost, regexp.MustCompile("/git-upload-pack"), (*Handler).ServeRPC},
}

var rpcServices = []struct {
	name string
	impl func(s *Handler, w http.ResponseWriter, r *http.Request)
}{
	{"git-upload-pack", (*Handler).UploadPack},
}

type Handler struct{}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only supports git wire protocol version 2.
	if r.Header.Get(protocolHeader) != protocolVersion2 {
		http.Error(w, "Requires git wire protocol version 2", http.StatusBadRequest)
		return
	}
	log.Printf("%s %s", r.Method, r.RequestURI)
	for _, service := range services {
		if r.Method == service.method && service.pattern.MatchString(r.URL.Path) {
			service.impl(s, w, r)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *Handler) InfoRefs(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service")
	if serviceName == "" {
		http.Error(w, "empty service", http.StatusBadRequest)
		return
	}
	for _, service := range rpcServices {
		if service.name == serviceName {
			w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", serviceName))
			service.impl(s, w, r)
			return
		}
	}
	http.Error(w, fmt.Sprintf("unknown service: %s", serviceName), http.StatusBadRequest)
}

func packetLineWriter(w http.ResponseWriter) *packetline.Writer {
	bw := (io.Writer)(w)
	if os.Getenv("DEBUG") == "true" {
		bw = io.MultiWriter(bw, os.Stderr)
	}
	return packetline.NewWriter(bw)
}

func (s *Handler) UploadPack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	plw := packetLineWriter(w)
	plw.WritePacketFmt("version 2\n").
		WritePacketFmt("ls-refs\n").
		WriteFlushPacket()
	if err := plw.Err(); err != nil {
		// TODO: logging
	}
}

func (s *Handler) ServeRPC(w http.ResponseWriter, r *http.Request) {
	plw := packetLineWriter(w)
	_ = plw.WriteFlushPacket()
	if err := plw.Err(); err != nil {
		// TODO: logging
	}
}
