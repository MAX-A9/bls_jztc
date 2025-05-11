<template>
	<view class="order-detail-container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">订单详情</text>
			</view>
		</view>
		
		<!-- 占位元素 -->
		<view class="placeholder" :style="{ height: `calc(${statusBarHeight}px + ${navBarHeight}px)` }"></view>
		
		<!-- 加载中 -->
		<view class="loading-container" v-if="loading">
			<uni-icons type="spinner-cycle" size="30" color="#999999"></uni-icons>
			<text class="loading-text">加载中...</text>
		</view>
		
		<!-- 订单内容 -->
		<view class="order-content" v-else-if="orderDetail">
			<!-- 订单状态 -->
			<view class="status-section">
				<text class="status-text" :class="getStatusClass(orderDetail.status)">{{ orderDetail.statusText }}</text>
				<text class="status-desc">{{ getStatusDesc(orderDetail.status) }}</text>
			</view>
			
			<!-- 商品信息 -->
			<view class="product-section">
				<view class="product-info">
					<view class="product-icon">
						<uni-icons type="shop" size="30" color="#fc3e2b"></uni-icons>
					</view>
					<view class="product-detail">
						<text class="product-name">{{ orderDetail.productName }}</text>
						<text class="product-price">¥{{ orderDetail.amount.toFixed(2) }}</text>
					</view>
				</view>
			</view>
			
			<!-- 订单信息 -->
			<view class="info-section">
				<view class="info-title">订单信息</view>
				<view class="info-item">
					<text class="info-label">订单编号</text>
					<view class="info-value-copy">
						<text class="info-value">{{ orderDetail.orderNo }}</text>
						<view class="copy-btn" @tap="copyOrderNo">复制</view>
					</view>
				</view>
				<view class="info-item">
					<text class="info-label">下单时间</text>
					<text class="info-value">{{ orderDetail.createdAt }}</text>
				</view>
				<view class="info-item" v-if="orderDetail.payTime">
					<text class="info-label">支付时间</text>
					<text class="info-value">{{ orderDetail.payTime }}</text>
				</view>
				<view class="info-item">
					<text class="info-label">支付方式</text>
					<text class="info-value">{{ getPaymentMethodText(orderDetail.paymentMethod) }}</text>
				</view>
				<view class="info-item" v-if="orderDetail.remark">
					<text class="info-label">备注</text>
					<text class="info-value">{{ orderDetail.remark }}</text>
				</view>
			</view>
			
			<!-- 订单金额 -->
			<view class="price-section">
				<view class="info-title">订单金额</view>
				<view class="price-item">
					<text class="price-label">实付金额</text>
					<text class="price-value">¥{{ orderDetail.amount.toFixed(2) }}</text>
				</view>
			</view>
			
			<!-- 底部操作区 -->
			<view class="action-section" v-if="orderDetail.status === 0">
				<view class="action-btn cancel" @tap="handleCancelOrder">取消订单</view>
				<view class="action-btn primary" @tap="handlePayOrder">立即支付</view>
			</view>
			
			<view class="action-section" v-else-if="orderDetail.status === 1 || orderDetail.status === 4">
				<view class="action-btn" @tap="handleContactService">联系客服</view>
			</view>
			
			<view class="action-section" v-else-if="orderDetail.status === 2 || orderDetail.status === 3 || orderDetail.status === 5">
				<view class="action-btn" @tap="handleRebuyOrder">再次购买</view>
			</view>
		</view>
		
		<!-- 空状态 -->
		<view class="empty-state" v-else>
			<image class="empty-icon" src="/static/images/empty.png" mode="aspectFit"></image>
			<text class="empty-text">订单信息不存在</text>
			<view class="return-btn" @tap="handleBack">返回上一页</view>
		</view>
	</view>
</template>

