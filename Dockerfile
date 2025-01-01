FROM php:8.2-apache

COPY ./src /var/www/html
RUN chown -R www-data:www-data /var/www/html

RUN mkdir /data \
  && echo -n "{}" > /data/data.json \
  && chown -R www-data:www-data /data

EXPOSE 80
