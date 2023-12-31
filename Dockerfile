FROM golang:1.21-alpine as builder

WORKDIR /project/pcr-backend

COPY src/go.* .

RUN go mod download

COPY src .
RUN go build -o /project/pcr-backend/build/main main.go

FROM alpine:latest
COPY --from=builder /project/pcr-backend/build/main /app/build/main
COPY templates /app/templates
COPY migrations /app/migrations

EXPOSE 8080

ENTRYPOINT [ "/app/build/main" ]