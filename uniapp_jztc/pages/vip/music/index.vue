<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#ffecd8"></uni-icons>
				</view>
				<text class="nav-title">天天领音乐会员</text>
			</view>
		</view>
		
		<!-- 顶部标题区域 -->
		<view class="header">
			<view class="header-content">
				<text class="title">不花钱 天天领会员</text>
				<view class="subtitle-wrapper">
					<text class="subtitle">说明</text>
				</view>
			</view>
		</view>
		
		<!-- 会员信息卡片 -->
		<view class="vip-info-card">
			<view class="user-info">
				<view class="title-bg">
					<view class="avatar-wrapper">
						<image class="avatar" :src="userAvatar" mode="aspectFill"></image>
					</view>
					<text class="user-title">我的会员时长</text>
				</view>
			</view>
			<view class="time-info">
				<text class="days">{{memberDuration.days}}</text>
				<text class="unit">天</text>
				<text class="days">{{memberDuration.hours}}</text>
				<text class="unit">小时</text>
				<uni-icons type="right" size="14" color="#862c13" class="time-arrow"></uni-icons>
				<view class="exchange-btn" @tap="goToExchange">
					<text>兑换会员</text>
					<uni-icons type="right" size="14" color="#862c13"></uni-icons>
				</view>
			</view>
		</view>
		
		<!-- 平台提示文字 -->
		<view class="platform-tip">
			<text>时长可任选兑换9大平台会员</text>
		</view>
		
		<!-- 平台列表 -->
		<scroll-view 
			class="platform-list" 
			scroll-x 
			:show-scrollbar="false"
			:enhanced="true"
			:bounces="true"
		>
			<view class="platform-scroll">
				<view 
					class="platform-item" 
					v-for="(item, index) in platforms" 
					:key="index"
					:class="{ active: currentCategoryId === item.id }"
					@tap="selectPlatform(item.id)"
				>
					<image :src="item.icon" mode="aspectFit"></image>
					<text>{{item.name}}</text>
				</view>
			</view>
		</scroll-view>
		
		<!-- 会员套餐列表 -->
		<scroll-view 
			class="package-list" 
			scroll-x 
			:show-scrollbar="false"
			:enhanced="true"
			:bounces="true"
		>
			<view class="package-scroll">
				<template v-if="packages.length > 0">
					<view 
						class="package-item" 
						v-for="(item, index) in packages" 
						:key="index"
					>
						<view class="package-tag">{{item.tag}}</view>
						<view class="package-content">
							<view class="name-row">
								<text class="package-name">{{item.name}}</text>
								<view class="month-tag">{{item.description}}</view>
							</view>
							<view class="price-info">
								<view class="price-row">
									<text class="price-symbol">需</text>
									<text class="package-price">{{item.days}}天</text>
								</view>
								<text class="duration-text">时长</text>
							</view>
						</view>
					</view>
				</template>
				<view v-if="packages.length === 0" class="empty-package">
					<text>暂无套餐</text>
				</view>
			</view>
		</scroll-view>
		
		<!-- 底部按钮 -->
		<view class="watch-btn" @tap="handleWatch">
			<text>看10秒 领会员</text>
		</view>
        		<!-- 套餐提示文字 -->
		<view class="package-tip">
			<uni-icons type="info-filled" size="20" color="#ffecd8"></uni-icons>
			<text>以上会员点击可兑</text>
		</view>
		
		<!-- 底部banner区域 -->
		<view class="banner-area">
			<image 
				class="banner-image" 
				src="/static/demo/1.png" 
				mode="widthFix"
			></image>
		</view>
	</view>
</template>

