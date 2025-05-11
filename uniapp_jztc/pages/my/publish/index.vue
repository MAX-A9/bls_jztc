<template>
	<view class="container">
		<!-- 选项卡 -->
		<view class="tab-container">
			<view 
				class="tab-item" 
				:class="{ active: activeType === 0 }" 
				@tap="switchTab(0)"
			>
				全部
			</view>
			<view 
				class="tab-item" 
				:class="{ active: activeType === 1 }" 
				@tap="switchTab(1)"
			>
				贴子
			</view>
			<view 
				class="tab-item" 
				:class="{ active: activeType === 2 }" 
				@tap="switchTab(2)"
			>
				闲置
			</view>
		</view>
		
		<!-- 内容区域 -->
		<scroll-view 
			class="content-scroll" 
			scroll-y 
			:refresher-enabled="true"
			:refresher-triggered="isRefreshing"
			@refresherrefresh="onRefresh"
			:show-scrollbar="false"
			@scrolltolower="loadMore"
		>
			<view class="publish-list" v-if="publishList.length > 0">
				<!-- 使用简化的内容列表样式 -->
				<view 
					class="content-item card" 
					v-for="(item, index) in publishList"
					:key="item.id"
					@tap="handleItemClick(item)"
				>
					<view class="content-main">
						<view class="tag-row">
							<!-- 状态标签 -->
							<view class="tag status-tag" :class="getStatusClass(item.status)">{{ item.status }}</view>
							<!-- 类目标签 -->
							<view class="tag category-tag">{{ item.category }}</view>
						</view>
						
						<!-- 标题 -->
						<view class="content-title">{{ item.title }}</view>
						
						<!-- 发布时间 -->
						<view class="time-row">
							<text class="publish-time">{{ item.publishedAt }}</text>
						</view>
					</view>
					
					<!-- 操作按钮 -->
					<view class="action-btns">
						<button class="action-btn edit-btn" @tap.stop="handleEdit(item)">编辑</button>
						<button class="action-btn delete-btn" @tap.stop="handleDelete(item, index)">删除</button>
					</view>
				</view>
			</view>
			
			<!-- 加载更多 -->
			<view class="loading-more" v-if="publishList.length > 0 && hasMore">
				<text class="loading-text">加载中...</text>
			</view>
			
			<!-- 没有更多 -->
			<view class="no-more" v-if="publishList.length > 0 && !hasMore">
				<text class="no-more-text">没有更多内容了</text>
			</view>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="publishList.length === 0 && !loading">
				<image class="empty-image" src="/static/images/empty.png" mode="aspectFit"></image>
				<text class="empty-text">暂无发布内容</text>
				<button class="publish-btn" @tap="handlePublish">去发布</button>
			</view>
		</scroll-view>
	</view>
</template>

