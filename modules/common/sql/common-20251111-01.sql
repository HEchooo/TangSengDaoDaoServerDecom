-- +migrate Up

-- 添加租户欢迎语配置字段
ALTER TABLE `app_config` ADD COLUMN `welcome_message_tenant` TEXT DEFAULT NULL COMMENT '租户欢迎语配置(JSON格式,key是tenantCode_language,如T18002704_zh_CN,value是welcomeMessage)';
