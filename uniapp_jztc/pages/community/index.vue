<template>
	<view class="community-container">
		<!-- 自定义顶部栏 -->
		<view class="custom-nav" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ 
				height: navBarHeight + 'px',
				paddingRight: rightPadding + 'px'
			}">
				<text class="page-title">闲置</text>
				<view class="search-box" @click="handleSearch" :style="{ width: searchBoxWidth + 'px' }">
					<uni-icons type="search" size="16" color="#666666"></uni-icons>
					<text class="placeholder-text">搜索闲置商品</text>
				</view>
			</view>
		</view>
		
		<!-- 内容区域 -->
		<view class="page-content" :style="{ paddingTop: navigationBarHeight + 'px' }">
			<!-- 菜单导航 -->
			<view class="menu-nav">
				<scroll-view 
					class="menu-scroll" 
					scroll-x="true" 
					:show-scrollbar="false"
				>
					<view class="menu-list">
						<view 
							v-for="(item, index) in menuList" 
							:key="index"
							class="menu-item"
							:class="{ active: currentMenu === index }"
							@click="switchMenu(index)"
						>
							<text class="menu-text">{{ item.name }}</text>
						</view>
					</view>
				</scroll-view>
			</view>
			
			<!-- 修改滚动区域结构 -->
			<scroll-view 
				class="content-scroll" 
				scroll-y="true"
				:enhanced="true"
				:bounces="true"
				:show-scrollbar="false"
				:refresher-enabled="true"
				:refresher-triggered="isRefreshing"
				@refresherrefresh="onRefresh"
				@refresherrestore="onRestore"
				@scrolltolower="handleScrollToLower"
			>
				<view class="content-wrapper">
					<community-news ref="newsComponent"></community-news>
				</view>
			</scroll-view>
		</view>
		
		<!-- 底部导航栏 -->
		<tab-bar :current-tab="tabIndex"></tab-bar>
	</view>
</template>

<script>
	import TabBar from '@/components/tab-bar/index.vue'
	import CommunityNews from '@/components/community-news/index.vue'
	import deviceMixin from '@/mixins/device.js'
	import shareMixin from '@/mixins/share.js'
	import { get } from '@/utils/request.js'
	import { addBrowseRecord } from '@/apis/content.js'
	import { getShareSettings } from '@/utils/share.js'
	
	export default {
		components: {
			TabBar,
			CommunityNews
		},
		mixins: [deviceMixin, shareMixin],
		data() {
			return {
				tabIndex: 1,
				isRefreshing: false,
				currentMenu: 0,
				searchBoxWidth: 200,
				rightPadding: 0,
				menuList: [],
				menuLoading: false,
				isHomePage: true, // 标记为首页类型，用于分享功能
				shareData: null
			}
		},
		onLoad() {
			// 获取闲置物品分类
			this.fetchCategories()
		},
		mounted() {
			// 计算搜索框宽度和右侧边距
			const menuButtonInfo = uni.getMenuButtonBoundingClientRect()
			const windowInfo = uni.getWindowInfo()
			
			// 计算搜索框宽度，减少标题占用的空间
			this.searchBoxWidth = menuButtonInfo.left - 100
			
			// 设置右侧边距，与胶囊对齐
			this.rightPadding = windowInfo.windowWidth - menuButtonInfo.left
		},
		onShow() {
			this.tabIndex = 1
			
			// 更新分享数据
			this.updateShareData()
		},
		methods: {
			// 获取闲置分类数据
			fetchCategories() {
				this.menuLoading = true
				
				// 添加默认全部选项
				const defaultMenu = {
					id: 0,
					name: '全部'
				}
				
				// 添加超时处理
				const fetchTimeout = setTimeout(() => {
					if (this.menuLoading) {
						// 超时后使用默认菜单并关闭加载状态
						this.useDefaultMenu()
						this.menuLoading = false
						console.log('获取闲置分类超时，使用默认值')
					}
				}, 3000)
				
				// 调用接口获取分类数据
				get('/wx/client/content/categories', { type: 2 })
					.then(res => {
						clearTimeout(fetchTimeout)
						if (res.code === 0) {
							// 处理返回数据
							const categoryList = res.data.list || []
							
							// 将全部选项添加到列表最前面
							this.menuList = [defaultMenu, ...categoryList]
							
							console.log('成功获取闲置分类:', this.menuList)
							
							// 更新分享数据
							this.updateShareData()
						} else {
							console.error('获取闲置分类失败:', res.message || '未知错误')
							// 加载失败时使用默认菜单
							this.useDefaultMenu()
						}
					})
					.catch(err => {
						clearTimeout(fetchTimeout)
						console.error('请求闲置分类接口出错:', err)
						// 加载失败时使用默认菜单
						this.useDefaultMenu()
					})
					.finally(() => {
						clearTimeout(fetchTimeout)
						this.menuLoading = false
					})
			},
			
			// 接口失败时使用默认菜单
			useDefaultMenu() {
				// 使用空数组，只保留"全部"选项
				this.menuList = [
					{ id: 0, name: '全部' }
				]
			},
			
			handleSearch() {
				uni.navigateTo({
					url: '/pages/search/index'
				})
			},
			switchMenu(index) {
				this.currentMenu = index
				
				// 获取选中分类的ID
				const selectedCategory = this.menuList[index]
				const categoryId = selectedCategory ? selectedCategory.id : 0
				
				// 发送分类变更事件，通知内容组件更新
				uni.$emit('idleCategoryChanged', { categoryId: categoryId })
				
				console.log('切换到分类:', selectedCategory.name, 'ID:', categoryId)
			},
			async onRefresh() {
				if (this.isRefreshing) return
				this.isRefreshing = true
				
				try {
					// 刷新分类数据
					this.fetchCategories()
					
					// 刷新内容组件数据
					if (this.$refs.newsComponent) {
						this.$refs.newsComponent.resetAndLoad()
					}
					
					await this.refreshData()
					
					// 更新分享数据
					this.updateShareData()
				} catch (error) {
					console.error('刷新失败:', error)
				} finally {
					this.isRefreshing = false
				}
			},
			async refreshData() {
				await new Promise(resolve => setTimeout(resolve, 1000))
			},
			// 处理上拉加载更多
			handleScrollToLower() {
				console.log('触发上拉加载更多')
				// 调用内容组件的加载更多方法
				if (this.$refs.newsComponent) {
					this.$refs.newsComponent.loadMore()
				}
			},
			// 更新分享数据
			async updateShareData() {
				try {
					// 获取分享配置
					const settings = await getShareSettings();
					
					if (settings) {
						// 使用闲置社区专用的分享设置
						const shareTitle = settings.content_share_text || '查看闲置社区';
						const shareImage = settings.content_share_image || '';
						
						// 直接设置分享数据
						this.shareData = {
							title: shareTitle,
							imageUrl: shareImage,
							path: '/pages/community/index'
						};
						
						console.log('闲置社区分享数据已更新:', this.shareData);
					} else {
						// 配置获取失败，调用混入的默认初始化
						this.initShareData();
					}
				} catch (error) {
					console.error('更新分享数据失败:', error);
					// 配置获取失败，调用混入的默认初始化
					this.initShareData();
				}
			},
			// 显示分享菜单
			handleShare() {
				// 更新分享数据
				this.updateShareData();
				// 显示分享菜单，包括分享朋友圈
				uni.showShareMenu({
					withShareTicket: true,
					menus: ['shareAppMessage', 'shareTimeline']
				});
			}
		}
	}
