<template>
	<view class="container" :style="pageStyle">
		<!-- 背景 -->
		<view class="header-bg"></view>
		
		<!-- 导航栏背景 -->
		<view class="nav-bg" :style="{ height: navigationBarHeight + 'px' }"></view>
		
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">{{pageTitle}}</text>
			</view>
		</view>

		<!-- 内容区域 -->
		<scroll-view 
			class="content-scroll" 
			scroll-y 
			:style="{ top: navigationBarHeight + 'px' }"
		>
			<view class="content-wrapper">
				<!-- 寄件流程卡片 -->
				<view class="process-card">
					<view class="process-list">
						<view class="process-item">
							<view class="process-number">1</view>
							<text class="process-text">填写寄件人信息</text>
						</view>
						<view class="process-arrow">
							<uni-icons type="right" size="14" color="#CCCCCC"></uni-icons>
						</view>
						<view class="process-item">
							<view class="process-number">2</view>
							<text class="process-text">复制商家退换地址</text>
						</view>
						<view class="process-arrow">
							<uni-icons type="right" size="14" color="#CCCCCC"></uni-icons>
						</view>
						<view class="process-item">
							<view class="process-number">3</view>
							<text class="process-text">匹配退货地址</text>
						</view>
					</view>
				</view>

				<!-- 寄件人收件人卡片 -->
				<view class="address-card">
					<!-- 寄件人 -->
					<view class="address-item" @tap="handleSelectAddress('send')">
						<view class="address-info">
							<view class="hexagon-wrapper">
								<view class="hexagon-label">
									<text>寄</text>
								</view>
							</view>
							<view class="info-content">
								<text class="title">请填写寄件人信息</text>
								<text class="desc">支持智能识别文本、图片中的地址</text>
							</view>
						</view>
						<view class="address-action">
							<view class="address-book">
								<text class="separator">|</text>
								<text class="book-text">地址簿</text>
							</view>
						</view>
					</view>
					
					<!-- 寄件人和收件人之间的分割线 -->
					<view class="divider"></view>
					
					<!-- 收件人 -->
					<view class="address-item" @tap="handleSelectAddress('receive')">
						<view class="address-info">
							<view class="hexagon-wrapper">
								<view class="hexagon-label">
									<text>收</text>
								</view>
							</view>
							<view class="info-content">
								<text class="title">请填写商家退货地址</text>
								<text class="desc">支持电商平台退货</text>
							</view>
						</view>
						<view class="address-action">
							<view class="upload-btn" @tap.stop="handleUploadImage">
								<uni-icons type="camera-filled" size="16" color="#007AFF"></uni-icons>
								<text class="upload-text">上传退货截图</text>
							</view>
						</view>
					</view>

					<!-- 智能粘贴地址上方的分割线 -->
					<view class="divider full-width"></view>

					<!-- 智能粘贴地址 -->
					<view class="paste-address" @tap="handlePasteAddress">
						<text class="paste-text">智能粘贴地址</text>
					</view>
				</view>

				<!-- 物品信息卡片 -->
				<view class="goods-card">
					<!-- 物品信息 -->
					<view class="goods-item" @tap="handleSelectGoods">
						<view class="goods-label">物品信息</view>
						<view class="goods-value">
							<text class="placeholder-text">请选择物品信息</text>
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
					
					<!-- 分割线 -->
					<view class="goods-divider"></view>
					
					<!-- 付款方式 -->
					<view class="goods-item" @tap="handleSelectPayment">
						<view class="goods-label">付款方式</view>
						<view class="goods-value">
							<text class="value-text">在线支付</text>
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
				</view>

				<!-- 快递公司选择卡片 -->
				<view class="express-card">
					<view class="express-header">
						<text class="express-title">选择快递公司</text>
					</view>
					<scroll-view 
						class="express-scroll" 
						scroll-x 
						:show-scrollbar="false"
					>
						<view class="express-list">
							<view 
								class="express-item"
								v-for="(item, index) in expressList"
								:key="index"
								:class="{ active: currentExpress === index }"
								@tap="selectExpress(index)"
							>
								<image class="express-logo" :src="item.logo" mode="aspectFit"></image>
								<text class="express-name">{{ item.name }}</text>
								<view class="express-price">
									<text class="current-price">¥{{ item.currentPrice }}</text>
									<text class="original-price">¥{{ item.originalPrice }}</text>
								</view>
								<text class="weight-info">续重{{ item.weight }}元/kg</text>
							</view>
						</view>
					</scroll-view>
				</view>

				<!-- 服务选项卡片 -->
				<view class="service-card">
					<!-- 期望上门时间 -->
					<view class="service-item" @tap="handleSelectTime">
						<view class="service-label">期望上门时间</view>
						<view class="service-value">
							<text class="placeholder-text">请选择上门时间</text>
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
					
					<!-- 分割线 -->
					<view class="service-divider"></view>
					
					<!-- 保价服务 -->
					<view class="service-item" @tap="handleInsurance">
						<view class="service-label">保费</view>
						<view class="service-value">
							<text class="insurance-tip">选填</text>
							<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
						</view>
					</view>
				</view>
			</view>
		</scroll-view>

		<!-- 底部结算区域 -->
		<view class="settlement-bar">
			<view class="price-section">
				<view class="price-row">
					<text class="price-label">预计</text>
					<text class="price-value">¥15.0</text>
					<view class="price-detail" @tap="handleShowDetail">
						<text>明细</text>
						<uni-icons type="bottom" size="12" color="#666666"></uni-icons>
					</view>
				</view>
				<view class="agreement-row">
					<view class="checkbox" @tap="toggleAgreement">
						<uni-icons 
							:type="isAgreed ? 'checkbox-filled' : 'circle'" 
							size="18" 
							:color="isAgreed ? '#007AFF' : '#CCCCCC'"
						></uni-icons>
					</view>
					<text class="agreement-text" @tap="handleReadAgreement">
						阅读并同意<text class="agreement-link">《快递服务协议》</text>
					</text>
				</view>
			</view>
			<view class="submit-btn" @tap="handleSubmitOrder">
				<text>立即下单</text>
			</view>
		</view>
	</view>
