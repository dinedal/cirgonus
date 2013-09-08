package web

import (
	"encoding/json"
	"net/http"
	"query"
	"types"
)

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