</script>

<style>
	/* 添加全局变量 */
	page {
		--menu-height: 88rpx;
		--tab-bar-height: 100rpx;
		--safe-bottom: env(safe-area-inset-bottom);
	}

	.community-container {
		display: flex;
		flex-direction: column;
		height: 100vh;
		background-color: #f5f5f5;
		overflow: hidden;
	}
	
	/* 自定义导航栏样式 */
	.custom-nav {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 100;
		background: linear-gradient(180deg, #fc3e2b 0%, #fa7154 101%);
	}
	
	.nav-content {
		display: flex;
		align-items: center;
		padding-left: 24rpx;
		gap: 12rpx;
	}
	
	.page-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #ffffff;
		flex-shrink: 0;
		margin-right: 12rpx;
		width: 64rpx;
	}
	
	.search-box {
		height: 64rpx;
		background-color: rgba(255, 255, 255, 0.9);
		border-radius: 32rpx;
		display: flex;
		align-items: center;
		padding: 0 24rpx;
		gap: 8rpx;
	}
	
	.placeholder-text {
		font-size: 26rpx;
		color: #999999;
	}
	
	.page-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		position: relative;
		z-index: 1;
		margin-top: -20rpx;
		background-color: #fa7154;
		/* 确保内容区域高度正确 */
		height: calc(100vh - var(--window-top));
	}
	
	/* 移除之前的渐变背景 */
	.page-content::before {
		display: none;
	}
	
	/* 菜单导航样式 */
	.menu-nav {
		background-color: #ffffff;
		border-top-left-radius: 30rpx;
		border-top-right-radius: 30rpx;
		width: 100%;
		position: relative;
		z-index: 2;
		margin-top: 20rpx;
	}
	
	.menu-scroll {
		width: 100%;
		white-space: nowrap;
	}
	
	.menu-list {
		display: flex;
		padding: 0 16rpx;
		height: 88rpx;
	}
	
	.menu-item {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0 12rpx;
		position: relative;
	}
	
	.menu-text {
		font-size: 28rpx;
		color: #333333;
		line-height: 88rpx;
		padding: 0 4rpx;
	}
	
	.menu-item.active .menu-text {
		color: #fc3e2b;
		font-weight: bold;
	}
	
	.menu-item.active::after {
		content: '';
		position: absolute;
		bottom: 16rpx;
		left: 50%;
		transform: translateX(-50%);
		width: 32rpx;
		height: 4rpx;
		background-color: #fc3e2b;
		border-radius: 2rpx;
	}
	
	.content-scroll {
		flex: 1;
		background-color: #ffffff;
		/* 调整上边距和内边距 */
		margin-top: -20rpx;
		padding-top: 20rpx;
		/* 使用更精确的高度计算 */
		height: calc(100% - var(--menu-height));
		/* 添加溢出处理 */
		overflow: hidden;
	}
	
	.content-wrapper {
		min-height: 100%;
		/* 使用更可靠的底部内边距计算 */
		padding-bottom: calc(var(--safe-bottom) + var(--tab-bar-height) + 40rpx);
	}
	
	.content-list {
		padding: 0;
	}
</style>
