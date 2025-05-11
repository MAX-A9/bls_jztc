<template>
	<view class="container">
		<!-- 顶部搜索区域 -->
		<view class="search-header">
			<view class="search-box">
				<uni-icons type="search" size="16" color="#666666"></uni-icons>
				<input 
					type="text" 
					v-model="searchKeyword"
					placeholder="搜索城市" 
					placeholder-class="placeholder"
					confirm-type="search"
					@input="handleSearch"
				/>
			</view>
		</view>
		
		<scroll-view 
			class="content-scroll" 
			scroll-y 
			:show-scrollbar="false" 
			:enhanced="true"
			:scroll-top="scrollTop"
		>
			<!-- 当前定位 -->
			<view class="location-section">
				<view class="section-title">当前地区</view>
				<view class="current-location" @tap="selectLocation(currentLocation)">
					<view class="location-icon">
						<uni-icons type="location" size="18" color="#fc3e2b"></uni-icons>
					</view>
					<view class="location-info">
						<text class="location-name">{{ currentLocation }}</text>
						<text class="location-status">当前地区</text>
					</view>
				</view>
			</view>
			
			<!-- 最近访问 -->
			<view class="location-section" v-if="recentLocations.length > 0">
				<view class="section-title">最近访问</view>
				<view class="location-list">
					<view 
						class="location-item" 
						v-for="(item, index) in recentLocations" 
						:key="index"
						@tap="selectLocation(item)"
					>
						<text>{{ item }}</text>
					</view>
				</view>
			</view>
			
			<!-- 热门城市 -->
			<view class="location-section">
				<view class="section-title">热门城市</view>
				<view class="location-grid">
					<view 
						class="grid-item" 
						v-for="item in hotRegions" 
						:key="item.id"
						@tap="selectLocation(item.name)"
					>
						<text>{{ item.name }}</text>
					</view>
				</view>
			</view>
			
			<!-- 城市列表 -->
			<view class="city-list">
				<view class="city-section" v-for="(section, letter) in filteredCityList" :key="letter">
					<view class="city-letter" :id="'letter-'+letter">{{ letter }}</view>
					<view class="city-items">
						<view 
							class="city-item" 
							v-for="city in section" 
							:key="city.id"
							@tap="selectLocation(city.name)"
						>
							<text>{{ city.name }}</text>
						</view>
					</view>
				</view>
			</view>
		</scroll-view>
		
		<!-- 字母索引 -->
		<view class="letter-index" 
			@touchstart="handleTouchStart" 
			@touchmove="handleTouchMove" 
			@touchend="handleTouchEnd"
		>
			<view 
				class="letter-item" 
				v-for="letter in indexLetters" 
				:key="letter"
				:data-letter="letter"
				@click="scrollToLetter(letter)"
			>
				<text>{{ letter }}</text>
			</view>
		</view>
		
		<!-- 字母提示 -->
		<view class="letter-indicator" v-if="showLetterIndicator">
			<text>{{ currentLetter }}</text>
		</view>
		
		<!-- 加载提示 -->
		<view v-if="loading" class="loading-tip">
			<view class="loading-icon"></view>
			<text>加载中...</text>
		</view>
		
		<!-- 错误提示 -->
		<view v-if="error" class="error-tip">
			<text>{{ error }}</text>
			<button class="retry-btn" @tap="fetchRegionList">重新加载</button>
		</view>
	</view>
</template>

