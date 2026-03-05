// Fix template - remove TableNumber field
const { MongoClient } = require('mongodb');

const uri = 'mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin';

async function fixTemplate() {
  const client = new MongoClient(uri);
  
  try {
    await client.connect();
    const db = client.db();
    
    // Get the template with [LOGO]
    const template = await db.collection('print_templates').findOne({
      name: 'bill logo'
    });
    
    if (!template) {
      console.log('Template not found');
      return;
    }
    
    console.log('Current template content:');
    console.log(template.content);
    
    // Remove TableNumber references
    let newContent = template.content;
    
    // Remove the line with TableNumber
    newContent = newContent.replace(/Bàn: \{\{if \.Order\.TableNumber\}\}\{\{\.Order\.TableNumber\}\}\{\{else\}\}Mang về\{\{end\}\}.*?\n/g, '');
    
    // Update template
    await db.collection('print_templates').updateOne(
      { _id: template._id },
      { $set: { content: newContent } }
    );
    
    console.log('\n✅ Template updated successfully!');
    console.log('\nNew content:');
    console.log(newContent);
    
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.close();
  }
}

fixTemplate();
