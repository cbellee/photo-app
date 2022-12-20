FROM golang:1.18.2-alpine3.15 AS builder

ARG SERVICE_NAME=""
ARG SERVICE_PORT=""

RUN mkdir /build
WORKDIR /build
COPY ./cmd .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./${SERVICE_NAME}/main.go

FROM scratch

ARG SERVICE_NAME=""
ARG SERVICE_PORT=""

WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/server .

# run server
EXPOSE $SERVICE_PORT
CMD ["./server"]