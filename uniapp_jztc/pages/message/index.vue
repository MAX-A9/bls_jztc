<template>
	<view class="message-container">
		<!-- 顶部标题栏 -->
		<view class="header">
			<view class="title-area">
				<text class="title">消息</text>
				<view class="clear-btn" @click="handleClearRead">
					<uni-icons type="trash" size="16" color="#666666"></uni-icons>
					<text class="clear-text">清除已读</text>
				</view>
			</view>
		</view>
		
		<!-- 占位元素 -->
		<view class="header-placeholder"></view>
		
		<!-- 消息列表 -->
		<scroll-view 
			class="content-scroll" 
			scroll-y
			:show-scrollbar="false"
			:enhanced="true"
			:bounces="true"
			:refresher-enabled="true"
			:refresher-triggered="isRefreshing"
			refresher-background="#ffffff"
			@refresherrefresh="onRefresh"
			@refresherrestore="onRestore"
		>
			<!-- 系统消息 - 已注释 -->
			<!-- 
			<view class="message-item system" @tap="handleSystemClick">
				<view class="avatar system-avatar">
					<uni-icons type="notification-filled" size="24" color="#fc3e2b"></uni-icons>
				</view>
				<view class="message-content">
					<view class="message-top">
						<text class="name">系统通知</text>
						<text class="time">{{systemMessage.time || ''}}</text>
					</view>
					<text class="preview">{{systemMessage.lastMessage || '暂无系统通知'}}</text>
				</view>
			</view>
			-->
			
			<!-- 聊天消息 -->
			<view 
				v-for="(item, index) in chatList"
				:key="item.id || index"
				class="swipe-action-box"
			>
				<uni-swipe-action>
					<uni-swipe-action-item 
						:right-options="swipeOptions" 
						@click="(e) => handleSwipeClick(e, item)"
					>
						<view 
							class="message-item chat"
							@tap="handleChatClick(item)"
						>
							<image class="avatar" :src="item.targetAvatar || '/static/demo/0.png'" mode="aspectFill"></image>
							<view class="message-content">
								<view class="message-top">
									<text class="name">{{ item.targetName || '用户' }}</text>
									<text class="time">{{ item.time || '' }}</text>
								</view>
								<view class="message-bottom">
									<text class="preview">{{ item.lastMessage || '' }}</text>
									<view class="unread" v-if="item.unreadCount && item.unreadCount > 0">{{ item.unreadCount }}</view>
								</view>
							</view>
						</view>
					</uni-swipe-action-item>
				</uni-swipe-action>
			</view>
			
			<!-- 空状态展示 -->
			<view class="empty-state" v-if="chatList.length === 0">
				<image class="empty-icon" src="/static/images/empty.png" mode="aspectFit"></image>
				<text class="empty-text">暂无聊天消息</text>
			</view>
			
			<!-- 底部安全区域 -->
			<view class="safe-area-bottom"></view>
		</scroll-view>
		
		<tab-bar :current-tab="tabIndex"></tab-bar>
	</view>
</template>

