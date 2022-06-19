### STAGE 1
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download

RUN go build -o main /app/

## STAGE 2
FROM alpine AS production

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app /app

ENTRYPOINT ["/app/main"]


