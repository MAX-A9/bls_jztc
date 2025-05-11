<template>
	<view class="swiper-banner" v-if="isGlobalEnabled">
		<swiper 
			class="swiper" 
			:circular="true" 
			:autoplay="true" 
			:interval="3000" 
			:duration="500"
			@change="handleChange"
		>
			<swiper-item v-for="(item, index) in bannerList" :key="index" @click="handleClick(item)">
				<image :src="item.image" mode="aspectFill" class="banner-image"></image>
			</swiper-item>
		</swiper>
		
		<!-- 指示点 -->
		<view class="dots">
			<view 
				class="dot" 
				v-for="(item, index) in bannerList" 
				:key="index"
				:class="{ active: current === index }"
			></view>
		</view>
	</view>
</template>

<script>
	import { get } from '@/utils/request.js';
	
	export default {
		name: 'SwiperBanner',
		data() {
			return {
				current: 0,
				bannerList: [],
				isGlobalEnabled: false
			}
		},
		created() {
			// 移除自动加载，由父组件控制加载时机
			// this.fetchBannerData();
		},
		methods: {
			async fetchBannerData() {
				try {
					const result = await get('/wx/banner/list');
					if (result && result.code === 0) {
						this.isGlobalEnabled = result.data.isGlobalEnabled;
						if (this.isGlobalEnabled && result.data.list) {
							// 按order排序
							this.bannerList = result.data.list
								.filter(item => item.isEnabled)
								.sort((a, b) => a.order - b.order);
						}
					}
				} catch (error) {
					console.error('获取轮播图数据失败', error);
				}
			},
			handleChange(e) {
				this.current = e.detail.current;
			},
			handleClick(item) {
				if (!item || !item.linkUrl) return;
				
				if (item.linkType === 'page') {
					uni.navigateTo({
						url: '/' + item.linkUrl,
						fail: (err) => {
							console.error('页面跳转失败', item.linkUrl, err);
						}
					});
				} else if (item.linkType === 'webview') {
					uni.navigateTo({
						url: '/pages/webview/index?url=' + encodeURIComponent(item.linkUrl),
						fail: (err) => {
							console.error('网页跳转失败', item.linkUrl, err);
						}
					});
				} else if (item.linkType === 'miniprogram') {
					// 跳转到其他小程序
					uni.navigateToMiniProgram({
						appId: item.linkUrl,
						fail: (err) => {
							console.error('小程序跳转失败', item.linkUrl, err);
							uni.showToast({
								title: '跳转失败',
								icon: 'none'
							});
						}
					});
				}
			}
		}
	}
</script>

<style>
	.swiper-banner {
		position: relative;
		padding: 0 20rpx;
		margin-top: 20rpx;
	}
	
	.swiper {
		height: 160rpx;
		border-radius: 20rpx;
		overflow: hidden;
	}
	
	.banner-image {
		width: 100%;
		height: 100%;
		border-radius: 20rpx;
	}
	
	/* 指示点样式 */
	.dots {
		position: absolute;
		bottom: 20rpx;
		right: 30rpx;
		display: flex;
		justify-content: flex-end;
		gap: 8rpx;
	}
	
	.dot {
		width: 12rpx;
		height: 12rpx;
		background-color: rgba(255, 255, 255, 0.5);
		border-radius: 50%;
		transition: all 0.3s;
	}
	
	.dot.active {
		width: 24rpx;
		background-color: #ffffff;
		border-radius: 6rpx;
	}
</style> 