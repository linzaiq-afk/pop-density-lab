# --- build stage -----------------------------------------------------------
FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

# --- runtime stage ---------------------------------------------------------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app
COPY --from=build /out/server ./server

USER app
EXPOSE 8081
ENV LISTEN_ADDR=:8081

ENTRYPOINT ["./server"]
