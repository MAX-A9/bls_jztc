# 数据库表结构说明文档

本文档包含系统所需的全部数据表SQL命令说明，包括管理员表、用户表、内容管理表、仪表盘表、订单表、兑换管理表、文件管理表、常规管理表和同城系统表的结构说明。

## 文件说明

- `admin_tables.sql` - 管理员相关表创建SQL脚本
- `user_tables.sql` - 用户相关表创建SQL脚本
- `content_tables.sql` - 内容管理相关表创建SQL脚本
- `dashboard_tables.sql` - 仪表盘相关表创建SQL脚本
- `order_tables.sql` - 订单管理相关表创建SQL脚本
- `exchange_tables.sql` - 兑换管理相关表创建SQL脚本
- `file_tables.sql` - 文件管理相关表创建SQL脚本
- `settings_tables.sql` - 常规管理相关表创建SQL脚本
- `city_system_tables.sql` - 同城系统相关表创建SQL脚本

## 管理员相关表结构

### 管理员表 (tb_admin)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 管理员ID (主键) |
| username | varchar(50) | 用户名 (唯一) |
| password | varchar(255) | 密码 (加密存储) |
| real_name | varchar(50) | 真实姓名 |
| phone | varchar(20) | 手机号 |
| role_group | varchar(50) | 角色组 |
| status | tinyint(1) | 状态：1正常，0禁用 |
| last_login_time | datetime | 最后登录时间 |
| last_login_ip | varchar(50) | 最后登录IP |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 管理员权限表 (tb_admin_permission)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 权限ID (主键) |
| admin_id | int | 管理员ID (外键) |
| permission_name | varchar(50) | 权限名称 |
| created_at | datetime | 创建时间 |

### 管理员操作日志表 (tb_admin_log)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 日志ID (主键) |
| admin_id | int | 管理员ID (外键) |
| admin_name | varchar(50) | 管理员用户名 |
| ip | varchar(50) | 操作IP |
| action | varchar(50) | 操作类型：登录、登出、添加、修改、删除、导出、导入、查询 |
| module | varchar(50) | 操作模块 |
| description | varchar(255) | 操作描述 |
| status | tinyint(1) | 操作状态：1成功，0失败 |
| details | text | 详细信息 |
| created_at | datetime | 操作时间 |

## 用户相关表结构

### 用户表 (tb_user)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 用户ID (主键) |
| username | varchar(50) | 用户名 (唯一) |
| password | varchar(255) | 密码 (加密存储) |
| real_name | varchar(50) | 真实姓名 |
| phone | varchar(20) | 手机号 (唯一) |
| role | varchar(20) | 角色：管理员、编辑员、普通用户 |
| status | tinyint(1) | 状态：1正常，0禁用 |
| identifier | varchar(20) | 用户来源标识：小程序、未知 |
| avatar | varchar(255) | 头像URL |
| department | varchar(50) | 部门 |
| position | varchar(50) | 职位 |
| last_login_time | datetime | 最后登录时间 |
| last_login_ip | varchar(50) | 最后登录IP |
| sms_notifications | tinyint(1) | 是否开启短信通知：1开启，0关闭 |
| theme | varchar(20) | 界面主题偏好 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 用户扩展信息表 (tb_user_profile)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | ID (主键) |
| user_id | int | 用户ID (外键) |
| email | varchar(100) | 邮箱 |
| gender | tinyint(1) | 性别：1男，2女，0未知 |
| birth_date | date | 出生日期 |
| address | varchar(255) | 地址 |
| id_card | varchar(20) | 身份证号 |
| bio | text | 个人简介 |
| tags | varchar(255) | 标签，逗号分隔 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 用户登录日志表 (tb_user_login_log)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 日志ID (主键) |
| user_id | int | 用户ID (外键) |
| username | varchar(50) | 用户名 |
| login_ip | varchar(50) | 登录IP |
| login_time | datetime | 登录时间 |
| login_type | varchar(20) | 登录类型：web、app、小程序 |
| login_status | tinyint(1) | 登录状态：1成功，0失败 |
| device_info | varchar(255) | 设备信息 |
| user_agent | varchar(500) | 浏览器UA信息 |
| remark | varchar(255) | 备注信息 |

### 用户操作日志表 (tb_user_operation_log)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 日志ID (主键) |
| user_id | int | 用户ID (外键) |
| username | varchar(50) | 用户名 |
| operation_ip | varchar(50) | 操作IP |
| operation_time | datetime | 操作时间 |
| module | varchar(50) | 操作模块 |
| operation | varchar(50) | 操作类型 |
| description | varchar(255) | 操作描述 |
| operation_result | tinyint(1) | 操作结果：1成功，0失败 |
| details | text | 操作详情 |

