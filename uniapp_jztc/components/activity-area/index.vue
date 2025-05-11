<template>
	<view class="activity-area" v-if="isGlobalEnabled">
		<view class="activity-grid" :class="{ 'small-screen': isSmallScreen, 'large-screen': isLargeScreen }">
			<!-- 左侧两个卡片 -->
			<view class="activity-left">
				<view 
					class="activity-card small" 
					v-for="(item, index) in leftItems" 
					:key="index" 
					@click="handleItemClick(item)"
				>
					<view :class="['card-pattern', getPatternClass(item.position)]"></view>
					<view class="card-content">
						<text class="card-title" :class="{ 'small-text': isSmallScreen }">{{ item.title }}</text>
						<text class="card-desc" :class="{ 'small-text': isSmallScreen }">{{ item.description }}</text>
					</view>
				</view>
			</view>
			
			<!-- 右侧大卡片 -->
			<view class="activity-right" v-if="rightItem">
				<view class="activity-card large" @click="handleItemClick(rightItem)">
					<view class="card-pattern video large"></view>
					<view class="card-content">
						<view class="title-wrapper">
							<text class="card-title" :class="{ 'small-text': isSmallScreen }">{{ rightItem.title }}</text>
							<view class="no-money-wrapper">
								<text class="no-money-text">免费</text>
							</view>
						</view>
						<text class="card-desc" :class="{ 'small-text': isSmallScreen }">{{ rightItem.description }}</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import { get } from '@/utils/request.js';
	import deviceAdapter from '@/mixins/device-adapter.js';
	
	export default {
		name: 'ActivityArea',
		mixins: [deviceAdapter],
		data() {
			return {
				activityList: [],
				isGlobalEnabled: false,
				leftItems: [],
				rightItem: null
			}
		},
		created() {
			// 移除自动加载，由父组件控制加载时机
			// this.fetchActivityData();
		},
		methods: {
			async fetchActivityData() {
				try {
					const result = await get('/wx/activity-area/get');
					if (result && result.code === 0) {
						this.isGlobalEnabled = result.data.isGlobalEnabled;
						if (this.isGlobalEnabled && result.data.list) {
							this.activityList = result.data.list;
							this.processActivityData();
						}
					}
				} catch (error) {
					console.error('获取活动区域数据失败', error);
				}
			},
			processActivityData() {
				// 分类处理活动数据
				this.leftItems = this.activityList.filter(item => 
					item.position === 'topLeft' || item.position === 'bottomLeft'
				);
				this.rightItem = this.activityList.find(item => item.position === 'right');
			},
			getPatternClass(position) {
				// 根据位置返回对应的样式类
				if (position === 'topLeft') return 'music';
				if (position === 'bottomLeft') return 'coupon';
				if (position === 'right') return 'video';
				return '';
			},
			handleItemClick(item) {
				if (item.linkType === 'page' && item.linkUrl) {
					uni.navigateTo({
						url: '/' + item.linkUrl
					});
				} else if (item.linkType === 'webview' && item.linkUrl) {
					// 处理外部链接
					uni.navigateTo({
						url: '/pages/webview/index?url=' + encodeURIComponent(item.linkUrl)
					});
				} else if (item.linkType === 'miniprogram' && item.linkUrl) {
					// 跳转到其他小程序
					uni.navigateToMiniProgram({
						appId: item.linkUrl,
						fail: (err) => {
							console.error('小程序跳转失败', item.linkUrl, err);
							uni.showToast({
								title: '跳转失败',
								icon: 'none'
							});
						}
					});
				}
			}
		}
	}
</script>

