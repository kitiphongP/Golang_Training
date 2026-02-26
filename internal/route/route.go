package route

import (
	"golang/internal/adapter/fiber/handlers"
	"net/http"

	"golang/internal/core/repository"
	"golang/internal/core/service"
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

	mux.HandleFunc("/github/skills", handlers.GetuserSkillHandler)
	mux.HandleFunc("/github/languages", handlers.GetUserLanguagesHandler)

	return mux
}