FROM php:8.2-apache

RUN apt-get update \
  && apt-get install -y libavif-dev libjpeg62-turbo-dev libpng-dev libxpm-dev libwebp-dev lm-sensors \
  && docker-php-ext-configure gd --with-avif --with-jpeg --with-xpm --with-webp \
  && docker-php-ext-install -j$(nproc) gd

COPY ./src /var/www/html
COPY ./color-extractor/src /var/www/html/vendor/League/ColorExtractor
RUN chown -R www-data:www-data /var/www/html

RUN mkdir /data \ 
  && mkdir /data/icons \
  && mkdir /data/backgrounds \
  && chown -R www-data:www-data /data \
  && echo 'Alias /dashboard-icons/ /data/icons/\n<Directory /data/icons/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>\n' >> /etc/apache2/conf-available/docker-php.conf \
  && echo 'Alias /backgrounds/ /data/backgrounds/\n<Directory /data/backgrounds/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>\n' >> /etc/apache2/conf-available/docker-php.conf

RUN echo 'memory_limit = 512M' >> /usr/local/etc/php/conf.d/docker-php-memlimit.ini

EXPOSE 80