## 内容管理相关表结构

### 内容表 (tb_content)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 内容ID (主键) |
| title | varchar(255) | 标题 |
| category | varchar(50) | 分类 |
| author | varchar(50) | 作者 |
| content | text | 内容详情(富文本) |
| status | varchar(20) | 状态：已发布、待审核、已下架 |
| views | int | 浏览量 |
| likes | int | 想要数量 |
| comments | int | 评论数 |
| is_recommended | tinyint(1) | 是否置顶推荐：1是，0否 |
| published_at | datetime | 发布时间 |
| expires_at | datetime | 到期时间 |
| top_until | datetime | 置顶截止时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 首页分类表 (tb_home_category)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 分类ID (主键) |
| name | varchar(50) | 分类名称 |
| sort_order | int | 排序顺序 |
| is_active | tinyint(1) | 是否启用：1启用，0禁用 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 闲置分类表 (tb_idle_category)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 分类ID (主键) |
| name | varchar(50) | 分类名称 |
| sort_order | int | 排序顺序 |
| is_active | tinyint(1) | 是否启用：1启用，0禁用 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 内容评论表 (tb_content_comment)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 评论ID (主键) |
| content_id | int | 内容ID (外键) |
| user_id | int | 用户ID |
| user_name | varchar(50) | 用户名 |
| comment | text | 评论内容 |
| status | varchar(20) | 状态：已审核、待审核、已拒绝 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 内容附件表 (tb_content_attachment)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 附件ID (主键) |
| content_id | int | 内容ID (外键) |
| file_name | varchar(255) | 文件名 |
| file_path | varchar(255) | 文件路径 |
| file_size | int | 文件大小(字节) |
| file_type | varchar(50) | 文件类型 |
| created_at | datetime | 创建时间 |

### 内容操作日志表 (tb_content_log)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 日志ID (主键) |
| content_id | int | 内容ID (外键) |
| admin_id | int | 管理员ID |
| admin_name | varchar(50) | 管理员名称 |
| action | varchar(50) | 操作类型：创建、编辑、删除、审核、推荐 |
| status_change | varchar(50) | 状态变更：变更前状态->变更后状态 |
| details | text | 操作详情 |
| ip | varchar(50) | 操作IP |
| created_at | datetime | 操作时间 |

## 仪表盘相关表结构

### 统计数据表 (tb_statistics)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 统计ID (主键) |
| date | date | 统计日期 |
| registration_count | int | 注册用户数 |
| exchange_count | int | 兑换次数 |
| publish_count | int | 发布内容数 |
| income_amount | decimal(10,2) | 收益金额 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 趋势数据表 (tb_trend_data)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 趋势ID (主键) |
| date | date | 日期 |
| category | varchar(20) | 类别：registration(注册用户)、exchange(兑换次数)、publish(发布内容)、income(收益金额) |
| value | int | 数值 |
| change_percent | decimal(5,2) | 变化百分比 |
| created_at | datetime | 创建时间 |

### 订单表 (tb_order)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 订单ID (主键) |
| order_no | varchar(30) | 订单编号 (唯一) |
| user_id | int | 用户ID |
| customer | varchar(50) | 客户名称 |
| amount | decimal(10,2) | 订单金额 |
| status | enum | 订单状态：pending(待处理)、paid(已支付)、shipped(已发货)、completed(已完成)、cancelled(已取消) |
| payment_method | varchar(20) | 支付方式 |
| product | varchar(255) | 商品名称 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 兑换表 (tb_exchange)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 兑换ID (主键) |
| exchange_no | varchar(30) | 兑换编号 (唯一) |
| user_id | int | 用户ID |
| user_name | varchar(50) | 用户名称 |
| product_id | int | 商品ID |
| product_name | varchar(255) | 商品名称 |
| points | int | 消耗积分 |
| status | enum | 兑换状态：pending(待处理)、processing(处理中)、completed(已完成)、cancelled(已取消) |
| remark | varchar(255) | 备注 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 仪表盘日志表 (tb_dashboard_log)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 日志ID (主键) |
| admin_id | int | 管理员ID |
| admin_name | varchar(50) | 管理员名称 |
| action | varchar(50) | 操作类型：view(查看)、export(导出)、filter(筛选) |
| time_range | varchar(20) | 时间范围：本周、本月、本年 |
| content_type | varchar(50) | 内容类型 |
| ip | varchar(50) | 操作IP |
| created_at | datetime | 操作时间 |

