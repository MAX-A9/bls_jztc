<template>
	<view class="my-container">
		<!-- 顶部用户信息区域 -->
		<view class="gradient-bg">
			<!-- 自定义状态栏 -->
			<view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
			
			<!-- 自定义导航栏 -->
			<view class="nav-bar" :style="{ height: navBarHeight + 'px' }">
				<view class="nav-left">
					<view class="setting-btn" @click="handleSetting">
						<uni-icons type="gear" size="24" color="#ffffff"></uni-icons>
					</view>
				</view>
				<text class="nav-title">我的</text>
			</view>
			
			<view class="user-header">
				<view class="user-info" @click="handleUserInfoClick">
					<image 
						class="avatar" 
						:src="userInfo.avatarUrl || '/static/demo/0.png'" 
						mode="aspectFill"
					></image>
					<view class="user-detail">
						<text class="nickname white">{{ userInfo.realName || userInfo.nickName || '未登录' }}</text>
						<text class="phone white-light">{{ userInfo.phone || (hasLogin ? '绑定手机号' : '点击登录') }}</text>
					</view>
				</view>
				<view class="sign-btn" @click="hasLogin ? handleSignIn() : handleLogin()">
					<uni-icons :type="hasLogin ? 'gift' : 'person'" size="16" color="#ffffff"></uni-icons>
					<text class="sign-text">{{ hasLogin ? '签到领福利' : '立即登录' }}</text>
				</view>
			</view>
			
			<!-- 功能入口区域 -->
			<view class="feature-wrapper">
				<view class="feature-bar">
					<view class="feature-item" @tap="handleFeatureClick('publish')">
						<text class="feature-num">{{ publishCount }}</text>
						<text class="feature-name">我的发布</text>
					</view>
					<view class="feature-item" @tap="handleFeatureClick('follow')">
						<text class="feature-num">{{ followingCount }}</text>
						<text class="feature-name">我的关注</text>
					</view>
					<view class="feature-item" @tap="handleFeatureClick('favorite')">
						<text class="feature-num">{{ favoriteCount }}</text>
						<text class="feature-name">我的收藏</text>
					</view>
					<view class="feature-item" @tap="handleFeatureClick('history')">
						<text class="feature-num">{{ browseHistoryCount }}</text>
						<text class="feature-name">浏览记录</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 内容卡片区域 -->
		<view class="content-card">
			<!-- 订单区域 -->
			<view class="card">
				<view class="order-section">
					<view class="section-header">
						<text class="section-title">我的订单</text>
						<view class="more-btn" @tap="handleViewAllOrders">
							<text class="more-text">全部订单</text>
							<uni-icons type="right" size="14" color="#999999"></uni-icons>
						</view>
					</view>
					<view class="order-types">
						<view class="order-type-item" @tap="handleOrderType('processing')">
							<uni-icons type="refresh" size="28" color="#fc3e2b"></uni-icons>
							<text class="type-name">进行中</text>
						</view>
						<view class="order-type-item" @tap="handleOrderType('unpaid')">
							<uni-icons type="wallet" size="28" color="#fc3e2b"></uni-icons>
							<text class="type-name">待支付</text>
						</view>
						<view class="order-type-item" @tap="handleOrderType('completed')">
							<uni-icons type="checkbox-filled" size="28" color="#fc3e2b"></uni-icons>
							<text class="type-name">已完成</text>
						</view>
						<view class="order-type-item" @tap="handleOrderType('all')">
							<uni-icons type="list" size="28" color="#fc3e2b"></uni-icons>
							<text class="type-name">全部</text>
						</view>
					</view>
				</view>
			</view>
			
			<!-- 专属管家卡片 -->
			<view class="card manager-card" @tap="showQrCode">
				<view class="manager-content">
					<text class="manager-text">添加专属管家、进专属社群</text>
					<uni-icons type="right" size="16" color="#999999"></uni-icons>
				</view>
			</view>
			
			<!-- 常用功能卡片 -->
			<view class="card function-card">
				<view class="order-section">
					<view class="section-header">
						<text class="section-title">常用功能</text>
					</view>
					<view class="function-grid">
						<!-- 注释掉师傅入驻
						<view class="function-item" @tap="handleMasterJoin">
							<uni-icons type="person" size="28" color="#fc3e2b" custom-prefix="custom"></uni-icons>
							<text class="function-name">师傅入驻</text>
						</view>
						-->
						<!-- 注释掉费用补缴
						<view class="function-item">
							<uni-icons type="wallet" size="28" color="#fc3e2b" custom-prefix="custom"></uni-icons>
							<text class="function-name">费用补缴</text>
						</view>
						-->
						<!-- 注释掉停发查询
						<view class="function-item">
							<uni-icons type="search" size="28" color="#fc3e2b"></uni-icons>
							<text class="function-name">停发查询</text>
						</view>
						-->
						<view class="function-item">
							<uni-icons type="shop" size="28" color="#fc3e2b"></uni-icons>
							<text class="function-name">商务合作</text>
						</view>
						<view class="function-item" @tap="handleContact">
							<uni-icons type="headphones" size="28" color="#fc3e2b"></uni-icons>
							<text class="function-name">在线客服</text>
						</view>
						<view class="function-item">
							<uni-icons type="help" size="28" color="#fc3e2b"></uni-icons>
							<text class="function-name">帮助中心</text>
						</view>
						<view class="function-item" @tap="handleFeedback">
							<uni-icons type="chat" size="28" color="#fc3e2b"></uni-icons>
							<text class="function-name">意见反馈</text>
						</view>
					</view>
				</view>
			</view>
		</view>
		
		<tab-bar :current-tab="tabIndex"></tab-bar>
		
		<!-- 二维码弹窗 -->
		<view class="qrcode-popup" v-if="showQrCodePopup" @tap.stop="closeQrCode">
			<view class="qrcode-container" @tap.stop>
				<image class="qrcode-image" :src="butlerImageUrl || '/static/images/qrcode.png'" mode="aspectFit"></image>
				<text class="qrcode-tip">长按识别二维码添加管家</text>
				<view class="close-btn" @tap="closeQrCode">
					<uni-icons type="closeempty" size="24" color="#666666"></uni-icons>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import TabBar from '@/components/tab-bar/index.vue'
	import deviceMixin from '@/mixins/device.js'
	import { getFavoriteCount, getBrowseHistoryCount, getFollowingCount } from '@/apis/content.js'
	import { getWxLoginCode, wxappLogin, getButlerImage } from '@/apis/user.js'
	import { setToken } from '@/utils/storage.js'
	import { mapGetters, mapActions } from 'vuex'
	import { user, publish } from '@/apis/index.js'
	import deviceInfo from '@/utils/device-info.js'
	
	export default {
		components: {
			TabBar
		},
		mixins: [deviceMixin],
		data() {
			return {
				tabIndex: 4,
				statusBarHeight: 0,
				navBarHeight: 44,
				showQrCodePopup: false,
				favoriteCount: 0, // 收藏数量
				browseHistoryCount: 0, // 浏览记录数量
				followingCount: 0, // 关注数量
				publishCount: 0, // 发布数量
				butlerImageUrl: '', // 管家二维码图片
			}
		},
		computed: {
			// 从Vuex中获取用户信息和登录状态
			userInfo() {
				return this.$store.getters['user/userInfo'];
			},
			hasLogin() {
				return this.$store.getters['user/isLoggedIn'];
			},
			loginLoading() {
				return this.$store.state.user.loginLoading;
			}
		},
		onLoad() {
			// 获取状态栏高度
			this.statusBarHeight = deviceInfo.getStatusBarHeight();
			
			// 从本地存储加载上次选择的城市
			this.loadFavoriteData();
		},
		onShow() {
			this.tabIndex = 4;
			// 确保每次打开页面时刷新用户信息
			this.checkAndRefreshUserInfo();
			this.getFavoriteCount();
			this.getBrowseHistoryCount();
			this.getFollowingCount();
			this.getPublishCount();
		},
		methods: {
			// 处理登录方法
			async handleLogin() {
				try {
					// 显示加载提示
					uni.showLoading({
						title: '登录中',
						mask: true
					});
					
					// 1. 获取微信登录凭证code
					const code = await getWxLoginCode();
					
					// 2. 直接调用登录接口
					const loginData = { code };
					const loginResult = await wxappLogin(loginData);
					
					// 3. 处理登录结果
					if (loginResult && loginResult.code === 0 && loginResult.data && loginResult.data.token) {
						// 设置token到store
						this.$store.commit('user/SET_TOKEN', loginResult.data.token);
						
						// 获取用户信息
						await this.$store.dispatch('user/getUserInfo');
						
						uni.hideLoading();
						uni.showToast({
							title: '登录成功',
							icon: 'success'
						});
						
						// 登录成功后刷新页面数据
						this.getFavoriteCount();
						this.getBrowseHistoryCount();
						this.getFollowingCount();
					} else {
						uni.hideLoading();
						uni.showToast({
							title: loginResult.message || '登录失败',
							icon: 'none',
							duration: 2000
						});
					}
				} catch (error) {
					uni.hideLoading();
					console.error('登录错误:', error);
					// 显示错误信息
					uni.showToast({
						title: error.message || '登录失败，请重试',
						icon: 'none',
						duration: 2000
					});
				}
			},
			
			// 处理用户信息区域点击
			handleUserInfoClick() {
				if (!this.hasLogin) {
					// 未登录时，触发登录
					this.handleLogin();
				} else if (!this.userInfo.phone) {
					// 已登录但未绑定手机号，跳转到绑定手机号页面
					uni.showToast({
						title: '该功能暂未开放',
						icon: 'none'
					});
					// 实际开发中可以跳转到手机绑定页面
					// uni.navigateTo({
					//   url: '/pages/my/bind-phone/index'
					// });
				}
			},
			
			// 检查并刷新用户信息
			async checkAndRefreshUserInfo() {
				// 先检查登录状态
				const isLoggedIn = this.$store.getters['user/isLoggedIn'];
				
				if (!isLoggedIn) {
					return;
				}
				
				// 先检查已有的用户信息是否完整
				const userInfo = this.$store.getters['user/userInfo'];
				const hasUserId = userInfo && (userInfo.clientId || userInfo.id);
				
				// 如果已经有用户ID，跳过请求
				if (hasUserId && userInfo.realName) {
					return;
				}
				
				// 用户已登录，但没有完整的用户信息，请求最新数据
				try {
					await this.$store.dispatch('user/getUserInfo');
				} catch (error) {
					// 如果是401错误，尝试重新登录
					if (error.code === 401) {
						try {
							await this.$store.dispatch('user/silentLogin');
						} catch (loginError) {
							// 重新登录失败
						}
					}
				}
			},
			
			handleSetting() {
				uni.navigateTo({
					url: '/pages/settings/index'
				});
			},
			
			handleFeatureClick(type) {
				// 判断是否登录
				if (!this.hasLogin) {
					uni.showToast({
						title: '请先登录',
						icon: 'none'
					});
					return;
				}
				
				const pages = {
					publish: '/pages/my/publish/index',
					follow: '/pages/my/follow/index',
					favorite: '/pages/my/favorite/index',
					history: '/pages/my/history/index'
				};
				
				if (pages[type]) {
					uni.navigateTo({
						url: pages[type]
					});
				} else {
					uni.showToast({
						title: '功能开发中',
						icon: 'none'
					});
				}
			},
			
			handleSignIn() {
				// 判断是否登录
				if (!this.hasLogin) {
					uni.showToast({
						title: '请先登录',
						icon: 'none'
					});
					return;
				}
				
				// 跳转到音乐会员福利页面
				uni.navigateTo({
					url: '/pages/vip/music/index'
				});
			},
			
			handleViewAllOrders() {
				// 判断是否登录
				if (!this.hasLogin) {
					uni.showToast({
						title: '请先登录',
						icon: 'none'
					});
					return;
				}
				
				uni.navigateTo({
					url: '/pages/my/orders/index'
				});
			},
			
			handleOrderType(type) {
				// 判断是否登录
				if (!this.hasLogin) {
					uni.showToast({
						title: '请先登录',
						icon: 'none'
					});
					return;
				}
				
				// 根据不同类型跳转到对应的订单列表页面
				let tabIndex = 0;
				switch(type) {
					case 'processing':
						tabIndex = 1;
						break;
					case 'unpaid':
						tabIndex = 2;
						break;
					case 'completed':
						tabIndex = 3;
						break;
					default:
						tabIndex = 0;
				}
				uni.navigateTo({
					url: `/pages/my/orders/index?tab=${tabIndex}`
				});
			},
			
			showQrCode() {
				// 显示二维码弹窗
				this.showQrCodePopup = true;
				
				// 获取二维码图片
				this.getButlerImage();
			},
			
			closeQrCode() {
				this.showQrCodePopup = false;
			},
			
			// 获取管家二维码图片
			async getButlerImage() {
				try {
					const res = await getButlerImage();
					if (res.code === 0 && res.data && res.data.imageUrl) {
						this.butlerImageUrl = res.data.imageUrl;
					} else {
						console.error('获取管家二维码失败', res);
					}
				} catch (error) {
					console.error('获取管家二维码出错', error);
				}
			},
			
			handleMasterJoin() {
				// 判断是否登录
				if (!this.hasLogin) {
					uni.showToast({
						title: '请先登录',
						icon: 'none'
					});
					return;
				}
				
				uni.navigateTo({
					url: '/pages/master/join/index'
				});
			},
			
			// 处理联系客服
			handleContact() {
				// #ifdef MP-WEIXIN
				if (wx.openCustomerServiceChat) {
					wx.openCustomerServiceChat({
						extInfo: {url: 'https://work.weixin.qq.com/'}, // 根据实际情况修改
						corpId: 'ww5823288888ed1111', // 需要替换为企业微信的企业ID
						success(res) {
							console.log('客服会话打开成功', res);
						},
						fail(err) {
							console.error('客服会话打开失败', err);
							// 降级方案，打开系统或小程序设置引导用户操作
							uni.showModal({
								title: '在线客服',
								content: '无法连接客服，您可以通过右上角胶囊按钮 -> 关于 -> 客服来联系我们',
								confirmText: '我知道了',
								showCancel: false
							});
						}
					});
				} else {
					// 旧版微信不支持此API，给用户提示
					uni.showModal({
						title: '在线客服',
						content: '您可以通过右上角胶囊按钮 -> 关于 -> 客服来联系我们',
						confirmText: '我知道了',
						showCancel: false
					});
				}
				// #endif
				
				// #ifndef MP-WEIXIN
				uni.showToast({
					title: '该功能仅支持微信小程序',
					icon: 'none'
				});
				// #endif
			},
			
			// 处理意见反馈
			handleFeedback() {
				// #ifdef MP-WEIXIN
				if (typeof wx.showFeedback === 'function') {
					// 使用官方API打开意见反馈页面
					wx.showFeedback();
				} else {
					// 降级处理：直接用按钮替代
					uni.showModal({
						title: '意见反馈',
						content: '您可以通过"右上角胶囊按钮 -> 关于 -> 反馈与投诉"来提交反馈',
						confirmText: '我知道了',
						showCancel: false
					});
				}
				// #endif
				
				// #ifndef MP-WEIXIN
				uni.showToast({
					title: '该功能仅支持微信小程序',
					icon: 'none'
				});
				// #endif
			},
			
			async getFavoriteCount() {
				// 检查是否登录
				if (!this.hasLogin) {
					this.favoriteCount = 0;
					return;
				}
				
				try {
					const res = await getFavoriteCount();
					if (res.code === 0 && res.data) {
						this.favoriteCount = res.data.total || 0;
					} else {
						this.favoriteCount = 0;
					}
				} catch (error) {
					console.error('获取收藏数量失败', error);
					this.favoriteCount = 0;
				}
			},
			
			async getBrowseHistoryCount() {
				// 检查是否登录
				if (!this.hasLogin) {
					this.browseHistoryCount = 0;
					return;
				}
				
				try {
					const res = await getBrowseHistoryCount();
					if (res.code === 0 && res.data) {
						this.browseHistoryCount = res.data.count || 0;
					} else {
						this.browseHistoryCount = 0;
					}
				} catch (error) {
					console.error('获取浏览记录数量失败', error);
					this.browseHistoryCount = 0;
				}
			},
			
			async getFollowingCount() {
				// 检查是否登录
				if (!this.hasLogin) {
					this.followingCount = 0;
					return;
				}
				
				try {
					const res = await getFollowingCount();
					if (res.code === 0 && res.data) {
						this.followingCount = res.data.count || 0;
					} else {
						this.followingCount = 0;
					}
				} catch (error) {
					console.error('获取关注数量失败', error);
					this.followingCount = 0;
				}
			},
			
			async getPublishCount() {
				// 检查是否登录
				if (!this.hasLogin) {
					this.publishCount = 0;
					return;
				}
				
				try {
					const res = await publish.getPublishCount();
					if (res.code === 0 && res.data) {
						this.publishCount = res.data.total || 0;
					} else {
						this.publishCount = 0;
					}
				} catch (error) {
					console.error('获取发布数量失败', error);
					this.publishCount = 0;
				}
			},
			
			loadFavoriteData() {
				// 实现从本地存储加载上次选择的城市
			}
		}
	}