<script>
	import deviceMixin from '@/mixins/device.js'
	import { user } from '@/apis/index.js'
	import { requestOrderPay } from '@/utils/pay.js'
	
	export default {
		mixins: [deviceMixin],
		data() {
			return {
				orderId: null,
				orderNo: null,
				loading: true,
				orderDetail: null
			}
		},
		onLoad(options) {
			// 获取订单ID和订单号
			if (options.id) {
				this.orderId = options.id
				// 加载订单详情
				this.loadOrderDetail()
			} else if (options.orderNo) {
				this.orderNo = options.orderNo
				// 加载订单详情
				this.loadOrderDetail()
			} else {
				this.loading = false
				uni.showToast({
					title: '订单参数不存在',
					icon: 'none'
				})
			}
		},
		methods: {
			// 返回上一页
			handleBack() {
				// 获取当前页面栈
				const pages = getCurrentPages();
				// 如果当前页面不是第一个页面，才可以返回
				if (pages.length > 1) {
					uni.navigateBack();
				} else {
					// 如果是第一个页面，则跳转到订单列表
					uni.redirectTo({
						url: '/pages/my/orders/index'
					});
				}
			},
			
			// 加载订单详情
			async loadOrderDetail() {
				// 标记加载状态
				let loadingShown = false;
				
				try {
					this.loading = true
					
					// 显示加载提示
					uni.showLoading({
						title: '加载中...',
						mask: true
					});
					loadingShown = true;
					
					// 调用接口获取订单详情
					const params = this.orderNo ? { orderNo: this.orderNo } : { id: this.orderId };
					const res = await user.getOrderDetail(params);
					
					// 请求成功后关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					if (res.code === 0 && res.data) {
						this.orderDetail = res.data;
						this.loading = false;
					} else {
						throw new Error(res.message || '获取订单详情失败');
					}
				} catch (error) {
					// 确保异常情况下也关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					console.error('获取订单详情失败:', error)
					uni.showToast({
						title: '获取订单失败',
						icon: 'none'
					})
					this.loading = false
				} finally {
					// 最后确保关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
					}
				}
			},
			
			// 根据订单状态获取CSS类名
			getStatusClass(status) {
				switch (status) {
					case 0: return 'unpaid'     // 待支付
					case 1: return 'paid'       // 已支付
					case 2: return 'cancelled'  // 已取消
					case 3: return 'refunded'   // 已退款
					case 4: return 'processing' // 进行中
					case 5: return 'completed'  // 已完成
					default: return ''
				}
			},
			
			// 根据订单状态获取描述文本
			getStatusDesc(status) {
				switch (status) {
					case 0: return '请及时完成支付'
					case 1: return '您的订单已支付'
					case 2: return '您的订单已取消'
					case 3: return '您的订单已退款'
					case 4: return '您的订单正在进行中'
					case 5: return '您的订单已完成'
					default: return ''
				}
			},
			
			// 获取支付方式文本
			getPaymentMethodText(method) {
				switch (method) {
					case 'wechat': return '微信支付'
					case 'alipay': return '支付宝'
					default: return '未知'
				}
			},
			
			// 复制订单号
			copyOrderNo() {
				if (!this.orderDetail || !this.orderDetail.orderNo) return;
				
				uni.setClipboardData({
					data: this.orderDetail.orderNo,
					success: () => {
						uni.showToast({
							title: '复制成功',
							icon: 'success'
						})
					}
				})
			},
			
			// 取消订单
			handleCancelOrder() {
				uni.showModal({
					title: '提示',
					content: '确定要取消该订单吗？',
					success: async (res) => {
						if (res.confirm) {
							// 标记加载状态
							let loadingShown = false;
							
							try {
								// 显示加载提示
								uni.showLoading({
									title: '处理中...',
									mask: true
								});
								loadingShown = true;
								
								// 调用取消订单接口
								await user.cancelOrder({ 
									orderNo: this.orderDetail.orderNo,
									reason: '用户主动取消'
								});
								
								// 请求成功后关闭加载提示
								if (loadingShown) {
									uni.hideLoading();
									loadingShown = false;
								}
								
								uni.showToast({
									title: '已取消订单',
									icon: 'success'
								})
								
								// 重新加载订单详情
								this.loadOrderDetail()
							} catch (error) {
								// 确保异常情况下也关闭加载提示
								if (loadingShown) {
									uni.hideLoading();
									loadingShown = false;
								}
								
								console.error('取消订单失败:', error)
								uni.showToast({
									title: error.message || '取消失败',
									icon: 'none'
								})
							} finally {
								// 最后确保关闭加载提示
								if (loadingShown) {
									uni.hideLoading();
								}
							}
						}
					}
				})
			},
			
			// 支付订单
			async handlePayOrder() {
				// 标记加载状态
				let loadingShown = false;
				
				try {
					// 显示加载提示
					uni.showLoading({
						title: '正在支付...',
						mask: true
					});
					loadingShown = true;
					
					// 调用订单支付接口
					const payResult = await requestOrderPay({
						orderNo: this.orderDetail.orderNo
					})
					
					// 请求成功后关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					if (payResult.success) {
						uni.showToast({
							title: '支付成功',
							icon: 'success'
						})
						
						// 重新加载订单详情
						this.loadOrderDetail()
					} else {
						uni.showToast({
							title: payResult.message || '支付已取消',
							icon: 'none'
						})
					}
				} catch (error) {
					// 确保异常情况下也关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					console.error('支付过程出错:', error)
					uni.showToast({
						title: error.message || '支付失败',
						icon: 'none'
					})
				} finally {
					// 最后确保关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
					}
				}
			},
			
			// 联系客服
			handleContactService() {
				// 实现联系客服的逻辑
				console.log('联系客服')
			},
			
			// 再次购买
			handleRebuyOrder() {
				// 实现再次购买的逻辑
				console.log('再次购买')
			}
		}
	}
