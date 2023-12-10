FROM golang:1.20

RUN git clone https://github.com/blablatov/gosberpay.git
WORKDIR gosberpay

RUN go mod download

COPY *.go ./
COPY *.conf ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /gosberpay
EXPOSE 8443

ENV DIRPATH=/mtls-grpc-gateway
WORKDIR $DIRPATH/gw-mtls-gate
RUN CGO_ENABLED=0 GOOS=linux go build -o $DIRPATH/gw-mtls-gate
EXPOSE 8444


WORKDIR $DIRPATH/gw-mtls-service
RUN CGO_ENABLED=0 GOOS=linux go build -o $DIRPATH/gw-mtls-service
EXPOSE 50051
CMD ["/gosberpay", "/gw-mtls-gate", "/gw-mtls-service"]