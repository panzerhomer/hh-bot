FROM golang:1.20 
RUN go install github.com/cosmtrek/air@latest
WORKDIR /usr/src/app
COPY . .
RUN go mod tidy