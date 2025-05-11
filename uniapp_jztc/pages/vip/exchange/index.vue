<template>
	<view class="container">
		<!-- 状态栏占位 -->
		<view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
		
		<!-- 自定义导航栏 -->
		<view class="nav-bar">
			<view class="nav-content">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#ffecd8"></uni-icons>
				</view>
				<text class="nav-title">兑换视频会员权益</text>
			</view>
		</view>

		<!-- 滚动提示区域 -->
		<view class="notice-area">
			<view class="notice-scroll">
				<view class="notice-text">
					<template v-if="recentExchanges.length > 0">
						<template v-for="(item, index) in recentExchanges" :key="index">
							<view class="notice-item">{{item.clientName}}兑换了{{item.productName}}</view>
							<text class="separator" v-if="index < recentExchanges.length - 1">|</text>
						</template>
					</template>
					<template v-else>
						<view class="notice-item">暂无兑换记录</view>
					</template>
				</view>
				<view class="notice-text" v-if="recentExchanges.length > 0">
					<template v-for="(item, index) in recentExchanges" :key="index">
						<view class="notice-item">{{item.clientName}}兑换了{{item.productName}}</view>
						<text class="separator" v-if="index < recentExchanges.length - 1">|</text>
					</template>
				</view>
			</view>
		</view>

		<!-- 会员时长卡片 -->
		<view class="vip-info-card">
			<view class="card-header">
				<text class="card-title">我的时长</text>
				<view class="exchange-record" @tap="handleExchangeRecord">
					<text>兑换记录</text>
				</view>
			</view>
			<view class="time-info">
				<text class="days">{{ memberDays }}</text>
				<text class="unit">天</text>
				<text class="days">{{ memberHours }}</text>
				<text class="unit">小时</text>
				<text class="days">{{ memberMinutes }}</text>
				<text class="unit">分钟</text>
			</view>
		</view>
		
		<!-- 分类选项卡 -->
		<view class="tab-container">
			<scroll-view class="tab-scroll" scroll-x :show-scrollbar="false">
				<view class="tab-list">
					<view 
						class="tab-item" 
						v-for="(item, index) in tabList" 
						:key="index"
						:class="{ active: currentTab === index }"
						@tap="switchTab(index)"
					>
						<text>{{ item }}</text>
					</view>
				</view>
			</scroll-view>
		</view>
		
		<!-- 会员卡片列表 -->
		<view class="vip-card-list">
			<block v-if="products.length > 0">
				<view class="vip-card" v-for="item in products" :key="item.id">
					<view class="card-tag" v-if="item.tags">{{item.tags}}</view>
					<view class="card-content">
						<view class="card-left">
							<image class="card-logo" :src="item.thumbnail" mode="aspectFit"></image>
							<view class="card-info">
								<text class="card-name">{{item.name}}</text>
								<view class="card-type">{{item.description}}</view>
							</view>
						</view>
						<view class="card-right">
							<view class="exchange-info">
								<text class="exchange-text">兑换</text>
								<text class="need-text">需{{item.duration}}天会员时长</text>
							</view>
							<view class="exchange-btn-small" @tap="handleExchangeVip(item)">
								<text>去兑换</text>
							</view>
						</view>
					</view>
				</view>
			</block>
			
			<!-- 没有产品时显示提示 -->
			<view v-if="products.length === 0" class="empty-tip">
				<text>暂无会员卡可兑换</text>
			</view>
		</view>
		
		<!-- 使用全局抽屉组件 -->
		<exchange-popup ref="exchangePopup" @confirm="handleConfirmExchange"></exchange-popup>
	</view>
</template>

