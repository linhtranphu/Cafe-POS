package order

import (
	"testing"
)

// TestParseRoleType tests the ParseRoleType function
func TestParseRoleType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RoleType
	}{
		{
			name:     "Parse waiter role",
			input:    "waiter",
			expected: RoleWaiter,
		},
		{
			name:     "Parse barista role",
			input:    "barista",
			expected: RoleBarista,
		},
		{
			name:     "Parse cashier role",
			input:    "cashier",
			expected: RoleCashier,
		},
		{
			name:     "Parse manager role (should default to waiter)",
			input:    "manager",
			expected: RoleWaiter,
		},
		{
			name:     "Parse empty string (should default to waiter)",
			input:    "",
			expected: RoleWaiter,
		},
		{
			name:     "Parse invalid role (should default to waiter)",
			input:    "invalid_role",
			expected: RoleWaiter,
		},
		{
			name:     "Parse uppercase role (should default to waiter)",
			input:    "WAITER",
			expected: RoleWaiter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseRoleType(tt.input)
			if result != tt.expected {
				t.Errorf("ParseRoleType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRoleType_IsValid tests the IsValid method
func TestRoleType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		roleType RoleType
		expected bool
	}{
		{
			name:     "Waiter role is valid",
			roleType: RoleWaiter,
			expected: true,
		},
		{
			name:     "Barista role is valid",
			roleType: RoleBarista,
			expected: true,
		},
		{
			name:     "Cashier role is valid",
			roleType: RoleCashier,
			expected: true,
		},
		{
			name:     "Empty role is invalid",
			roleType: RoleType(""),
			expected: false,
		},
		{
			name:     "Invalid role is invalid",
			roleType: RoleType("invalid"),
			expected: false,
		},
		{
			name:     "Manager role is invalid (not in RoleType)",
			roleType: RoleType("manager"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.roleType.IsValid()
			if result != tt.expected {
				t.Errorf("RoleType(%q).IsValid() = %v, want %v", tt.roleType, result, tt.expected)
			}
		})
	}
}

// TestRoleType_String tests the String method
func TestRoleType_String(t *testing.T) {
	tests := []struct {
		name     string
		roleType RoleType
		expected string
	}{
		{
			name:     "Waiter role to string",
			roleType: RoleWaiter,
			expected: "waiter",
		},
		{
			name:     "Barista role to string",
			roleType: RoleBarista,
			expected: "barista",
		},
		{
			name:     "Cashier role to string",
			roleType: RoleCashier,
			expected: "cashier",
		},
		{
			name:     "Empty role to string",
			roleType: RoleType(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.roleType.String()
			if result != tt.expected {
				t.Errorf("RoleType(%q).String() = %q, want %q", tt.roleType, result, tt.expected)
			}
		})
	}
}

// TestParseRoleType_FromUserRole simulates the conversion from user.Role
func TestParseRoleType_FromUserRole(t *testing.T) {
	// Simulate user.Role type (which is also a string type)
	type UserRole string
	
	const (
		UserRoleWaiter  UserRole = "waiter"
		UserRoleBarista UserRole = "barista"
		UserRoleCashier UserRole = "cashier"
		UserRoleManager UserRole = "manager"
	)

	tests := []struct {
		name     string
		userRole UserRole
		expected RoleType
	}{
		{
			name:     "Convert user.Role waiter to order.RoleType",
			userRole: UserRoleWaiter,
			expected: RoleWaiter,
		},
		{
			name:     "Convert user.Role barista to order.RoleType",
			userRole: UserRoleBarista,
			expected: RoleBarista,
		},
		{
			name:     "Convert user.Role cashier to order.RoleType",
			userRole: UserRoleCashier,
			expected: RoleCashier,
		},
		{
			name:     "Convert user.Role manager to order.RoleType (defaults to waiter)",
			userRole: UserRoleManager,
			expected: RoleWaiter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This simulates the conversion in shift_handler.go:
			// roleType := order.ParseRoleType(string(role.(user.Role)))
			result := ParseRoleType(string(tt.userRole))
			
			if result != tt.expected {
				t.Errorf("ParseRoleType(string(%q)) = %v, want %v", tt.userRole, result, tt.expected)
			}
		})
	}
}

// TestRoleTypeConversion_InterfaceAssertion tests the interface assertion pattern
func TestRoleTypeConversion_InterfaceAssertion(t *testing.T) {
	// Simulate user.Role type
	type UserRole string
	const UserRoleBarista UserRole = "barista"

	// Simulate context.Get("role") returning interface{}
	var contextValue interface{} = UserRoleBarista

	// This is what was failing before the fix:
	// role.(string) would panic because role is UserRole, not string
	
	// Correct way: assert to UserRole first, then convert to string
	userRole, ok := contextValue.(UserRole)
	if !ok {
		t.Fatal("Failed to assert interface{} to UserRole")
	}

	// Convert to string
	roleString := string(userRole)
	if roleString != "barista" {
		t.Errorf("Expected 'barista', got %q", roleString)
	}

	// Parse to RoleType
	roleType := ParseRoleType(roleString)
	if roleType != RoleBarista {
		t.Errorf("Expected RoleBarista, got %v", roleType)
	}
}

// TestRoleTypeConversion_WrongAssertion demonstrates the bug
func TestRoleTypeConversion_WrongAssertion(t *testing.T) {
	// Simulate user.Role type
	type UserRole string
	const UserRoleBarista UserRole = "barista"

	// Simulate context.Get("role") returning interface{}
	var contextValue interface{} = UserRoleBarista

	// This should panic (the original bug)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when asserting UserRole as string, but didn't panic")
		} else {
			// Expected panic
			t.Logf("Correctly panicked with: %v", r)
		}
	}()

	// This will panic: interface conversion: interface {} is order.UserRole, not string
	_ = contextValue.(string)
}

// BenchmarkParseRoleType benchmarks the ParseRoleType function
func BenchmarkParseRoleType(b *testing.B) {
	roles := []string{"waiter", "barista", "cashier", "manager", "invalid"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		role := roles[i%len(roles)]
		_ = ParseRoleType(role)
	}
}

// BenchmarkRoleType_IsValid benchmarks the IsValid method
func BenchmarkRoleType_IsValid(b *testing.B) {
	roleTypes := []RoleType{RoleWaiter, RoleBarista, RoleCashier, RoleType("invalid")}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		roleType := roleTypes[i%len(roleTypes)]
		_ = roleType.IsValid()
	}
}
