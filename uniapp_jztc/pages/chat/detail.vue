<template>
	<view class="chat-container" @tap="hideMorePanel">
		<!-- 聊天内容区域 -->
		<scroll-view 
			class="chat-content" 
			scroll-y 
			:scroll-into-view="'msg-' + (messageList.length - 1)"
			:scroll-with-animation="false"
			:show-scrollbar="false"
			:enhanced="true"
			:bounces="true"
			ref="chatScroll"
			@tap="hideMorePanel"
		>
			<!-- 日期分割线 -->
			<view class="date-divider">
				<text>{{ currentDate }}</text>
			</view>
			
			<!-- 聊天消息列表 -->
			<view class="message-list">
				<view 
					v-for="(message, index) in messageList" 
					:key="index"
					:class="['message-item', message.isSelf ? 'self' : 'other']"
					:id="'msg-' + index"
				>
					<!-- 左侧头像区域（对方） -->
					<view v-if="!message.isSelf" class="avatar-container left">
						<image 
							class="avatar" 
							:src="message.senderAvatar" 
							mode="aspectFill"
						></image>
					</view>
					
					<!-- 消息气泡 -->
					<view class="message-bubble">
						<!-- 文本消息 -->
						<view v-if="message.type === 'text'" class="text-message">
							<text>{{ message.content }}</text>
						</view>
						
						<!-- 图片消息 -->
						<image 
							v-else-if="message.type === 'image'" 
							class="image-message" 
							:src="message.content" 
							mode="widthFix"
							@tap.stop="previewImage(message.content)"
						></image>
						
						<!-- 商品卡片 -->
						<view v-else-if="message.type === 'product'" class="product-card">
							<image class="product-image" :src="message.content.image" mode="aspectFill"></image>
							<view class="product-info">
								<text class="product-title">{{ message.content.title }}</text>
								<text class="product-price">¥{{ message.content.price }}</text>
							</view>
						</view>
					</view>
					
					<!-- 右侧头像区域（自己） -->
					<view v-if="message.isSelf" class="avatar-container right">
						<image 
							class="avatar" 
							:src="message.senderAvatar" 
							mode="aspectFill"
						></image>
					</view>
				</view>
			</view>
			
			<!-- 底部留白，确保最后一条消息可以完全显示 -->
			<view class="bottom-space"></view>
		</scroll-view>
		
		<!-- 输入区域 -->
		<view class="footer">
			<view class="input-area">
				<view class="input-box" @tap.stop="focusInput">
					<input 
						type="text" 
						v-model="inputMessage" 
						placeholder="想跟TA说点什么..." 
						confirm-type="send"
						:focus="inputFocus"
						@confirm="sendTextMessage"
						@focus="onInputFocus"
						@blur="onInputBlur"
						:adjust-position="true"
						:cursor-spacing="10"
					/>
				</view>
				<view 
					class="action-btn" 
					:class="{'send-btn': inputMessage.trim()}"
					@tap.stop="inputMessage.trim() ? sendTextMessage() : toggleMorePanel()"
				>
					<text v-if="inputMessage.trim()">发送</text>
					<uni-icons v-else type="plusempty" size="24" color="#333333"></uni-icons>
				</view>
			</view>
			
			<!-- 更多功能面板 -->
			<view class="more-panel" v-if="showMorePanel && !keyboardVisible" @tap.stop>
				<view class="more-grid">
					<view class="more-item" @tap.stop="chooseImage">
						<view class="more-icon">
							<uni-icons type="image" size="24" color="#ffffff"></uni-icons>
						</view>
						<text class="more-text">图片</text>
					</view>
					
					<view class="more-item" @tap.stop="chooseLocation">
						<view class="more-icon location-icon">
							<uni-icons type="location" size="24" color="#ffffff"></uni-icons>
						</view>
						<text class="more-text">位置</text>
					</view>
					
					<view class="more-item" @tap.stop="shareProduct">
						<view class="more-icon product-icon">
							<uni-icons type="shop" size="24" color="#ffffff"></uni-icons>
						</view>
						<text class="more-text">商品</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import { sendMessage, getMessageList, markMessageRead } from '@/apis/message.js'
	import messagePollingService from '@/utils/message-polling.js'
	
	export default {
		data() {
			return {
				chatId: null,
				chatInfo: {
					id: 0,
					name: '',
					avatar: ''
				},
				messageList: [],
				inputMessage: '',
				showMorePanel: false,
				inputFocus: false,
				currentDate: this.formatDate(new Date()),
				keyboardVisible: false,
				userAvatar: '',
				pagination: {
					page: 1,
					size: 20,
					totalCount: 0,
					totalPage: 0
				},
				isLoading: false,
				nickName: '', // 聊天对象昵称
				_markingAsRead: false,
				chatMessageUnsubscribe: null // 用于存储取消订阅函数
			}
		},
		onLoad(options) {
			// 获取聊天ID或会话ID
			if (options.id || options.sessionId) {
				this.chatId = options.id || options.sessionId
				
				// 设置聊天对象名称（如果在URL参数中传递）
				if (options.nickName) {
					this.nickName = decodeURIComponent(options.nickName)
				}
				
				this.loadChatInfo()
				this.loadChatHistory()
				
				// 设置用户头像
				const userInfo = uni.getStorageSync('USER_INFO') || {};
				this.userAvatar = userInfo.avatarUrl || '/static/demo/0.png';
				
				// 设置当前页面路径，通知消息轮询服务
				messagePollingService.setCurrentPage('/pages/chat/detail');
			}
		},
		onReady() {
			// 设置导航栏标题
			if (this.chatInfo.name) {
				uni.setNavigationBarTitle({
					title: this.chatInfo.name
				})
			}
		},
		// 监听键盘高度变化
		onKeyboardHeightChange(e) {
			this.keyboardVisible = e.height > 0
			
			// 如果键盘弹出，隐藏更多面板
			if (this.keyboardVisible) {
				this.showMorePanel = false
			}
		},
		// 页面显示时标记消息为已读
		onShow() {
			// 确保chatId存在后再开始操作
			if (this.chatId) {
				console.log('聊天详情页面显示，chatId:', this.chatId);
				
				// 再次设置当前页面路径，确保轮询服务能识别白名单页面
				messagePollingService.setCurrentPage('/pages/chat/detail');
				
				// 只在消息列表已加载的情况下标记为已读
				if (this.messageList.length > 0) {
					this.markMessagesAsRead();
				}
				
				// 页面显示时开始聊天轮询
				this.startChatPolling();
			}
		},
		// 显示导航栏右侧按钮
		onNavigationBarButtonTap(e) {
			if (e.index === 0) {
				this.showMoreOptions()
			}
		},
		onHide() {
			// 不再停止轮询，让白名单机制处理
			console.log('聊天详情页面隐藏');
		},
		onUnload() {
			// 页面卸载时通知轮询服务
			console.log('聊天详情页面卸载');
			this.stopChatPolling();
			messagePollingService.setCurrentPage('');
		},
		methods: {
			// 加载聊天信息
			loadChatInfo() {
				if (this.nickName) {
					// 使用传递的昵称更新聊天信息
					this.chatInfo = {
						id: this.chatId,
						name: this.nickName,
						avatar: '/static/demo/0.png' // 使用默认头像，后续会从消息中获取并更新
					};
					
					// 设置导航栏标题
					uni.setNavigationBarTitle({
						title: this.chatInfo.name
					});
				} else {
					// 可以通过API获取用户信息，这里使用简单处理
					this.chatInfo = {
						id: this.chatId,
						name: '用户' + this.chatId,
						avatar: '/static/demo/0.png' // 使用默认头像，后续会从消息中获取并更新
					};
					
					// 设置导航栏标题
					uni.setNavigationBarTitle({
						title: this.chatInfo.name
					});
				}
			},
			
			// 加载聊天历史
			async loadChatHistory() {
				if (!this.chatId || this.isLoading) return;
				
				this.isLoading = true;
				
				try {
					uni.showLoading({
						title: '加载中...',
						mask: true
					});
					
					const res = await getMessageList({
						targetId: this.chatId,
						page: this.pagination.page,
						size: this.pagination.size
					});
					
					if (res && res.code === 0) {
						const { list, totalCount, totalPage, currentPage, size } = res.data;
						
						// 更新分页信息
						this.pagination = {
							page: currentPage,
							size: size,
							totalCount: totalCount,
							totalPage: totalPage
						};
						
						// 格式化消息列表
						if (list && list.length > 0) {
							// 先格式化消息
							const formattedMessages = list.map(item => {
								// 将日期字符串转换为兼容 iOS 的格式
								const formattedDate = item.createdAt ? item.createdAt.replace(/-/g, '/') : '';
								const createdAtTime = formattedDate ? new Date(formattedDate).getTime() : 0;
								
								return {
									type: 'text', // 默认都是文本消息
									content: item.content,
									time: this.formatMessageTime(item.createdAt),
									isSelf: item.isSelf,
									id: item.id,
									senderId: item.senderId,
									senderName: item.senderName,
									senderAvatar: item.senderAvatar || '/static/demo/0.png',
									receiverId: item.receiverId,
									receiverName: item.receiverName,
									receiverAvatar: item.receiverAvatar || '/static/demo/0.png',
									isRead: item.isRead,
									createdAt: createdAtTime // 用于排序
								};
							});
							
							// 按照消息ID或创建时间升序排序，确保旧消息在上，新消息在下
							formattedMessages.sort((a, b) => a.id - b.id);
							
							this.messageList = formattedMessages;
							
							// 标记消息为已读（只在这里调用一次）
							this.markMessagesAsRead();
						} else {
							// 如果没有消息，则显示空数组
							this.messageList = [];
						}
					} else {
						throw new Error(res?.message || '获取消息列表失败');
					}
				} catch (error) {
					console.error('加载聊天记录失败:', error);
					uni.showToast({
						title: '加载聊天记录失败',
						icon: 'none'
					});
				} finally {
					this.isLoading = false;
					uni.hideLoading();
				}
			},
			
			// 格式化消息时间
			formatMessageTime(timestamp) {
				if (!timestamp) return '';
				
				// 对日期字符串进行兼容性处理
				let dateObj;
				if (typeof timestamp === 'string') {
					// 将 "yyyy-MM-dd HH:mm:ss" 转换为 "yyyy/MM/dd HH:mm:ss" 格式，兼容iOS
					const formattedTimestamp = timestamp.replace(/-/g, '/');
					dateObj = new Date(formattedTimestamp);
				} else {
					dateObj = new Date(timestamp);
				}
				
				// 检查是否为有效日期
				if (isNaN(dateObj.getTime())) {
					console.error('Invalid date format:', timestamp);
					return '';
				}
				
				return this.formatTime(dateObj);
			},
			
			// 显示更多选项
			showMoreOptions() {
				uni.showActionSheet({
					itemList: ['清空聊天记录', '投诉', '加入黑名单'],
					success: (res) => {
						switch (res.tapIndex) {
							case 0:
								this.clearChatHistory()
								break
							case 1:
								this.reportUser()
								break
							case 2:
								this.blockUser()
								break
						}
					}
				})
			},
			
			// 清空聊天记录
			clearChatHistory() {
				uni.showModal({
					title: '提示',
					content: '确定要清空聊天记录吗？',
					success: (res) => {
						if (res.confirm) {
							this.messageList = []
							uni.showToast({
								title: '已清空聊天记录',
								icon: 'success'
							})
						}
					}
				})
			},
			
			// 投诉用户
			reportUser() {
				// TODO: 实现投诉用户功能
				uni.showToast({
					title: '已提交投诉',
					icon: 'success'
				})
			},
			
			// 拉黑用户
			blockUser() {
				uni.showModal({
					title: '提示',
					content: '确定要将该用户加入黑名单吗？',
					success: (res) => {
						if (res.confirm) {
							// TODO: 实现拉黑用户功能
							uni.showToast({
								title: '已加入黑名单',
								icon: 'success'
							})
							setTimeout(() => {
								uni.navigateBack()
							}, 1500)
						}
					}
				})
			},
			
			// 聊天消息处理函数
			handleChatMessages(data) {
				if (!data || !data.list) return;
				
				const { list } = data;
				
				// 如果有新消息
				if (list && list.length > 0) {
					// 格式化消息
					const newMessages = list
						.filter(item => {
							// 过滤掉已有的消息（通过ID判断）
							const exists = this.messageList.some(existing => existing.id === item.id);
							return !exists;
						})
						.map(item => {
							// 将日期字符串转换为兼容 iOS 的格式
							const formattedDate = item.createdAt ? item.createdAt.replace(/-/g, '/') : '';
							const createdAtTime = formattedDate ? new Date(formattedDate).getTime() : 0;
							
							return {
								type: 'text', // 默认都是文本消息
								content: item.content,
								time: this.formatMessageTime(item.createdAt),
								isSelf: item.isSelf,
								id: item.id,
								senderId: item.senderId,
								senderName: item.senderName,
								senderAvatar: item.senderAvatar || '/static/demo/0.png',
								receiverId: item.receiverId,
								receiverName: item.receiverName,
								receiverAvatar: item.receiverAvatar || '/static/demo/0.png',
								isRead: item.isRead,
								createdAt: createdAtTime // 用于排序
							};
						});
					
					// 如果有新消息，添加到列表
					if (newMessages.length > 0) {
						// 按照消息ID或创建时间升序排序
						newMessages.sort((a, b) => a.id - b.id);
						
						// 添加新消息到列表
						this.messageList = [...this.messageList, ...newMessages];
						
						// 标记消息为已读
						this.markMessagesAsRead();
					}
				}
			},
			
			// 开始聊天轮询
			startChatPolling() {
				if (!this.chatId) return;
				
				// 获取最后一条消息ID
				let lastMsgId = 0;
				if (this.messageList.length > 0) {
					const lastMsg = this.messageList[this.messageList.length - 1];
					lastMsgId = lastMsg.id || 0;
				}
				
				// 取消之前的订阅（如果有）
				this.stopChatPolling();
				
				// 订阅聊天消息
				this.chatMessageUnsubscribe = messagePollingService.setChatParams(
					{
						targetId: this.chatId,
						lastId: lastMsgId
					},
					this.handleChatMessages
				);
			},
			
			// 停止聊天轮询
			stopChatPolling() {
				if (this.chatMessageUnsubscribe) {
					this.chatMessageUnsubscribe();
					this.chatMessageUnsubscribe = null;
				}
			},
			
			// 发送文本消息成功后更新lastId
			updateLastId(messageId) {
				if (messageId && this.chatId) {
					messagePollingService.updateChatParams({ lastId: messageId });
				}
			},
			
			// 发送文本消息
			async sendTextMessage() {
				if (!this.inputMessage.trim()) return
				
				// 保存输入内容
				const messageContent = this.inputMessage.trim();
				
				// 先清空输入框
				this.inputMessage = ''
				
				// 尝试发送消息
				try {
					// 获取当前时间
					const now = new Date();
					
					// 先添加到本地消息列表（乐观更新）
					const tempMessage = {
						type: 'text',
						content: messageContent,
						time: this.formatTime(now),
						isSelf: true,
						sending: true, // 标记为发送中
						createdAt: now.getTime(), // 用于排序
						senderAvatar: this.userAvatar || '/static/demo/0.png'
					};
					
					this.messageList.push(tempMessage);
					
					// 隐藏更多面板
					this.showMorePanel = false;
					
					// 调用发送消息API
					const response = await sendMessage({
						content: messageContent,
						receiverId: parseInt(this.chatId)
					});
					
					if (response && response.code === 0) {
						// 发送成功，更新消息状态
						const lastIndex = this.messageList.length - 1;
						if (lastIndex >= 0) {
							this.messageList[lastIndex].sending = false;
							if (response.data && response.data.id) {
								this.messageList[lastIndex].id = response.data.id;
								// 更新lastId
								this.updateLastId(response.data.id);
							}
						}
					} else {
						throw new Error(response?.message || '发送消息失败');
					}
				} catch (error) {
					console.error('发送消息失败:', error);
					
					// 移除最后一条消息（如果存在）或标记为发送失败
					const lastIndex = this.messageList.length - 1;
					if (lastIndex >= 0 && this.messageList[lastIndex].sending) {
						this.messageList[lastIndex].sending = false;
						this.messageList[lastIndex].failed = true;
					}
					
					uni.showToast({
						title: '发送失败，请重试',
						icon: 'none'
					});
				}
			},
			
			// 切换更多功能面板
			toggleMorePanel() {
				this.showMorePanel = !this.showMorePanel
				this.inputFocus = false
			},
			
			// 隐藏更多功能面板
			hideMorePanel() {
				if (this.showMorePanel) {
					this.showMorePanel = false
				}
			},
			
			// 聚焦输入框
			focusInput() {
				this.inputFocus = true
				this.hideMorePanel()
			},
			
			// 输入框获得焦点
			onInputFocus(e) {
				this.keyboardVisible = true
			},
			
			// 输入框失去焦点
			onInputBlur() {
				setTimeout(() => {
					this.keyboardVisible = false
				}, 100)
			},
			
			// 选择图片
			chooseImage() {
				uni.chooseImage({
					count: 1,
					success: (res) => {
						const tempFilePath = res.tempFilePaths[0]
						
						// 获取当前时间
						const now = new Date();
						
						// 添加图片消息
						this.messageList.push({
							type: 'image',
							content: tempFilePath,
							time: this.formatTime(now),
							isSelf: true,
							sending: true,
							createdAt: now.getTime() // 用于排序
						})
						
						// 隐藏更多面板
						this.showMorePanel = false
						
						// TODO: 实际上传图片并发送消息的逻辑
						uni.showToast({
							title: '图片功能开发中',
							icon: 'none'
						})
					}
				})
			},
			
			// 选择位置
			chooseLocation() {
				// TODO: 实现选择位置功能
				uni.showToast({
					title: '选择位置功能开发中',
					icon: 'none'
				})
				this.showMorePanel = false
			},
			
			// 分享商品
			shareProduct() {
				// TODO: 实现分享商品功能
				uni.showToast({
					title: '分享商品功能开发中',
					icon: 'none'
				})
				this.showMorePanel = false
			},
			
			// 预览图片
			previewImage(url) {
				// 收集所有图片URL
				const imageUrls = this.messageList
					.filter(msg => msg.type === 'image')
					.map(msg => msg.content)
				
				uni.previewImage({
					current: url,
					urls: imageUrls
				})
			},
			
			// 格式化日期
			formatDate(date) {
				const year = date.getFullYear()
				const month = date.getMonth() + 1
				const day = date.getDate()
				return `${year}年${month}月${day}日`
			},
			
			// 格式化时间
			formatTime(date) {
				const hours = date.getHours().toString().padStart(2, '0')
				const minutes = date.getMinutes().toString().padStart(2, '0')
				return `${hours}:${minutes}`
			},
			
			// 标记消息为已读
			async markMessagesAsRead() {
				if (!this.chatId) {
					console.error('无法标记消息为已读：chatId不存在');
					return;
				}
				
				// 避免重复调用标记为已读
				if (this._markingAsRead) return;
				this._markingAsRead = true;
				
				try {
					console.log('标记消息为已读，targetId:', this.chatId);
					
					// 确保targetId是数字类型
					let targetId = this.chatId;
					if (typeof targetId === 'string') {
						targetId = parseInt(targetId);
					}
					
					// 调用API标记为已读
					const response = await markMessageRead({
						targetId: targetId
					});
					
					if (response && response.code === 0) {
						console.log('消息已成功标记为已读');
						
						// 更新本地消息列表中的已读状态
						this.messageList.forEach(msg => {
							if (!msg.isSelf) {
								msg.isRead = true;
							}
						});
					} else {
						console.error('标记消息为已读失败:', response?.message || '未知错误');
					}
				} catch (error) {
					console.error('标记消息为已读请求失败:', error);
				} finally {
					this._markingAsRead = false;
				}
			}
		}
	}
