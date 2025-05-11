<template>
	<view class="fixed-bottom">
		<view class="bottom-left">
			<view class="action-btn" @tap="handleComment">
				<uni-icons type="chat" size="20" color="#666666"></uni-icons>
				<text class="action-text">留言</text>
			</view>
			<view class="action-btn" @tap="handleCollect">
				<uni-icons :type="isCollected ? 'star-filled' : 'star'" size="20" :color="isCollected ? '#fc3e2b' : '#666666'"></uni-icons>
				<text class="action-text" :class="{ active: isCollected }">收藏</text>
			</view>
			<button class="action-btn share-button" open-type="share">
				<uni-icons type="redo" size="20" color="#666666"></uni-icons>
				<text class="action-text">分享</text>
			</button>
		</view>
		<view class="bottom-right">
			<button class="action-button" @tap="handleMessage">
				<uni-icons type="chat" size="20" color="#FFFFFF"></uni-icons>
				<text class="button-text">私信</text>
			</button>
		</view>
	</view>
</template>

<script>
	export default {
		name: 'ActionBar',
		props: {
			// 是否已收藏
			isCollected: {
				type: Boolean,
				default: false
			},
			// 发布者信息
			publisher: {
				type: Object,
				default: () => ({
					id: '',
					name: ''
				})
			},
			// 是否显示留言按钮
			showComment: {
				type: Boolean,
				default: true
			},
			// 按钮文本
			messageText: {
				type: String,
				default: '私信'
			}
		},
		data() {
			return {
			}
		},
		methods: {
			// 留言
			handleComment() {
				this.$emit('comment');
			},
			// 收藏
			handleCollect() {
				this.$emit('collect');
			},
			// 私信
			handleMessage() {
				// 如果父组件未提供发布者ID
				if (!this.publisher.id) {
					uni.showToast({
						title: '无法获取发布者信息',
						icon: 'none'
					});
					return;
				}
				
				// 触发事件由父组件处理
				this.$emit('message', this.publisher);
			}
		}
	}
</script>

<style>
	/* 底部操作栏样式 */
	.fixed-bottom {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		height: 100rpx;
		background-color: #FFFFFF;
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0 20rpx 0 12rpx;
		z-index: 99;
		padding-bottom: env(safe-area-inset-bottom);
	}
	
	.bottom-left {
		display: flex;
		gap: 8rpx;
	}
	
	.action-btn {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 4rpx;
		min-width: 60rpx;
		padding: 0 4rpx;
	}
	
	.share-button {
		background: transparent;
		border: none;
		margin: 0;
		padding: 0 4rpx;
		line-height: normal;
	}
	
	.share-button::after {
		border: none;
	}
	
	.action-text {
		font-size: 22rpx;
		color: #666666;
		line-height: 1;
	}
	
	.action-text.active {
		color: #fc3e2b;
	}
	
	.action-button {
		background: linear-gradient(135deg, #fc3e2b 0%, #fa7154 100%);
		padding: 0 40rpx;
		border-radius: 40rpx;
		border: none;
		height: 72rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8rpx;
	}
	
	.button-text {
		font-size: 28rpx;
		color: #FFFFFF;
		font-weight: 500;
	}
</style> 