# Scripts Directory

This directory contains utility and testing scripts for the Cafe POS project.

## Deployment Scripts

- **build-push.sh** - Build and push Docker images to Docker Hub
- **deploy-ec2.sh** - Deploy application to AWS EC2
- **deploy-from-hub.sh** - Deploy from Docker Hub images

## Backend Management Scripts

- **restart-backend.sh** - Restart the backend server
- **fix-admin-role.sh** - Fix admin role in database
- **fix-admin-role-mongo.js** - MongoDB script to fix admin role

## Data Cleanup Scripts

- **clean-ingredients.sh** - Clean all ingredient data
- **clean-ingredients.js** - Node.js script to clean ingredients
- **clean-ingredients-mongo.js** - MongoDB script to clean ingredients
- **clean-incorrect-maintenance.js** - Clean incorrect maintenance records
- **clear-ingredient-localstorage.html** - HTML utility to clear localStorage

## Testing Scripts

### State Machine Tests
- **test-state-machine-validation.sh** - Validate state machine implementation
- **test-shift-state-machine.sh** - Test shift state machine
- **test-order-state-machine.sh** - Test order state machine

### Feature Tests
- **test-barista-shift-validation.sh** - Test barista shift validation
- **test-barista-simple.sh** - Simple barista workflow test
- **test-order-workflow-simple.sh** - Simple order workflow test
- **test-facility-maintenance.sh** - Test facility maintenance features
- **test-create-facility.sh** - Test facility creation
- **test-formatters-utility.sh** - Test formatter utilities
- **test-ingredient-endpoints.sh** - Test ingredient API endpoints

## Verification Scripts

- **verify-facilities-clean.sh** - Verify facilities data is clean

## Usage

Most scripts can be run directly from the project root:

```bash
# Example: Restart backend
./scripts/restart-backend.sh

# Example: Run tests
./scripts/test-order-state-machine.sh
```

Make sure scripts have execute permissions:
```bash
chmod +x scripts/*.sh
```
