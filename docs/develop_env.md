# Own's best practices to start local development environment
+ Start db
`docker-compose up -d db`
+ Start api
`docker-compose run -p 80:80 api bash`
+ Install vi
```
apt-get update
apt-get install vim
```
+ Install missing go packages (if any)
`go get <pacakge-location>`

# Debugging database
+ Start bash session of postgres container
```
# List all running container
sudo docker ps -a
# Get the ContainerID of postgres container, start a bash session of it
sudo docker exec -it <container-id> bash
# Start the psql to start interactive session sql of Postgres
psql -U postgres
```
+ Use pgAdmin
```
workon pgadmin
pgadmin
-> connect to `http://localhost:5050`
```
