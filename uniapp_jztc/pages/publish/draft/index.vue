<template>
	<view class="draft-container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">草稿箱</text>
				<view class="right-btn" @tap="handleClearAll" v-if="draftList.length > 0">
					<text class="clear-text">清空</text>
				</view>
			</view>
		</view>
		
		<!-- 占位元素 -->
		<view class="placeholder" :style="{ height: `calc(${statusBarHeight}px + ${navBarHeight}px)` }"></view>
		
		<!-- 草稿列表 -->
		<view class="draft-list" v-if="draftList.length > 0">
			<view
				class="draft-item"
				v-for="item in draftList"
				:key="item.id"
				@tap="handleDraftClick(item)"
			>
				<view class="draft-info">
					<view class="draft-title">{{ getTitlePreview(item) }}</view>
					<view class="draft-meta">
						<text class="draft-type">{{ item.type === 'idle' ? '闲置' : '信息' }}</text>
						<text class="draft-time">{{ formatTime(item.updateTime) }}</text>
					</view>
				</view>
				
				<view class="draft-actions">
					<view class="delete-btn" @tap.stop="handleDeleteDraft(item)">
						<uni-icons type="trash" size="18" color="#ff4500"></uni-icons>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 空状态 -->
		<view class="empty-state" v-else>
			<image class="empty-icon" src="/static/images/icons/empty.png" mode="aspectFit"></image>
			<text class="empty-text">暂无草稿</text>
		</view>
	</view>
</template>

<script>
	import deviceMixin from '@/mixins/device.js'
	import { getIdleDrafts, getInfoDrafts, deleteIdleDraft, deleteInfoDraft, getAllDraftList } from '@/utils/storage.js'
	import deviceInfo from '@/utils/device-info.js'
	
	export default {
		mixins: [deviceMixin],
		data() {
			return {
				statusBarHeight: 0,
				navBarHeight: 44,
				draftList: []
			}
		},
		onLoad() {
			// 获取状态栏高度
			this.statusBarHeight = deviceInfo.getStatusBarHeight();
			
			// 加载草稿数据
			this.loadDraftList()
		},
		methods: {
			// 返回上一页
			handleBack() {
				uni.navigateBack()
			},
			
			// 加载草稿列表
			loadDraftList() {
				this.draftList = getAllDraftList()
			},
			
			// 格式化时间显示
			formatTime(timeStr) {
				if (!timeStr) return '未知时间'
				
				const date = new Date(timeStr)
				const now = new Date()
				const diff = now - date
				
				// 一小时内
				if (diff < 3600000) {
					const minutes = Math.floor(diff / 60000)
					return `${minutes === 0 ? '刚刚' : minutes + '分钟前'}`
				}
				
				// 一天内
				if (diff < 86400000) {
					const hours = Math.floor(diff / 3600000)
					return `${hours}小时前`
				}
				
				// 一周内
				if (diff < 604800000) {
					const days = Math.floor(diff / 86400000)
					return `${days}天前`
				}
				
				// 超过一周
				return `${date.getMonth() + 1}月${date.getDate()}日`
			},
			
			// 获取标题预览（如果没有标题则显示描述的一部分）
			getTitlePreview(draft) {
				if (draft.title && draft.title.trim()) {
					return draft.title
				}
				
				if (draft.description && draft.description.trim()) {
					const preview = draft.description.slice(0, 20)
					return preview.length >= 20 ? preview + '...' : preview
				}
				
				return draft.type === 'idle' ? '闲置物品草稿' : '信息发布草稿'
			},
			
			// 处理草稿点击，根据类型跳转到不同页面
			handleDraftClick(draft) {
				if (draft.type === 'idle') {
					uni.navigateTo({
						url: `/pages/publish/idle/index?draftId=${draft.id}`
					})
				} else {
					uni.navigateTo({
						url: `/pages/publish/info/index?draftId=${draft.id}`
					})
				}
			},
			
			// 处理删除草稿
			handleDeleteDraft(draft) {
				uni.showModal({
					title: '提示',
					content: '确定要删除该草稿吗？',
					success: (res) => {
						if (res.confirm) {
							let success = false
							
							if (draft.type === 'idle') {
								success = deleteIdleDraft(draft.id)
							} else {
								success = deleteInfoDraft(draft.id)
							}
							
							if (success) {
								uni.showToast({
									title: '删除成功',
									icon: 'success'
								})
								// 重新加载列表
								this.loadDraftList()
							} else {
								uni.showToast({
									title: '删除失败',
									icon: 'none'
								})
							}
						}
					}
				})
			},
			
			// 清空所有草稿
			handleClearAll() {
				uni.showModal({
					title: '提示',
					content: '确定要清空所有草稿吗？',
					success: (res) => {
						if (res.confirm) {
							// 遍历删除所有草稿
							this.draftList.forEach(draft => {
								if (draft.type === 'idle') {
									deleteIdleDraft(draft.id)
								} else {
									deleteInfoDraft(draft.id)
								}
							})
							
							// 清空列表
							this.draftList = []
							
							uni.showToast({
								title: '清空成功',
								icon: 'success'
							})
						}
					}
				})
			}
		}
	}
</script>

<style>
	.draft-container {
		padding-bottom: 30rpx;
		background-color: #f8f8f8;
		min-height: 100vh;
	}
	
	.nav-bar {
		background-color: #ffffff;
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 99;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
	}
	
	.nav-content {
		display: flex;
		align-items: center;
		position: relative;
		justify-content: center;
	}
	
	.back-btn {
		position: absolute;
		left: 20rpx;
		padding: 10rpx;
	}
	
	.right-btn {
		position: absolute;
		right: 30rpx;
		z-index: 10;
	}
	
	.clear-text {
		font-size: 28rpx;
		color: #666666;
	}
	
	.nav-title {
		font-size: 34rpx;
		font-weight: 500;
		color: #333333;
		max-width: 350rpx;
		text-align: center;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.placeholder {
		width: 100%;
	}
	
	.draft-list {
		padding: 20rpx 30rpx;
	}
	
	.draft-item {
		margin-bottom: 20rpx;
		background-color: #ffffff;
		border-radius: 12rpx;
		padding: 30rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
	}
	
	.draft-info {
		flex: 1;
		overflow: hidden;
	}
	
	.draft-title {
		font-size: 32rpx;
		color: #333333;
		margin-bottom: 16rpx;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.draft-meta {
		display: flex;
		align-items: center;
	}
	
	.draft-type {
		font-size: 24rpx;
		color: #ffffff;
		background-color: #ff8c00;
		padding: 6rpx 12rpx;
		border-radius: 6rpx;
		margin-right: 16rpx;
	}
	
	.draft-time {
		font-size: 24rpx;
		color: #999999;
	}
	
	.draft-actions {
		margin-left: 20rpx;
	}
	
	.delete-btn {
		padding: 10rpx;
	}
	
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding-top: 200rpx;
	}
	
	.empty-icon {
		width: 200rpx;
		height: 200rpx;
		margin-bottom: 30rpx;
	}
	
	.empty-text {
		font-size: 30rpx;
		color: #999999;
	}
</style> 