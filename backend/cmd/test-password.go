package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	hash := "$2a$14$rdkCc7VOwAJ115Qp7LJzp.0eBLX27d4LN5.CgL9YC6NReFU0.bwey"
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("❌ Password does NOT match:", err)
	} else {
		fmt.Println("✅ Password matches!")
	}
}
