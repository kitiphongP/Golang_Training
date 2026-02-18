package route

import (
	"golang/internal/handlers"
	"net/http"

	"golang/internal/repository"
	"golang/internal/service"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	repo := repository.NewUserRepository()
	userService := service.NewUserService(repo)
	handler := handlers.NewUserHandler(userService)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("search") != "" || r.URL.Query().Get("keyword") != "" {
				handler.SearchUserHandler(w, r)
				return
			}
			handler.GetUserHandler(w, r)
		case http.MethodPost:
			handler.CreateUserHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}