# ===== BUILD STAGE =====
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /http-server

# ===== RUN STAGE =====
FROM alpine:latest

WORKDIR /app

COPY --from=build /http-server /http-server
COPY config.production.yaml ./config.production.yaml
COPY config.development.yaml ./config.development.yaml

EXPOSE 8080

ENTRYPOINT ["/http-server"]
