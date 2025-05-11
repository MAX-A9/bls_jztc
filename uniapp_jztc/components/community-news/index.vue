<template>
	<view class="news-container">
		<!-- 加载状态 -->
		<view class="loading-container" v-if="loading && newsList.length === 0">
			<view class="loading-indicator"></view>
			<text class="loading-text">加载中...</text>
		</view>
		
		<!-- 闲置列表 -->
		<view class="news-list" v-if="!loading || newsList.length > 0">
			<view 
				class="news-item" 
				v-for="(item, index) in newsList" 
				:key="index"
				@click="handleNewsClick(item)"
			>
				<image 
					:src="item.image || defaultImage" 
					mode="aspectFill" 
					class="news-image"
					@error="handleImageError(index)"
				></image>
				<view class="news-content">
					<text class="news-text">{{ item.title }}</text>
					<text class="news-desc">{{ item.desc }}</text>
					<view class="location">
						<text class="location-text">{{ item.location }}</text>
					</view>
					<view class="price-want">
						<text class="price">¥{{ item.price }}</text>
						<text class="want-count">{{ item.wantCount }}人想要</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 加载更多提示 -->
		<view class="load-more" v-if="newsList.length > 0">
			<text class="load-more-text" v-if="hasMore && !loading">上拉加载更多</text>
			<text class="load-more-text" v-if="loading && newsList.length > 0">加载中...</text>
			<text class="load-more-text" v-if="!hasMore && newsList.length > 0">没有更多数据了</text>
		</view>
		
		<!-- 无数据提示 -->
		<view class="empty-container" v-if="!loading && newsList.length === 0">
			<image src="/static/images/empty.png" mode="aspectFit" class="empty-image"></image>
			<text class="empty-text">暂无闲置物品</text>
		</view>
	</view>
</template>

<script>
	import { get } from '@/utils/request.js'
	import { addBrowseRecord, getRegionIdleList } from '@/apis/content.js'
	
	export default {
		name: 'CommunityNews',
		props: {
			categoryId: {
				type: [Number, String],
				default: 0
			}
		},
		data() {
			return {
				newsList: [],
				loading: false,
				hasMore: true,
				page: 1,
				pageSize: 10,
				regionId: 0, // 当前区域ID
				currentCategoryId: 0, // 当前选中的分类ID
				defaultImage: '/static/images/default-idle.png' // 默认图片路径，使用统一的默认图片
			}
		},
		created() {
			// 初始化currentCategoryId为props中的值
			this.currentCategoryId = this.categoryId
			
			// 从本地存储获取当前的位置ID并初始化
			const currentLocationId = uni.getStorageSync('currentLocationId') || 0
			if (currentLocationId) {
				this.regionId = currentLocationId
				console.log('闲置列表组件创建时，从本地存储初始化区域ID:', this.regionId)
			}
			
			// 监听分类变更事件
			uni.$on('idleCategoryChanged', this.handleCategoryChanged)
			
			// 监听位置变更事件
			uni.$on('locationChanged', this.handleLocationChanged)
			
			// 初始加载数据
			this.loadData()
		},
		beforeDestroy() {
			// 移除事件监听
			uni.$off('idleCategoryChanged', this.handleCategoryChanged)
			uni.$off('locationChanged', this.handleLocationChanged)
		},
		methods: {
			// 处理分类变更
			handleCategoryChanged(data) {
				console.log('接收到分类变更事件:', data)
				if (data && data.categoryId !== undefined) {
					// 更新当前分类ID并重置列表
					this.currentCategoryId = data.categoryId
					this.resetAndLoad()
				}
			},
			
			// 处理位置变更
			handleLocationChanged(data) {
				console.log('位置变更事件触发，接收到的regionId:', data.regionId)
				const newRegionId = data.regionId || 0
				
				// 确保regionId有效
				if (!newRegionId) {
					console.error('接收到无效的regionId:', newRegionId)
					return
				}
				
				// 更新当前地区ID并输出日志
				if (this.regionId !== newRegionId) {
					console.log('地区ID已更新:', this.regionId, '->', newRegionId)
					this.regionId = newRegionId
					// 重置数据并重新加载
					this.resetAndLoad()
				} else {
					console.log('地区ID未变化，保持当前值:', this.regionId)
				}
			},
			
			// 重置列表并重新加载数据
			resetAndLoad() {
				this.page = 1
				this.newsList = []
				this.hasMore = true
				this.loadData()
			},
			
			// 加载闲置物品数据
			async loadData() {
				if (this.loading) return
				
				try {
					this.loading = true
					
					// 从本地存储获取当前的位置ID
					const currentLocationId = uni.getStorageSync('currentLocationId') || 0
					
					// 更新当前地区ID
					if (!this.regionId && currentLocationId) {
						this.regionId = currentLocationId
					}
					
					// 构建请求参数
					const params = {
						regionId: this.regionId,
						page: this.page, // 后端API要求页码最小值为1
						pageSize: this.pageSize,
						category: this.currentCategoryId // 使用当前分类ID
					}
					
					console.log('加载闲置物品列表，参数:', params)
					
					// 调用API获取数据
					const result = await getRegionIdleList(params)
					
					if (result.code === 0 && result.data) {
						// 获取数据成功
						const { list = [], total = 0 } = result.data
						
						// 根据接口返回数据结构处理数据
						const formattedList = list.map(item => {
							return {
								id: item.id,
								title: item.title || '',
								desc: item.summary || '暂无描述',
								location: item.tradePlace || '未知地点',
								price: item.price || 0,
								image: item.image || this.defaultImage,
								wantCount: item.likes || 0
							}
						})
						
						// 第一页时直接替换，加载更多时追加
						if (this.page === 1) {
							this.newsList = formattedList
						} else {
							this.newsList = [...this.newsList, ...formattedList]
						}
						
						// 判断是否还有更多数据
						this.hasMore = this.newsList.length < total
						
						console.log(`加载了${list.length}条闲置物品数据，总共${total}条`)
						
						// 通知父组件刷新完成
						this.$parent && this.$parent.isRefreshing === true && this.notifyRefreshComplete()
					} else {
						throw new Error(result.message || '获取闲置物品列表失败')
					}
				} catch (error) {
					console.error('加载闲置物品列表失败:', error)
					
					// 如果是第一页并且出错，显示默认数据
					if (this.page === 1 && this.newsList.length === 0) {
						console.log('使用默认数据')
						this.useDefaultData()
					}
					
					uni.showToast({
						title: '加载失败，请重试',
						icon: 'none',
						duration: 2000
					})
				} finally {
					this.loading = false
				}
			},
			
			// 通知父组件刷新完成
			notifyRefreshComplete() {
				// 如果父组件是index页面且正在刷新，关闭刷新状态
				if (this.$parent && this.$parent.isRefreshing === true) {
					this.$parent.isRefreshing = false
					console.log('内容加载完成，关闭刷新状态')
				}
			},
			
			// 加载更多数据
			loadMore() {
				if (this.loading || !this.hasMore) return
				this.page++
				this.loadData()
			},
			
			// 图片加载错误处理
			handleImageError(index) {
				// 替换为默认图片
				this.newsList[index].image = this.defaultImage
			},
			
			// 使用默认数据（接口失败时的备选方案，实际应用中应删除）
			useDefaultData() {
				this.newsList = []
				this.hasMore = false
			},
			
			handleNewsClick(news) {
				// 记录浏览历史
				this.recordBrowseHistory(news.id)
				
				// 跳转到详情页
				uni.navigateTo({
					url: `/pages/community/detail?id=${news.id}`,
					success: () => {
						console.log('跳转到闲置详情页:', news)
					},
					fail: (err) => {
						console.error('跳转失败:', err)
					}
				})
			},
			
			// 记录浏览历史
			async recordBrowseHistory(contentId) {
				if (!contentId) {
					console.error('记录浏览历史失败: 内容ID为空')
					return
				}
				
				try {
					// 使用idle类型记录闲置物品浏览
					const result = await addBrowseRecord(contentId, 'idle')
					if (result.code !== 0) {
						// 在开发环境中显示错误信息，方便调试
						if (process.env.NODE_ENV === 'development') {
							console.error('记录闲置浏览历史失败:', result.message)
						}
					}
				} catch (error) {
					console.error('记录闲置浏览历史异常:', error)
				}
			}
		}
	}
