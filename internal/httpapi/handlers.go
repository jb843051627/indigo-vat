package httpapi

import (
	"errors"
	"github.com/jb843051627/indigo-vat/internal/model"
	"github.com/jb843051627/indigo-vat/internal/service"
	"net/http"
)

func (s *Server) vats(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.ListVats(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, items)
}
func (s *Server) createVat(w http.ResponseWriter, r *http.Request) {
	var in model.VatInput
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	item, err := s.service.CreateVat(r.Context(), in)
	if err != nil {
		writeJSON(w, 422, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, item)
}
func (s *Server) cycles(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.ListCycles(r.Context(), r.URL.Query().Get("state"))
	if err != nil {
		writeJSON(w, 422, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, items)
}
func (s *Server) createCycle(w http.ResponseWriter, r *http.Request) {
	var in model.CycleInput
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	item, err := s.service.StartCycle(r.Context(), in)
	if err != nil {
		if errors.Is(err, service.ErrVatNotActive) {
			writeJSON(w, 409, map[string]string{"error": err.Error(), "code": "vat_not_active"})
			return
		}
		writeJSON(w, 422, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, item)
}
func (s *Server) report(w http.ResponseWriter, r *http.Request) {
	item, err := s.service.BuildReport(r.Context(), r.PathValue("id"))
	if err != nil {
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 404, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, item)
}
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte("<!doctype html><html><head><title>indigo-vat</title></head><body><h1>Indigo Vat Operations</h1><p>Fermentation cycles and release dossiers</p></body></html>"))
}
