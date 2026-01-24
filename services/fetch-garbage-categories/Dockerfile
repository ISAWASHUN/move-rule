FROM golang:1.25.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap ./cmd

# Lambda用の軽量イメージ
FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap

CMD ["bootstrap"]
