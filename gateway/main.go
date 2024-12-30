package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/okusarobert/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

// decorator pattern or high order functions
type apiFunc func(w http.ResponseWriter, r *http.Request) (int, error)

func main() {
	httpListenAddr := flag.String("httpAddr", ":6000", "http transport server listen port")
	flag.Parse()
	var (
		c              = client.NewHTTPClient("http://localhost:3000")
		invoiceHandler = newInvoiceHandler(c)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("http gateway server started on port %s", *httpListenAddr)
	log.Fatal(http.ListenAndServe(*httpListenAddr, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) (int, error) {
	obuID := r.URL.Query().Get("obu")
	if obuID == "" {
		return http.StatusBadRequest, fmt.Errorf("invalid OBU ID")
	} else {
		id, err := strconv.Atoi(obuID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("internal server error")
		}
		invoice, err := h.client.GetInvoice(r.Context(), id)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = writeJSON(w, http.StatusOK, invoice)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("api request")
		}(time.Now())
		if status, err := fn(w, r); err != nil {
			writeJSON(w, status, map[string]string{"error": err.Error()})
		}
	}
}
