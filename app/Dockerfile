FROM golang:1.23 AS builder 

# RUN echo 'memory_limit = 512M' >> $PHP_INI_DIR/conf.d/docker-php-memlimit.ini \
#   && echo 'upload_max_filesize = 200M' >> $PHP_INI_DIR/conf.d/docker-php-uploadsize.ini \
#   && echo 'post_max_size = 200M' >> $PHP_INI_DIR/conf.d/docker-php-uploadsize.ini \

RUN mkdir -p /usr/src/flimsy
WORKDIR /usr/src/flimsy

# use cached go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY ./src .

# use this instead to create or update go.mod and go.sum
# COPY ./src .
# RUN go mod init BeringLogic/flimsy \
#   && go mod tidy


ENV CGO_ENABLED=0
RUN go build -v -o /usr/local/bin/flimsy



FROM alpine:3.21 AS runner

COPY --from=builder /usr/local/bin/flimsy /usr/local/bin/flimsy

RUN mkdir /data
VOLUME /data

EXPOSE 8080

CMD ["flimsy"]
