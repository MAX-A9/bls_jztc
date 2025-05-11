<template>
	<view v-if="visible" class="login-modal">
		<view class="modal-mask" @tap="cancel"></view>
		<view class="modal-content">
			<view class="modal-header">
				<text class="modal-title">登录</text>
				<view class="close-btn" @tap="cancel">
					<uni-icons type="close" size="20" color="#999999"></uni-icons>
				</view>
			</view>
			
			<view class="modal-body">
				<view class="login-icon">
					<image src="/static/images/login-icon.png" mode="aspectFit"></image>
				</view>
				
				<text class="login-tips">请登录后继续操作</text>
				
				<button 
					class="wx-login-btn" 
					open-type="getUserInfo" 
					@tap="handleLogin" 
					:loading="loading"
				>
					<uni-icons v-if="!loading" type="weixin" size="20" color="#ffffff"></uni-icons>
					<text class="btn-text">{{ loading ? '登录中...' : '微信登录' }}</text>
				</button>
				
				<view class="agreement-area">
					<view class="checkbox" @tap="toggleAgreement">
						<view :class="['checkbox-inner', { checked: isAgreed }]">
							<view v-if="isAgreed" class="check-icon"></view>
						</view>
					</view>
					<view class="agreement-text">
						<text>我已阅读并同意</text>
						<text class="link" @tap.stop="goToAgreement('user')">《用户协议》</text>
						<text>和</text>
						<text class="link" @tap.stop="goToAgreement('privacy')">《隐私政策》</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	export default {
		name: 'LoginModal',
		data() {
			return {
				loading: false,
				isAgreed: true
			}
		},
		computed: {
			visible() {
				return this.$store.getters['user/showLoginModal'];
			}
		},
		methods: {
			// 处理登录
			async handleLogin() {
				if (!this.isAgreed) {
					uni.showToast({
						title: '请先同意用户协议和隐私政策',
						icon: 'none'
					});
					return;
				}
				
				if (this.loading) return;
				
				this.loading = true;
				try {
					await this.$store.dispatch('user/login');
					uni.showToast({
						title: '登录成功',
						icon: 'success'
					});
				} catch (error) {
					console.error('登录失败:', error);
					uni.showToast({
						title: '登录失败，请重试',
						icon: 'none'
					});
				} finally {
					this.loading = false;
				}
			},
			
			// 切换协议勾选状态
			toggleAgreement() {
				this.isAgreed = !this.isAgreed;
			},
			
			// 跳转至协议页面
			goToAgreement(type) {
				uni.navigateTo({
					url: `/pages/agreement/index?type=${type}`
				});
			},
			
			// 取消登录
			cancel() {
				this.$store.dispatch('user/hideLoginModal');
			}
		}
	}
</script>

<style>
	.login-modal {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 9999;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.modal-mask {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.6);
	}
	
	.modal-content {
		width: 80%;
		max-width: 600rpx;
		background-color: #FFFFFF;
		border-radius: 12rpx;
		overflow: hidden;
		position: relative;
		z-index: 1;
		box-shadow: 0 10rpx 20rpx rgba(0, 0, 0, 0.1);
	}
	
	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 30rpx;
		border-bottom: 1rpx solid #EEEEEE;
	}
	
	.modal-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333333;
	}
	
	.close-btn {
		padding: 10rpx;
	}
	
	.modal-body {
		padding: 40rpx 30rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
	.login-icon {
		width: 120rpx;
		height: 120rpx;
		margin-bottom: 30rpx;
	}
	
	.login-icon image {
		width: 100%;
		height: 100%;
	}
	
	.login-tips {
		font-size: 28rpx;
		color: #666666;
		margin-bottom: 40rpx;
		text-align: center;
	}
	
	.wx-login-btn {
		width: 100%;
		height: 80rpx;
		background: linear-gradient(135deg, #12B70D 0%, #3A9838 100%);
		border-radius: 40rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 30rpx;
		border: none;
	}
	
	.btn-text {
		font-size: 28rpx;
		color: #FFFFFF;
		margin-left: 10rpx;
	}
	
	.agreement-area {
		display: flex;
		align-items: center;
		margin-top: 20rpx;
	}
	
	.checkbox {
		width: 40rpx;
		height: 40rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-right: 10rpx;
	}
	
	.checkbox-inner {
		width: 32rpx;
		height: 32rpx;
		border: 1rpx solid #DDDDDD;
		border-radius: 6rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.checkbox-inner.checked {
		background-color: #fa5e44;
		border-color: #fa5e44;
	}
	
	.check-icon {
		width: 20rpx;
		height: 12rpx;
		border: 2rpx solid #FFFFFF;
		border-top: none;
		border-right: none;
		transform: rotate(-45deg);
	}
	
	.agreement-text {
		font-size: 24rpx;
		color: #999999;
		line-height: 1.4;
		flex: 1;
	}
	
	.link {
		color: #fa5e44;
	}
</style> 