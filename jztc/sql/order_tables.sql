-- 订单主表
CREATE TABLE IF NOT EXISTS `tb_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(30) NOT NULL COMMENT '订单编号',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `customer` varchar(50) NOT NULL COMMENT '客户名称',
  `amount` decimal(10,2) NOT NULL COMMENT '订单金额',
  `status` enum('pending','paid','shipped','completed','cancelled') NOT NULL DEFAULT 'pending' COMMENT '订单状态:pending(待付款),paid(已付款),shipped(已发货),completed(已完成),cancelled(已取消)',
  `payment_method` varchar(20) DEFAULT NULL COMMENT '支付方式',
  `product_id` int(11) DEFAULT NULL COMMENT '商品ID',
  `product` varchar(255) NOT NULL COMMENT '商品名称',
  `order_sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序字段',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `ship_time` datetime DEFAULT NULL COMMENT '发货时间',
  `complete_time` datetime DEFAULT NULL COMMENT '完成时间',
  `cancel_time` datetime DEFAULT NULL COMMENT '取消时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单主表';

-- 订单支付信息表
CREATE TABLE IF NOT EXISTS `tb_order_payment` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '支付ID',
  `order_id` int(11) NOT NULL COMMENT '订单ID',
  `payment_no` varchar(50) NOT NULL COMMENT '支付流水号',
  `payment_method` varchar(20) NOT NULL COMMENT '支付方式',
  `amount` decimal(10,2) NOT NULL COMMENT '支付金额',
  `status` enum('pending','success','failed','refunded') NOT NULL DEFAULT 'pending' COMMENT '支付状态',
  `transaction_id` varchar(100) DEFAULT NULL COMMENT '第三方交易ID',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `refund_time` datetime DEFAULT NULL COMMENT '退款时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_payment_no` (`payment_no`),
  KEY `idx_order_id` (`order_id`),
  CONSTRAINT `fk_payment_order` FOREIGN KEY (`order_id`) REFERENCES `tb_order` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单支付信息表';

-- 订单商品表
CREATE TABLE IF NOT EXISTS `tb_order_product` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '订单商品ID',
  `order_id` int(11) NOT NULL COMMENT '订单ID',
  `product_id` int(11) NOT NULL COMMENT '商品ID',
  `product_name` varchar(255) NOT NULL COMMENT '商品名称',
  `product_price` decimal(10,2) NOT NULL COMMENT '商品价格',
  `quantity` int(11) NOT NULL DEFAULT '1' COMMENT '数量',
  `total_price` decimal(10,2) NOT NULL COMMENT '总价',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_product_id` (`product_id`),
  CONSTRAINT `fk_order_product_order` FOREIGN KEY (`order_id`) REFERENCES `tb_order` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单商品表';

-- 订单日志表
CREATE TABLE IF NOT EXISTS `tb_order_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `order_id` int(11) NOT NULL COMMENT '订单ID',
  `admin_id` int(11) DEFAULT NULL COMMENT '管理员ID',
  `admin_name` varchar(50) NOT NULL COMMENT '管理员名称',
  `action` varchar(50) NOT NULL COMMENT '操作类型：view(查看)、edit(编辑)、delete(删除)、approve(批准)、reject(拒绝)、ship(发货)、complete(完成)、cancel(取消)',
  `status_before` varchar(50) DEFAULT NULL COMMENT '操作前状态',
  `status_after` varchar(50) DEFAULT NULL COMMENT '操作后状态',
  `description` varchar(255) NOT NULL COMMENT '操作描述',
  `ip` varchar(50) NOT NULL COMMENT '操作IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_order_log_order` FOREIGN KEY (`order_id`) REFERENCES `tb_order` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单操作日志表';

