.PHONY: help build up down restart logs clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

logs: ## Show logs from all services
	docker-compose logs -f

logs-backend: ## Show backend logs
	docker-compose logs -f backend

logs-frontend: ## Show frontend logs
	docker-compose logs -f frontend

logs-mongodb: ## Show MongoDB logs
	docker-compose logs -f mongodb

clean: ## Remove all containers, volumes, and images
	docker-compose down -v
	docker system prune -af

ps: ## Show running containers
	docker-compose ps

exec-backend: ## Execute shell in backend container
	docker-compose exec backend sh

exec-mongodb: ## Execute MongoDB shell
	docker-compose exec mongodb mongosh -u admin -p admin123

seed-admin: ## Seed admin user only
	cd backend && go run cmd/seed-admin/main.go

build-backend: ## Build only backend image
	docker-compose build backend

build-frontend: ## Build only frontend image
	docker-compose build frontend

deploy: build up ## Build and deploy all services
	@echo "✅ Deployment complete!"
	@echo "Frontend: http://localhost"
	@echo "Backend API: http://localhost:8080"
	@echo "MongoDB: mongodb://localhost:27017"
