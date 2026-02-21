/**
 * Test script to verify printer connection and print functionality
 * Usage: node src/test-printer.js <printer-ip> [port]
 */

const printerService = require('./services/printerService')

// Get printer IP from command line
const printerIP = process.argv[2]
const printerPort = parseInt(process.argv[3]) || 9100

if (!printerIP) {
  console.error('Usage: node src/test-printer.js <printer-ip> [port]')
  console.error('Example: node src/test-printer.js 192.168.1.100 9100')
  process.exit(1)
}

// ESC/POS test content
const testContent = `
================================
   PRINTER CONNECTION TEST
================================
Date: ${new Date().toLocaleString()}
Printer: ${printerIP}:${printerPort}
================================

This is a test print from
Local Print Bridge

If you can read this, the
printer is working correctly!

================================
        TEST SUCCESSFUL
================================


`

async function runTest() {
  console.log('🖨️  Local Print Bridge - Printer Test')
  console.log('='.repeat(50))
  console.log(`Testing printer: ${printerIP}:${printerPort}`)
  console.log('='.repeat(50))

  // Test 1: Connection test
  console.log('\n📡 Test 1: Testing connection...')
  try {
    await printerService.testConnection(printerIP, printerPort)
    console.log('✅ Connection test PASSED')
  } catch (error) {
    console.error('❌ Connection test FAILED:', error.message)
    console.log('\nTroubleshooting:')
    console.log('1. Check if printer is powered on')
    console.log('2. Verify printer IP address is correct')
    console.log('3. Ensure printer is connected to the same network')
    console.log('4. Check if port 9100 is open on the printer')
    process.exit(1)
  }

  // Test 2: Print test
  console.log('\n🖨️  Test 2: Sending test print...')
  try {
    await printerService.print(testContent, printerIP, printerPort)
    console.log('✅ Print test PASSED')
    console.log('\n📄 Check your printer for the test printout!')
  } catch (error) {
    console.error('❌ Print test FAILED:', error.message)
    process.exit(1)
  }

  // Show stats
  console.log('\n📊 Statistics:')
  const stats = printerService.getStats()
  console.log(`   Total prints: ${stats.totalPrints}`)
  console.log(`   Successful: ${stats.successfulPrints}`)
  console.log(`   Failed: ${stats.failedPrints}`)
  console.log(`   Success rate: ${stats.successRate}`)

  console.log('\n' + '='.repeat(50))
  console.log('✅ All tests completed successfully!')
  console.log('='.repeat(50))
}

runTest().catch(error => {
  console.error('\n❌ Test failed:', error.message)
  process.exit(1)
})
