package database

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)
// ตัวแปร DB เป็นตัวเก็บการเชื่อมต่อกับฐานข้อมูล MongoDB
var DB *mongo.Database

func ConnectMongoDB(uri string, dbName string){
	// สร้าง context พร้อม timeout เพื่อใช้ในการเชื่อมต่อกับ MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// เชื่อมต่อกับ MongoDB โดยใช้ URI ที่กำหนด
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Mongo connection error:", err)
	}

	DB = client.Database(dbName)
	log.Println("MongoDB connected 🚀")
}