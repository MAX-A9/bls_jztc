<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">设置</text>
			</view>
		</view>
		
		<!-- 隐藏的canvas用于图片处理 -->
		<canvas canvas-id="avatarCanvas" style="width: 200px; height: 200px; position: absolute; left: -1000px; top: -1000px;"></canvas>
		
		<!-- 设置内容 -->
		<scroll-view 
			class="content-scroll" 
			scroll-y
			:style="{ top: `calc(${statusBarHeight}px + ${navBarHeight}px)` }"
		>
			<view class="settings-wrapper">
				<!-- 个人设置 -->
				<view class="settings-section">
					<view class="section-title">个人设置</view>
					
					<view class="settings-item">
						<text class="item-label">头像</text>
						<view class="item-right">
							<button class="avatar-wrapper" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
								<image class="avatar-preview" :src="userInfo.avatar" mode="aspectFill"></image>
							</button>
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
					
					<view class="settings-item">
						<text class="item-label">昵称</text>
						<view class="item-right nickname-input-wrapper">
							<input 
								class="nickname-input" 
								type="nickname" 
								:placeholder="userInfo.nickname || '请输入昵称'" 
								@blur="onNicknameInput"
							/>
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
					
					<view class="settings-item" @tap="handleEditPhone">
						<text class="item-label">手机号</text>
						<view class="item-right">
							<text class="item-value">{{ userInfo.phone || '未设置' }}</text>
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
				</view>
				
				<!-- 订阅通知 -->
				<view class="settings-section">
					<view class="section-title">订阅通知</view>
					
					<view class="settings-item">
						<text class="item-label">消息通知</text>
						<switch 
							:checked="notifications.message" 
							color="#fc3e2b"
							@change="handleSwitchChange('message', $event)"
						/>
					</view>
					
					<view class="settings-item">
						<text class="item-label">留言通知</text>
						<switch 
							:checked="notifications.comment" 
							color="#fc3e2b"
							@change="handleSwitchChange('comment', $event)"
						/>
					</view>
				</view>
				
				<!-- 其他设置 -->
				<view class="settings-section">
					<view class="section-title">其他设置</view>
					
					<view class="settings-item" @tap="handleClearCache">
						<text class="item-label">清除缓存</text>
						<view class="item-right">
							<text class="cache-size">{{ cacheSize }}</text>
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
					
					<view class="settings-item" @tap="handleAbout">
						<text class="item-label">关于我们</text>
						<view class="item-right">
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
				</view>
				
				<!-- 底部协议 -->
				<view class="footer">
					<view class="agreement-links">
						<text class="agreement-link" @tap="handleAgreement('service')">用户协议</text>
						<text class="divider">|</text>
						<text class="agreement-link" @tap="handleAgreement('privacy')">隐私政策</text>
					</view>
				</view>
			</view>
		</scroll-view>
		
		<!-- 自定义关于我们弹窗 -->
		<view class="about-popup" v-if="showAboutPopup">
			<view class="about-container">
				<view class="about-header">
					<text class="about-title">关于我们</text>
				</view>
				<view class="about-content">
					<text class="app-name">{{ appInfo.name }} v1.0.0</text>
					<text class="app-desc">{{ appInfo.description }}</text>
				</view>
				<view class="about-footer">
					<view class="confirm-btn" @tap="showAboutPopup = false">确定</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import { mapGetters, mapActions } from 'vuex';
	import { user } from '@/apis/index.js';
	import { agreement } from '../../apis/index.js';
	import { settings } from '../../apis/index.js';
	import { getUserInfo } from '@/utils/auth.js';
	import deviceInfo from '@/utils/device-info.js';
	
	export default {
		data() {
			return {
				statusBarHeight: 0,
				navBarHeight: 44,
				userInfo: {
					avatar: '/static/demo/0.png',
					nickname: '橘子用户',
					phone: ''
				},
				notifications: {
					message: true,
					comment: true
				},
				cacheSize: '8.5MB',
				updateData: {}, // 临时存储要更新的数据
				appInfo: {
					name: '橘子同城',
					description: '为您提供便捷的本地生活服务'
				},
				showAboutPopup: false // 控制关于我们弹窗显示
			}
		},
		computed: {
			...mapGetters('user', ['isLoggedIn'])
		},
		onLoad() {
			// 获取状态栏高度 - 使用新的设备信息工具
			this.statusBarHeight = deviceInfo.getStatusBarHeight();
			
			// 加载用户信息
			this.loadUserInfo();
			
			// 获取缓存大小
			this.updateCacheSize();
			
			// 获取应用基础信息
			this.fetchAppBaseInfo();
		},
		methods: {
			...mapActions('user', ['getUserInfo']),
			
			// 加载用户信息
			async loadUserInfo() {
				if (this.isLoggedIn) {
					try {
						// 从Vuex加载用户信息
						const userInfo = await this.getUserInfo();
						
						// 从本地存储获取，确保数据完整
						const localUserInfo = getUserInfo();
						
						this.userInfo = {
							avatar: localUserInfo.avatarUrl || '/static/demo/0.png',
							nickname: localUserInfo.realName || '橘子用户',
							phone: localUserInfo.phone || ''
						};
					} catch (err) {
						console.error('加载用户信息失败:', err);
						uni.showToast({
							title: '加载用户信息失败',
							icon: 'none'
						});
					}
				}
			},
			
			// 处理返回按钮点击
			handleBack() {
				uni.navigateBack()
			},
			
			// 头像选择处理
			async onChooseAvatar(e) {
				try {
					const { avatarUrl } = e.detail;
					
					if (!avatarUrl) {
						return;
					}
					
					// 显示加载提示
					uni.showLoading({ title: '处理中...' });
					
					// 1. 更新本地显示
					this.userInfo.avatar = avatarUrl;
					
					try {
						// 2. 调用接口保存头像（会自动转换为base64）
						const result = await user.uploadAvatar(avatarUrl);
						
						// 3. 更新Vuex中的用户信息
						await this.getUserInfo();
						
						uni.hideLoading();
						uni.showToast({
							title: '头像已更新',
							icon: 'success'
						});
					} catch (err) {
						uni.hideLoading();
						uni.showToast({
							title: '头像更新失败: ' + (err.message || '未知错误'),
							icon: 'none',
							duration: 3000
						});
					}
				} catch (err) {
					uni.hideLoading();
					uni.showToast({
						title: '更新头像失败',
						icon: 'none'
					});
				}
			},
			
			// 昵称输入处理
			async onNicknameInput(e) {
				try {
					const nickname = e.detail.value;
					
					if (!nickname || nickname === this.userInfo.nickname) {
						return;
					}
					
					// 更新本地显示
					this.userInfo.nickname = nickname;
					this.updateData.realName = nickname;
					
					// 保存到后端
					await this.saveUserInfo();
					
					uni.showToast({
						title: '昵称已更新',
						icon: 'success'
					});
				} catch (err) {
					console.error('更新昵称失败:', err);
					uni.showToast({
						title: '更新昵称失败',
						icon: 'none'
					});
				}
			},
			
			// 修改手机号
			handleEditPhone() {
				uni.showModal({
					title: '修改手机号',
					editable: true,
					placeholderText: '请输入11位手机号',
					content: this.userInfo.phone || '',
					success: async (res) => {
						if (res.confirm && res.content) {
							const phone = res.content;
							
							// 简单的手机号验证
							if (!/^1\d{10}$/.test(phone)) {
								uni.showToast({
									title: '请输入有效的手机号',
									icon: 'none'
								});
								return;
							}
							
							if (phone === this.userInfo.phone) {
								return;
							}
							
							try {
								// 更新本地显示
								this.userInfo.phone = phone;
								this.updateData.phone = phone;
								
								// 保存到后端
								uni.showLoading({ title: '更新中...' });
								await this.saveUserInfo();
								uni.hideLoading();
								
								uni.showToast({
									title: '手机号已更新',
									icon: 'success'
								});
							} catch (err) {
								uni.hideLoading();
								console.error('更新手机号失败:', err);
								uni.showToast({
									title: '更新手机号失败',
									icon: 'none'
								});
							}
						}
					}
				});
			},
			
			// 保存用户信息到后端
			async saveUserInfo() {
				if (Object.keys(this.updateData).length === 0) {
					return;
				}
				
				try {
					// 调用接口保存用户信息
					// 接口字段说明：
					// avatarUrl - 用户头像URL
					// realName - 用户昵称
					// phone - 用户手机号
					await user.updateClientProfile(this.updateData);
					
					// 更新Vuex中的用户信息
					this.getUserInfo();
					
					// 清空临时数据
					this.updateData = {};
					
					return true;
				} catch (err) {
					console.error('保存用户信息失败:', err);
					throw err;
				}
			},
			
			// 处理开关变化
			handleSwitchChange(type, event) {
				this.notifications[type] = event.detail.value
			},
			
			// 处理清除缓存
			handleClearCache() {
				uni.showModal({
					title: '提示',
					content: '确定要清除缓存吗？',
					success: (res) => {
						if (res.confirm) {
							try {
								// 显示加载提示
								uni.showLoading({
									title: '清除中...',
									mask: true
								});
								
								// 清除本地存储中的非必要数据
								// 保留登录信息和用户数据
								const keysToPreserve = ['token', 'USER_INFO', 'currentLocation'];
								
								// 获取所有存储的keys
								uni.getStorageInfo({
									success: (res) => {
										const allKeys = res.keys || [];
										
										// 过滤出需要删除的keys
										const keysToRemove = allKeys.filter(key => !keysToPreserve.includes(key));
										
										// 逐个删除存储
										keysToRemove.forEach(key => {
											uni.removeStorageSync(key);
										});
										
										// 清除图片缓存
										uni.getSavedFileList({
											success: (res) => {
												const fileList = res.fileList || [];
												fileList.forEach(file => {
													uni.removeSavedFile({
														filePath: file.filePath,
														fail: () => {}
													});
												});
											},
											complete: () => {}
										});
										
										// 更新缓存大小显示
										setTimeout(() => {
											this.updateCacheSize();
											
											uni.hideLoading();
											uni.showToast({
												title: '缓存已清除',
												icon: 'success'
											});
										}, 500);
									},
									fail: () => {
										uni.hideLoading();
										uni.showToast({
											title: '清除缓存失败',
											icon: 'none'
										});
									}
								});
							} catch (error) {
								console.error('清除缓存出错:', error);
								uni.hideLoading();
								uni.showToast({
									title: '清除缓存失败',
									icon: 'none'
								});
							}
						}
					}
				});
			},
			
			// 更新缓存大小
			updateCacheSize() {
				try {
					uni.getStorageInfo({
						success: (res) => {
							// 计算存储大小（KB）
							const sizeInKB = res.currentSize || 0;
							
							// 转换为合适的单位
							if (sizeInKB < 1024) {
								this.cacheSize = sizeInKB + 'KB';
							} else {
								const sizeInMB = (sizeInKB / 1024).toFixed(1);
								this.cacheSize = sizeInMB + 'MB';
							}
						},
						fail: () => {
							this.cacheSize = '0KB';
						}
					});
				} catch (error) {
					console.error('获取缓存大小失败:', error);
					this.cacheSize = '未知';
				}
			},
			
			// 处理关于我们
			handleAbout() {
				this.showAboutPopup = true;
			},
			
			// 获取应用基础信息
			async fetchAppBaseInfo() {
				try {
					const res = await settings.getBaseSettings();
					if (res.code === 0 && res.data) {
						// 更新应用信息
						if (res.data.name) {
							this.appInfo.name = res.data.name;
						}
						if (res.data.description) {
							this.appInfo.description = res.data.description;
						}
					} else {
						console.error('获取应用基础信息失败:', res.message || '未知错误');
					}
				} catch (err) {
					console.error('请求应用基础信息接口出错:', err);
				}
			},
			
			// 处理协议点击
			handleAgreement(type) {
				const apiType = type === 'service' ? 'user' : 'privacy'
				
				// 跳转到协议页面
				uni.navigateTo({
					url: `/pages/agreement/index?type=${apiType}`
				})
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background-color: #f5f5f5;
		position: relative;
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
	}
	
	.back-btn {
		width: 88rpx;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-left: 7px;
	}
	
	.nav-title {
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
		font-size: 34rpx;
		font-weight: 500;
		color: #333333;
		max-width: 350rpx;
		text-align: center;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.content-scroll {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: #f5f5f5;
		z-index: 1;
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
	}
	
	.settings-wrapper {
		padding: 20rpx 0;
		padding-bottom: calc(env(safe-area-inset-bottom) + 40rpx);
	}
	
	.settings-section {
		margin-bottom: 20rpx;
	}
	
	.section-title {
		font-size: 28rpx;
		color: #999999;
		padding: 20rpx 30rpx;
	}
	
	.settings-item {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 30rpx;
		background-color: #ffffff;
		border-bottom: 1rpx solid #f5f5f5;
	}
	
	.settings-item:last-child {
		border-bottom: none;
	}
	
	.item-label {
		font-size: 30rpx;
		color: #333333;
	}
	
	.item-right {
		display: flex;
		align-items: center;
	}
	
	.item-value {
		font-size: 28rpx;
		color: #999999;
		margin-right: 10rpx;
	}
	
	.avatar-wrapper {
		background-color: transparent;
		padding: 0;
		margin: 0;
		line-height: 1;
		border: none;
	}
	
	.avatar-wrapper::after {
		border: none;
	}
	
	.avatar-preview {
		width: 80rpx;
		height: 80rpx;
		border-radius: 50%;
		margin-right: 10rpx;
	}
	
	.nickname-input-wrapper, .phone-input-wrapper {
		display: flex;
		align-items: center;
	}
	
	.nickname-input, .phone-input {
		text-align: right;
		font-size: 28rpx;
		color: #333333;
		padding-right: 10rpx;
		width: 200rpx;
	}
	
	.cache-size {
		font-size: 28rpx;
		color: #999999;
		margin-right: 10rpx;
	}
	
	.footer {
		padding: 40rpx 30rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
	.agreement-links {
		display: flex;
		align-items: center;
		margin-bottom: 40rpx;
	}
	
	.agreement-link {
		font-size: 26rpx;
		color: #666666;
	}
	
	.divider {
		margin: 0 20rpx;
		color: #cccccc;
	}
	
	/* 自定义关于我们弹窗样式 */
	.about-popup {
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
	
	.about-container {
		width: 600rpx;
		background-color: #ffffff;
		border-radius: 20rpx;
		overflow: hidden;
		display: flex;
		flex-direction: column;
	}
	
	.about-header {
		position: relative;
		padding: 30rpx 0;
		text-align: center;
		border-bottom: 1rpx solid #f5f5f5;
	}
	
	.about-title {
		font-size: 32rpx;
		font-weight: 500;
		color: #333333;
	}
	
	.about-content {
		padding: 40rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
	.app-name {
		font-size: 30rpx;
		font-weight: 500;
		color: #333333;
		margin-bottom: 30rpx;
		text-align: center;
	}
	
	.app-desc {
		font-size: 28rpx;
		color: #666666;
		text-align: center;
		line-height: 1.5;
	}
	
	.about-footer {
		padding: 20rpx 40rpx 40rpx;
		display: flex;
		justify-content: center;
	}
	
	.confirm-btn {
		width: 240rpx;
		height: 80rpx;
		background: linear-gradient(135deg, #fc3e2b 0%, #fa7154 100%);
		border-radius: 40rpx;
		color: #ffffff;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 30rpx;
	}
</style> 