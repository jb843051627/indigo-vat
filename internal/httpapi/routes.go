package httpapi

import "net/http"

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /api/v1/vats", s.vats)
	s.mux.HandleFunc("POST /api/v1/vats", s.createVat)
	s.mux.HandleFunc("GET /api/v1/cycles", s.cycles)
	s.mux.HandleFunc("POST /api/v1/cycles", s.createCycle)
	s.mux.HandleFunc("GET /api/v1/cycles/{id}/report", s.report)
	s.mux.HandleFunc("GET /", s.index)
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }
