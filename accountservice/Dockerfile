FROM alpine:3.6

EXPOSE 6767
ADD accountservice/accountservice-linux-amd64 /
CMD ["/accountservice-linux-amd64", "-configServerUrl=http://configserver:8888", "-profile=test", "-configBranch=P8"]

ADD healthchecker/healthchecker-linux-amd64 /
HEALTHCHECK --interval=5s --timeout=5s CMD ["./healthchecker-linux-amd64", "-port=6767"]
