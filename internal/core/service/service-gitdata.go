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
// ฟังก์ชันสำหรับปรับปรุงข้อมูล GitHubLanguageReport ให้สมบูรณ์และมีค่าเริ่มต้นที่เหมาะสม
func normalizeGitHubLanguageReport(report *models.GitHubLanguageReport) *models.GitHubLanguageReport {
	if report == nil {
		return nil
	}
	if report.Repos == nil {
		report.Repos = []models.Repo{}
	}
	// ตรวจสอบว่า Languages ใน report เป็น nil หรือไม่ หากเป็น nil ให้กำหนดค่าเป็นแผนที่ว่างเพื่อป้องกันการเกิด panic เมื่อเข้าถึง Languages ในภายหลัง
	if len(report.Languages) == 0 {
		report.Languages = AnalyzeSkills(report.Repos)
	}
	// ตรวจสอบว่า UpdatedAt ใน report เป็นค่าเริ่มต้นหรือไม่ หากเป็นค่าเริ่มต้นให้กำหนดค่าเป็นเวลาปัจจุบันเพื่อให้ข้อมูลมีความสมบูรณ์มากขึ้น
	if report.UpdatedAt.IsZero() {
		report.UpdatedAt = time.Now()
	}

	return report
}

// ฟังก์ชันสำหรับเชื่อมต่อกับ MongoDB collection ที่ใช้เก็บข้อมูล GitHub reports
func githubCollection() *mongo.Collection {
	if database.DB == nil {
		return nil
	}
	return database.DB.Collection("usergithub")
}
// ฟังก์ชันสำหรับค้นหา GitHubLanguageReport จาก MongoDB โดยใช้ filter ที่กำหนด
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
// ฟังก์ชันสำหรับบันทึก GitHubLanguageReport ลงใน MongoDB โดยใช้ ReplaceOne เพื่ออัปเดตข้อมูลหากมีอยู่แล้วหรือสร้างใหม่หากไม่มี
func saveGitHubReportToDB(report *models.GitHubLanguageReport) error {
	collection := githubCollection()
	if collection == nil {
		return fmt.Errorf("mongodb is not connected")
	}
	normalized := normalizeGitHubLanguageReport(report)
	// กำหนดค่า UpdatedAt เป็นเวลาปัจจุบัน
	normalized.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// ใช้ ReplaceOne เพื่ออัปเดตข้อมูลใน MongoDB หากมีอยู่แล้วหรือสร้างใหม่หากไม่มี โดยใช้ filter ที่ค้นหาจาก username ของ report
	filter := bson.M{"username": normalized.Username}
	_, err := collection.ReplaceOne(ctx, filter, normalized, options.Replace().SetUpsert(true))
	return err
}

//======================================User Name========================================

// // ฟังก์ชันสำหรับดึงข้อมูลโปรไฟล์ผู้ใช้จาก GitHub API
func FetchGitHubUser(username string) (*models.GitHubUser, error) {
	// ถ้า username เป็นค่าว่าง ให้ส่งคืนข้อผิดพลาดทันทีเพื่อป้องกันการเรียก API 
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	// endpoint ของ GitHub API สำหรับดึงข้อมูลโปรไฟล์ผู้ใช้
	endpoint := fmt.Sprintf("https://api.github.com/users/%s", username)
	// เรียก API เพื่อดึงข้อมูลโปรไฟล์ผู้ใช้
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch github user profile: %v", err)
	}
	defer resp.Body.Close()
	// ตรวจสอบสถานะการตอบกลับจาก API หากไม่ใช่ 200 OK ให้ส่งคืนข้อผิดพลาด
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user api returned status: %s", resp.Status)
	}
	// ตัวแปรสำหรับเก็บข้อมูลโปรไฟล์ผู้ใช้ที่ได้รับจาก API
	var user models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
// ฟังก์ชันหลักสำหรับดึงข้อมูลรายงานภาษาที่ใช้ใน repos ของผู้ใช้จาก GitHub API และจัดการกับการแคชข้อมูลใน MongoDB
func GetGitHubResponse(username string) (*models.GitHubLanguageReport, error) {
	// ตรวจสอบว่า username เป็นค่าว่างหรือไม่ 
	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("username is required")
	}
	// พยายามค้นหาข้อมูลรายงานจาก MongoDB โดยใช้ username เป็น filter
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
	// map ข้อมูลที่ได้รับจาก API ไปยัง struct GitHubLanguageReport เพื่อจัดเก็บข้อมูลรายงานภาษาที่ใช้ใน repos ของผู้ใช้
	report := &models.GitHubLanguageReport{
		Username:  user.Login,
		Name:      user.Name,
		Email:     user.Email,
		Repos:     repos,
		Languages: AnalyzeSkills(repos),
	}
	// บันทึกข้อมูลรายงานลงใน MongoDB เพื่อใช้ในการแคชข้อมูลสำหรับการเรียกครั้งถัดไป
	if err := saveGitHubReportToDB(report); err != nil {
		return nil, err
	}

	return normalizeGitHubLanguageReport(report), nil
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

	return repos, nil
}

//======================================E Mail========================================

// ฟังก์ชันสำหรับดึงข้อมูลรายงานภาษาที่ใช้ใน repos ของผู้ใช้จาก GitHub API โดยใช้ username หรือ email เป็นตัวระบุ และจัดการกับการแคชข้อมูลใน MongoDB
func GetGitHubResponseByIdentifier(username string, email string) (*models.GitHubLanguageReport, error) {
	if strings.TrimSpace(username) != "" {
		return GetGitHubResponse(username)
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

	report, err := GetGitHubResponse(resolvedUsername)
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
// ฟังก์ชันสำหรับค้นหา username ของผู้ใช้ GitHub โดยใช้ email เป็นตัวระบุผ่าน GitHub Search API
func ResolveGitHubUsernameByEmail(email string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email is required")
	}
	// endpoint ของ GitHub Search API สำหรับค้นหาผู้ใช้โดยใช้ email เป็นตัวระบุ
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