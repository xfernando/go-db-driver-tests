FROM golang:1.8

RUN apt-get update && apt-get install -y libpq-dev

WORKDIR /go/src/app
COPY ./ .
WORKDIR /go/src/app/sqltest
RUN go get -v ./...

CMD ["go", "test", "-v"]