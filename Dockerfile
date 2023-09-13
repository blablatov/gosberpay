FROM golang:1.20

RUN git clone https://github.com/blablatov/gosberpay.git
WORKDIR gosberpay

RUN go mod download

COPY *.go ./
COPY *.conf ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /gosberpay

EXPOSE 8443

CMD ["/gosberpay"]