<script>
	import deviceMixin from '@/mixins/device.js'
	import shareMixin from '@/mixins/share.js'
	import { get, post } from '@/utils/request.js'
	import exchangePopup from '@/components/exchange-popup/index.vue'
	import uniIcons from '@dcloudio/uni-ui/lib/uni-icons/uni-icons'
	import { vip } from '@/apis/index.js'
	import { getShareSettings } from '@/utils/share.js'
	
	export default {
		components: {
			uniIcons,
			exchangePopup
		},
		mixins: [deviceMixin, shareMixin],
		data() {
			return {
				recentExchanges: [], // 最近兑换记录
				currentTab: 0,
				tabList: ['全部'], // 初始只有全部选项，后面从接口获取
				memberDays: '0',
				memberHours: '0',
				memberMinutes: '0',
				categories: [], // 保存完整的分类数据，包括ID
				products: [], // 保存产品列表数据
				page: 1,
				size: 10,
				currentCategoryId: null, // 当前选中的分类ID
				currentProduct: null, // 当前选中的产品
				shareImageUrl: '' // 保存分享图片地址
			}
		},
		onLoad() {
			// 获取会员时长
			this.getMemberDuration()
			// 获取平台分类列表
			this.getCategoryList()
			// 获取最近兑换记录
			this.getRecentExchangeRecords()
			// 获取分享图片
			this.getShareImage()
		},
		// 自定义分享给好友
		onShareAppMessage() {
			return {
				title: "兑换视频会员权益", // 使用页面导航栏标题
				path: '/pages/vip/exchange/index',
				imageUrl: this.shareImageUrl
			}
		},
		// 自定义分享到朋友圈
		onShareTimeline() {
			return {
				title: "兑换视频会员权益", // 使用页面导航栏标题
				query: '',
				imageUrl: this.shareImageUrl
			}
		},
		methods: {
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
			},
			handleBack() {
				uni.navigateBack()
			},
			handleExchangeRecord() {
				// 跳转到兑换记录页面
				uni.navigateTo({
					url: '/pages/vip/exchange/record'
				})
			},
			// 获取最近兑换记录
			async getRecentExchangeRecords() {
				try {
					const res = await vip.getRecentExchangeRecords()
					if (res.code === 0 && res.data) {
						this.recentExchanges = res.data.list || []
					} else {
						console.error('获取最近兑换记录失败:', res?.message)
					}
				} catch (error) {
					console.error('获取最近兑换记录异常:', error)
				}
			},
			switchTab(index) {
				this.currentTab = index
				
				// 设置当前分类ID
				if (index === 0) {
					// 全部分类
					this.currentCategoryId = null
				} else if (this.categories.length > 0 && index - 1 < this.categories.length) {
					// 特定分类
					this.currentCategoryId = this.categories[index - 1].id
				}
				
				// 重新加载产品列表
				this.getProductList()
			},
			handleExchangeVip(item) {
				// 保存当前选中的产品
				this.currentProduct = item
				
				// 打开兑换弹窗并传递产品信息
				this.$nextTick(() => {
					if (this.$refs.exchangePopup) {
						this.$refs.exchangePopup.open(item)
					} else {
						uni.showToast({
							title: '组件加载中，请稍后再试',
							icon: 'none'
						})
					}
				})
			},
			// 处理确认兑换
			handleConfirmExchange(data) {
				// 执行兑换
				this.exchangeProduct(data)
			},
			async exchangeProduct(data) {
				if (!data) return
				
				try {
					uni.showLoading({ title: '兑换中...' })
					
					// 发送兑换请求 - 使用vip API
					const res = await vip.createExchangeRecord(data)
					
					uni.hideLoading()
					
					if (res.code === 0) {
						uni.showToast({
							title: '兑换成功',
							icon: 'success',
							duration: 2000
						})
						
						// 刷新会员时长和最近兑换记录
						setTimeout(() => {
							this.getMemberDuration()
							this.getRecentExchangeRecords()
						}, 1000)
					} else {
						uni.showToast({
							title: res.msg || '兑换失败',
							icon: 'none'
						})
					}
				} catch (error) {
					uni.hideLoading()
					console.error('兑换失败:', error)
					uni.showToast({
						title: '兑换失败，请重试',
						icon: 'none'
					})
				}
			},
			// 获取会员时长
			async getMemberDuration() {
				try {
					const res = await get('/wx/client/duration/remaining')
					if (res.code === 0 && res.data) {
						// 获取到会员时长信息
						const remainingDuration = res.data.remainingDuration || '0天0小时0分钟'
						
						// 提取天数
						this.memberDays = remainingDuration.match(/(\d+)天/) ? remainingDuration.match(/(\d+)天/)[1] : '0'
						this.memberHours = remainingDuration.match(/(\d+)小时/) ? remainingDuration.match(/(\d+)小时/)[1] : '0'
						this.memberMinutes = remainingDuration.match(/(\d+)分钟/) ? remainingDuration.match(/(\d+)分钟/)[1] : '0'
					} else {
						console.error('获取会员时长失败:', res?.message)
					}
				} catch (error) {
					console.error('获取会员时长异常:', error)
				}
			},
			// 获取平台分类列表
			async getCategoryList() {
				try {
					const res = await get('/wx/client/shop-category/list')
					if (res.code === 0 && res.data && res.data.list && res.data.list.length > 0) {
						// 保存完整的分类数据
						this.categories = res.data.list
						
						// 保留"全部"选项，添加从接口获取的分类
						this.tabList = ['全部', ...res.data.list.map(item => item.name)]
						
						// 获取产品列表（默认全部）
						this.getProductList()
					}
				} catch (error) {
					console.error('获取平台分类列表失败:', error)
				}
			},
			// 获取产品列表
			async getProductList() {
				try {
					uni.showLoading({ title: '加载中...' })
					
					// 构建请求参数
					const params = {
						page: this.page,
						size: this.size
					}
					
					// 如果有选中特定分类，添加分类ID
					if (this.currentCategoryId) {
						params.categoryId = this.currentCategoryId
					}
					
					const res = await get('/wx/product/list', params)
					if (res.code === 0 && res.data) {
						this.products = res.data.list || []
					} else {
						console.error('获取产品列表失败:', res?.message)
						this.products = []
					}
				} catch (error) {
					console.error('获取产品列表异常:', error)
					this.products = []
				} finally {
					uni.hideLoading()
				}
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background: #000000;
	}
	
	.status-bar {
		background: #000000;
	}
	
	.nav-bar {
		background: #000000;
	}
	
	.nav-content {
		height: 44px;
		display: flex;
		align-items: center;
		padding: 0 30rpx;
		position: relative;
	}
	
	.back-btn {
		position: absolute;
		left: 30rpx;
		height: 100%;
		display: flex;
		align-items: center;
	}
	
	.nav-title {
		color: #ffecd8;
		font-size: 32rpx;
		font-weight: bold;
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
	}
	
	.notice-area {
		padding: 20rpx 30rpx;
		overflow: hidden;
	}
	
	.notice-scroll {
		width: 100%;
		height: 60rpx;
		overflow: hidden;
		position: relative;
	}
	
	.notice-text {
		display: flex;
		align-items: center;
		height: 60rpx;
		position: absolute;
		white-space: nowrap;
		animation: scrollText 15s linear infinite;
		left: 100%;
	}
	
	.notice-text:nth-child(1) {
		animation-delay: 0s;
	}
	
	.notice-text:nth-child(2) {
		animation-delay: 7.5s;
	}
	
	@keyframes scrollText {
		0% {
			transform: translateX(0%);
		}
		100% {
			transform: translateX(-200%);
		}
	}
	
	.notice-item {
		display: inline-block;
		font-size: 24rpx;
		color: #ffecd8;
		margin-right: 10rpx;
	}
	
	.separator {
		display: inline-block;
		font-size: 24rpx;
		color: #ffecd8;
		margin: 0 10rpx;
	}
	
	.vip-info-card {
		margin: 20rpx 30rpx;
		background: linear-gradient(to right, #ffd4c2, #ffbeb0);
		border-radius: 16rpx;
		padding: 20rpx;
		position: relative;
	}
	
	.card-title {
		font-size: 28rpx;
		color: #a04a39;
		font-weight: bold;
	}
	
	.exchange-record {
		background: rgba(255, 255, 255, 0.6);
		border-radius: 20rpx;
		padding: 4rpx 16rpx;
		display: flex;
		align-items: center;
		height: 40rpx;
		position: absolute;
		right: 20rpx;
		top: 50%;
		transform: translateY(-50%);
	}
	
	.exchange-record text {
		font-size: 24rpx;
		color: #862c13;
		line-height: 1;
	}
	
	.time-info {
		display: flex;
		align-items: flex-end;
		margin-top: 20rpx;
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
		margin-bottom: 8rpx;
	}
	
	/* 分类选项卡样式 */
	.tab-container {
		margin: 20rpx 0;
		background: rgba(255, 255, 255, 0.05);
	}
	
	.tab-scroll {
		white-space: nowrap;
	}
	
	.tab-list {
		display: flex;
		padding: 0 20rpx;
	}
	
	.tab-item {
		display: inline-block;
		padding: 20rpx 30rpx;
		position: relative;
	}
	
	.tab-item text {
		font-size: 28rpx;
		color: #999999;
	}
	
	.tab-item.active text {
		color: #ffecd8;
		font-weight: bold;
	}
	
	.tab-item.active::after {
		content: '';
		position: absolute;
		bottom: 10rpx;
		left: 50%;
		transform: translateX(-50%);
		width: 40rpx;
		height: 4rpx;
		background-color: #ffecd8;
		border-radius: 2rpx;
	}
	
	/* 会员卡片列表样式 */
	.vip-card-list {
		padding: 0 30rpx;
	}
	
	.vip-card {
		position: relative;
		background: linear-gradient(to right, #f8f8f8, #ffffff);
		border-radius: 16rpx;
		margin-bottom: 20rpx;
		overflow: hidden;
	}
	
	.card-tag {
		position: absolute;
		left: 0;
		top: 0;
		background: #fc3e2b;
		color: #ffffff;
		font-size: 22rpx;
		padding: 4rpx 12rpx;
		border-radius: 0 0 12rpx 0;
		z-index: 1;
	}
	
	.card-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20rpx;
	}
	
	.card-left {
		display: flex;
		align-items: center;
	}
	
	.card-logo {
		width: 80rpx;
		height: 80rpx;
		margin-right: 16rpx;
		border-radius: 12rpx;
	}
	
	.card-info {
		display: flex;
		flex-direction: column;
	}
	
	.card-name {
		font-size: 28rpx;
		color: #333333;
		font-weight: bold;
		margin-bottom: 8rpx;
	}
	
	.card-type {
		display: inline-block;
		font-size: 22rpx;
		color: #666666;
		background: #f0f0f0;
		padding: 2rpx 12rpx;
		border-radius: 10rpx;
	}
	
	.card-right {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
	}
	
	.exchange-info {
		display: flex;
		align-items: center;
		margin-bottom: 12rpx;
	}
	
	.exchange-text {
		font-size: 24rpx;
		color: #333333;
		margin-right: 8rpx;
	}
	
	.need-text {
		font-size: 24rpx;
		color: #fc3e2b;
	}
	
	.exchange-btn-small {
		background: linear-gradient(135deg, #fc3e2b 0%, #fa7154 100%);
		border-radius: 30rpx;
		padding: 8rpx 30rpx;
	}
	
	.exchange-btn-small text {
		font-size: 24rpx;
		color: #ffffff;
	}
	
	.empty-tip {
		text-align: center;
		padding: 20rpx;
		color: #ffecd8;
	}
</style> 