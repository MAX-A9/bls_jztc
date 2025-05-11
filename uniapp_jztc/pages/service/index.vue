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
				<view class="service-type-card">
					<view class="card-title">选择服务类型</view>
					<view class="service-type-list">
						<view 
							class="service-type-item" 
							:class="{ active: currentService === 'install' }"
							@tap="selectService('install')"
						>
							<view class="icon-wrapper">
								<uni-icons type="tools" size="32" :color="currentService === 'install' ? '#007AFF' : '#666666'"></uni-icons>
							</view>
							<text class="type-name">安装</text>
						</view>
						
						<view 
							class="service-type-item" 
							:class="{ active: currentService === 'repair' }"
							@tap="selectService('repair')"
						>
							<view class="icon-wrapper">
								<uni-icons type="refresh" size="32" :color="currentService === 'repair' ? '#007AFF' : '#666666'"></uni-icons>
							</view>
							<text class="type-name">维修</text>
						</view>
						
						<view 
							class="service-type-item" 
							:class="{ active: currentService === 'clean' }"
							@tap="selectService('clean')"
						>
							<view class="icon-wrapper">
								<uni-icons type="clear" size="32" :color="currentService === 'clean' ? '#007AFF' : '#666666'"></uni-icons>
							</view>
							<text class="type-name">清洗</text>
						</view>
					</view>
				</view>
				
				<!-- 安装类型卡片 -->
				<view class="install-type-card">
					<view class="card-title">安装需求</view>
					<view class="input-section">
						<textarea
							class="demand-input"
							v-model="demandText"
							placeholder="请输入您的需求"
							:maxlength="200"
							auto-height
						></textarea>
						<view class="word-count">{{demandText.length}}/200</view>
					</view>
					<view class="keywords-section">
						<view class="keywords-title">常见需求</view>
						<view class="keywords-list">
							<view 
								class="keyword-item"
								v-for="(item, index) in keywords"
								:key="index"
								@tap="selectKeyword(item)"
							>
								<text>{{item}}</text>
							</view>
						</view>
					</view>
				</view>

				<!-- 在安装类型卡片后添加服务地址卡片 -->
				<view class="address-card">
					<view class="card-title">服务地址</view>
					<view class="address-content" @tap="handleSelectAddress">
						<view class="address-placeholder" v-if="!selectedAddress">
							<uni-icons type="location" size="20" color="#666666"></uni-icons>
							<text>请选择服务地址</text>
						</view>
						<view class="address-info" v-else>
							<view class="address-row">
								<text class="name">{{selectedAddress.name}}</text>
								<text class="phone">{{selectedAddress.phone}}</text>
							</view>
							<view class="address-detail">
								<text>{{selectedAddress.address}}</text>
							</view>
						</view>
						<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
					</view>
				</view>

				<!-- 在服务地址卡片后添加预约时间和物品图片备注卡片 -->
				<view class="time-card">
					<view class="card-title">预约时间</view>
					<view class="time-content" @tap="handleSelectTime">
						<view class="time-placeholder" v-if="!selectedTime">
							<uni-icons type="calendar" size="20" color="#666666"></uni-icons>
							<text>请选择预约时间</text>
						</view>
						<view class="time-info" v-else>
							<text>{{selectedTime}}</text>
						</view>
						<uni-icons type="right" size="16" color="#CCCCCC"></uni-icons>
					</view>
				</view>

				<!-- 将物品图片和备注分成两个卡片 -->
				<view class="image-card">
					<view class="card-title">物品图片</view>
					<view class="image-section">
						<view class="image-list">
							<view 
								class="image-item" 
								v-for="(item, index) in imageList" 
								:key="index"
								@tap="previewImage(index)"
							>
								<image :src="item" mode="aspectFill"></image>
								<view class="delete-btn" @tap.stop="deleteImage(index)">
									<uni-icons type="clear" size="12" color="#FFFFFF"></uni-icons>
								</view>
							</view>
							<view 
								class="upload-btn" 
								v-if="imageList.length < 4"
								@tap="chooseImage"
							>
								<uni-icons type="camera" size="24" color="#666666"></uni-icons>
								<text>上传图片</text>
							</view>
						</view>
						<text class="image-tip">最多上传4张图片</text>
					</view>
				</view>

				<view class="note-card">
					<view class="card-title">备注信息</view>
					<view class="note-section">
						<textarea
							class="note-input"
							v-model="noteText"
							placeholder="请输入备注信息（选填）"
							:maxlength="100"
							auto-height
						></textarea>
						<view class="word-count">{{noteText.length}}/100</view>
					</view>
				</view>
			</view>
		</scroll-view>

		<!-- 底部结算区域 -->
		<view class="settlement-bar">
			<view class="submit-btn" @tap="handleSubmitOrder">
				<view class="btn-content">
					<text class="main-text">立即预约</text>
					<view class="sub-text">
						<text>查看附近师傅报价</text>
						<uni-icons type="right" size="12" color="#E6F0FF"></uni-icons>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import deviceAdapter from '@/mixins/device-adapter.js'
	
	export default {
		mixins: [deviceAdapter],
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
		data() {
			return {
				pageTitle: '安装维修',
				isAgreed: false, // 是否同意协议
				currentService: '', // 当前选中的服务类型
				demandText: '', // 安装需求输入框的值
				keywords: ['柜子', '空调', '洗衣机', '电视', '床', '沙发', '衣柜', '冰箱', '热水器', '油烟机'], // 关键词列表
				selectedAddress: null, // 选中的地址
				selectedTime: '', // 选中的预约时间
				imageList: [], // 物品图片列表
				noteText: '', // 备注信息输入框的值
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
			handleShowDetail() {
				// 显示价格明细
				uni.showModal({
					title: '费用明细',
					content: '服务费：¥15.0\n总计：¥15.0',
					showCancel: false
				})
			},
			handleReadAgreement() {
				// 跳转到服务协议页面
				uni.navigateTo({
					url: '/pages/agreement/service'
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
			selectService(type) {
				this.currentService = type
			},
			selectKeyword(keyword) {
				if (this.demandText) {
					// 如果已有文本，在末尾添加逗号和关键词
					this.demandText = this.demandText.trim() + '，' + keyword
				} else {
					// 如果没有文本，直接设置关键词
					this.demandText = keyword
				}
			},
			handleSelectAddress() {
				uni.navigateTo({
					url: '/pages/address/list?type=service'
				})
			},
			handleSelectTime() {
				uni.navigateTo({
					url: '/pages/service/time'
				})
			},
			previewImage(index) {
				uni.previewImage({
					urls: this.imageList,
					current: index
				})
			},
			deleteImage(index) {
				this.imageList.splice(index, 1)
			},
			chooseImage() {
				uni.chooseImage({
					count: 4 - this.imageList.length,
					sizeType: ['compressed'],
					sourceType: ['album', 'camera'],
					success: (res) => {
						this.imageList = [...this.imageList, ...res.tempFilePaths]
					}
				})
			},
			handleCheckPrice() {
				uni.navigateTo({
					url: '/pages/service/nearby-price'
				})
			}
		}
	}
</script>

<style>
	/* 保留基础容器和导航栏样式 */
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
		justify-content: center;
		align-items: center;
		box-shadow: 0 -2rpx 8rpx rgba(0, 0, 0, 0.04);
		z-index: 100;
	}
	
	.submit-btn {
		background: linear-gradient(135deg, #FF9500 0%, #FF5E3A 100%);
		padding: 12rpx 40rpx;
		border-radius: 48rpx;
		width: 600rpx;
		display: flex;
		justify-content: center;
	}
	
	.btn-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 2rpx;
	}
	
	.main-text {
		font-size: 28rpx;
		color: #FFFFFF;
		font-weight: 500;
		line-height: 1.2;
	}
	
	.sub-text {
		display: flex;
		align-items: center;
		gap: 4rpx;
	}
	
	.sub-text text {
		font-size: 20rpx;
		color: #E6F0FF;
		line-height: 1.2;
	}
	
	/* 服务类型卡片样式 */
	.service-type-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 30rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
		margin-top: var(--card-gap);
	}
	
	.card-title {
		font-size: 28rpx;
		color: #262626;
		font-weight: 600;
		margin-bottom: 30rpx;
	}
	
	.service-type-list {
		display: flex;
		justify-content: space-between;
		gap: 20rpx;
	}
	
	.service-type-item {
		flex: 1;
		background-color: #F8F8F8;
		border-radius: 12rpx;
		padding: 20rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		transition: all 0.3s ease;
	}
	
	.service-type-item.active {
		background-color: #E6F0FF;
	}
	
	.icon-wrapper {
		width: 80rpx;
		height: 80rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 16rpx;
	}
	
	.type-name {
		font-size: 28rpx;
		color: #262626;
		font-weight: 500;
	}
	
	.service-type-item.active .type-name {
		color: #007AFF;
	}
	
	/* 安装类型卡片样式 */
	.install-type-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-top: var(--card-gap);
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.card-title {
		font-size: 28rpx;
		color: #262626;
		font-weight: 600;
		margin-bottom: 30rpx;
	}
	
	.input-section {
		position: relative;
		margin-bottom: 20rpx;
		width: 100%;
	}
	
	.demand-input {
		width: 100%;
		min-height: 120rpx;
		background-color: #F8F8F8;
		border-radius: 12rpx;
		padding: 20rpx;
		font-size: 28rpx;
		color: #262626;
		line-height: 1.5;
		box-sizing: border-box;
	}
	
	.word-count {
		position: absolute;
		right: 20rpx;
		bottom: 20rpx;
		font-size: 24rpx;
		color: #999999;
		background-color: #F8F8F8;
		padding: 0 4rpx;
		z-index: 1;
	}
	
	.keywords-section {
		margin-top: 30rpx;
	}
	
	.keywords-title {
		font-size: 26rpx;
		color: #666666;
		margin-bottom: 16rpx;
	}
	
	.keywords-list {
		display: flex;
		flex-wrap: wrap;
		gap: 16rpx;
	}
	
	.keyword-item {
		padding: 12rpx 24rpx;
		background-color: #F8F8F8;
		border-radius: 26rpx;
		transition: all 0.3s ease;
	}
	
	.keyword-item:active {
		background-color: #E6F0FF;
	}
	
	.keyword-item text {
		font-size: 26rpx;
		color: #666666;
	}

	/* 服务地址卡片样式 */
	.address-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-top: var(--card-gap);
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.card-title {
		font-size: 28rpx;
		color: #262626;
		font-weight: 600;
		margin-bottom: 30rpx;
	}
	
	.address-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		min-height: 120rpx;
		background-color: #F8F8F8;
		border-radius: 12rpx;
		padding: 20rpx;
	}
	
	.address-placeholder {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 12rpx;
	}
	
	.address-placeholder text {
		font-size: 28rpx;
		color: #999999;
	}
	
	.address-info {
		flex: 1;
		margin-right: 20rpx;
	}
	
	.address-row {
		display: flex;
		align-items: center;
		gap: 20rpx;
		margin-bottom: 8rpx;
	}
	
	.name {
		font-size: 28rpx;
		color: #262626;
		font-weight: 500;
	}
	
	.phone {
		font-size: 28rpx;
		color: #666666;
	}
	
	.address-detail {
		font-size: 26rpx;
		color: #666666;
		line-height: 1.4;
	}

	/* 预约时间卡片样式 */
	.time-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-top: var(--card-gap);
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.card-title {
		font-size: 28rpx;
		color: #262626;
		font-weight: 600;
		margin-bottom: 30rpx;
	}
	
	.time-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		min-height: 100rpx;
		background-color: #F8F8F8;
		border-radius: 12rpx;
		padding: 20rpx;
	}
	
	.time-placeholder {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 12rpx;
	}
	
	.time-placeholder text {
		font-size: 28rpx;
		color: #999999;
	}
	
	.time-info {
		flex: 1;
		font-size: 28rpx;
		color: #262626;
	}

	/* 物品图片卡片样式 */
	.image-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-top: var(--card-gap);
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.card-title {
		font-size: 28rpx;
		color: #262626;
		font-weight: 600;
		margin-bottom: 30rpx;
	}
	
	.image-section {
		margin-bottom: 30rpx;
	}
	
	.image-list {
		display: flex;
		flex-wrap: wrap;
		gap: 20rpx;
		margin-bottom: 12rpx;
	}
	
	.image-item {
		width: 160rpx;
		height: 160rpx;
		position: relative;
	}
	
	.image-item image {
		width: 100%;
		height: 100%;
		border-radius: 8rpx;
	}
	
	.delete-btn {
		position: absolute;
		top: -10rpx;
		right: -10rpx;
		width: 36rpx;
		height: 36rpx;
		background-color: rgba(0, 0, 0, 0.5);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.upload-btn {
		width: 160rpx;
		height: 160rpx;
		background-color: #F8F8F8;
		border-radius: 8rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 8rpx;
	}
	
	.upload-btn text {
		font-size: 24rpx;
		color: #666666;
	}
	
	.image-tip {
		font-size: 24rpx;
		color: #999999;
	}
	
	/* 备注信息卡片样式 */
	.note-card {
		background-color: #FFFFFF;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-top: var(--card-gap);
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	}
	
	.card-title {
		font-size: 28rpx;
		color: #262626;
		font-weight: 600;
		margin-bottom: 30rpx;
	}
	
	.note-section {
		position: relative;
	}
	
	.note-input {
		width: 100%;
		min-height: 100rpx;
		background-color: #F8F8F8;
		border-radius: 12rpx;
		padding: 20rpx;
		font-size: 28rpx;
		color: #262626;
		line-height: 1.5;
		box-sizing: border-box;
	}
	
	.word-count {
		position: absolute;
		right: 20rpx;
		bottom: 20rpx;
		font-size: 24rpx;
		color: #999999;
		background-color: #F8F8F8;
		padding: 0 4rpx;
		z-index: 1;
	}
</style> 