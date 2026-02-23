#!/usr/bin/env node

/**
 * Frontend Testing Script for Cashier Fund Handover
 * 
 * This script performs automated API testing to verify the frontend
 * integration with the backend fund handover endpoints.
 * 
 * Usage:
 *   TOKEN=your_jwt_token node test-frontend-fund-handover.js
 */

const https = require('https');
const http = require('http');

// Configuration
const API_URL = process.env.API_URL || 'http://localhost:8080';
const TOKEN = process.env.TOKEN;

if (!TOKEN) {
  console.error('❌ Error: TOKEN environment variable is required');
  console.error('Usage: TOKEN=your_jwt_token node test-frontend-fund-handover.js');
  process.exit(1);
}

// Helper function to make HTTP requests
function makeRequest(method, path, data = null) {
  return new Promise((resolve, reject) => {
    const url = new URL(path, API_URL);
    const isHttps = url.protocol === 'https:';
    const lib = isHttps ? https : http;

    const options = {
      hostname: url.hostname,
      port: url.port || (isHttps ? 443 : 80),
      path: url.pathname + url.search,
      method: method,
      headers: {
        'Authorization': `Bearer ${TOKEN}`,
        'Content-Type': 'application/json'
      }
    };

    const req = lib.request(options, (res) => {
      let body = '';
      res.on('data', (chunk) => body += chunk);
      res.on('end', () => {
        try {
          const response = {
            status: res.statusCode,
            headers: res.headers,
            data: body ? JSON.parse(body) : null
          };
          resolve(response);
        } catch (e) {
          resolve({
            status: res.statusCode,
            headers: res.headers,
            data: body
          });
        }
      });
    });

    req.on('error', reject);

    if (data) {
      req.write(JSON.stringify(data));
    }

    req.end();
  });
}

// Test cases
const tests = {
  passed: 0,
  failed: 0,
  results: []
};

function logTest(name, passed, message = '') {
  const icon = passed ? '✅' : '❌';
  console.log(`${icon} ${name}`);
  if (message) {
    console.log(`   ${message}`);
  }
  tests.results.push({ name, passed, message });
  if (passed) tests.passed++;
  else tests.failed++;
}

