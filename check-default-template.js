// Check default template
const { MongoClient } = require('mongodb');

const uri = 'mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin';

async function checkTemplate() {
  const client = new MongoClient(uri);
  
  try {
    await client.connect();
    const db = client.db();
    
    // Get default BILL template
    const template = await db.collection('print_templates').findOne({
      type: 'BILL',
      is_default: true
    });
    
    if (!template) {
      console.log('❌ No default BILL template found!');
      return;
    }
    
    console.log('=== Default BILL Template ===');
    console.log('Name:', template.name);
    console.log('Type:', template.type);
    console.log('Is Default:', template.is_default);
    console.log('\n=== Content ===');
    console.log(template.content);
    console.log('\n=== Checking for [LOGO] marker ===');
    console.log('Has [LOGO]:', template.content.includes('[LOGO]'));
    
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.close();
  }
}

checkTemplate();
