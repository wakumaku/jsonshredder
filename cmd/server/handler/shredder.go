package handler

import (
	"fmt"
	"io"
	"net/http"

	"wakumaku/jsonshredder/internal/service"

	"github.com/gorilla/mux"
)

// ShredderForwarder proxies the result
func ShredderForwarder(svc *service.Shredder, forwardSvc *service.Forwarder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"message": "error reading body: %s"}`, err)
			return
		}
		defer r.Body.Close()

		out, err := svc.Shred(vars["transformation"], body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"message": "error shredding: %s"}`, err)
			return
		}

		if forwardSvc != nil {
			if err := forwardSvc.Forward(r.Context(), vars["forwarder"], out); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"message": "error forwarding: %s"}`, err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("X-APPLICATION", "jsonshredder")
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", out)
	}
}
