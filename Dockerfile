FROM php:8.2-apache

COPY ./src /var/www/html
RUN chown -R www-data:www-data /var/www/html

RUN mkdir -p /data/icons \
  && echo -n "{}" > /data/data.json \
  && chown -R www-data:www-data /data

RUN echo 'Alias /dashboard-icons/ /data/icons/\n<Directory /data/icons/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>' >> /etc/apache2/conf-available/docker-php.conf

EXPOSE 80