</template>

<script>
	import deviceAdapter from '@/mixins/device-adapter.js'
	
	export default {
		mixins: [deviceAdapter],
		data() {
			return {
				pageTitle: '网购退货',
				senderInfo: {
					name: '',
					phone: '',
					address: ''
				},
				receiverInfo: {
					name: '',
					phone: '',
					address: ''
				},
				currentExpress: -1,
				expressList: [
					{
						logo: '/static/express/sf.png',
						name: '顺丰快递',
						originalPrice: '18.0',
						currentPrice: '15.0',
						weight: '5'
					},
					{
						logo: '/static/express/zt.png',
						name: '中通快递',
						originalPrice: '15.0',
						currentPrice: '12.0',
						weight: '4'
					},
					{
						logo: '/static/express/yd.png',
						name: '韵达快递',
						originalPrice: '15.0',
						currentPrice: '12.0',
						weight: '4'
					},
					// ... 可以添加更多快递公司
				],
				hasInsurance: false, // 是否选择保价服务
				isAgreed: false, // 是否同意协议
			}
		},
		computed: {
			// 使用统一的布局尺寸
			pageStyle() {
				return {
					'--nav-height': `${this.layoutSize.navHeight}px`,
					'--settlement-height': `${this.layoutSize.settlementHeight}rpx`,
					'--content-padding': `${this.layoutSize.contentPadding}rpx`,
					'--card-gap': `${this.layoutSize.cardGap}rpx`,
					'--content-bottom': `${this.layoutSize.contentBottom}rpx`
				}
			}
		},
		methods: {
			handleBack() {
				uni.navigateBack({
					delta: 1,
					fail: () => {
						uni.switchTab({
							url: '/pages/index/index'
						})
					}
				})
			},
			handleSelectAddress(type) {
				uni.navigateTo({
					url: `/pages/address/list?type=${type}`
				})
			},
			handleSelectGoods() {
				uni.navigateTo({
					url: '/pages/express/goods/select'
				})
			},
			handleSelectPayment() {
				uni.navigateTo({
					url: '/pages/express/payment/select'
				})
			},
			selectExpress(index) {
				this.currentExpress = index
			},
			handleSelectTime() {
				uni.navigateTo({
					url: '/pages/express/time/select'
				})
			},
			handleInsurance() {
				uni.navigateTo({
					url: '/pages/express/insurance/set'
				})
			},
			handleShowDetail() {
				// 显示价格明细
				uni.showModal({
					title: '费用明细',
					content: '快递费：¥15.0\n保价费：¥0.0\n总计：¥15.0',
					showCancel: false
				})
			},
			handleReadAgreement() {
				// 跳转到服务协议页面
				uni.navigateTo({
					url: '/pages/agreement/express'
				})
			},
			toggleAgreement() {
				this.isAgreed = !this.isAgreed
			},
			handleSubmitOrder() {
				if (!this.isAgreed) {
					uni.showToast({
						title: '请先阅读并同意服务协议',
						icon: 'none'
					})
					return
				}
				uni.showLoading({
					title: '提交中...'
				})
				// 这里添加提交订单的逻辑
			},
			handleUploadImage() {
				uni.chooseImage({
					count: 1,
					success: (res) => {
						console.log('上传退货截图:', res.tempFilePaths[0])
						// 这里可以添加处理图片的逻辑
					}
				})
			},
			handlePasteAddress() {
				// 获取剪贴板内容
				uni.getClipboardData({
					success: (res) => {
						if (res.data) {
							console.log('剪贴板内容:', res.data)
							// 这里可以添加解析地址的逻辑
							uni.showToast({
								title: '已识别地址',
								icon: 'success'
							})
						} else {
							uni.showToast({
								title: '剪贴板为空',
								icon: 'none'
							})
						}
					}
				})
			}
		}
	}
