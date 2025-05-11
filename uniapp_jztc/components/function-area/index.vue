<template>
	<view class="grid-container" v-if="isGlobalEnabled">
		<scroll-view 
			class="scroll-row scroll-view-container" 
			scroll-x="true" 
			:enhanced="true"
			@scroll="handleScroll"
		>
			<!-- 小程序列表 -->
			<view 
				class="grid-item" 
				v-for="(item) in miniProgramList" 
				:key="'mini-' + item.id"
				@click="navigateToMiniProgram(item)"
			>
				<image :src="item.logo" mode="aspectFill" class="logo-image"></image>
				<text class="item-text">{{ item.name }}</text>
			</view>
		</scroll-view>
		
		<!-- 滚动指示器 -->
		<view class="scroll-indicator" v-if="miniProgramList.length > 0">
			<view 
				class="indicator-bar"
				:style="{ transform: `translateX(${currentPage * 100}%)` }"
			></view>
		</view>
	</view>
</template>

<script>
	import { getMiniProgramList } from '@/apis/content.js'
	
	export default {
		name: 'FunctionArea',
		data() {
			return {
				scrollLeft: 0,
				contentWidth: 0, // 内容总宽度
				viewWidth: 0, // 可视区域宽度
				currentPage: 0, // 当前页码
				isGlobalEnabled: true, // 默认显示，加载完成后根据接口返回值更新
				miniProgramList: [] // 小程序列表
			}
		},
		created() {
			// 移除自动加载，由父组件控制加载时机
			// this.loadMiniProgramList()
		},
		mounted() {
			this.initScrollBar()
		},
		methods: {
			// 加载小程序列表
			async loadMiniProgramList() {
				try {
					const res = await getMiniProgramList()
					if (res.code === 0 && res.data) {
						// 更新全局启用状态
						this.isGlobalEnabled = res.data.isGlobalEnabled
						
						// 如果启用，则更新小程序列表
						if (this.isGlobalEnabled && res.data.list) {
							// 按order字段排序
							this.miniProgramList = res.data.list
								.filter(item => item.isEnabled) // 只展示启用的小程序
								.sort((a, b) => a.order - b.order)
						}
					} else {
						console.error('获取小程序列表失败:', res.message)
					}
				} catch (error) {
					console.error('加载小程序列表异常:', error)
				}
			},
			// 跳转到小程序
			navigateToMiniProgram(miniProgram) {
				uni.navigateToMiniProgram({
					appId: miniProgram.appId,
					success(res) {
						console.log('跳转成功', res)
					},
					fail(err) {
						console.error('跳转失败', err)
						uni.showToast({
							title: '跳转失败',
							icon: 'none'
						})
					}
				})
			},
			initScrollBar() {
				const query = uni.createSelectorQuery().in(this)
				query.select('.scroll-row').boundingClientRect(data => {
					if (data) {
						this.viewWidth = data.width
					}
				}).exec()
			},
			handleScroll(e) {
				const { scrollLeft, scrollWidth } = e.detail
				// 计算滚动比例
				const maxScroll = scrollWidth - this.viewWidth
				// 计算当前页码（0到1之间的小数）
				this.currentPage = maxScroll > 0 ? scrollLeft / maxScroll : 0
			}
		}
	}
</script>

<style>
	.grid-container {
		padding: 30rpx 0;
		background-color: #ffffff;
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
	.scroll-row {
		width: 100%;
		white-space: nowrap;
		padding: 0 30rpx;
		box-sizing: border-box;
	}
	
	.grid-item {
		display: inline-flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		width: 20%;
		box-sizing: border-box;
		padding: 0 10rpx;
		vertical-align: top;
	}
	
	.grid-item:nth-child(n+6) {
		width: auto;
		padding: 0 20rpx;
	}
	
	.logo-image {
		width: 110rpx;
		height: 110rpx;
		border-radius: 50%;
		margin-bottom: 12rpx;
	}
	
	.item-text {
		font-size: 24rpx;
		color: #333333;
		line-height: 1.4;
		white-space: nowrap;
		text-align: center;
		height: 34rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	/* 滚动指示器样式 */
	.scroll-indicator {
		display: flex;
		justify-content: center;
		align-items: center;
		margin-top: 16rpx;
		position: relative;
		width: 80rpx;
		height: 8rpx;
		background-color: #dddddd;
		border-radius: 4rpx;
		margin-left: auto;
		margin-right: auto;
		overflow: hidden;
	}
	
	.indicator-bar {
		position: absolute;
		left: 0;
		top: 0;
		width: 50%;
		height: 100%;
		border-radius: 4rpx;
		background-color: #fc3e2b;
		transition: transform 0.3s;
	}
</style> 