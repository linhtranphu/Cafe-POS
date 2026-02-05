# 🚀 Cash Handover Deployment Checklist

**Feature:** Cash Handover (Waiter → Cashier)  
**Version:** 1.0.0  
**Date:** 2026-02-04

---

## ✅ Pre-Deployment Checklist

### Code Quality
- [x] All code written and committed
- [x] Backend compiles without errors
- [x] Frontend builds without errors
- [ ] Code review completed by team
- [ ] No critical bugs in issue tracker
- [ ] All TODOs resolved or documented

### Testing
- [x] Unit tests written (11 test cases)
- [x] Integration tests created (12 scenarios)
- [x] E2E test scenarios documented (10 scenarios)
- [ ] All tests passing
- [ ] Performance testing completed
- [ ] Security audit completed
- [ ] Mobile testing completed
- [ ] Cross-browser testing completed

### Documentation
- [x] API documentation complete
- [x] User guide complete (Vietnamese)
- [x] Technical documentation complete
- [x] Deployment guide ready
- [x] Rollback plan documented
- [ ] Training materials prepared

### Database
- [x] Migration script ready
- [x] Index creation script ready
- [ ] Migration tested on staging
- [ ] Backup strategy confirmed
- [ ] Rollback script prepared

### Infrastructure
- [ ] Server resources verified (CPU, RAM, Disk)
- [ ] Database capacity checked
- [ ] Network bandwidth sufficient
- [ ] Load balancer configured
- [ ] SSL certificates valid
- [ ] Monitoring tools ready

---

## 🔧 Deployment Steps

### Step 1: Backup (CRITICAL)
```bash
# Backup database
mongodump --db cafe_pos --out backup_$(date +%Y%m%d_%H%M%S)

# Backup current code
tar -czf code_backup_$(date +%Y%m%d_%H%M%S).tar.gz backend/ frontend/

# Verify backups
ls -lh backup_* code_backup_*
```

- [ ] Database backup completed
- [ ] Code backup completed
- [ ] Backups verified and stored safely

### Step 2: Database Migration
```bash
# Run migration script
mongosh cafe_pos < scripts/mongodb-handover-migration.js

# Verify migration
mongosh cafe_pos --eval "db.shifts.findOne({type:'WAITER'})"
mongosh cafe_pos --eval "db.cashier_shifts.findOne()"
```

- [ ] Migration script executed
- [ ] New fields added to shifts
- [ ] New fields added to cashier_shifts
- [ ] Sample documents verified

### Step 3: Create Indexes
```bash
# Create indexes
mongosh cafe_pos < scripts/mongodb-handover-indexes.js

# Verify indexes
mongosh cafe_pos --eval "db.cash_handovers.getIndexes()"
mongosh cafe_pos --eval "db.cash_discrepancies.getIndexes()"
```

- [ ] Indexes created for cash_handovers
- [ ] Indexes created for cash_discrepancies
- [ ] Index creation verified

### Step 4: Deploy Backend
```bash
# Build backend
cd backend
go build -o cafe-pos-server

# Stop current server
sudo systemctl stop cafe-pos-backend

# Replace binary
sudo cp cafe-pos-server /usr/local/bin/

# Start server
sudo systemctl start cafe-pos-backend

# Check status
sudo systemctl status cafe-pos-backend
```

- [ ] Backend built successfully
- [ ] Old server stopped
- [ ] New binary deployed
- [ ] Server started successfully
- [ ] Health check passed

### Step 5: Deploy Frontend
```bash
# Build frontend
cd frontend
npm run build

# Backup current frontend
sudo mv /var/www/cafe-pos /var/www/cafe-pos.backup

# Deploy new frontend
sudo cp -r dist /var/www/cafe-pos

# Restart web server
sudo systemctl restart nginx
```

- [ ] Frontend built successfully
- [ ] Old frontend backed up
- [ ] New frontend deployed
- [ ] Web server restarted
- [ ] Static files accessible

### Step 6: Verify Deployment
```bash
# Check backend health
curl http://localhost:8080/health

# Check frontend
curl http://localhost:3000

# Check API endpoints
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/cash-handovers/pending
```

- [ ] Backend health check passed
- [ ] Frontend loads correctly
- [ ] API endpoints responding
- [ ] No errors in logs

---

## 🧪 Post-Deployment Testing

### Smoke Tests (CRITICAL)
- [ ] Waiter can login
- [ ] Waiter can create handover
- [ ] Cashier can see pending handovers
- [ ] Cashier can confirm handover
- [ ] Cash amounts update correctly
- [ ] History displays correctly

