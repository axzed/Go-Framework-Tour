CREATE DATABASE IF NOT EXISTS `userapp`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `userapp`.user (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL,
    `avatar` varchar(512) NOT NULL ,
    `email` varchar(128) NOT NULL,
    `password` CHAR(128) NOT NULL,
    `create_time` INT UNSIGNED NOT NULL,
    `update_time` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE (email)
) CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `userapp`.user_extend (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED,
    `birthday` INT UNSIGNED NOT NULL,
    `nationality` SMALLINT NOT NULL ,
    `create_time` INT UNSIGNED NOT NULL,
    `update_time` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE (user_id)
) CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;