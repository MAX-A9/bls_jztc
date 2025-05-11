<template>
	<view class="tab-bar" :style="{ paddingBottom: safeAreaInsets.bottom + 'px' }">
		<!-- 统一的背景 -->
		<view class="tab-bg" :style="{ paddingBottom: safeAreaInsets.bottom + 'px' }"></view>
		
		<!-- 导航内容 -->
		<view class="tab-content">
			<view class="tab-item" @click="switchTab(0)">
				<image class="tab-image" :src="currentTab === 0 ? '/static/tabbar/sy-liang.png' : '/static/tabbar/sy.png'" mode="aspectFit"></image>
				<text class="tab-text" :class="{ active: currentTab === 0 }">首页</text>
			</view>
			
			<view class="tab-item" @click="switchTab(1)">
				<image class="tab-image" :src="currentTab === 1 ? '/static/tabbar/sq-liang.png' : '/static/tabbar/sq.png'" mode="aspectFit"></image>
				<text class="tab-text" :class="{ active: currentTab === 1 }">闲置</text>
			</view>
			
			<!-- 中间凹陷区域 -->
			<view class="center-area">
				<view class="center-circle"></view>
				<!-- 发布按钮 -->
				<view class="publish-button" @click="switchTab(2)">
					<image class="publish-icon" src="/static/tabbar/fb.png" mode="aspectFit"></image>
					<text class="publish-text">发布</text>
				</view>
			</view>
			
			<view class="tab-item" @click="switchTab(3)">
				<image class="tab-image" :src="currentTab === 3 ? '/static/tabbar/xiaoxi-liang.png' : '/static/tabbar/xiaoxi.png'" mode="aspectFit"></image>
				<text class="tab-text" :class="{ active: currentTab === 3 }">消息</text>
				<view v-if="unreadCount > 0" class="message-badge">{{unreadCount > 99 ? '99+' : unreadCount}}</view>
			</view>
			
			<view class="tab-item" @click="switchTab(4)">
				<image class="tab-image" :src="currentTab === 4 ? '/static/tabbar/my-liang.png' : '/static/tabbar/my.png'" mode="aspectFit"></image>
				<text class="tab-text" :class="{ active: currentTab === 4 }">我的</text>
			</view>
		</view>
	</view>
</template>

<script>
	import deviceMixin from '@/mixins/device.js'
	import messagePollingService from '@/utils/message-polling.js'
	
	export default {
		name: 'TabBar',
		mixins: [deviceMixin],
		props: {
			currentTab: {
				type: Number,
				default: 0
			}
		},
		data() {
			return {
				unreadCount: 0,
				unreadCountUnsubscribe: null // 用于存储取消订阅函数
			}
		},
		created() {
			// 订阅未读消息数量变化
			this.unreadCountUnsubscribe = messagePollingService.addUnreadCountListener(
				unreadCount => {
					this.unreadCount = unreadCount;
				}
			);
		},
		beforeDestroy() {
			// 组件销毁前取消订阅
			if (this.unreadCountUnsubscribe) {
				this.unreadCountUnsubscribe();
				this.unreadCountUnsubscribe = null;
			}
		},
		methods: {
			switchTab(index) {
				const routes = [
					'/pages/index/index',
					'/pages/community/index',
					'/pages/publish/index',
					'/pages/message/index',
					'/pages/my/index'
				];
				
				if (index >= 0 && index < routes.length) {
					uni.switchTab({
						url: routes[index]
					});
				}
			}
		}
	}
</script>

<style>
	.tab-bar {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		height: 100rpx;
		z-index: 999;
	}
	
	.tab-bg {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: #ffffff;
		box-shadow: 0 -2rpx 6rpx rgba(0, 0, 0, 0.04);
		border-top-left-radius: 30rpx;
		border-top-right-radius: 30rpx;
	}
	
	.tab-content {
		position: relative;
		display: flex;
		align-items: center;
		height: 100%;
		padding: 0 20rpx;
	}
	
	.tab-item {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		position: relative;
	}
	
	.tab-item .tab-image {
		width: 48rpx;
		height: 48rpx;
		margin-bottom: 4rpx;
	}
	
	.tab-item .tab-text {
		font-size: 24rpx;
		color: #666666;
	}
	
	.tab-item .tab-text.active {
		color: #fa5e44;
	}
	
	/* 消息未读数角标 */
	.message-badge {
		position: absolute;
		top: -6rpx;
		right: 20rpx;
		min-width: 32rpx;
		height: 32rpx;
		padding: 0 6rpx;
		background-color: #fa5e44;
		color: #ffffff;
		border-radius: 16rpx;
		font-size: 20rpx;
		line-height: 32rpx;
		text-align: center;
		z-index: 10;
	}
	
	/* 中间凹陷区域 */
	.center-area {
		width: 120rpx;
		height: 120rpx;
		position: relative;
		margin-top: -40rpx;
	}
	
	/* 凹陷的圆形背景 */
	.center-circle {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: #ffffff;
		border-radius: 50%;
		box-shadow: inset 0 4rpx 8rpx rgba(0, 0, 0, 0.1);
	}
	
	/* 发布按钮 */
	.publish-button {
		position: absolute;
		top: 10rpx;
		left: 10rpx;
		right: 10rpx;
		bottom: 10rpx;
		background: linear-gradient(135deg, #fc3e2b 0%, #fa7154 100%);
		border-radius: 50%;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		box-shadow: 0 4rpx 16rpx rgba(250, 113, 84, 0.3);
	}
	
	.publish-button .publish-icon {
		width: 40rpx;
		height: 40rpx;
		margin-bottom: 4rpx;
	}
	
	.publish-button .publish-text {
		font-size: 24rpx;
		color: #ffffff;
	}
</style>