<script>
	import TabBar from '@/components/tab-bar/index.vue'
	import deviceAdapter from '@/mixins/device-adapter.js'
	import shareMixin from '@/mixins/share.js'
	import { message } from '@/apis/index.js'
	import messagePollingService from '@/utils/message-polling.js'
	
	export default {
		components: {
			TabBar
		},
		mixins: [deviceAdapter, shareMixin],
		data() {
			return {
				tabIndex: 2,
				isRefreshing: false,
				totalUnreadCount: 0,
				// 系统消息 - 保留但不显示
				systemMessage: {
					lastMessage: '暂无系统通知',
					time: ''
				},
				chatList: [],
				pagination: {
					page: 1,
					size: 20,
					totalCount: 0,
					totalPage: 0
				},
				swipeOptions: [
					{
						text: '已读',
						style: {
							backgroundColor: '#8F8F94'
						}
					},
					{
						text: '删除',
						style: {
							backgroundColor: '#dd524d'
						}
					}
				],
				conversationUnsubscribe: null, // 用于存储取消订阅函数
				_readSessionIds: {} // 用于记录已标记为已读的会话ID
			}
		},
		onLoad() {
			// 初始加载会话列表
			this.fetchSessionList()
			
			// 订阅会话列表变化
			this.startConversationPolling()
		},
		onShow() {
			this.tabIndex = 3
			// 页面显示时开始轮询
			this.startConversationPolling()
		},
		onHide() {
			// 页面隐藏时停止轮询
			this.stopConversationPolling()
		},
		onUnload() {
			// 页面卸载时停止轮询
			this.stopConversationPolling()
		},
		methods: {
			// 开始会话列表轮询
			startConversationPolling() {
				// 如果已经有订阅，先取消
				this.stopConversationPolling()
				
				// 设置会话列表轮询为2秒一次，更快地获取最新数据
				messagePollingService.setPollingInterval(2000)
				
				// 创建新的订阅
				this.conversationUnsubscribe = messagePollingService.setConversationParams(
					{
						page: this.pagination.page,
						size: this.pagination.size
					},
					this.handleConversationUpdate
				)
			},
			
			// 停止会话列表轮询
			stopConversationPolling() {
				if (this.conversationUnsubscribe) {
					this.conversationUnsubscribe()
					this.conversationUnsubscribe = null
					
					// 恢复默认间隔
					messagePollingService.setPollingInterval(5000)
				}
			},
			
			// 处理会话列表更新
			handleConversationUpdate(data) {
				if (!data) return
				
				const { list, totalCount, totalPage, currentPage, size } = data
				
				// 更新分页信息
				this.pagination = {
					page: currentPage,
					size: size,
					totalCount: totalCount,
					totalPage: totalPage
				}
				
				// 直接更新会话列表，确保每次都显示最新数据
				if (list && list.length > 0) {
					this.chatList = list.map(item => {
						return {
							...item,
							time: this.formatTime(item.lastTime)
						}
					})
				} else {
					this.chatList = []
				}
				
				// 已读消息处理：将本地标记为已读的消息与新获取的数据同步
				this.chatList.forEach(item => {
					const readItem = this._readSessionIds && this._readSessionIds[item.id]
					if (readItem) {
						item.unreadCount = 0
					}
				})
			},
			
			// 处理滑动操作点击
			async handleSwipeClick(e, item) {
				if (e.index === 0) {
					// 标记为已读
					await this.markAsRead(item)
				} else if (e.index === 1) {
					// 删除会话
					await this.deleteSession(item)
				}
			},
			
			// 标记会话为已读
			async markAsRead(item) {
				if (!item.unreadCount || item.unreadCount <= 0) {
					uni.showToast({
						title: '没有未读消息',
						icon: 'none'
					})
					return
				}
				
				try {
					uni.showLoading({
						title: '处理中...',
						mask: true
					})
					
					// 调用标记为已读API
					const response = await message.markMessageRead({
						targetId: item.targetId
					})
					
					if (response && response.code === 0) {
						// 更新本地未读消息数
						const index = this.chatList.findIndex(chat => chat.id === item.id)
						if (index !== -1) {
							// 保存原来的未读消息数，用于更新总计数
							const previousUnread = this.chatList[index].unreadCount || 0
							
							// 更新会话的未读消息数为0
							this.chatList[index].unreadCount = 0
							
							// 记录这个会话ID为已读
							this._readSessionIds[item.id] = true
							
							// 更新总的未读消息计数
							if (previousUnread > 0) {
								this.totalUnreadCount = Math.max(0, this.totalUnreadCount - previousUnread)
							}
						}
						
						uni.showToast({
							title: '已标记为已读',
							icon: 'success'
						})
					} else {
						throw new Error(response?.message || '标记为已读失败')
					}
				} catch (error) {
					console.error('标记为已读失败:', error)
					uni.showToast({
						title: '标记为已读失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			},
			
			// 删除会话
			async deleteSession(item) {
				uni.showModal({
					title: '提示',
					content: '确定要删除此会话吗？',
					success: async (res) => {
						if (res.confirm) {
							try {
								uni.showLoading({
									title: '处理中...',
									mask: true
								})
								
								// 调用删除会话API
								const response = await message.deleteConversation(item.id)
								
								uni.hideLoading()
								
								if (response && response.code === 0 && response.data && response.data.success) {
									// 更新本地会话列表
									const index = this.chatList.findIndex(chat => chat.id === item.id)
									if (index !== -1) {
										// 如果有未读消息，减少总未读数
										const previousUnread = this.chatList[index].unreadCount || 0
										if (previousUnread > 0) {
											this.totalUnreadCount = Math.max(0, this.totalUnreadCount - previousUnread)
										}
										// 从本地列表中移除
										this.chatList.splice(index, 1)
									}
									
									uni.showToast({
										title: '已删除会话',
										icon: 'success'
									})
								} else {
									throw new Error(response?.message || '删除会话失败')
								}
							} catch (error) {
								uni.hideLoading()
								console.error('删除会话失败:', error)
								uni.showToast({
									title: '删除会话失败',
									icon: 'none'
								})
							}
						}
					}
				})
			},
			
			// 获取聊天会话列表
			async fetchSessionList() {
				try {
					uni.showLoading({
						title: '加载中...',
						mask: true
					})
					
					const res = await message.getConversationList({
						page: this.pagination.page,
						size: this.pagination.size
					})
					
					if (res && res.code === 0) {
						const { list, totalCount, totalPage, currentPage, size } = res.data
						
						// 更新分页信息
						this.pagination = {
							page: currentPage,
							size: size,
							totalCount: totalCount,
							totalPage: totalPage
						}
						
						// 更新轮询服务的参数
						messagePollingService.updateConversationParams({
							page: currentPage,
							size: size
						})
						
						// 更新会话列表
						if (list && list.length > 0) {
							this.chatList = list.map(item => {
								// 不在这里直接创建Date对象，而是传递原始字符串到formatTime中处理
								return {
									...item,
									time: this.formatTime(item.lastTime)
								}
							})
						} else {
							this.chatList = []
						}
					} else {
						throw new Error(res?.message || '获取会话列表失败')
					}
					
					uni.hideLoading()
				} catch (error) {
					console.error('获取聊天会话列表失败:', error)
					uni.showToast({
						title: '获取会话列表失败',
						icon: 'none'
					})
					uni.hideLoading()
				}
			},
			
			// 清除所有已读消息
			handleClearRead() {
				uni.showModal({
					title: '提示',
					content: '确定要清除所有已读消息吗？',
					success: async (res) => {
						if (res.confirm) {
							uni.showLoading({
								title: '处理中...',
								mask: true
							})
							
							try {
								// 获取目前的会话列表
								const hasReadItems = this.chatList.filter(item => !item.unreadCount || item.unreadCount === 0);
								
								// 如果没有已读消息，提示用户
								if (hasReadItems.length === 0) {
									uni.hideLoading();
									uni.showToast({
										title: '没有可清除的已读消息',
										icon: 'none'
									});
									return;
								}
								
								// 逐个删除已读会话
								let successCount = 0;
								for (const item of hasReadItems) {
									try {
										const response = await message.deleteConversation(item.id);
										if (response && response.code === 0 && response.data && response.data.success) {
											successCount++;
										}
									} catch (err) {
										console.error(`删除会话 ${item.id} 失败:`, err);
									}
								}
								
								// 更新本地会话列表，移除所有已读会话
								this.chatList = this.chatList.filter(item => item.unreadCount && item.unreadCount > 0);
								
								uni.hideLoading();
								uni.showToast({
									title: `已清除${successCount}个已读会话`,
									icon: 'success'
								});
							} catch (error) {
								uni.hideLoading();
								console.error('清除已读消息失败:', error);
								uni.showToast({
									title: '清除已读消息失败',
									icon: 'none'
								});
							}
						}
					}
				});
			},
			
			// 聊天项点击处理
			async handleChatClick(chat) {
				// 跳转到聊天详情页
				uni.navigateTo({
					url: `/pages/chat/detail?id=${chat.targetId}&nickName=${encodeURIComponent(chat.targetName || '用户')}`,
					success: async () => {
						// 清除该聊天的未读消息
						if (chat.unreadCount && chat.unreadCount > 0) {
							try {
								const response = await message.clearSessionUnread(chat.id)
								
								if (response && response.code === 0) {
									// 更新本地未读消息数
									const index = this.chatList.findIndex(item => item.id === chat.id)
									if (index !== -1) {
										// 保存原来的未读消息数，用于更新总计数
										const previousUnread = this.chatList[index].unreadCount || 0
										
										// 更新会话的未读消息数为0
										this.chatList[index].unreadCount = 0
										
										// 记录这个会话ID为已读
										this._readSessionIds[chat.id] = true
										
										// 更新总的未读消息计数
										if (previousUnread > 0) {
											this.totalUnreadCount = Math.max(0, this.totalUnreadCount - previousUnread)
										}
									}
								}
							} catch (error) {
								console.error('清除会话未读消息失败:', error)
							}
						}
					}
				})
			},
			
			// 处理系统消息点击 - 已注释
			/*
			handleSystemClick() {
				uni.navigateTo({
					url: '/pages/notification/index'
				})
			},
			*/
			
			// 下拉刷新处理
			async onRefresh() {
				if (this.isRefreshing) return
				this.isRefreshing = true
				
				try {
					// 重置分页
					this.pagination.page = 1
					
					// 重新获取会话列表
					await this.fetchSessionList()
					
					uni.showToast({
						title: '刷新成功',
						icon: 'success'
					})
				} catch (error) {
					console.error('刷新失败:', error)
					uni.showToast({
						title: '刷新失败',
						icon: 'none'
					})
				} finally {
					this.isRefreshing = false
				}
			},
			
			// 刷新复位
			onRestore() {
				console.log('刷新复位')
			},
			
			// 格式化时间显示
			formatTime(timestamp) {
				if (!timestamp) return ''
				
				// 处理日期字符串兼容性
				let msgDate;
				if (typeof timestamp === 'string') {
					// 将 "yyyy-MM-dd HH:mm:ss" 转换为 "yyyy/MM/dd HH:mm:ss" 格式，兼容iOS
					const formattedTimestamp = timestamp.replace(/-/g, '/');
					msgDate = new Date(formattedTimestamp);
				} else {
					msgDate = new Date(timestamp);
				}
				
				// 检查是否为有效日期
				if (isNaN(msgDate.getTime())) {
					console.error('Invalid date format:', timestamp);
					return '';
				}
				
				const now = new Date()
				const diffDays = Math.floor((now - msgDate) / (24 * 60 * 60 * 1000))
				
				// 获取时间部分
				const hours = msgDate.getHours().toString().padStart(2, '0')
				const minutes = msgDate.getMinutes().toString().padStart(2, '0')
				const timeStr = `${hours}:${minutes}`
				
				// 今天的消息只显示时间
				if (diffDays === 0) {
					return timeStr
				}
				
				// 昨天的消息
				if (diffDays === 1) {
					return '昨天'
				}
				
				// 一周内的消息
				if (diffDays < 7) {
					const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
					return weekdays[msgDate.getDay()]
				}
				
				// 更早的消息显示日期
				const month = (msgDate.getMonth() + 1).toString().padStart(2, '0')
				const day = msgDate.getDate().toString().padStart(2, '0')
				return `${month}-${day}`
			}
		}
	}
</script>

<style>
	.message-container {
		display: flex;
		flex-direction: column;
		height: 100vh;
		background-color: #ffffff;
	}
	
	.header {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 99;
		background-color: #ffffff;
		padding: 0 30rpx;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
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
	
	.header-placeholder {
		height: 88rpx;
		flex-shrink: 0;
	}
	
	.content-scroll {
		flex: 1;
		background-color: #ffffff;
		z-index: 1;
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
	}
	
	.message-item {
		display: flex;
		align-items: center;
		padding: 20rpx 30rpx;
		background-color: #ffffff;
	}
	
	.swipe-action-box {
		border-bottom: 1rpx solid #f5f5f5;
	}
	
	.swipe-action-box:last-child {
		border-bottom: none;
	}
	
	.avatar {
		width: 88rpx;
		height: 88rpx;
		border-radius: 50%;
		margin-right: 20rpx;
	}
	
	.system-avatar {
		background-color: #fff2f0;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 88rpx;
		height: 88rpx;
		border-radius: 50%;
		margin-right: 20rpx;
	}
	
	.message-content {
		flex: 1;
	}
	
	.message-top {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 8rpx;
	}
	
	.name {
		font-size: 30rpx;
		color: #333333;
		font-weight: 500;
	}
	
	.time {
		font-size: 24rpx;
		color: #999999;
	}
	
	.message-bottom {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.preview {
		font-size: 26rpx;
		color: #666666;
		flex: 1;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
	
	.unread {
		min-width: 36rpx;
		height: 36rpx;
		padding: 0 10rpx;
		background-color: #fc3e2b;
		border-radius: 18rpx;
		color: #ffffff;
		font-size: 24rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-left: 12rpx;
		flex-shrink: 0;
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
	
	.safe-area-bottom {
		height: calc(100rpx + env(safe-area-inset-bottom));
	}
	
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 100rpx 0;
	}
	
	.empty-icon {
		width: 200rpx;
		height: 200rpx;
		margin-bottom: 20rpx;
	}
	
	.empty-text {
		font-size: 28rpx;
		color: #999999;
	}
</style>
