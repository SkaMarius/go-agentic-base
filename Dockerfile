FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM alpine:3.20 AS runtime
RUN adduser -D -H -u 10001 appuser
COPY --from=build /out/server /usr/local/bin/server
USER appuser

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/server"]