</script>

<style>
	.news-container {
		background-color: #ffffff;
		padding: 20rpx 20rpx 120rpx;
		height: 100%;
	}
	
	.news-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;
		padding: 0 10rpx;
	}
	
	.news-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333333;
	}
	
	.news-list {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 20rpx;
		padding: 0 10rpx;
	}
	
	.news-item {
		display: flex;
		flex-direction: column;
		background: #ffffff;
		border-radius: 12rpx;
		overflow: hidden;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
	}
	
	.news-image {
		width: 100%;
		height: 320rpx;
		object-fit: cover;
	}
	
	.news-content {
		padding: 16rpx;
		display: flex;
		flex-direction: column;
		gap: 8rpx;
	}
	
	.news-text {
		font-size: 28rpx;
		color: #333333;
		font-weight: bold;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.news-desc {
		font-size: 24rpx;
		color: #666666;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		margin-bottom: 4rpx;
	}
	
	.location {
		overflow: hidden;
		margin-bottom: 8rpx;
	}
	
	.location-text {
		font-size: 22rpx;
		color: #999999;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.price-want {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-top: 4rpx;
	}
	
	.price {
		font-size: 30rpx;
		color: #fc3e2b;
		font-weight: bold;
	}
	
	.want-count {
		font-size: 24rpx;
		color: #999999;
	}
	
	/* 加载中样式 */
	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 300rpx;
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
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}
	
	.loading-text {
		font-size: 26rpx;
		color: #999999;
	}
	
	/* 加载更多提示 */
	.load-more {
		text-align: center;
		padding: 30rpx 0;
	}
	
	.load-more-text {
		font-size: 26rpx;
		color: #999999;
	}
	
	/* 空数据提示 */
	.empty-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 100rpx 0;
	}
	
	.empty-image {
		width: 200rpx;
		height: 200rpx;
		margin-bottom: 30rpx;
	}
	
	.empty-text {
		font-size: 28rpx;
		color: #999999;
	}
</style> 