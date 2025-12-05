# =========================
# 1️⃣ Build Stage
# =========================
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/file-upload


# =========================
# 2️⃣ Runtime Stage
# =========================
FROM alpine:3.20

# ✅ Install certificates + timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# ✅ Set timezone manually to fix "unknown time zone" issue
ENV TZ=Asia/Kolkata
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
