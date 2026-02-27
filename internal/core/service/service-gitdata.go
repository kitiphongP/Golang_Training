package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang/internal/adapter/storage/database"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang/internal/core/models"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func githubCollection() *mongo.Collection {
	if database.DB == nil {
		return nil
	}
	return database.DB.Collection("usergithub")
}

func normalizeGitHubLanguageReport(report *models.GitHubLanguageReport) *models.GitHubLanguageReport {
	if report == nil {
		return nil
	}

	if report.Repos == nil {
		report.Repos = []models.Repo{}
	}

	if report.Languages == nil || len(report.Languages) == 0 {
		report.Languages = AnalyzeSkills(report.Repos)
	}

	if report.UpdatedAt.IsZero() {
		report.UpdatedAt = time.Now()
	}

	return report
}

func findGitHubReportFromDB(filter bson.M) (*models.GitHubLanguageReport, error) {
	collection := githubCollection()
	if collection == nil {
		return nil, fmt.Errorf("mongodb is not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var report models.GitHubLanguageReport
	err := collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		return nil, err
	}

	return normalizeGitHubLanguageReport(&report), nil
}

func saveGitHubReportToDB(report *models.GitHubLanguageReport) error {
	collection := githubCollection()
	if collection == nil {
		return fmt.Errorf("mongodb is not connected")
	}

	normalized := normalizeGitHubLanguageReport(report)
	normalized.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": normalized.Username}
	_, err := collection.ReplaceOne(ctx, filter, normalized, options.Replace().SetUpsert(true))
	return err
}

// ฟังก์ชันสำหรับดึงข้อมูล repos ของผู้ใช้จาก GitHub API
func FetchGitHubRepos(username string) ([]models.Repo, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	// endpoint ของ GitHub API สำหรับดึงข้อมูล repos ของผู้ใช้
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)

	// เรียก API เพื่อดึงข้อมูล repos
	resp, err := http.Get(url)
	// ตรวจสอบข้อผิดพลาดจากการเรียก API
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %v", err)
	}
	// ปิดการเชื่อมต่อหลังจากใช้งานเสร็จ
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status: %s", resp.Status)
	}

	// ตัวแปรสำหรับเก็บข้อมูล repos ที่ได้รับจาก API
	var repos []models.Repo

	// แปลงข้อมูล JSON ที่ได้รับจาก API เป็น struct Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	// ส่งคืนข้อมูล repos ที่ได้รับจาก API
	return repos, nil
}

func GetGitHubLanguageSummary(username string) (map[string]int, error) {
	repos, err := FetchGitHubRepos(username)
	if err != nil {
		return nil, err
	}

	return AnalyzeSkills(repos), nil
}

func FetchGitHubUser(username string) (*models.GitHubUser, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	endpoint := fmt.Sprintf("https://api.github.com/users/%s", username)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch github user profile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user api returned status: %s", resp.Status)
	}

	var user models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetGitHubLanguageReport(username string) (*models.GitHubLanguageReport, error) {
	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("username is required")
	}

	if cached, err := findGitHubReportFromDB(bson.M{"username": username}); err == nil {
		return cached, nil
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	user, err := FetchGitHubUser(username)
	if err != nil {
		return nil, err
	}

	repos, err := FetchGitHubRepos(user.Login)
	if err != nil {
		return nil, err
	}

	report := &models.GitHubLanguageReport{
		Username:  user.Login,
		Name:      user.Name,
		Email:     user.Email,
		Repos:     repos,
		Languages: AnalyzeSkills(repos),
	}

	if err := saveGitHubReportToDB(report); err != nil {
		return nil, err
	}

	return normalizeGitHubLanguageReport(report), nil
}

func GetGitHubLanguageReportByIdentifier(username string, email string) (*models.GitHubLanguageReport, error) {
	if strings.TrimSpace(username) != "" {
		return GetGitHubLanguageReport(username)
	}

	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("username or email is required")
	}

	if cached, err := findGitHubReportFromDB(bson.M{"email": email}); err == nil {
		return cached, nil
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	resolvedUsername, err := ResolveGitHubUsernameByEmail(email)
	if err != nil {
		return nil, err
	}

	report, err := GetGitHubLanguageReport(resolvedUsername)
	if err != nil {
		return nil, err
	}

	if report.Email == "" {
		report.Email = email
		if err := saveGitHubReportToDB(report); err != nil {
			return nil, err
		}
	}

	return normalizeGitHubLanguageReport(report), nil
}

func ResolveGitHubUsernameByEmail(email string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email is required")
	}

	endpoint := fmt.Sprintf("https://api.github.com/search/users?q=%s+in:email", url.QueryEscape(email))
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to search github user by email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github search api returned status: %s", resp.Status)
	}

	var result models.GitHubUserSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Items) == 0 || result.Items[0].Login == "" {
		return "", fmt.Errorf("no github account found for this public email")
	}

	return result.Items[0].Login, nil
}

func FetchGitHubReposByEmail(email string) (string, []models.Repo, error) {
	username, err := ResolveGitHubUsernameByEmail(email)
	if err != nil {
		return "", nil, err
	}

	repos, err := FetchGitHubRepos(username)
	if err != nil {
		return "", nil, err
	}

	return username, repos, nil
}

func GetGitHubLanguageSummaryByEmail(email string) (string, map[string]int, error) {
	username, repos, err := FetchGitHubReposByEmail(email)
	if err != nil {
		return "", nil, err
	}

	return username, AnalyzeSkills(repos), nil
}

// ฟังก์ชันสำหรับวิเคราะห์ทักษะจาก repos ที่ได้รับ
func AnalyzeSkills(repos []models.Repo) map[string]int {
	// สร้างแผนที่เพื่อเก็บจำนวนครั้งที่แต่ละภาษาโปรแกรมถูกใช้ใน repos
	skills := make(map[string]int)

	// วนลูปผ่าน repos และนับจำนวนครั้งที่แต่ละภาษาโปรแกรมถูกใช้
	for _, repo := range repos {
		if repo.Language != "" {
			skills[repo.Language]++
		}
	}

	return skills
}