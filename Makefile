up:
	docker-compose up

test:
	docker-compose -f docker-compose.tests.yml up --build
	docker wait avito-e2e-1
	docker logs avito-e2e-1
	docker-compose -f docker-compose.tests.yml down -v