## 订单管理相关表结构

### 订单主表 (tb_order)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 订单ID (主键) |
| order_no | varchar(30) | 订单编号 (唯一) |
| user_id | int | 用户ID |
| customer | varchar(50) | 客户名称 |
| amount | decimal(10,2) | 订单金额 |
| status | enum | 订单状态:pending(待付款)、paid(已付款)、shipped(已发货)、completed(已完成)、cancelled(已取消) |
| payment_method | varchar(20) | 支付方式 |
| product_id | int | 商品ID |
| product | varchar(255) | 商品名称 |
| order_sort | int | 排序字段 |
| remark | varchar(255) | 备注 |
| pay_time | datetime | 支付时间 |
| ship_time | datetime | 发货时间 |
| complete_time | datetime | 完成时间 |
| cancel_time | datetime | 取消时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 订单支付信息表 (tb_order_payment)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 支付ID (主键) |
| order_id | int | 订单ID (外键) |
| payment_no | varchar(50) | 支付流水号 |
| payment_method | varchar(20) | 支付方式 |
| amount | decimal(10,2) | 支付金额 |
| status | enum | 支付状态:pending(待支付)、success(成功)、failed(失败)、refunded(已退款) |
| transaction_id | varchar(100) | 第三方交易ID |
| pay_time | datetime | 支付时间 |
| refund_time | datetime | 退款时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 订单商品表 (tb_order_product)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 订单商品ID (主键) |
| order_id | int | 订单ID (外键) |
| product_id | int | 商品ID |
| product_name | varchar(255) | 商品名称 |
| product_price | decimal(10,2) | 商品价格 |
| quantity | int | 数量 |
| total_price | decimal(10,2) | 总价 |
| created_at | datetime | 创建时间 |

### 订单日志表 (tb_order_log)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 日志ID (主键) |
| order_id | int | 订单ID (外键) |
| admin_id | int | 管理员ID |
| admin_name | varchar(50) | 管理员名称 |
| action | varchar(50) | 操作类型:view(查看)、edit(编辑)、delete(删除)、approve(批准)、reject(拒绝)、ship(发货)、complete(完成)、cancel(取消) |
| status_before | varchar(50) | 操作前状态 |
| status_after | varchar(50) | 操作后状态 |
| description | varchar(255) | 操作描述 |
| ip | varchar(50) | 操作IP |
| created_at | datetime | 操作时间 |

### 订单配送信息表 (tb_order_shipping)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 配送ID (主键) |
| order_id | int | 订单ID (外键) |
| shipping_method | varchar(50) | 配送方式 |
| shipping_company | varchar(50) | 快递公司 |
| tracking_number | varchar(50) | 快递单号 |
| shipping_fee | decimal(10,2) | 配送费用 |
| receiver_name | varchar(50) | 收件人姓名 |
| receiver_phone | varchar(20) | 收件人电话 |
| receiver_address | varchar(255) | 收件人地址 |
| ship_time | datetime | 发货时间 |
| receive_time | datetime | 收货时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 兑换管理相关表结构

### 兑换商品表 (tb_exchange_product)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 商品ID (主键) |
| name | varchar(100) | 商品名称 |
| description | text | 商品描述 |
| image | varchar(255) | 商品图片URL |
| price | decimal(10,2) | 商品价格 |
| stock | int | 库存数量 |
| category_id | int | 商品分类ID |
| status | enum | 商品状态：active(已上架)、inactive(未上架)、soldout(已售罄) |
| time_required | int | 兑换所需时长(小时) |
| order_sort | int | 排序字段 |
| exchange_count | int | 已兑换次数 |
| is_virtual | tinyint(1) | 是否虚拟商品：1是，0否 |
| is_coupon | tinyint(1) | 是否优惠券：1是，0否 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 商城分类表 (tb_exchange_category)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 分类ID (主键) |
| name | varchar(50) | 分类名称 |
| description | varchar(255) | 分类描述 |
| is_enabled | tinyint(1) | 是否启用：1启用，0禁用 |
| order_sort | int | 排序字段 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 兑换记录表 (tb_exchange_record)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 兑换记录ID (主键) |
| exchange_no | varchar(50) | 兑换单号 |
| user_id | int | 用户ID |
| user_name | varchar(50) | 用户名称 |
| product_id | int | 商品ID |
| product_name | varchar(100) | 商品名称 |
| quantity | int | 兑换数量 |
| phone | varchar(20) | 充值账号/手机号 |
| address | text | 收货地址 |
| status | enum | 兑换状态：pending(处理中)、processing(处理中)、completed(已完成)、failed(失败) |
| remark | varchar(255) | 备注 |
| admin_id | int | 处理管理员ID |
| admin_name | varchar(50) | 处理管理员姓名 |
| process_time | datetime | 处理时间 |
| exchange_time | datetime | 兑换时间 |
| complete_time | datetime | 完成时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 兑换日志表 (tb_exchange_log)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 日志ID (主键) |
| exchange_id | int | 兑换记录ID (外键) |
| admin_id | int | 管理员ID |
| admin_name | varchar(50) | 管理员名称 |
| action | varchar(50) | 操作类型：view(查看)、process(处理)、complete(完成)、fail(标记失败) |
| status_before | varchar(50) | 操作前状态 |
| status_after | varchar(50) | 操作后状态 |
| description | varchar(255) | 操作描述 |
| ip | varchar(50) | 操作IP |
| created_at | datetime | 操作时间 |

