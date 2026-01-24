FROM golang:1.25.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o db-migrator

FROM gcr.io/distroless/base-debian12
WORKDIR /app

USER 65532

COPY --from=builder /app/db-migrator .

ENTRYPOINT ["./db-migrator"]
