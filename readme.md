> Приложение тестировалось на сервере с характеристиками:\
> OC – Ubuntu 22.04\
> CPU – 1 x 3.3 ГГц\
> RAM – 2 ГБ\
> Диск – 30 ГБ NVMe

1. Обновляем систему

```shell
apt update && apt upgrade -y
```

2. Устанавливаем git

```shell
apt install -y git
```

3. Клонируем проект

```shell
git clone https://github.com/wonderfau1t/tpu-practice-searcher.git
cd tpu-practice-searcher
```

4. Запускаем скрипт установки Docker и Docker Compose

```shell
sh ./script.sh
```

5. Настраиваем ``docker-compose.yaml``
> Комментарии по настройке указаны внутри файла
6. Билдим проект
```shell
docker compose up -d --build
```
7. Создаем папку, где будет храниться наше SPA-приложение (фронтенд)
```shell
mkdir /var/www/tg-practice
```
8. Копируем билд фронтенда в папку
```shell
cp frontend-build/* /var/www/tg-practice/ -r
```
9. Настраиваем nginx для SPA-приложения
```text
server {
    listen 443 ssl;
    server_name <ваш домен>;
    # тут настроенный ssl-сертификат
    
    location / {
        root /var/www/tg-practice;
        try_files $uri /index.html;
    }
    
    location /api/ {
        proxy_pass http://localhost:8001/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```