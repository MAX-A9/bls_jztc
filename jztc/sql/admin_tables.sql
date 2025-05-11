-- 管理员数据表
CREATE TABLE IF NOT EXISTS `tb_admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码（加密存储）',
  `real_name` varchar(50) NOT NULL COMMENT '真实姓名',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `role_group` varchar(50) NOT NULL COMMENT '角色组',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1正常，0禁用',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

-- 管理员权限表（用于存储管理员的权限，与管理员是多对多关系）
CREATE TABLE IF NOT EXISTS `tb_admin_permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `admin_id` int(11) NOT NULL COMMENT '管理员ID',
  `permission_name` varchar(50) NOT NULL COMMENT '权限名称',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  CONSTRAINT `fk_admin_permission_admin` FOREIGN KEY (`admin_id`) REFERENCES `tb_admin` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员权限表';

-- 权限基础数据
INSERT INTO `tb_admin_permission` (`admin_id`, `permission_name`) VALUES 
(1, '用户管理'),
(1, '内容管理'),
(1, '系统设置'),
(1, '数据统计'),
(1, '权限管理');

-- 管理员操作日志表
CREATE TABLE IF NOT EXISTS `tb_admin_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `admin_id` int(11) DEFAULT NULL COMMENT '管理员ID',
  `admin_name` varchar(50) NOT NULL COMMENT '管理员用户名',
  `ip` varchar(50) NOT NULL COMMENT '操作IP',
  `action` varchar(50) NOT NULL COMMENT '操作类型：登录、登出、添加、修改、删除、导出、导入、查询',
  `module` varchar(50) NOT NULL COMMENT '操作模块',
  `description` varchar(255) NOT NULL COMMENT '操作描述',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '操作状态：1成功，0失败',
  `details` text DEFAULT NULL COMMENT '详细信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_admin_log_admin` FOREIGN KEY (`admin_id`) REFERENCES `tb_admin` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员操作日志表';

-- 创建超级管理员账号
INSERT INTO `tb_admin` (`username`, `password`, `real_name`, `phone`, `role_group`, `status`, `created_at`) VALUES 
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsNSGoTJNfarVq8q6/s.q', '系统管理员', '13800000000', '超级管理员', 1, NOW()); 