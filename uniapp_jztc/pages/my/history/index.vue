<template>
	<view class="container">
		<!-- 日期筛选 -->
		<view class="date-filter-container">
			<view class="date-filter">
				<view 
					class="date-item" 
					v-for="(item, index) in dateFilters" 
					:key="index"
					:class="{ active: currentFilter === index }"
					@tap="switchFilter(index)"
				>
					<text class="date-text">{{ item }}</text>
				</view>
			</view>
			<view class="clear-btn" @tap="handleClearHistory">
				<uni-icons type="trash" size="16" color="#666666"></uni-icons>
				<text class="clear-text">清空{{ currentFilter > 0 ? dateFilters[currentFilter] : '全部' }}</text>
			</view>
		</view>
		
		<scroll-view 
			class="content-scroll" 
			scroll-y 
			:refresher-enabled="true"
			:refresher-triggered="isRefreshing"
			@refresherrefresh="onRefresh"
			:show-scrollbar="false"
			:enhanced="true"
			:bounces="true"
			:fast-deceleration="true"
			@scrolltolower="loadMore"
		>
			<!-- 加载中 -->
			<view class="loading-container" v-if="loading && historyList.length === 0">
				<view class="loading-indicator"></view>
				<text class="loading-text">数据加载中...</text>
			</view>
			
			<!-- 浏览历史列表 -->
			<view class="history-list" v-if="historyList.length > 0">
				<view class="date-group" v-for="(group, date) in groupedHistory" :key="date">
					<view class="date-header">
						<text class="date-label">{{ date }}</text>
					</view>
					
					<view 
						class="history-item"
						v-for="(item, index) in group"
						:key="index"
						@tap="handleItemClick(item)"
						:class="{ 'offline-item': item.contentStatus === 0 || item.contentStatus === 2 }"
					>
						<view class="offline-mask" v-if="item.contentStatus === 0 || item.contentStatus === 2">
							<text class="offline-text">已下架</text>
						</view>
						<image 
							class="item-image" 
							:src="item.image" 
							mode="aspectFill" 
							v-if="item.image"
							@error="handleImageError(item)"
						></image>
						<view class="item-content" :class="{ 'full-width': !item.image }">
							<text class="item-title" :class="{ 'offline-title': item.contentStatus === 0 || item.contentStatus === 2 }">{{ item.title }}</text>
							<text class="item-price" v-if="item.price > 0">¥{{ item.price }}</text>
							<view class="item-footer">
								<text class="item-time">{{ item.time }}</text>
								<view class="item-type" :class="{ 'item-type-offline': item.contentStatus === 0 || item.contentStatus === 2 }">
									<text>{{ item.category }}</text>
									<view class="item-status" v-if="item.contentStatus === 0 || item.contentStatus === 2">已下架</view>
								</view>
							</view>
						</view>
					</view>
				</view>
				
				<!-- 加载更多提示 -->
				<view class="load-more" v-if="historyList.length > 0">
					<view v-if="loading && historyList.length > 0" class="loading">
						<view class="loading-indicator-small"></view>
						<text>正在加载...</text>
					</view>
					<text v-else-if="hasMore">上拉加载更多</text>
					<text v-else>没有更多数据了</text>
				</view>
			</view>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="!loading && historyList.length === 0">
				<image class="empty-image" src="/static/empty/history.png" mode="aspectFit"></image>
				<text class="empty-text">暂无浏览记录</text>
				<button class="action-btn" @tap="handleExplore">去发现</button>
			</view>
		</scroll-view>
	</view>
</template>

