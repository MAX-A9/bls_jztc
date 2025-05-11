<template>
	<view class="notification-container">
		<!-- 顶部导航栏 -->
		<view class="header">
			<view class="title-area">
				<text class="title">系统通知</text>
				<view class="clear-btn" @click="handleClearAll">
					<uni-icons type="trash" size="16" color="#666666"></uni-icons>
					<text class="clear-text">清除全部</text>
				</view>
			</view>
		</view>
		
		<!-- 占位元素 -->
		<view class="header-placeholder"></view>
		
		<!-- 通知列表 -->
		<scroll-view 
			class="notification-content" 
			scroll-y 
			:show-scrollbar="false"
			:enhanced="true"
			:bounces="true"
			:refresher-enabled="true"
			:refresher-triggered="isRefreshing"
			refresher-background="#f5f5f5"
			@refresherrefresh="onRefresh"
			@refresherrestore="onRestore"
		>
			<view v-if="notificationList.length === 0" class="empty-container">
				<image class="empty-image" src="/static/images/empty-notification.png" mode="aspectFit"></image>
				<text class="empty-text">暂无系统通知</text>
			</view>
			
			<view 
				v-else
				class="notification-item"
				v-for="(item, index) in notificationList"
				:key="index"
				@tap="handleNotificationClick(item)"
			>
				<view class="notification-icon" :class="[item.type]">
					<uni-icons :type="getIconByType(item.type)" size="24" color="#ffffff"></uni-icons>
				</view>
				<view class="notification-content">
					<view class="notification-top">
						<text class="title">{{ item.title }}</text>
						<text class="time">{{ item.time }}</text>
					</view>
					<text class="content">{{ item.content }}</text>
				</view>
				<view class="unread-dot" v-if="!item.isRead"></view>
			</view>
			
			<!-- 底部安全区域 -->
			<view class="safe-area-bottom"></view>
		</scroll-view>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				isRefreshing: false,
				notificationList: [
					{
						id: 1,
						type: 'system',
						title: '系统更新通知',
						content: '亲爱的用户，我们的应用已更新到最新版本，新增了多项实用功能，欢迎体验！',
						time: '今天 12:30',
						isRead: false
					},
					{
						id: 2,
						type: 'comment',
						title: '评论通知',
						content: '用户"张先生"评论了您的闲置物品："这个价格很合理，我很感兴趣"',
						time: '今天 10:15',
						isRead: true
					},
					{
						id: 3,
						type: 'like',
						title: '点赞通知',
						content: '您的闲置物品"iPhone 12 128G 蓝色"收到了3个新的点赞',
						time: '昨天 15:40',
						isRead: false
					},
					{
						id: 4,
						type: 'order',
						title: '订单通知',
						content: '您的订单#20230615001已发货，预计3天内送达',
						time: '昨天 09:20',
						isRead: true
					},
					{
						id: 5,
						type: 'system',
						title: '账号安全提醒',
						content: '您的账号于6月10日在新设备上登录，如非本人操作，请及时修改密码',
						time: '3天前',
						isRead: true
					},
					{
						id: 6,
						type: 'comment',
						title: '评论通知',
						content: '用户"李女士"回复了您的评论："好的，那就这个价格成交"',
						time: '4天前',
						isRead: true
					},
					{
						id: 7,
						type: 'like',
						title: '点赞通知',
						content: '您的闲置物品"九成新小米平板5 Pro"收到了5个新的点赞',
						time: '5天前',
						isRead: true
					},
					{
						id: 8,
						type: 'order',
						title: '订单通知',
						content: '您的订单#20230608002已完成，感谢您的使用',
						time: '一周前',
						isRead: true
					}
				]
			}
		},
		methods: {
			getIconByType(type) {
				const iconMap = {
					'system': 'notification-filled',
					'comment': 'chat-filled',
					'like': 'heart-filled',
					'order': 'cart-filled'
				}
				return iconMap[type] || 'notification-filled'
			},
			
			handleClearAll() {
				uni.showModal({
					title: '提示',
					content: '确定要清空所有通知吗？',
					success: (res) => {
						if (res.confirm) {
							this.notificationList = []
							uni.showToast({
								title: '已清空所有通知',
								icon: 'success'
							})
						}
					}
				})
			},
			
			handleNotificationClick(notification) {
				// 标记为已读
				const index = this.notificationList.findIndex(item => item.id === notification.id)
				if (index !== -1 && !this.notificationList[index].isRead) {
					this.notificationList[index].isRead = true
				}
				
				// 根据通知类型跳转到不同页面
				switch (notification.type) {
					case 'comment':
						// 跳转到评论详情
						uni.navigateTo({
							url: `/pages/content/detail?id=${notification.id}&type=comment`
						})
						break
					case 'like':
						// 跳转到点赞详情
						uni.navigateTo({
							url: `/pages/my/publish/index?highlight=${notification.id}`
						})
						break
					case 'order':
						// 跳转到订单详情
						uni.navigateTo({
							url: `/pages/my/orders/index?id=${notification.id}`
						})
						break
					case 'system':
					default:
						// 显示通知详情
						uni.showModal({
							title: notification.title,
							content: notification.content,
							showCancel: false,
							confirmText: '知道了'
						})
						break
				}
			},
			
			async onRefresh() {
				if (this.isRefreshing) return
				this.isRefreshing = true
				
				try {
					await new Promise(resolve => setTimeout(resolve, 1000))
					
					// 更新通知列表
					this.notificationList = [
						{
							id: 9,
							type: 'system',
							title: '新活动通知',
							content: '618年中大促即将开始，多款商品低至5折，先到先得！',
							time: '刚刚',
							isRead: false
						},
						...this.notificationList
					]
					
					uni.showToast({
						title: '刷新成功',
						icon: 'success'
					})
				} catch (error) {
					console.error('刷新失败:', error)
					uni.showToast({
						title: '刷新失败',
						icon: 'error'
					})
				} finally {
					this.isRefreshing = false
				}
			},
			
			onRestore() {
				console.log('刷新复位')
			}
		}
	}
