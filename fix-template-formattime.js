// Fix template - fix formatTime calls
const { MongoClient } = require('mongodb');

const uri = 'mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin';

async function fixTemplate() {
  const client = new MongoClient(uri);
  
  try {
    await client.connect();
    const db = client.db();
    
    const template = await db.collection('print_templates').findOne({
      name: 'bill logo'
    });
    
    if (!template) {
      console.log('Template not found');
      return;
    }
    
    let newContent = template.content;
    
    // Fix formatTime calls - add layout parameter
    newContent = newContent.replace(/\{\{formatTime \.Order\.CreatedAt\s*\}\}/g, '{{formatTime .Order.CreatedAt "02/01 15:04"}}');
    
    // Update template
    await db.collection('print_templates').updateOne(
      { _id: template._id },
      { $set: { content: newContent } }
    );
    
    console.log('✅ Template fixed!');
    
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.close();
  }
}

fixTemplate();