<script>
	import { mapState, mapGetters, mapActions } from 'vuex';
	import { getStringFirstLetter, getPinyin } from '@/utils/pinyin.js';
	
	export default {
		data() {
			return {
				currentLocation: '',
				searchKeyword: '',
				recentLocations: [],
				scrollTop: 0,
				letterPositions: {},
				showLetterIndicator: false,
				currentLetter: '',
				touchStartY: 0,
				touchStartTime: 0,
				isTouching: false,
				letterIndexRect: null, // 字母索引容器的位置信息
			}
		},
		computed: {
			...mapState('region', ['regionList', 'loading', 'error']),
			...mapGetters('region', ['hotRegions']),
			
			// 过滤后的城市列表，分组显示
			filteredCityList() {
				// 从Vuex中获取分组后的地区列表
				const groupedList = this.$store.getters['region/groupedRegionList'];
				
				// 如果有搜索关键词，则过滤
				if (!this.searchKeyword) {
					return groupedList;
				}
				
				const keyword = this.searchKeyword.toLowerCase();
				const result = {};
				
				// 过滤包含关键词的城市
				for (const letter in groupedList) {
					const filteredCities = groupedList[letter].filter(city => {
						// 城市名称包含关键词
						if (city.name.toLowerCase().includes(keyword)) {
							return true;
						}
						
						// 城市拼音首字母匹配
						if (getStringFirstLetter(city.name).toLowerCase() === keyword.charAt(0).toLowerCase()) {
							return true;
						}
						
						// 城市全拼包含关键词
						const fullPinyin = getPinyin(city.name, { separator: '' }).toLowerCase();
						if (fullPinyin.includes(keyword)) {
							return true;
						}
						
						// 城市拼音首字母缩写包含关键词
						const pinyinInitials = getPinyin(city.name, { pattern: 'first', separator: '' }).toLowerCase();
						if (pinyinInitials.includes(keyword)) {
							return true;
						}
						
						return false;
					});
					
					if (filteredCities.length > 0) {
						result[letter] = filteredCities;
					}
				}
				
				return result;
			},
			
			// 索引字母列表
			indexLetters() {
				return Object.keys(this.filteredCityList).sort();
			}
		},
		onLoad() {
			// 获取缓存中的当前位置
			const savedLocation = uni.getStorageSync('currentLocation');
			if (savedLocation) {
				this.currentLocation = savedLocation;
				console.log('使用已保存的位置:', savedLocation);
			} else {
				// 如果没有保存位置，显示加载中状态
				this.currentLocation = '加载中...';
			}
			
			// 获取最近访问的城市
			const recentLocations = uni.getStorageSync('recentLocations');
			if (recentLocations) {
				this.recentLocations = JSON.parse(recentLocations);
			}
			
			// 加载地区列表数据
			this.fetchRegionList();
		},
		onReady() {
			// 获取字母索引容器的位置信息
			setTimeout(() => {
				const query = uni.createSelectorQuery().in(this);
				query.select('.letter-index').boundingClientRect(rect => {
					this.letterIndexRect = rect;
				}).exec();
			}, 300);
		},
		methods: {
			...mapActions('region', {
				getRegionList: 'getRegionList'
			}),
			
			// 获取地区列表
			fetchRegionList() {
				this.getRegionList({ status: 0 });
				
				// 添加回调，确保在地区列表加载完成后再设置默认位置
				this.$nextTick(() => {
					// 使用普通的watch方式监听regionList变化
					let watchExecuted = false;
					const checkRegionList = () => {
						// 只在第一次有数据时执行
						if (this.regionList && this.regionList.length > 0 && !watchExecuted) {
							watchExecuted = true;
							
							// 检查是否需要更新当前位置
							const savedLocation = uni.getStorageSync('currentLocation');
							if (!savedLocation || savedLocation === '未知地区' || savedLocation === '加载中...' || savedLocation === '请选择地区') {
								// 如果没有保存的位置或使用了默认位置，则使用第一个可用地区
								if (this.regionList.length > 0) {
									this.currentLocation = this.regionList[0].name;
									uni.setStorageSync('currentLocation', this.currentLocation);
									console.log('位置选择页更新默认位置:', this.currentLocation);
								}
							}
							
							// 获取字母索引位置
							setTimeout(() => {
								this.getLetterPositions();
							}, 500);
						}
					};
					
					// 初始检查
					checkRegionList();
					
					// 如果初始检查没有数据，设置定时器定期检查
					if (!watchExecuted) {
						const timer = setInterval(() => {
							checkRegionList();
							// 如果执行了或30秒后仍未执行，则清除定时器
							if (watchExecuted || !this.regionList) {
								clearInterval(timer);
							}
						}, 500);
						
						// 30秒后自动清除定时器防止内存泄漏
						setTimeout(() => {
							clearInterval(timer);
						}, 30000);
					}
				});
			},
			
			handleSearch(e) {
				// 搜索逻辑已通过计算属性实现
			},
			
			// 选择城市
			selectLocation(cityName) {
				// 查找选择的城市对象，获取完整信息
				let selectedRegion = null;
				
				// 在所有地区中查找匹配的城市
				this.regionList.forEach(region => {
					if (region.name === cityName) {
						selectedRegion = region;
					}
				});
				
				if (!selectedRegion) {
					console.error('未找到选择的城市对象:', cityName);
					// 尝试根据名称创建一个城市对象
					selectedRegion = {
						id: 0,
						name: cityName
					};
				}
				
				console.log('选择城市:', selectedRegion);
				
				// 更新当前位置信息
				uni.setStorageSync('currentLocation', selectedRegion.name);
				// 更新地区ID
				uni.setStorageSync('currentLocationId', selectedRegion.id);
				
				console.log('保存地区ID到本地存储:', selectedRegion.id);
				
				// 更新最近访问城市
				let recent = this.recentLocations.filter(item => item !== selectedRegion.name);
				recent.unshift(selectedRegion.name);
				if (recent.length > 3) {
					recent = recent.slice(0, 3);
				}
				this.recentLocations = recent;
				uni.setStorageSync('recentLocations', JSON.stringify(recent));
				
				// 通过eventChannel将选择的地区信息传回给首页
				const eventChannel = this.getOpenerEventChannel();
				if (eventChannel && eventChannel.emit) {
					eventChannel.emit('locationSelected', {
						id: selectedRegion.id, 
						name: selectedRegion.name
					});
					console.log('通过事件通道发送地区信息:', selectedRegion.id, selectedRegion.name);
				} else {
					console.error('无法获取事件通道');
				}
				
				// 返回上一页
				uni.navigateBack();
				
				// 显示提示
				uni.showToast({
					title: '已切换到' + selectedRegion.name,
					icon: 'success'
				});
			},
			
			// 获取所有字母索引的位置信息
			getLetterPositions() {
				const query = uni.createSelectorQuery().in(this);
				this.indexLetters.forEach(letter => {
					query.select('#letter-' + letter).boundingClientRect();
				});
				query.select('.content-scroll').boundingClientRect();
				query.exec(res => {
					if (res && res.length > 0) {
						// 最后一个结果是scroll-view的位置信息
						const scrollViewInfo = res[res.length - 1];
						// 前面的结果是各个字母的位置信息
						this.letterPositions = {};
						
						for (let i = 0; i < this.indexLetters.length; i++) {
							const letter = this.indexLetters[i];
							const position = res[i];
							if (position) {
								// 计算相对于scroll-view的位置
								this.letterPositions[letter] = position.top - scrollViewInfo.top;
							}
						}
					}
				});
			},
			
			// 滚动到指定字母区域
			scrollToLetter(letter) {
				// 使用缓存的位置信息
				if (this.letterPositions && this.letterPositions[letter] !== undefined) {
					this.scrollTop = this.letterPositions[letter];
					// 添加触感反馈
					if (uni.vibrateShort) {
						uni.vibrateShort();
					}
					return;
				}
				
				// 修复字母索引功能
				const query = uni.createSelectorQuery().in(this);
				query.select('#letter-' + letter).boundingClientRect(data => {
					if (data) {
						// 获取scroll-view组件
						const scrollView = uni.createSelectorQuery().in(this).select('.content-scroll');
						scrollView.boundingClientRect(scrollData => {
							if (scrollData) {
								// 计算在scroll-view中的相对位置
								const scrollTop = data.top - scrollData.top;
								// 设置scroll-view的滚动位置
								this.scrollTop = scrollTop;
								
								// 添加触感反馈
								if (uni.vibrateShort) {
									uni.vibrateShort();
								}
							}
						}).exec();
					}
				}).exec();
			},
			
			// 处理触摸开始事件
			handleTouchStart(e) {
				this.touchStartY = e.touches[0].clientY;
				this.touchStartTime = Date.now();
				this.isTouching = true;
				
				// 获取当前触摸的字母
				const letter = this.getLetterFromTouch(e.touches[0].clientY);
				if (letter) {
					this.currentLetter = letter;
					this.showLetterIndicator = true;
					this.scrollToLetter(letter);
				}
			},
			
			// 处理触摸移动事件
			handleTouchMove(e) {
				if (!this.isTouching) return;
				
				// 获取当前触摸的字母
				const letter = this.getLetterFromTouch(e.touches[0].clientY);
				if (letter && letter !== this.currentLetter) {
					this.currentLetter = letter;
					this.scrollToLetter(letter);
				}
			},
			
			handleTouchEnd(e) {
				this.showLetterIndicator = false;
			},
			
			// 根据触摸位置获取对应的字母
			getLetterFromTouch(touchY) {
				// 使用缓存的字母索引容器位置信息
				if (!this.letterIndexRect) return null;
				
				// 计算触摸点在容器内的相对位置
				const offsetY = touchY - this.letterIndexRect.top;
				
				// 计算索引
				const itemHeight = this.letterIndexRect.height / this.indexLetters.length;
				let index = Math.floor(offsetY / itemHeight);
				
				// 边界处理
				if (index < 0) index = 0;
				if (index >= this.indexLetters.length) index = this.indexLetters.length - 1;
				
				// 获取对应字母
				return this.indexLetters[index];
			},
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background-color: #f5f5f5;
		position: relative;
	}
	
	.search-header {
		padding: 20rpx 30rpx;
		background-color: #ffffff;
	}
	
	.search-box {
		display: flex;
		align-items: center;
		background-color: #f5f5f5;
		padding: 0 24rpx;
		border-radius: 36rpx;
		height: 72rpx;
	}
	
	.search-box input {
		flex: 1;
		margin-left: 16rpx;
		font-size: 28rpx;
		height: 100%;
	}
	
	.placeholder {
		color: #999999;
		font-size: 28rpx;
	}
	
	.content-scroll {
		height: calc(100vh - 112rpx);
	}
	
	.location-section {
		background-color: #ffffff;
		margin-bottom: 20rpx;
		padding: 0 30rpx;
	}
	
	.section-title {
		font-size: 28rpx;
		color: #999999;
		padding: 20rpx 0;
	}
	
	.current-location {
		display: flex;
		align-items: center;
		padding: 20rpx 0;
		border-top: 1rpx solid #f5f5f5;
	}
	
	.location-icon {
		margin-right: 16rpx;
	}
	
	.location-info {
		display: flex;
		flex-direction: column;
	}
	
	.location-name {
		font-size: 30rpx;
		color: #333333;
		margin-bottom: 6rpx;
	}
	
	.location-status {
		font-size: 24rpx;
		color: #999999;
	}
	
	.location-list {
		display: flex;
		flex-wrap: wrap;
	}
	
	.location-item {
		padding: 20rpx 0;
		margin-right: 40rpx;
		font-size: 30rpx;
		color: #333333;
	}
	
	.location-grid {
		display: flex;
		flex-wrap: wrap;
		padding-bottom: 20rpx;
	}
	
	.grid-item {
		width: 25%;
		height: 80rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 28rpx;
		color: #333333;
		margin-bottom: 20rpx;
	}
	
	.grid-item text {
		padding: 10rpx 30rpx;
		background-color: #f5f5f5;
		border-radius: 8rpx;
	}
	
	.city-list {
		background-color: #ffffff;
	}
	
	.city-section {
		padding: 0 30rpx;
	}
	
	.city-letter {
		font-size: 28rpx;
		color: #999999;
		padding: 20rpx 0;
		background-color: #f5f5f5;
		margin: 0 -30rpx;
		padding-left: 30rpx;
	}
	
	.city-items {
		padding: 10rpx 0;
	}
	
	.city-item {
		padding: 20rpx 0;
		font-size: 30rpx;
		color: #333333;
		border-bottom: 1rpx solid #f5f5f5;
	}
	
	.letter-index {
		position: fixed;
		right: 20rpx;
		top: 50%;
		transform: translateY(-50%);
		display: flex;
		flex-direction: column;
		background-color: rgba(255, 255, 255, 0.7);
		border-radius: 30rpx;
		padding: 10rpx 0;
		z-index: 10;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.1);
	}
	
	.letter-item {
		width: 40rpx;
		height: 40rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
		padding: 4rpx 0;
	}
	
	.letter-item:active {
		background-color: rgba(252, 62, 43, 0.1);
	}
	
	.letter-item:active text {
		color: #fc3e2b;
	}
	
	.letter-item text {
		font-size: 24rpx;
		color: #666666;
		font-weight: 500;
	}
	
	/* 优化字母索引区域的触摸体验 */
	.letter-index {
		padding: 10rpx 6rpx;
	}
	
	.letter-item {
		width: 50rpx;
	}
	
	.letter-indicator {
		position: fixed;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		width: 120rpx;
		height: 120rpx;
		background-color: rgba(0, 0, 0, 0.6);
		border-radius: 16rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
	}
	
	.letter-indicator text {
		font-size: 80rpx;
		color: #ffffff;
		font-weight: bold;
	}
	
	/* 加载提示样式 */
	.loading-tip {
		padding: 30rpx;
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}
	
	.loading-icon {
		width: 40rpx;
		height: 40rpx;
		margin-bottom: 10rpx;
		border: 4rpx solid #f5f5f5;
		border-top: 4rpx solid #fc3e2b;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}
	
	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}
	
	.loading-tip text {
		font-size: 28rpx;
		color: #999;
	}
	
	/* 错误提示样式 */
	.error-tip {
		padding: 40rpx;
		text-align: center;
	}
	
	.error-tip text {
		font-size: 28rpx;
		color: #999;
		display: block;
		margin-bottom: 20rpx;
	}
	
	.retry-btn {
		width: 200rpx;
		height: 70rpx;
		font-size: 28rpx;
		line-height: 70rpx;
		border-radius: 35rpx;
		background-color: #fc3e2b;
		color: #fff;
	}
</style> 