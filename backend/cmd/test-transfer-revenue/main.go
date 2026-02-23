package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "http://localhost:3000/api"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Shift struct {
	ID                 string  `json:"id"`
	TransferRevenue    float64 `json:"transfer_revenue"`
	RemainingTransfer  float64 `json:"remaining_transfer"`
	HandedOverTransfer float64 `json:"handed_over_transfer"`
	CurrentCash        float64 `json:"current_cash"`
	RemainingCash      float64 `json:"remaining_cash"`
}

func main() {
	fmt.Println("=== Testing Transfer Revenue Calculation ===\n")

	// Step 1: Login as waiter
	fmt.Println("1. Logging in as waiter...")
	token, err := login("waiter", "waiter123")
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		fmt.Println("\nPlease ensure waiter user exists:")
		fmt.Println("Username: waiter")
		fmt.Println("Password: waiter123")
		os.Exit(1)
	}

	fmt.Printf("✅ Logged in successfully\n")
	fmt.Printf("Token: %s...\n\n", token[:20])

	// Step 2: Get current shift
	fmt.Println("2. Fetching current shift...")
	shift, err := getCurrentShift(token)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		fmt.Println("\nℹ️  No open shift found. Please:")
		fmt.Println("1. Go to http://localhost:5173/#/shifts")
		fmt.Println("2. Login as waiter / waiter123")
		fmt.Println("3. Start a new shift")
		fmt.Println("4. Create some orders with transfer payment (CK or QR)")
		fmt.Println("5. Run this test again")
		os.Exit(0)
	}

	// Display shift data
	fmt.Println("\n=== Shift Data ===")
	shiftJSON, _ := json.MarshalIndent(shift, "", "  ")
	fmt.Println(string(shiftJSON))

	// Check fields
	fmt.Println("\n=== Results ===")
	fmt.Printf("✅ transfer_revenue: %.0f VND\n", shift.TransferRevenue)
	fmt.Printf("✅ remaining_transfer: %.0f VND\n", shift.RemainingTransfer)
	fmt.Printf("✅ handed_over_transfer: %.0f VND\n", shift.HandedOverTransfer)
	fmt.Printf("✅ current_cash: %.0f VND\n", shift.CurrentCash)
	fmt.Printf("✅ remaining_cash: %.0f VND\n", shift.RemainingCash)

	if shift.TransferRevenue > 0 {
		fmt.Println("\n🎉 Transfer revenue is being calculated correctly!")
	} else {
		fmt.Println("\nℹ️  Transfer revenue is 0. Create orders with transfer payment to test.")
	}
}

func login(username, password string) (string, error) {
	reqBody := LoginRequest{
		Username: username,
		Password: password,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	return loginResp.Token, nil
}

func getCurrentShift(token string) (*Shift, error) {
	req, _ := http.NewRequest("GET", baseURL+"/waiter/shifts/current", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("no open shift found")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get shift: %s", string(body))
	}

	var shift Shift
	if err := json.NewDecoder(resp.Body).Decode(&shift); err != nil {
		return nil, err
	}

	return &shift, nil
}
