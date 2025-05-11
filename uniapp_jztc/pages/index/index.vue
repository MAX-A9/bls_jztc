<template>
	<view class="container">
		<!-- 骨架屏，在页面加载时显示 -->
		<view class="skeleton-screen" v-if="pageLoading">
			<view class="skeleton-header">
				<view class="skeleton-nav"></view>
				<view class="skeleton-search"></view>
			</view>
			<view class="skeleton-content">
				<view class="skeleton-item" v-for="i in 4" :key="i"></view>
			</view>
		</view>
		
		<!-- 固定顶部区域 -->
		<view class="fixed-header">
			<!-- 渐变背景区域 -->
			<view class="gradient-bg">
				<!-- 自定义状态栏 -->
				<view class="custom-nav" :style="{ paddingTop: statusBarHeight + 'px' }">
					<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
						<view class="nav-logo" :style="{ height: isIOS ? '28px' : '32px' }">
							<image :src="logoUrl" mode="heightFix" class="logo-img" />
						</view>
						<view class="nav-center">
							<view class="nav-location" @tap="navigateToLocationSelect">
								<text class="location-text">{{ currentLocation || '定位中...' }}</text>
								<uni-icons type="bottom" size="12" color="#ffffff"></uni-icons>
							</view>
						</view>
						<view class="nav-right"></view>
					</view>
				</view>
				
				<!-- 搜索框/菜单导航切换区域 -->
				<view class="menu-wrapper">
					<!-- 搜索框，当未滚动到指定位置时显示 -->
					<view class="search-box" v-show="!showFixedMenu">
						<view class="search-input">
							<uni-icons type="search" size="16" color="#666666"></uni-icons>
							<input 
								type="text" 
								placeholder="搜索" 
								placeholder-class="placeholder"
								confirm-type="search"
								@confirm="handleSearch"
							/>
						</view>
					</view>
					
					<!-- 菜单导航，当滚动到内容区域时显示 -->
					<scroll-view 
						class="nav-menu fixed-nav-menu scroll-view-container"
						scroll-x="true" 
						:enhanced="true"
						:class="{ 'menu-show': showFixedMenu }"
						v-show="showFixedMenu"
					>
						<view 
							class="menu-item" 
							v-for="(item, index) in menuList" 
							:key="index"
							:class="{ active: currentMenu === index }"
							@click="switchMenu(index)"
						>
							<text class="fixed-menu-text">{{ item.name }}</text>
						</view>
						<view v-if="menuLoading" class="menu-loading">
							<text class="loading-text">加载中...</text>
						</view>
					</scroll-view>
				</view>
			</view>
		</view>
		
		<!-- 使用scroll-view包裹可滚动内容 -->
		<scroll-view 
			class="scroll-content scroll-view-container" 
			scroll-y="true"
			:style="{ paddingTop: headerHeight + 'px' }"
			:enhanced="true"
			:bounces="true"
			@scroll="handleScroll"
			:scroll-top="scrollTop"
			scroll-with-animation
			ref="scrollView"
			@scrolltolower="handleScrollToLower"
			@refresherrefresh="handleRefresh"
			refresher-enabled
			:refresher-triggered="refreshing"
		>
			<!-- 页面内容区域 -->
			<view class="page-content">
				<!-- 功能区域 -->
				<function-area ref="functionArea"></function-area>
				
				<!-- 活动区域 -->
				<activity-area ref="activityArea"></activity-area>
				
				<!-- 轮播图 -->
				<swiper-banner ref="swiperBanner"></swiper-banner>
				
				<!-- 菜单导航 -->
				<view class="category-menu" id="category-menu">
					<scroll-view 
						class="nav-menu scroll-view-container" 
						scroll-x="true" 
						:enhanced="true"
					>
						<view 
							class="menu-item" 
							v-for="(item, index) in menuList" 
							:key="index"
							:class="{ active: currentMenu === index }"
							@click="switchMenu(index)"
						>
							<text class="menu-text">{{ item.name }}</text>
						</view>
						<view v-if="menuLoading" class="menu-loading">
							<text class="loading-text">加载中...</text>
						</view>
					</scroll-view>
				</view>
				
				<!-- 内容列表 -->
				<view class="content-area" id="content-area">
					<content-list ref="contentList"></content-list>
				</view>
			</view>
		</scroll-view>
		
		<!-- 返回顶部按钮 -->
		<view 
			class="back-to-top" 
			v-show="showBackToTop"
			@tap="scrollToTop"
		>
			<uni-icons type="top" size="24" color="#ffffff"></uni-icons>
		</view>
		
		<!-- 添加 TabBar 组件 -->
		<tab-bar @tabChange="handleTabChange" @publish="handlePublish"></tab-bar>
		
		<!-- 位置选择引导弹窗 -->
		<view class="location-guide-mask" v-if="showLocationGuide" @tap="closeLocationGuide">
			<view class="location-guide-content" @tap.stop>
				<view class="guide-header">
					<text class="guide-title">选择所在地区</text>
				</view>
				<view class="guide-body">
					<text class="guide-text">请选择您所在的地区以获取更精准的服务</text>
					<text class="guide-current">当前：{{ currentLocation }}</text>
				</view>
				<view class="guide-footer">
					<button class="guide-btn guide-btn-confirm" @tap="goToLocationSelectFromGuide">去选择</button>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import deviceMixin from '@/mixins/device.js'
	import shareMixin from '@/mixins/share.js'
	import FunctionArea from '@/components/function-area/index.vue'
	import SwiperBanner from '@/components/swiper-banner/index.vue'
	import ContentList from '@/components/content-list/index.vue'
	import ActivityArea from '@/components/activity-area/index.vue'
	import TabBar from '@/components/tab-bar/index.vue'
	import { get } from '@/utils/request.js'
	import { settings } from '@/apis/index.js'
	
	export default {
		components: {
			FunctionArea,
			SwiperBanner,
			ContentList,
			ActivityArea,
			TabBar
		},
		mixins: [deviceMixin, shareMixin],
		data() {
			return {
				tabIndex: 0,
				currentLocation: '',
				isLocating: false,
				menuList: [], // 改为空数组，通过API获取
				menuLoading: false, // 添加加载状态
				currentMenu: 0,
				headerHeight: 0,
				showFixedMenu: false, // 是否显示固定菜单
				menuTop: 0, // 菜单导航的顶部位置
				contentTop: 0, // 内容区域顶部位置
				showBackToTop: false, // 是否显示返回顶部按钮
				scrollTop: 0, // 用于控制滚动位置
				oldScrollTop: 0, // 添加这个变量用于强制更新 scroll-top
				showLocationGuide: false, // 是否显示位置引导弹窗
				pageLoading: true, // 页面加载状态
				refreshing: false, // 添加下拉刷新状态
				logoUrl: '', // 默认logo地址
				contentListInitialized: false, // 标记内容列表是否已初始化
				isHomePage: true, // 标记为首页，用于分享功能
			}
		},
		onLoad() {
			// 设置页面加载超时处理
			setTimeout(() => {
				if (this.pageLoading) {
					this.pageLoading = false;
					console.log('页面加载超时，强制显示页面');
				}
			}, 3000);
			
			// 优先初始化位置，这是必须的
			this.initLocation();
			
			// 使用延迟加载策略
			this.staggeredLoading();
		},
		mounted() {
			this.calculateHeaderHeight()
			// 延迟获取位置信息，确保页面完全渲染
			setTimeout(() => {
				this.getMenuPosition()
			}, 300)
		},
		onReady() {
			// 页面渲染完成后重新计算位置
			this.getMenuPosition()
		},
		onResize() {
			// 屏幕尺寸变化时重新计算位置
			this.calculateHeaderHeight()
			this.getMenuPosition()
		},
		onShow() {
			// 设置当前页面的导航索引
			this.tabIndex = 0;
			
			// 获取并更新最新的位置信息
			const savedLocation = uni.getStorageSync('currentLocation');
			if (savedLocation) {
				this.currentLocation = savedLocation;
				console.log('页面显示时更新位置信息:', savedLocation);
			}
			
			// 初始化内容列表组件（如果尚未初始化）
			this.$nextTick(() => {
				if (this.$refs.contentList && !this.contentListInitialized) {
					console.log('页面显示时，初始化内容列表组件');
					this.$refs.contentList.init();
					this.contentListInitialized = true;
				} else if (this.$refs.contentList) {
					// 已初始化过的情况下，不重复完整初始化
					console.log('内容列表已初始化，使用已有数据');
				}
			});
			
			// 检查是否是首次打开应用
			const isFirstLaunch = !uni.getStorageSync('notFirstLaunch');
			if (isFirstLaunch) {
				// 标记已非首次启动
				uni.setStorageSync('notFirstLaunch', true);
				// 延迟显示引导弹窗，确保界面已加载完成
				setTimeout(() => {
					this.showLocationGuide = true;
				}, 1000);
			}
			
			// 显示分享菜单
			this.showShareMenu();
		},
		onNavigationBarButtonTap(e) {
			if (e.index === 0) {
				this.navigateToLocationSelect()
			}
		},
		methods: {
			// 添加一个分阶段加载的方法
			staggeredLoading() {
				// 第一阶段：必要数据（分类菜单）
				setTimeout(() => {
					this.fetchMenuCategories();
				}, 200);
				
				// 第二阶段：次要数据（基础设置）
				setTimeout(() => {
					this.fetchBaseSettings();
				}, 800);
				
				// 第三阶段：可见区域内的组件（延迟到页面可见后）
				setTimeout(() => {
					this.initVisibleComponents();
				}, 1500);
			},
			
			// 初始化可见区域内的组件
			initVisibleComponents() {
				// 先初始化功能区域
				if (this.$refs.functionArea) {
					this.$refs.functionArea.loadMiniProgramList();
				}
				
				// 再初始化轮播图
				setTimeout(() => {
					if (this.$refs.swiperBanner) {
						this.$refs.swiperBanner.fetchBannerData();
					}
				}, 300);
				
				// 最后初始化活动区域
				setTimeout(() => {
					if (this.$refs.activityArea) {
						this.$refs.activityArea.fetchActivityData();
					}
				}, 600);
			},
			
			// 获取菜单分类数据
			fetchMenuCategories() {
				this.menuLoading = true;
				
				// 添加默认推荐选项
				const defaultMenu = {
					id: 0,
					name: '推荐'
				};
				
				// 添加超时处理
				const fetchTimeout = setTimeout(() => {
					if (this.menuLoading) {
						// 超时后使用默认菜单并关闭加载状态
						this.useDefaultMenu();
						this.menuLoading = false;
						this.pageLoading = false;
						console.log('获取菜单分类超时，使用默认值');
					}
				}, 3000);
				
				// 调用接口获取分类数据
				get('/wx/client/content/categories', { type: 1 })
					.then(res => {
						clearTimeout(fetchTimeout);
						if (res.code === 0) {
							// 处理返回数据
							const categoryList = res.data.list || [];
							
							// 将推荐选项添加到列表最前面
							this.menuList = [defaultMenu, ...categoryList];
							
							console.log('成功获取菜单分类:', this.menuList);
						} else {
							console.error('获取菜单分类失败:', res.message || '未知错误');
							// 加载失败时使用默认菜单
							this.useDefaultMenu();
						}
					})
					.catch(err => {
						clearTimeout(fetchTimeout);
						console.error('请求菜单分类接口出错:', err);
						// 加载失败时使用默认菜单
						this.useDefaultMenu();
					})
					.finally(() => {
						clearTimeout(fetchTimeout);
						this.menuLoading = false;
						this.pageLoading = false;
					});
			},
			
			// 接口失败时使用默认菜单
			useDefaultMenu() {
				this.menuList = [];
			},
			
			// 初始化位置信息
			initLocation() {
				// 优先使用本地存储的位置信息
				const savedLocation = uni.getStorageSync('currentLocation');
				if (savedLocation) {
					this.currentLocation = savedLocation;
					console.log('使用本地存储的位置信息:', savedLocation);
					
					// 检查是否有关联的地区ID
					const locationId = uni.getStorageSync('currentLocationId');
					if (locationId) {
						console.log('使用已保存的地区ID:', locationId);
						// 通知其他组件位置已变更
						uni.$emit('locationChanged', { regionId: locationId });
						
						// 仅在未初始化时才调用内容列表初始化方法
						this.$nextTick(() => {
							if (this.$refs.contentList && !this.contentListInitialized) {
								console.log('位置初始化后，调用内容列表初始化方法');
								this.$refs.contentList.init();
								this.contentListInitialized = true; // 标记为已初始化
							}
						});
					}
					
					return;
				}
				
				// 如果没有本地存储的位置，则使用默认位置
				this.useDefaultLocation();
			},
			
			// 使用默认位置
			useDefaultLocation() {
				// 检查是否已有存储的位置信息
				const savedLocation = uni.getStorageSync('currentLocation');
				if (savedLocation) {
					this.currentLocation = savedLocation;
					console.log('使用已保存的位置:', savedLocation);
					
					// 检查是否有关联的地区ID
					const locationId = uni.getStorageSync('currentLocationId');
					if (locationId) {
						console.log('使用已保存的地区ID:', locationId);
						// 通知其他组件位置已变更
						uni.$emit('locationChanged', { regionId: locationId });
						
						// 仅在未初始化时才调用内容列表初始化方法
						this.$nextTick(() => {
							if (this.$refs.contentList && !this.contentListInitialized) {
								console.log('使用已保存位置后，调用内容列表初始化方法');
								this.$refs.contentList.init();
								this.contentListInitialized = true; // 标记为已初始化
							}
						});
					}
					
					return;
				}
				
				// 尝试从store中获取地区列表
				const store = this.$store;
				if (store && store.state.region && store.state.region.regionList.length > 0) {
					// 使用第一个地区作为默认位置
					const firstRegion = store.state.region.regionList[0];
					this.currentLocation = firstRegion.name;
					
					// 保存区域ID到本地存储
					console.log('从store中设置默认地区ID:', firstRegion.id);
					uni.setStorageSync('currentLocationId', firstRegion.id);
					uni.setStorageSync('currentLocation', firstRegion.name);
					
					// 触发位置更改事件
					uni.$emit('locationChanged', { regionId: firstRegion.id });
					
					// 仅在未初始化时才调用内容列表初始化方法
					this.$nextTick(() => {
						if (this.$refs.contentList && !this.contentListInitialized) {
							console.log('设置默认位置后，调用内容列表初始化方法');
							this.$refs.contentList.init();
							this.contentListInitialized = true; // 标记为已初始化
						}
					});
					
					console.log('设置默认位置:', firstRegion.name, '来自', firstRegion.location);
				} else {
					// 如果没有地区列表，先使用默认值，同时发起获取地区列表的请求
					this.currentLocation = '请选择地区';
					if (store && store.dispatch) {
						store.dispatch('region/getRegionList').then(() => {
							// 获取到地区列表后更新位置
							if (store.state.region.regionList.length > 0) {
								const firstRegion = store.state.region.regionList[0];
								this.currentLocation = firstRegion.name;
								uni.setStorageSync('currentLocation', this.currentLocation);
								
								// 保存区域ID到本地存储
								console.log('异步获取后设置地区ID:', firstRegion.id);
								uni.setStorageSync('currentLocationId', firstRegion.id);
								
								// 触发位置更改事件
								uni.$emit('locationChanged', { regionId: firstRegion.id });
								
								// 仅在未初始化时才调用内容列表初始化方法
								this.$nextTick(() => {
									if (this.$refs.contentList && !this.contentListInitialized) {
										console.log('异步获取位置后，调用内容列表初始化方法');
										this.$refs.contentList.init();
										this.contentListInitialized = true; // 标记为已初始化
									}
								});
								
								console.log('异步更新默认位置:', firstRegion.name);
							}
						}).catch(err => {
							console.error('获取地区列表失败:', err);
						});
					}
				}
				
				// 保存到本地存储，仅当没有已保存的位置时
				if (!savedLocation) {
					uni.setStorageSync('currentLocation', this.currentLocation);
				}
			},
			
			// 关闭位置引导弹窗
			closeLocationGuide() {
				this.showLocationGuide = false;
			},
			
			// 引导页面跳转到位置选择
			goToLocationSelectFromGuide() {
				this.showLocationGuide = false;
				this.navigateToLocationSelect();
			},
			
			// 跳转到位置选择页面
			navigateToLocationSelect() {
				uni.navigateTo({
					url: '/pages/location/select/index',
					events: {
						// 选择位置后的回调
						locationSelected: (data) => {
							if (data && data.name) {
								// 更新当前位置
								this.currentLocation = data.name;
								// 保存位置到本地存储
								uni.setStorageSync('currentLocation', data.name);
								
								// 保存区域ID到本地存储并触发位置更改事件
								if (data.id) {
									// 打印当前数据，确认ID存在
									console.log('位置选择数据:', data);
									
									// 先清除旧的区域ID
									uni.removeStorageSync('currentLocationId');
									// 保存区域ID到本地存储
									uni.setStorageSync('currentLocationId', data.id);
									
									// 强制刷新内容列表
									if (this.$refs.contentList) {
										console.log('直接调用内容列表组件方法更新区域ID:', data.id);
										this.$refs.contentList.regionId = data.id;
										this.$refs.contentList.resetAndLoad();
									}
									
									// 触发位置更改事件
									console.log('发送locationChanged事件，regionId:', data.id);
									uni.$emit('locationChanged', { regionId: data.id });
								}
								
								console.log('用户选择了新位置:', data.name, 'ID:', data.id);
							}
						}
					}
				});
			},
			switchMenu(index) {
				// 防止菜单未加载完成时点击
				if (this.menuLoading || !this.menuList.length) {
					return;
				}
				
				if (this.currentMenu !== index) {
					// 更新当前选中的菜单
					this.currentMenu = index;
					
					// 获取选中分类的ID
					const selectedCategory = this.menuList[index];
					const categoryId = selectedCategory ? selectedCategory.id : '';
					
					// 通知内容列表组件更新分类
					if (this.$refs.contentList) {
						this.$refs.contentList.handleCategoryChange(categoryId);
					}
				}
			},
			calculateHeaderHeight() {
				const query = uni.createSelectorQuery().in(this)
				query.select('.fixed-header').boundingClientRect(data => {
					this.headerHeight = data.height
				}).exec()
			},
			handleSearch(e) {
				const keyword = e.detail.value
				console.log('搜索关键词:', keyword)
				// 这里添加搜索处理逻辑
			},
			// 获取菜单导航的位置
			getMenuPosition() {
				const query = uni.createSelectorQuery().in(this)
				
				// 获取固定头部高度
				query.select('.fixed-header').boundingClientRect(data => {
					if (data) {
						this.headerHeight = data.height
					}
				}).exec()
				
				// 获取分类菜单位置
				query.select('#category-menu').boundingClientRect(data => {
					if (data) {
						this.menuTop = data.top
					}
				}).exec()
				
				// 获取内容列表区域的位置
				query.select('#content-area').boundingClientRect(data => {
					if (data) {
						// 内容区域的绝对位置减去头部高度，得到滚动时需要的相对位置
						this.contentTop = data.top - this.headerHeight
						console.log('内容区域位置:', this.contentTop)
					}
				}).exec()
			},
			// 处理滚动事件
			handleScroll(e) {
				const scrollTop = e.detail.scrollTop
				this.oldScrollTop = scrollTop // 记录当前滚动位置
				
				// 当滚动超过一定距离时显示返回顶部按钮
				this.showBackToTop = scrollTop > 500
				
				// 内容区域位置存在且有效时才进行比较
				if (this.contentTop && this.contentTop > 0) {
					// 更新固定菜单显示状态 - 改为滚动到内容区域时显示
					this.showFixedMenu = scrollTop >= this.contentTop
				} else {
					// 防止位置计算错误时的降级处理
					this.showFixedMenu = scrollTop >= this.menuTop
				}
			},
			
			// 返回顶部
			scrollToTop() {
				// 先将 scrollTop 设置为 1
				this.scrollTop = 1
				
				// 使用 nextTick 确保视图更新后再设置为 0
				this.$nextTick(() => {
					this.scrollTop = 0
				})
			},
			// 处理底部导航切换
			handleTabChange(index) {
				console.log('切换到标签:', index)
				// 这里可以添加导航逻辑
			},
			// 处理发布按钮点击
			handlePublish() {
				console.log('点击发布按钮')
				// 这里可以添加发布逻辑
			},
			// 处理下拉刷新
			handleRefresh() {
				this.refreshing = true;
				// 刷新菜单分类
				this.fetchMenuCategories();
				
				// 刷新功能区域
				if (this.$refs.functionArea) {
					this.$refs.functionArea.loadMiniProgramList();
				}
				
				// 刷新活动区域
				if (this.$refs.activityArea) {
					this.$refs.activityArea.fetchActivityData();
				}
				
				// 刷新轮播图
				if (this.$refs.swiperBanner) {
					this.$refs.swiperBanner.fetchBannerData();
				}
				
				// 刷新内容列表（如果存在）
				if (this.$refs.contentList) {
					this.$refs.contentList.resetAndLoad();
				}
				
				// 刷新分享数据
				this.initShareData();
				
				// 设置超时自动关闭刷新状态
				setTimeout(() => {
					if (this.refreshing) {
						this.refreshing = false;
						console.log('刷新状态超时自动关闭');
					}
				}, 3000);
			},
			// 处理上拉加载更多
			handleScrollToLower() {
				console.log('触发上拉加载更多');
				// 调用内容列表组件的加载更多方法
				if (this.$refs.contentList) {
					this.$refs.contentList.loadMore();
				}
			},
			// 获取小程序基础设置
			fetchBaseSettings() {
				settings.getBaseSettings()
					.then(res => {
						if (res.code === 0 && res.data) {
							// 只更新logo字段
							if (res.data.logo) {
								this.logoUrl = res.data.logo;
								console.log('成功获取小程序Logo:', this.logoUrl);
							}
						} else {
							console.error('获取小程序基础设置失败:', res.message || '未知错误');
						}
					})
					.catch(err => {
						console.error('请求小程序基础设置接口出错:', err);
					});
			},
			// 显示分享菜单，包括分享朋友圈
			showShareMenu() {
				// 先更新分享数据
				this.initShareData();
				
				// 显示分享菜单
				uni.showShareMenu({
					withShareTicket: true,
					menus: ['shareAppMessage', 'shareTimeline']
				});
			},
		}
	}
