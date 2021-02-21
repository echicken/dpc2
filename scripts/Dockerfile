# Build
FROM golang:1.15-alpine3.12 AS build

WORKDIR /src
ADD . .
RUN mkdir /app
RUN go mod vendor
RUN go build -o /app/doorparty-connector ./cmd/doorparty-connector.go
RUN chmod +x /app/doorparty-connector

# Play
FROM alpine:3.12 AS runtime
WORKDIR /app
COPY --from=build /app /app

CMD ["/app/doorparty-connector"]