-- 订单配送信息表
CREATE TABLE IF NOT EXISTS `tb_order_shipping` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '配送ID',
  `order_id` int(11) NOT NULL COMMENT '订单ID',
  `shipping_method` varchar(50) NOT NULL COMMENT '配送方式',
  `shipping_company` varchar(50) DEFAULT NULL COMMENT '快递公司',
  `tracking_number` varchar(50) DEFAULT NULL COMMENT '快递单号',
  `shipping_fee` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '配送费用',
  `receiver_name` varchar(50) NOT NULL COMMENT '收件人姓名',
  `receiver_phone` varchar(20) NOT NULL COMMENT '收件人电话',
  `receiver_address` varchar(255) NOT NULL COMMENT '收件人地址',
  `ship_time` datetime DEFAULT NULL COMMENT '发货时间',
  `receive_time` datetime DEFAULT NULL COMMENT '收货时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  CONSTRAINT `fk_shipping_order` FOREIGN KEY (`order_id`) REFERENCES `tb_order` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单配送信息表';

-- 示例数据：订单主表
INSERT INTO `tb_order` (`order_no`, `user_id`, `customer`, `amount`, `status`, `payment_method`, `product`, `order_sort`, `created_at`) VALUES 
('ORD20230501001', 1001, '张三', 128.00, 'completed', '微信支付', '置顶一小时', 1, '2023-05-01 10:23:45'),
('ORD20230501002', 1002, '李四', 256.50, 'pending', '支付宝', '商品B', 2, '2023-05-01 11:34:56'),
('ORD20230502001', 1003, '王五', 64.80, 'paid', '微信支付', '商品C', 3, '2023-05-02 09:12:34'),
('ORD20230502002', 1004, '赵六', 199.00, 'shipped', '银行卡', '商品D', 4, '2023-05-02 14:45:23'),
('ORD20230503001', 1005, '钱七', 88.88, 'cancelled', '支付宝', '商品E', 5, '2023-05-03 16:54:32');

-- 示例数据：订单支付信息
INSERT INTO `tb_order_payment` (`order_id`, `payment_no`, `payment_method`, `amount`, `status`, `transaction_id`, `pay_time`) VALUES 
(1, 'PAY20230501001', '微信支付', 128.00, 'success', 'WX123456789', '2023-05-01 10:30:45'),
(3, 'PAY20230502001', '微信支付', 64.80, 'success', 'WX987654321', '2023-05-02 09:20:34');

-- 示例数据：订单商品
INSERT INTO `tb_order_product` (`order_id`, `product_id`, `product_name`, `product_price`, `quantity`, `total_price`) VALUES 
(1, 101, '置顶一小时', 128.00, 1, 128.00),
(2, 102, '商品B', 256.50, 1, 256.50),
(3, 103, '商品C', 64.80, 1, 64.80),
(4, 104, '商品D', 199.00, 1, 199.00),
(5, 105, '商品E', 88.88, 1, 88.88);

-- 示例数据：订单日志
INSERT INTO `tb_order_log` (`order_id`, `admin_id`, `admin_name`, `action`, `status_before`, `status_after`, `description`, `ip`) VALUES 
(1, 1, '管理员', 'approve', 'paid', 'completed', '确认订单完成', '192.168.1.100'),
(3, 1, '管理员', 'approve', 'pending', 'paid', '确认支付成功', '192.168.1.100'),
(4, 1, '管理员', 'ship', 'paid', 'shipped', '订单已发货', '192.168.1.100'),
(5, 1, '管理员', 'cancel', 'pending', 'cancelled', '用户申请取消订单', '192.168.1.100');

-- 示例数据：订单配送信息
INSERT INTO `tb_order_shipping` (`order_id`, `shipping_method`, `shipping_company`, `tracking_number`, `shipping_fee`, `receiver_name`, `receiver_phone`, `receiver_address`, `ship_time`) VALUES 
(4, '快递', '顺丰速运', 'SF12345678', 0.00, '赵六', '13900000004', '北京市海淀区中关村大街1号', '2023-05-02 16:30:00'); 