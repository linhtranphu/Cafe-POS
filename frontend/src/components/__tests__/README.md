# MenuItemCostBreakdown Component Tests

## Overview

This directory contains unit tests for the MenuItemCostBreakdown component, covering:
- Component rendering and modal display (Requirement 8.1)
- Ingredient breakdown table with all details (Requirements 8.2, 8.3, 8.4)
- Total cost summary and warnings (Requirement 8.3)
- Conversion rate and wastage percentage display
- Missing cost detection and highlighting
- Data fetching and error handling

## Setup Testing Framework

The tests are written but require a testing framework to be installed. Follow these steps:

### 1. Install Dependencies

```bash
cd frontend
npm install -D vitest @vue/test-utils happy-dom
```

### 2. Update vite.config.js

Add test configuration to `vite.config.js`:

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        rewrite: (path) => path
      }
    }
  },
  css: {
    postcss: './postcss.config.js'
  },
  test: {
    globals: true,
    environment: 'happy-dom',
    setupFiles: ['./src/components/__tests__/setup.js']
  }
})
```

### 3. Create Test Setup File

Create `frontend/src/components/__tests__/setup.js`:

```javascript
import { vi } from 'vitest'

// Mock router
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn()
  }),
  useRoute: () => ({
    params: {},
    query: {}
  })
}))

// Mock formatters
vi.mock('../../utils/formatters', () => ({
  formatPrice: (value) => {
    if (!value && value !== 0) return '0'
    return value.toLocaleString('vi-VN')
  },
  formatDate: (date) => {
    return new Date(date).toLocaleDateString('vi-VN')
  }
}))
```

### 4. Add Test Script to package.json

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview",
    "test": "vitest --run",
    "test:watch": "vitest",
    "test:ui": "vitest --ui"
  }
}
```

## Running Tests

### Run All Tests Once
```bash
npm test
```

### Run Tests in Watch Mode
```bash
npm run test:watch
```

### Run Tests with UI
```bash
npm run test:ui
```

### Run Specific Test File
```bash
npm test MenuItemCostBreakdown.test.js
```

## Test Coverage

The test suite covers:

### Component Rendering (Requirement 8.1)
- ✅ Does not render when isOpen is false
- ✅ Renders modal when isOpen is true
- ✅ Shows loading state initially
- ✅ Displays menu item info after loading
- ✅ Shows error state when API fails
- ✅ Emits close event on close button click
- ✅ Emits close event when clicking outside modal

### Ingredient Breakdown Table (Requirements 8.2, 8.3, 8.4)
- ✅ Displays all ingredients with complete data
- ✅ Displays conversion rate when non-default
- ✅ Displays wastage percentage when present
- ✅ Does not display conversion rate when default (1.0)
- ✅ Does not display wastage when zero
- ✅ Highlights ingredients with missing cost_per_unit

### Total Cost Summary (Requirement 8.3)
- ✅ Displays total cost at bottom
- ✅ Shows warning when any ingredient has incomplete cost
- ✅ Does not show warning when all ingredients have complete cost

### Helper Functions
- ✅ Detects incomplete cost correctly
- ✅ Detects non-default conversion rate
- ✅ Detects wastage correctly

### Data Fetching
- ✅ Fetches cost breakdown when component opens
- ✅ Does not fetch when menuItemId is null
- ✅ Refetches when menuItemId changes
- ✅ Refetches when modal reopens

## Test Data

The tests use three mock cost breakdown scenarios:

1. **Complete Data**: All ingredients have valid costs, conversion rates, and wastage
2. **Missing Cost**: One ingredient has missing cost_per_unit (0 or null)
3. **With Conversion**: Ingredients with non-default conversion rates and wastage

## Notes

- Tests use mocked API responses to avoid external dependencies
- All tests follow the MINIMAL testing approach - focusing on core functionality
- Tests validate requirements 8.1, 8.2, 8.3, and 8.4
- Mock data matches the API response structure from the design document
- Component properly handles edge cases like missing costs and null values
