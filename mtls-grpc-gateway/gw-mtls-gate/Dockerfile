# Multi-stage gw-mtls-client build
# Многоэтапная сборка gw-mtls-client

FROM golang AS build

ENV location /go/src/github.com/blablatov/mtls-grpc-gateway

WORKDIR ${location}/gw-mtls-client

ADD ./main.go ${location}/gw-mtls-client

RUN go mod init github.com/blablatov/mtls-grpc-gateway/gw-mtls-client

RUN CGO_ENABLED=0 go build -o gw-mtls-client

# Go binaries are self-contained executables. Используя директиву FROM scratch - 
# Go образы  не должны содержать ничего, кроме одного двоичного исполняемого файла.
FROM scratch
COPY --from=build ./gw-mtls-client ./gw-mtls-client

ENTRYPOINT ["./gw-mtls-client"]
EXPOSE 50051