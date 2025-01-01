FROM php:8.2-apache

COPY ./src /var/www/html

RUN chown -R www-data:www-data /var/www/html

USER www-data

EXPOSE 80

CMD ["apache2-foreground"]
