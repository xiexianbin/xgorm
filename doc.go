/*
Package xgorm implements the GORM framework's generic CRUD (add, delete, change, and retrieve)
operations based on Go 1.18+'s generic features.

The simplest way to use xgorm:

	package main

	import (
		"context"
		"fmt"
		"log"

		"go.xiexianbin.cn/xgorm"
		"gorm.io/driver/sqlite"
		"gorm.io/gorm"
	)

	type User struct {
		gorm.Model
		Name     string `gorm:"size:255"`
		Email    string `gorm:"size:255;uniqueIndex"`
		Password string `gorm:"size:255"`
	}

	func main() {
		db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		db.AutoMigrate(&User{})

		userRepo := xgorm.NewRepository[User](db)

		ctx := context.Background()

		user := &User{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "securepassword",
		}
		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("new user: %+v\n", user)

		foundUser, err := userRepo.FindByID(ctx, user.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found user: %+v\n", foundUser)

		foundUser.Name = "John Updated"
		if err := userRepo.Update(ctx, foundUser); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Update user: %+v\n", foundUser)

		users, err := userRepo.FindByCondition(ctx, "name LIKE ?", "%John%")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Users with name containing 'John':")
		for _, _user := range users {
			fmt.Printf("  %+v\n", _user)
		}

		err = userRepo.Transaction(ctx, func(txRepo xgorm.IRepository[User]) error {
			newUser := &User{
				Name:     "Transaction User",
				Email:    "transaction@example.com",
				Password: "password",
			}
			if err := txRepo.Create(ctx, newUser); err != nil {
				return err
			}
			fmt.Printf("tx new user: %+v\n", newUser)

			foundUser.Name = "Updated in transaction"
			return txRepo.Update(ctx, foundUser)
		})
		if err != nil {
			log.Fatal("Transaction failed:", err)
		}
	}

For a full guide visit https://github.com/xiexianbin/xgorm
*/
package xgorm
