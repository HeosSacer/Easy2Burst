CREATE DATABASE easy_db;
CREATE USER 'easy_db'@'localhost' IDENTIFIED BY 'easy';
GRANT ALL PRIVILEGES ON easy_db.* TO 'easy_db'@'localhost';
USE easy_db SOURCE C:/Users/Heos/AppData/Roaming/Easy2Burst/BurstWallet/init-mysql.sql;
ALTER DATABASE easy_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;