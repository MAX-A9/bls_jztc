<template>
	<view class="container">
		<!-- 状态栏占位 -->
		<view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
		
		<!-- 自定义导航栏 -->
		<view class="nav-bar">
			<view class="nav-content">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">兑换记录</text>
			</view>
		</view>
		
		<!-- 记录列表 -->
		<view class="record-list">
			<view class="record-item" v-for="item in recordList" :key="item.id">
				<view class="record-top">
					<view class="record-date">{{item.exchangeTime}}</view>
					<view class="record-status" :class="getStatusClass(item.status)">{{formatStatus(item.status)}}</view>
				</view>
				<view class="record-content">
					<view class="product-info">
						<text class="product-name">{{item.productName}}</text>
						<text class="duration">支出时长：{{item.duration}}天</text>
					</view>
					<view class="record-account">
						<text class="account-label">兑换账号：</text>
						<text class="account-value">{{item.rechargeAccount}}</text>
					</view>
					<view class="record-remark" v-if="item.remark">
						<text class="remark-label">备注：</text>
						<text class="remark-content">{{item.remark}}</text>
					</view>
				</view>
			</view>
			
			<!-- 加载更多 -->
			<uni-load-more :status="loadMoreStatus" v-if="recordList.length > 0"></uni-load-more>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="recordList.length === 0 && !isLoading">
				<image class="empty-image" src="/static/images/empty-record.png" mode="aspectFit"></image>
				<text class="empty-text">暂无兑换记录</text>
			</view>
		</view>
	</view>
</template>

<script>
	import deviceMixin from '@/mixins/device.js'
	import { vip } from '@/apis/index.js'
	import uniIcons from '@dcloudio/uni-ui/lib/uni-icons/uni-icons'
	import uniLoadMore from '@dcloudio/uni-ui/lib/uni-load-more/uni-load-more'
	
	export default {
		components: {
			uniIcons,
			uniLoadMore
		},
		mixins: [deviceMixin],
		data() {
			return {
				recordList: [],
				page: 1,
				size: 10,
				totalPages: 0,
				isLoading: false,
				loadMoreStatus: 'more' // more-加载前 loading-加载中 noMore-没有更多
			}
		},
		onLoad() {
			this.getExchangeRecords()
		},
		// 下拉刷新
		onPullDownRefresh() {
			this.page = 1
			this.recordList = []
			this.getExchangeRecords().then(() => {
				uni.stopPullDownRefresh()
			})
		},
		// 上拉加载更多
		onReachBottom() {
			if (this.page < this.totalPages && this.loadMoreStatus !== 'loading') {
				this.page++
				this.loadMoreStatus = 'loading'
				this.getExchangeRecords()
			}
		},
		methods: {
			handleBack() {
				uni.navigateBack()
			},
			// 格式化状态为中文
			formatStatus(status) {
				// 处理英文状态转中文
				if (!status) return '未知状态'
				
				// 转为小写进行比较
				const statusLower = typeof status === 'string' ? status.toLowerCase() : status
				
				if (statusLower === 'processing' || statusLower === '处理中') {
					return '处理中'
				} else if (statusLower === 'completed' || statusLower === '已完成') {
					return '已完成'
				} else if (statusLower === 'failed' || statusLower === '失败') {
					return '失败'
				} else {
					return status // 返回原始状态
				}
			},
			// 获取状态样式类
			getStatusClass(status) {
				// 转为小写进行比较
				const statusLower = typeof status === 'string' ? status.toLowerCase() : status
				
				if (statusLower === 'processing' || statusLower === '处理中') {
					return 'status-processing'
				} else if (statusLower === 'completed' || statusLower === '已完成') {
					return 'status-success'
				} else if (statusLower === 'failed' || statusLower === '失败') {
					return 'status-failed'
				} else {
					return 'status-processing' // 默认处理中样式
				}
			},
			// 获取兑换记录
			async getExchangeRecords() {
				if (this.isLoading) return
				
				try {
					this.isLoading = true
					
					const params = {
						page: this.page,
						size: this.size
					}
					
					const res = await vip.getExchangeRecords(params)
					
					if (res.code === 0 && res.data) {
						if (this.page === 1) {
							this.recordList = res.data.list || []
						} else {
							this.recordList = [...this.recordList, ...(res.data.list || [])]
						}
						
						this.totalPages = res.data.pages || 0
						
						// 更新加载状态
						if (this.page >= this.totalPages) {
							this.loadMoreStatus = 'noMore'
						} else {
							this.loadMoreStatus = 'more'
						}
					} else {
						uni.showToast({
							title: res.message || '获取兑换记录失败',
							icon: 'none'
						})
					}
				} catch (error) {
					console.error('获取兑换记录失败:', error)
					uni.showToast({
						title: '获取记录失败，请重试',
						icon: 'none'
					})
				} finally {
					this.isLoading = false
				}
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background: #f5f5f5;
	}
	
	.status-bar {
		background: #ffffff;
	}
	
	.nav-bar {
		background: #ffffff;
		box-shadow: 0 2rpx 6rpx rgba(0, 0, 0, 0.05);
	}
	
	.nav-content {
		height: 44px;
		display: flex;
		align-items: center;
		padding: 0 30rpx;
		position: relative;
	}
	
	.back-btn {
		position: absolute;
		left: 30rpx;
		height: 100%;
		display: flex;
		align-items: center;
	}
	
	.nav-title {
		color: #333333;
		font-size: 32rpx;
		font-weight: bold;
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
	}
	
	.record-list {
		padding: 20rpx;
	}
	
	.record-item {
		background: #ffffff;
		border-radius: 12rpx;
		padding: 24rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 6rpx rgba(0, 0, 0, 0.05);
	}
	
	.record-top {
		display: flex;
		justify-content: space-between;
		padding-bottom: 16rpx;
		border-bottom: 1rpx solid #f0f0f0;
	}
	
	.record-date {
		font-size: 24rpx;
		color: #666666;
	}
	
	.record-status {
		font-size: 24rpx;
		padding: 2rpx 12rpx;
		border-radius: 20rpx;
	}
	
	.status-success {
		color: #27ae60;
		background-color: rgba(39, 174, 96, 0.1);
	}
	
	.status-processing {
		color: #f39c12;
		background-color: rgba(243, 156, 18, 0.1);
	}
	
	.status-failed {
		color: #e74c3c;
		background-color: rgba(231, 76, 60, 0.1);
	}
	
	.record-content {
		padding-top: 16rpx;
	}
	
	.product-info {
		margin-bottom: 12rpx;
	}
	
	.product-name {
		font-size: 28rpx;
		color: #333333;
		font-weight: bold;
		margin-right: 20rpx;
	}
	
	.duration {
		font-size: 24rpx;
		color: #666666;
	}
	
	.record-account {
		font-size: 26rpx;
		color: #333333;
		margin-bottom: 8rpx;
	}
	
	.account-label {
		color: #666666;
	}
	
	.record-remark {
		font-size: 26rpx;
		color: #333333;
	}
	
	.remark-label {
		color: #666666;
	}
	
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 100rpx 0;
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