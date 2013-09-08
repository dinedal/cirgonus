package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"query"
	"types"
)

type Request struct {
	Name string
}

type WebHandler struct {
	Config types.CirconusConfig
}

func (wh *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case "GET":
		{
			out, err := json.Marshal(query.AllPlugins(wh.Config))

			if err != nil {
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
