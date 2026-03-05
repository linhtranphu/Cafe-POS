#!/bin/bash

# Function to update layout in a file
update_file() {
    local file=$1
    
    # Update address text
    sed -i '' 's/"address":   "Đ\/c: 10\/8 Trần Nhật Duật, Tân Định, Q1"/"address":   "Đ\/c: 10\/8 Trần Nhật Duật, P. Tân Định, Quận 1"/' "$file"
    
    # Replace the logo and address section
    sed -i '' '/logoPath := /,/y += 30/c\
	logoPath := "../../../backend/uploads/logos/logo_24094.jpeg"\
	logo, err := loadImage(logoPath)\
	var logoH int\
	if err == nil {\
		resizedLogo := resize.Resize(150, 0, logo, resize.Lanczos3)\
		logoW := resizedLogo.Bounds().Dx()\
		logoH = resizedLogo.Bounds().Dy()\
		\
		// Vẽ logo bên trái\
		dc.DrawImage(resizedLogo, Margin, int(y))\
	}\
\
	// Vẽ địa chỉ bên phải logo\
	textX := float64(Margin + 160)\
	dc.LoadFontFace(fontPath, 16)\
	// Wrap địa chỉ nếu cần\
	maxWidth := float64(ImageWidthPixels - Margin - 170)\
	addressLines := wrapText(dc, data["address"], maxWidth)\
	textY := y + 20.0\
	for i, line := range addressLines {\
		dc.DrawString(line, textX, textY+float64(i*18))\
	}\
\
	// Cập nhật y sau logo\
	if logoH > 60 {\
		y += float64(logoH) + 15\
	} else {\
		y += 75\
	}
' "$file"
}

# Update both files
update_file "main.go"
update_file "preview.go"

echo "✓ Updated layout in main.go and preview.go"
