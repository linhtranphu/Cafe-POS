// Create a test order to trigger logo rendering
const { MongoClient, ObjectId } = require('mongodb');

const uri = 'mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin';

async function createTestOrder() {
  const client = new MongoClient(uri);
  
  try {
    await client.connect();
    const db = client.db();
    
    // Create test order
    const order = {
      _id: new ObjectId(),
      order_number: `TEST-${Date.now()}`,
      items: [{
        name: 'Test Coffee',
        variant_name: '',
        quantity: 1,
        price: 30000,
        subtotal: 30000,
        note: ''
      }],
      subtotal: 30000,
      discount: 0,
      total: 30000,
      waiter_name: 'Test Waiter',
      table_number: '1',
      status: 'COMPLETED',
      created_at: new Date(),
      updated_at: new Date()
    };
    
    await db.collection('orders').insertOne(order);
    console.log('✅ Test order created:', order.order_number);
    console.log('Order ID:', order._id);
    console.log('\nNow check backend logs for [LOGO DEBUG] messages');
    console.log('Or wait for auto-print to trigger...');
    
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.close();
  }
}

createTestOrder();
