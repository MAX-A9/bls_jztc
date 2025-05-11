<template>
	<view class="container">
		<!-- 选项卡 -->
		<view class="tab-bar">
			<view 
				class="tab-item" 
				v-for="(item, index) in tabList" 
				:key="index"
				:class="{ active: currentTab === index }"
				@tap="switchTab(index)"
			>
				<text class="tab-text">{{ item }}</text>
			</view>
		</view>
		
		<!-- 占位元素，为固定的选项卡腾出空间 -->
		<view class="tab-placeholder"></view>
		
		<scroll-view 
			class="content-scroll" 
			scroll-y 
			:refresher-enabled="true"
			:refresher-triggered="isRefreshing"
			@refresherrefresh="onRefresh"
			:show-scrollbar="false"
			:enhanced="true"
			:bounces="true"
			@scrolltolower="loadMore"
		>
			<!-- 闲置收藏列表 -->
			<view class="favorite-list" v-if="currentTab === 0 && favoriteList.length > 0">
				<view 
					class="favorite-item"
					v-for="(item, index) in favoriteList"
					:key="item.id"
					@tap="handleItemClick(item)"
				>
					<image class="item-image" :src="item.image" mode="aspectFill"></image>
					<view class="item-content">
						<text class="item-title">{{ item.title }}</text>
						<text class="item-price" v-if="item.price > 0">¥{{ item.price }}</text>
						<view class="item-footer">
							<view class="user-info">
								<text class="publisher-name">{{ item.publisher }}</text>
							</view>
							<view class="favorite-btn" @tap.stop="toggleFavorite(item.id, 2, index)">
								<uni-icons type="star-filled" size="20" color="#fc3e2b"></uni-icons>
							</view>
						</view>
					</view>
				</view>
				<!-- 加载更多 -->
				<view class="loading-more" v-if="hasMore && !loading">
					<text class="loading-text">上拉加载更多</text>
				</view>
				<view class="loading-more" v-if="loading">
					<text class="loading-text">加载中...</text>
				</view>
				<view class="loading-more" v-if="!hasMore && favoriteList.length > 0">
					<text class="loading-text">没有更多数据了</text>
				</view>
			</view>
			
			<!-- 帖子收藏列表 -->
			<view class="favorite-list" v-if="currentTab === 1 && postList.length > 0">
				<view 
					class="favorite-item"
					v-for="(item, index) in postList"
					:key="item.id"
					@tap="handleItemClick(item)"
				>
					<image class="item-image" :src="item.image" mode="aspectFill" v-if="item.image"></image>
					<view class="item-content" :class="{'full-width': !item.image}">
						<text class="item-title">{{ item.title }}</text>
						<view class="item-footer">
							<view class="user-info">
								<text class="publisher-name">{{ item.publisher }}</text>
							</view>
							<view class="favorite-btn" @tap.stop="toggleFavorite(item.id, 1, index)">
								<uni-icons type="star-filled" size="20" color="#fc3e2b"></uni-icons>
							</view>
						</view>
					</view>
				</view>
				<!-- 加载更多 -->
				<view class="loading-more" v-if="hasMore && !loading">
					<text class="loading-text">上拉加载更多</text>
				</view>
				<view class="loading-more" v-if="loading">
					<text class="loading-text">加载中...</text>
				</view>
				<view class="loading-more" v-if="!hasMore && postList.length > 0">
					<text class="loading-text">没有更多数据了</text>
				</view>
			</view>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="(currentTab === 0 && favoriteList.length === 0 && !loading) || (currentTab === 1 && postList.length === 0 && !loading)">
				<image class="empty-image" src="/static/images/empty.png" mode="aspectFit"></image>
				<text class="empty-text">暂无收藏内容</text>
				<button class="action-btn" @tap="handleExplore">去发现</button>
			</view>
		</scroll-view>
	</view>
</template>

