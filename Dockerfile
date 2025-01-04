FROM php:8.2-apache

COPY ./src /var/www/html
RUN chown -R www-data:www-data /var/www/html

RUN mkdir /data \ 
  && chown -R www-data:www-data /data \
  && mkdir /data/icons \
  && mkdir /data/backgrounds \
  && echo 'Alias /dashboard-icons/ /data/icons/\n<Directory /data/icons/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>\n' >> /etc/apache2/conf-available/docker-php.conf \
  && echo 'Alias /backgrounds/ /data/backgrounds/\n<Directory /data/backgrounds/>\n\tOrder allow,deny\n\tAllow from all\nRequire all granted\n</Directory>\n' >> /etc/apache2/conf-available/docker-php.conf

EXPOSE 80
