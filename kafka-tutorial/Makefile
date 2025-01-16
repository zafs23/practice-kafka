# default make commands with no arguments
.DEFAULT_GOAL := help

export DOCKER_BUILDKIT=1

DOCKER_COMPOSE := docker-compose.yml

COMPOSE_TARGETS := up down build restart ps logs

$(COMPOSE_TARGETS): vet

up:
	docker compose -f $(DOCKER_COMPOSE) up -d

down:
	docker compose -f $(DOCKER_COMPOSE) down

build:
	docker compose -f $(DOCKER_COMPOSE) build

logs:
	docker compose -f $(DOCKER_COMPOSE) logs -f

ps:
	docker compose -f $(DOCKER_COMPOSE) ps

restart: down up

run-producer: up
	docker compose -f $(DOCKER_COMPOSE) up producer

run-consumer: up 
	docker compose -f $(DOCKER_COMPOSE) up consumer


vet: 
	@echo "Running go vet on producer..."
	cd producer && go vet -v ./...
	@echo "Running go test on consumer..."
	cd consumer && go test -v ./...

help:
	@echo "Usage:"
	@echo  "   --------------------------------------------------------"     
	@echo  "   make [target]"
	@echo  "   Targets: "
	@echo  "   up          Start all containers and run the kafka cluster"
	@echo  "   down        Stop and remove all containers"
	@echo  "   build       Build or rebuild all containers images (with BuildKit)"
	@echo  "   logs        Get all the logs from the containers"
	@echo  "   ps          Show containers status"
	@echo  "   restart     Restart containers (up + down)"
	@echo  "   vet         Run 'go vet' on producer and consumer"
	@echo  "   test        Run 'go test' on producer and consumer"
	@echo  "   help        make usage"
	@echo  "   --------------------------------------------------------"



