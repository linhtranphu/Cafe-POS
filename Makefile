.PHONY: help build up down restart restart-no-bridge docker-restart logs clean dev dev-stop dev-backend dev-frontend \
        docker-build docker-build-backend docker-build-frontend docker-build-bridge \
        docker-push docker-push-backend docker-push-push-frontend docker-push-bridge \
        docker-publish docker-publish-backend docker-publish-frontend docker-publish-bridge \
        docker-clean

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

restart: ## Restart local dev environment (MongoDB + backend + frontend)
	@bash restart_local.sh

restart-no-bridge: ## Restart local dev without Print Bridge
	@bash restart_local.sh --no-bridge

docker-restart: ## Restart all Docker services
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

# ── Local Development ─────────────────────────────────────────────────────────

dev: ## Start local dev environment (MongoDB + backend + frontend)
	@bash restart_local.sh

dev-no-bridge: ## Start local dev without Print Bridge
	@bash restart_local.sh --no-bridge

dev-stop: ## Stop local dev (backend + frontend)
	@bash stop_local.sh

dev-backend: ## Build and run backend only (local)
	@bash start-backend.sh

dev-frontend: ## Start frontend dev server only
	@cd frontend && npm run dev -- --host

dev-logs-backend: ## Tail backend log
	@tail -f backend.log

dev-logs-frontend: ## Tail frontend log
	@tail -f frontend.log

dev-logs-mongo: ## Tail MongoDB container log
	@docker logs -f cafe-pos-mongodb

dev-logs-bridge: ## Tail Print Bridge container log
	@docker logs -f local-print-bridge

# ── Docker Hub Build & Push ───────────────────────────────────────────────────

DOCKER_USER := linhtranphu

docker-build-backend: ## Build backend Docker image
	cd backend && docker build --no-cache -t $(DOCKER_USER)/cafe-pos-backend:latest .

docker-build-frontend: ## Build frontend Docker image
	cd frontend && docker build --no-cache -t $(DOCKER_USER)/cafe-pos-frontend:latest .

docker-build-bridge: ## Build Print Bridge Docker image
	cd local-print-bridge && docker build --no-cache -t $(DOCKER_USER)/local-print-bridge:latest .

docker-build: docker-build-backend docker-build-frontend ## Build backend + frontend images

docker-build-all: docker-build-backend docker-build-frontend docker-build-bridge ## Build all images

docker-push-backend: ## Push backend image to Docker Hub
	docker push $(DOCKER_USER)/cafe-pos-backend:latest

docker-push-frontend: ## Push frontend image to Docker Hub
	docker push $(DOCKER_USER)/cafe-pos-frontend:latest

docker-push-bridge: ## Push Print Bridge image to Docker Hub
	docker push $(DOCKER_USER)/local-print-bridge:latest

docker-push: docker-push-backend docker-push-frontend ## Push backend + frontend to Docker Hub

docker-push-all: docker-push-backend docker-push-frontend docker-push-bridge ## Push all images to Docker Hub

docker-publish: docker-build docker-push docker-clean ## Build and push backend + frontend
	@echo "✅ Backend + Frontend published to Docker Hub"

docker-publish-backend: docker-build-backend docker-push-backend docker-clean ## Build and push backend only
	@echo "✅ Backend published to Docker Hub"

docker-publish-frontend: docker-build-frontend docker-push-frontend docker-clean ## Build and push frontend only
	@echo "✅ Frontend published to Docker Hub"

docker-publish-bridge: docker-build-bridge docker-push-bridge docker-clean ## Build and push Print Bridge only
	@echo "✅ Print Bridge published to Docker Hub"

docker-publish-all: docker-build-all docker-push-all docker-clean ## Build and push all images
	@echo "✅ All images published to Docker Hub"

docker-clean: ## Remove dangling images, unused cache, and stopped containers
	docker image prune -f
	docker builder prune -f
	docker container prune -f
	docker system df
