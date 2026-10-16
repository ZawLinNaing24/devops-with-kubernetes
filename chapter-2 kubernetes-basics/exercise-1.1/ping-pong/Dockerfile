# ===== Stage 1: Build Stage =====
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# ၁။ go.mod နဲ့ go.sum ကို အရင် COPY လုပ်ပါ (Dependency တွေကို Cache လုပ်ဖို့အတွက်ပါ)
COPY go.mod go.sum ./

# ၂။ Dependencies တွေကို အရင် ဆွဲယူပါ (ကုဒ်အပြောင်းအလဲမရှိရင် ဒီအဆင့်က Cache ကနေ ချက်ချင်း ပြီးသွားပါမယ်)
RUN go mod download

# ၃။ မှီခိုနေတဲ့ ဖိုင်တွေ အားလုံးကို နောက်မှ COPY ကူးပါ
COPY . .

# ၄ - ပိုမိုလုံခြုံပြီး ပေါ့ပါးသော Static Binary အဖြစ် Compile လုပ်ပါ
RUN CGO_ENABLED=0 GOOS=linux go build -o pingpong main.go


# ===== Stage 2: Production Stage =====
FROM alpine:latest

WORKDIR /app

# ၅။ Builder Stage ထဲက Compile လုပ်ပြီးသား Binary ဖိုင်ကိုသာ သန့်သန့်ရှင်းရှင်း ကူးယူပါ
COPY --from=builder /app/pingpong .

# ၆။ App ကို စတင် Run ပါ
CMD ["./pingpong"]