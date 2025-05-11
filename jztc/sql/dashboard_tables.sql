-- 统计数据表
CREATE TABLE IF NOT EXISTS `tb_statistics` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '统计ID',
  `date` date NOT NULL COMMENT '统计日期',
  `registration_count` int(11) NOT NULL DEFAULT '0' COMMENT '注册用户数',
  `exchange_count` int(11) NOT NULL DEFAULT '0' COMMENT '兑换次数',
  `publish_count` int(11) NOT NULL DEFAULT '0' COMMENT '发布内容数',
  `income_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '收益金额',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date` (`date`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='每日统计数据表';

-- 趋势数据表
CREATE TABLE IF NOT EXISTS `tb_trend_data` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '趋势ID',
  `date` date NOT NULL COMMENT '日期',
  `category` varchar(20) NOT NULL COMMENT '类别：registration(注册用户)、exchange(兑换次数)、publish(发布内容)、income(收益金额)',
  `value` int(11) NOT NULL DEFAULT '0' COMMENT '数值',
  `change_percent` decimal(5,2) DEFAULT NULL COMMENT '变化百分比',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_category` (`date`,`category`),
  KEY `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='趋势数据表';

-- 订单表
CREATE TABLE IF NOT EXISTS `tb_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(30) NOT NULL COMMENT '订单编号',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `customer` varchar(50) NOT NULL COMMENT '客户名称',
  `amount` decimal(10,2) NOT NULL COMMENT '订单金额',
  `status` enum('pending','paid','shipped','completed','cancelled') NOT NULL DEFAULT 'pending' COMMENT '订单状态',
  `payment_method` varchar(20) DEFAULT NULL COMMENT '支付方式',
  `product` varchar(255) NOT NULL COMMENT '商品名称',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 兑换表
CREATE TABLE IF NOT EXISTS `tb_exchange` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '兑换ID',
  `exchange_no` varchar(30) NOT NULL COMMENT '兑换编号',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `user_name` varchar(50) NOT NULL COMMENT '用户名称',
  `product_id` int(11) NOT NULL COMMENT '商品ID',
  `product_name` varchar(255) NOT NULL COMMENT '商品名称',
  `points` int(11) NOT NULL COMMENT '消耗积分',
  `status` enum('pending','processing','completed','cancelled') NOT NULL DEFAULT 'pending' COMMENT '兑换状态',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_exchange_no` (`exchange_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分兑换表';

-- 仪表盘日志表
CREATE TABLE IF NOT EXISTS `tb_dashboard_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `admin_id` int(11) DEFAULT NULL COMMENT '管理员ID',
  `admin_name` varchar(50) NOT NULL COMMENT '管理员名称',
  `action` varchar(50) NOT NULL COMMENT '操作类型：view(查看)、export(导出)、filter(筛选)',
  `time_range` varchar(20) NOT NULL COMMENT '时间范围：本周、本月、本年',
  `content_type` varchar(50) DEFAULT NULL COMMENT '内容类型',
  `ip` varchar(50) NOT NULL COMMENT '操作IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_action` (`action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='仪表盘操作日志表';

-- 样例数据：统计数据
INSERT INTO `tb_statistics` (`date`, `registration_count`, `exchange_count`, `publish_count`, `income_amount`) VALUES 
(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 456, 15, 652, 3500.00),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 487, 18, 598, 3800.00),
(DATE_SUB(CURDATE(), INTERVAL 4 DAY), 512, 22, 612, 4100.50),
(DATE_SUB(CURDATE(), INTERVAL 3 DAY), 498, 12, 578, 3650.00),
(DATE_SUB(CURDATE(), INTERVAL 2 DAY), 523, 17, 634, 4200.00),
(DATE_SUB(CURDATE(), INTERVAL 1 DAY), 545, 13, 648, 4425.50),
(CURDATE(), 576, 24, 678, 5200.00);

-- 样例数据：趋势数据
INSERT INTO `tb_trend_data` (`date`, `category`, `value`, `change_percent`) VALUES 
(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 'registration', 456, 5.20),
(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 'exchange', 15, 3.80),
(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 'publish', 652, -2.40),
(DATE_SUB(CURDATE(), INTERVAL 6 DAY), 'income', 3500, 7.50),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 'registration', 487, 6.80),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 'exchange', 18, 5.20),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 'publish', 598, -1.80),
(DATE_SUB(CURDATE(), INTERVAL 5 DAY), 'income', 3800, 8.20),
(CURDATE(), 'registration', 576, 15.20),
(CURDATE(), 'exchange', 24, 8.00),
(CURDATE(), 'publish', 678, -5.00),
(CURDATE(), 'income', 5200, 12.00);

-- 样例数据：订单
INSERT INTO `tb_order` (`order_no`, `user_id`, `customer`, `amount`, `status`, `payment_method`, `product`, `created_at`) VALUES 
('ORD20230501001', 1001, '张三', 128.00, 'completed', '微信支付', '置顶一小时', DATE_SUB(NOW(), INTERVAL 2 DAY)),
('ORD20230501002', 1002, '李四', 256.50, 'pending', '支付宝', '商品B', DATE_SUB(NOW(), INTERVAL 2 DAY)),
('ORD20230502001', 1003, '王五', 78.50, 'paid', '微信支付', '商品C', DATE_SUB(NOW(), INTERVAL 1 DAY)),
('ORD20230502002', 1004, '赵六', 199.00, 'shipped', '支付宝', '商品D', DATE_SUB(NOW(), INTERVAL 1 DAY)),
('ORD20230503001', 1005, '钱七', 59.90, 'cancelled', '银联', '商品E', CURDATE());

-- 样例数据：兑换
INSERT INTO `tb_exchange` (`exchange_no`, `user_id`, `user_name`, `product_id`, `product_name`, `points`, `status`, `created_at`) VALUES 
('EX20230501001', 1001, '张三', 101, '高级会员月卡', 1200, 'completed', DATE_SUB(NOW(), INTERVAL 2 DAY)),
('EX20230501002', 1002, '李四', 102, '精品课程礼包', 800, 'pending', DATE_SUB(NOW(), INTERVAL 2 DAY)),
('EX20230502001', 1003, '王五', 103, '限量版数字藏品', 2000, 'processing', DATE_SUB(NOW(), INTERVAL 1 DAY)),
('EX20230502002', 1004, '赵六', 104, '专属头像框', 500, 'completed', DATE_SUB(NOW(), INTERVAL 1 DAY)),
('EX20230503001', 1005, '钱七', 105, '优惠券礼包', 300, 'cancelled', CURDATE()); 