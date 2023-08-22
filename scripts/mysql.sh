docker run -d --name mysql-container \
    -e MYSQL_ROOT_PASSWORD=tiktokDB \
    -e MYSQL_USER=tiktokDB \
    -e MYSQL_PASSWORD=tiktokDB \
    -e MYSQL_DATABASE=db_tiktok \
    -p 3309:3306 \
    mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