### Functional Tests
- [ ] Partial handover works
- [ ] Handover and end shift works
- [ ] Discrepancy handling works
- [ ] Manager approval works
- [ ] Cancel handover works
- [ ] Quick confirm works

### Integration Tests
- [ ] Shift integration works
- [ ] Order integration works
- [ ] User authentication works
- [ ] Role authorization works

### Performance Tests
- [ ] Page load time < 2s
- [ ] API response time < 500ms
- [ ] No memory leaks
- [ ] Database queries optimized

### Mobile Tests
- [ ] Works on iOS Safari
- [ ] Works on Android Chrome
- [ ] Touch interactions work
- [ ] Modals display correctly
- [ ] Forms are usable

---

## 📊 Monitoring Setup

### Application Monitoring
- [ ] Error logging configured
- [ ] Performance monitoring active
- [ ] Uptime monitoring enabled
- [ ] Alert thresholds set

### Database Monitoring
- [ ] Query performance tracking
- [ ] Index usage monitoring
- [ ] Connection pool monitoring
- [ ] Disk space alerts

### Business Metrics
- [ ] Handover count tracking
- [ ] Discrepancy rate tracking
- [ ] Average processing time
- [ ] User adoption rate

---

## 👥 User Communication

### Before Deployment
- [ ] Notify users of maintenance window
- [ ] Send feature announcement
- [ ] Schedule training sessions
- [ ] Prepare support team

### During Deployment
- [ ] Display maintenance message
- [ ] Update status page
- [ ] Monitor support channels

### After Deployment
- [ ] Send deployment complete notification
- [ ] Share user guide link
- [ ] Announce training schedule
- [ ] Collect initial feedback

---

## 🆘 Rollback Plan

### If Critical Issues Found

**Step 1: Stop New Services**
```bash
sudo systemctl stop cafe-pos-backend
sudo systemctl stop nginx
```

**Step 2: Restore Database**
```bash
# Restore from backup
mongorestore --db cafe_pos backup_YYYYMMDD_HHMMSS/cafe_pos/
```

**Step 3: Restore Code**
```bash
# Restore backend
sudo cp /backup/cafe-pos-server.old /usr/local/bin/cafe-pos-server

# Restore frontend
sudo rm -rf /var/www/cafe-pos
sudo mv /var/www/cafe-pos.backup /var/www/cafe-pos
```

**Step 4: Restart Services**
```bash
sudo systemctl start cafe-pos-backend
sudo systemctl start nginx
```

**Step 5: Verify Rollback**
- [ ] Old version running
- [ ] Database restored
- [ ] Users can access system
- [ ] No data loss

---

## 📝 Post-Deployment Tasks

### Immediate (Day 1)
- [ ] Monitor error logs
- [ ] Check performance metrics
- [ ] Respond to user issues
- [ ] Document any problems

### Short-term (Week 1)
- [ ] Conduct user training
- [ ] Gather user feedback
- [ ] Fix minor bugs
- [ ] Optimize performance

### Medium-term (Month 1)
- [ ] Analyze usage patterns
- [ ] Review discrepancy trends
- [ ] Plan improvements
- [ ] Update documentation

---

## 📞 Support Contacts

### Technical Issues
- **Backend:** [Developer Name] - [Email] - [Phone]
- **Frontend:** [Developer Name] - [Email] - [Phone]
- **Database:** [DBA Name] - [Email] - [Phone]
- **DevOps:** [DevOps Name] - [Email] - [Phone]

### Business Issues
- **Product Owner:** [Name] - [Email] - [Phone]
- **Project Manager:** [Name] - [Email] - [Phone]

### Emergency
- **On-Call:** [Phone Number]
- **Escalation:** [Manager Phone]

---

## ✅ Sign-Off

### Development Team
- [ ] Backend Developer: _________________ Date: _______
- [ ] Frontend Developer: ________________ Date: _______
- [ ] QA Engineer: ______________________ Date: _______

### Management
- [ ] Technical Lead: ____________________ Date: _______
- [ ] Product Owner: ____________________ Date: _______
- [ ] Project Manager: ___________________ Date: _______

### Operations
- [ ] DevOps Engineer: ___________________ Date: _______
- [ ] Database Admin: ____________________ Date: _______

---

## 📊 Deployment Summary

**Deployment Date:** _______________  
**Deployment Time:** _______________  
**Downtime:** _______________  
**Issues Found:** _______________  
**Rollback Required:** Yes / No  
**Overall Status:** Success / Partial / Failed

**Notes:**
_____________________________________________
_____________________________________________
_____________________________________________

---

**Prepared by:** Development Team  
**Last Updated:** 2026-02-04  
**Version:** 1.0