### 优惠券表 (tb_coupon)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 优惠券ID (主键) |
| code | varchar(50) | 优惠券码 |
| name | varchar(100) | 优惠券名称 |
| type | enum | 优惠券类型：fixed(固定金额)、percent(百分比) |
| value | decimal(10,2) | 优惠券值：固定金额或百分比 |
| min_order_amount | decimal(10,2) | 最低订单金额 |
| valid_from | datetime | 有效期开始 |
| valid_to | datetime | 有效期结束 |
| user_id | int | 用户ID |
| exchange_id | int | 关联的兑换记录ID |
| status | enum | 状态：unused(未使用)、used(已使用)、expired(已过期) |
| used_time | datetime | 使用时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 文件管理相关表结构

### 文件表 (tb_file)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 文件ID (主键) |
| file_name | varchar(255) | 文件名 |
| original_name | varchar(255) | 原始文件名 |
| file_type | varchar(50) | 文件类型 |
| file_size | bigint | 文件大小(字节) |
| file_path | varchar(255) | 文件路径 |
| thumbnail_path | varchar(255) | 缩略图路径(图片类型) |
| category | enum | 文件分类：image(图片)、document(文档)、media(媒体)、other(其他) |
| description | text | 文件描述 |
| user_id | int | 上传用户ID |
| uploader | varchar(50) | 上传者名称 |
| status | enum | 文件状态：normal(正常)、processing(处理中)、deleted(已删除) |
| download_count | int | 下载次数 |
| view_count | int | 浏览次数 |
| is_public | tinyint(1) | 是否公开：1是，0否 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 文件分享表 (tb_file_share)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 分享ID (主键) |
| file_id | int | 文件ID |
| user_id | int | 分享用户ID |
| share_code | varchar(32) | 分享码 |
| expire_time | datetime | 过期时间 |
| view_count | int | 查看次数 |
| download_count | int | 下载次数 |
| is_password | tinyint(1) | 是否设置密码：1是，0否 |
| password | varchar(32) | 访问密码 |
| status | tinyint(1) | 状态：1有效，0无效 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 文件标签表 (tb_file_tag)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 标签ID (主键) |
| tag_name | varchar(50) | 标签名称 |
| created_at | datetime | 创建时间 |

### 文件标签关联表 (tb_file_tag_relation)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 关联ID (主键) |
| file_id | int | 文件ID |
| tag_id | int | 标签ID |
| created_at | datetime | 创建时间 |

### 文件操作日志表 (tb_file_log)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 日志ID (主键) |
| file_id | int | 文件ID |
| user_id | int | 操作用户ID |
| username | varchar(50) | 操作用户名 |
| action | enum | 操作类型：upload(上传)、download(下载)、view(查看)、delete(删除)、rename(重命名)、share(分享)、update(更新) |
| description | varchar(255) | 操作描述 |
| ip | varchar(50) | 操作IP |
| created_at | datetime | 操作时间 |

### 文件夹表 (tb_folder)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 文件夹ID (主键) |
| folder_name | varchar(100) | 文件夹名称 |
| parent_id | int | 父文件夹ID |
| user_id | int | 创建用户ID |
| description | varchar(255) | 文件夹描述 |
| is_public | tinyint(1) | 是否公开：1是，0否 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 文件与文件夹关联表 (tb_file_folder)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 关联ID (主键) |
| file_id | int | 文件ID |
| folder_id | int | 文件夹ID |
| created_at | datetime | 创建时间 |

