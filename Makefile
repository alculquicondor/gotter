GOFLAGS := --ldflags '-w -linkmode external'
CC := $(shell which musl-clang)
STACK := gotter

all: images

images: accountservice/accountservice-linux-amd64 healthchecker/healthchecker-linux-amd64
	docker-compose build

accountservice/accountservice-linux-amd64: \
		accountservice/dbclient/boltclient.go \
		accountservice/model/account.go \
		accountservice/service/handlers.go \
		accountservice/service/router.go \
		accountservice/service/routes.go \
		accountservice/service/webserver.go \
		accountservice/main.go
	cd accountservice && \
	CC=${CC} go build ${GOFLAGS} -o accountservice-linux-amd64

healthchecker/healthchecker-linux-amd64: healthchecker/main.go
	cd healthchecker && \
	CC=${CC} go build ${GOFLAGS} -o healthchecker-linux-amd64

clean:
	rm -f \
	accountservice/accountservice-linux-amd64
	healthchecker/healthchecker-linux-amd64

deploy:
	docker deploy -c docker-compose.yml ${STACK}

rm_stack
	docker stack rm ${STACK}
