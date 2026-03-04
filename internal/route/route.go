package route

import (
	"golang/internal/adapter/fiber/handlers"
	"net/http"
)

func Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/github/skills", handlers.GetuserSkillHandler)
	mux.HandleFunc("/api/github/languages", handlers.GetUserLanguagesHandler)

	return mux
}