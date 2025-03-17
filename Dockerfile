FROM alpine:3.21 AS builder

RUN apk add --no-cache git go gcc

RUN mkdir -p /usr/src/flimsy
WORKDIR /usr/src/flimsy

COPY ./src/ .

RUN go mod init github.com/BeringLogic/flimsy
RUN go mod tidy

RUN go build -o bin/flimsy cmd/flimsy/main.go



FROM alpine:3.21 AS runner

RUN apk add --no-cache lm-sensors tzdata

COPY --from=builder /usr/src/flimsy/bin/flimsy /usr/local/bin/flimsy
COPY --from=builder /usr/src/flimsy/static /var/lib/flimsy/static
COPY --from=builder /usr/src/flimsy/templates /var/lib/flimsy/templates

RUN mkdir /data
VOLUME /data

CMD ["flimsy"]
