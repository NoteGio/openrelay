FROM nginx

ADD docker-cfg/proxy_pass.template /etc/nginx/conf.d/proxy_pass.template

CMD /bin/bash -c "envsubst < /etc/nginx/conf.d/proxy_pass.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"
