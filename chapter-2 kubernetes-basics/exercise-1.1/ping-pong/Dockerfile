# ===== Stage 1: Build Stage =====
FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o pingpong main.go


# ===== Stage 2: Production Stage =====

FROM alpine
WORKDIR /app
COPY --from=builder /app/pingpong .
CMD [ "./pingpong" ]