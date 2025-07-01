# 📦 Деплой TPU Practice Searcher

> Приложение тестировалось на сервере с характеристиками:

- **ОС:** Ubuntu 22.04
- **CPU:** 1 x 3.3 ГГц
- **RAM:** 2 ГБ
- **Диск:** 30 ГБ NVMe

---

## ⚙️ Установка и запуск

### 1. Обновление системы

```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Установка необходимых пакетов

```bash
sudo apt install -y git
```

### 3. Клонирование проекта

```bash
git clone https://github.com/wonderfau1t/tpu-practice-searcher.git
cd tpu-practice-searcher
```

### 4. Установка Docker и Docker Compose

```bash
sh ./script.sh
```

---

## 🔐 Настройка переменных окружения

Создайте `.env` файл в корне проекта:

```bash
cp .env.example .env
```

Заполните его соответствующими значениями:

```env
# .env
POSTGRES_USER=notdefaultuser   # Задание пользователя для базы данных
POSTGRES_PASSWORD=notdefaultuserpassword   # Задание пароля
POSTGRES_DB=db   # Установка имени базы данных

TELEGRAM_BOT_TOKEN=your_telegram_bot_token   # Токен телеграм-бота
JWT_SECRET_TOKEN=your_jwt_secret_key   # Ключ для подписания JWT-токена
```

---

## 🐳 Сборка и запуск

### 5. Сборка фронтенда (один раз)

```bash
docker compose --profile build-only up frontend --build
```

### 6. Копирование собранного фронтенда в NGINX-директорию

```bash
sudo mkdir -p /var/www/tg-practice
sudo cp -r frontend-build/* /var/www/tg-practice/
```

### 7. Запуск всех сервисов

```bash
docker compose up -d --build
```

---

## 🌐 Настройка NGINX

### 8. Конфигурация NGINX

Создайте файл `/etc/nginx/sites-available/tpu-practice`:

```nginx
server {
    listen 443 ssl;
    server_name <ваш домен>;
    
    # Должен быть настроен ssl-сертификат
    # ssl_certificate ...;
    # ssl_certificate_key ...;


    location / {
        root /var/www/tg-practice;
        try_files $uri /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:8000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Активируйте сайт:

```bash
sudo ln -s /etc/nginx/sites-available/tpu-practice /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```
