// Check logo settings in MongoDB
const { MongoClient } = require('mongodb');

const uri = process.env.MONGODB_URI || 'mongodb://localhost:27017/cafe_pos';

async function checkLogoSettings() {
  const client = new MongoClient(uri);
  
  try {
    await client.connect();
    console.log('Connected to MongoDB');
    
    const db = client.db();
    const settings = await db.collection('shop_settings').findOne({});
    
    console.log('\n=== Shop Settings ===');
    console.log('show_logo:', settings?.show_logo);
    console.log('logo_url:', settings?.logo_url);
    console.log('\nFull settings:', JSON.stringify(settings, null, 2));
    
    // Check templates with [LOGO] marker
    console.log('\n=== Templates with [LOGO] ===');
    const templates = await db.collection('print_templates').find({
      content: { $regex: '\\[LOGO\\]' }
    }).toArray();
    
    console.log(`Found ${templates.length} templates with [LOGO] marker`);
    templates.forEach(t => {
      console.log(`- ${t.name} (${t.type}, default: ${t.is_default})`);
    });
    
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.close();
  }
}

checkLogoSettings();
