package server

import (
	"net/http"

	"meme-generator/internal/auth"
	"meme-generator/internal/controller"
)

func (s *Server) registerRoutes() {
	apiGroup := s.router.Group("/api")

	// unauthorized routes
	unauthedApiGroup := apiGroup.Group("/")

	// auth
	unauthedApiGroup.Handle(http.MethodPost, "/register", controller.Register(s.store))
	unauthedApiGroup.Handle(http.MethodPost, "/login", controller.Login(s.store))

	// authorized routes
	authedApiGroup := apiGroup.Group("/")
	authedApiGroup.Use(auth.AuthRequired)

	// auth
	authedApiGroup.Handle(http.MethodPost, "/logout", controller.Logout(s.store))

	// user
	authedApiGroup.Handle(http.MethodGet, "/users", controller.ListUsers(s.store))
	authedApiGroup.Handle(http.MethodGet, "/user/me", controller.GetUserSelf(s.store))
	authedApiGroup.Handle(http.MethodGet, "/user/:id", controller.GetUser(s.store))
	authedApiGroup.Handle(http.MethodGet, "/user/:id/templates", controller.ListUserTemplates(s.store))
	authedApiGroup.Handle(http.MethodGet, "/user/:id/memes", controller.ListUserMemes(s.store))

	// template
	authedApiGroup.Handle(http.MethodPost, "/template", controller.CreateTemplate(s.store, s.memeManager))
	authedApiGroup.Handle(http.MethodGet, "/template/:id", controller.GetTemplate(s.store, s.memeManager))
	authedApiGroup.Handle(http.MethodGet, "/template/:id/image", controller.GetTemplateImage(s.store, s.memeManager))
	authedApiGroup.Handle(http.MethodGet, "/templates/last", controller.GetLastTemplates(s.store, s.memeManager))

	// meme
	authedApiGroup.Handle(http.MethodPost, "/meme", controller.CreateMeme(s.store, s.memeManager))
	authedApiGroup.Handle(http.MethodGet, "/meme/:id", controller.GetMeme(s.store, s.memeManager))
	authedApiGroup.Handle(http.MethodGet, "/meme/:id/image", controller.GetMemeImage(s.store, s.memeManager))
	authedApiGroup.Handle(http.MethodGet, "/memes/last", controller.GetLastMemes(s.store, s.memeManager))
	authedApiGroup.Handle(http.MethodPost, "/meme/preview", controller.PreviewMeme(s.store, s.memeManager))
}