<script>
	import { getBrowseHistoryList, clearBrowseHistory } from '@/apis/content.js'
	import { formatTimeAgo } from '@/utils/date.js'
	
	export default {
		data() {
			return {
				currentFilter: 0,
				isRefreshing: false,
				dateFilters: ['全部', '今天', '昨天', '本周', '本月'],
				timeTypeMap: {
					0: 'all',
					1: 'today',
					2: 'yesterday',
					3: 'this_week',
					4: 'this_month'
				},
				historyList: [],
				page: 1,
				pageSize: 10,
				hasMore: true,
				loading: false,
				total: 0
			}
		},
		created() {
			// 加载历史记录
			this.loadBrowseHistory()
		},
		computed: {
			groupedHistory() {
				// 按日期分组
				const grouped = {}
				
				this.historyList.forEach(item => {
					// 提取日期部分作为分组键
					const date = item.browseTime.split(' ')[0]
					if (!grouped[date]) {
						grouped[date] = []
					}
					grouped[date].push(item)
				})
				
				return grouped
			}
		},
		methods: {
			switchFilter(index) {
				if (this.currentFilter === index) return
				this.currentFilter = index
				this.page = 1
				this.historyList = []
				this.hasMore = true
				this.loadBrowseHistory()
			},
			
			// 加载浏览历史数据
			async loadBrowseHistory() {
				if (this.loading) return
				
				try {
					this.loading = true
					
					const params = {
						timeType: this.timeTypeMap[this.currentFilter],
						page: this.page,
						pageSize: this.pageSize
					}
					
					const result = await getBrowseHistoryList(params)
					
					if (result.code === 0 && result.data) {
						const { list = [], total = 0, page = 1 } = result.data
						
						// 格式化数据
						const formattedList = list.map(item => {
							return {
								id: item.id,
								contentId: item.contentId,
								title: item.contentTitle || '未知标题',
								image: item.contentCover || '',
								price: item.price || 0,
								contentType: item.contentType,
								browseTime: item.browseTime,
								time: item.browseTime.split(' ')[1],
								category: item.category || '未分类',
								contentStatus: item.contentStatus
							}
						})
						
						// 第一页直接替换，之后追加
						if (this.page === 1) {
							this.historyList = formattedList
						} else {
							this.historyList = [...this.historyList, ...formattedList]
						}
						
						this.total = total
						this.hasMore = this.historyList.length < total
						this.page++
					} else {
						throw new Error(result.message || '获取浏览历史失败')
					}
				} catch (error) {
					console.error('加载浏览历史失败:', error)
					uni.showToast({
						title: '加载失败，请重试',
						icon: 'none',
						duration: 2000
					})
				} finally {
					this.loading = false
					if (this.isRefreshing) {
						this.isRefreshing = false
					}
				}
			},
			
			async onRefresh() {
				this.isRefreshing = true
				this.page = 1
				this.hasMore = true
				await this.loadBrowseHistory()
				uni.showToast({
					title: '刷新成功',
					icon: 'success'
				})
			},
			
			handleItemClick(item) {
				// 已下架内容不进行跳转，显示提示信息
				if (item.contentStatus === 0 || item.contentStatus === 2) {
					uni.showToast({
						title: '该内容已下架',
						icon: 'none',
						duration: 2000
					});
					return;
				}
				
				// 根据不同内容类型跳转到不同页面
				let url = ''
				switch (item.contentType) {
					case 'idle':
						url = `/pages/community/detail?id=${item.contentId}`
						break
					case 'article':
					case 'info':
					default:
						url = `/pages/content/detail?id=${item.contentId}`
						break
				}
				
				uni.navigateTo({ url })
			},
			
			async handleClearHistory() {
				uni.showModal({
					title: '提示',
					content: '确定要清空浏览记录吗？',
					success: async (res) => {
						if (res.confirm) {
							try {
								// 获取当前筛选类型
								const timeType = this.timeTypeMap[this.currentFilter];
								// 调用清空接口，传递当前时间筛选类型
								const result = await clearBrowseHistory(timeType);
								
								if (result.code === 0) {
									this.historyList = [];
									this.page = 1;
									this.hasMore = false;
									uni.showToast({
										title: '已清空记录',
										icon: 'success'
									});
								} else {
									throw new Error(result.message || '清空失败');
								}
							} catch (error) {
								console.error('清空浏览记录失败:', error);
								uni.showToast({
									title: '清空失败，请重试',
									icon: 'none'
								});
							}
						}
					}
				});
			},
			
			handleExplore() {
				uni.switchTab({
					url: '/pages/index/index'
				})
			},
			
			// 加载更多
			loadMore() {
				if (this.hasMore && !this.loading) {
					this.loadBrowseHistory()
				}
			},
			
			handleImageError(item) {
				// 图片加载错误时使用空字符串，让布局自动调整
				item.image = '';
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background-color: #f5f5f5;
		display: flex;
		flex-direction: column;
		height: 100vh;
		overflow: hidden;
		will-change: transform;
	}
	
	.date-filter-container {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20rpx 30rpx;
		background-color: #ffffff;
		border-bottom: 1rpx solid #f0f0f0;
		position: sticky;
		top: 0;
		z-index: 100;
	}
	
	.date-filter {
		display: flex;
		background-color: #f5f5f5;
		border-radius: 32rpx;
		padding: 4rpx;
	}
	
	.date-item {
		padding: 10rpx 20rpx;
		border-radius: 28rpx;
		transition: all 0.3s;
		position: relative;
	}
	
	.date-item.active {
		background-color: #ffffff;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
	}
	
	.date-text {
		font-size: 24rpx;
		color: #666666;
	}
	
	.date-item.active .date-text {
		color: #fc3e2b;
		font-weight: 500;
	}
	
	.clear-btn {
		display: flex;
		align-items: center;
		padding: 8rpx 16rpx;
		background-color: #f5f5f5;
		border-radius: 24rpx;
		transition: all 0.3s;
	}
	
	.clear-btn:active {
		transform: scale(0.95);
		opacity: 0.9;
	}
	
	.clear-text {
		font-size: 24rpx;
		color: #666666;
		margin-left: 4rpx;
	}
	
	.content-scroll {
		flex: 1;
		height: 0;
		box-sizing: border-box;
		-webkit-overflow-scrolling: touch;
	}
	
	.history-list {
		padding: 20rpx;
		transform: translateZ(0);
	}
	
	.date-header {
		display: flex;
		align-items: center;
		padding: 20rpx 0 16rpx;
	}
	
	.date-header::before, 
	.date-header::after {
		content: '';
		flex: 1;
		height: 1px;
		background-color: #f0f0f0;
	}
	
	.date-header::before {
		margin-right: 20rpx;
	}
	
	.date-header::after {
		margin-left: 20rpx;
	}
	
	.date-label {
		font-size: 28rpx;
		color: #666666;
		font-weight: 500;
		white-space: nowrap;
	}
	
	.history-item {
		display: flex;
		background-color: #ffffff;
		border-radius: 12rpx;
	padding: 20rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
	will-change: transform;
	position: relative;
	transition: all 0.3s;
	overflow: hidden;
}

.history-item:active {
	transform: scale(0.98);
	opacity: 0.9;
}

/* 已下架商品样式 */
.offline-item {
	background-color: #f9f9f9;
}

.offline-mask {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.03);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 1;
	pointer-events: none;
}

.offline-text {
	font-size: 32rpx;
	color: rgba(153, 153, 153, 0.8);
	transform: rotate(-15deg);
	font-weight: bold;
	border: 2rpx solid rgba(153, 153, 153, 0.3);
	padding: 8rpx 30rpx;
	border-radius: 8rpx;
	background-color: rgba(255, 255, 255, 0.8);
}

.offline-title {
	color: #999999;
}

.item-image {
	width: 160rpx;
	height: 160rpx;
	border-radius: 8rpx;
	margin-right: 20rpx;
	background-color: #f5f5f5;
	position: relative;
	z-index: 0;
}

.item-content {
	flex: 1;
	display: flex;
	flex-direction: column;
	position: relative;
	z-index: 0;
}

.item-content.full-width {
	width: 100%;
}

.item-title {
	font-size: 28rpx;
	color: #333333;
	margin-bottom: 16rpx;
	line-height: 1.4;
	display: -webkit-box;
	-webkit-box-orient: vertical;
	-webkit-line-clamp: 2;
	line-clamp: 2; /* 标准属性，提高兼容性 */
	overflow: hidden;
	transition: color 0.3s;
}

.item-price {
	font-size: 32rpx;
	color: #fc3e2b;
	font-weight: 500;
	margin-bottom: 16rpx;
}

.item-footer {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-top: auto;
}

.item-time {
	font-size: 24rpx;
	color: #999999;
}

.item-type {
	font-size: 24rpx;
	color: #666666;
	background-color: #f5f5f5;
	padding: 4rpx 12rpx;
	border-radius: 20rpx;
	display: flex;
	align-items: center;
}

.item-type-offline {
	background-color: #f0f0f0;
	color: #999999;
}

.item-status {
	font-size: 22rpx;
	color: #ff6b6b;
	margin-left: 10rpx;
	background-color: rgba(255, 107, 107, 0.1);
	padding: 2rpx 8rpx;
	border-radius: 16rpx;
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding-top: 200rpx;
}

.empty-image {
	width: 240rpx;
	height: 240rpx;
	margin-bottom: 30rpx;
}

.empty-text {
	font-size: 30rpx;
	color: #999999;
	margin-bottom: 40rpx;
}

.action-btn {
	width: 240rpx;
	height: 80rpx;
	background: linear-gradient(90deg, #ff7c4d, #fc3e2b);
	color: #ffffff;
	font-size: 28rpx;
	border-radius: 40rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 4rpx 12rpx rgba(252, 62, 43, 0.2);
}

.loading-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 200rpx 0;
}

.loading-indicator {
	width: 48rpx;
	height: 48rpx;
	border: 4rpx solid #fc3e2b;
	border-radius: 50%;
	border-top-color: transparent;
	animation: spin 0.75s linear infinite;
	margin-bottom: 20rpx;
}

.loading-text {
	font-size: 28rpx;
	color: #666666;
}

.load-more {
	display: flex;
	justify-content: center;
	align-items: center;
	padding: 20rpx 0;
	color: #999999;
	font-size: 26rpx;
	margin-top: 10rpx;
}

.loading {
	display: flex;
	align-items: center;
	justify-content: center;
	height: 60rpx;
}

.loading-indicator-small {
	width: 24rpx;
	height: 24rpx;
	border: 2rpx solid #fc3e2b;
	border-radius: 50%;
	border-top-color: transparent;
	animation: spin 0.75s linear infinite;
	margin-right: 12rpx;
}

@keyframes spin {
	0% {
		transform: rotate(0deg);
	}
	100% {
		transform: rotate(360deg);
	}
}
</style> 