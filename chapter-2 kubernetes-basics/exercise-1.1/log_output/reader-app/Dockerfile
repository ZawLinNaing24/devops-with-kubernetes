# Stage 1: Build Stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY main.go .

RUN go mod init reader-app
RUN go build -o reader main.go

# Stage 2: Run Stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/reader .

# Reader က Web Server အဖြစ် အလုပ်လုပ်မည်ဖြစ်၍ Port 3000 ကို ဖွင့်ထားခြင်း
EXPOSE 3000

CMD ["./reader"]