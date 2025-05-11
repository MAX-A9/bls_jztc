<template>
	<view class="content-list">
		<!-- 加载状态 -->
		<view class="loading-container" v-if="loading && contentList.length === 0">
			<view class="loading-indicator"></view>
			<text class="loading-text">加载中...</text>
		</view>
		
		<!-- 内容列表 -->
		<view 
			class="content-item card" 
			v-for="(item, index) in contentList" 
			:key="index"
			@tap="handleItemClick(item)"
		>
			<!-- 左侧内容区域 -->
			<view :class="['content-left', { 'full-width': !item.image }]">
				<view class="tag-row">
					<!-- 置顶标签 -->
					<view class="tag top-tag" v-if="item.isTop">置顶</view>
					<!-- 类目标签 -->
					<view class="tag category-tag">{{ item.category }}</view>
				</view>
				
				<!-- 标题 -->
				<view class="content-title">{{ item.title }}</view>
				
				<!-- 底部信息 -->
				<view class="content-footer">
					<text class="publisher">{{ item.publisher }}</text>
					<text class="publish-time">{{ item.publishTimeFormatted }}</text>
				</view>
			</view>
			
			<!-- 右侧图片区域 -->
			<view class="content-right" v-if="item.image">
				<image :src="item.image" mode="aspectFill" class="content-image"></image>
			</view>
		</view>
		
		<!-- 加载更多提示 -->
		<view class="load-more" v-if="contentList.length > 0">
			<text class="load-more-text" v-if="hasMore && !loading">上拉加载更多</text>
			<text class="load-more-text" v-if="loading && contentList.length > 0">加载中...</text>
			<text class="load-more-text" v-if="!hasMore && contentList.length > 0">没有更多数据了</text>
		</view>
		
		<!-- 无数据提示 -->
		<view class="empty-container" v-if="!loading && contentList.length === 0">
			<image src="/static/images/empty.png" mode="aspectFit" class="empty-image"></image>
			<text class="empty-text">暂无内容</text>
		</view>
	</view>
</template>