## 常规管理相关表结构

### 系统配置表 (tb_system_config)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 配置ID (主键) |
| config_key | varchar(50) | 配置键名 (唯一) |
| config_value | text | 配置值 |
| config_group | varchar(50) | 配置分组：system(系统)、file(文件)、locale(本地化)、storage(存储) |
| remark | varchar(255) | 备注说明 |
| is_system | tinyint(1) | 是否系统配置：1是，0否 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 用户偏好设置表 (tb_user_preference)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 设置ID (主键) |
| user_id | int | 用户ID |
| preference_key | varchar(50) | 偏好设置键名 |
| preference_value | text | 偏好设置值 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 短信渠道表 (tb_sms_channel)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 渠道ID (主键) |
| name | varchar(50) | 渠道名称 |
| app_id | varchar(100) | 应用ID |
| app_key | varchar(255) | 应用密钥 |
| sign_name | varchar(50) | 短信签名 |
| status | tinyint(1) | 状态：1启用，0禁用 |
| is_default | tinyint(1) | 是否默认：1是，0否 |
| config | text | 额外配置JSON |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 短信模板表 (tb_sms_template)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 模板ID (主键) |
| name | varchar(50) | 模板名称 |
| code | varchar(100) | 模板编码 |
| content | text | 模板内容 |
| channel_id | int | 所属渠道ID |
| type | varchar(20) | 模板类型：验证码、通知、营销 |
| status | tinyint(1) | 状态：1启用，0禁用 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 短信发送记录表 (tb_sms_record)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 记录ID (主键) |
| mobile | varchar(20) | 手机号 |
| content | text | 短信内容 |
| template_id | int | 模板ID |
| channel_id | int | 渠道ID |
| status | tinyint(1) | 发送状态：1成功，0失败 |
| error_msg | varchar(255) | 错误消息 |
| send_time | datetime | 发送时间 |
| created_at | datetime | 创建时间 |

### 短信验证码表 (tb_sms_code)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | ID (主键) |
| mobile | varchar(20) | 手机号 |
| code | varchar(10) | 验证码 |
| type | varchar(20) | 验证码类型：登录、注册、找回密码 |
| expire_time | datetime | 过期时间 |
| is_used | tinyint(1) | 是否已使用：1是，0否 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 个人设置表 (tb_user_profile_setting)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 设置ID (主键) |
| user_id | int | 用户ID |
| avatar | varchar(255) | 头像URL |
| position | varchar(50) | 职位 |
| department | varchar(50) | 部门 |
| theme | varchar(20) | 界面主题：light(浅色)、dark(深色) |
| sms_notifications | tinyint(1) | 是否开启短信通知：1开启，0关闭 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 同城系统相关表结构