</script>

<style>
	.my-container {
		min-height: 100vh;
		background-color: #f5f5f5;
		padding-bottom: env(safe-area-inset-bottom);
		position: relative;
	}
	
	.gradient-bg {
		background: linear-gradient(180deg, rgba(252, 62, 43, 0.95) 0%, rgba(255, 100, 51, 0.9) 100%);
		padding-bottom: 45rpx;
	}
	
	.content-card {
		position: relative;
		margin-top: -35rpx;
		background-color: #ffffff;
		min-height: calc(100vh - 350rpx);
		border-radius: 40rpx 40rpx 0 0;
		padding: 20rpx 30rpx;
		padding-bottom: calc(120rpx + env(safe-area-inset-bottom));
		z-index: 1;
		box-shadow: 0 -2rpx 8rpx rgba(0, 0, 0, 0.02);
	}
	
	.status-bar {
		width: 100%;
	}
	
	.nav-bar {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
	}
	
	.nav-left {
		position: absolute;
		left: 30rpx;
		display: flex;
		align-items: center;
		z-index: 1;
	}
	
	.nav-title {
		text-align: center;
		font-size: 34rpx;
		font-weight: 500;
		color: #ffffff;
	}
	
	.user-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 20rpx 30rpx;
	}
	
	.user-info {
		display: flex;
		align-items: center;
		flex: 1;
		margin-right: 20rpx;
	}
	
	.avatar {
		width: 100rpx;
		height: 100rpx;
		border-radius: 50%;
		margin-right: 24rpx;
		border: 4rpx solid rgba(255, 255, 255, 0.3);
	}
	
	.user-detail {
		display: flex;
		flex-direction: column;
	}
	
	.nickname {
		font-size: 32rpx;
		font-weight: 500;
		color: #ffffff;
		margin-bottom: 8rpx;
	}
	
	.phone {
		font-size: 26rpx;
		color: rgba(255, 255, 255, 0.8);
	}
	
	.white {
		color: #ffffff;
	}
	
	.white-light {
		color: rgba(255, 255, 255, 0.8);
	}
	
	.setting-btn {
		width: 80rpx;
		height: 80rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.feature-bar {
		display: flex;
		justify-content: space-around;
	}
	
	.feature-wrapper {
		margin: 0 30rpx;
		padding: 30rpx 20rpx;
		background: rgba(255, 255, 255, 0.25);
		border-radius: 16rpx;
		backdrop-filter: blur(10px);
	}
	
	.feature-item {
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
	.feature-num {
		font-size: 32rpx;
		font-weight: 500;
		color: #ffffff;
		margin-bottom: 8rpx;
	}
	
	.feature-name {
		font-size: 24rpx;
		color: rgba(255, 255, 255, 0.8);
	}
	
	.sign-btn {
		display: flex;
		align-items: center;
		padding: 12rpx 20rpx;
		background: rgba(255, 255, 255, 0.2);
		border-radius: 30rpx;
		border: 1rpx solid rgba(255, 255, 255, 0.3);
	}
	
	.sign-text {
		font-size: 24rpx;
		color: #ffffff;
		margin-left: 8rpx;
	}
	
	.card {
		background-color: #ffffff;
		border-radius: 24rpx;
		box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
		border: 1rpx solid rgba(0, 0, 0, 0.03);
	}
	
	.order-section {
		padding: 20rpx;
	}
	
	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;
	}
	
	.section-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #000000;
	}
	
	.more-btn {
		display: flex;
		align-items: center;
	}
	
	.more-text {
		font-size: 26rpx;
		color: #999999;
		margin-right: 4rpx;
	}
	
	.order-types {
		display: flex;
		justify-content: space-around;
		padding: 10rpx 0;
		width: 100%;
		margin: 0 auto;
	}
	
	.order-type-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 12rpx;
		flex: 1;
	}
	
	.type-name {
		font-size: 26rpx;
		color: #000000;
		margin-top: 12rpx;
	}
	
	.manager-card {
		margin-top: 30rpx;
	}
	
	.manager-content {
		padding: 30rpx;
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	
	.manager-text {
		font-size: 30rpx;
		color: #000000;
	}
	
	.function-card {
		margin-top: 30rpx;
	}
	
	.function-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		padding: 10rpx 0;
		gap: 0;
	}
	
	.function-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 10rpx 0;
	}
	
	.function-name {
		font-size: 26rpx;
		color: #000000;
		margin-top: 8rpx;
	}
	
	.qrcode-popup {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.6);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 999;
	}
	
	.qrcode-container {
		width: 600rpx;
		background-color: #ffffff;
		border-radius: 20rpx;
		padding: 40rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		position: relative;
	}
	
	.qrcode-image {
		width: 400rpx;
		height: 400rpx;
		margin-bottom: 30rpx;
	}
	
	.qrcode-tip {
		font-size: 28rpx;
		color: #666666;
		margin-bottom: 20rpx;
	}
	
	.close-btn {
		position: absolute;
		top: 20rpx;
		right: 20rpx;
		width: 60rpx;
		height: 60rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	/* 添加点击样式 */
	.user-info:active {
		opacity: 0.8;
	}
</style>
