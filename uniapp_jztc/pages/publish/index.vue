<template>
	<view class="publish-container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="draft-btn-nav" @tap="handleDraftClick">
					<text class="draft-text-nav">草稿</text>
				</view>
				<text class="nav-title">请选择发布类目</text>
			</view>
		</view>
		
		<!-- 占位元素 -->
		<view class="placeholder" :style="{ height: `calc(${statusBarHeight}px + ${navBarHeight}px)` }"></view>
		
		<view class="category-container">
			<view 
				class="category-card idle-card" 
				@click="handleCategoryClick(categoryList[0])"
			>
				<view class="card-content">
					<text class="category-name">{{ categoryList[0].name }}</text>
					<text class="category-desc">发布闲置物品，快速出售</text>
				</view>
			</view>
			
			<view 
				class="category-card info-card" 
				@click="handleCategoryClick(categoryList[1])"
			>
				<view class="card-content">
					<text class="category-name">{{ categoryList[1].name }}</text>
					<text class="category-desc">发布各类信息，满足需求</text>
				</view>
			</view>
		</view>
		<tab-bar :current-tab="tabIndex"></tab-bar>
	</view>
</template>

<script>
	import TabBar from '@/components/tab-bar/index.vue'
	import deviceInfo from '@/utils/device-info.js'
	
	export default {
		components: {
			TabBar
		},
		data() {
			return {
				tabIndex: -1,
				statusBarHeight: 0,
				navBarHeight: 44,
				categoryList: [
					{ name: '闲置发布', type: 'idle' },
					{ name: '信息发布', type: 'info' }
				]
			}
		},
		onLoad() {
			// 获取状态栏高度
			this.statusBarHeight = deviceInfo.getStatusBarHeight();
		},
		onShow() {
			this.tabIndex = 2
		},
		methods: {
			handleCategoryClick(category) {
				uni.navigateTo({
					url: `/pages/publish/${category.type}/index`
				})
			},
			handleDraftClick() {
				// 跳转到草稿页面
				uni.navigateTo({
					url: '/pages/publish/draft/index'
				})
			}
		}
	}
</script>

<style>
	.publish-container {
		padding-bottom: 120rpx;
		background-color: #ffffff;
		min-height: 100vh;
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
		justify-content: center;
	}
	
	.draft-btn-nav {
		position: absolute;
		left: 15rpx;
		font-size: 28rpx;
		color: #666666;
		z-index: 10;
	}
	
	.draft-text-nav {
		font-size: 28rpx;
		color: #666666;
	}
	
	.nav-title {
		font-size: 34rpx;
		font-weight: 500;
		color: #333333;
		max-width: 350rpx;
		text-align: center;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.placeholder {
		width: 100%;
	}
	
	.category-container {
		display: flex;
		flex-direction: column;
		padding: 40rpx 30rpx;
		gap: 40rpx;
	}
	
	.category-card {
		height: 240rpx;
		border-radius: 20rpx;
		position: relative;
		overflow: hidden;
		box-shadow: 0 8rpx 20rpx rgba(0, 0, 0, 0.1);
	}
	
	.idle-card {
		background-image: linear-gradient(to right, #ff8c00, #ff4500);
	}
	
	.info-card {
		background-image: linear-gradient(to right, #1e90ff, #4169e1);
	}
	
	.card-content {
		position: relative;
		z-index: 2;
		height: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 0 60rpx;
		text-align: center;
	}
	
	.category-name {
		font-size: 42rpx;
		font-weight: bold;
		color: #ffffff;
		margin-bottom: 20rpx;
		letter-spacing: 2rpx;
	}
	
	.category-desc {
		font-size: 28rpx;
		color: rgba(255, 255, 255, 0.9);
		max-width: 80%;
	}
</style>
