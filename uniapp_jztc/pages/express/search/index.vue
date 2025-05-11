<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">查快递</text>
			</view>
		</view>

		<!-- 搜索区域 -->
		<view class="search-area">
			<view class="search-box">
				<uni-icons type="search" size="18" color="#666666"></uni-icons>
				<input 
					type="text" 
					v-model="searchText"
					placeholder="请输入运单号" 
					placeholder-class="placeholder"
					@confirm="handleSearch"
				/>
				<view class="scan-btn" @tap="handleScan">
					<uni-icons type="scan" size="20" color="#007AFF"></uni-icons>
				</view>
			</view>
		</view>

		<!-- 历史记录 -->
		<view class="history-section" v-if="historyList.length > 0">
			<view class="section-header">
				<text class="section-title">历史记录</text>
				<view class="clear-btn" @tap="handleClearHistory">
					<uni-icons type="trash" size="14" color="#999999"></uni-icons>
					<text>清空</text>
				</view>
			</view>
			<view class="history-list">
				<view 
					class="history-item"
					v-for="(item, index) in historyList"
					:key="index"
					@tap="handleHistoryClick(item)"
				>
					<view class="item-info">
						<text class="express-number">{{ item.number }}</text>
						<text class="express-company">{{ item.company }}</text>
					</view>
					<text class="express-status" :class="item.status">{{ item.statusText }}</text>
				</view>
			</view>
		</view>

		<!-- 快递公司列表 -->
		<view class="company-section">
			<view class="section-header">
				<text class="section-title">常用快递号码大全</text>
			</view>
			<view class="company-grid">
				<view 
					class="company-item"
					v-for="(item, index) in companyList"
					:key="index"
					@tap="handleCompanyClick(item)"
				>
					<image :src="item.logo" mode="aspectFit" class="company-logo"></image>
					<text class="company-name">{{ item.name }}</text>
					<text class="company-phone">{{ item.phone }}</text>
				</view>
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
				searchText: '',
				historyList: [
					{
						number: 'SF1234567890',
						company: '顺丰快递',
						status: 'delivered',
						statusText: '已签收'
					},
					{
						number: 'YT9876543210',
						company: '圆通快递',
						status: 'shipping',
						statusText: '运输中'
					}
				],
				companyList: [
					{ name: '顺丰快递', logo: '/static/express/sf.png', phone: '95338' },
					{ name: '中通快递', logo: '/static/express/zt.png', phone: '95311' },
					{ name: '韵达快递', logo: '/static/express/yd.png', phone: '95546' },
					{ name: '申通快递', logo: '/static/express/st.png', phone: '95543' },
					{ name: '圆通快递', logo: '/static/express/yt.png', phone: '95554' },
					{ name: '百世快递', logo: '/static/express/bs.png', phone: '95320' },
					{ name: '京东快递', logo: '/static/express/jd.png', phone: '950616' },
					{ name: '邮政快递', logo: '/static/express/yz.png', phone: '11183' }
				]
			}
		},
		methods: {
			handleBack() {
				uni.navigateBack()
			},
			handleSearch() {
				if (!this.searchText) {
					uni.showToast({
						title: '请输入运单号',
						icon: 'none'
					})
					return
				}
				console.log('搜索运单号:', this.searchText)
			},
			handleScan() {
				uni.scanCode({
					success: (res) => {
						this.searchText = res.result
						this.handleSearch()
					}
				})
			},
			handleClearHistory() {
				uni.showModal({
					title: '提示',
					content: '确定要清空历史记录吗？',
					success: (res) => {
						if (res.confirm) {
							this.historyList = []
						}
					}
				})
			},
			handleHistoryClick(item) {
				uni.navigateTo({
					url: `/pages/express/detail?number=${item.number}`
				})
			},
			handleCompanyClick(item) {
				uni.makePhoneCall({
					phoneNumber: item.phone,
					fail(err) {
						console.log('拨打电话失败:', err)
					}
				})
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background-color: #F8F8F8;
		padding-bottom: env(safe-area-inset-bottom);
	}
	
	.nav-bar {
		background-color: #FFFFFF;
		padding: 0 30rpx;
	}
	
	.nav-content {
		height: 44px;
		display: flex;
		align-items: center;
		position: relative;
	}
	
	.back-btn {
		padding: 20rpx;
		margin-left: -20rpx;
		position: absolute;
		left: 0;
		z-index: 1;
	}
	
	.nav-title {
		font-size: 34rpx;
		font-weight: 500;
		color: #333333;
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
		width: 400rpx;
		text-align: center;
	}
	
	.search-area {
		background-color: #FFFFFF;
		padding: 20rpx 30rpx;
	}
	
	.search-box {
		display: flex;
		align-items: center;
		background-color: #F5F5F5;
		border-radius: 32rpx;
		padding: 0 24rpx;
		height: 64rpx;
	}
	
	.search-box input {
		flex: 1;
		height: 100%;
		margin: 0 16rpx;
		font-size: 28rpx;
	}
	
	.placeholder {
		color: #999999;
		font-size: 28rpx;
	}
	
	.scan-btn {
		padding: 10rpx;
	}
	
	.history-section {
		margin-top: 20rpx;
		background-color: #FFFFFF;
		padding: 30rpx;
	}
	
	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;
	}
	
	.section-title {
		font-size: 30rpx;
		font-weight: 600;
		color: #333333;
	}
	
	.clear-btn {
		display: flex;
		align-items: center;
		gap: 4rpx;
	}
	
	.clear-btn text {
		font-size: 24rpx;
		color: #999999;
	}
	
	.history-list {
		display: flex;
		flex-direction: column;
		gap: 20rpx;
	}
	
	.history-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20rpx;
		background-color: #F8F8F8;
		border-radius: 12rpx;
	}
	
	.item-info {
		display: flex;
		flex-direction: column;
		gap: 8rpx;
	}
	
	.express-number {
		font-size: 28rpx;
		color: #333333;
		font-weight: 500;
	}
	
	.express-company {
		font-size: 24rpx;
		color: #666666;
	}
	
	.express-status {
		font-size: 26rpx;
	}
	
	.express-status.delivered {
		color: #67C23A;
	}
	
	.express-status.shipping {
		color: #409EFF;
	}
	
	.company-section {
		margin-top: 20rpx;
		background-color: #FFFFFF;
		padding: 30rpx;
	}
	
	.company-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 30rpx;
		padding: 20rpx 0;
	}
	
	.company-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8rpx;
	}
	
	.company-logo {
		width: 80rpx;
		height: 80rpx;
		margin-bottom: 4rpx;
	}
	
	.company-name {
		font-size: 24rpx;
		color: #333333;
		margin-bottom: 2rpx;
	}
	
	.company-phone {
		font-size: 24rpx;
		color: #007AFF;
	}
</style> 