</script>

<style>
	.chat-container {
		display: flex;
		flex-direction: column;
		height: 100vh;
		background-color: #f5f5f5;
		position: relative;
	}
	
	.chat-content {
		flex: 1;
		height: 0;
		padding: 20rpx 30rpx;
		padding-bottom: 120rpx;
		box-sizing: border-box;
	}
	
	.date-divider {
		display: flex;
		justify-content: center;
		margin: 20rpx 0;
	}
	
	.date-divider text {
		font-size: 24rpx;
		color: #999999;
		background-color: rgba(0, 0, 0, 0.05);
		padding: 6rpx 20rpx;
		border-radius: 24rpx;
	}
	
	.message-list {
		display: flex;
		flex-direction: column;
	}
	
	.message-item {
		display: flex;
		margin-bottom: 30rpx;
		align-items: flex-start;
	}
	
	.message-item.self {
		justify-content: flex-end;
	}
	
	.avatar-container {
		width: 80rpx;
		flex-shrink: 0;
	}
	
	.avatar-container.left {
		margin-right: 20rpx;
	}
	
	.avatar-container.right {
		margin-left: 20rpx;
	}
	
	.avatar {
		width: 80rpx;
		height: 80rpx;
		border-radius: 50%;
		flex-shrink: 0;
	}
	
	.message-bubble {
		max-width: 70%;
	}
	
	.text-message {
		padding: 20rpx;
		border-radius: 12rpx;
		font-size: 28rpx;
		line-height: 1.4;
		word-break: break-all;
	}
	
	.message-item.other .text-message {
		background-color: #ffffff;
		color: #333333;
		border-top-left-radius: 0;
	}
	
	.message-item.self .text-message {
		background-color: #fc3e2b;
		color: #ffffff;
		border-top-right-radius: 0;
	}
	
	.image-message {
		max-width: 100%;
		border-radius: 12rpx;
	}
	
	.product-card {
		width: 400rpx;
		background-color: #ffffff;
		border-radius: 12rpx;
		overflow: hidden;
	}
	
	.product-image {
		width: 100%;
		height: 200rpx;
	}
	
	.product-info {
		padding: 16rpx;
	}
	
	.product-title {
		font-size: 28rpx;
		color: #333333;
		line-height: 1.4;
		margin-bottom: 8rpx;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
	
	.product-price {
		font-size: 32rpx;
		color: #fc3e2b;
		font-weight: 500;
	}
	
	.bottom-space {
		height: 120rpx;
		width: 100%;
	}
	
	.footer {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 100;
	}
	
	.input-area {
		display: flex;
		align-items: center;
		background-color: #f8f8f8;
		border-top: 1rpx solid #e5e5e5;
		padding: 20rpx 30rpx;
		padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	}
	
	.input-box {
		flex: 1;
		background-color: #ffffff;
		border-radius: 36rpx;
		padding: 0 20rpx;
		height: 72rpx;
		border: 1rpx solid #e5e5e5;
	}
	
	.input-box input {
		width: 100%;
		height: 100%;
		font-size: 28rpx;
	}
	
	.action-btn {
		height: 72rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		background-color: #f0f0f0;
		border-radius: 36rpx;
		padding: 0 24rpx;
		margin-left: 16rpx;
		min-width: 72rpx;
	}
	
	.action-btn text {
		font-size: 28rpx;
		color: #666666;
	}
	
	.action-btn.send-btn {
		background-color: #ff8c00;
		background-image: linear-gradient(to right, #ff8c00, #ff4500);
	}
	
	.action-btn.send-btn text {
		color: #ffffff;
	}
	
	.more-panel {
		background-color: #ffffff;
		padding: 20rpx 0;
	}
	
	.more-grid {
		display: flex;
		padding: 20rpx;
	}
	
	.more-item {
		width: 120rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-right: 40rpx;
	}
	
	.more-icon {
		width: 100rpx;
		height: 100rpx;
		background-color: #fc3e2b;
		border-radius: 20rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 16rpx;
	}
	
	.location-icon {
		background-color: #4cd964;
	}
	
	.product-icon {
		background-color: #007aff;
	}
	
	.more-text {
		font-size: 24rpx;
		color: #666666;
	}
</style> 