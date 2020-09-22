package controllers

import "github.com/Muh-Sidik/BlogSimpel-Gomux/api/middleware"

func (s *Server) initializeRoutes() {

	s.Router.HandleFunc("/", middleware.SetMiddlewareJSON(s.Home)).Methods("GET")

	//LOgin
	s.Router.HandleFunc("/login", middleware.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users

	s.Router.HandleFunc("/users", middleware.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.CreateUser))).Methods("POST")
	s.Router.HandleFunc("/users/{id}", middleware.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middleware.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Post

	s.Router.HandleFunc("/posts", middleware.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/posts", middleware.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	s.Router.HandleFunc("/posts/{id}", middleware.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middleware.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

}
