package settings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewShopSettings(t *testing.T) {
	shopName := "Test Cafe"
	settings := NewShopSettings(shopName)

	assert.NotNil(t, settings)
	assert.Equal(t, shopName, settings.ShopName)
	assert.Equal(t, 20.0, settings.LowMarginThreshold, "Default threshold should be 20.0")
	assert.False(t, settings.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, settings.UpdatedAt.IsZero(), "UpdatedAt should be set")
}

func TestUpdateLowMarginThreshold_ValidThreshold(t *testing.T) {
	settings := NewShopSettings("Test Cafe")
	originalUpdatedAt := settings.UpdatedAt

	// Wait a bit to ensure timestamp changes
	time.Sleep(10 * time.Millisecond)

	err := settings.UpdateLowMarginThreshold(25.0)

	assert.NoError(t, err)
	assert.Equal(t, 25.0, settings.LowMarginThreshold)
	assert.True(t, settings.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be updated")
}

func TestUpdateLowMarginThreshold_ZeroThreshold(t *testing.T) {
	settings := NewShopSettings("Test Cafe")

	err := settings.UpdateLowMarginThreshold(0.0)

	assert.NoError(t, err, "Zero threshold should be valid")
	assert.Equal(t, 0.0, settings.LowMarginThreshold)
}

func TestUpdateLowMarginThreshold_NegativeThreshold(t *testing.T) {
	settings := NewShopSettings("Test Cafe")
	originalThreshold := settings.LowMarginThreshold

	err := settings.UpdateLowMarginThreshold(-5.0)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidThreshold, err)
	assert.Equal(t, originalThreshold, settings.LowMarginThreshold, "Threshold should not change on error")
}

func TestUpdateLowMarginThreshold_HighThreshold(t *testing.T) {
	settings := NewShopSettings("Test Cafe")

	err := settings.UpdateLowMarginThreshold(50.0)

	assert.NoError(t, err, "High threshold should be valid")
	assert.Equal(t, 50.0, settings.LowMarginThreshold)
}

func TestUpdateLowMarginThreshold_DecimalThreshold(t *testing.T) {
	settings := NewShopSettings("Test Cafe")

	err := settings.UpdateLowMarginThreshold(22.5)

	assert.NoError(t, err)
	assert.Equal(t, 22.5, settings.LowMarginThreshold)
}
