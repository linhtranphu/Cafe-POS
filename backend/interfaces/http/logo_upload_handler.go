package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cafe-pos/backend/infrastructure/mongodb"

	"github.com/gin-gonic/gin"
)

// LogoUploadHandler handles logo upload and deletion
type LogoUploadHandler struct {
	shopSettingsRepo *mongodb.ShopSettingsRepository
	uploadDir        string
}

// NewLogoUploadHandler creates a new logo upload handler
func NewLogoUploadHandler(shopSettingsRepo *mongodb.ShopSettingsRepository, uploadDir string) *LogoUploadHandler {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Warning: Failed to create upload directory %s: %v", uploadDir, err)
	}
	
	return &LogoUploadHandler{
		shopSettingsRepo: shopSettingsRepo,
		uploadDir:        uploadDir,
	}
}

// UploadLogo handles logo file upload
// POST /api/settings/logo
func (h *LogoUploadHandler) UploadLogo(c *gin.Context) {
	// Parse multipart form
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file size (max 2MB)
	const maxSize = 2 * 1024 * 1024 // 2MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kích thước file tối đa 2MB"})
		return
	}

	// Validate file format (PNG, JPG, JPEG)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chỉ hỗ trợ định dạng PNG, JPG, JPEG"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("logo_%d%s", c.GetInt64("timestamp"), ext)
	if c.GetInt64("timestamp") == 0 {
		// Fallback to current timestamp if not set by middleware
		filename = fmt.Sprintf("logo_%d%s", os.Getpid(), ext)
	}
	
	filepath := filepath.Join(h.uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		log.Printf("Failed to save logo file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu file"})
		return
	}

	// Update shop settings with logo URL
	shopSettings, err := h.shopSettingsRepo.GetSettings(c.Request.Context())
	if err != nil {
		log.Printf("Failed to get shop settings: %v", err)
		// Clean up uploaded file
		os.Remove(filepath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật cài đặt"})
		return
	}

	// Delete old logo file if exists
	if shopSettings.LogoURL != "" {
		oldPath := shopSettings.LogoURL
		if strings.HasPrefix(oldPath, "/uploads/") {
			oldPath = "." + oldPath // Convert to relative path
		}
		if _, err := os.Stat(oldPath); err == nil {
			os.Remove(oldPath)
		}
	}

	// Update logo URL (store as relative path for serving)
	logoURL := "/uploads/logos/" + filename
	shopSettings.UpdatePrintSettings(
		shopSettings.ShopAddress,
		shopSettings.ShopPhone,
		logoURL,
		shopSettings.CustomMessage,
	)

	if err := h.shopSettingsRepo.UpdateSettings(c.Request.Context(), shopSettings); err != nil {
		log.Printf("Failed to update shop settings: %v", err)
		// Clean up uploaded file
		os.Remove(filepath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật cài đặt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logo_url": logoURL,
		"message":  "Upload logo thành công",
	})
}

// DeleteLogo handles logo deletion
// DELETE /api/settings/logo
func (h *LogoUploadHandler) DeleteLogo(c *gin.Context) {
	// Get shop settings
	shopSettings, err := h.shopSettingsRepo.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy cài đặt"})
		return
	}

	// Delete logo file if exists
	if shopSettings.LogoURL != "" {
		logoPath := shopSettings.LogoURL
		if strings.HasPrefix(logoPath, "/uploads/") {
			logoPath = "." + logoPath // Convert to relative path
		}
		if _, err := os.Stat(logoPath); err == nil {
			if err := os.Remove(logoPath); err != nil {
				log.Printf("Failed to delete logo file: %v", err)
			}
		}
	}

	// Update shop settings to remove logo URL
	shopSettings.UpdatePrintSettings(
		shopSettings.ShopAddress,
		shopSettings.ShopPhone,
		"", // Clear logo URL
		shopSettings.CustomMessage,
	)

	if err := h.shopSettingsRepo.UpdateSettings(c.Request.Context(), shopSettings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật cài đặt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đã xóa logo thành công",
	})
}
