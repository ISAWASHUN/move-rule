FROM golang:1.25.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd

# Lambda用のprovidedイメージ（Lambda Runtime Interface Client含む）
FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap

# ENTRYPOINTはベースイメージに含まれているため、CMDのみ指定
CMD ["bootstrap"]
