FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server ./cmd

FROM alpine:latest AS final
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /bin/server /bin/server
EXPOSE 8080
ENTRYPOINT ["/bin/server"]
