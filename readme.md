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

5. Настраиваем docker-compose.yaml
```yaml

```