FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o travelsphere main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/travelsphere .
COPY --from=builder /app/conf ./conf
COPY --from=builder /app/views ./views
COPY --from=builder /app/static ./static

EXPOSE 8080

ENV RUNMODE=prod

CMD ["./travelsphere"]