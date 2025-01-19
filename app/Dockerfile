FROM php:8.4-fpm-alpine

WORKDIR /var/www/html

RUN apk add --no-cache libpng-dev libjpeg-turbo-dev libwebp-dev libavif-dev libxpm-dev lm-sensors \
  && docker-php-ext-configure gd --with-avif --with-jpeg --with-xpm --with-webp \
  && docker-php-ext-install -j$(nproc) gd \
  && echo 'memory_limit = 512M' >> $PHP_INI_DIR/conf.d/docker-php-memlimit.ini \
  && mv "$PHP_INI_DIR/php.ini-development" "$PHP_INI_DIR/php.ini"

#RUN pecl install xdebug \
#  && docker-php-ext-enable xdebug

#RUN apk add --no-cache nginx

#RUN echo 'Alias /dashboard-icons/ /data/icons/\n<Directory /data/icons/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>\n' >> /etc/apache2/conf-available/docker-php.conf \
#  && echo 'Alias /backgrounds/ /data/backgrounds/\n<Directory /data/backgrounds/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>\n' >> /etc/apache2/conf-available/docker-php.conf

COPY ./src /var/www
COPY ./color-extractor/src /var/www/vendor/League/ColorExtractor

EXPOSE 9000