<script>
	import { getRegionContentList } from '@/apis/content.js'
	import { formatTimeAgo, getTimestamp } from '@/utils/date.js'
	
	export default {
		name: 'ContentList',
		props: {
			// 可选参数，允许从父组件传入
			categoryId: {
				type: [Number, String],
				default: ''
			}
		},
		data() {
			return {
				contentList: [],
				loading: false,
				hasMore: true,
				page: 1,
				pageSize: 10,
				regionId: 0, // 当前区域ID
				currentCategory: '', // 当前分类
				isInitialized: false // 添加初始化标志
			}
		},
		created() {
			// 从本地存储获取当前的位置ID并初始化
			const currentLocationId = uni.getStorageSync('currentLocationId') || 0;
			if (currentLocationId) {
				this.regionId = currentLocationId;
				console.log('组件创建时，从本地存储初始化区域ID:', this.regionId);
			}
			
			// 监听页面显示，更新区域ID
			uni.$on('locationChanged', this.handleLocationChanged);
			console.log('已注册locationChanged事件监听');
		},
		mounted() {
			// 组件挂载后初始化数据
			this.$nextTick(() => {
				this.init();
			});
		},
		beforeDestroy() {
			// 移除事件监听
			uni.$off('locationChanged', this.handleLocationChanged);
		},
		methods: {
			// 初始化方法，用于组件挂载后或页面显示时调用
			init() {
				// 避免重复初始化
				if (this.isInitialized && this.contentList.length > 0) {
					console.log('内容列表已初始化，跳过重复初始化');
					return;
				}
				
				// 从本地存储获取当前的位置ID
				const currentLocationId = uni.getStorageSync('currentLocationId') || 0;
				
				// 如果没有regionId但有本地存储的值，则使用本地存储的值
				if ((!this.regionId || this.regionId === 0) && currentLocationId) {
					this.regionId = currentLocationId;
					console.log('初始化时，从本地存储更新区域ID:', this.regionId);
				}
				
				// 在有regionId的情况下加载数据
				if (this.regionId && this.regionId !== 0) {
					console.log('初始化时，开始加载内容列表数据...');
					this.resetAndLoad();
					// 标记为已初始化
					this.isInitialized = true;
				} else {
					console.log('初始化时，未找到有效的区域ID');
				}
			},
			// 格式化发布时间
			formatPublishTime(time) {
				return formatTimeAgo(time);
			},
			
			handleItemClick(item) {
				// 跳转到详情页，并传递必要的参数
				uni.navigateTo({
					url: `/pages/content/detail?id=${item.id}&category=${item.category}`,
					success: () => {
						console.log('跳转到详情页:', item)
					},
					fail: (err) => {
						console.error('跳转失败:', err)
					}
				})
			},
			
			// 处理分类切换
			handleCategoryChange(categoryId) {
				this.currentCategory = categoryId;
				// 重置数据并重新加载
				this.resetAndLoad();
			},
			
			// 处理位置变更
			handleLocationChanged(data) {
				console.log('位置变更事件触发，接收到的regionId:', data.regionId);
				const newRegionId = data.regionId || 0;
				
				// 确保regionId有效
				if (!newRegionId) {
					console.error('接收到无效的regionId:', newRegionId);
					return;
				}
				
				// 更新当前地区ID并输出日志
				if (this.regionId !== newRegionId) {
					console.log('地区ID已更新:', this.regionId, '->', newRegionId);
					this.regionId = newRegionId;
					// 重置数据并重新加载
					this.resetAndLoad();
				} else {
					console.log('地区ID未变化，保持当前值:', this.regionId);
				}
			},
			
			// 重置数据并重新加载
			resetAndLoad() {
				console.log('重置内容列表，当前区域ID:', this.regionId);
				this.page = 1;
				this.contentList = [];
				this.hasMore = true;
				this.loadData();
			},
			
			// 加载更多数据
			loadMore() {
				if (this.loading || !this.hasMore) return;
				this.page++;
				this.loadData();
			},
			
			// 加载数据
			async loadData() {
				if (this.loading) return;
				
				try {
					this.loading = true;
					
					// 从本地存储获取当前的位置ID
					const currentLocationId = uni.getStorageSync('currentLocationId') || 0;
					console.log('从本地存储获取的区域ID:', currentLocationId);
					
					// 更新当前地区ID
					if (!this.regionId && currentLocationId) {
						console.log('使用本地存储区域ID更新组件regionId:', currentLocationId);
						this.regionId = currentLocationId;
					}
					
					// 确保regionId有效，如果无效则提示用户选择地区
					if (!this.regionId) {
						uni.showToast({
							title: '请先选择地区',
							icon: 'none'
						});
						throw new Error('未选择地区');
					}
					
					// 构建参数
					const params = {
						regionId: this.regionId,
						page: this.page,
						pageSize: this.pageSize
					};
					
					// 如果有分类ID，添加到请求参数
					if (this.currentCategory) {
						params.category = this.currentCategory;
					}
					
					console.log('加载内容列表，参数:', params);
					
					// 调用API获取数据
					const result = await getRegionContentList(params);
					
					if (result.code === 0 && result.data) {
						// 获取数据成功
						const { list = [], total = 0 } = result.data;
						
						// 处理列表数据，格式化发布时间并进行排序
						let formattedList = list.map(item => {
							// 确保item是一个新对象，避免修改原始数据
							return {
								...item,
								// 格式化发布时间为"多久前"的格式
								publishTimeFormatted: this.formatPublishTime(item.publishTime)
							};
						});
						
						// 对列表进行排序，置顶内容排在前面
						formattedList.sort((a, b) => {
							// 先按照isTop排序（置顶在前）
							if (a.isTop && !b.isTop) return -1;
							if (!a.isTop && b.isTop) return 1;
							
							// 如果都是置顶或都不是置顶，则按照发布时间排序（新的在前）
							let timeA = 0, timeB = 0;
							
							// 预先处理日期格式，避免在排序时出错
							try {
								// 如果是 "yyyy-MM-dd HH:mm:ss" 格式，手动转换为ISO格式
								if (a.publishTime && typeof a.publishTime === 'string' && a.publishTime.includes(' ')) {
									const aDateParts = a.publishTime.split(' ')[0].split('-');
									const aTimeParts = a.publishTime.split(' ')[1].split(':');
									// 使用标准Date构造函数创建日期对象
									const aDate = new Date(
										parseInt(aDateParts[0]), // 年
										parseInt(aDateParts[1]) - 1, // 月(0-11)
										parseInt(aDateParts[2]), // 日
										parseInt(aTimeParts[0]), // 时
										parseInt(aTimeParts[1]), // 分
										parseInt(aTimeParts[2]) // 秒
									);
									timeA = aDate.getTime();
								} else {
									timeA = getTimestamp(a.publishTime);
								}
								
								if (b.publishTime && typeof b.publishTime === 'string' && b.publishTime.includes(' ')) {
									const bDateParts = b.publishTime.split(' ')[0].split('-');
									const bTimeParts = b.publishTime.split(' ')[1].split(':');
									// 使用标准Date构造函数创建日期对象
									const bDate = new Date(
										parseInt(bDateParts[0]), // 年
										parseInt(bDateParts[1]) - 1, // 月(0-11)
										parseInt(bDateParts[2]), // 日
										parseInt(bTimeParts[0]), // 时
										parseInt(bTimeParts[1]), // 分
										parseInt(bTimeParts[2]) // 秒
									);
									timeB = bDate.getTime();
								} else {
									timeB = getTimestamp(b.publishTime);
								}
							} catch (e) {
								console.error('日期解析错误:', e, a.publishTime, b.publishTime);
								// 发生错误时使用当前时间作为默认值
								timeA = timeB = Date.now();
							}
							
							return timeB - timeA; // 降序排列，最新的在前面
						});
						
						// 追加数据到列表
						// 第一页时直接替换，加载更多时追加
						if (this.page === 1) {
							this.contentList = formattedList;
						} else {
							this.contentList = [...this.contentList, ...formattedList];
						}
						
						// 判断是否还有更多数据
						this.hasMore = this.contentList.length < total;
						
						console.log(`加载了${list.length}条数据，总共${total}条，使用的区域ID:${this.regionId}`);
						
						// 通知父组件刷新完成
						this.$parent && this.$parent.refreshing === true && this.notifyRefreshComplete();
					} else {
						throw new Error(result.message || '获取内容列表失败');
					}
				} catch (error) {
					console.error('加载内容列表失败:', error);
					uni.showToast({
						title: '加载失败，请重试',
						icon: 'none',
						duration: 2000
					});
				} finally {
					this.loading = false;
				}
			},
			
			// 通知父组件刷新完成
			notifyRefreshComplete() {
				// 发送事件给父组件
				this.$emit('refreshComplete');
				
				// 直接修改父组件的刷新状态（如果父组件是index页面）
				if (this.$parent && this.$parent.refreshing === true) {
					this.$parent.refreshing = false;
					console.log('内容加载完成，关闭刷新状态');
				}
			},
		}
	}
</script>

<style>
	.content-list {
		padding: 20rpx;
		background-color: #ffffff;
		padding-bottom: 120rpx;
		min-height: 300rpx;
	}
	
	/* 卡片样式 */
	.content-item.card {
		display: flex;
		padding: 24rpx;
		background-color: #ffffff;
		border-radius: 16rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 16rpx rgba(0, 0, 0, 0.06);
	}
	
	.content-left {
		flex: 1;
		margin-right: 20rpx;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
	}
	
	.content-left.full-width {
		margin-right: 0;
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
	
	.top-tag {
		background-color: #fc3e2b;
		color: #ffffff;
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
	
	.content-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: 24rpx;
		color: #999999;
		margin-top: auto; /* 将底部信息推到底部 */
	}
	
	.content-right {
		width: 160rpx;
		height: 160rpx;
		flex-shrink: 0;
	}
	
	.content-image {
		width: 100%;
		height: 100%;
		border-radius: 12rpx; /* 增加图片圆角 */
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
	
	/* 加载更多样式 */
	.load-more {
		text-align: center;
		padding: 30rpx 0;
	}
	
	.load-more-text {
		font-size: 26rpx;
		color: #999999;
	}
	
	/* 空数据样式 */
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