<script>
	import { publish } from '@/apis'
	
	export default {
		data() {
			return {
				activeType: 0, // 0:全部, 1:普通信息, 2:闲置
				isRefreshing: false,
				loading: false,
				publishList: [],
				page: 1,
				pageSize: 10,
				total: 0,
				hasMore: true
			}
		},
		onLoad() {
			// 加载数据
			this.loadPublishList()
		},
		methods: {
			// 切换选项卡
			switchTab(type) {
				if (this.activeType === type) return
				this.activeType = type
				this.resetAndLoadList()
			},
			// 重置列表并重新加载
			resetAndLoadList() {
				this.publishList = []
				this.page = 1
				this.hasMore = true
				this.loadPublishList()
			},
			// 加载发布列表
			async loadPublishList() {
				if (!this.hasMore || this.loading) return
				
				this.loading = true
				try {
					const params = {
						page: this.page,
						pageSize: this.pageSize,
						type: this.activeType
					}
					
					const res = await publish.getMyPublishList(params)
					
					if (res.code === 0 && res.data) {
						const { list, total, page, pages } = res.data
						
						// 追加数据
						this.publishList = this.page === 1 ? list : [...this.publishList, ...list]
						this.total = total
						
						// 判断是否还有更多数据
						this.hasMore = page < pages
						
						// 页码加1，为下次加载做准备
						if (this.hasMore) {
							this.page++
						}
					} else {
						uni.showToast({
							title: res.message || '加载失败',
							icon: 'none'
						})
					}
				} catch (error) {
					console.error('加载发布列表失败', error)
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					})
				} finally {
					this.loading = false
				}
			},
			// 下拉刷新
			async onRefresh() {
				this.isRefreshing = true
				this.resetAndLoadList()
				setTimeout(() => {
					this.isRefreshing = false
				}, 1000)
			},
			// 上拉加载更多
			loadMore() {
				this.loadPublishList()
			},
			// 获取状态样式
			getStatusClass(status) {
				// 根据状态返回对应的class
				switch(status) {
					case '已发布':
						return 'published'
					case '已下架':
						return 'offline'
					case '已售出':
						return 'sold'
					case '审核中':
						return 'pending'
					default:
						return 'default'
				}
			},
			// 点击内容项
			handleItemClick(item) {
				uni.navigateTo({
					url: `/pages/content/detail?id=${item.id}`
				})
			},
			// 编辑
			handleEdit(item) {
				uni.navigateTo({
					url: `/pages/publish/edit?id=${item.id}`
				})
			},
			// 删除
			handleDelete(item, index) {
				uni.showModal({
					title: '提示',
					content: '确定要删除该发布吗？',
					success: async (res) => {
						if (res.confirm) {
							try {
								// 这里应该调用删除API
								const result = await publish.deleteContent(item.id)
								if (result.code === 0) {
									this.publishList.splice(index, 1)
									uni.showToast({
										title: '删除成功',
										icon: 'success'
									})
								} else {
									uni.showToast({
										title: result.message || '删除失败',
										icon: 'none'
									})
								}
							} catch (error) {
								console.error('删除失败', error)
								uni.showToast({
									title: '删除失败',
									icon: 'none'
								})
							}
						}
					}
				})
			},
			// 去发布
			handlePublish() {
				uni.switchTab({
					url: '/pages/publish/index'
				})
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
		height: 100vh; /* 添加固定高度 */
	}
	
	/* 选项卡样式 */
	.tab-container {
		display: flex;
		background-color: #fff;
		padding: 0 30rpx;
		height: 88rpx;
		border-bottom: 1rpx solid #f0f0f0;
		position: sticky;
		top: 0;
		z-index: 10;
		flex-shrink: 0; /* 防止内容压缩选项卡 */
	}
	
	.tab-item {
		position: relative;
		flex: 1;
		height: 88rpx;
		display: flex;
		justify-content: center;
		align-items: center;
		font-size: 30rpx;
		color: #666;
	}
	
	.tab-item.active {
		color: #333;
		font-weight: bold;
	}
	
	.tab-item.active::after {
		content: '';
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 48rpx;
		height: 4rpx;
		background-color: #fc3e2b;
		border-radius: 2rpx;
	}
	
	.content-scroll {
		flex: 1;
		box-sizing: border-box;
		height: calc(100vh - 88rpx); /* 计算剩余高度 */
		overflow-y: auto; /* 确保可以滚动 */
	}
	
	.publish-list {
		padding: 20rpx;
		padding-bottom: 20rpx;
	}
	
	/* 内容项样式 */
	.content-item.card {
		display: flex;
		flex-direction: column;
		padding: 24rpx;
		background-color: #ffffff;
		border-radius: 16rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.06);
	}
	
	.content-main {
		flex: 1;
	}
	
	.tag-row {
		display: flex;
		gap: 12rpx;
		margin-bottom: 12rpx;
	}
	
	.tag {
		padding: 4rpx 12rpx;
		border-radius: 4rpx;
		font-size: 22rpx;
	}
	
	.status-tag {
		color: #ffffff;
	}
	
	.status-tag.published {
		background-color: #409eff;
	}
	
	.status-tag.offline {
		background-color: #909399;
	}
	
	.status-tag.sold {
		background-color: #67c23a;
	}
	
	.status-tag.pending {
		background-color: #e6a23c;
	}
	
	.status-tag.default {
		background-color: #909399;
	}
	
	.category-tag {
		background-color: #f5f5f5;
		color: #666666;
	}
	
	.content-title {
		font-size: 32rpx;
		font-weight: 600;
		color: #333333;
		line-height: 1.4;
		margin-bottom: 12rpx;
	}
	
	.time-row {
		margin-bottom: 16rpx;
	}
	
	.publish-time {
		font-size: 24rpx;
		color: #999999;
	}
	
	.action-btns {
		display: flex;
		justify-content: flex-end;
		gap: 16rpx;
		margin-top: 16rpx;
		padding-top: 16rpx;
		border-top: 1rpx solid #f0f0f0;
	}
	
	.action-btn {
		font-size: 26rpx;
		padding: 8rpx 24rpx;
		border-radius: 30rpx;
		margin: 0;
		line-height: 1.5;
	}
	
	.edit-btn {
		color: #409eff;
		background-color: rgba(64, 158, 255, 0.1);
		border: 1rpx solid #409eff;
	}
	
	.delete-btn {
		color: #f56c6c;
		background-color: rgba(245, 108, 108, 0.1);
		border: 1rpx solid #f56c6c;
	}
	
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding-top: 200rpx;
	}
	
	.empty-image {
		width: 200rpx;
		height: 200rpx;
		margin-bottom: 30rpx;
	}
	
	.empty-text {
		font-size: 28rpx;
		color: #999999;
		margin-bottom: 40rpx;
	}
	
	.publish-btn {
		width: 240rpx;
		height: 80rpx;
		background-color: #fc3e2b;
		color: #ffffff;
		font-size: 28rpx;
		border-radius: 40rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.loading-more, .no-more {
		text-align: center;
		padding: 20rpx 0;
	}
	
	.loading-text, .no-more-text {
		font-size: 24rpx;
		color: #999;
	}
</style> 