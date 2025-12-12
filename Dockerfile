# ===== BUILD STAGE =====
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /http-server cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate cmd/migrate/main.go

# ===== RUN STAGE =====
FROM alpine:latest

WORKDIR /app

COPY --from=build /http-server /http-server
COPY --from=build /migrate /migrate
COPY --from=build /app/migrations/*.sql /migrations/
COPY --from=build /app/resources ./resources
COPY config.production.yaml ./config.production.yaml
COPY config.development.yaml ./config.development.yaml

EXPOSE 8080

RUN echo "\
    /migrate && /http-server \
    " > /app/entrypoint.sh

ENTRYPOINT ["sh", "/app/entrypoint.sh"]
