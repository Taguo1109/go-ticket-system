package main

/**
 * @File: main.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/18 下午8:49
 * @Software: GoLand
 * @Version:  1.0
 */

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
	CTX = context.Background()
)

func initEnv() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("⚠️  .env not found, using system env")
	}
}

func initMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ MySQL connection failed: %v", err)
	}
	DB = db
	log.Println("✅ Connected to MySQL")
}

func initRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // e.g., "localhost:6379"
		Password: os.Getenv("REDIS_PASS"), // "" if none
		DB:       0,
	})

	_, err := rdb.Ping(CTX).Result()
	if err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}
	RDB = rdb
	log.Println("✅ Connected to Redis")
}

func main() {
	initEnv()
	initMySQL()
	initRedis()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	log.Println("🚀 API Server running on :8080")
	log.Fatal(r.Run(":8080"))
}
