FROM golang:1.22.1

ENV CGO_ENABLED=1\
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

RUN go mod download

COPY . .

RUN  go build -o /todo_app main.go

EXPOSE ${TODO_PORT}

CMD ["/todo_app"] 