FROM golang:1.23 

# RUN echo 'memory_limit = 512M' >> $PHP_INI_DIR/conf.d/docker-php-memlimit.ini \
#   && echo 'upload_max_filesize = 200M' >> $PHP_INI_DIR/conf.d/docker-php-uploadsize.ini \
#   && echo 'post_max_size = 200M' >> $PHP_INI_DIR/conf.d/docker-php-uploadsize.ini \

RUN addgroup flimsy \
  && adduser --no-create-home --disabled-password --ingroup flimsy flimsy \
  && mkdir -p /data \
  && chown -R flimsy:flimsy /data
VOLUME /data

EXPOSE 8080

RUN mkdir -p /usr/src/flimsy
WORKDIR /usr/src/flimsy

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./src .
RUN go build -v -o /usr/local/bin/flimsy

USER flimsy
CMD ["flimsy"]
