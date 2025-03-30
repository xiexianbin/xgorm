package main

import (
	"context"
	"fmt"
	"log"

	"go.xiexianbin.cn/xgorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "gorm.io/driver/mysql"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255"`
	Email    string `gorm:"size:255;uniqueIndex"`
	Password string `gorm:"size:255"`
}

func main() {
	// 初始化Gorm
	// dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})

	// 创建Repository
	userRepo := xgorm.NewRepository[User](db)

	ctx := context.Background()

	// 创建用户
	user := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepassword",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new user: %+v\n", user)

	// 查询用户
	foundUser, err := userRepo.FindByID(ctx, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found user: %+v\n", foundUser)

	// 更新用户
	foundUser.Name = "John Updated"
	if err := userRepo.Update(ctx, foundUser); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Update user: %+v\n", foundUser)

	// 条件查询
	users, err := userRepo.FindByCondition(ctx, "name LIKE ?", "%John%")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Users with name containing 'John':")
	for _, _user := range users {
		fmt.Printf("  %+v\n", _user)
	}

	// 事务示例
	err = userRepo.Transaction(ctx, func(txRepo xgorm.IRepository[User]) error {
		// create new user in tx
		newUser := &User{
			Name:     "Transaction User",
			Email:    "transaction@example.com",
			Password: "password",
		}
		if err := txRepo.Create(ctx, newUser); err != nil {
			return err
		}
		fmt.Printf("tx new user: %+v\n", newUser)

		// 更新第一个用户
		foundUser.Name = "Updated in transaction"
		return txRepo.Update(ctx, foundUser)
	})
	if err != nil {
		log.Fatal("Transaction failed:", err)
	}
}
