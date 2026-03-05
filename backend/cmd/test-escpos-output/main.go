package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("=== Test ESC/POS Output ===")
	
	// Read the binary file
	data, err := os.ReadFile("test_uploaded_logo.bin")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	
	log.Printf("Binary file size: %d bytes", len(data))
	
	// Show first 100 bytes in hex
	log.Println("\nFirst 100 bytes (hex):")
	for i := 0; i < 100 && i < len(data); i++ {
		if i%16 == 0 {
			fmt.Printf("\n%04x: ", i)
		}
		fmt.Printf("%02x ", data[i])
	}
	fmt.Println()
	
	// Check for ESC/POS commands
	log.Println("\nChecking for ESC/POS commands:")
	
	// ESC @ (Initialize)
	if len(data) >= 2 && data[0] == 0x1B && data[1] == 0x40 {
		log.Println("✅ Found ESC @ (Initialize) at start")
	} else {
		log.Println("❌ Missing ESC @ (Initialize) at start")
	}
	
	// Look for GS commands
	gsCount := 0
	for i := 0; i < len(data)-1; i++ {
		if data[i] == 0x1D {
			gsCount++
		}
	}
	log.Printf("Found %d GS commands", gsCount)
	
	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(data)
	log.Printf("\nBase64 encoded size: %d bytes", len(encoded))
	log.Printf("Base64 preview (first 100 chars): %s...", encoded[:100])
	
	// Test decode
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}
	
	if len(decoded) == len(data) {
		log.Println("✅ Base64 encode/decode works correctly")
	} else {
		log.Printf("❌ Size mismatch after decode: %d vs %d", len(decoded), len(data))
	}
	
	// Save base64 to file
	base64File := "test_uploaded_logo_base64.txt"
	if err := os.WriteFile(base64File, []byte(encoded), 0644); err != nil {
		log.Printf("Failed to save base64: %v", err)
	} else {
		log.Printf("✅ Saved base64 to %s", base64File)
	}
	
	log.Println("\n=== Test completed ===")
}
