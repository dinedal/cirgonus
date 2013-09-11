package web

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"query"
	"strings"
	"types"
)

type Request struct {
	Name string
}

type WebHandler struct {
	Config types.CirconusConfig
}

func (wh *WebHandler) showUnauthorized(w http.ResponseWriter) {
	/* FIXME log */
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

func (wh *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !wh.handleAuth(r) {
		wh.showUnauthorized(w)
		return
	}

	switch r.Method {
	case "GET":
		{
			out, err := json.Marshal(query.AllPlugins(wh.Config))

			if err != nil {
				/* FIXME log */
				w.WriteHeader(500)
			} else {
				w.Write(out)
			}

			return
		}
	case "POST":
		{
			req := Request{}
			in, err := ioutil.ReadAll(r.Body)

			if err != nil {
				/* FIXME log */
				w.WriteHeader(500)
			}

			json.Unmarshal(in, &req)

			if req.Name != "" {
				out, err := json.Marshal(query.Plugin(req.Name, wh.Config))
				if err != nil {
					/* FIXME log */
					w.WriteHeader(500)
				} else {
					w.Write(out)
				}
			} else {
				w.WriteHeader(404)
			}
		}
	}
}

func Start(listen string, config types.CirconusConfig) error {
	s := &http.Server{
		Addr:    listen,
		Handler: &WebHandler{Config: config},
	}

	return s.ListenAndServe()
}
