FROM golang:1.24.4-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /todo_list ./cmd/app/
RUN set -a && . /app/.env && set +a
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
EXPOSE 8080
CMD ["/todo_list"]