### Banner表 (tb_banner)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | Banner ID (主键) |
| image | varchar(255) | 图片路径 |
| link_type | varchar(20) | 链接类型：page(页面)、miniprogram(小程序)、url(外部链接) |
| link_url | varchar(255) | 链接地址 |
| is_enabled | tinyint(1) | 是否启用：1是，0否 |
| order_sort | int | 排序字段 |
| position | varchar(20) | 位置：home(首页)、category(分类页)等 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 导航小程序表 (tb_miniprogram_nav)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 导航ID (主键) |
| name | varchar(50) | 小程序名称 |
| app_id | varchar(50) | 小程序AppID |
| logo | varchar(255) | Logo图片路径 |
| is_enabled | tinyint(1) | 是否启用：1是，0否 |
| order_sort | int | 排序字段 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 内页设置表 (tb_page_settings)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 设置ID (主键) |
| page_key | varchar(50) | 页面键名 |
| setting_key | varchar(50) | 设置键名 |
| setting_value | text | 设置值 |
| remark | varchar(255) | 备注说明 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 底部Tab设置表 (tb_tab_bar)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | Tab ID (主键) |
| title | varchar(20) | Tab标题 |
| icon | varchar(255) | 图标路径 |
| is_active | tinyint(1) | 是否启用：1是，0否 |
| order_sort | int | 排序字段 |
| page_path | varchar(100) | 页面路径 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 小程序配置表 (tb_miniprogram_config)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 配置ID (主键) |
| app_id | varchar(50) | 小程序AppID |
| app_secret | varchar(100) | 小程序AppSecret |
| name | varchar(50) | 小程序名称 |
| description | text | 小程序描述 |
| logo | varchar(255) | Logo图片路径 |
| is_enabled | tinyint(1) | 是否启用：1是，0否 |
| privacy_policy | text | 隐私政策 |
| user_agreement | text | 用户协议 |
| mch_id | varchar(50) | 商户号 |
| mch_key | varchar(100) | 商户密钥 |
| pay_notify_url | varchar(255) | 支付回调URL |
| enable_payment | tinyint(1) | 是否启用支付：1是，0否 |
| enable_ads | tinyint(1) | 是否启用广告：1是，0否 |
| ad_app_id | varchar(50) | 广告AppID |
| banner_ad_unit_id | varchar(50) | Banner广告位ID |
| interstitial_ad_unit_id | varchar(50) | 插屏广告位ID |
| video_ad_unit_id | varchar(50) | 视频广告位ID |
| native_ad_unit_id | varchar(50) | 原生广告位ID |
| ad_frequency | varchar(20) | 广告频率：low(低)、medium(中)、high(高)、custom(自定义) |
| custom_ad_frequency | int | 自定义广告频率 |
| ad_positions | varchar(255) | 广告位置，英文逗号分隔 |
| enable_rewards | tinyint(1) | 是否启用奖励：1是，0否 |
| reward_min_minutes | int | 最小奖励分钟数 |
| reward_max_days | int | 最大奖励天数 |
| daily_reward_limit | int | 每日奖励上限 |
| max_daily_reward_time | int | 每日最大奖励小时数 |
| reward_expire_days | int | 奖励过期天数 |
| first_time_min_minutes | int | 首次使用最小分钟数 |
| first_time_max_value | int | 首次使用最大奖励值 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 置顶套餐表 (tb_top_package)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 套餐ID (主键) |
| title | varchar(50) | 套餐名称 |
| description | varchar(255) | 套餐描述 |
| price | decimal(10,2) | 套餐价格 |
| duration | int | 有效时长(小时) |
| is_enabled | tinyint(1) | 是否启用：1是，0否 |
| order_sort | int | 排序字段 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 发布套餐表 (tb_publish_package)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 套餐ID (主键) |
| title | varchar(50) | 套餐名称 |
| description | varchar(255) | 套餐描述 |
| price | decimal(10,2) | 套餐价格 |
| duration | int | 有效时长(天) |
| is_enabled | tinyint(1) | 是否启用：1是，0否 |
| order_sort | int | 排序字段 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## 使用说明

### 手动执行SQL文件

1. 登录到您的MySQL客户端或管理工具
2. 选择您的数据库
3. 执行SQL文件中的SQL命令

```sql
-- 创建管理员相关表
SOURCE /path/to/admin_tables.sql;

-- 创建用户相关表
SOURCE /path/to/user_tables.sql;

-- 创建内容管理相关表
SOURCE /path/to/content_tables.sql;

-- 创建仪表盘相关表
SOURCE /path/to/dashboard_tables.sql;

-- 创建订单管理相关表
SOURCE /path/to/order_tables.sql;

-- 创建兑换管理相关表
SOURCE /path/to/exchange_tables.sql;

-- 创建文件管理相关表
SOURCE /path/to/file_tables.sql;

-- 创建常规管理相关表
SOURCE /path/to/settings_tables.sql;

-- 创建同城系统相关表
SOURCE /path/to/city_system_tables.sql;
```

或者使用mysql命令行工具：

```bash
# 创建管理员相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < admin_tables.sql

# 创建用户相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < user_tables.sql

# 创建内容管理相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < content_tables.sql

# 创建仪表盘相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < dashboard_tables.sql

# 创建订单管理相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < order_tables.sql

# 创建兑换管理相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < exchange_tables.sql

# 创建文件管理相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < file_tables.sql

# 创建常规管理相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < settings_tables.sql

# 创建同城系统相关表
mysql -h主机名 -P端口 -u用户名 -p密码 数据库名 < city_system_tables.sql
```

## 初始账号说明

### 管理员初始账号

SQL文件会自动创建一个超级管理员账号：

- 用户名：admin
- 密码：123456 (加密存储在数据库中)
- 角色组：超级管理员

### 用户初始账号

SQL文件会自动创建两个测试用户：

1. 普通用户：
   - 用户名：testuser
   - 密码：123456 (加密存储在数据库中)
   - 角色：普通用户

2. 编辑员用户：
   - 用户名：editor
   - 密码：123456 (加密存储在数据库中)
   - 角色：编辑员

**注意：** 请在生产环境中删除或修改这些测试账号。 