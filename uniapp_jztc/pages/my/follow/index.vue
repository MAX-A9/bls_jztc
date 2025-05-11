<template>
	<view class="container">
		<!-- 用户关注列表 -->
		<scroll-view 
			class="content-scroll" 
			scroll-y 
			:refresher-enabled="true"
			:refresher-triggered="isRefreshing"
			@refresherrefresh="onRefresh"
			:show-scrollbar="false"
			@scrolltolower="loadMore"
		>
			<!-- 加载中 -->
			<view class="loading-container" v-if="loading && userList.length === 0">
				<view class="loading-indicator"></view>
				<text class="loading-text">数据加载中...</text>
			</view>
			
			<view class="follow-list" v-if="userList.length > 0">
				<view 
					class="user-item"
					v-for="(item, index) in userList"
					:key="item.client_id"
					@tap="handleUserClick(item)"
				>
					<image class="user-avatar" :src="item.avatar_url || '/static/avatar/default.png'" mode="aspectFill"></image>
					<view class="user-info">
						<text class="user-name">{{ item.real_name }}</text>
						<view class="user-meta">
							<text class="publish-count">已发布 {{ item.publish_count || 0 }} 条</text>
						</view>
					</view>
					<button class="follow-btn" @tap.stop="toggleFollow(item, index)">已关注</button>
				</view>
				
				<!-- 加载更多 -->
				<view class="loading-more" v-if="hasMore && !loading">
					<text class="loading-text">上拉加载更多</text>
				</view>
				<view class="loading-more" v-if="loading && userList.length > 0">
					<text class="loading-text">加载中...</text>
				</view>
				<view class="loading-more" v-if="!hasMore && userList.length > 0">
					<text class="loading-text">没有更多数据了</text>
				</view>
			</view>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="!loading && userList.length === 0">
				<image class="empty-image" src="/static/images/empty.png" mode="aspectFit"></image>
				<text class="empty-text">暂无关注内容</text>
				<button class="action-btn" @tap="handleExplore">去发现</button>
			</view>
		</scroll-view>
	</view>
</template>

<script>
	import { getFollowingList, unfollowUser } from '@/apis/content.js'
	
	export default {
		data() {
			return {
				isRefreshing: false,
				loading: false,
				userList: [],
				currentPage: 1,
				pageSize: 10,
				hasMore: true,
				total: 0
			}
		},
		onLoad() {
			// 加载关注列表
			this.loadFollowingList()
		},
		methods: {
			// 加载关注列表
			async loadFollowingList() {
				if (this.loading) return
				
				this.loading = true
				try {
					const result = await getFollowingList({
						page: this.currentPage,
						size: this.pageSize
					})
					
					if (result.code === 0 && result.data) {
						// 第一页时直接赋值，加载更多时追加
						if (this.currentPage === 1) {
							this.userList = result.data.list || []
						} else {
							this.userList = [...this.userList, ...(result.data.list || [])]
						}
						
						// 更新总数和判断是否有更多
						this.total = result.data.total || 0
						this.hasMore = this.userList.length < this.total
					} else {
						uni.showToast({
							title: result.message || '加载失败',
							icon: 'none'
						})
					}
				} catch (error) {
					console.error('加载关注列表失败:', error)
					uni.showToast({
						title: '网络异常，请稍后重试',
						icon: 'none'
					})
				} finally {
					this.loading = false
					if (this.isRefreshing) {
						this.isRefreshing = false
					}
				}
			},
			
			// 下拉刷新
			async onRefresh() {
				this.isRefreshing = true
				this.currentPage = 1
				await this.loadFollowingList()
				
				// 防止加载成功但数据为空时无法触发刷新完成
				if (this.isRefreshing) {
					this.isRefreshing = false
				}
			},
			
			// 加载更多
			loadMore() {
				if (this.loading || !this.hasMore) return
				
				this.currentPage++
				this.loadFollowingList()
			},
			
			// 点击用户进入详情
			handleUserClick(user) {
				if (!user || !user.client_id) return
				
				uni.navigateTo({
					url: `/pages/user/detail?id=${user.client_id}`
				})
			},
			
			// 取消关注
			async toggleFollow(user, index) {
				if (!user || !user.client_id) return
				
				uni.showModal({
					title: '提示',
					content: '确定要取消关注吗？',
					success: async (res) => {
						if (res.confirm) {
							try {
								// 显示加载提示
								uni.showLoading({
									title: '处理中...',
									mask: true
								})
								
								// 调用取消关注接口
								const result = await unfollowUser({
									publisher_id: user.client_id
								})
								
								// 关闭加载提示
								uni.hideLoading()
								
								if (result && result.code === 0) {
									// 成功后从列表中移除
									this.userList.splice(index, 1)
									this.total--
									
									uni.showToast({
										title: '已取消关注',
										icon: 'success'
									})
								} else {
									throw new Error(result?.message || '操作失败')
								}
							} catch (error) {
								uni.hideLoading()
								console.error('取消关注失败:', error)
								uni.showToast({
									title: error.message || '操作失败，请重试',
									icon: 'none'
								})
							}
						}
					}
				})
			},
			
			// 跳转到发现页
			handleExplore() {
				uni.switchTab({
					url: '/pages/index/index'
				})
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background-color: #ffffff;
		display: flex;
		flex-direction: column;
	}
	
	.content-scroll {
		flex: 1;
		height: 100vh;
		box-sizing: border-box;
	}
	
	.follow-list {
		padding: 0;
	}
	
	.user-item {
		display: flex;
		align-items: center;
		padding: 24rpx 30rpx;
		border-bottom: 1rpx solid #f5f5f5;
	}
	
	.user-avatar {
		width: 80rpx;
		height: 80rpx;
		border-radius: 50%;
		margin-right: 20rpx;
		background-color: #f5f5f5;
	}
	
	.user-info {
		flex: 1;
		display: flex;
		flex-direction: column;
	}
	
	.user-name {
		font-size: 30rpx;
		color: #333333;
		font-weight: 500;
		margin-bottom: 6rpx;
	}
	
	.user-meta {
		display: flex;
		align-items: center;
	}
	
	.publish-count {
		font-size: 24rpx;
		color: #999999;
		margin-right: 16rpx;
	}
	
	.follow-btn {
		min-width: 100rpx;
		height: 56rpx;
		font-size: 24rpx;
		border-radius: 28rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0;
		padding: 0 20rpx;
		background: linear-gradient(135deg, #fc3e2b 0%, #fa7154 100%);
		color: #ffffff;
		border: none;
		transition: all 0.2s;
	}
	
	.follow-btn:active {
		opacity: 0.8;
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
	
	.action-btn {
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
	
	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 60rpx 0;
	}
	
	.loading-indicator {
		width: 40rpx;
		height: 40rpx;
		border: 4rpx solid #f0f0f0;
		border-radius: 50%;
		border-top-color: #fc3e2b;
		animation: spin 1s linear infinite;
		margin-bottom: 20rpx;
	}
	
	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}
	
	.loading-text {
		font-size: 26rpx;
		color: #999999;
	}
	
	.loading-more {
		text-align: center;
		padding: 20rpx 0;
	}
</style> 