package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ---------------- MODELS ----------------

// User → users table
type User struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string
	Email  string  `gorm:"unique"`
	Orders []Order // One-to-Many relationship
}

// Product → products table
type Product struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string
	Price float64
	Orders []Order `gorm:"many2many:order_products"` // Many-to-Many through join table
}

// Order → orders table
type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User
	Products  []Product `gorm:"many2many:order_products"` // Many-to-Many
	Quantity  int
}

func main() {
	// ---------------- DATABASE CONNECTION ----------------
	dsn := "host=localhost user=saalim password=949535 dbname=footballers port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Database connected ✅")

	// ---------------- AUTO MIGRATE ----------------
	err = db.AutoMigrate(&User{}, &Product{}, &Order{})
	if err != nil {
		panic(err)
	}
	fmt.Println("AutoMigrate completed ✅")

	// ---------------- SAMPLE DATA ----------------
	user := User{Name: "Lamine Yamal", Email: "lamine@example.com"}
	db.Create(&user)

	product := Product{Name: "Football Shoe", Price: 99.99}
	db.Create(&product)

	order := Order{
		UserID:   user.ID,
		Products: []Product{product}, // Many-to-Many
		Quantity: 2,
	}
	db.Create(&order)

	// ---------------- FETCH DATA ----------------
	var users []User
	db.Preload("Orders.Products").Find(&users)

	fmt.Println("\nUsers and their Orders:")
	for _, u := range users {
		fmt.Printf("User: %d | %s | %s\n", u.ID, u.Name, u.Email)
		for _, o := range u.Orders {
			fmt.Printf("  OrderID: %d | Quantity: %d | Products: ", o.ID, o.Quantity)
			for _, p := range o.Products {
				fmt.Printf("%s ", p.Name)
			}
			fmt.Println()
		}
	}

	var products []Product
	db.Find(&products)
	fmt.Println("\nProducts:")
	for _, p := range products {
		fmt.Printf("Product: %d | %s | Price: %.2f\n", p.ID, p.Name, p.Price)
	}
}