</script>

<style>
	.container {
		width: 100%;
		min-height: 100vh;
		position: relative;
		background-color: #F8F8F8;
	}
	
	/* 背景渐变 */
	.header-bg {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		height: 600rpx;
		background: linear-gradient(
			to right,
			#d7e8fe,
			#aacffc
		);
		mask-image: linear-gradient(
			to bottom,
			rgba(0, 0, 0, 1) 70%,
			rgba(0, 0, 0, 0)
		);
		-webkit-mask-image: linear-gradient(
			to bottom,
			rgba(0, 0, 0, 1) 70%,
			rgba(0, 0, 0, 0)
		);
		z-index: 0;
	}
	
	/* 导航栏背景 */
	.nav-bg {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		background: linear-gradient(
			to right,
			#d7e8fe,
			#aacffc
		);
		z-index: 2;
	}
	
	.nav-bar {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 2;
	}
	
	.nav-content {
		height: 44px;
		display: flex;
		align-items: center;
		position: relative;
	}
	
	.back-btn {
		position: absolute;
		left: 0;
		width: 44px;
		height: 44px;
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 2;
	}
	
	.nav-title {
		flex: 1;
		text-align: center;
		font-size: 16px;
		color: #262626;
		font-weight: 500;
	}
	
	/* 内容区域样式 */
	.content-scroll {
		flex: 1;
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 1;
		height: calc(100vh - var(--nav-height) - var(--settlement-height));
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
	}
	
	.content-wrapper {
		padding: 0 var(--content-padding);
		position: relative;
		padding-bottom: var(--content-bottom);
	}
	
	/* 地址卡片样式 */
	.address-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 0 30rpx;
			box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
		position: relative;
		z-index: 1;
		margin-top: var(--card-gap);
	}
	
	.address-item {
		display: flex;
		align-items: center;
		padding: 30rpx 0;
		justify-content: space-between;
	}
	
	.address-info {
		display: flex;
		align-items: center;
		flex: 1;
	}
	
	.hexagon-wrapper {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-right: 20rpx;
	}
	
	.hexagon-label {
		width: 40rpx;
		height: 40rpx;
		flex-shrink: 0;
		position: relative;
		background: linear-gradient(135deg, #007AFF 0%, #409EFF 100%);
		border-radius: 8rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
		aspect-ratio: 1;
	}
	
	/* 添加光泽效果 */
	.hexagon-label::after {
		content: '';
		position: absolute;
		top: 0;
		left: -50%;
		width: 200%;
		height: 100%;
		background: linear-gradient(
			to bottom,
			rgba(255, 255, 255, 0.2) 0%,
			rgba(255, 255, 255, 0.1) 30%,
			rgba(255, 255, 255, 0) 50%
		);
		transform: rotate(-30deg);
	}
	
	.hexagon-label text {
		color: #FFFFFF;
		font-size: 22rpx;
		font-weight: 500;
		position: relative;
		z-index: 1;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
		line-height: 1;
	}
	
	.info-content {
		flex: 1;
		display: flex;
		flex-direction: column;
	}
	
	.title {
		font-size: 28rpx;
		color: #262626;
		font-weight: 600;
		margin-bottom: 8rpx;
	}
	
	.desc {
		font-size: 24rpx;
		color: #999999;
	}
	
	.address-action {
		display: flex;
		align-items: center;
		margin-left: 12rpx;
	}
	
	.upload-btn {
		display: flex;
		align-items: center;
		gap: 4rpx;
		background-color: #F0F7FF;
		padding: 12rpx 24rpx;
		border-radius: 8rpx;
		height: 72rpx;
		min-width: 180rpx;
		justify-content: center;
	}
	
	.upload-text {
		font-size: 24rpx;
		color: #007AFF;
		font-weight: 500;
	}
	
	.divider {
		height: 1rpx;
		background-color: #EEEEEE;
		width: calc(100% - 60rpx);
		margin-left: 60rpx;
	}
	
	.divider.full-width {
		width: 100%;
		margin-left: 0;
	}
	
	/* 物品信息卡片样式 */
	.goods-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		margin-top: var(--card-gap);
		padding: 0 30rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.goods-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 30rpx 0;
	}
	
	.goods-label {
		font-size: 28rpx;
		color: #262626;
		font-weight: 500;
	}
	
	.goods-value {
		display: flex;
		align-items: center;
		gap: 8rpx;
	}
	
	.placeholder-text {
		font-size: 28rpx;
		color: #999999;
	}
	
	.value-text {
		font-size: 28rpx;
		color: #262626;
	}
	
	.goods-divider {
		height: 1rpx;
		background-color: #EEEEEE;
		width: 100%;
	}
	
	/* 快递公司选择卡片样式 */
	.express-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		margin-top: var(--card-gap);
		padding: 20rpx 0;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.express-header {
		padding: 0 30rpx;
		margin-bottom: 20rpx;
	}
	
	.express-title {
		font-size: 28rpx;
		color: #262626;
		font-weight: 600;
	}
	
	.express-scroll {
		width: 100%;
		white-space: nowrap;
	}
	
	.express-list {
		display: inline-flex;
		padding: 0 20rpx;
	}
	
	.express-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 20rpx;
		margin-right: 20rpx;
		background-color: #F8F8F8;
		border-radius: 12rpx;
		width: 200rpx;
	}
	
	.express-item.active {
		background-color: #E6F0FF;
	}
	
	.express-logo {
		width: 80rpx;
		height: 80rpx;
		margin-bottom: 12rpx;
	}
	
	.express-name {
		font-size: 26rpx;
		color: #262626;
		margin-bottom: 8rpx;
	}
	
	.express-price {
		display: flex;
		align-items: center;
		gap: 8rpx;
		margin-bottom: 8rpx;
	}
	
	.current-price {
		font-size: 28rpx;
		color: #FC3E2B;
		font-weight: 600;
	}
	
	.original-price {
		font-size: 24rpx;
		color: #999999;
		text-decoration: line-through;
	}
	
	.weight-info {
		font-size: 22rpx;
		color: #666666;
	}
	
	/* 服务选项卡片样式 */
	.service-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		margin-top: var(--card-gap);
		padding: 0 30rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.service-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 30rpx 0;
	}
	
	.service-label {
		font-size: 28rpx;
		color: #262626;
		font-weight: 500;
	}
	
	.service-value {
		display: flex;
		align-items: center;
		gap: 8rpx;
	}
	
	.insurance-tip {
		font-size: 26rpx;
		color: #999999;
		margin-right: 8rpx;
	}
	
	.service-divider {
		height: 1rpx;
		background-color: #EEEEEE;
		width: 100%;
	}
	
	/* 底部结算栏样式 */
	.settlement-bar {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		height: var(--settlement-height);
		background-color: #FFFFFF;
		padding: 0 30rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 -2rpx 8rpx rgba(0, 0, 0, 0.04);
		z-index: 100;
	}
	
	.price-section {
		flex: 1;
		margin-right: 30rpx;
		padding-left: 4rpx;
	}
	
	.price-row {
		display: flex;
		align-items: center;
		margin-bottom: 4rpx;
	}
	
	.price-label {
		font-size: 26rpx;
		color: #666666;
		margin-right: 8rpx;
	}
	
	.price-value {
		font-size: 32rpx;
		color: #FC3E2B;
		font-weight: 600;
		margin-right: 12rpx;
	}
	
	.price-detail {
		display: flex;
		align-items: center;
		gap: 4rpx;
	}
	
	.price-detail text {
		font-size: 24rpx;
		color: #666666;
	}
	
	.agreement-row {
		margin-top: 4rpx;
		display: flex;
		align-items: center;
		gap: 0;
	}
	
	.checkbox {
		display: flex;
		align-items: center;
		padding: 4rpx 2rpx 4rpx 0;
	}
	
	.agreement-text {
		font-size: 24rpx;
		color: #999999;
		flex: 1;
		margin-left: 0;
	}
	
	.agreement-link {
		color: #007AFF;
	}
	
	.submit-btn {
		background: linear-gradient(135deg, #007AFF 0%, #409EFF 100%);
		padding: 20rpx 40rpx;
		border-radius: 40rpx;
	}
	
	.submit-btn text {
		font-size: 28rpx;
		color: #FFFFFF;
		font-weight: 500;
	}
	
	/* 寄件流程卡片样式 */
	.process-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.process-list {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	
	.process-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12rpx;
	}
	
	.process-number {
		width: 40rpx;
		height: 40rpx;
		background: linear-gradient(135deg, #007AFF 0%, #409EFF 100%);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #FFFFFF;
		font-size: 24rpx;
		font-weight: 500;
	}
	
	.process-text {
		font-size: 24rpx;
		color: #333333;
		white-space: nowrap;
	}
	
	.process-arrow {
		display: flex;
		align-items: center;
		padding: 0 20rpx;
		margin-bottom: 20rpx;
	}
	
	.address-book {
		display: flex;
		align-items: center;
		margin-left: 20rpx;
		height: 40rpx;
	}
	
	.separator {
		font-size: 24rpx;
		color: #EEEEEE;
		margin: 0 12rpx;
		line-height: 40rpx;
	}
	
	.book-text {
		font-size: 24rpx;
		color: #262626;
		line-height: 40rpx;
		font-weight: 600;
	}
	
	/* 智能粘贴地址样式 */
	.paste-address {
		display: flex;
		align-items: center;
		padding: 24rpx 0;
		justify-content: center;
	}
	
	.paste-text {
		font-size: 26rpx;
		color: #007AFF;
		font-weight: 500;
	}
</style> 