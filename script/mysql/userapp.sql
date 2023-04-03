CREATE DATABASE IF NOT EXISTS `userapp`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `userapp`.user (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL,
    `avatar` varchar(512) NOT NULL ,
#     email 128 已经够长了，不行你还可以用 256
    `email` varchar(128) NOT NULL,
#     加密后固定是 128 字节
    `password` CHAR(128) NOT NULL,
#     UUID 的字符串形式是固定 36 字节
    `salt` CHAR(36) NOT NULL,
    `create_time` INT UNSIGNED NOT NULL,
    `update_time` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE (email)
) CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

# CREATE TABLE IF NOT EXISTS `userapp`.user_extend (
#     `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
#     `user_id` BIGINT UNSIGNED,
#     `birthday` INT UNSIGNED NOT NULL,
#     `nationality` SMALLINT NOT NULL ,
#     `create_time` INT UNSIGNED NOT NULL,
#     `update_time` INT UNSIGNED NOT NULL,
#     PRIMARY KEY (`id`),
#     UNIQUE (user_id)
# ) CHARACTER SET utf8mb4
# COLLATE utf8mb4_unicode_ci;