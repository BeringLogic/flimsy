FROM alpine:3.21 AS builder 

RUN apk add --no-cache git go gcc
RUN apk add --no-cache lm-sensors tzdata

RUN mkdir -p /usr/src/flimsy
WORKDIR /usr/src/flimsy

RUN go install github.com/air-verse/air@latest

COPY ./src/ .
COPY ./src/static /var/lib/flimsy/static
COPY ./src/templates /var/lib/flimsy/templates

RUN go mod init github.com/BeringLogic/flimsy
RUN go mod tidy 

RUN /root/go/bin/air init

RUN mkdir /data
VOLUME /data

EXPOSE 8080

# CMD ["/root/go/bin/air", "--build.cmd", "go build -o bin/flimsy cmd/flimsy/main.go", "--build.bin", "./bin/flimsy", "-build.include_dir", "/var/lib/flimsy/static,/var/lib/flimsy/templates", "-build.stop_on_error", "true"]
CMD ["/root/go/bin/air", "--build.cmd", "go build -o bin/flimsy cmd/flimsy/main.go", "--build.bin", "./bin/flimsy", "-root", "../../..", "-build.include_dir", "usr/src/flimsy,var/lib/flimsy/static,var/lib/flimsy/templates", "-build.stop_on_error", "true"]
