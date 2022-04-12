# syntax=docker/dockerfile:1

FROM golang:1.16-alpine
RUN apk add --update curl && \
    rm -rf /var/cache/apk/*
WORKDIR /app

COPY . .
COPY .env .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /subscriptions-app


EXPOSE 8080

CMD [ "/subscriptions-app" ]