<style>
	.activity-area {
		padding: 20rpx;
		background-color: #ffffff;
	}
	
	.activity-grid {
		display: flex;
		gap: 20rpx;
	}
	
	/* 小屏幕适配 */
	.activity-grid.small-screen {
		gap: 10rpx;
	}
	
	/* 大屏幕适配 */
	.activity-grid.large-screen {
		gap: 30rpx;
		max-width: 1200rpx;
		margin: 0 auto;
	}
	
	.activity-left {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 20rpx;
	}
	
	.activity-grid.small-screen .activity-left {
		gap: 10rpx;
	}
	
	.activity-grid.large-screen .activity-left {
		gap: 30rpx;
	}
	
	.activity-right {
		flex: 1;
	}
	
	.activity-card {
		background: #f8f8f8;
		border-radius: 16rpx;
		position: relative;
		overflow: hidden;
		transition: transform 0.2s ease;
	}
	
	.activity-card:active {
		transform: scale(0.98);
	}
	
	.activity-card.small {
		height: 160rpx;
	}
	
	.activity-grid.small-screen .activity-card.small {
		height: 140rpx;
	}
	
	.activity-grid.large-screen .activity-card.small {
		height: 180rpx;
	}
	
	.activity-card.large {
		height: 340rpx;
		height: calc(160rpx * 2 + 20rpx);
	}
	
	.activity-grid.small-screen .activity-card.large {
		height: calc(140rpx * 2 + 10rpx);
	}
	
	.activity-grid.large-screen .activity-card.large {
		height: calc(180rpx * 2 + 30rpx);
	}
	
	.card-content {
		position: absolute;
		top: 20rpx;
		left: 20rpx;
		z-index: 1;
	}
	
	.activity-grid.small-screen .card-content {
		top: 10rpx;
		left: 10rpx;
	}
	
	.activity-grid.large-screen .card-content {
		top: 30rpx;
		left: 30rpx;
	}
	
	.card-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333333;
		display: block;
		margin-bottom: 8rpx;
	}
	
	.card-title.small-text {
		font-size: 28rpx;
	}
	
	.right-column .card-title {
		display: inline-block;
		margin-bottom: 0;
		line-height: 1;
	}
	
	.arrow-wrapper {
		margin-top: 8rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24rpx;
		height: 24rpx;
		background: rgba(102, 102, 102, 0.1);
		border-radius: 50%;
	}
	
	.no-money-wrapper {
		background: linear-gradient(to right, #ff6b6b, #ff3366);
		padding: 4rpx 12rpx;
		border-radius: 8rpx;
		display: inline-flex;
		align-items: center;
		transform: scale(0.95);
		line-height: 1;
	}
	
	.no-money-text {
		color: #ffffff;
		font-size: 24rpx;
		line-height: 1;
	}
	
	.title-wrapper {
		display: flex;
		align-items: center;
		gap: 8rpx;
		line-height: 1;
	}
	
	.card-pattern {
		position: absolute;
		right: 0;
		bottom: 0;
		width: 120rpx;
		height: 120rpx;
		background: linear-gradient(135deg, transparent 40%, rgba(252, 62, 43, 0.15) 40%);
		border-radius: 0 0 16rpx 0;
		opacity: 0.8;
	}
	
	.activity-grid.small-screen .card-pattern {
		width: 100rpx;
		height: 100rpx;
	}
	
	.activity-grid.large-screen .card-pattern {
		width: 140rpx;
		height: 140rpx;
	}
	
	.card-pattern.music {
		background: linear-gradient(135deg, transparent 40%, rgba(26, 173, 25, 0.15) 40%);
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.card-pattern.music::before {
		content: 'VIP';
		position: absolute;
		right: 20rpx;
		bottom: 20rpx;
		color: rgba(26, 173, 25, 0.5);
		font-size: 32rpx;
		font-weight: bold;
		transform: rotate(-45deg);
		animation: pulse 2s infinite;
	}
	
	.activity-grid.small-screen .card-pattern.music::before {
		font-size: 28rpx;
		right: 16rpx;
		bottom: 16rpx;
	}
	
	@keyframes pulse {
		0% {
			opacity: 0.5;
		}
		50% {
			opacity: 1;
		}
		100% {
			opacity: 0.5;
		}
	}
	
	.card-pattern.video.large {
		width: 160rpx;
		height: 160rpx;
		background: linear-gradient(135deg, transparent 40%, rgba(0, 161, 214, 0.15) 40%);
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.activity-grid.small-screen .card-pattern.video.large {
		width: 140rpx;
		height: 140rpx;
	}
	
	.activity-grid.large-screen .card-pattern.video.large {
		width: 180rpx;
		height: 180rpx;
	}
	
	.card-pattern.video::before {
		content: '视频';
		position: absolute;
		right: 30rpx;
		bottom: 30rpx;
		color: rgba(0, 161, 214, 0.5);
		font-size: 36rpx;
		font-weight: bold;
		transform: rotate(-45deg);
	}
	
	.activity-grid.small-screen .card-pattern.video::before {
		font-size: 32rpx;
		right: 24rpx;
		bottom: 24rpx;
	}
	
	.card-pattern.coupon {
		background: linear-gradient(135deg, transparent 40%, rgba(255, 102, 0, 0.15) 40%);
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.card-pattern.coupon::before {
		content: '兑';
		position: absolute;
		right: 20rpx;
		bottom: 20rpx;
		color: rgba(255, 102, 0, 0.5);
		font-size: 32rpx;
		font-weight: bold;
		transform: rotate(-45deg);
	}
	
	.activity-grid.small-screen .card-pattern.coupon::before {
		font-size: 28rpx;
		right: 16rpx;
		bottom: 16rpx;
	}
	
	.card-desc {
		font-size: 24rpx;
		color: #666666;
	}
	
	.card-desc.small-text {
		font-size: 22rpx;
	}
</style>