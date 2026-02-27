package route

import (
	"golang/internal/adapter/fiber/handlers"
	"net/http"
)

func Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/github/skills", handlers.GetuserSkillHandler)
	mux.HandleFunc("/github/languages", handlers.GetUserLanguagesHandler)

	return mux
}