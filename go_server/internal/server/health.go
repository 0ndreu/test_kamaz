package server

import "net/http"

func (s *server) CheckLiveness(w http.ResponseWriter, _ *http.Request) {
	err := s.svc.IsAlive(s.ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
