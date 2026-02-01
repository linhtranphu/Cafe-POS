# Best Practices Guide

This folder contains documented best practices for the Cafe POS project, covering frontend, backend, and deployment patterns.

## Frontend Best Practices

### Form & Input Fields
- **[DATE_INPUT_FIELDS.md](./DATE_INPUT_FIELDS.md)** - Date input styling and mobile optimization
  - Use `appearance: none` to normalize browser styling
  - Apply `px-3 py-3 text-sm` for consistent sizing
  - Prevents width overflow on mobile devices

### Layout & Responsive Design
- Mobile-first approach with Tailwind CSS
- Use `h-screen w-screen overflow-hidden flex flex-col` for full-screen mobile layouts
- Fixed headers with `flex-shrink-0`
- Scrollable content with `flex-1 overflow-y-auto`
- Responsive grids: `grid-cols-1 sm:grid-cols-2`

### Component Patterns
- Slide-in modals from right side for mobile UX
- Bottom navigation for mobile navigation
- Inline forms instead of popups for better mobile experience
- Use hash-based routing (`createWebHashHistory`) for SPAs

## Backend Best Practices

### Environment Configuration
- Use `.env` file for local development
- Export environment variables before running backend:
  - `MONGODB_URI=mongodb://admin:password123@localhost:27017`
  - `MONGODB_DATABASE=cafe_pos`
  - `JWT_SECRET=your-jwt-secret-key-min-32-chars-long`
- Keep MongoDB credentials consistent between `MONGO_INITDB_ROOT_PASSWORD` and `MONGODB_URI`

### Database
- MongoDB runs in Docker container for consistency
- Use `restart_local.sh` to manage local development environment
- Always check MongoDB connection before starting backend

## Deployment Best Practices

### Docker
- Always use `--no-cache` flag when building images
- Use `docker-compose.yml` for production (Docker Hub images)
- Use `docker-compose.local.yml` for local development
- Tag images with version numbers and `latest`

### Port Configuration
- Backend: port 3000
- Frontend: port 5173 (Vite dev server)
- MongoDB: port 27017
- Nginx (production): port 80/443

## Quick Reference

### Common Tasks

**Start local development:**
```bash
./restart_local.sh
```

**Build and push Docker images:**
```bash
./build_docker_hub.sh
```

**Build and run locally:**
```bash
./build_docker_hub_local.sh
```

**Update admin password:**
```bash
cd backend/cmd/update-admin-password
go run main.go
```

## Contributing

When adding new best practices:
1. Create a new markdown file in this folder
2. Document the problem, solution, and examples
3. Update this README with a link to the new file
4. Include implementation examples and testing checklist
