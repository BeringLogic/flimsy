FROM nginx:alpine

COPY ./web/default.conf /etc/nginx/conf.d/default.conf

COPY ./app/src/* /var/www/

RUN mkdir -p /var/www/vendor/League/ColorExtractor
COPY ./app/color-extractor/src/* /var/www/vendor/League/ColorExtractor

