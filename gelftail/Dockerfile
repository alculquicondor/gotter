FROM alpine:3.6

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* /tmp/*

EXPOSE 12202/udp
ADD gelftail-linux-amd64 /
ADD token.txt /
CMD ["/gelftail-linux-amd64", "-port=12202"]
