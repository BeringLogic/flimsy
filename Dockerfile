FROM alpine:3.21 AS builder 

RUN apk add --no-cache git go gcc

# That stuff will probably need [this](https://gin-gonic.com/docs/examples/custom-http-config/)
# RUN echo 'memory_limit = 512M' >> $PHP_INI_DIR/conf.d/docker-php-memlimit.ini \
#   && echo 'upload_max_filesize = 200M' >> $PHP_INI_DIR/conf.d/docker-php-uploadsize.ini \
#   && echo 'post_max_size = 200M' >> $PHP_INI_DIR/conf.d/docker-php-uploadsize.ini \

RUN mkdir -p /usr/src/flimsy
WORKDIR /usr/src/flimsy

# use cached go.mod and go.sum
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download && go mod verify

COPY ./src/ .

# use this instead to create go.mod and go.sum
# RUN go mod init github.com/BeringLogic/flimsy

# use this instead to update go.mod and go.sum
# RUN go mod tidy


RUN go build -v -o /usr/local/bin/flimsy ./cmd/flimsy



FROM alpine:3.21 AS runner

RUN apk add --no-cache lm-sensors tzdata

COPY --from=builder /usr/local/bin/flimsy /usr/local/bin/flimsy
COPY --from=builder /usr/src/flimsy/static /var/lib/flimsy/static
COPY --from=builder /usr/src/flimsy/templates /var/lib/flimsy/templates

RUN mkdir /data
VOLUME /data

EXPOSE 8080

CMD ["flimsy"]
