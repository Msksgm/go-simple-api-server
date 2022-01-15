package server

import "os"

func (s *Server) routes() {
	s.router.Use(Logger(os.Stdout))
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	apiRouter.Handle("/health", healthCheck())
	apiRouter.Handle("/users", s.createUser()).Methods("POST")
	apiRouter.Handle("/users/login", s.loginUser()).Methods("POST")
}
