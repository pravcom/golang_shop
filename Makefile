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
