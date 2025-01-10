FROM php:8.2-apache

WORKDIR /var/www/html

RUN apt-get update \
  && apt-get install -y --no-install-recommends libavif-dev libjpeg62-turbo-dev libpng-dev libxpm-dev libwebp-dev lm-sensors \
  && docker-php-ext-configure gd --with-avif --with-jpeg --with-xpm --with-webp \
  && docker-php-ext-install -j$(nproc) gd \
  && echo 'memory_limit = 512M' >> /usr/local/etc/php/conf.d/docker-php-memlimit.ini \
  && rm -rf /var/lib/apt/lists/* \
  && mkdir /data \ 
  && mkdir /data/icons \
  && mkdir /data/backgrounds \
  && chown -R www-data:www-data /data \
  && echo 'Alias /dashboard-icons/ /data/icons/\n<Directory /data/icons/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>\n' >> /etc/apache2/conf-available/docker-php.conf \
  && echo 'Alias /backgrounds/ /data/backgrounds/\n<Directory /data/backgrounds/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>\n' >> /etc/apache2/conf-available/docker-php.conf

COPY ./src /var/www/html
COPY ./color-extractor/src /var/www/html/vendor/League/ColorExtractor

EXPOSE 80
