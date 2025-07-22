-- +migrate Up

ALTER TABLE `app_config` ADD COLUMN welcome_message_en TEXT NULL COMMENT '登录欢迎语 English';
