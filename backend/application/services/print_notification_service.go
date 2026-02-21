package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"cafe-pos/backend/domain/printing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrintNotificationService handles print-related notifications
type PrintNotificationService interface {
	// NotifyPrintFailure creates a notification for a failed print job
	NotifyPrintFailure(job *printing.PrintJob, printerName string, err error)

	// NotifyPrinterOffline creates a notification for an offline printer
	NotifyPrinterOffline(printerID primitive.ObjectID, printerName string)

	// NotifyHardwareError creates a notification for hardware errors
	NotifyHardwareError(printerID primitive.ObjectID, printerName string, errorType printing.NotificationType, errorMsg string)

	// GetUnreadNotifications returns all unread notifications
	GetUnreadNotifications() ([]*printing.PrintNotification, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(id primitive.ObjectID) error

	// MarkAllAsRead marks all notifications as read
	MarkAllAsRead() error
}

type printNotificationService struct {
	notificationRepo printing.PrintNotificationRepository
}

// NewPrintNotificationService creates a new notification service
func NewPrintNotificationService(notificationRepo printing.PrintNotificationRepository) PrintNotificationService {
	return &printNotificationService{
		notificationRepo: notificationRepo,
	}
}

// NotifyPrintFailure creates a notification for a failed print job
func (s *printNotificationService) NotifyPrintFailure(job *printing.PrintJob, printerName string, err error) {
	if job == nil {
		return
	}

	// Determine notification type based on error message
	notifType := printing.NotificationTypePrintFailure
	errorMsg := err.Error()

	// Check for specific error types
	if strings.Contains(strings.ToLower(errorMsg), "connection") ||
		strings.Contains(strings.ToLower(errorMsg), "offline") ||
		strings.Contains(strings.ToLower(errorMsg), "unreachable") {
		notifType = printing.NotificationTypePrinterOffline
	}

	notification := &printing.PrintNotification{
		Type:        notifType,
		Severity:    printing.NotificationSeverityError,
		Title:       s.getNotificationTitle(notifType, job.Type),
		Message:     s.getNotificationMessage(notifType, job, printerName),
		JobID:       job.ID,
		OrderID:     job.OrderID,
		OrderNumber: job.OrderNumber,
		PrinterID:   job.PrinterID,
		PrinterName: printerName,
		ErrorMsg:    errorMsg,
		Read:        false,
		CreatedAt:   time.Now(),
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		log.Printf("[NOTIFICATION ERROR] Failed to create notification: %v", err)
	} else {
		log.Printf("[NOTIFICATION] Created notification - type=%s, job_id=%s, order=%s, printer=%s",
			notifType, job.ID.Hex(), job.OrderNumber, printerName)
	}
}

// NotifyPrinterOffline creates a notification for an offline printer
func (s *printNotificationService) NotifyPrinterOffline(printerID primitive.ObjectID, printerName string) {
	notification := &printing.PrintNotification{
		Type:        printing.NotificationTypePrinterOffline,
		Severity:    printing.NotificationSeverityWarning,
		Title:       "Máy in không kết nối được",
		Message:     fmt.Sprintf("Máy in '%s' đang offline hoặc không thể kết nối. Vui lòng kiểm tra kết nối mạng và trạng thái máy in.", printerName),
		PrinterID:   printerID,
		PrinterName: printerName,
		Read:        false,
		CreatedAt:   time.Now(),
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		log.Printf("[NOTIFICATION ERROR] Failed to create printer offline notification: %v", err)
	} else {
		log.Printf("[NOTIFICATION] Printer offline - printer=%s", printerName)
	}
}

// NotifyHardwareError creates a notification for hardware errors
func (s *printNotificationService) NotifyHardwareError(printerID primitive.ObjectID, printerName string, errorType printing.NotificationType, errorMsg string) {
	var title, message string
	var severity printing.NotificationSeverity

	switch errorType {
	case printing.NotificationTypePaperOut:
		title = "Máy in hết giấy"
		message = fmt.Sprintf("Máy in '%s' đã hết giấy. Vui lòng thêm giấy và thử lại.", printerName)
		severity = printing.NotificationSeverityWarning
	case printing.NotificationTypePaperJam:
		title = "Máy in bị kẹt giấy"
		message = fmt.Sprintf("Máy in '%s' bị kẹt giấy. Vui lòng kiểm tra và xử lý.", printerName)
		severity = printing.NotificationSeverityError
	case printing.NotificationTypeCoverOpen:
		title = "Nắp máy in đang mở"
		message = fmt.Sprintf("Nắp máy in '%s' đang mở. Vui lòng đóng nắp và thử lại.", printerName)
		severity = printing.NotificationSeverityWarning
	default:
		title = "Lỗi phần cứng máy in"
		message = fmt.Sprintf("Máy in '%s' gặp lỗi phần cứng: %s", printerName, errorMsg)
		severity = printing.NotificationSeverityError
	}

	notification := &printing.PrintNotification{
		Type:        errorType,
		Severity:    severity,
		Title:       title,
		Message:     message,
		PrinterID:   printerID,
		PrinterName: printerName,
		ErrorMsg:    errorMsg,
		Read:        false,
		CreatedAt:   time.Now(),
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		log.Printf("[NOTIFICATION ERROR] Failed to create hardware error notification: %v", err)
	} else {
		log.Printf("[NOTIFICATION] Hardware error - type=%s, printer=%s", errorType, printerName)
	}
}

// GetUnreadNotifications returns all unread notifications
func (s *printNotificationService) GetUnreadNotifications() ([]*printing.PrintNotification, error) {
	return s.notificationRepo.FindUnread()
}

// MarkAsRead marks a notification as read
func (s *printNotificationService) MarkAsRead(id primitive.ObjectID) error {
	return s.notificationRepo.MarkAsRead(id)
}

// MarkAllAsRead marks all notifications as read
func (s *printNotificationService) MarkAllAsRead() error {
	return s.notificationRepo.MarkAllAsRead()
}

// getNotificationTitle returns the appropriate title for a notification type
func (s *printNotificationService) getNotificationTitle(notifType printing.NotificationType, jobType printing.PrintJobType) string {
	switch notifType {
	case printing.NotificationTypePrinterOffline:
		return "Máy in không kết nối được"
	case printing.NotificationTypePrintFailure:
		if jobType == printing.PrintJobTypeBill {
			return "In bill thất bại"
		}
		return "In tem thất bại"
	default:
		return "Lỗi in ấn"
	}
}

// getNotificationMessage returns the appropriate message for a notification
func (s *printNotificationService) getNotificationMessage(notifType printing.NotificationType, job *printing.PrintJob, printerName string) string {
	jobTypeStr := "bill"
	if job.Type == printing.PrintJobTypeLabel {
		jobTypeStr = "tem"
	}

	switch notifType {
	case printing.NotificationTypePrinterOffline:
		return fmt.Sprintf("Không thể in %s cho đơn hàng %s. Máy in '%s' đang offline hoặc không thể kết nối.",
			jobTypeStr, job.OrderNumber, printerName)
	case printing.NotificationTypePrintFailure:
		return fmt.Sprintf("Không thể in %s cho đơn hàng %s trên máy in '%s'. Vui lòng kiểm tra và thử lại.",
			jobTypeStr, job.OrderNumber, printerName)
	default:
		return fmt.Sprintf("Lỗi khi in %s cho đơn hàng %s.", jobTypeStr, job.OrderNumber)
	}
}
