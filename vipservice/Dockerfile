FROM alpine:3.6

EXPOSE 6868
ADD vipservice/vipservice-linux-amd64 /
CMD ["/vipservice-linux-amd64", "-configServerUrl=http://configserver:8888", "-profile=test", "-configBranch=P9"] || exit 1

ADD healthchecker/healthchecker-linux-amd64 /
HEALTHCHECK --interval=5s --timeout=5s CMD ["./healthchecker-linux-amd64", "-port=6868"]
