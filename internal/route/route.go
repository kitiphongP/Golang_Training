package route

import (
	"net/http"
	"golang/internal/repository"
	"golang/internal/service"
	"golang/internal/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	repo := repository.NewUserRepository()
	service := service.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	mux.HandleFunc("/users", handler.GetUserHandler)
	mux.HandleFunc("/users/create", handler.CreateUserHandler)
	mux.HandleFunc("/users/search", handler.SearchUserHandler)

	return mux
}