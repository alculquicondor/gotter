GOFLAGS := --ldflags '-w -linkmode external'
CC := $(shell which musl-clang)
STACK := gotter

all: images

images: accountservice/accountservice-linux-amd64 \
		healthchecker/healthchecker-linux-amd64 \
		vipservice/vipservice-linux-amd64 \
		gelftail/gelftail-linux-amd64 \
		gelftail/token.txt \
		support/config-server/build/libs/config-server-0.0.1-SNAPSHOT.jar \
		docker-compose.yml
	docker-compose build


accountservice/accountservice-linux-amd64: \
		accountservice/dbclient/boltclient.go \
		accountservice/model/account.go \
		accountservice/model/healthcheck.go \
		accountservice/model/vipnotification.go \
		accountservice/service/handlers.go \
		accountservice/service/router.go \
		accountservice/service/routes.go \
		accountservice/service/webserver.go \
		accountservice/main.go \
		common/config/events.go \
		common/config/loader.go \
		common/messaging/messagingclient.go \
		common/netutils/utils.go
	cd accountservice && \
	CC=${CC} go build ${GOFLAGS} -o accountservice-linux-amd64


vipservice/vipservice-linux-amd64: \
		vipservice/service/router.go \
		vipservice/service/routes.go \
		vipservice/service/webserver.go \
		vipservice/main.go \
		common/config/events.go \
		common/config/loader.go \
		common/messaging/messagingclient.go
	cd vipservice && \
	CC=${CC} go build ${GOFLAGS} -o vipservice-linux-amd64

healthchecker/healthchecker-linux-amd64: healthchecker/main.go
	cd healthchecker && \
	CC=${CC} go build ${GOFLAGS} -o healthchecker-linux-amd64

gelftail/gelftail-linux-amd64: \
		gelftail/aggregator/aggregator.go \
		gelftail/transformer/transformer.go \
		gelftail/transformer/transformer.go \
		gelftail/gelftail.go
	cd gelftail && \
	CC=${CC} go build ${GOFLAGS} -o gelftail-linux-amd64


support/config-server/build/libs/config-server-0.0.1-SNAPSHOT.jar: \
		support/config-server/src/main/resources/application.yml
	cd support/config-server && \
	./gradlew build


loadtest/loadtest: loadtest/main.go
	cd loadtest && \
	CC=${CC} go build ${GOFLAGS} -o loadtest


clean:
	rm -f \
	accountservice/accountservice-linux-amd64
	healthchecker/healthchecker-linux-amd64

deploy:
	docker stack deploy -c docker-compose.yml ${STACK}

rm_stack:
	docker stack rm ${STACK}
