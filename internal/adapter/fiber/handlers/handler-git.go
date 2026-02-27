package handlers

import (
	"encoding/json"
	"golang/internal/core/service"
	"net/http"
)

// ฟังก์ชัน Resolve(แปลง) จาก query parameters (username หรือ email)
func resolveIdentifierFromQuery(r *http.Request) (string, string, error) {
	username := r.URL.Query().Get("username")
	email := r.URL.Query().Get("email")

	if username == "" && email == "" {
		return "", "", http.ErrMissingFile
	}

	return username, email, nil
}
// ฟังก์ชันสำหรับดึงข้อมูล Repositories ของผู้ใช้จาก GitHub API
func GetUserLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// เรียกใช้ฟังก์ชัน resolveIdentifierFromQuery เพื่อดึงข้อมูล username หรือ email จาก query parameters
	username, email, err := resolveIdentifierFromQuery(r)

	// ถ้าไม่มีทั้ง username และ email จะส่ง error กลับไปยัง client
	if err == http.ErrMissingFile {
		http.Error(w, "username or email is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// เรียกใช้ฟังก์ชัน GetGitHubLanguageReportByIdentifier เพื่อดึงข้อมูลภาษาที่ผู้ใช้ถนัดจาก GitHub API
	report, err := service.GetGitHubLanguageReportByIdentifier(username, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// ฟังก์ชันสำหรับดึงข้อมูลจำนวนภาษาของผู้ใช้จาก GitHub API
func GetuserSkillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// เรียกใช้ฟังก์ชัน resolveIdentifierFromQuery เพื่อดึงข้อมูล username หรือ email จาก query parameters
	username, email, err := resolveIdentifierFromQuery(r)
	// ถ้าไม่มีทั้ง username และ email จะส่ง error กลับไปยัง client
	if err == http.ErrMissingFile {
		http.Error(w, "username or email is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// เรียกใช้ฟังก์ชัน GetGitHubLanguageReportByIdentifier เพื่อดึงข้อมูลภาษาที่ผู้ใช้ถนัดจาก GitHub API
	report, err := service.GetGitHubLanguageReportByIdentifier(username, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report.Languages)
}