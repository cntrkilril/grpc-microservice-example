FROM golang:1.20.2-alpine3.17 AS builder

ENV CGO_ENABLED 0

ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata git

WORKDIR /build

ADD go.mod .

ADD go.sum .

ARG GITPASSWORD
ARG GITUSERNAME
RUN go env GOPRIVATE="gitlab.com/d1zero-online-booking/*" && echo "machine gitlab.com" > ~/.netrc &&\
    echo "login $GITUSERNAME" >> ~/.netrc && echo "password $GITPASSWORD" >>  ~/.netrc

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/main cmd/app/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /build/db /app/db

CMD ["./main"]
