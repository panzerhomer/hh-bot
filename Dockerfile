# FROM golang:latest

# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o main .

# EXPOSE 3000

# ENTRYPOINT ["go", "run", "./main.go"]

FROM golang:1.17 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev

RUN go build -o /server

FROM debian:buster

EXPOSE 3000

WORKDIR /
COPY --from=build-env /server /

CMD ["/server"]