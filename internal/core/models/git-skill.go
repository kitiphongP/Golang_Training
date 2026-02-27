package models

import "time"

// UserSkills เป็นโครงสร้างที่ใช้สำหรับเก็บข้อมูลทักษะของผู้ใช้ โดยมีฟิลด์ Username สำหรับเก็บชื่อผู้ใช้, Skills สำหรับเก็บแผนที่ของทักษะและคะแนน และ LastUpdated สำหรับเก็บเวลาที่ข้อมูลถูกอัปเดตล่าสุด
type UserSkills struct {
	Username string 	`bson:"username" json:"username"`
	Skills   map[string]float64 `bson:"skills" json:"skills"`
	LastUpdated time.Time 	`bson:"last_updated" json:"last_updated"`
}