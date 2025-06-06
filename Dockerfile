FROM alpine:latest AS builder

RUN apk add --no-cache git go gcc

RUN mkdir -p /usr/src/flimsy
WORKDIR /usr/src/flimsy

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o bin/flimsy cmd/flimsy/main.go



FROM alpine:latest AS runner

RUN apk add --no-cache lm-sensors tzdata

COPY --from=builder /usr/src/flimsy/bin/flimsy /usr/local/bin/flimsy

RUN mkdir /data
VOLUME /data

CMD ["flimsy"]
