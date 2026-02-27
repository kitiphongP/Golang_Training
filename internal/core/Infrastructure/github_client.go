package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)
// GitHubClient เป็นโครงสร้างที่ใช้สำหรับเชื่อมต่อและดึงข้อมูลจาก GitHub API
type GitHubClient struct {
	BaseURL string
	Token   string
}
// Repo เป็นโครงสร้างที่ใช้สำหรับเก็บข้อมูลของ Repository ที่ดึงมาจาก GitHub API
type Repo struct {
	Name string `json:"name"`
	Language string `json:"language"`
}
// ฟังก์ชัน GetRepos เป็นฟังก์ชันที่ใช้สำหรับดึงข้อมูล Repository ของผู้ใช้จาก GitHub API โดยรับพารามิเตอร์เป็น context และ username ของผู้ใช้
func (g *GitHubClient) GetRepos(ctx context.Context, username string) ([]Repo, error) {
	// สร้าง URL สำหรับเรียก GitHub API โดยใช้ BaseURL และ username ที่ได้รับมา
	url := fmt.Sprintf("%s/users/%s/repos", g.BaseURL, username)
	// สร้าง HTTP request โดยใช้ context ที่ได้รับมา และ URL ที่สร้างขึ้น
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	// ถ้ามี Token ให้เพิ่ม Authorization header ใน HTTP request
	if g.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.Token)
	}
	// ส่ง HTTP request ไปยัง GitHub API และรับ response กลับมา
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var repos []Repo
	err = json.NewDecoder(resp.Body).Decode(&repos)
	return repos, err
}