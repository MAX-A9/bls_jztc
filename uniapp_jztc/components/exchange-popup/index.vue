<template>
	<uni-popup ref="popup" type="bottom">
		<view class="exchange-popup">
			<view class="popup-header">
				<text class="popup-title">兑换会员</text>
				<view class="close-btn" @tap="closePopup">
					<uni-icons type="close" size="20" color="#333"></uni-icons>
				</view>
			</view>
			
			<view class="popup-content">
				<view class="product-info" v-if="productInfo">
					<image class="product-icon" :src="productInfo.thumbnail" mode="aspectFit"></image>
					<text class="product-name">{{productInfo.name}}</text>
					<text class="product-desc">需{{productInfo.duration}}天会员时长</text>
				</view>
				
				<view class="input-group">
					<view class="input-item">
						<input 
							class="phone-input"
							type="number" 
							v-model="phone" 
							maxlength="11" 
							placeholder="请输入领取该权益的手机号" 
							placeholder-class="placeholder"
						/>
					</view>
					<view class="input-item">
						<input 
							class="phone-input"
							type="number" 
							v-model="confirmPhone" 
							maxlength="11" 
							placeholder="请再次输入领取该权益的手机号" 
							placeholder-class="placeholder"
						/>
					</view>
				</view>
				
				<view class="confirm-btn" @tap="handleConfirm">
					<text class="btn-text">确认兑换</text>
				</view>
			</view>
		</view>
	</uni-popup>
</template>

<script>
	// 直接导入uni组件
	import uniPopup from '@dcloudio/uni-ui/lib/uni-popup/uni-popup'
	import uniIcons from '@dcloudio/uni-ui/lib/uni-icons/uni-icons'
	
	export default {
		name: 'exchange-popup',
		components: {
			uniPopup,
			uniIcons
		},
		data() {
			return {
				phone: '',
				confirmPhone: '',
				productInfo: null
			}
		},
		mounted() {
			// 确保组件初始化完成
			console.log('Exchange popup component mounted')
		},
		methods: {
			// 打开弹窗
			open(product) {
				console.log('Try to open popup with product:', product)
				this.productInfo = product
				// 清空输入框
				this.phone = ''
				this.confirmPhone = ''
				// 打开弹窗
				this.$nextTick(() => {
					if (this.$refs.popup) {
						this.$refs.popup.open('bottom')
					} else {
						console.error('Popup ref not found')
					}
				})
			},
			
			// 关闭弹窗
			closePopup() {
				this.$refs.popup.close()
			},
			
			// 确认兑换
			handleConfirm() {
				// 验证手机号
				if (!this.phone) {
					uni.showToast({
						title: '请输入手机号',
						icon: 'none'
					})
					return
				}
				
				if (!/^1\d{10}$/.test(this.phone)) {
					uni.showToast({
						title: '请输入正确的手机号',
						icon: 'none'
					})
					return
				}
				
				if (this.phone !== this.confirmPhone) {
					uni.showToast({
						title: '两次输入的手机号不一致',
						icon: 'none'
					})
					return
				}
				
				// 适配接口参数格式
				const data = {
					duration: this.productInfo.duration || 0,
					productName: this.productInfo.name || '',
					rechargeAccount: this.phone
				}
				
				// 触发确认事件
				this.$emit('confirm', data)
				
				// 关闭弹窗
				this.closePopup()
			}
		}
	}
</script>

<style>
	/* 兑换弹窗样式 */
	.exchange-popup {
		background-color: #ffffff;
		border-radius: 24rpx 24rpx 0 0;
		width: 100%;
	}
	
	.popup-header {
		position: relative;
		text-align: center;
		padding: 30rpx 0;
		border-bottom: 1px solid #f0f0f0;
	}
	
	.popup-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333333;
	}
	
	.close-btn {
		position: absolute;
		right: 30rpx;
		top: 50%;
		transform: translateY(-50%);
		width: 40rpx;
		height: 40rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.popup-content {
		padding: 30rpx;
	}
	
	.product-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-bottom: 40rpx;
	}
	
	.product-icon {
		width: 100rpx;
		height: 100rpx;
		border-radius: 16rpx;
		margin-bottom: 16rpx;
	}
	
	.product-name {
		font-size: 32rpx;
		color: #333333;
		font-weight: bold;
		margin-bottom: 8rpx;
	}
	
	.product-desc {
		font-size: 26rpx;
		color: #fc3e2b;
	}
	
	.input-group {
		margin-bottom: 40rpx;
	}
	
	.input-item {
		background: #f8f8f8;
		border-radius: 12rpx;
		padding: 20rpx;
		margin-bottom: 20rpx;
	}
	
	.phone-input {
		width: 100%;
		height: 40rpx;
		font-size: 28rpx;
		color: #333333;
	}
	
	.placeholder {
		color: #999999;
	}
	
	.confirm-btn {
		background: linear-gradient(135deg, #fc3e2b 0%, #fa7154 100%);
		border-radius: 12rpx;
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.btn-text {
		font-size: 32rpx;
		color: #ffffff;
		font-weight: 500;
	}
</style> 