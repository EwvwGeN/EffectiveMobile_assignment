FROM golang:1.21.5-alpine AS builder
WORKDIR /app
RUN apk update --no-cache \
    && apk add --no-cache \
        git 
COPY cmd/ ./cmd/
COPY internal/ ./internal
COPY http/ ./http
COPY go.mod .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go generate ./cmd/server/.
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/server/

FROM alpine AS runner
WORKDIR /
COPY --from=builder /app/main /main

ENTRYPOINT ["./main"]