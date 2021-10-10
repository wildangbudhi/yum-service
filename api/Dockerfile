# Build Stage
FROM golang:1.17.2-alpine3.14 as build

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -v -o main-app

# Run Stage
FROM alpine:3.14
COPY --from=build go/src/app app/

CMD ["./app/main-app"]