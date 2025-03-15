FROM golang:1.22-alpine as build

WORKDIR /app
COPY . .

RUN apk add --update make npm nodejs gcc musl-dev upx &&\
    npm -D install tailwindcss@3 &&\
    make build

FROM alpine:latest as run

WORKDIR /app
COPY --from=build /app/bin/qpoll ./qpoll

RUN mkdir -p /app/db
EXPOSE 8080

CMD ["./qpoll"]
