# Test Chromedp Bill Printer

## Mô tả

Tool test in bill bằng Chromedp (HTML → Ảnh → ESC/POS).

## Chạy test

```bash
cd backend/cmd/test-chromedp-print
go run main.go
```

## Kết quả

- Tạo file `test_bill_escpos.bin` chứa ESC/POS commands
- In ra console số bytes đã generate
- (Optional) Gửi trực tiếp đến máy in nếu uncomment code

## Gửi đến máy in

Uncomment dòng này trong `main.go`:

```go
printerIP := "192.168.1.100:9100"
if err := services.SendToPrinter(printerIP, escposData); err != nil {
    log.Printf("Failed: %v", err)
}
```

Thay `192.168.1.100` bằng IP máy in của bạn.

## Xem preview

Để xem ảnh preview trước khi in:

```go
renderer.SavePreviewImage(ord, shopSettings, "preview.png")
```

## Dependencies

```bash
go get github.com/chromedp/chromedp
```

Chromedp sẽ tự động download Chrome binary lần đầu chạy.
