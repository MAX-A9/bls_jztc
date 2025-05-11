-- 兑换商品表
CREATE TABLE IF NOT EXISTS `tb_exchange_product` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `name` varchar(100) NOT NULL COMMENT '商品名称',
  `description` text DEFAULT NULL COMMENT '商品描述',
  `image` varchar(255) DEFAULT NULL COMMENT '商品图片URL',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '商品价格',
  `stock` int(11) NOT NULL DEFAULT '0' COMMENT '库存数量',
  `category_id` int(11) NOT NULL COMMENT '商品分类ID',
  `status` enum('active','inactive','soldout') NOT NULL DEFAULT 'inactive' COMMENT '商品状态：active(已上架)、inactive(未上架)、soldout(已售罄)',
  `time_required` int(11) NOT NULL DEFAULT '0' COMMENT '兑换所需时长(小时)',
  `order_sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序字段',
  `exchange_count` int(11) NOT NULL DEFAULT '0' COMMENT '已兑换次数',
  `is_virtual` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否虚拟商品：1是，0否',
  `is_coupon` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否优惠券：1是，0否',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_order_sort` (`order_sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='兑换商品表';

-- 商城分类表
CREATE TABLE IF NOT EXISTS `tb_exchange_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `description` varchar(255) DEFAULT NULL COMMENT '分类描述',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1启用，0禁用',
  `order_sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序字段',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_order_sort` (`order_sort`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商城分类表';

-- 兑换记录表
CREATE TABLE IF NOT EXISTS `tb_exchange_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '兑换记录ID',
  `exchange_no` varchar(50) NOT NULL COMMENT '兑换单号',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `user_name` varchar(50) NOT NULL COMMENT '用户名称',
  `product_id` int(11) NOT NULL COMMENT '商品ID',
  `product_name` varchar(100) NOT NULL COMMENT '商品名称',
  `quantity` int(11) NOT NULL DEFAULT '1' COMMENT '兑换数量',
  `phone` varchar(20) DEFAULT NULL COMMENT '充值账号/手机号',
  `address` text DEFAULT NULL COMMENT '收货地址',
  `status` enum('pending','processing','completed','failed') NOT NULL DEFAULT 'pending' COMMENT '兑换状态：pending(处理中)、processing(处理中)、completed(已完成)、failed(失败)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `admin_id` int(11) DEFAULT NULL COMMENT '处理管理员ID',
  `admin_name` varchar(50) DEFAULT NULL COMMENT '处理管理员姓名',
  `process_time` datetime DEFAULT NULL COMMENT '处理时间',
  `exchange_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '兑换时间',
  `complete_time` datetime DEFAULT NULL COMMENT '完成时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_exchange_no` (`exchange_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`),
  KEY `idx_exchange_time` (`exchange_time`),
  KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='兑换记录表';

-- 兑换日志表
CREATE TABLE IF NOT EXISTS `tb_exchange_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `exchange_id` int(11) NOT NULL COMMENT '兑换记录ID',
  `admin_id` int(11) DEFAULT NULL COMMENT '管理员ID',
  `admin_name` varchar(50) NOT NULL COMMENT '管理员名称',
  `action` varchar(50) NOT NULL COMMENT '操作类型：view(查看)、process(处理)、complete(完成)、fail(标记失败)',
  `status_before` varchar(50) DEFAULT NULL COMMENT '操作前状态',
  `status_after` varchar(50) DEFAULT NULL COMMENT '操作后状态',
  `description` varchar(255) NOT NULL COMMENT '操作描述',
  `ip` varchar(50) NOT NULL COMMENT '操作IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_exchange_id` (`exchange_id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_exchange_log_record` FOREIGN KEY (`exchange_id`) REFERENCES `tb_exchange_record` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='兑换操作日志表';

-- 优惠券表
CREATE TABLE IF NOT EXISTS `tb_coupon` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '优惠券ID',
  `code` varchar(50) NOT NULL COMMENT '优惠券码',
  `name` varchar(100) NOT NULL COMMENT '优惠券名称',
  `type` enum('fixed','percent') NOT NULL DEFAULT 'fixed' COMMENT '优惠券类型：fixed(固定金额)、percent(百分比)',
  `value` decimal(10,2) NOT NULL COMMENT '优惠券值：固定金额或百分比',
  `min_order_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '最低订单金额',
  `valid_from` datetime DEFAULT NULL COMMENT '有效期开始',
  `valid_to` datetime DEFAULT NULL COMMENT '有效期结束',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `exchange_id` int(11) DEFAULT NULL COMMENT '关联的兑换记录ID',
  `status` enum('unused','used','expired') NOT NULL DEFAULT 'unused' COMMENT '状态：unused(未使用)、used(已使用)、expired(已过期)',
  `used_time` datetime DEFAULT NULL COMMENT '使用时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_exchange_id` (`exchange_id`),
  KEY `idx_status` (`status`),
  KEY `idx_valid_from_to` (`valid_from`,`valid_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券表';

-- 示例数据：商城分类
INSERT INTO `tb_exchange_category` (`name`, `description`, `is_enabled`, `order_sort`) VALUES 
('优惠券', '各类优惠券', 1, 1),
('票券', '电影票、演出票等', 1, 2),
('餐饮', '餐饮类礼品券', 1, 3),
('实体商品', '各类实体奖品', 1, 4),
('限时活动', '限时活动奖品', 0, 5);

-- 示例数据：兑换商品
INSERT INTO `tb_exchange_product` (`name`, `description`, `image`, `price`, `stock`, `category_id`, `status`, `time_required`, `order_sort`, `is_virtual`, `is_coupon`) VALUES 
('优惠券-10元', '可在商城中使用的10元优惠券', 'https://img.example.com/coupon_10.jpg', 0, 1000, 1, 'active', 2, 1, 1, 1),
('优惠券-50元', '可在商城中使用的50元优惠券', 'https://img.example.com/coupon_50.jpg', 0, 500, 1, 'active', 8, 2, 1, 1),
('电影票兑换券', '可兑换一张电影票', 'https://img.example.com/movie_ticket.jpg', 0, 200, 2, 'active', 12, 3, 1, 0),
('咖啡券', '可兑换一杯中杯咖啡', 'https://img.example.com/coffee.jpg', 0, 300, 3, 'active', 5, 4, 1, 0),
('实体商品A', '精美实体商品', 'https://img.example.com/product_a.jpg', 0, 50, 4, 'active', 72, 5, 0, 0);

-- 示例数据：兑换记录
INSERT INTO `tb_exchange_record` (`exchange_no`, `user_id`, `user_name`, `product_id`, `product_name`, `quantity`, `phone`, `status`, `remark`, `exchange_time`) VALUES 
('EX20240305001', 101, '张三', 1, '优惠券-50元', 1, '13800138000', 'completed', '自动发放成功', '2024-03-05 14:22:35'),
('EX20240304001', 102, '李四', 2, '高级会员月卡', 1, '13800138000', 'completed', '自动充值成功', '2024-03-04 09:15:42'),
('EX20240303001', 103, '王五', 3, '精美文具套装', 2, '13800138000', 'pending', '周末送货', '2024-03-03 16:45:18'),
('EX20240302001', 104, '赵六', 4, '限量版徽章', 3, '13800138000', 'failed', '库存不足', '2024-03-02 11:30:25'),
('EX20240301001', 105, '钱七', 5, '定制T恤', 1, '13800138000', 'pending', '尺码：XL', '2024-03-01 15:20:33');

-- 示例数据：优惠券
INSERT INTO `tb_coupon` (`code`, `name`, `type`, `value`, `min_order_amount`, `valid_from`, `valid_to`, `user_id`, `exchange_id`, `status`) VALUES 
('COUPON1001', '优惠券-10元', 'fixed', 10.00, 100.00, '2024-03-01 00:00:00', '2024-06-30 23:59:59', 101, 1, 'unused'),
('COUPON1002', '优惠券-50元', 'fixed', 50.00, 300.00, '2024-03-01 00:00:00', '2024-06-30 23:59:59', 102, 2, 'unused'),
('COUPON1003', '优惠券-10%', 'percent', 10.00, 50.00, '2024-03-01 00:00:00', '2024-06-30 23:59:59', 103, 3, 'unused'); 