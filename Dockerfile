FROM golang:alpine AS backend-app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o nusastra cmd/app/main.go

FROM alpine:latest AS prod

WORKDIR /app

COPY --from=backend-app . .

EXPOSE 8080
ENTRYPOINT ["./app/nusastra"]