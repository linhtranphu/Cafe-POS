package expense

// Default category names (Vietnamese)
const (
	CategoryIngredient  = "Nguyên liệu"
	CategoryFacility    = "Cơ sở vật chất"
	CategoryMaintenance = "Bảo trì"
	CategoryUtility     = "Tiện ích"
	CategorySalary      = "Nhân sự"
	CategoryMarketing   = "Marketing"
	CategoryOther       = "Khác"
)

// GetDefaultCategories returns list of default expense categories
// These categories will be auto-created when the system starts
func GetDefaultCategories() []string {
	return []string{
		CategoryIngredient,
		CategoryFacility,
		CategoryMaintenance,
		CategoryUtility,
		CategorySalary,
		CategoryMarketing,
		CategoryOther,
	}
}

// GetCategoryDescription returns a description for each category
func GetCategoryDescription(categoryName string) string {
	descriptions := map[string]string{
		CategoryIngredient:  "Chi phí mua nguyên liệu, thực phẩm",
		CategoryFacility:    "Chi phí mua sắm thiết bị, cơ sở vật chất",
		CategoryMaintenance: "Chi phí bảo trì, sửa chữa thiết bị",
		CategoryUtility:     "Chi phí điện, nước, internet, điện thoại",
		CategorySalary:      "Chi phí lương, thưởng nhân viên",
		CategoryMarketing:   "Chi phí quảng cáo, marketing",
		CategoryOther:       "Chi phí khác",
	}
	
	if desc, ok := descriptions[categoryName]; ok {
		return desc
	}
	return ""
}
