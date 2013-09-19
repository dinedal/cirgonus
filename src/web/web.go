package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"plugins/record"
	"query"
	"strings"
	"types"
)

type Request struct {
	Name  string
	Value interface{}
}

type WebHandler struct {
	Config types.CirconusConfig
	Logger *syslog.Writer
}

func (wh *WebHandler) showUnauthorized(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", "Basic realm=\"cirgonus\"")
	w.WriteHeader(401)
}

func (wh *WebHandler) handleAuth(r *http.Request) bool {
	header, ok := r.Header["Authorization"]

	if !ok {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.Split(header[0], " ")[1])

	if err != nil {
		return false
	}

	credentials := strings.Split(string(decoded), ":")

	if credentials[0] != wh.Config.Username || credentials[1] != wh.Config.Password {
		return false
	}

	return true
}

func (wh *WebHandler) readAndUnmarshal(w http.ResponseWriter, r *http.Request, requestType string) Request {
	req := Request{}
	in, err := ioutil.ReadAll(r.Body)

	wh.Logger.Debug(fmt.Sprintf("Handling %s with payload '%s'", requestType, in))

	if err != nil {
		wh.Logger.Crit(fmt.Sprintf("Error encountered reading: %s", err))
		w.WriteHeader(500)
	}

	json.Unmarshal(in, &req)

	return req
}

func (wh *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !wh.handleAuth(r) {
		wh.Logger.Info(fmt.Sprintf("Unauthorized access from %s", r.RemoteAddr))
		wh.showUnauthorized(w)
		return
	}

	switch r.Method {
	case "GET":
		{
			wh.Logger.Debug("Handling GET")

			out, err := json.Marshal(query.AllPlugins(wh.Config, wh.Logger))

			if err != nil {
				wh.Logger.Crit(fmt.Sprintf("Error marshalling all metrics: %s", err))
				w.WriteHeader(500)
			} else {
				wh.Logger.Debug(fmt.Sprintf("Writing all metrics to %s", r.RemoteAddr))
				w.Write(out)
			}

			return
		}
	case "POST":
		{
			req := wh.readAndUnmarshal(w, r, "POST")

			if req.Name != "" {
				out, err := json.Marshal(query.Plugin(req.Name, wh.Config, wh.Logger))
				if err != nil {
					wh.Logger.Crit(fmt.Sprintf("Error gathering metrics for %s: %s", req.Name, err))
					w.WriteHeader(500)
				} else {
					wh.Logger.Debug(fmt.Sprintf("Handling POST for metric '%s'", req.Name))
					w.Write(out)
				}
			} else {
				wh.Logger.Debug(fmt.Sprintf("404ing because no payload from %s", r.RemoteAddr))
				w.WriteHeader(404)
			}
		}
	case "PUT":
		{
			req := wh.readAndUnmarshal(w, r, "PUT")
			if req.Name == "" {
				wh.Logger.Crit(fmt.Sprintf("Cannot write record with an empty value"))
				w.WriteHeader(500)
			} else {
				wh.Logger.Debug("here")
				record.RecordMetric(req.Name, req.Value, wh.Logger)
				wh.Logger.Debug("here")
				w.WriteHeader(200)
				w.Write([]byte("OK"))
			}
		}
	}
}

func Start(listen string, config types.CirconusConfig, log *syslog.Writer) error {
	log.Info("Starting Web Service")

	s := &http.Server{
		Addr:    listen,
		Handler: &WebHandler{Config: config, Logger: log},
	}

	return s.ListenAndServe()
}
