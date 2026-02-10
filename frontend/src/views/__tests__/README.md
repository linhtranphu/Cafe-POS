# MenuCostView Component Tests

## Overview

This directory contains unit tests for the MenuCostView component, covering:
- Component rendering with mock data
- Category filtering (Requirement 4.3)
- Sorting functionality (Requirement 4.4)
- Warning color coding (Requirements 7.2, 7.3)
- Search functionality
- Cost breakdown modal
- Helper functions

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
    setupFiles: ['./src/views/__tests__/setup.js']
  }
})
```

### 3. Create Test Setup File

Create `frontend/src/views/__tests__/setup.js`:

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
    if (!value) return '0đ'
    return `${value.toLocaleString('vi-VN')}đ`
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
npm test MenuCostView.test.js
```

## Test Coverage

The test suite covers:

### Component Rendering (Requirement 4.1)
- ✅ Renders component with menu items
- ✅ Displays summary statistics
- ✅ Shows loading state
- ✅ Shows error state

### Category Filtering (Requirement 4.3)
- ✅ Displays all items when no filter selected
- ✅ Filters items by Coffee category
- ✅ Filters items by Tea category
- ✅ Shows all items when filter cleared

### Sorting Functionality (Requirement 4.4)
- ✅ Sorts by profit_margin descending (default)
- ✅ Sorts by profit_margin ascending
- ✅ Sorts by absolute_profit descending
- ✅ Sorts by name ascending
- ✅ Toggles sort order

### Warning Color Coding (Requirements 7.2, 7.3)
- ✅ Green border for profitable items
- ✅ Yellow border for low margin items
- ✅ Red border for loss items
- ✅ Correct text colors for profit margins

### Search Functionality
- ✅ Filters by name search
- ✅ Filters by category search
- ✅ Case insensitive search

### Cost Breakdown Modal
- ✅ Opens modal on item click
- ✅ Displays ingredient breakdown
- ✅ Closes modal on close button

### Helper Functions
- ✅ Formats percentages correctly
- ✅ Gets correct cost status labels
- ✅ Gets correct warning messages

## Notes

- Tests use mocked API responses to avoid external dependencies
- All tests follow the MINIMAL testing approach - focusing on core functionality
- Tests validate requirements 4.1, 4.3, 4.4, 7.2, and 7.3
- Mock data matches the API response structure from the design document
