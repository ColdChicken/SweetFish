CREATE DATABASE IF NOT EXISTS sweetfish DEFAULT CHARACTER SET utf8;
SET NAMES 'utf8';
USE sweetfish;

CREATE TABLE TOKEN(
  id INT NOT NULL AUTO_INCREMENT COMMENT '唯一id',
  token VARCHAR(128) NOT NULL COMMENT 'token',
  username VARCHAR(128) NOT NULL COMMENT 'token对应的username',
  expire_time DATETIME NOT NULL COMMENT 'token过期时间',
  PRIMARY KEY (id),
  KEY (token)
) ENGINE = INNODB;

-- 基于用户名密码的用户信息，兼容openid
CREATE TABLE USER_AUTH_INFO (
    id INT NOT NULL AUTO_INCREMENT COMMENT 'id',
    username VARCHAR(128) NOT NULL COMMENT '用户名',
    password VARCHAR(1024) NOT NULL COMMENT '密码',
    PRIMARY KEY (id),
    UNIQUE KEY(username)
) ENGINE = INNODB;

