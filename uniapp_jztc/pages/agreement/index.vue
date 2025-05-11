<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">{{ title }}</text>
			</view>
		</view>
		
		<!-- 协议内容 -->
		<scroll-view 
			class="content-scroll" 
			scroll-y
			:style="{ top: `calc(${statusBarHeight}px + ${navBarHeight}px)` }"
		>
			<view class="agreement-content" v-if="!loading">
				<rich-text :nodes="content"></rich-text>
			</view>
			
			<!-- 加载提示 -->
			<view class="loading-container" v-if="loading">
				<uni-icons type="spinner-cycle" size="30" color="#fc3e2b"></uni-icons>
				<text class="loading-text">加载中...</text>
			</view>
			
			<!-- 错误提示 -->
			<view class="error-container" v-if="errorMsg">
				<uni-icons type="info" size="50" color="#cccccc"></uni-icons>
				<text class="error-text">{{ errorMsg }}</text>
				<button class="retry-btn" @tap="loadAgreementContent">重新加载</button>
			</view>
		</scroll-view>
	</view>
</template>

<script>
	import { agreement } from '@/apis/index.js';
	import deviceInfo from '@/utils/device-info.js';
	
	export default {
		data() {
			return {
				statusBarHeight: 0,
				navBarHeight: 44,
				title: '协议详情',
				content: '',
				type: '',
				loading: true,
				errorMsg: ''
			}
		},
		onLoad(options) {
			this.type = options.type || 'privacy';
			
			// 获取状态栏高度
			this.statusBarHeight = deviceInfo.getStatusBarHeight();
			
			// 获取协议内容
			this.loadAgreementContent();
		},
		methods: {
			// 处理返回按钮点击
			handleBack() {
				uni.navigateBack();
			},
			
			// 加载协议内容
			async loadAgreementContent() {
				this.loading = true;
				this.errorMsg = '';
				
				try {
					// 调用协议接口获取内容
					const res = await agreement.getAgreement({ type: this.type });
					
					if (res.code === 0 && res.data) {
						// 处理协议内容，支持富文本显示
						// 如果后端返回的是HTML格式，则直接使用
						this.content = res.data.content || `暂无${this.title}内容`;
						this.loading = false;
					} else {
						throw new Error(res.message || '获取协议内容失败');
					}
				} catch (error) {
					console.error('获取协议内容失败:', error);
					this.errorMsg = error.message || '获取协议内容失败，请重试';
					this.loading = false;
				}
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background-color: #ffffff;
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
		background-color: #ffffff;
		z-index: 1;
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
	}
	
	.agreement-content {
		padding: 30rpx;
		font-size: 28rpx;
		color: #333333;
		line-height: 1.8;
	}
	
	.loading-container, .error-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding-top: 200rpx;
	}
	
	.loading-text, .error-text {
		font-size: 28rpx;
		color: #999999;
		margin-top: 20rpx;
		text-align: center;
	}
	
	.error-text {
		margin-bottom: 40rpx;
	}
	
	.retry-btn {
		font-size: 28rpx;
		color: #ffffff;
		background-color: #fc3e2b;
		padding: 16rpx 60rpx;
		border-radius: 40rpx;
	}
</style> 