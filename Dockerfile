FROM golang:1.23-alpine

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY ./app/go.mod ./app/go.sum ./

RUN go mod download

COPY ./app .

EXPOSE 8080

ENTRYPOINT [ "air", "-c", ".air.toml"]