</script>

<style>
	.order-detail-container {
		min-height: 100vh;
		background-color: #f5f5f5;
		position: relative;
	}
	
	.nav-bar {
		background-color: #ffffff;
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 99;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
	}
	
	.nav-content {
		display: flex;
		align-items: center;
		position: relative;
	}
	
	.back-btn {
		width: 88rpx;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-left: 7px;
	}
	
	.nav-title {
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
		font-size: 34rpx;
		font-weight: 500;
		color: #333333;
		max-width: 350rpx;
		text-align: center;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.placeholder {
		width: 100%;
	}
	
	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 500rpx;
	}
	
	.loading-text {
		font-size: 28rpx;
		color: #999999;
		margin-top: 20rpx;
	}
	
	.order-content {
		padding: 30rpx;
	}
	
	.status-section {
		background-color: #ffffff;
		border-radius: 16rpx;
		padding: 40rpx 30rpx;
		margin-bottom: 20rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
	.status-text {
		font-size: 36rpx;
		font-weight: 500;
		margin-bottom: 16rpx;
	}
	
	.status-text.unpaid {
		color: #fc3e2b;
	}
	
	.status-text.paid {
		color: #67c23a;
	}
	
	.status-text.cancelled {
		color: #909399;
	}
	
	.status-text.refunded {
		color: #e6a23c;
	}
	
	.status-text.processing {
		color: #409eff;
	}
	
	.status-text.completed {
		color: #67c23a;
	}
	
	.status-desc {
		font-size: 28rpx;
		color: #666666;
	}
	
	.product-section {
		background-color: #ffffff;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-bottom: 20rpx;
	}
	
	.product-info {
		display: flex;
		align-items: center;
	}
	
	.product-icon {
		width: 120rpx;
		height: 120rpx;
		border-radius: 8rpx;
		margin-right: 20rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: #f8f8f8;
	}
	
	.product-detail {
		flex: 1;
		display: flex;
		flex-direction: column;
	}
	
	.product-name {
		font-size: 30rpx;
		color: #333333;
		font-weight: 500;
		margin-bottom: 16rpx;
	}
	
	.product-price {
		font-size: 32rpx;
		color: #fc3e2b;
		font-weight: 500;
	}
	
	.info-section, .price-section {
		background-color: #ffffff;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-bottom: 20rpx;
	}
	
	.info-title {
		font-size: 32rpx;
		font-weight: 500;
		color: #333333;
		margin-bottom: 24rpx;
	}
	
	.info-item {
		display: flex;
		justify-content: space-between;
		margin-bottom: 16rpx;
	}
	
	.info-item:last-child {
		margin-bottom: 0;
	}
	
	.info-label {
		font-size: 28rpx;
		color: #666666;
	}
	
	.info-value {
		font-size: 28rpx;
		color: #333333;
	}
	
	.info-value-copy {
		display: flex;
		align-items: center;
	}
	
	.copy-btn {
		font-size: 24rpx;
		color: #1e90ff;
		margin-left: 16rpx;
	}
	
	.price-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.price-label {
		font-size: 28rpx;
		color: #666666;
	}
	
	.price-value {
		font-size: 34rpx;
		color: #fc3e2b;
		font-weight: 500;
	}
	
	.action-section {
		display: flex;
		justify-content: flex-end;
		gap: 20rpx;
		margin-top: 40rpx;
	}
	
	.action-btn {
		padding: 16rpx 40rpx;
		border-radius: 40rpx;
		font-size: 28rpx;
	}
	
	.action-btn.cancel {
		border: 1rpx solid #dddddd;
		color: #666666;
	}
	
	.action-btn.primary {
		background-color: #fc3e2b;
		color: #ffffff;
	}
	
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding-top: 120rpx;
	}
	
	.empty-icon {
		width: 200rpx;
		height: 200rpx;
		margin-bottom: 30rpx;
	}
	
	.empty-text {
		font-size: 30rpx;
		color: #999999;
		margin-bottom: 40rpx;
	}
	
	.return-btn {
		padding: 16rpx 40rpx;
		border-radius: 40rpx;
		background-color: #fc3e2b;
		color: #ffffff;
		font-size: 28rpx;
	}
</style> 