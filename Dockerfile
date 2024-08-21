# Minimal va kerakli paketlarni o'rnatish uchun build-base ni qo'shamiz
FROM golang:1.22.5-alpine

# Build asboblarini o'rnatish
RUN apk add --no-cache build-base

WORKDIR /app

# go.mod va go.sum fayllarini avval nusxalab, modullarni yuklab olish
COPY go.mod go.sum ./
RUN go mod download

# Qolgan barcha fayllarni nusxalash
COPY . .

# Binariy faylni build qilish
RUN go build -o main .

# 8080 portini ochamiz
EXPOSE 8080

# Bajariladigan faylni ishga tushirish
CMD ["./main"]
