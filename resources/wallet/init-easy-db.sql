SELECT @@version;
CREATE DATABASE easy_db;
DROP USER 'easy_db'@'localhost';
flush privileges;
CREATE USER 'easy_db'@'localhost' IDENTIFIED BY 'easy';
GRANT ALL PRIVILEGES ON easy_db.* TO 'easy_db'@'localhost';
USE easy_db SOURCE {PATH_TO_SQL_SCRIPT};
ALTER DATABASE easy_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;