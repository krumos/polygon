createdb: 
	sudo docker run -d --name PolygonDB -p 5432:5432 -e POSTGRES_USERNAME=postgres -e POSTGRES_PASSWORD=password postgres:13.5

psql:
	sudo docker exec -it PolygonDB psql -U postgres postgres

creatermq:
	sudo docker run -d --hostname request-host -e RABBITMQ_DEFAULT_USER=guest -e RABBITMQ_DEFAULT_PASS=guest --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
