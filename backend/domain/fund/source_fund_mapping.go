package fund

// SourceFundTypeMap định nghĩa quỹ mặc định cho từng loại nguồn giao dịch.
// Thay đổi ở đây sẽ ảnh hưởng toàn bộ hệ thống.
var SourceFundTypeMap = map[SourceType]FundType{
	SourceTypeExpense:      FundTypeOperating,  // Chi tiêu thủ công → quỹ vận hành
	SourceTypeIngredient:   FundTypeInventory,  // Mua nguyên liệu   → quỹ hàng hóa
	SourceTypeFacility:     FundTypeOperating,  // Mua cơ sở vật chất → quỹ vận hành
	SourceTypeHandover:     FundTypeCashDrawer, // Bàn giao ca        → ngăn kéo tiền
	SourceTypeManual:       FundTypeOperating,  // Thủ công            → quỹ vận hành
	SourceTypeFundTransfer: FundTypeOperating,  // Chuyển quỹ          → xác định lúc runtime
	SourceTypeWaiterStart:  FundTypeCashDrawer, // Tiền đầu ca waiter  → ngăn kéo tiền
}

// FundTypeForSource trả về quỹ mặc định cho một loại nguồn.
// Nếu không có mapping, trả về FundTypeOperating.
func FundTypeForSource(sourceType SourceType) FundType {
	if ft, ok := SourceFundTypeMap[sourceType]; ok {
		return ft
	}
	return FundTypeOperating
}