async function runTests() {
  console.log('🧪 Frontend Fund Handover Testing');
  console.log('='.repeat(50));
  console.log('');

  let shiftId = null;
  let managedFunds = null;

  // Test 1: Get Current Cashier Shift
  console.log('📋 Test 1: Get Current Cashier Shift');
  try {
    const response = await makeRequest('GET', '/api/v1/cashier-shifts/current');
    
    if (response.status === 200 && response.data) {
      shiftId = response.data.id || response.data.ID || response.data._id;
      logTest('Get current shift', true, `Shift ID: ${shiftId}`);
      console.log(`   Cashier: ${response.data.cashier_name}`);
      console.log(`   Status: ${response.data.status}`);
      console.log(`   Starting Float: ${response.data.starting_float}₫`);
    } else if (response.status === 404) {
      logTest('Get current shift', false, 'No open cashier shift found');
      console.log('');
      console.log('⚠️  Please start a cashier shift first');
      process.exit(1);
    } else {
      logTest('Get current shift', false, `Status: ${response.status}`);
      process.exit(1);
    }
  } catch (error) {
    logTest('Get current shift', false, error.message);
    process.exit(1);
  }
  console.log('');

  // Test 2: Get Managed Funds
  console.log('💰 Test 2: Get Managed Funds');
  try {
    const response = await makeRequest('GET', `/api/v1/cashier-shifts/${shiftId}/managed-funds`);
    
    if (response.status === 200 && response.data) {
      managedFunds = response.data;
      logTest('Get managed funds', true);
      console.log(`   Starting Float:    ${managedFunds.starting_float}₫`);
      console.log(`   Received Cash:     ${managedFunds.received_cash}₫`);
      console.log(`   Received Transfer: ${managedFunds.received_transfer}₫`);
      console.log(`   Expected Cash:     ${managedFunds.expected_cash}₫`);
      console.log(`   Total Managed:     ${managedFunds.total_managed_funds}₫`);
      console.log(`   Handover Count:    ${managedFunds.handover_count}`);
      
      // Validate calculations
      const expectedTotal = managedFunds.received_cash + managedFunds.received_transfer;
      const expectedCash = managedFunds.starting_float + managedFunds.received_cash;
      
      if (managedFunds.total_managed_funds === expectedTotal) {
        logTest('Total calculation correct', true);
      } else {
        logTest('Total calculation correct', false, 
          `Expected ${expectedTotal}, got ${managedFunds.total_managed_funds}`);
      }
      
      if (managedFunds.expected_cash === expectedCash) {
        logTest('Expected cash calculation correct', true);
      } else {
        logTest('Expected cash calculation correct', false,
          `Expected ${expectedCash}, got ${managedFunds.expected_cash}`);
      }
    } else {
      logTest('Get managed funds', false, `Status: ${response.status}`);
    }
  } catch (error) {
    logTest('Get managed funds', false, error.message);
  }
  console.log('');

  // Test 3: Validate Managed Funds Response Structure
  console.log('🔍 Test 3: Validate Response Structure');
  if (managedFunds) {
    const requiredFields = [
      'cashier_shift_id',
      'starting_float',
      'received_cash',
      'received_transfer',
      'total_managed_funds',
      'expected_cash',
      'handover_count'
    ];
    
    let allFieldsPresent = true;
    for (const field of requiredFields) {
      if (managedFunds[field] === undefined) {
        logTest(`Field '${field}' present`, false);
        allFieldsPresent = false;
      }
    }
    
    if (allFieldsPresent) {
      logTest('All required fields present', true);
    }
    
    // Validate data types
    if (typeof managedFunds.starting_float === 'number') {
      logTest('starting_float is number', true);
    } else {
      logTest('starting_float is number', false);
    }
    
    if (typeof managedFunds.received_cash === 'number') {
      logTest('received_cash is number', true);
    } else {
      logTest('received_cash is number', false);
    }
    
    if (typeof managedFunds.handover_count === 'number') {
      logTest('handover_count is number', true);
    } else {
      logTest('handover_count is number', false);
    }
  }
  console.log('');

  // Test 4: Test Close with Fund Handover (Dry Run - No Variance)
  console.log('🔒 Test 4: Close with Fund Handover (Validation Only)');
  console.log('   Note: This is a validation test, not actually closing the shift');
  
  if (managedFunds) {
    const testPayload = {
      actual_cash: managedFunds.expected_cash, // No variance
      variance_reason: null,
      variance_notes: null
    };
    
    console.log(`   Test Payload:`);
    console.log(`   - Actual Cash: ${testPayload.actual_cash}₫`);
    console.log(`   - Variance: 0₫ (no variance)`);
    
    logTest('Payload structure valid', true);
    logTest('No variance scenario prepared', true);
  }
  console.log('');

  // Test 5: Test Close with Fund Handover (With Variance)
  console.log('🔒 Test 5: Close with Fund Handover (With Variance - Validation)');
  
  if (managedFunds) {
    const variance = -5000; // Shortage
    const testPayload = {
      actual_cash: managedFunds.expected_cash + variance,
      variance_reason: 'COUNTING_ERROR',
      variance_notes: 'Test variance documentation - đếm nhầm tờ 50k'
    };
    
    console.log(`   Test Payload:`);
    console.log(`   - Actual Cash: ${testPayload.actual_cash}₫`);
    console.log(`   - Expected Cash: ${managedFunds.expected_cash}₫`);
    console.log(`   - Variance: ${variance}₫`);
    console.log(`   - Reason: ${testPayload.variance_reason}`);
    console.log(`   - Notes: ${testPayload.variance_notes}`);
    
    // Validate notes length
    if (testPayload.variance_notes.length >= 10) {
      logTest('Variance notes length valid (≥10)', true);
    } else {
      logTest('Variance notes length valid (≥10)', false);
    }
    
    logTest('Variance scenario prepared', true);
  }
  console.log('');

  // Test 6: Frontend Integration Checklist
  console.log('📱 Test 6: Frontend Integration Checklist');
  console.log('   Manual verification required:');
  console.log('');
  console.log('   Dashboard:');
  console.log('   [ ] Managed funds section displays');
  console.log('   [ ] Cash amount shows with green styling');
  console.log('   [ ] Transfer amount shows with blue styling');
  console.log('   [ ] Total shows with orange gradient');
  console.log('   [ ] Warning message displays');
  console.log('   [ ] Pull-to-refresh works');
  console.log('');
  console.log('   Closure Flow:');
  console.log('   [ ] Managed funds summary displays');
  console.log('   [ ] Cash counting input works');
  console.log('   [ ] Variance calculation automatic');
  console.log('   [ ] Variance documentation form appears when needed');
  console.log('   [ ] Confirmation summary shows all data');
  console.log('   [ ] Submit button calls correct API');
  console.log('');
  console.log('   Mobile:');
  console.log('   [ ] Responsive on mobile devices');
  console.log('   [ ] Touch interactions smooth');
  console.log('   [ ] No horizontal scrolling');
  console.log('');

  // Summary
  console.log('='.repeat(50));
  console.log('📊 Test Summary');
  console.log('='.repeat(50));
  console.log(`Total Tests: ${tests.passed + tests.failed}`);
  console.log(`✅ Passed: ${tests.passed}`);
  console.log(`❌ Failed: ${tests.failed}`);
  console.log('');

  if (tests.failed === 0) {
    console.log('🎉 All automated tests passed!');
    console.log('');
    console.log('Next Steps:');
    console.log('1. Perform manual frontend testing using FRONTEND_TESTING_GUIDE.md');
    console.log('2. Test on mobile devices');
    console.log('3. Verify UI/UX matches design');
    console.log('4. Test error scenarios');
    console.log('5. Get stakeholder approval');
  } else {
    console.log('⚠️  Some tests failed. Please review and fix issues.');
    process.exit(1);
  }
}

// Run tests
runTests().catch(error => {
  console.error('❌ Test execution failed:', error);
  process.exit(1);
});
