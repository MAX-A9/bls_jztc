<template>
	<view class="container">
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
				<!-- 快递模板分类导航 -->
				<view class="template-nav">
					<scroll-view 
						class="template-scroll" 
						scroll-x 
						:show-scrollbar="false"
					>
						<view class="template-list">
							<view 
								v-for="(item, index) in templateList" 
								:key="index"
								class="template-item"
								:class="{ active: currentTemplate === index }"
								@tap="switchTemplate(index)"
							>
								<text>{{ item }}</text>
							</view>
							<!-- 更多模板作为最后一个选项 -->
							<view class="template-item more" @tap="handleMoreTemplate">
								<text>更多模板</text>
								<uni-icons type="right" size="12" color="#666666"></uni-icons>
							</view>
						</view>
					</scroll-view>
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
							<text class="action-separator">|</text>
							<text class="action-text">地址簿</text>
						</view>
					</view>
					
					<!-- 分割线 -->
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
								<text class="title">请填写收件人信息</text>
								<text class="desc">支持智能识别文本、图片中的地址</text>
							</view>
						</view>
						<view class="address-action">
							<text class="action-separator">|</text>
							<text class="action-text">地址簿</text>
						</view>
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
	import deviceMixin from '@/mixins/device.js'
	
	export default {
		mixins: [deviceMixin],
		data() {
			return {
				pageTitle: '寄快递',
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
				currentTemplate: 0,
				templateList: ['常规', '箱包', '酒类', '蔬菜', '钟表首饰', '医药保健', '文体娱乐'],
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
			switchTemplate(index) {
				this.currentTemplate = index
			},
			handleMoreTemplate() {
				uni.navigateTo({
					url: '/pages/express/template/more'
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
		height: calc(100vh - var(--window-top) - 140rpx - env(safe-area-inset-bottom));
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
	}
	
	.content-wrapper {
		padding: 0 20rpx;
		position: relative;
		padding-bottom: calc(40rpx + env(safe-area-inset-bottom) + 100rpx);
	}
	
	/* 地址卡片样式 */
	.address-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 0 30rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
		position: relative;
		z-index: 1;
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
		margin-left: 20rpx;
		height: 40rpx;
	}
	
	.action-separator {
		font-size: 24rpx;
		color: #EEEEEE;
		margin-right: 12rpx;
		margin-left: 4rpx;
		line-height: 40rpx;
	}
	
	.action-text {
		font-size: 24rpx;
		color: #262626;
		line-height: 40rpx;
		font-weight: 600;
	}
	
	.divider {
		height: 1rpx;
		background-color: #EEEEEE;
		width: calc(100% - 60rpx);
		margin-left: 60rpx;
	}
	
	/* 模板导航样式 */
	.template-nav {
		display: flex;
		align-items: center;
		padding: 20rpx 0;
		margin-bottom: 20rpx;
		position: relative;
		height: 80rpx;
	}
	
	.template-scroll {
		flex: 1;
		white-space: nowrap;
		padding-right: 20rpx;
		height: 100%;
	}
	
	.template-list {
		display: inline-flex;
		padding-left: 4rpx;
		height: 100%;
		align-items: center;
	}
	
	.template-item {
		padding: 12rpx 24rpx;
		margin-right: 8rpx;
		border-radius: 26rpx;
		background-color: #F5F5F5;
		display: inline-flex;
		align-items: center;
		height: 52rpx;
	}
	
	.template-item text {
		font-size: 26rpx;
		color: #666666;
		line-height: 1;
	}
	
	.template-item.active {
		background-color: #E6F0FF;
	}
	
	.template-item.active text {
		color: #007AFF;
		font-weight: 500;
	}
	
	/* 更多模板选项样式 */
	.template-item.more {
		display: inline-flex;
		align-items: center;
		gap: 4rpx;
		height: 52rpx;
	}
	
	/* 物品信息卡片样式 */
	.goods-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		margin-top: 20rpx;
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
		margin-top: 20rpx;
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
		margin-top: 20rpx;
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
		background-color: #FFFFFF;
		padding: 16rpx 30rpx calc(16rpx + env(safe-area-inset-bottom));
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
</style> 