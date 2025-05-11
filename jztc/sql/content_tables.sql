-- 内容数据表
CREATE TABLE IF NOT EXISTS `tb_content` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '内容ID',
  `title` varchar(255) NOT NULL COMMENT '标题',
  `category` varchar(50) NOT NULL COMMENT '分类',
  `author` varchar(50) NOT NULL COMMENT '作者',
  `content` text NOT NULL COMMENT '内容详情(富文本)',
  `status` varchar(20) NOT NULL DEFAULT '待审核' COMMENT '状态：已发布、待审核、已下架',
  `views` int(11) NOT NULL DEFAULT '0' COMMENT '浏览量',
  `likes` int(11) NOT NULL DEFAULT '0' COMMENT '想要数量',
  `comments` int(11) NOT NULL DEFAULT '0' COMMENT '评论数',
  `is_recommended` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否置顶推荐：1是，0否',
  `published_at` datetime DEFAULT NULL COMMENT '发布时间',
  `expires_at` datetime DEFAULT NULL COMMENT '到期时间',
  `top_until` datetime DEFAULT NULL COMMENT '置顶截止时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_is_recommended` (`is_recommended`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容表';

-- 首页分类表
CREATE TABLE IF NOT EXISTS `tb_home_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序顺序',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1启用，0禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页分类表';

-- 闲置分类表
CREATE TABLE IF NOT EXISTS `tb_idle_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序顺序',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1启用，0禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='闲置分类表';

-- 内容评论表
CREATE TABLE IF NOT EXISTS `tb_content_comment` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `content_id` int(11) NOT NULL COMMENT '内容ID',
  `user_id` int(11) DEFAULT NULL COMMENT '用户ID',
  `user_name` varchar(50) NOT NULL COMMENT '用户名',
  `comment` text NOT NULL COMMENT '评论内容',
  `status` varchar(20) NOT NULL DEFAULT '待审核' COMMENT '状态：已审核、待审核、已拒绝',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_comment_content` FOREIGN KEY (`content_id`) REFERENCES `tb_content` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容评论表';

-- 内容附件表
CREATE TABLE IF NOT EXISTS `tb_content_attachment` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '附件ID',
  `content_id` int(11) NOT NULL COMMENT '内容ID',
  `file_name` varchar(255) NOT NULL COMMENT '文件名',
  `file_path` varchar(255) NOT NULL COMMENT '文件路径',
  `file_size` int(11) NOT NULL COMMENT '文件大小(字节)',
  `file_type` varchar(50) NOT NULL COMMENT '文件类型',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_content_id` (`content_id`),
  CONSTRAINT `fk_attachment_content` FOREIGN KEY (`content_id`) REFERENCES `tb_content` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容附件表';

-- 内容操作日志表
CREATE TABLE IF NOT EXISTS `tb_content_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `content_id` int(11) NOT NULL COMMENT '内容ID',
  `admin_id` int(11) DEFAULT NULL COMMENT '管理员ID',
  `admin_name` varchar(50) NOT NULL COMMENT '管理员名称',
  `action` varchar(50) NOT NULL COMMENT '操作类型：创建、编辑、删除、审核、推荐',
  `status_change` varchar(50) DEFAULT NULL COMMENT '状态变更：变更前状态->变更后状态',
  `details` text DEFAULT NULL COMMENT '操作详情',
  `ip` varchar(50) NOT NULL COMMENT '操作IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_action` (`action`),
  CONSTRAINT `fk_log_content` FOREIGN KEY (`content_id`) REFERENCES `tb_content` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容操作日志表';

-- 创建示例数据：首页分类
INSERT INTO `tb_home_category` (`name`, `sort_order`, `is_active`) VALUES 
('手机数码', 1, 1),
('家用电器', 2, 1),
('电脑办公', 3, 1),
('服装鞋帽', 4, 0),
('图书音像', 5, 1);

-- 创建示例数据：闲置分类
INSERT INTO `tb_idle_category` (`name`, `sort_order`, `is_active`) VALUES 
('二手手机', 1, 1),
('二手电脑', 2, 1),
('闲置家电', 3, 1),
('二手书籍', 4, 0),
('闲置服装', 5, 1);

-- 创建示例数据：内容
INSERT INTO `tb_content` (`title`, `category`, `author`, `content`, `status`, `views`, `likes`, `comments`, `is_recommended`, `published_at`, `created_at`) VALUES 
('示例内容标题1', '首页', '测试作者1', '<p>这是示例内容1的详细描述，支持<strong>富文本</strong>编辑。</p><p>可以插入图片、表格等多媒体内容。</p>', '已发布', 1200, 45, 5, 1, NOW(), NOW()),
('示例内容标题2', '闲置', '测试作者2', '<p>这是示例内容2的详细描述，支持<strong>富文本</strong>编辑。</p><p>可以插入图片、表格等多媒体内容。</p>', '待审核', 0, 0, 0, 0, NULL, NOW()); 