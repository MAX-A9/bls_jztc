<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">我的订单</text>
			</view>
			
			<!-- 订单状态选项卡 -->
			<view class="tab-bar">
				<view 
					class="tab-item" 
					v-for="(item, index) in tabList" 
					:key="index"
					:class="{ active: currentTab === index }"
					@tap="switchTab(index)"
				>
					<text class="tab-text">{{ item.name }}</text>
				</view>
			</view>
		</view>
		
		<!-- 订单列表 -->
		<scroll-view 
			class="order-scroll" 
			scroll-y 
			:style="{ top: `calc(${statusBarHeight}px + ${navBarHeight}px + 80rpx)` }"
			:refresher-enabled="true"
			:refresher-triggered="isRefreshing"
			@refresherrefresh="onRefresh"
			@scrolltolower="loadMore"
		>
			<view class="order-list">
				<view 
					class="order-item"
					v-for="item in orderList"
					:key="item.id"
					@tap="handleOrderClick(item)"
				>
					<view class="order-header">
						<text class="order-time">{{ item.createdAt }}</text>
						<text class="order-status" :class="getStatusClass(item.status)">{{ item.statusText }}</text>
					</view>
					<view class="order-content">
						<view class="order-icon">
							<uni-icons type="shop" size="30" color="#fc3e2b"></uni-icons>
						</view>
						<view class="order-info">
							<text class="order-title">{{ item.productName }}</text>
							<text class="order-no">订单号：{{ item.orderNo }}</text>
							<view class="price-row">
								<text class="price-label">实付款</text>
								<text class="price-value">¥{{ item.amount.toFixed(2) }}</text>
							</view>
						</view>
					</view>
					<view class="order-footer">
						<view class="btn-group">
							<block v-if="item.status === 0">
								<view class="order-btn" @tap.stop="handleCancelOrder(item)">
									<text>取消订单</text>
								</view>
								<view class="order-btn primary" @tap.stop="handlePayOrder(item)">
									<text>立即支付</text>
								</view>
							</block>
							<block v-else-if="item.status === 1">
								<view class="order-btn" @tap.stop="handleContactService(item)">
									<text>联系客服</text>
								</view>
							</block>
							<block v-else-if="item.status === 2 || item.status === 3">
								<view class="order-btn primary" @tap.stop="handleRebuyOrder(item)">
									<text>再次购买</text>
								</view>
							</block>
						</view>
					</view>
				</view>
			</view>
			
			<!-- 加载更多 -->
			<view class="loading-more" v-if="orderList.length > 0 && hasMore">
				<text class="loading-text">加载中...</text>
			</view>
			
			<!-- 没有更多数据 -->
			<view class="no-more" v-if="orderList.length > 0 && !hasMore">
				<text class="no-more-text">没有更多数据了</text>
			</view>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="orderList.length === 0 && !isLoading">
				<image class="empty-image" src="/static/images/empty.png" mode="aspectFit"></image>
				<text class="empty-text">暂无相关订单</text>
			</view>
		</scroll-view>
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
				currentTab: 0,
				isRefreshing: false,
				isLoading: false,
				hasMore: true,
				page: 1,
				pageSize: 10,
				total: 0,
				tabList: [
					{ name: '全部', status: 'all' },
					{ name: '进行中', status: 'process' },
					{ name: '待支付', status: 'unpaid' },
					{ name: '已完成', status: 'completed' }
				],
				orderList: []
			}
		},
		onLoad(options) {
			if (options.tab) {
				this.currentTab = parseInt(options.tab)
			}
			
			// 获取订单列表
			this.loadOrderList()
		},
		methods: {
			// 获取订单列表
			async loadOrderList(isRefresh = false) {
				if (this.isLoading) return
				
				// 标记加载状态
				let loadingShown = false;
				
				try {
					this.isLoading = true
					
					// 如果不是下拉刷新，则显示加载提示
					if (!isRefresh) {
						uni.showLoading({
							title: '加载中...',
							mask: true
						});
						loadingShown = true;
					}
					
					// 如果是刷新，重置页码
					if (isRefresh) {
						this.page = 1
						this.orderList = []
					}
					
					// 获取当前选择的状态
					const status = this.tabList[this.currentTab].status
					
					// 调用接口获取订单列表
					const res = await user.getOrderList({
						page: this.page,
						pageSize: this.pageSize,
						status: status
					})
					
					// 关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					if (res.code === 0 && res.data) {
						// 如果是第一页，直接替换数据
						if (this.page === 1) {
							this.orderList = res.data.list || []
						} else {
							// 否则追加数据
							this.orderList = [...this.orderList, ...(res.data.list || [])]
						}
						
						// 更新总数
						this.total = res.data.total || 0
						
						// 判断是否还有更多数据
						this.hasMore = this.orderList.length < this.total
					} else {
						throw new Error(res.message || '获取订单列表失败')
					}
				} catch (error) {
					// 确保异常情况下也关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					console.error('获取订单列表失败:', error)
					uni.showToast({
						title: '获取订单失败',
						icon: 'none'
					})
				} finally {
					this.isLoading = false
					
					// 如果是下拉刷新，结束刷新状态
					if (this.isRefreshing) {
						this.isRefreshing = false
					}
					
					// 最后确保关闭加载提示
					if (loadingShown) {
						uni.hideLoading();
					}
				}
			},
			
			// 根据订单状态获取CSS类名
			getStatusClass(status) {
				switch (status) {
					case 0: return 'unpaid'
					case 1: return 'processing'
					case 2: return 'completed'
					case 3: return 'cancelled'
					case 4: return 'refunded'
					default: return ''
				}
			},
			
			// 返回上一页
			handleBack() {
				// 获取当前页面栈
				const pages = getCurrentPages();
				// 如果当前页面不是第一个页面，才可以返回
				if (pages.length > 1) {
					uni.navigateBack();
				} else {
					// 如果是第一个页面，则跳转到个人中心页
					uni.switchTab({
						url: '/pages/my/index'
					});
				}
			},
			
			// 切换选项卡
			switchTab(index) {
				if (this.currentTab === index) return
				this.currentTab = index
				this.page = 1
				this.hasMore = true
				this.loadOrderList(true)
			},
			
			// 下拉刷新
			onRefresh() {
				this.isRefreshing = true
				this.page = 1
				this.hasMore = true
				this.loadOrderList(true)
			},
			
			// 加载更多
			loadMore() {
				if (!this.hasMore || this.isLoading) return
				this.page++
				this.loadOrderList()
			},
			
			// 点击订单
			handleOrderClick(order) {
				uni.navigateTo({
					url: `/pages/my/orders/detail?orderNo=${order.orderNo}`
				})
			},
			
			// 联系客服
			handleContactService(order) {
				uni.makePhoneCall({
					phoneNumber: '10086',
					fail: () => {
						uni.showToast({
							title: '拨打电话失败',
							icon: 'none'
						})
					}
				})
			},
			
			// 取消订单
			handleCancelOrder(order) {
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
									orderNo: order.orderNo,
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
								
								// 刷新订单列表
								this.onRefresh()
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
			async handlePayOrder(order) {
				// 标记加载状态
				let loadingShown = false;
				
				try {
					// 显示加载提示
					uni.showLoading({
						title: '正在支付...',
						mask: true // 使用蒙层防止用户点击
					});
					loadingShown = true;
					
					// 调用订单支付接口
					const payResult = await requestOrderPay({
						orderNo: order.orderNo
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
						
						// 刷新订单列表
						this.onRefresh()
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
			
			// 再次购买
			handleRebuyOrder(order) {
				uni.navigateTo({
					url: '/pages/publish/info/index'
				})
			}
		}
	}
</script>

<style>
	.container {
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
	
	.tab-bar {
		display: flex;
		background-color: #ffffff;
		border-bottom: 1rpx solid #f5f5f5;
		height: 80rpx;
		width: 100%;
	}
	
	.tab-item {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
	}
	
	.tab-text {
		font-size: 28rpx;
		color: #666666;
		padding: 0 20rpx;
	}
	
	.tab-item.active .tab-text {
		color: #fc3e2b;
		font-weight: 500;
	}
	
	.tab-item.active::after {
		content: '';
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 40rpx;
		height: 4rpx;
		background-color: #fc3e2b;
		border-radius: 2rpx;
	}
	
	.order-scroll {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: #f5f5f5;
		z-index: 1;
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
		height: calc(100vh - var(--status-bar-height) - var(--nav-bar-height) - 80rpx);
	}
	
	.order-list {
		padding: 10rpx;
		min-height: calc(100vh - var(--status-bar-height) - var(--nav-bar-height) - 80rpx - 10rpx);
		box-sizing: border-box;
		padding-bottom: calc(env(safe-area-inset-bottom) + 30rpx);
	}
	
	.order-item {
		background-color: #ffffff;
		border-radius: 16rpx;
		margin-bottom: 12rpx;
		padding: 20rpx;
	}
	
	.order-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;
	}
	
	.order-time {
		font-size: 26rpx;
		color: #999999;
	}
	
	.order-status {
		font-size: 26rpx;
	}
	
	.order-status.processing {
		color: #409eff;
	}
	
	.order-status.unpaid {
		color: #fc3e2b;
	}
	
	.order-status.completed {
		color: #67c23a;
	}
	
	.order-status.cancelled {
		color: #909399;
	}
	
	.order-status.refunded {
		color: #e6a23c;
	}
	
	.order-content {
		display: flex;
		margin-bottom: 20rpx;
	}
	
	.order-icon {
		width: 160rpx;
		height: 160rpx;
		border-radius: 8rpx;
		margin-right: 20rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: #f8f8f8;
	}
	
	.order-info {
		flex: 1;
		display: flex;
		flex-direction: column;
	}
	
	.order-title {
		font-size: 28rpx;
		color: #333333;
		font-weight: 500;
		margin-bottom: 8rpx;
	}
	
	.order-no {
		font-size: 26rpx;
		color: #666666;
		margin-bottom: 16rpx;
	}
	
	.price-row {
		display: flex;
		align-items: baseline;
		margin-top: auto;
	}
	
	.price-label {
		font-size: 26rpx;
		color: #999999;
		margin-right: 8rpx;
	}
	
	.price-value {
		font-size: 32rpx;
		color: #333333;
		font-weight: 500;
	}
	
	.order-footer {
		border-top: 1rpx solid #f5f5f5;
		padding-top: 20rpx;
	}
	
	.btn-group {
		display: flex;
		justify-content: flex-end;
		gap: 20rpx;
	}
	
	.order-btn {
		padding: 12rpx 24rpx;
		border-radius: 32rpx;
		border: 1rpx solid #dddddd;
	}
	
	.order-btn text {
		font-size: 26rpx;
		color: #666666;
	}
	
	.order-btn.primary {
		background-color: #fc3e2b;
		border-color: #fc3e2b;
	}
	
	.order-btn.primary text {
		color: #ffffff;
	}
	
	.loading-more, .no-more {
		text-align: center;
		padding: 30rpx 0;
	}
	
	.loading-text, .no-more-text {
		font-size: 26rpx;
		color: #999999;
	}
	
	.empty-state {
		min-height: calc(100vh - var(--status-bar-height) - var(--nav-bar-height) - 80rpx - 200rpx);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding-top: 100rpx;
	}
	
	.empty-image {
		width: 200rpx;
		height: 200rpx;
		margin-bottom: 20rpx;
	}
	
	.empty-text {
		font-size: 28rpx;
		color: #999999;
	}
</style> 