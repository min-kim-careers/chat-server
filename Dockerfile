FROM golang:1.24-alpine

RUN go install github.com/air-verse/air@latest
RUN go install github.com/becheran/roumon@latest

RUN export PATH="$PATH:$HOME/go/bin"

WORKDIR /app

COPY ./app/go.mod ./app/go.sum ./

RUN go mod download

COPY ./app .

ENTRYPOINT [ "air", "-c", ".air.toml"]