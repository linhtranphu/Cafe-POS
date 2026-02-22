const io = require('socket.io-client');

console.log('Testing WebSocket connection to backend...');

const socket = io('http://localhost:3000', {
  transports: ['websocket', 'polling'],
  reconnection: false,
  timeout: 5000
});

socket.on('connect', () => {
  console.log('✅ Connected to backend!');
  console.log('Socket ID:', socket.id);
  
  // Test listening for print-job-created event
  socket.on('print-job-created', (data) => {
    console.log('📨 Received print-job-created event:', data);
  });
  
  console.log('Waiting for events...');
  
  setTimeout(() => {
    console.log('Test complete. Disconnecting...');
    socket.disconnect();
    process.exit(0);
  }, 10000);
});

socket.on('connect_error', (error) => {
  console.error('❌ Connection error:', error.message);
  process.exit(1);
});

socket.on('disconnect', (reason) => {
  console.log('Disconnected:', reason);
});

socket.on('error', (error) => {
  console.error('Socket error:', error);
});
