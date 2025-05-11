-- 用户数据表
CREATE TABLE IF NOT EXISTS `tb_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码（加密存储）',
  `real_name` varchar(50) NOT NULL COMMENT '真实姓名',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `role` varchar(20) NOT NULL DEFAULT '普通用户' COMMENT '角色：管理员、编辑员、普通用户',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1正常，0禁用',
  `identifier` varchar(20) NOT NULL DEFAULT '未知' COMMENT '用户来源标识：小程序、未知',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `department` varchar(50) DEFAULT NULL COMMENT '部门',
  `position` varchar(50) DEFAULT NULL COMMENT '职位',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `sms_notifications` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否开启短信通知：1开启，0关闭',
  `theme` varchar(20) DEFAULT 'light' COMMENT '界面主题偏好',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 用户登录日志表
CREATE TABLE IF NOT EXISTS `tb_user_login_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` int(11) DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `login_ip` varchar(50) NOT NULL COMMENT '登录IP',
  `login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `login_type` varchar(20) NOT NULL DEFAULT 'web' COMMENT '登录类型：web、app、小程序',
  `login_status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '登录状态：1成功，0失败',
  `device_info` varchar(255) DEFAULT NULL COMMENT '设备信息',
  `user_agent` varchar(500) DEFAULT NULL COMMENT '浏览器UA信息',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注信息',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_login_time` (`login_time`),
  CONSTRAINT `fk_login_log_user` FOREIGN KEY (`user_id`) REFERENCES `tb_user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户登录日志表';

-- 用户操作日志表
CREATE TABLE IF NOT EXISTS `tb_user_operation_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` int(11) DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `operation_ip` varchar(50) NOT NULL COMMENT '操作IP',
  `operation_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  `module` varchar(50) NOT NULL COMMENT '操作模块',
  `operation` varchar(50) NOT NULL COMMENT '操作类型',
  `description` varchar(255) NOT NULL COMMENT '操作描述',
  `operation_result` tinyint(1) NOT NULL DEFAULT '1' COMMENT '操作结果：1成功，0失败',
  `details` text DEFAULT NULL COMMENT '操作详情',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_operation_time` (`operation_time`),
  CONSTRAINT `fk_operation_log_user` FOREIGN KEY (`user_id`) REFERENCES `tb_user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户操作日志表';

-- 用户扩展信息表（存储额外的用户信息）
CREATE TABLE IF NOT EXISTS `tb_user_profile` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `gender` tinyint(1) DEFAULT NULL COMMENT '性别：1男，2女，0未知',
  `birth_date` date DEFAULT NULL COMMENT '出生日期',
  `address` varchar(255) DEFAULT NULL COMMENT '地址',
  `id_card` varchar(20) DEFAULT NULL COMMENT '身份证号',
  `bio` text DEFAULT NULL COMMENT '个人简介',
  `tags` varchar(255) DEFAULT NULL COMMENT '标签，逗号分隔',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_user_profile` FOREIGN KEY (`user_id`) REFERENCES `tb_user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户扩展信息表';

-- 创建测试用户数据
INSERT INTO `tb_user` (`username`, `password`, `real_name`, `phone`, `role`, `status`, `identifier`, `department`, `position`, `created_at`) VALUES 
('testuser', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsNSGoTJNfarVq8q6/s.q', '测试用户', '13900000000', '普通用户', 1, '小程序', '测试部', '测试工程师', NOW()),
('editor', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsNSGoTJNfarVq8q6/s.q', '编辑人员', '13900000001', '编辑员', 1, '小程序', '内容部', '内容编辑', NOW()); 