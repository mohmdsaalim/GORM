package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Test struct {
	ID int `gorm:"primaryKey"`
	Name string
	Email string `gorm:"unique"`
}

// connecting func
func setupDB() *gorm.DB {
	dsn := "host=localhost user=saalim password=949535 dbname=gintest port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil{
		panic("Falied to connect database")
	}
	fmt.Println("Database Connected")
	db.AutoMigrate(&Test{})

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour) // tommarrowwwww!!!
	return db
}

func main() {
	db := setupDB()
	r := gin.Default()
	// showing the data of users
	var users []Test
	r.GET("/users", func(ctx *gin.Context) {
		db.Find(&users)
		ctx.JSON(200, users)
	})

	r.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var user Test
		if err := db.First(&user, id).Error; err != nil{
			ctx.JSON(404, gin.H{"error":"User Not Found"})
			return
		}
		ctx.JSON(200, user)
	})

	r.POST("/users", func(ctx *gin.Context) {
		var input Test
		if err := ctx.ShouldBindJSON(&input);err != nil{
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user := Test{Name: input.Name, Email: input.Email}
		if err := db.Create(&user).Error; err != nil{
			ctx.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}
		ctx.JSON(201, user)
	})

	r.PUT("user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
			var user Test
			if err := db.First(&user, id).Error; err != nil{
				ctx.JSON(400, gin.H{"error":err.Error()})
				return
			}

			var input Test
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
			db.Model(&user).Updates(Test{Name: input.Name, Email: input.Email})
			ctx.JSON(200, user)
	})

	r.DELETE("user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var user Test
		if err := db.First(&user, id).Error; err != nil{
			ctx.JSON(404, gin.H{
				"error":"user not found",
			})
			return
		}
		db.Delete(&user)
		ctx.JSON(200, gin.H{
			"message":"User deleted",
		})
	})
	r.Run(":8080")
}