// Test logo rendering in bill
const { MongoClient, ObjectId } = require('mongodb');

const uri = process.env.MONGODB_URI || 'mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin';

async function testLogoRendering() {
  const client = new MongoClient(uri);
  
  try {
    await client.connect();
    console.log('Connected to MongoDB\n');
    
    const db = client.db();
    
    // Get latest print job
    const printJobs = await db.collection('print_jobs').find({
      type: 'BILL'
    }).sort({ created_at: -1 }).limit(1).toArray();
    
    if (printJobs.length === 0) {
      console.log('No print jobs found. Create an order first.');
      return;
    }
    
    const job = printJobs[0];
    console.log('=== Latest Bill Print Job ===');
    console.log('ID:', job._id);
    console.log('Order:', job.order_id);
    console.log('Status:', job.status);
    console.log('Created:', job.created_at);
    console.log('\n=== Content Preview (first 500 chars) ===');
    console.log(job.content.substring(0, 500));
    console.log('\n=== Checking for Logo ===');
    
    // Check if content has ESC/POS image commands (GS v 0)
    const hasEscPosImage = job.content.includes('\x1D\x76\x30');
    console.log('Has ESC/POS image command (GS v 0):', hasEscPosImage);
    
    // Check if [LOGO] marker still exists (should not)
    const hasLogoMarker = job.content.includes('[LOGO]');
    console.log('Has [LOGO] marker (should be false):', hasLogoMarker);
    
    // Check content length
    console.log('Content length:', job.content.length, 'bytes');
    
    if (!hasEscPosImage && !hasLogoMarker) {
      console.log('\n❌ Logo not rendered! Neither ESC/POS commands nor [LOGO] marker found.');
      console.log('This means processTemplateContent() removed the marker but did not add logo.');
    } else if (hasLogoMarker) {
      console.log('\n⚠️  [LOGO] marker still in content! Logo was not processed.');
    } else if (hasEscPosImage) {
      console.log('\n✅ Logo ESC/POS commands found in content!');
    }
    
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.close();
  }
}

testLogoRendering();
