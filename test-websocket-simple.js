#!/usr/bin/env node

/**
 * Simple WebSocket Connection Test
 * Tests if Print Bridge can connect to backend Socket.IO server
 */

const io = require('socket.io-client');

console.log('🧪 Testing WebSocket Connection');
console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
console.log('');

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:3000';
console.log(`📡 Connecting to: ${BACKEND_URL}`);
console.log('');

const socket = io(BACKEND_URL, {
  path: '/socket.io/',
  transports: ['websocket', 'polling'],
  reconnection: false,
  timeout: 10000
});

let connected = false;

socket.on('connect', () => {
  connected = true;
  console.log('✅ SUCCESS! Connected to backend');
  console.log(`   Socket ID: ${socket.id}`);
  console.log('');
  console.log('🎧 Listening for events...');
  console.log('   - print-job-created');
  console.log('   - print-job-status-changed');
  console.log('');
  console.log('💡 To test print job broadcast:');
  console.log('   1. Create an order in frontend');
  console.log('   2. Or use API to create print job');
  console.log('');
  console.log('⏱️  Waiting 30 seconds for events...');
  
  setTimeout(() => {
    console.log('');
    console.log('⏰ Test timeout reached');
    if (connected) {
      console.log('✅ Connection stable for 30 seconds');
    }
    socket.disconnect();
    process.exit(0);
  }, 30000);
});

socket.on('print-job-created', (data) => {
  console.log('');
  console.log('📨 Received: print-job-created');
  console.log('   Data:', JSON.stringify(data, null, 2));
  console.log('');
});

socket.on('print-job-status-changed', (data) => {
  console.log('');
  console.log('📊 Received: print-job-status-changed');
  console.log('   Data:', JSON.stringify(data, null, 2));
  console.log('');
});

socket.on('connect_error', (error) => {
  console.error('❌ Connection Error:', error.message);
  console.log('');
  console.log('🔍 Troubleshooting:');
  console.log('   1. Is backend running? curl http://localhost:3000/health');
  console.log('   2. Is Socket.IO endpoint working?');
  console.log('      curl "http://localhost:3000/socket.io/?EIO=3&transport=polling"');
  console.log('   3. Check backend logs: tail -f backend.log | grep Socket');
  console.log('');
  process.exit(1);
});

socket.on('disconnect', (reason) => {
  console.log('');
  console.log('🔌 Disconnected:', reason);
  if (!connected) {
    console.log('❌ Never connected successfully');
    process.exit(1);
  }
});

socket.on('error', (error) => {
  console.error('');
  console.error('⚠️  Socket Error:', error);
  console.log('');
});

// Handle Ctrl+C
process.on('SIGINT', () => {
  console.log('');
  console.log('🛑 Test interrupted by user');
  socket.disconnect();
  process.exit(0);
});
