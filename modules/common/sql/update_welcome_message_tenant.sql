-- 更新 app_config 表中的租户欢迎语配置示例
-- 请根据实际需求修改配置内容后再执行

-- 方式1: 直接设置完整的JSON配置
UPDATE `app_config` 
SET `welcome_message_tenant` = '{
  "T18002704_zh_CN": "欢迎来到 cn2u.ai!",
  "T18002704_en_US": "Welcome to cn2u.ai!",
  "T18002705_zh_CN": "欢迎使用我们的服务!",
  "T18002705_en_US": "Welcome to our service!"
}'
WHERE `id` = 1;

-- 方式2: 使用JSON_SET增量更新(推荐,不会覆盖已有配置)
-- UPDATE `app_config` 
-- SET `welcome_message_tenant` = JSON_SET(
--   COALESCE(`welcome_message_tenant`, '{}'),
--   '$.T18002704_zh_CN', '欢迎来到 cn2u.ai!',
--   '$.T18002704_en_US', 'Welcome to cn2u.ai!',
--   '$.T18002705_zh_CN', '欢迎使用我们的服务!',
--   '$.T18002705_en_US', 'Welcome to our service!'
-- )
-- WHERE `id` = 1;

-- 查看更新后的配置
SELECT `id`, `welcome_message_tenant` FROM `app_config` WHERE `id` = 1;
