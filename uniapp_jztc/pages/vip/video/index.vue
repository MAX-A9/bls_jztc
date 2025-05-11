<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#ffecd8"></uni-icons>
				</view>
				<text class="nav-title">0元领视频会员</text>
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
						<image class="avatar" src="/static/avatar/default.png" mode="aspectFill"></image>
					</view>
					<text class="user-title">我的会员时长</text>
				</view>
			</view>
			<view class="time-info">
				<text class="days">3</text>
				<text class="unit">天</text>
				<uni-icons type="right" size="14" color="#862c13" class="time-arrow"></uni-icons>
				<view class="exchange-btn">
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
				<view class="platform-item" v-for="(item, index) in platforms" :key="index">
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
				<view class="package-item" 
					v-for="(item, index) in packages" 
					:key="index"
				>
					<view class="package-tag">{{item.tag}}</view>
					<view class="package-content">
						<view class="name-row">
							<text class="package-name">{{item.name}}</text>
							<view class="month-tag">月卡</view>
						</view>
						<view class="price-info">
							<view class="price-row">
								<text class="price-symbol">¥</text>
								<text class="package-price">{{item.price}}</text>
							</view>
							<text class="duration-text">+{{item.days}}天时长</text>
						</view>
					</view>
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
	
	export default {
		mixins: [deviceMixin],
		data() {
			return {
				platforms: [
					{ name: '爱奇艺', icon: '/static/platform/iqy.png' },
					{ name: '腾讯视频', icon: '/static/platform/tx.png' },
					{ name: '优酷', icon: '/static/platform/yk.png' },
					{ name: '芒果TV', icon: '/static/platform/mg.png' },
					{ name: '哔哩哔哩', icon: '/static/platform/bl.png' },
					{ name: '搜狐视频', icon: '/static/platform/sh.png' }
				],
				packages: [
					{
						tag: '特惠专享',
						name: '爱奇艺VIP',
						price: '0.01',
						days: '30'
					},
					{
						tag: '全网低价',
						name: '腾讯视频VIP',
						price: '9.66',
						days: '1'
					},
					{
						tag: '6.1折',
						name: '优酷视频VIP',
						price: '11',
						days: '30'
					}
				]
			}
		},
		methods: {
			handleBack() {
				uni.navigateBack()
			},
			handleWatch() {
				uni.showLoading({
					title: '加载中...'
				})
				setTimeout(() => {
					uni.hideLoading()
					uni.showToast({
						title: '领取成功',
						icon: 'success'
					})
				}, 1000)
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
</style> 