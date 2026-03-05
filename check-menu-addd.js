const { MongoClient } = require('mongodb');

async function checkMenuAddd() {
  const uri = process.env.MONGODB_URI || 'mongodb://localhost:27017';
  const client = new MongoClient(uri);

  try {
    await client.connect();
    console.log('✅ Connected to MongoDB');

    const db = client.db('cafe_pos');
    const menuCollection = db.collection('menu_items');

    // Find menu item with name containing "áddd" or "addd"
    const menuItem = await menuCollection.findOne({
      name: { $regex: /[aá]ddd/i }
    });

    if (!menuItem) {
      console.log('❌ Không tìm thấy món "áddd"');
      return;
    }

    console.log('\n📋 THÔNG TIN MÓN "áddd":');
    console.log('=====================================');
    console.log('ID:', menuItem._id);
    console.log('Tên:', menuItem.name);
    console.log('Giá:', menuItem.price);
    console.log('Danh mục:', menuItem.category);
    console.log('Has Variants:', menuItem.has_variants);
    console.log('\n💰 CHI PHÍ:');
    console.log('Current Cost:', menuItem.current_cost);
    console.log('Cost Status:', menuItem.cost_status);
    console.log('Cost Last Calculated:', menuItem.cost_last_calculated_at);

    console.log('\n🥬 NGUYÊN LIỆU:');
    if (menuItem.ingredients && menuItem.ingredients.length > 0) {
      console.log(`Số lượng: ${menuItem.ingredients.length}`);
      for (const ing of menuItem.ingredients) {
        console.log(`  - ${ing.name}: ${ing.quantity} ${ing.unit}`);
        console.log(`    ID: ${ing.ingredient_id}`);
        
        // Check if ingredient exists in database
        const ingredientDoc = await db.collection('ingredients').findOne({
          _id: ing.ingredient_id
        });
        
        if (ingredientDoc) {
          console.log(`    ✅ Tồn tại - Giá: ${ingredientDoc.cost_per_unit} ${ingredientDoc.unit}`);
          console.log(`    Stock: ${ingredientDoc.current_stock}`);
        } else {
          console.log(`    ❌ KHÔNG TÌM THẤY TRONG DATABASE`);
        }
      }
    } else {
      console.log('❌ KHÔNG CÓ NGUYÊN LIỆU NÀO');
    }

    console.log('\n🔍 VARIANTS:');
    if (menuItem.variants && menuItem.variants.length > 0) {
      console.log(`Số lượng: ${menuItem.variants.length}`);
      for (const variant of menuItem.variants) {
        console.log(`\n  📦 ${variant.name}:`);
        console.log(`    ID: ${variant.id}`);
        console.log(`    Giá: ${variant.price}`);
        console.log(`    Chi phí: ${variant.current_cost}`);
        console.log(`    Cost Status: ${variant.cost_status}`);
        console.log(`    Nguyên liệu: ${variant.ingredients?.length || 0}`);
        
        if (variant.ingredients && variant.ingredients.length > 0) {
          for (const ing of variant.ingredients) {
            console.log(`      - ${ing.name}: ${ing.quantity} ${ing.unit}`);
          }
        }
      }
    } else {
      console.log('Không có variants');
    }

    // Check if there are any missing ingredients
    console.log('\n⚠️  PHÂN TÍCH VẤN ĐỀ:');
    
    if (!menuItem.ingredients || menuItem.ingredients.length === 0) {
      console.log('🔴 Món này KHÔNG CÓ NGUYÊN LIỆU nào được định nghĩa');
      console.log('   → Không thể tính chi phí');
    } else {
      let missingCount = 0;
      for (const ing of menuItem.ingredients) {
        const ingredientDoc = await db.collection('ingredients').findOne({
          _id: ing.ingredient_id
        });
        if (!ingredientDoc) {
          missingCount++;
        }
      }
      
      if (missingCount > 0) {
        console.log(`🔴 Có ${missingCount}/${menuItem.ingredients.length} nguyên liệu KHÔNG TỒN TẠI trong database`);
        console.log('   → Cost Status sẽ là INCOMPLETE');
      } else {
        console.log('✅ Tất cả nguyên liệu đều tồn tại');
        
        // Check if all ingredients have cost
        let noCostCount = 0;
        for (const ing of menuItem.ingredients) {
          const ingredientDoc = await db.collection('ingredients').findOne({
            _id: ing.ingredient_id
          });
          if (ingredientDoc && (!ingredientDoc.cost_per_unit || ingredientDoc.cost_per_unit === 0)) {
            noCostCount++;
            console.log(`⚠️  Nguyên liệu "${ingredientDoc.name}" chưa có giá (cost_per_unit = ${ingredientDoc.cost_per_unit})`);
          }
        }
        
        if (noCostCount > 0) {
          console.log(`🟡 Có ${noCostCount}/${menuItem.ingredients.length} nguyên liệu CHƯA CÓ GIÁ`);
          console.log('   → Cost Status sẽ là ESTIMATED hoặc INCOMPLETE');
        }
      }
    }

  } catch (error) {
    console.error('❌ Error:', error);
  } finally {
    await client.close();
  }
}

checkMenuAddd();
