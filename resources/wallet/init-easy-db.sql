SELECT @@version;
CREATE DATABASE easy_db;
DROP USER 'easy_db'@'localhost';
flush privileges;
CREATE USER 'easy_db'@'localhost' IDENTIFIED BY 'easy';
GRANT ALL PRIVILEGES ON easy_db.* TO 'easy_db'@'localhost';
USE easy_db SOURCE C:/Users/Heos/AppData/Roaming/Easy2Burst/MariaDB/MariaDbMinimalist/bin/init-mysql.sql;