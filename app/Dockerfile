FROM composer:lts AS composer

WORKDIR /var/www
RUN composer require --ignore-platform-req=ext-gd league/color-extractor 0.4.0 

FROM php:8.4-fpm-alpine

RUN apk add --no-cache lm-sensors tzdata

RUN apk add --no-cache libpng-dev libjpeg-turbo-dev libwebp-dev libavif-dev libxpm-dev \
  && docker-php-ext-configure gd --with-avif --with-jpeg --with-xpm --with-webp \
  && docker-php-ext-install -j$(nproc) gd

RUN echo 'memory_limit = 512M' >> $PHP_INI_DIR/conf.d/docker-php-memlimit.ini \
  && echo 'upload_max_filesize = 200M' >> $PHP_INI_DIR/conf.d/docker-php-uploadsize.ini \
  && echo 'post_max_size = 200M' >> $PHP_INI_DIR/conf.d/docker-php-uploadsize.ini \
  && cp "$PHP_INI_DIR/php.ini-development" "$PHP_INI_DIR/php.ini"

RUN echo -e "[www]\naccess.log = /dev/null" >> /usr/local/etc/php-fpm.d/flimsy.conf

#RUN pecl install xdebug \
#  && docker-php-ext-enable xdebug

RUN adduser -H -D flimsy \
  && addgroup -S flimsy flimsy \
  && mkdir -p /data \
  && chown -R flimsy:flimsy /data
USER flimsy
VOLUME /data
WORKDIR /var/www

COPY ./src/ /var/www
COPY --from=composer /var/www/vendor /var/www/vendor

EXPOSE 9000
