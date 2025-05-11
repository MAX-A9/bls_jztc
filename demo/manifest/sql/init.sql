-- 创建用户表
CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `nickname` varchar(50) NOT NULL COMMENT '昵称',
  `role` varchar(20) NOT NULL DEFAULT 'admin' COMMENT '角色 admin:管理员 user:普通用户',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0:禁用 1:正常',
  `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 插入管理员账号
INSERT INTO `user` (`username`, `password`, `nickname`, `role`, `status`)
VALUES ('admin', MD5('Admin123'), '系统管理员', 'admin', 1)
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;

-- 创建客户表
CREATE TABLE IF NOT EXISTS `client` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '客户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `real_name` varchar(50) NOT NULL COMMENT '真实姓名',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0:禁用 1:正常',
  `identifier` varchar(20) NOT NULL DEFAULT 'unknown' COMMENT '来源标识 wxapp:小程序 unknown:未知',
  `open_id` varchar(50) DEFAULT NULL COMMENT '微信openid',
  `session_key` varchar(50) DEFAULT NULL COMMENT '微信session_key',
  `avatar_url` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户表';

-- 创建文件表
CREATE TABLE IF NOT EXISTS `file` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `name` varchar(255) NOT NULL COMMENT '文件名',
  `path` varchar(255) NOT NULL COMMENT '文件路径',
  `size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小(字节)',
  `type` varchar(50) NOT NULL COMMENT '文件类型(image/document/video等)',
  `content_type` varchar(100) DEFAULT NULL COMMENT '内容类型',
  `extension` varchar(20) DEFAULT NULL COMMENT '扩展名',
  `is_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否公开 0:私有 1:公开',
  `user_id` int(11) DEFAULT NULL COMMENT '上传用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '上传用户名',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_is_public` (`is_public`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表';

-- 创建订单表
CREATE TABLE IF NOT EXISTS `order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(32) NOT NULL COMMENT '订单编号',
  `client_name` varchar(50) NOT NULL COMMENT '客户名称',
  `content_id` int(11) DEFAULT NULL COMMENT '内容ID',
  `client_id` int(11) NOT NULL COMMENT '客户ID',
  `product_name` varchar(100) NOT NULL COMMENT '商品名称',
  `amount` decimal(10,2) NOT NULL COMMENT '订单金额',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0:待支付 1:已支付 2:已取消',
  `payment_method` varchar(20) DEFAULT NULL COMMENT '支付方式 wechat:微信 alipay:支付宝',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `expire_time` datetime DEFAULT NULL COMMENT '订单过期时间',
  `transaction_id` varchar(64) DEFAULT NULL COMMENT '支付交易号',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_client_id` (`client_id`),
  KEY `idx_status` (`status`),
  KEY `idx_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';


-- 创建内容表
CREATE TABLE IF NOT EXISTS `content` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '内容ID',
  `title` varchar(255) NOT NULL COMMENT '标题',
  `category` varchar(50) NOT NULL COMMENT '分类',
  `author` varchar(50) NOT NULL COMMENT '作者',
  `region_id` int(11) DEFAULT NULL COMMENT '所属地区ID',
  `client_id` int(11) DEFAULT NULL COMMENT '客户ID',
  `content` text NOT NULL COMMENT '内容详情(富文本)',
  `extend` text DEFAULT NULL COMMENT '扩展字段(JSON格式)',
  `status` varchar(20) NOT NULL DEFAULT '待审核' COMMENT '状态：已发布、待审核、已下架',
  `views` int(11) NOT NULL DEFAULT '0' COMMENT '浏览量',
  `likes` int(11) NOT NULL DEFAULT '0' COMMENT '想要数量',
  `comments` int(11) NOT NULL DEFAULT '0' COMMENT '评论数',
  `is_recommended` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否置顶推荐：1是，0否',
  `published_at` datetime DEFAULT NULL COMMENT '发布时间',
  `expires_at` datetime DEFAULT NULL COMMENT '到期时间',
  `top_until` datetime DEFAULT NULL COMMENT '置顶截止时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_is_recommended` (`is_recommended`),
  KEY `idx_region_id` (`region_id`),
  KEY `idx_client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容表';

-- 创建首页分类表
CREATE TABLE IF NOT EXISTS `home_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1是，0否',
  `icon` varchar(255) DEFAULT NULL COMMENT '分类图标URL',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页分类表';

-- 创建闲置分类表
CREATE TABLE IF NOT EXISTS `idle_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1是，0否',
  `icon` varchar(255) DEFAULT NULL COMMENT '分类图标URL',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='闲置分类表';

-- 插入首页分类测试数据
INSERT INTO `home_category` (`name`, `sort_order`, `is_active`, `icon`)
VALUES 
('手机数码', 1, 1, 'https://img.icons8.com/color/48/000000/smartphone.png'),
('家用电器', 2, 1, 'https://img.icons8.com/color/48/000000/washing-machine.png'),
('电脑办公', 3, 1, 'https://img.icons8.com/color/48/000000/laptop.png'),
('服装鞋帽', 4, 0, 'https://img.icons8.com/color/48/000000/t-shirt.png'),
('图书音像', 5, 1, 'https://img.icons8.com/color/48/000000/books.png')
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;

-- 插入闲置分类测试数据
INSERT INTO `idle_category` (`name`, `sort_order`, `is_active`, `icon`)
VALUES 
('二手手机', 1, 1, 'https://img.icons8.com/color/48/000000/old-phone.png'),
('二手电脑', 2, 1, 'https://img.icons8.com/color/48/000000/laptop-e-mail.png'),
('闲置家电', 3, 1, 'https://img.icons8.com/color/48/000000/fridge.png'),
('二手书籍', 4, 0, 'https://img.icons8.com/color/48/000000/book-shelf.png'),
('闲置服装', 5, 1, 'https://img.icons8.com/color/48/000000/clothes.png')
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;

-- 创建套餐表
CREATE TABLE IF NOT EXISTS `package` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '套餐ID',
  `title` varchar(50) NOT NULL COMMENT '套餐名称',
  `description` varchar(500) NOT NULL COMMENT '套餐简介',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格(元)',
  `type` varchar(20) NOT NULL COMMENT '套餐类型: top-置顶套餐, publish-发布套餐',
  `duration` int(11) NOT NULL DEFAULT '1' COMMENT '时长值',
  `duration_type` varchar(10) NOT NULL DEFAULT 'day' COMMENT '时长单位: hour-小时, day-天, month-月',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值，数字越小排序越靠前',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='套餐表';

-- 插入默认套餐数据 - 置顶套餐
INSERT INTO `package` (`title`, `description`, `price`, `type`, `duration`, `duration_type`, `sort_order`, `created_at`)
VALUES 
('置顶一小时', '内容将在列表顶部显示一小时', 2.00, 'top', 1, 'hour', 1, CURRENT_TIMESTAMP),
('置顶一天', '内容将在列表顶部显示24小时', 10.00, 'top', 1, 'day', 2, CURRENT_TIMESTAMP),
('置顶三天', '内容将在列表顶部显示3天', 25.00, 'top', 3, 'day', 3, CURRENT_TIMESTAMP),
('置顶一周', '内容将在列表顶部显示7天', 50.00, 'top', 7, 'day', 4, CURRENT_TIMESTAMP)
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;

-- 插入默认套餐数据 - 发布套餐
INSERT INTO `package` (`title`, `description`, `price`, `type`, `duration`, `duration_type`, `sort_order`, `created_at`)
VALUES 
('7天展示', '发布内容可展示7天', 3.00, 'publish', 7, 'day', 1, CURRENT_TIMESTAMP),
('15天展示', '发布内容可展示15天', 8.00, 'publish', 15, 'day', 2, CURRENT_TIMESTAMP),
('30天展示', '发布内容可展示30天', 15.00, 'publish', 30, 'day', 3, CURRENT_TIMESTAMP),
('60天展示', '发布内容可展示60天', 25.00, 'publish', 60, 'day', 4, CURRENT_TIMESTAMP)
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;

-- 创建地区表
CREATE TABLE IF NOT EXISTS `region` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '地区ID',
  `location` varchar(255) NOT NULL COMMENT '所在地区，如：北京市/市辖区/东城区',
  `name` varchar(50) NOT NULL COMMENT '地区名称',
  `level` varchar(10) NOT NULL COMMENT '级别: 省,县,乡',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0:启用 1:禁用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_location` (`location`(191)),
  KEY `idx_level` (`level`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='地区表';

-- 创建评论表
CREATE TABLE IF NOT EXISTS `content_comment` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `content_id` int(11) NOT NULL COMMENT '内容ID',
  `client_id` int(11) NOT NULL COMMENT '客户ID',
  `real_name` varchar(50) NOT NULL COMMENT '真实姓名',
  `comment` text NOT NULL COMMENT '评论内容',
  `status` varchar(20) NOT NULL DEFAULT '待审核' COMMENT '状态：已审核、待审核、已拒绝',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_client_id` (`client_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容评论表';

-- 添加外键约束
ALTER TABLE `content` ADD CONSTRAINT `fk_content_region` FOREIGN KEY (`region_id`) REFERENCES `region` (`id`) ON DELETE SET NULL;
ALTER TABLE `content` ADD CONSTRAINT `fk_content_client` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE SET NULL;

-- 评论表外键约束
ALTER TABLE `content_comment` ADD CONSTRAINT `fk_comment_content` FOREIGN KEY (`content_id`) REFERENCES `content` (`id`) ON DELETE CASCADE;
ALTER TABLE `content_comment` ADD CONSTRAINT `fk_comment_client` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE CASCADE;

-- 创建收藏表
CREATE TABLE IF NOT EXISTS `favorite` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
  `client_id` int(11) NOT NULL COMMENT '客户ID',
  `content_id` int(11) NOT NULL COMMENT '内容ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_client_content` (`client_id`, `content_id`),
  KEY `idx_content_id` (`content_id`),
  KEY `idx_client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏表';

-- 收藏表外键约束
ALTER TABLE `favorite` ADD CONSTRAINT `fk_favorite_content` FOREIGN KEY (`content_id`) REFERENCES `content` (`id`) ON DELETE CASCADE;
ALTER TABLE `favorite` ADD CONSTRAINT `fk_favorite_client` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE CASCADE;

-- 创建浏览历史记录表
CREATE TABLE IF NOT EXISTS `browse_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `client_id` int(11) NOT NULL COMMENT '客户ID',
  `content_id` int(11) NOT NULL COMMENT '内容ID',
  `content_type` varchar(20) NOT NULL COMMENT '内容类型',
  `browse_time` datetime NOT NULL COMMENT '浏览时间',
  PRIMARY KEY (`id`),
  KEY `idx_client_time` (`client_id`, `browse_time` DESC),
  KEY `idx_content` (`content_id`, `content_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='浏览历史记录表';

-- 浏览历史记录表外键约束
ALTER TABLE `browse_history` ADD CONSTRAINT `fk_browse_history_client` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE CASCADE;
ALTER TABLE `browse_history` ADD CONSTRAINT `fk_browse_history_content` FOREIGN KEY (`content_id`) REFERENCES `content` (`id`) ON DELETE CASCADE;

-- 创建发布人关注表
CREATE TABLE IF NOT EXISTS `publisher_follow` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '关注ID',
  `client_id` int(11) NOT NULL COMMENT '关注者ID',
  `publisher_id` int(11) NOT NULL COMMENT '被关注发布人ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '关注时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_client_publisher` (`client_id`, `publisher_id`),
  KEY `idx_publisher_id` (`publisher_id`),
  KEY `idx_client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发布人关注表';

-- 发布人关注表外键约束
ALTER TABLE `publisher_follow` ADD CONSTRAINT `fk_publisher_follow_client` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE CASCADE;
ALTER TABLE `publisher_follow` ADD CONSTRAINT `fk_publisher_follow_publisher` FOREIGN KEY (`publisher_id`) REFERENCES `client` (`id`) ON DELETE CASCADE;

-- 导航小程序表
CREATE TABLE IF NOT EXISTS `mini_program` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '导航小程序ID',
  `name` varchar(50) NOT NULL COMMENT '小程序名称',
  `app_id` varchar(50) NOT NULL COMMENT '小程序AppID',
  `logo` varchar(255) DEFAULT NULL COMMENT '小程序图标URL',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用 0:禁用 1:启用',
  `order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值，数字越小排序越靠前',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order` (`order`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='导航小程序表';

-- 初始化示例导航小程序数据
INSERT INTO `mini_program` (`name`, `app_id`, `logo`, `is_enabled`, `order`, `created_at`, `updated_at`) VALUES
('美团外卖', 'wx123456789', 'https://example.com/meituan.png', 1, 1, NOW(), NOW()),
('滴滴出行', 'wx234567890', 'https://example.com/didi.png', 1, 2, NOW(), NOW()),
('京东购物', 'wx345678901', 'https://example.com/jd.png', 1, 3, NOW(), NOW());

-- 轮播图表
CREATE TABLE IF NOT EXISTS `banner` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '轮播图ID',
  `image` varchar(255) NOT NULL COMMENT '轮播图片URL',
  `link_type` varchar(20) NOT NULL COMMENT '跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页',
  `link_url` varchar(255) NOT NULL COMMENT '跳转地址',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用 0:禁用 1:启用',
  `order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值，数字越小排序越靠前',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order` (`order`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='轮播图表';

-- 初始化示例轮播图数据
INSERT INTO `banner` (`image`, `link_type`, `link_url`, `is_enabled`, `order`, `created_at`, `updated_at`) VALUES
('https://example.com/banner1.jpg', 'page', '/pages/index/detail?id=1', 1, 1, NOW(), NOW()),
('https://example.com/banner2.jpg', 'page', '/pages/activity/detail?id=2', 1, 2, NOW(), NOW()),
('https://example.com/banner3.jpg', 'page', '/pages/news/detail?id=3', 1, 3, NOW(), NOW());

-- 活动区域表
CREATE TABLE IF NOT EXISTS `activity_area` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '活动区域ID',
  
  -- 左上模块
  `top_left_title` varchar(50) NOT NULL DEFAULT '' COMMENT '左上模块标题',
  `top_left_description` varchar(200) NOT NULL DEFAULT '' COMMENT '左上模块描述',
  `top_left_link_type` varchar(20) NOT NULL DEFAULT 'page' COMMENT '左上模块跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页',
  `top_left_link_url` varchar(255) NOT NULL DEFAULT '' COMMENT '左上模块跳转地址',
  
  -- 左下模块
  `bottom_left_title` varchar(50) NOT NULL DEFAULT '' COMMENT '左下模块标题',
  `bottom_left_description` varchar(200) NOT NULL DEFAULT '' COMMENT '左下模块描述',
  `bottom_left_link_type` varchar(20) NOT NULL DEFAULT 'page' COMMENT '左下模块跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页',
  `bottom_left_link_url` varchar(255) NOT NULL DEFAULT '' COMMENT '左下模块跳转地址',
  
  -- 右侧模块
  `right_title` varchar(50) NOT NULL DEFAULT '' COMMENT '右侧模块标题',
  `right_description` varchar(200) NOT NULL DEFAULT '' COMMENT '右侧模块描述',
  `right_link_type` varchar(20) NOT NULL DEFAULT 'page' COMMENT '右侧模块跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页',
  `right_link_url` varchar(255) NOT NULL DEFAULT '' COMMENT '右侧模块跳转地址',
  
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='活动区域表';

-- 初始化一条默认记录
INSERT INTO `activity_area` (`id`, `top_left_title`, `top_left_description`, `top_left_link_type`, `top_left_link_url`, 
`bottom_left_title`, `bottom_left_description`, `bottom_left_link_type`, `bottom_left_link_url`, 
`right_title`, `right_description`, `right_link_type`, `right_link_url`, `created_at`, `updated_at`) 
VALUES (1, '天天领音乐会员', '享QQ音乐VIP曲库', 'page', 'pages/vip/music/index', 
'兑换会员权益', '十大会员特权', 'page', 'pages/vip/exchange/index', 
'0元领视频会员', '十大会员特权', 'page', 'pages/vip/video/index', NOW(), NOW());

-- 系统配置表
CREATE TABLE IF NOT EXISTS `system_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `module` varchar(50) NOT NULL COMMENT '模块名称',
  `key` varchar(50) NOT NULL COMMENT '配置键名',
  `value` varchar(255) NOT NULL DEFAULT '' COMMENT '配置值',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '配置描述',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_module_key` (`module`, `key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 初始化导航小程序总开关配置
INSERT INTO `system_config` (`module`, `key`, `value`, `description`) 
VALUES ('mini_program', 'enabled', 'true', '导航小程序总开关：true-启用，false-禁用');

-- 初始化轮播图总开关配置
INSERT INTO `system_config` (`module`, `key`, `value`, `description`) 
VALUES ('banner', 'enabled', 'true', '轮播图总开关：true-启用，false-禁用');

-- 初始化活动区域总开关配置
INSERT INTO `system_config` (`module`, `key`, `value`, `description`) 
VALUES ('activity_area', 'enabled', 'true', '活动区域总开关：true-启用，false-禁用');

-- 初始化首页内页轮播图总开关配置
INSERT INTO `system_config` (`module`, `key`, `value`, `description`) 
VALUES ('inner_banner', 'home_enabled', 'true', '首页内页轮播图总开关：true-启用，false-禁用');

-- 初始化闲置页内页轮播图总开关配置
INSERT INTO `system_config` (`module`, `key`, `value`, `description`) 
VALUES ('inner_banner', 'idle_enabled', 'true', '闲置页内页轮播图总开关：true-启用，false-禁用');

-- 初始化套餐总开关配置
INSERT INTO `system_config` (`module`, `key`, `value`, `description`) 
VALUES ('package', 'top_enabled', 'true', '置顶套餐总开关：true-启用，false-禁用'),
       ('package', 'publish_enabled', 'true', '发布套餐总开关：true-启用，false-禁用');

-- 创建分享设置表
CREATE TABLE IF NOT EXISTS `share_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设置ID',
  `default_share_text` varchar(255) NOT NULL DEFAULT '快来看看这个小程序吧，非常好用！' COMMENT '默认分享语',
  `default_share_image` varchar(255) NOT NULL DEFAULT 'https://example.com/share/default.png' COMMENT '默认分享图片',
  `content_share_text` varchar(255) NOT NULL DEFAULT '我发现了一篇好内容，快来看看吧！' COMMENT '内容页分享语',
  `content_share_image` varchar(255) NOT NULL DEFAULT 'https://example.com/share/content_default.png' COMMENT '内容默认分享图片',
  `home_share_text` varchar(255) NOT NULL DEFAULT '推荐这个小程序，各种信息一应俱全！' COMMENT '首页分享语',
  `home_share_image` varchar(255) NOT NULL DEFAULT 'https://example.com/share/home_default.png' COMMENT '首页默认分享图片',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分享设置表';

-- 初始化分享设置
INSERT INTO `share_settings` (`id`, `default_share_text`, `default_share_image`, `content_share_text`, `content_share_image`, `home_share_text`, `home_share_image`)
VALUES (1, '快来看看这个小程序吧，非常好用！', 'https://example.com/share/default.png', '我发现了一篇好内容，快来看看吧！', 'https://example.com/share/content_default.png', '推荐这个小程序，各种信息一应俱全！', 'https://example.com/share/home_default.png');

-- 内页轮播图表
CREATE TABLE IF NOT EXISTS `inner_page_banner` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `image` varchar(255) NOT NULL COMMENT '轮播图片地址',
  `link_type` varchar(50) NOT NULL COMMENT '跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页',
  `link_url` varchar(255) NOT NULL COMMENT '跳转地址',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用：0-禁用，1-启用',
  `order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值，越小越靠前',
  `banner_type` varchar(20) NOT NULL COMMENT '轮播图类型：home-首页轮播，idle-闲置轮播',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order` (`order`),
  KEY `idx_type` (`banner_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内页轮播图';

-- 添加首页内页轮播图示例数据
INSERT INTO `inner_page_banner` (`image`, `link_type`, `link_url`, `is_enabled`, `order`, `banner_type`, `created_at`, `updated_at`) VALUES
('https://example.com/inner_banner1.jpg', 'page', '/pages/content/detail?id=1', 1, 1, 'home', NOW(), NOW()),
('https://example.com/inner_banner2.jpg', 'page', '/pages/content/detail?id=2', 1, 2, 'home', NOW(), NOW()),
('https://example.com/inner_banner3.jpg', 'page', '/pages/promotion/detail?id=1', 1, 3, 'home', NOW(), NOW());

-- 添加闲置页内页轮播图示例数据
INSERT INTO `inner_page_banner` (`image`, `link_type`, `link_url`, `is_enabled`, `order`, `banner_type`, `created_at`, `updated_at`) VALUES
('https://example.com/idle_banner1.jpg', 'page', '/pages/idle/detail?id=1', 1, 1, 'idle', NOW(), NOW()),
('https://example.com/idle_banner2.jpg', 'page', '/pages/idle/detail?id=2', 1, 2, 'idle', NOW(), NOW()),
('https://example.com/idle_banner3.jpg', 'page', '/pages/idle/category?id=1', 1, 3, 'idle', NOW(), NOW());

-- 创建底部导航栏表
CREATE TABLE IF NOT EXISTS `bottom_tab` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) NOT NULL COMMENT 'Tab名称',
  `icon` varchar(255) NOT NULL COMMENT '未选中状态图标地址',
  `selected_icon` varchar(255) NOT NULL COMMENT '选中状态图标地址',
  `path` varchar(255) NOT NULL COMMENT '页面路径',
  `order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值，越小越靠前',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：0-禁用，1-启用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order` (`order`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='底部导航栏';

-- 添加底部导航栏示例数据
INSERT INTO `bottom_tab` (`name`, `icon`, `selected_icon`, `path`, `order`, `is_enabled`, `created_at`, `updated_at`) VALUES
('首页', 'https://example.com/icons/home.png', 'https://example.com/icons/home_selected.png', '/pages/index/index', 1, 1, NOW(), NOW()),
('分类', 'https://example.com/icons/category.png', 'https://example.com/icons/category_selected.png', '/pages/category/index', 2, 1, NOW(), NOW()),
('发布', 'https://example.com/icons/publish.png', 'https://example.com/icons/publish_selected.png', '/pages/publish/index', 3, 1, NOW(), NOW()),
('消息', 'https://example.com/icons/message.png', 'https://example.com/icons/message_selected.png', '/pages/message/index', 4, 1, NOW(), NOW()),
('我的', 'https://example.com/icons/profile.png', 'https://example.com/icons/profile_selected.png', '/pages/profile/index', 5, 1, NOW(), NOW());

-- 创建小程序基础设置表
CREATE TABLE IF NOT EXISTS `mini_program_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设置ID',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '小程序名称',
  `description` varchar(500) NOT NULL DEFAULT '' COMMENT '小程序描述',
  `logo` varchar(255) NOT NULL DEFAULT '' COMMENT '小程序Logo地址',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小程序基础设置表';

-- 初始化小程序基础设置
INSERT INTO `mini_program_settings` (`id`, `name`, `description`, `logo`, `created_at`, `updated_at`)
VALUES (1, '城市服务小程序', '提供本地化城市生活服务的小程序平台', 'https://example.com/logo.png', NOW(), NOW());

-- 创建广告设置表
CREATE TABLE IF NOT EXISTS `ad_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设置ID',
  `enable_wx_ad` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用微信广告：0-禁用，1-启用',
  `rewarded_video_ad_id` varchar(100) NOT NULL DEFAULT '' COMMENT '激励视频广告位ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='广告设置表';

-- 初始化广告设置
INSERT INTO `ad_settings` (`id`, `enable_wx_ad`, `rewarded_video_ad_id`, `created_at`, `updated_at`)
VALUES (1, 0, '', NOW(), NOW());

-- 创建奖励设置表
CREATE TABLE IF NOT EXISTS `reward_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设置ID',
  `enable_reward` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用奖励功能：0-禁用，1-启用',
  `first_view_min_reward_min` int(11) NOT NULL DEFAULT '60' COMMENT '首次观看广告最小奖励(分钟)',
  `first_view_max_reward_day` decimal(5,2) NOT NULL DEFAULT '7.00' COMMENT '首次观看广告最大奖励(天)',
  `single_ad_min_reward_min` int(11) NOT NULL DEFAULT '5' COMMENT '单次广告最小奖励(分钟)',
  `single_ad_max_reward_day` decimal(5,2) NOT NULL DEFAULT '1.00' COMMENT '单次广告最大奖励(天)',
  `daily_reward_limit` int(11) NOT NULL DEFAULT '10' COMMENT '每日奖励次数上限',
  `daily_max_accumulated_day` decimal(5,2) NOT NULL DEFAULT '2.00' COMMENT '每日最大累计奖励(天)',
  `reward_expiration_days` int(11) NOT NULL DEFAULT '30' COMMENT '奖励过期天数',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='奖励设置表';

-- 初始化奖励设置
INSERT INTO `reward_settings` (`id`, `enable_reward`, `first_view_min_reward_min`, `first_view_max_reward_day`, 
                             `single_ad_min_reward_min`, `single_ad_max_reward_day`, 
                             `daily_reward_limit`, `daily_max_accumulated_day`, `reward_expiration_days`, 
                             `created_at`, `updated_at`)
VALUES (1, 0, 60, 7.00, 5, 1.00, 10, 2.00, 30, NOW(), NOW());

-- 创建订单支付映射表
CREATE TABLE IF NOT EXISTS `order_pay_map` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `original_order_no` varchar(32) NOT NULL COMMENT '原始订单号',
  `temp_order_no` varchar(64) NOT NULL COMMENT '临时订单号',
  `client_id` int(11) NOT NULL COMMENT '客户ID',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态：0-未处理 1-已处理',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_temp_order_no` (`temp_order_no`),
  KEY `idx_original_order_no` (`original_order_no`),
  KEY `idx_client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单支付映射表';

-- 创建协议设置表
CREATE TABLE IF NOT EXISTS `agreement_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设置ID',
  `privacy_policy` text COMMENT '隐私政策内容',
  `user_agreement` text COMMENT '用户协议内容',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='协议设置表';

-- 初始化默认协议设置
INSERT INTO `agreement_settings` (`id`, `privacy_policy`, `user_agreement`)
VALUES (1, '这是默认隐私政策内容，请在管理后台修改', '这是默认用户协议内容，请在管理后台修改')
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;

-- 创建奖励记录表
CREATE TABLE IF NOT EXISTS `reward_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `client_id` int(11) NOT NULL COMMENT '客户ID',
  `reward_minutes` int(11) NOT NULL COMMENT '奖励时长(分钟)',
  `reward_days` decimal(5,2) NOT NULL COMMENT '奖励时长(天)',
  `is_first_view` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否首次观看：0-否，1-是',
  `remaining_minutes` int(11) NOT NULL COMMENT '剩余时长(分钟)',
  `total_reward_minutes` int(11) NOT NULL COMMENT '累计获得奖励(分钟)',
  `used_minutes` int(11) NOT NULL DEFAULT '0' COMMENT '已使用时长(分钟)',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0-已过期，1-有效',
  `expire_at` datetime NOT NULL COMMENT '过期时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_client_id` (`client_id`),
  KEY `idx_status` (`status`),
  KEY `idx_expire_at` (`expire_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='奖励记录表';

-- 创建客户时长表
CREATE TABLE IF NOT EXISTS `client_duration` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `client_id` int(11) NOT NULL COMMENT '客户ID',
  `client_name` varchar(50) NOT NULL COMMENT '客户名称',
  `remaining_duration` varchar(50) NOT NULL COMMENT '剩余时长，格式如: 3天18小时42分钟',
  `total_duration` varchar(50) NOT NULL COMMENT '累计获得时长，格式如: 3天18小时42分钟',
  `used_duration` varchar(50) NOT NULL COMMENT '已使用时长，格式如: 4天8小时59分钟',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户时长表';

-- 添加客户时长测试数据
INSERT INTO `client_duration` (`client_id`, `client_name`, `remaining_duration`, `total_duration`, `used_duration`, `created_at`) VALUES
(1, '测试客户', '3天18小时42分钟', '5天10小时30分钟', '1天15小时48分钟', NOW());

-- 创建商品表
CREATE TABLE IF NOT EXISTS `product` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `name` varchar(100) NOT NULL COMMENT '商品名称',
  `category_id` int(11) NOT NULL COMMENT '分类ID',
  `category_name` varchar(50) NOT NULL COMMENT '分类名称',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '商品价格',
  `duration` int(11) NOT NULL DEFAULT '0' COMMENT '所需时长(天)',
  `stock` int(11) NOT NULL DEFAULT '0' COMMENT '库存数量',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0:未上架 1:已上架 2:已售罄',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值，数字越小排序越靠前',
  `description` text COMMENT '商品描述',
  `thumbnail` varchar(255) DEFAULT NULL COMMENT '缩略图URL',
  `images` text COMMENT '商品图片URLs，JSON格式',
  `tags` varchar(255) DEFAULT NULL COMMENT '商品标签，多个标签用逗号分隔',
  `sales` int(11) NOT NULL DEFAULT '0' COMMENT '销量',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- 插入示例商品数据
INSERT INTO `product` (`name`, `category_id`, `category_name`, `price`, `duration`, `stock`, `status`, `sort_order`, `description`, `thumbnail`, `images`, `tags`, `sales`) VALUES
('示例商品1', 1, '手机数码', 199.00, 3, 100, 1, 1, '这是示例商品1的详细描述', 'https://example.com/products/1/thumbnail.jpg', '[\"https://example.com/products/1/image1.jpg\",\"https://example.com/products/1/image2.jpg\"]', '手机,数码', 10),
('示例商品2', 2, '家用电器', 299.00, 7, 50, 1, 2, '这是示例商品2的详细描述', 'https://example.com/products/2/thumbnail.jpg', '[\"https://example.com/products/2/image1.jpg\",\"https://example.com/products/2/image2.jpg\"]', '家用,电器', 5),
('示例商品3', 3, '电脑办公', 399.00, 15, 0, 2, 3, '这是示例商品3的详细描述', 'https://example.com/products/3/thumbnail.jpg', '[\"https://example.com/products/3/image1.jpg\",\"https://example.com/products/3/image2.jpg\"]', '电脑,办公', 20),
('示例商品4', 4, '服装鞋帽', 99.00, 1, 200, 0, 4, '这是示例商品4的详细描述', 'https://example.com/products/4/thumbnail.jpg', '[\"https://example.com/products/4/image1.jpg\",\"https://example.com/products/4/image2.jpg\"]', '服装,鞋帽', 0);

-- 创建兑换列表表
CREATE TABLE IF NOT EXISTS `exchange_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '兑换记录ID',
  `client_id` int(11) NOT NULL COMMENT '客户ID',
  `client_name` varchar(50) NOT NULL COMMENT '客户名称',
  `recharge_account` varchar(100) NOT NULL COMMENT '充值账号',
  `product_name` varchar(100) NOT NULL COMMENT '商品名称',
  `duration` int(11) NOT NULL COMMENT '消耗时长(天)',
  `exchange_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '兑换时间',
  `status` varchar(20) NOT NULL DEFAULT '处理中' COMMENT '状态：处理中、已完成、失败',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_client_id` (`client_id`),
  KEY `idx_status` (`status`),
  KEY `idx_exchange_time` (`exchange_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='兑换记录表';

-- 创建商城分类表
CREATE TABLE IF NOT EXISTS `shop_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序值',
  `product_count` int(11) NOT NULL DEFAULT '0' COMMENT '商品数量',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1启用，0禁用',
  `image` varchar(255) DEFAULT NULL COMMENT '分类图片URL',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商城分类表';

-- 添加外键约束
ALTER TABLE `exchange_record` ADD CONSTRAINT `fk_exchange_client` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE CASCADE;

-- 创建用户消息表
CREATE TABLE IF NOT EXISTS `client_message` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '消息ID',
  `sender_id` int(11) NOT NULL COMMENT '发送者ID',
  `sender_name` varchar(50) NOT NULL COMMENT '发送者名称',
  `receiver_id` int(11) NOT NULL COMMENT '接收者ID',
  `receiver_name` varchar(50) NOT NULL COMMENT '接收者名称',
  `content` text NOT NULL COMMENT '消息内容',
  `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读：0未读，1已读',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0删除，1正常',
  `sender_visible` tinyint(1) NOT NULL DEFAULT '1' COMMENT '对发送者是否可见：0不可见，1可见',
  `receiver_visible` tinyint(1) NOT NULL DEFAULT '1' COMMENT '对接收者是否可见：0不可见，1可见',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_sender_id` (`sender_id`),
  KEY `idx_receiver_id` (`receiver_id`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户消息表';

-- 创建会话表
CREATE TABLE IF NOT EXISTS `client_conversation` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '会话ID',
  `client_id` int(11) NOT NULL COMMENT '用户ID',
  `target_id` int(11) NOT NULL COMMENT '对方ID',
  `target_name` varchar(50) NOT NULL COMMENT '对方名称',
  `target_avatar` varchar(255) DEFAULT NULL COMMENT '对方头像',
  `last_message` varchar(255) DEFAULT NULL COMMENT '最后一条消息内容',
  `unread_count` int(11) NOT NULL DEFAULT '0' COMMENT '未读消息数',
  `last_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '最后消息时间',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0删除，1正常',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_client_target` (`client_id`,`target_id`),
  KEY `idx_client_id` (`client_id`),
  KEY `idx_target_id` (`target_id`),
  KEY `idx_last_time` (`last_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话表';

-- 创建专属管家表
CREATE TABLE IF NOT EXISTS `butler` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '管家ID',
  `image_url` varchar(255) NOT NULL COMMENT '管家图片URL',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='专属管家表';