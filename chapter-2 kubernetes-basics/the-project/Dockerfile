# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app main.go

# Stage 2: Final Image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/todo-app .
ENV PORT=3000
EXPOSE 3000
CMD ["./todo-app"]