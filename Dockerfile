# 使用官方 Golang base image
FROM golang:1.23.8-alpine AS builder


# 設定工作目錄
WORKDIR /app

# 安裝 git（copier、swag等用得到）
RUN apk add --no-cache git

# 複製 go.mod、go.sum 並預先下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製全部程式碼
COPY . .

# 編譯 binary
RUN go build -o server main.go

# 對外開放 port
EXPOSE 8080

# 啟動應用程式
CMD ["./server"]