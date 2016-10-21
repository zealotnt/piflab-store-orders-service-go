# PIFLab Store API
[![CircleCI](https://circleci.com/gh/zealotnt/piflab-store-api-go.svg?style=svg)](https://circleci.com/gh/zealotnt/piflab-store-api-go)  
[![Coverage Status](https://coveralls.io/repos/github/zealotnt/piflab-store-api-go/badge.svg)](https://coveralls.io/github/zealotnt/piflab-store-api-go)  

## API Docs
http://docs.piflabstore.apiary.io/

## Dependencies

- **GO 1.5**

## 3rd parties

## Framework

- **Dependency**: [Godep](https://github.com/tools/godep)
- **Router**: [Gorilla Mux](https://github.com/gorilla/mux)

## Build Docker image

`docker build -t piflab-store-api-go .`

## Run a command inside container

``docker run -v `pwd`:/go/src/github.com/o0khoiclub0o/piflab-store-api-go piflab-store-api-go <command you want to run>``

## Add package

- ``docker run -it -v `pwd`:/go/src/github.com/o0khoiclub0o/piflab-store-api-go piflab-store-api-go bash``
- `go get <package>`
- `import "<package>"`
- `godep save ./...`

## Development

``docker run -p 80:80 -v `pwd`:/go/src/github.com/o0khoiclub0o/piflab-store-api-go --env-file=.env --rm -it piflab-store-api-go bash``

``docker run -p 80:80 -v `pwd`:/go/src/github.com/o0khoiclub0o/piflab-store-api-go --env-file=.env --rm  piflab-store-api-go gin -p 80 run``

## Testing

``docker run -v `pwd`:/go/src/github.com/o0khoiclub0o/piflab-store-api-go piflab-store-api-go ginkgo -r``

## Migration

### Migrate
`goose up`

### Rollback
`goose down`

### Seed
`go run db/seeds/main.go`

## Docker-compose

### Run entire app
`docker-compose up`

### Manually start database
`docker-compose up -d db`

### Manually run piflab-store-api command
`docker-compose run -p 80:80 api <command you want to run>`

### Practical Docker-compose use
- cd to root of project directory
- Start DB `docker-compose up -d db`
- Start API `docker-compose run -p 80:80 api bash`
- Start bash instance of DB `sudo docker exec -it <container_id> bash`

### Golang dump
```
import "github.com/davecgh/go-spew/spew"
spew.Dump(structure_to_be_dump)
```

### Save package
If use `godep save .` it will remove the pacakges that aren't used in source code (such as `ginkgo`, `gomega`, ...), so use `godep save ./...` save only (no remove).