FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/.

RUN ls -l /app

CMD [ "/main" ]
