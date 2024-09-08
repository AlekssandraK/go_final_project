FROM golang:1.22.1

WORKDIR /app

ENV CGO_ENABLED=1
ENV GOOS=linux 
ENV TODO_PORT=7540
ENV TODO_DBFILE="scheduler.db"
ENV TODO_PASSWORD=0330

COPY . .

RUN go mod download

RUN  go build -o ./todo_app main.go

EXPOSE ${TODO_PORT}

CMD ["./todo_app"] 