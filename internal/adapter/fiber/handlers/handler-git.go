package handlers

import (
	"encoding/json"
	"golang/internal/core/service"
	"net/http"
)

func resolveIdentifierFromQuery(r *http.Request) (string, string, error) {
	username := r.URL.Query().Get("username")
	email := r.URL.Query().Get("email")

	if username == "" && email == "" {
		return "", "", http.ErrMissingFile
	}

	return username, email, nil
}

func GetUserLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, email, err := resolveIdentifierFromQuery(r)
	if err == http.ErrMissingFile {
		http.Error(w, "username or email is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	report, err := service.GetGitHubLanguageReportByIdentifier(username, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// ฟังก์ชันสำหรับดึงข้อมูลทักษะของผู้ใช้จาก GitHub API
func GetuserSkillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, email, err := resolveIdentifierFromQuery(r)
	if err == http.ErrMissingFile {
		http.Error(w, "username or email is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	report, err := service.GetGitHubLanguageReportByIdentifier(username, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report.Languages)
}