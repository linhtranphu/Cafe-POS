#!/bin/bash
cd backend
# Remove the incorrectly indented lines
sed -i '' '/ttt\/\/ Reprint routes/,/ttt$/d' main.go

# Add correctly indented lines before "// Waiter routes"
awk '/\t\t\t\/\/ Waiter routes/ {
    print "\t\t\t// Reprint routes (only for cashier and manager)"
    print "\t\t\treprintGroup := protected.Group(\"/orders/:id\")"
    print "\t\t\treprintGroup.Use(http.RequireRole(user.RoleCashier, user.RoleManager))"
    print "\t\t\t{"
    print "\t\t\t\treprintGroup.POST(\"/reprint-bill\", orderHandler.ReprintBill)"
    print "\t\t\t\treprintGroup.POST(\"/reprint-label\", orderHandler.ReprintLabel)"
    print "\t\t\t}"
    print ""
}
{print}' main.go > main.go.tmp && mv main.go.tmp main.go
