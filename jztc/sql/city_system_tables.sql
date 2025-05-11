-- 同城系统相关表

-- 首页布局设置表
CREATE TABLE IF NOT EXISTS `tb_banner` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Banner ID',
  `image` varchar(255) NOT NULL COMMENT '图片路径',
  `link_type` varchar(20) NOT NULL DEFAULT 'page' COMMENT '链接类型：page(页面)、miniprogram(小程序)、url(外部链接)',
  `link_url` varchar(255) NOT NULL COMMENT '链接地址',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1是，0否',
  `order_sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序字段',
  `position` varchar(20) NOT NULL DEFAULT 'home' COMMENT '位置：home(首页)、category(分类页)等',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_position` (`position`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_order_sort` (`order_sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Banner表';

-- 导航小程序表
CREATE TABLE IF NOT EXISTS `tb_miniprogram_nav` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '导航ID',
  `name` varchar(50) NOT NULL COMMENT '小程序名称',
  `app_id` varchar(50) NOT NULL COMMENT '小程序AppID',
  `logo` varchar(255) NOT NULL COMMENT 'Logo图片路径',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1是，0否',
  `order_sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序字段',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_order_sort` (`order_sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='导航小程序表';

-- 内页设置表
CREATE TABLE IF NOT EXISTS `tb_page_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设置ID',
  `page_key` varchar(50) NOT NULL COMMENT '页面键名',
  `setting_key` varchar(50) NOT NULL COMMENT '设置键名',
  `setting_value` text COMMENT '设置值',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注说明',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_page_setting_key` (`page_key`,`setting_key`),
  KEY `idx_page_key` (`page_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内页设置表';

-- 底部Tab设置表
CREATE TABLE IF NOT EXISTS `tb_tab_bar` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Tab ID',
  `title` varchar(20) NOT NULL COMMENT 'Tab标题',
  `icon` varchar(255) NOT NULL COMMENT '图标路径',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1是，0否',
  `order_sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序字段',
  `page_path` varchar(100) NOT NULL COMMENT '页面路径',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_is_active` (`is_active`),
  KEY `idx_order_sort` (`order_sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='底部Tab设置表';

-- 基础设置表
CREATE TABLE IF NOT EXISTS `tb_miniprogram_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `app_id` varchar(50) NOT NULL COMMENT '小程序AppID',
  `app_secret` varchar(100) NOT NULL COMMENT '小程序AppSecret',
  `name` varchar(50) NOT NULL COMMENT '小程序名称',
  `description` text COMMENT '小程序描述',
  `logo` varchar(255) DEFAULT NULL COMMENT 'Logo图片路径',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1是，0否',
  `privacy_policy` text COMMENT '隐私政策',
  `user_agreement` text COMMENT '用户协议',
  `mch_id` varchar(50) DEFAULT NULL COMMENT '商户号',
  `mch_key` varchar(100) DEFAULT NULL COMMENT '商户密钥',
  `pay_notify_url` varchar(255) DEFAULT NULL COMMENT '支付回调URL',
  `enable_payment` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用支付：1是，0否',
  `enable_ads` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用广告：1是，0否',
  `ad_app_id` varchar(50) DEFAULT NULL COMMENT '广告AppID',
  `banner_ad_unit_id` varchar(50) DEFAULT NULL COMMENT 'Banner广告位ID',
  `interstitial_ad_unit_id` varchar(50) DEFAULT NULL COMMENT '插屏广告位ID',
  `video_ad_unit_id` varchar(50) DEFAULT NULL COMMENT '视频广告位ID',
  `native_ad_unit_id` varchar(50) DEFAULT NULL COMMENT '原生广告位ID',
  `ad_frequency` varchar(20) DEFAULT 'low' COMMENT '广告频率：low(低)、medium(中)、high(高)、custom(自定义)',
  `custom_ad_frequency` int(11) DEFAULT '20' COMMENT '自定义广告频率',
  `ad_positions` varchar(255) DEFAULT 'home,detail' COMMENT '广告位置，英文逗号分隔',
  `enable_rewards` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用奖励：1是，0否',
  `reward_min_minutes` int(11) DEFAULT '5' COMMENT '最小奖励分钟数',
  `reward_max_days` int(11) DEFAULT '1' COMMENT '最大奖励天数',
  `daily_reward_limit` int(11) DEFAULT '10' COMMENT '每日奖励上限',
  `max_daily_reward_time` int(11) DEFAULT '2' COMMENT '每日最大奖励小时数',
  `reward_expire_days` int(11) DEFAULT '30' COMMENT '奖励过期天数',
  `first_time_min_minutes` int(11) DEFAULT '60' COMMENT '首次使用最小分钟数',
  `first_time_max_value` int(11) DEFAULT '7' COMMENT '首次使用最大奖励值',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_app_id` (`app_id`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小程序配置表';

-- 发布设置 - 置顶套餐表
CREATE TABLE IF NOT EXISTS `tb_top_package` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '套餐ID',
  `title` varchar(50) NOT NULL COMMENT '套餐名称',
  `description` varchar(255) DEFAULT NULL COMMENT '套餐描述',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '套餐价格',
  `duration` int(11) NOT NULL COMMENT '有效时长(小时)',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1是，0否',
  `order_sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序字段',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_order_sort` (`order_sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='置顶套餐表';

-- 发布设置 - 发布套餐表
CREATE TABLE IF NOT EXISTS `tb_publish_package` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '套餐ID',
  `title` varchar(50) NOT NULL COMMENT '套餐名称',
  `description` varchar(255) DEFAULT NULL COMMENT '套餐描述',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '套餐价格',
  `duration` int(11) NOT NULL COMMENT '有效时长(天)',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1是，0否',
  `order_sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序字段',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_is_enabled` (`is_enabled`),
  KEY `idx_order_sort` (`order_sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发布套餐表';

-- 添加演示数据 - Banner表
INSERT INTO `tb_banner` (`image`, `link_type`, `link_url`, `is_enabled`, `order_sort`, `position`) VALUES
('https://img.example.com/banner1.png', 'page', '/pages/index/detail?id=1', 1, 1, 'home'),
('https://img.example.com/banner2.png', 'page', '/pages/activity/detail?id=5', 1, 2, 'home'),
('https://img.example.com/banner3.png', 'miniprogram', 'wx123456789', 1, 3, 'home'),
('https://img.example.com/home-banner1.png', 'page', '/pages/home/detail?id=1', 1, 1, 'home_page'),
('https://img.example.com/home-banner2.png', 'page', '/pages/home/activity?id=5', 1, 2, 'home_page'),
('https://img.example.com/idle-banner1.png', 'page', '/pages/idle/detail?id=1', 1, 1, 'idle_page'),
('https://img.example.com/idle-banner2.png', 'page', '/pages/idle/category?id=5', 1, 2, 'idle_page');

-- 添加演示数据 - 导航小程序表
INSERT INTO `tb_miniprogram_nav` (`name`, `app_id`, `logo`, `is_enabled`, `order_sort`) VALUES
('美团外卖', 'wx123456789', 'https://img.example.com/meituan.png', 1, 1),
('滴滴出行', 'wx234567890', 'https://img.example.com/didi.png', 1, 2),
('京东购物', 'wx345678901', 'https://img.example.com/jd.png', 0, 3),
('饿了么', 'wx456789012', 'https://img.example.com/eleme.png', 1, 4);

-- 添加演示数据 - 内页设置表
INSERT INTO `tb_page_settings` (`page_key`, `setting_key`, `setting_value`, `remark`) VALUES
('home', 'enable_banner', '1', '启用首页Banner'),
('home', 'banner_height', '200', 'Banner高度(px)'),
('home', 'show_search', '1', '显示搜索框'),
('idle', 'enable_banner', '1', '启用闲置页Banner'),
('idle', 'banner_height', '180', 'Banner高度(px)'),
('idle', 'show_filter', '1', '显示筛选组件');

-- 添加演示数据 - 底部Tab设置表
INSERT INTO `tb_tab_bar` (`title`, `icon`, `is_active`, `order_sort`, `page_path`) VALUES
('首页', 'home-icon.png', 1, 1, 'pages/index/index'),
('闲置', 'idle-icon.png', 1, 2, 'pages/idle/index'),
('发布', 'publish-icon.png', 1, 3, 'pages/publish/index'),
('消息', 'message-icon.png', 1, 4, 'pages/message/index'),
('我的', 'profile-icon.png', 1, 5, 'pages/profile/index');

-- 添加演示数据 - 小程序配置表
INSERT INTO `tb_miniprogram_config` (`app_id`, `app_secret`, `name`, `description`, `logo`, `is_enabled`, `privacy_policy`, `user_agreement`, `mch_id`, `pay_notify_url`, `enable_payment`) VALUES
('wx1234567890abcdef', '071d454a693afd029cb4e2d851e9a23f', '同城生活服务', '提供本地生活服务、商家信息和社区互动的小程序', 'https://example.com/miniprogram-logo.png', 1, '# 隐私政策\n\n我们非常重视您的隐私...', '# 用户协议\n\n欢迎使用我们的服务...', '1230000109', 'https://api.example.com/pay/notify', 0);

-- 添加演示数据 - 置顶套餐表
INSERT INTO `tb_top_package` (`title`, `description`, `price`, `duration`, `is_enabled`, `order_sort`) VALUES
('置顶一小时', '内容将在列表顶部显示一小时', 2.00, 1, 1, 1),
('置顶一天', '内容将在列表顶部显示24小时', 10.00, 24, 1, 2),
('置顶一周', '内容将在列表顶部显示7天', 50.00, 168, 1, 3);

-- 添加演示数据 - 发布套餐表
INSERT INTO `tb_publish_package` (`title`, `description`, `price`, `duration`, `is_enabled`, `order_sort`) VALUES
('10天', '发布内容可展示10天', 5.00, 10, 1, 1),
('20天', '发布内容可展示20天', 10.00, 20, 1, 2),
('30天', '发布内容可展示30天', 15.00, 30, 1, 3); 