</script>

<style>
	.container {
		width: 100%;
		height: 100vh;
		background-color: #ffffff;
		position: relative;
		padding-bottom: 100rpx;
	}
	
	.scroll-content {
		height: 100vh;
		box-sizing: border-box;
		padding-bottom: env(safe-area-inset-bottom);
	}
	
	.fixed-header {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 100;
		background-color: #ffffff;
	}
	
	.gradient-bg {
		background: linear-gradient(180deg, #fc3e2b 0%, #fa7154 70%);
		border-bottom-left-radius: 40rpx;
		border-bottom-right-radius: 40rpx;
	}
	
	/* 自定义导航栏样式 */
	.custom-nav {
		width: 100%;
	}
	
	.nav-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 30rpx;
	}
	
	.nav-logo {
		flex: 1;
		display: flex;
		align-items: center;
	}
	
	.logo-img {
		height: 100%;
	}
	
	.nav-center {
		flex: 1;
		display: flex;
		justify-content: center;
	}
	
	.nav-right {
		flex: 1;
	}
	
	.nav-location {
		display: flex;
		align-items: center;
		padding: 4rpx 20rpx;
		background-color: rgba(255, 255, 255, 0.1);
		border-radius: 24rpx;
	}
	
	.location-text {
		font-size: 28rpx;
		color: #ffffff;
		margin-right: 6rpx;
	}
	
	/* 搜索框样式 */
	.menu-wrapper {
		padding: 20rpx 30rpx;
	}
	
	.search-box {
		width: 100%;
	}
	
	.search-input {
		display: flex;
		align-items: center;
		background-color: rgba(255, 255, 255, 0.9);
		padding: 0 24rpx;
		border-radius: 36rpx;
		height: 72rpx;
	}
	
	.search-input input {
		flex: 1;
		margin-left: 16rpx;
		font-size: 28rpx;
		height: 100%;
	}
	
	.placeholder {
		color: #999999;
		font-size: 28rpx;
	}
	
	/* 新增分类菜单样式 */
	.category-menu {
		background-color: #ffffff;
		padding: 20rpx 0;
		margin-top: 20rpx;
		width: 100%;
		overflow: hidden;
	}
	
	.nav-menu {
		width: 100%;
		white-space: nowrap;
		padding: 0 30rpx;
		box-sizing: border-box;
		padding-right: 40rpx;
	}
	
	.menu-item {
		padding: 12rpx 24rpx;
		height: 100%;
		display: inline-flex;
		align-items: center;
		font-size: 30rpx;
		white-space: nowrap;
		position: relative;
	}
	
	.menu-item:last-child {
		margin-right: 10rpx;
	}
	
	.menu-text {
		font-size: 28rpx;
		color: #333333;
		line-height: 40rpx;
		font-weight: 400;
	}
	
	.menu-item.active .menu-text {
		color: #fc3e2b;
		font-weight: 700;
	}
	
	.menu-item.active::after {
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
	
	/* 添加内容区域样式 */
	.page-content {
		padding-top: 0;
		background-color: #ffffff;
		min-height: 100vh;
	}
	
	/* 固定菜单样式 */
	.fixed-nav-menu {
		width: 100%;
		white-space: nowrap;
		padding: 0 20rpx;
		box-sizing: border-box;
		transform: translateY(-100%);
		opacity: 0;
		transition: all 0.25s ease-out;
	}
	
	.fixed-nav-menu.menu-show {
		transform: translateY(0);
		opacity: 1;
	}
	
	.fixed-nav-menu .menu-item {
		display: inline-block;
		padding: 12rpx 24rpx;
		position: relative;
	}
	
	.fixed-menu-text {
		font-size: 28rpx;
		color: #ffffff;
		line-height: 40rpx;
		font-weight: 400;
	}
	
	.fixed-nav-menu .menu-item.active .fixed-menu-text {
		font-weight: 700;
	}
	
	.fixed-nav-menu .menu-item.active::after {
		content: '';
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 40rpx;
		height: 4rpx;
		background-color: #ffffff;
		border-radius: 2rpx;
	}
	
	/* 返回顶部按钮样式 */
	.back-to-top {
		position: fixed;
		left: 30rpx;
		bottom: 180rpx;
		width: 80rpx;
		height: 80rpx;
		background: #fa7154;
		border-radius: 50%;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.1);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 99;
	}
	
	/* 图标容器样式 */
	.back-to-top .uni-icons {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
		color: #ffffff;
	}
	
	/* 确保图标本身居中 */
	.back-to-top :deep(.uni-icons__icon) {
		margin: auto;
	}
	
	/* 菜单加载状态样式 */
	.menu-loading {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0 30rpx;
		height: 100%;
	}
	
	.loading-text {
		font-size: 28rpx;
		color: #999999;
	}
	
	.fixed-nav-menu .loading-text {
		color: #ffffff;
	}
	
	/* 位置引导弹窗样式 */
	.location-guide-mask {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.6);
		z-index: 999;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.location-guide-content {
		width: 80%;
		max-width: 600rpx;
		background-color: #ffffff;
		border-radius: 16rpx;
		overflow: hidden;
		box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.15);
	}
	
	.guide-header {
		padding: 30rpx;
		display: flex;
		justify-content: center;
		align-items: center;
		border-bottom: 1px solid #f5f5f5;
	}
	
	.guide-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
	}
	
	.guide-body {
		padding: 40rpx 30rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
	.guide-text {
		font-size: 28rpx;
		color: #333;
		text-align: center;
		line-height: 1.5;
		margin-bottom: 20rpx;
	}
	
	.guide-current {
		font-size: 26rpx;
		color: #666;
		margin-top: 20rpx;
		padding: 10rpx 24rpx;
		background-color: #f8f8f8;
		border-radius: 24rpx;
	}
	
	.guide-footer {
		padding: 30rpx;
		display: flex;
		justify-content: center;
		border-top: 1px solid #f5f5f5;
	}
	
	.guide-btn {
		width: 70%;
		height: 80rpx;
		line-height: 80rpx;
		text-align: center;
		border-radius: 40rpx;
		font-size: 28rpx;
		margin: 0;
	}
	
	.guide-btn-confirm {
		background-color: #fc3e2b;
		color: #ffffff;
	}
	
	/* 骨架屏样式 */
	.skeleton-screen {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100vh;
		background-color: #f8f8f8;
		z-index: 999;
		padding: 0 30rpx;
		box-sizing: border-box;
	}
	
	.skeleton-header {
		padding-top: var(--status-bar-height);
	}
	
	.skeleton-nav {
		height: 44px;
		background-color: #e0e0e0;
		border-radius: 8rpx;
		margin: 20rpx 0;
		animation: skeleton-loading 1.5s infinite ease-in-out;
	}
	
	.skeleton-search {
		height: 72rpx;
		background-color: #e0e0e0;
		border-radius: 36rpx;
		margin: 20rpx 0;
		animation: skeleton-loading 1.5s infinite ease-in-out;
	}
	
	.skeleton-content {
		margin-top: 40rpx;
	}
	
	.skeleton-item {
		height: 200rpx;
		background-color: #e0e0e0;
		border-radius: 12rpx;
		margin-bottom: 30rpx;
		animation: skeleton-loading 1.5s infinite ease-in-out;
	}
	
	@keyframes skeleton-loading {
		0% {
			opacity: 0.7;
		}
		50% {
			opacity: 0.4;
		}
		100% {
			opacity: 0.7;
		}
	}
</style>
