SHELL=/bin/bash
COMPOSE_ARGS=-f deployments/docker-compose.yaml --project-directory ./
REQUEST_COUNT=110

run:
	docker-compose $(COMPOSE_ARGS) down
	docker-compose $(COMPOSE_ARGS) up -d --build

rerun:
	docker-compose $(COMPOSE_ARGS) up -d --build

logs:
	docker-compose $(COMPOSE_ARGS) logs -f service

zk:
	docker-compose $(COMPOSE_ARGS) logs -f zk01 zk02 zk03

test-local:
	scripts/run.sh $(REQUEST_COUNT)