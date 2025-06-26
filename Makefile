build:
	docker build -t shop_db .

create:
	docker run -d --name shop_db -e POSTGRES_PASSWORD=postgres -p 8091:5432 -v pgdata:/var/lib/postgresql/data shop_db

start:
	docker start shop_db

stop:
	docker stop shop_db

remove:
	docker rm shop_db

see:
	docker ps -a

createdb:
	docker exec -it shop_db createdb -U postgres --owner=postgres simple_shop
removedb:
	docker exec -it shop_db dropdb -U postgres -f simple_shop

.PHONY: cover
cover:
	go test -count=1 -race -coverprofile=coverage ./...
	go tool cover -html=coverage
	timeout /t 5 /nobreak >nul
	del coverage

test:
	go test -count=1 ./...

gen:
	mockgen -source=./internal/repository/repository.go -destination=./internal/repository/mocks/mock_repository.go
	mockgen -source=./internal/repository/order.go -destination=./internal/repository/mocks/mock_order.go