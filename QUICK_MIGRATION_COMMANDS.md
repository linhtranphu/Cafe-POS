# Quick Migration Commands - Copy & Paste

## Option 1: Using Script (Recommended)

### Step 1: Copy script to EC2
```bash
scp migrate-v2.0-mongodb.sh ubuntu@your-ec2-ip:/home/ubuntu/
```

### Step 2: Run on EC2
```bash
ssh ubuntu@your-ec2-ip
sudo bash /home/ubuntu/migrate-v2.0-mongodb.sh
```

Type `yes` when prompted.

---

## Option 2: Manual Commands (If script fails)

Copy and paste these commands directly on EC2:

### 1. Create shop_settings
```bash
sudo docker exec mongodb mongosh cafe_pos --quiet --eval '
var existing = db.shop_settings.findOne();
if (existing) {
    print("✅ Shop settings already exist");
} else {
    var now = new Date();
    db.shop_settings.insertOne({
        shop_name: "Cafe POS",
        shop_address: "123 Main Street",
        shop_phone: "0123-456-789",
        logo_url: "",
        custom_message: "Cảm ơn quý khách! Hẹn gặp lại!",
        show_logo: false,
        show_address: true,
        show_phone: true,
        show_custom_message: true,
        auto_print_enabled: true,
        created_at: now,
        updated_at: now
    });
    print("✅ Created shop settings");
}
'
```

### 2. Create indexes for print_jobs
```bash
sudo docker exec mongodb mongosh cafe_pos --quiet --eval '
db.print_jobs.createIndex({ "status": 1, "created_at": -1 });
db.print_jobs.createIndex({ "order_id": 1 });
db.print_jobs.createIndex({ "printer_id": 1 });
print("✅ Created print_jobs indexes");
'
```

### 3. Create indexes for printer_configs
```bash
sudo docker exec mongodb mongosh cafe_pos --quiet --eval '
db.printer_configs.createIndex({ "type": 1, "is_default": 1 });
print("✅ Created printer_configs indexes");
'
```

### 4. Create indexes for print_templates
```bash
sudo docker exec mongodb mongosh cafe_pos --quiet --eval '
db.print_templates.createIndex({ "type": 1, "is_default": 1 });
print("✅ Created print_templates indexes");
'
```

### 5. Verify migration
```bash
sudo docker exec mongodb mongosh cafe_pos --quiet --eval '
print("Shop settings: " + db.shop_settings.countDocuments());
print("Orders: " + db.orders.countDocuments());
print("Users: " + db.users.countDocuments());
'
```

Expected output:
```
Shop settings: 1
Orders: [your count]
Users: [your count]
```

### 6. Restart backend
```bash
sudo docker-compose restart backend
```

### 7. Check logs
```bash
sudo docker-compose logs backend | tail -50
```

Look for:
```
✅ MongoDB connected successfully
Server starting on :3000
```

---

## Verify API Works

```bash
# Test shop settings API (replace TOKEN with actual token)
curl -s https://tacafe.store/api/manager/shop-settings \
  -H "Authorization: Bearer YOUR_TOKEN" | jq
```

Expected: JSON with shop settings

---

## Troubleshooting

### If MongoDB commands fail with authentication error:

```bash
# Try with authentication
sudo docker exec mongodb mongosh cafe_pos \
  -u admin -p password123 --authenticationDatabase admin \
  --eval 'db.shop_settings.find().pretty()'
```

### If container name is different:

```bash
# Find MongoDB container name
sudo docker ps | grep mongo

# Use the actual name (replace mongodb with actual name)
sudo docker exec [ACTUAL_CONTAINER_NAME] mongosh cafe_pos --eval '...'
```

### Check if shop_settings was created:

```bash
sudo docker exec mongodb mongosh cafe_pos --quiet --eval 'db.shop_settings.find().pretty()'
```

Should show 1 document with shop_name, shop_address, etc.

---

## Complete Deployment Checklist

- [ ] Backup database
- [ ] Build and push Docker images (`./build_docker_hub.sh`)
- [ ] Pull new images on EC2 (`sudo docker-compose pull`)
- [ ] Stop services (`sudo docker-compose down`)
- [ ] Start services (`sudo docker-compose up -d`)
- [ ] Run migration (this script)
- [ ] Restart backend (`sudo docker-compose restart backend`)
- [ ] Verify API works
- [ ] Test frontend (https://tacafe.store)
- [ ] Configure shop settings in UI

---

## Rollback (if needed)

```bash
# Remove shop_settings (to re-run migration)
sudo docker exec mongodb mongosh cafe_pos --quiet --eval 'db.shop_settings.deleteMany({})'

# Then re-run migration
sudo bash /home/ubuntu/migrate-v2.0-mongodb.sh
```
