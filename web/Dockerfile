FROM nginx:alpine

COPY ./web/nginx.conf /etc/nginx/nginx.conf
COPY ./web/default.conf /etc/nginx/conf.d/default.conf

COPY ./app/src/ /var/www/
COPY ./web/favicon.ico /var/www/