<script>
	import { getFavoriteList, cancelFavorite } from '@/apis/content.js'
	
	export default {
		data() {
			return {
				currentTab: 0,
				tabList: ['闲置', '帖子'],
				isRefreshing: false,
				loading: false,
				favoriteList: [],
				postList: [],
				favoriteParams: {
					page: 1,
					pageSize: 10,
					type: 2
				},
				postParams: {
					page: 1,
					pageSize: 10,
					type: 1
				},
				hasMore: true,
				favoriteTotal: 0,
				postTotal: 0
			}
		},
		onLoad() {
			// 加载收藏数据
			this.loadFavoriteData()
		},
		methods: {
			// 切换选项卡
			switchTab(index) {
				if (this.currentTab === index) return
				this.currentTab = index
				// 切换到闲置选项卡
				if (index === 0 && this.favoriteList.length === 0) {
					this.loadFavoriteData()
				} 
				// 切换到帖子选项卡
				else if (index === 1 && this.postList.length === 0) {
					this.loadPostData()
				}
			},
			
			// 加载闲置收藏数据
			async loadFavoriteData() {
				if (this.loading) return
				
				this.loading = true
				try {
					const result = await getFavoriteList({
						page: this.favoriteParams.page,
						pageSize: this.favoriteParams.pageSize,
						type: this.favoriteParams.type
					})
					
					if (result.code === 0 && result.data) {
						// 第一页时直接赋值，加载更多时追加
						if (this.favoriteParams.page === 1) {
							this.favoriteList = result.data.list || []
						} else {
							this.favoriteList = [...this.favoriteList, ...(result.data.list || [])]
						}
						
						// 更新总数和判断是否有更多
						this.favoriteTotal = result.data.total || 0
						this.hasMore = this.favoriteList.length < this.favoriteTotal
					} else {
						uni.showToast({
							title: result.message || '加载失败',
							icon: 'none'
						})
					}
				} catch (error) {
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
			
			// 加载帖子收藏数据
			async loadPostData() {
				if (this.loading) return
				
				this.loading = true
				try {
					const result = await getFavoriteList({
						page: this.postParams.page,
						pageSize: this.postParams.pageSize,
						type: this.postParams.type
					})
					
					if (result.code === 0 && result.data) {
						// 第一页时直接赋值，加载更多时追加
						if (this.postParams.page === 1) {
							this.postList = result.data.list || []
						} else {
							this.postList = [...this.postList, ...(result.data.list || [])]
						}
						
						// 更新总数和判断是否有更多
						this.postTotal = result.data.total || 0
						this.hasMore = this.postList.length < this.postTotal
					} else {
						uni.showToast({
							title: result.message || '加载失败',
							icon: 'none'
						})
					}
				} catch (error) {
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
				
				// 重置页码
				if (this.currentTab === 0) {
					this.favoriteParams.page = 1
					await this.loadFavoriteData()
				} else {
					this.postParams.page = 1
					await this.loadPostData()
				}
			},
			
			// 加载更多
			loadMore() {
				if (this.loading || !this.hasMore) return
				
				// 增加页码
				if (this.currentTab === 0) {
					this.favoriteParams.page++
					this.loadFavoriteData()
				} else {
					this.postParams.page++
					this.loadPostData()
				}
			},
			
			// 点击收藏项
			handleItemClick(item) {
				// 根据内容类型跳转到不同页面
				if (this.currentTab === 0 || item.type === 2) {
					// 闲置物品详情
					uni.navigateTo({
						url: `/pages/community/detail?id=${item.id}`
					})
				} else {
					// 帖子内容详情
					uni.navigateTo({
						url: `/pages/content/detail?id=${item.id}`
					})
				}
			},
			
			// 取消收藏
			async toggleFavorite(contentId, type, index) {
				uni.showModal({
					title: '提示',
					content: '确定要取消收藏吗？',
					success: async (res) => {
						if (res.confirm) {
							try {
								const result = await cancelFavorite(contentId)
								
								if (result.code === 0) {
									// 从列表中移除
									if (type === 1) {
										this.postList.splice(index, 1)
									} else {
										this.favoriteList.splice(index, 1)
									}
									
									uni.showToast({
										title: '已取消收藏',
										icon: 'success'
									})
								} else {
									uni.showToast({
										title: result.message || '操作失败',
										icon: 'none'
									})
								}
							} catch (error) {
								uni.showToast({
									title: '网络异常，请稍后重试',
									icon: 'none'
								})
							}
						}
					}
				})
			},
			
			// 去发现
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
		background-color: #f5f5f5;
		display: flex;
		flex-direction: column;
		height: 100vh;
		overflow: hidden;
	}
	
	.tab-bar {
		display: flex;
		background-color: #ffffff;
		height: 88rpx;
		border-bottom: 1rpx solid #f0f0f0;
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 10;
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
	
	.tab-placeholder {
		height: 88rpx;
		flex-shrink: 0;
	}
	
	.content-scroll {
		flex: 1;
		height: calc(100vh - 88rpx);
		box-sizing: border-box;
		-webkit-overflow-scrolling: touch;
	}
	
	.favorite-list {
		padding: 20rpx;
	}
	
	.favorite-item {
		display: flex;
		background-color: #ffffff;
		border-radius: 12rpx;
		padding: 20rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
	}
	
	.item-image {
		width: 200rpx;
		height: 200rpx;
		border-radius: 8rpx;
		margin-right: 20rpx;
	}
	
	.item-content {
		flex: 1;
		display: flex;
		flex-direction: column;
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
		overflow: hidden;
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
	
	.user-info {
		display: flex;
		align-items: center;
	}
	
	.user-avatar {
		width: 40rpx;
		height: 40rpx;
		border-radius: 50%;
		margin-right: 10rpx;
	}
	
	.user-name {
		font-size: 24rpx;
		color: #999999;
	}
	
	.favorite-btn {
		width: 60rpx;
		height: 60rpx;
		display: flex;
		align-items: center;
		justify-content: center;
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
	
	.category-tag {
		font-size: 24rpx;
		color: #666666;
		background-color: #f5f5f5;
		padding: 4rpx 12rpx;
		border-radius: 4rpx;
		display: inline-block;
		margin-bottom: 16rpx;
	}
	
	.publisher-name {
		font-size: 24rpx;
		color: #999999;
	}
	
	.loading-more {
		text-align: center;
		padding: 20rpx 0;
	}
	
	.loading-text {
		font-size: 24rpx;
		color: #999999;
	}
</style> 