</script>

<style>
	.notification-container {
		display: flex;
		flex-direction: column;
		height: 100vh;
		background-color: #f5f5f5;
	}
	
	.header {
		position: fixed;
		top: var(--window-top);
		left: 0;
		right: 0;
		z-index: 99;
		background-color: #ffffff;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
		padding: 0 30rpx;
	}
	
	.title-area {
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 88rpx;
	}
	
	.title {
		font-size: 36rpx;
		font-weight: bold;
		color: #333333;
	}
	
	.clear-btn {
		display: flex;
		align-items: center;
		padding: 8rpx 16rpx;
		background-color: #f5f5f5;
		border-radius: 24rpx;
	}
	
	.clear-text {
		font-size: 24rpx;
		color: #666666;
		margin-left: 4rpx;
	}
	
	.header-placeholder {
		height: calc(var(--window-top) + 88rpx);
		flex-shrink: 0;
	}
	
	.notification-content {
		flex: 1;
		box-sizing: border-box;
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
	}
	
	.notification-item {
		display: flex;
		padding: 24rpx 30rpx;
		background-color: #ffffff;
		margin-bottom: 20rpx;
		position: relative;
	}
	
	.notification-icon {
		width: 80rpx;
		height: 80rpx;
		border-radius: 50%;
		margin-right: 20rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}
	
	.notification-icon.system {
		background-color: #fc3e2b;
	}
	
	.notification-icon.comment {
		background-color: #007aff;
	}
	
	.notification-icon.like {
		background-color: #ff2d55;
	}
	
	.notification-icon.order {
		background-color: #4cd964;
	}
	
	.notification-content {
		flex: 1;
	}
	
	.notification-top {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 8rpx;
	}
	
	.notification-top .title {
		font-size: 30rpx;
		font-weight: 500;
	}
	
	.notification-top .time {
		font-size: 24rpx;
		color: #999999;
	}
	
	.notification-content .content {
		font-size: 26rpx;
		color: #666666;
		line-height: 1.5;
	}
	
	.unread-dot {
		position: absolute;
		top: 24rpx;
		right: 30rpx;
		width: 16rpx;
		height: 16rpx;
		border-radius: 50%;
		background-color: #fc3e2b;
	}
	
	.empty-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding-top: 200rpx;
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
	
	.safe-area-bottom {
		height: 40rpx;
	}
</style> 