<script>
	import deviceMixin from '@/mixins/device.js'
	import shareMixin from '@/mixins/share.js'
	import { request, get, post } from '@/utils/request.js'
	import { getUserInfo } from '@/utils/storage.js'
	import { mapState } from 'vuex'
	import { getShareSettings } from '@/utils/share.js'
	
	export default {
		mixins: [deviceMixin, shareMixin],
		data() {
			return {
				platforms: [],
				packages: [],
				currentCategoryId: null,
				adInfo: {
					enableWxAd: false,
					rewardedVideoAdId: ''
				},
				videoAd: null,
				isLoading: false,
				userAvatar: '/static/images/default-avatar.png',
				memberDuration: {
					days: '0',
					hours: '0'
				},
				shareImageUrl: ''
			}
		},
		computed: {
			...mapState('user', ['userInfo', 'isLogin'])
		},
		onLoad() {
			// 获取广告设置
			this.getAdSettings()
			// 获取平台列表
			this.getPlatformList()
			// 获取用户头像
			this.getUserAvatar()
			// 获取会员时长
			this.getMemberDuration()
			// 获取分享图片
			this.getShareImage()
		},
		// 自定义分享给好友
		onShareAppMessage() {
			return {
				title: "天天领音乐会员", // 使用页面导航栏标题
				path: '/pages/vip/music/index',
				imageUrl: this.shareImageUrl
			}
		},
		// 自定义分享到朋友圈
		onShareTimeline() {
			return {
				title: "天天领音乐会员", // 使用页面导航栏标题
				query: '',
				imageUrl: this.shareImageUrl
			}
		},
		methods: {
			// 获取用户头像
			getUserAvatar() {
				try {
					// 从Vuex获取用户信息
					if (this.isLogin && this.userInfo && this.userInfo.avatarUrl) {
						this.userAvatar = this.userInfo.avatarUrl
						return
					}
					
					// 从本地存储获取用户信息
					const userInfo = getUserInfo()
					if (userInfo && userInfo.avatarUrl) {
						this.userAvatar = userInfo.avatarUrl
					} else {
						// 如果没有头像，使用微信默认头像
						this.userAvatar = 'https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0'
					}
				} catch (error) {
					console.error('获取用户头像失败', error)
					// 出错时也使用微信默认头像
					this.userAvatar = 'https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0'
				}
			},
			// 获取平台列表
			async getPlatformList() {
				try {
					const res = await get('/wx/client/shop-category/list')
					if (res.code === 0 && res.data && res.data.list && res.data.list.length > 0) {
						this.platforms = res.data.list.map(item => {
							return {
								id: item.id,
								name: item.name,
								icon: item.image
							}
						})
						
						// 获取第一个平台的会员套餐
						if (this.platforms.length > 0) {
							this.currentCategoryId = this.platforms[0].id
							this.getPackageList(this.currentCategoryId)
						}
					}
				} catch (error) {
					console.error('获取平台列表失败', error)
				}
			},
			
			// 获取会员套餐列表
			async getPackageList(categoryId) {
				try {
					console.log('开始获取会员套餐列表，分类ID:', categoryId)
					uni.showLoading({ title: '加载中...' })
					
					// 使用request代替get
					const res = await request({
						url: '/wx/product/list',
						method: 'GET',
						data: { categoryId: categoryId }
					})
					
					uni.hideLoading()
					
					console.log('套餐接口响应:', JSON.stringify(res))
					
					// 先清空数组
					this.packages = []
					
					if (res.code === 0) {
						// 处理数据，如果list是null或空数组，保持packages为空数组
						if (res.data && res.data.list && res.data.list.length > 0) {
							setTimeout(() => {
								this.packages = res.data.list.map(item => {
									return {
										id: item.id,
										tag: item.tags,
										name: item.name,
										price: item.price.toFixed(2),
										days: item.duration,
										description: item.description
									}
								})
								console.log('处理后的套餐数据:', JSON.stringify(this.packages))
							}, 100)
						} else {
							console.log('当前平台没有可用套餐')
						}
					} else {
						console.log('接口请求失败:', res.message || '未知错误')
					}
				} catch (error) {
					uni.hideLoading()
					console.error('获取会员套餐列表失败', error)
				}
			},
			
			// 选择平台
			selectPlatform(categoryId) {
				console.log('选择平台:', categoryId, '当前平台:', this.currentCategoryId)
				if (this.currentCategoryId !== categoryId) {
					this.currentCategoryId = categoryId
					this.getPackageList(categoryId)
				}
			},
			
			// 获取广告设置
			async getAdSettings() {
				let loading = false
				try {
					loading = true
					uni.showLoading({ title: '加载中...' })
					
					const res = await request({
						url: '/wx/ad/settings',
						method: 'GET'
					})
					
					if (res && res.code === 0) {
						this.adInfo = res.data
						
						// 如果启用了微信广告，创建广告实例
						if (this.adInfo.enableWxAd && this.adInfo.rewardedVideoAdId) {
							this.createRewardedVideoAd()
						}
					} else {
						console.error('获取广告设置失败:', res?.message)
					}
				} catch (error) {
					console.error('获取广告设置异常:', error)
				} finally {
					if (loading) {
						uni.hideLoading()
					}
				}
			},
			
			// 创建激励视频广告
			createRewardedVideoAd() {
				// 检查环境是否支持
				if (!wx || !wx.createRewardedVideoAd) {
					console.error('当前环境不支持激励视频广告')
					return
				}
				
				// 创建广告实例
				this.videoAd = wx.createRewardedVideoAd({
					adUnitId: this.adInfo.rewardedVideoAdId
				})
				
				// 监听加载事件
				this.videoAd.onLoad(() => {
					console.log('激励视频广告加载成功')
				})
				
				// 监听错误事件
				this.videoAd.onError(err => {
					console.error('激励视频广告出错:', err)
					uni.showToast({
						title: '广告加载失败',
						icon: 'none'
					})
				})
				
				// 监听关闭事件
				this.videoAd.onClose(res => {
					// 用户点击了"关闭广告"按钮
					if (res && res.isEnded) {
						// 正常播放结束，可以下发奖励
						this.reportAdViewed()
					} else {
						// 播放中途退出
						uni.showToast({
							title: '请完整观看广告才能获得奖励',
							icon: 'none'
						})
					}
				})
			},
			
			// 显示激励视频广告
			showRewardedVideoAd() {
				if (!this.videoAd) {
					uni.showToast({
						title: '广告加载中，请稍后再试',
						icon: 'none'
					})
					return
				}
				
				this.videoAd.show().catch(() => {
					// 失败重试
					this.videoAd.load()
						.then(() => this.videoAd.show())
						.catch(err => {
							console.error('激励视频广告显示失败:', err)
							uni.showToast({
								title: '广告加载失败，请稍后再试',
								icon: 'none'
							})
						})
				})
			},
			
			// 上报广告观看完成
			async reportAdViewed() {
				if (this.isLoading) return
				
				this.isLoading = true
				let loading = false
				
				try {
					loading = true
					uni.showLoading({ title: '领取奖励中...' })
					
					const res = await post('/wx/ad/reward/viewed')
					
					if (res && res.code === 0) {
						// 使用接口返回的消息提示
						const message = res.data?.message || '恭喜获得会员奖励！'
						
						// 确保先隐藏loading再显示toast
						if (loading) {
							uni.hideLoading()
							loading = false
						}
						
						uni.showToast({
							title: message,
							icon: 'success',
							duration: 3000
						})
						
						// 延迟后刷新页面数据
						setTimeout(() => {
							this.refreshUserInfo()
						}, 2000)
					} else {
						// 确保先隐藏loading再显示toast
						if (loading) {
							uni.hideLoading()
							loading = false
						}
						
						uni.showToast({
							title: res?.message || '领取失败，请重试',
							icon: 'none'
						})
					}
				} catch (error) {
					console.error('上报广告观看失败:', error)
					
					// 确保先隐藏loading再显示toast
					if (loading) {
						uni.hideLoading()
						loading = false
					}
					
					uni.showToast({
						title: '网络异常，请重试',
						icon: 'none'
					})
				} finally {
					// 最后确保loading已关闭
					if (loading) {
						uni.hideLoading()
					}
					this.isLoading = false
				}
			},
			
			// 刷新用户信息
			refreshUserInfo() {
				// 这里可以添加刷新用户会员信息的逻辑
				// 例如重新请求会员天数等
				let loading = false
				try {
					loading = true
					uni.showLoading({ title: '刷新数据...' })
					
					// 获取最新会员时长
					this.getMemberDuration()
					
					// 模拟刷新延迟
					setTimeout(() => {
						// 可以添加实际刷新逻辑
					}, 1000)
				} catch (error) {
					console.error('刷新用户信息失败:', error)
				} finally {
					// 确保至少显示1秒的加载状态
					setTimeout(() => {
						if (loading) {
							uni.hideLoading()
						}
					}, 1000)
				}
			},
			
			// 获取会员时长
			async getMemberDuration() {
				try {
					const res = await get('/wx/client/duration/remaining')
					if (res.code === 0 && res.data) {
						// 将接口返回的完整时长提取出天数和小时数部分
						const durationStr = res.data.remainingDuration || '0天0小时0分钟'
						
						// 提取天数和小时数，默认为0
						const days = durationStr.match(/(\d+)天/) ? durationStr.match(/(\d+)天/)[1] : '0'
						const hours = durationStr.match(/(\d+)小时/) ? durationStr.match(/(\d+)小时/)[1] : '0'
						
						this.memberDuration = {
							days,
							hours
						}
					} else {
						console.error('获取会员时长失败:', res?.message)
					}
				} catch (error) {
					console.error('获取会员时长异常:', error)
				}
			},
			
			handleBack() {
				uni.navigateBack()
			},
			
			// 修改观看按钮点击事件
			handleWatch() {
				// 检查是否启用了广告
				if (this.adInfo.enableWxAd && this.adInfo.rewardedVideoAdId) {
					this.showRewardedVideoAd()
				} else {
					uni.showToast({
						title: '广告功能暂未开放',
						icon: 'none'
					})
				}
			},
			
			goToExchange() {
				uni.navigateTo({
					url: '/pages/vip/exchange/index'
				})
			},
			
			// 获取分享图片
			async getShareImage() {
				try {
					// 使用共享的getShareSettings方法
					const settings = await getShareSettings();
					if (settings && settings.default_share_image) {
						this.shareImageUrl = settings.default_share_image;
					}
				} catch (error) {
					console.error('获取分享图片失败:', error);
				}
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background: linear-gradient(180deg, 
			#21242d 0%, 
			#21242d 40%, 
			#1f2123 45%,
			#1d1f1e 48%,
			#1b1d1a 50%,
			#191b17 52%,
			#181a16 55%,
			#181a15 60%,
			#181a15 100%
		);
		padding-bottom: calc(env(safe-area-inset-bottom) + 120rpx);
	}
	
	/* 导航栏样式 */
	.nav-bar {
		position: relative;
		z-index: 100;
	}
	
	.nav-content {
		height: 44px;
		display: flex;
		align-items: center;
		position: relative;
	}
	
	.back-btn {
		position: absolute;
		left: 0;
		top: 50%;
		transform: translateY(-50%);
		padding: 20rpx;
		margin-left: 10rpx;
		z-index: 1;
	}
	
	.nav-title {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		font-size: 34rpx;
		color: #ffecd8;
		font-weight: 500;
	}
	
	.header {
		margin-top: 40rpx;
		margin-bottom: 50rpx;
		padding: 0 30rpx;
	}
	
	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-right: -30rpx;
		position: relative;
		padding: 20rpx 0;
	}
	
	.title {
		font-size: 58rpx;
		color: #ffecd8;
		font-weight: bold;
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
		white-space: nowrap;
	}
	
	.subtitle-wrapper {
		background: #302f35;
		padding: 6rpx 16rpx;
		border-radius: 24rpx;
		border-top-right-radius: 0;
		border-bottom-right-radius: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		width: fit-content;
		margin-left: auto;
	}
	
	.subtitle {
		font-size: 28rpx;
		color: #ffecd8;
		opacity: 0.8;
		text-align: center;
		white-space: nowrap;
	}
	
	.vip-info-card {
		background: linear-gradient(135deg, #ffe4d6 0%, #ffc4b7 100%);
		border-radius: 24rpx;
		padding: 30rpx;
		margin-bottom: 30rpx;
		margin-left: 30rpx;
		margin-right: 30rpx;
	}
	
	.user-info {
		display: flex;
		align-items: center;
		margin-bottom: 20rpx;
	}
	
	.title-bg {
		display: flex;
		align-items: center;
		background: rgba(0, 0, 0, 0.06);
		padding: 8rpx 20rpx 8rpx 12rpx;
		border-radius: 50rpx 0 0 50rpx;
	}
	
	.avatar-wrapper {
		width: 48rpx;
		height: 48rpx;
		border-radius: 50%;
		background: #ffffff;
		padding: 2rpx;
		margin-right: 8rpx;
	}
	
	.avatar {
		width: 100%;
		height: 100%;
		border-radius: 50%;
	}
	
	.user-title {
		font-size: 28rpx;
		color: #a04a39;
		font-weight: bold;
	}
	
	.time-info {
		display: flex;
		align-items: center;
		margin-bottom: 20rpx;
		padding-left: 12rpx;
	}
	
	.days {
		font-size: 64rpx;
		font-weight: bold;
		color: #862c13;
		line-height: 1;
	}
	
	.unit {
		font-size: 32rpx;
		color: #862c13;
		margin-left: 8rpx;
		align-self: flex-end;
		margin-bottom: 8rpx;
	}
	
	.time-arrow {
		margin-left: 4rpx;
		align-self: flex-end;
		margin-bottom: 8rpx;
	}
	
	.exchange-btn {
		margin-left: auto;
		background: #ffffff;
		padding: 12rpx 24rpx;
		border-radius: 32rpx;
		display: flex;
		align-items: center;
		gap: 4rpx;
	}
	
	.exchange-btn text {
		font-size: 28rpx;
		color: #862c13;
		font-weight: bold;
	}
	
	.platform-tip {
		text-align: center;
		margin-bottom: 20rpx;
		margin-left: 30rpx;
		margin-right: 30rpx;
	}
	
	.platform-tip text {
		font-size: 24rpx;
		color: #ffecd8;
	}
	
	.platform-list {
		width: 100%;
		margin-bottom: 30rpx;
		margin-right: 30rpx;
		white-space: nowrap;
	}
	
	.platform-scroll {
		display: inline-flex;
		padding: 10rpx 0;
		padding-right: 30rpx;
        margin-left: 30rpx;
	}
	
	.platform-item {
		display: inline-flex;
		flex-direction: column;
		align-items: center;
		margin-right: 30rpx;
		width: 100rpx;
		position: relative;
	}
	
	.platform-item.active::after {
		content: '';
		position: absolute;
		bottom: -10rpx;
		left: 50%;
		transform: translateX(-50%);
		width: 30rpx;
		height: 6rpx;
		background-color: #ffecd8;
		border-radius: 3rpx;
	}
	
	.platform-item image {
		width: 80rpx;
		height: 80rpx;
		border-radius: 16rpx;
		margin-bottom: 8rpx;
	}
	
	.platform-item text {
		font-size: 24rpx;
		color: #ffffff;
		white-space: nowrap;
	}
	
	.package-list {
		width: 100%;
		white-space: nowrap;
		margin-bottom: 20rpx;
	}
	
	.package-scroll {
		display: inline-flex;
		padding: 10rpx 0;
		padding-right: 30rpx;
		margin-left: 30rpx;
	}
	
	.package-item {
		position: relative;
		padding: 0;
		border-radius: 16rpx;
		background: linear-gradient(45deg, 
			#ffffff 0%, 
			#ffffff 60%, 
			#f8d5af 100%
		);
		margin-right: 20rpx;
		width: 220rpx;
		height: 140rpx;
		box-sizing: border-box;
		border: 2rpx solid #ffecd8;
	}
	
	.package-tag {
		position: absolute;
		left: -2rpx;
		top: -2rpx;
		background: #4c4945;
		color: #ffecd8;
		font-size: 22rpx;
		padding: 4rpx 12rpx;
		border-radius: 16rpx 0 12rpx 0;
	}
	
	.package-content {
		flex-direction: column;
		position: relative;
		padding: 46rpx 20rpx 24rpx;
		z-index: 1;
	}
	
	.name-row {
		display: flex;
		align-items: center;
		gap: 8rpx;
	}
	
	.package-name {
		font-size: 22rpx;
		color: #333333;
	}
	
	.month-tag {
		font-size: 20rpx;
		color: #6d542b;
		font-weight: bold;
		background-color: #efcd9b;
		padding: 2rpx 8rpx;
		border-radius: 10rpx 0 15rpx 0;
	}
	
	.price-info {
		margin-top: auto;
		display: flex;
		align-items: baseline;
	}
	
	.price-row {
		display: flex;
		align-items: baseline;
	}
	
	.price-symbol {
		font-size: 24rpx;
		color: #fc3e2b;
	}
	
	.package-price {
		font-size: 36rpx;
		color: #fc3e2b;
		font-weight: bold;
	}
	
	.duration-text {
		font-size: 20rpx;
		color: #c2615a;
		margin-left: 4rpx;
	}
	
	.watch-btn {
		background: linear-gradient(135deg, #fc3e2b 0%, #fa7154 100%);
		height: 88rpx;
		border-radius: 44rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 20rpx;
		margin-left: 30rpx;
		margin-right: 30rpx;
	}
	
	.watch-btn text {
		font-size: 32rpx;
		color: #ffffff;
		font-weight: 500;
	}
	
	.package-tip {
		text-align: center;
		margin-bottom: 40rpx;
		margin-left: 30rpx;
		margin-right: 30rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8rpx;
	}
	
	.package-tip text {
		font-size: 24rpx;
		color: #ffecd8;
		opacity: 0.8;
	}
	
	/* 调整图标透明度以匹配文字 */
	.package-tip :deep(.uni-icons) {
		opacity: 0.8;
	}
	
	.banner-area {
		padding: 0 30rpx;
		margin-bottom: 40rpx;
	}
	
	.banner-image {
		width: 100%;
		border-radius: 16rpx;
	}
	
	.empty-package {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 220rpx;
		height: 140rpx;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 16rpx;
		margin-right: 20rpx;
	}
	
	.empty-package text {
		font-size: 28rpx;
		color: #ffecd8;
		opacity: 0.8;
	}
</style> 