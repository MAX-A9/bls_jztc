-- 文件管理相关表

-- 文件表
CREATE TABLE IF NOT EXISTS `tb_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `file_name` varchar(255) NOT NULL COMMENT '文件名',
  `original_name` varchar(255) NOT NULL COMMENT '原始文件名',
  `file_type` varchar(50) NOT NULL COMMENT '文件类型',
  `file_size` bigint(20) NOT NULL COMMENT '文件大小(字节)',
  `file_path` varchar(255) NOT NULL COMMENT '文件路径',
  `thumbnail_path` varchar(255) DEFAULT NULL COMMENT '缩略图路径(图片类型)',
  `category` enum('image','document','media','other') NOT NULL DEFAULT 'other' COMMENT '文件分类',
  `description` text DEFAULT NULL COMMENT '文件描述',
  `user_id` int(11) NOT NULL COMMENT '上传用户ID',
  `uploader` varchar(50) NOT NULL COMMENT '上传者名称',
  `status` enum('normal','processing','deleted') NOT NULL DEFAULT 'normal' COMMENT '文件状态：normal(正常)、processing(处理中)、deleted(已删除)',
  `download_count` int(11) NOT NULL DEFAULT '0' COMMENT '下载次数',
  `view_count` int(11) NOT NULL DEFAULT '0' COMMENT '浏览次数',
  `is_public` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否公开：1是，0否',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_file_type` (`file_type`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表';

-- 文件分享表
CREATE TABLE IF NOT EXISTS `tb_file_share` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分享ID',
  `file_id` int(11) NOT NULL COMMENT '文件ID',
  `user_id` int(11) NOT NULL COMMENT '分享用户ID',
  `share_code` varchar(32) NOT NULL COMMENT '分享码',
  `expire_time` datetime DEFAULT NULL COMMENT '过期时间',
  `view_count` int(11) NOT NULL DEFAULT '0' COMMENT '查看次数',
  `download_count` int(11) NOT NULL DEFAULT '0' COMMENT '下载次数',
  `is_password` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否设置密码：1是，0否',
  `password` varchar(32) DEFAULT NULL COMMENT '访问密码',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1有效，0无效',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_share_code` (`share_code`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件分享表';

-- 文件标签表
CREATE TABLE IF NOT EXISTS `tb_file_tag` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '标签ID',
  `tag_name` varchar(50) NOT NULL COMMENT '标签名称',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件标签表';

-- 文件标签关联表
CREATE TABLE IF NOT EXISTS `tb_file_tag_relation` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `file_id` int(11) NOT NULL COMMENT '文件ID',
  `tag_id` int(11) NOT NULL COMMENT '标签ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_tag` (`file_id`,`tag_id`),
  KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件标签关联表';

-- 文件操作日志表
CREATE TABLE IF NOT EXISTS `tb_file_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `file_id` int(11) NOT NULL COMMENT '文件ID',
  `user_id` int(11) NOT NULL COMMENT '操作用户ID',
  `username` varchar(50) NOT NULL COMMENT '操作用户名',
  `action` enum('upload','download','view','delete','rename','share','update') NOT NULL COMMENT '操作类型',
  `description` varchar(255) DEFAULT NULL COMMENT '操作描述',
  `ip` varchar(50) DEFAULT NULL COMMENT '操作IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件操作日志表';

-- 文件夹表
CREATE TABLE IF NOT EXISTS `tb_folder` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '文件夹ID',
  `folder_name` varchar(100) NOT NULL COMMENT '文件夹名称',
  `parent_id` int(11) DEFAULT NULL COMMENT '父文件夹ID',
  `user_id` int(11) NOT NULL COMMENT '创建用户ID',
  `description` varchar(255) DEFAULT NULL COMMENT '文件夹描述',
  `is_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否公开：1是，0否',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件夹表';

-- 文件与文件夹关联表
CREATE TABLE IF NOT EXISTS `tb_file_folder` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `file_id` int(11) NOT NULL COMMENT '文件ID',
  `folder_id` int(11) NOT NULL COMMENT '文件夹ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_folder` (`file_id`,`folder_id`),
  KEY `idx_folder_id` (`folder_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件与文件夹关联表';

-- 添加测试数据 (文件夹)
INSERT INTO `tb_folder` (`id`, `folder_name`, `parent_id`, `user_id`, `description`, `is_public`) VALUES
(1, '文档', NULL, 1, '重要文档文件夹', 1),
(2, '图片', NULL, 1, '图片文件夹', 1),
(3, '合同', 1, 1, '合同文档', 0),
(4, '产品资料', NULL, 2, '产品相关文档', 1),
(5, '系统文件', NULL, 1, '系统相关文件', 0);

-- 添加测试数据 (文件)
INSERT INTO `tb_file` (`id`, `file_name`, `original_name`, `file_type`, `file_size`, `file_path`, `thumbnail_path`, `category`, `description`, `user_id`, `uploader`, `status`, `download_count`, `view_count`, `is_public`) VALUES
(1, 'product_manual_v1.2.pdf', '产品手册V1.2.pdf', '文档', 1258231, '/files/documents/product_manual_v1.2.pdf', NULL, 'document', '产品使用手册最新版', 1, '管理员', 'normal', 23, 145, 1),
(2, 'company_logo.png', '公司Logo.png', '图片', 256879, '/files/images/company_logo.png', '/files/thumbnails/company_logo_thumb.png', 'image', '公司标准Logo', 1, '管理员', 'normal', 12, 78, 1),
(3, 'quarterly_report_2023Q1.xlsx', '2023年第一季度报告.xlsx', '表格', 785412, '/files/documents/quarterly_report_2023Q1.xlsx', NULL, 'document', '2023年第一季度财务报告', 2, '财务主管', 'normal', 5, 19, 0),
(4, 'product_video.mp4', '产品宣传视频.mp4', '视频', 25698745, '/files/media/product_video.mp4', '/files/thumbnails/product_video_thumb.jpg', 'media', '产品宣传片', 3, '市场经理', 'normal', 45, 230, 1),
(5, 'contract_template.docx', '合同模板.docx', '文档', 358741, '/files/documents/contract_template.docx', NULL, 'document', '标准合同模板', 1, '管理员', 'normal', 18, 95, 0),
(6, 'meeting_notes.txt', '会议纪要.txt', '文本', 15879, '/files/documents/meeting_notes.txt', NULL, 'document', '2023年5月15日项目进度会议纪要', 2, '项目经理', 'normal', 7, 25, 0),
(7, 'system_architecture.png', '系统架构图.png', '图片', 458796, '/files/images/system_architecture.png', '/files/thumbnails/system_architecture_thumb.png', 'image', '系统架构设计图', 1, '技术总监', 'normal', 15, 68, 1),
(8, 'user_manual.pdf', '用户手册.pdf', '文档', 958741, '/files/documents/user_manual.pdf', NULL, 'document', '用户使用手册', 3, '产品经理', 'normal', 28, 132, 1);

-- 添加测试数据 (文件标签)
INSERT INTO `tb_file_tag` (`id`, `tag_name`) VALUES
(1, '重要'),
(2, '产品'),
(3, '合同'),
(4, '报告'),
(5, '宣传'),
(6, '文档'),
(7, '设计'),
(8, '公开');

-- 添加测试数据 (文件标签关联)
INSERT INTO `tb_file_tag_relation` (`file_id`, `tag_id`) VALUES
(1, 2), (1, 6), (1, 8),
(2, 5), (2, 7), (2, 8),
(3, 1), (3, 4), 
(4, 2), (4, 5), (4, 8),
(5, 1), (5, 3), (5, 6),
(6, 4), (6, 6),
(7, 1), (7, 7),
(8, 2), (8, 6), (8, 8);

-- 添加测试数据 (文件与文件夹关联)
INSERT INTO `tb_file_folder` (`file_id`, `folder_id`) VALUES
(1, 4), 
(2, 2),
(3, 1),
(4, 2),
(5, 3),
(6, 1),
(7, 5),
(8, 4);

-- 添加测试数据 (文件分享)
INSERT INTO `tb_file_share` (`file_id`, `user_id`, `share_code`, `expire_time`, `view_count`, `download_count`, `is_password`, `password`, `status`) VALUES
(1, 1, 'a1b2c3d4e5f6g7h8', DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 7 DAY), 15, 8, 0, NULL, 1),
(2, 1, 'z1x2c3v4b5n6m7a8', NULL, 25, 10, 1, '123456', 1),
(4, 3, 'q1w2e3r4t5y6u7i8', DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 30 DAY), 42, 18, 0, NULL, 1),
(8, 3, 'p1o2i3u4y5t6r7e8', DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 14 DAY), 12, 5, 1, 'abc123', 1);

-- 添加测试数据 (文件日志)
INSERT INTO `tb_file_log` (`file_id`, `user_id`, `username`, `action`, `description`, `ip`) VALUES
(1, 1, '管理员', 'upload', '上传了文件', '192.168.1.100'),
(1, 2, '用户A', 'download', '下载了文件', '192.168.1.101'),
(1, 3, '用户B', 'view', '查看了文件', '192.168.1.102'),
(2, 1, '管理员', 'upload', '上传了文件', '192.168.1.100'),
(2, 1, '管理员', 'share', '分享了文件', '192.168.1.100'),
(3, 2, '财务主管', 'upload', '上传了文件', '192.168.1.103'),
(3, 1, '管理员', 'view', '查看了文件', '192.168.1.100'),
(4, 3, '市场经理', 'upload', '上传了文件', '192.168.1.104'),
(4, 3, '市场经理', 'share', '分享了文件', '192.168.1.104'),
(5, 1, '管理员', 'upload', '上传了文件', '192.168.1.100'),
(6, 2, '项目经理', 'upload', '上传了文件', '192.168.1.105'),
(7, 1, '技术总监', 'upload', '上传了文件', '192.168.1.106'),
(8, 3, '产品经理', 'upload', '上传了文件', '192.168.1.107'),
(8, 3, '产品经理', 'share', '分享了文件', '192.168.1.107'),
(8, 1, '管理员', 'download', '下载了文件', '192.168.1.100'); 