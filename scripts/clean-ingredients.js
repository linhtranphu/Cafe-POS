// Script to clean all ingredient data from MongoDB
const { MongoClient } = require('mongodb');

const MONGO_URI = process.env.MONGO_URI || 'mongodb://localhost:27017';
const DB_NAME = 'cafe_pos';

async function cleanIngredients() {
  const client = new MongoClient(MONGO_URI);
  
  try {
    await client.connect();
    console.log('✅ Connected to MongoDB');
    
    const db = client.db(DB_NAME);
    const ingredientsCollection = db.collection('ingredients');
    const stockHistoryCollection = db.collection('stock_history');
    
    // Delete all ingredients
    const ingredientsResult = await ingredientsCollection.deleteMany({});
    console.log(`🗑️  Deleted ${ingredientsResult.deletedCount} ingredients`);
    
    // Delete all stock history
    const historyResult = await stockHistoryCollection.deleteMany({});
    console.log(`🗑️  Deleted ${historyResult.deletedCount} stock history records`);
    
    console.log('✅ All ingredient data cleaned successfully!');
    
  } catch (error) {
    console.error('❌ Error cleaning ingredients:', error);
    process.exit(1);
  } finally {
    await client.close();
    console.log('👋 Disconnected from MongoDB');
  }
}

cleanIngredients();
