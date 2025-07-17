FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
RUN go build -o /todo_list ./cmd/app/
FROM alpine:latest
WORKDIR /
COPY --from=builder /todo_list .
RUN apk --no-cache add sqlite-libs
EXPOSE 8080
VOLUME [ "/data" ]
CMD ["./todo_list"]