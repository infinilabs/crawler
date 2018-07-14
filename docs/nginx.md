## Working with Nginx

```
server {
listen 80;
listen 443;
server_name example.org;

      access_log /home/data/example.org.log;


location / {
    proxy_http_version 1.1;
    proxy_pass http://127.0.0.1:9001;
}

location /ws {
    proxy_buffering off;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Real-IP $remote_addr;

    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "Upgrade";
    proxy_pass http://127.0.0.1:9001;
}

#Some protect here
#location /admin/ {
#    proxy_http_version 1.1;
#    proxy_pass http://127.0.0.1:9001/;
#}

}
```
