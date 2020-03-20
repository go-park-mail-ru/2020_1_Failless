# Как накатить базу локально  

1. Скачать postgresql-12  
Если у вас Ubuntu 16+, запустить скрипт [./install.sh](https://github.com/go-park-mail-ru/2020_1_Failless/blob/feature/users-search/scripts/install.sh)  
Если у вас не Ubuntu, то можете попробовать выполнить похожие команды  
2. Скачать расширение posgis  
Если вы запускали `install.sh`, оно уже у вас есть, если нет, то на сайте должна быть инструкция для вашей OS. Для macOS есть статья [тут](https://medium.com/@Umesh_Kafle/postgresql-and-postgis-installation-in-mac-os-87fa98a6814d)  
3. Скачать словари русского языка  
В случае Ubuntu:  
```sh
sudo apt install -y myspell-ru
cd /usr/share/postgresql/12/tsearch_data  
DICT=/usr/share/hunspell/ru_RU  
sudo iconv -f koi8-r -t utf-8 -o /usr/share/postgresql/12/tsearch_data/russian.affix $DICT.aff
sudo iconv -f koi8-r -t utf-8 -o /usr/share/postgresql/12/tsearch_data/russian.dict  $DICT.dic
```  
Иначе попробовать что-то похожее  
4. Зайти в postgresql:
```bash  
sudo su - postgres
psql  
\. /path/to/init.sql  -- перед этим прокомментировать верхние строчки  
```
По идее готово, но если не получится, то создайте пользователя и базу руками:  
```bash  
sudo su - postgres
psql  
CREATE USER $YOURUSERNAME WITH SUPERUSER PASSWORD '$YOURPASSWORD';
CREATE DATABASE eventum
\c eventum
CREATE EXTENSION IF NOT EXISTS citext;
\. ./scripts/init.sql  
