# Stage 1: Build Stage (Go Code ကို Compile လုပ်မည့်အပိုင်း)
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Go ဖိုင်ကို Copy ကူးယူခြင်း
COPY main.go .

# လိုအပ်သော Packages များ (uuid) ကို ဒေါင်းလုဒ်ဆွဲရန် Go Module စတင်ခြင်း
RUN go mod init writer-app
RUN go get github.com/google/uuid

# 'writer' အမည်ဖြင့် Binary File အဖြစ် Compile လုပ်ခြင်း
RUN go build -o writer main.go

# Stage 2: Run Stage (တကယ် Run မည့် ပေါ့ပါးသော Image)
FROM alpine:latest
WORKDIR /app

# Stage 1 မှ Compile လုပ်ပြီးသား Binary File ကိုသာ ကူးယူခြင်း
COPY --from=builder /app/writer .

# Application ကို စတင် မောင်းနှင်ခြင်း
CMD ["./writer"]