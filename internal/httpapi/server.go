package httpapi

import (
	"github.com/jb843051627/indigo-vat/internal/service"
	"net/http"
)

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func New(s *service.Service) *Server {
	server := &Server{service: s, mux: http.NewServeMux()}
	server.routes()
	return server
}
func (s *Server) Handler() http.Handler { return logging(s.mux) }
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Indigo-Vat", "web")
		next.ServeHTTP(w, r)
	})
}
