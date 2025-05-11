-- 常规管理相关表

-- 系统配置表
CREATE TABLE IF NOT EXISTS `tb_system_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `config_key` varchar(50) NOT NULL COMMENT '配置键名',
  `config_value` text COMMENT '配置值',
  `config_group` varchar(50) NOT NULL DEFAULT 'system' COMMENT '配置分组',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注说明',
  `is_system` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否系统配置：1是，0否',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_config_key` (`config_key`),
  KEY `idx_config_group` (`config_group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 用户偏好设置表
CREATE TABLE IF NOT EXISTS `tb_user_preference` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设置ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `preference_key` varchar(50) NOT NULL COMMENT '偏好设置键名',
  `preference_value` text COMMENT '偏好设置值',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_key` (`user_id`,`preference_key`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户偏好设置表';

-- 短信渠道表
CREATE TABLE IF NOT EXISTS `tb_sms_channel` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '渠道ID',
  `name` varchar(50) NOT NULL COMMENT '渠道名称',
  `app_id` varchar(100) NOT NULL COMMENT '应用ID',
  `app_key` varchar(255) NOT NULL COMMENT '应用密钥',
  `sign_name` varchar(50) DEFAULT NULL COMMENT '短信签名',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1启用，0禁用',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认：1是，0否',
  `config` text COMMENT '额外配置JSON',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信渠道表';

-- 短信模板表
CREATE TABLE IF NOT EXISTS `tb_sms_template` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `name` varchar(50) NOT NULL COMMENT '模板名称',
  `code` varchar(100) NOT NULL COMMENT '模板编码',
  `content` text NOT NULL COMMENT '模板内容',
  `channel_id` int(11) NOT NULL COMMENT '所属渠道ID',
  `type` varchar(20) NOT NULL COMMENT '模板类型：验证码、通知、营销',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1启用，0禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_channel_id` (`channel_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信模板表';

-- 短信发送记录表
CREATE TABLE IF NOT EXISTS `tb_sms_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `mobile` varchar(20) NOT NULL COMMENT '手机号',
  `content` text NOT NULL COMMENT '短信内容',
  `template_id` int(11) DEFAULT NULL COMMENT '模板ID',
  `channel_id` int(11) NOT NULL COMMENT '渠道ID',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '发送状态：1成功，0失败',
  `error_msg` varchar(255) DEFAULT NULL COMMENT '错误消息',
  `send_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_mobile` (`mobile`),
  KEY `idx_status` (`status`),
  KEY `idx_send_time` (`send_time`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_channel_id` (`channel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信发送记录表';

-- 短信验证码表
CREATE TABLE IF NOT EXISTS `tb_sms_code` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `mobile` varchar(20) NOT NULL COMMENT '手机号',
  `code` varchar(10) NOT NULL COMMENT '验证码',
  `type` varchar(20) NOT NULL COMMENT '验证码类型：登录、注册、找回密码',
  `expire_time` datetime NOT NULL COMMENT '过期时间',
  `is_used` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已使用：1是，0否',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_mobile_type` (`mobile`,`type`),
  KEY `idx_expire_time` (`expire_time`),
  KEY `idx_is_used` (`is_used`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信验证码表';

-- 个人设置表
CREATE TABLE IF NOT EXISTS `tb_user_profile_setting` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设置ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `position` varchar(50) DEFAULT NULL COMMENT '职位',
  `department` varchar(50) DEFAULT NULL COMMENT '部门',
  `theme` varchar(20) DEFAULT 'light' COMMENT '界面主题：light(浅色)、dark(深色)',
  `sms_notifications` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否开启短信通知：1开启，0关闭',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='个人设置表';

-- 添加系统配置测试数据
INSERT INTO `tb_system_config` (`config_key`, `config_value`, `config_group`, `remark`, `is_system`) VALUES
('site_name', '后台管理系统', 'system', '站点名称', 1),
('site_logo', '/logo.png', 'system', '站点LOGO', 1),
('admin_email', 'admin@example.com', 'system', '管理员邮箱', 1),
('file_upload_limit', '10', 'file', '文件上传大小限制（MB）', 1),
('cache_duration', '3600', 'system', '缓存时长（秒）', 1),
('default_language', 'zh_CN', 'locale', '默认语言', 1),
('date_format', 'YYYY-MM-DD', 'locale', '日期格式', 1),
('time_format', 'HH:mm:ss', 'locale', '时间格式', 1),
('oss_enabled', '0', 'storage', '是否启用OSS存储', 1),
('oss_endpoint', 'oss-cn-beijing.aliyuncs.com', 'storage', 'OSS域名', 1),
('oss_bucket', 'my-bucket', 'storage', 'OSS存储桶', 1),
('oss_access_key_id', 'LTAI4xxxxxxxxxxxxx', 'storage', 'OSS访问密钥ID', 1),
('oss_access_key_secret', '********', 'storage', 'OSS访问密钥', 1),
('oss_region', 'oss-cn-beijing', 'storage', 'OSS区域', 1),
('oss_directory', 'uploads/', 'storage', 'OSS存储目录', 1),
('oss_public_access', '1', 'storage', '是否允许公开访问', 1);

-- 添加短信渠道测试数据
INSERT INTO `tb_sms_channel` (`name`, `app_id`, `app_key`, `sign_name`, `status`, `is_default`) VALUES
('阿里云短信', 'LTAI5txxxxxxxxxxxxxxx', '********', '测试签名', 1, 1),
('腾讯云短信', 'SDK_APPIDxxxxxxxx', '********', '测试签名2', 0, 0);

-- 添加短信模板测试数据
INSERT INTO `tb_sms_template` (`name`, `code`, `content`, `channel_id`, `type`, `status`) VALUES
('注册验证码', 'SMS_123456789', '您的验证码为：${code}，有效期10分钟，请勿泄露给他人。', 1, '验证码', 1),
('找回密码验证码', 'SMS_987654321', '您正在找回密码，验证码为：${code}，有效期10分钟，请勿泄露给他人。', 1, '验证码', 1),
('订单支付成功通知', 'SMS_567891234', '您的订单${orderNo}已支付成功，感谢您的购买！', 2, '通知', 1);

-- 添加短信发送记录测试数据
INSERT INTO `tb_sms_record` (`mobile`, `content`, `template_id`, `channel_id`, `status`, `error_msg`, `send_time`) VALUES
('13800138000', '您的验证码为：123456，有效期10分钟，请勿泄露给他人。', 1, 1, 1, NULL, DATE_SUB(NOW(), INTERVAL 2 DAY)),
('13900139000', '您正在找回密码，验证码为：654321，有效期10分钟，请勿泄露给他人。', 2, 1, 1, NULL, DATE_SUB(NOW(), INTERVAL 1 DAY)),
('13700137000', '您的订单20230918001已支付成功，感谢您的购买！', 3, 2, 0, '短信余额不足', DATE_SUB(NOW(), INTERVAL 12 HOUR));

-- 添加短信验证码测试数据
INSERT INTO `tb_sms_code` (`mobile`, `code`, `type`, `expire_time`, `is_used`) VALUES
('13800138000', '123456', '登录', DATE_ADD(NOW(), INTERVAL 10 MINUTE), 0),
('13900139000', '654321', '找回密码', DATE_ADD(NOW(), INTERVAL 10 MINUTE), 1);

-- 添加个人设置测试数据
INSERT INTO `tb_user_profile_setting` (`user_id`, `avatar`, `position`, `department`, `theme`, `sms_notifications`) VALUES
(1, 'https://randomuser.me/api/portraits/men/42.jpg', '系统管理员', '技术部', 'light', 0); 