<template>
	<view class="idle-publish-container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">闲置发布</text>
			</view>
		</view>
		
		<!-- 占位元素 -->
		<view class="placeholder" :style="{ height: `calc(${statusBarHeight}px + ${navBarHeight}px)` }"></view>
		
		<!-- 表单内容 -->
		<scroll-view class="form-scroll" scroll-y>
			<!-- 闲置类型选择 -->
			<view class="form-section">
				<view class="section-title">闲置类型</view>
				<view class="type-grid">
					<view 
						v-for="(type, index) in idleTypes" 
						:key="index"
						class="type-item"
						:class="{ active: formData.type === type.value }"
						@tap="selectIdleType(type.value)"
					>
						<image v-if="type.iconUrl" class="type-icon" :src="type.iconUrl" mode="aspectFit"></image>
						<uni-icons v-else type="info" size="24" :color="formData.type === type.value ? '#ffffff' : '#666666'"></uni-icons>
						<text class="type-text">{{ type.label }}</text>
					</view>
				</view>
			</view>
			
			<!-- 标题描述区域 -->
			<view class="form-section">
				<view class="input-item">
					<input 
						class="title-input" 
						v-model="formData.title" 
						placeholder="请输入闲置物品名称" 
						maxlength="30"
					/>
					<text class="count-text">{{formData.title.length}}/30</text>
				</view>
				
				<view class="input-item">
					<textarea 
						class="desc-input" 
						v-model="formData.description" 
						placeholder="描述一下物品的细节和转手原因，更容易出售" 
						maxlength="1000"
						auto-height
					/>
					<text class="count-text">{{formData.description.length}}/1000</text>
				</view>
			</view>
			
			<!-- 图片上传区域 -->
			<view class="form-section">
				<view class="section-title">添加图片</view>
				<view class="upload-area">
					<view class="upload-grid">
						<view 
							v-for="(item, index) in images" 
							:key="index" 
							class="image-item"
						>
							<image class="preview-image" :src="item" mode="aspectFill"></image>
							<view class="delete-btn" @tap.stop="deleteImage(index)">
								<uni-icons type="closeempty" size="14" color="#ffffff"></uni-icons>
							</view>
						</view>
						
						<view class="upload-btn" @tap="chooseImage" v-if="images.length < 9">
							<uni-icons type="camera" size="24" color="#999999"></uni-icons>
							<text class="upload-text">{{images.length}}/9</text>
						</view>
					</view>
					<text class="tip-text">第一张图片将作为主图</text>
				</view>
			</view>
			
			<!-- 价格区域 -->
			<view class="form-section">
				<view class="price-item">
					<text class="label">价格</text>
					<view class="price-input-wrapper">
						<text class="price-symbol">¥</text>
						<input 
							class="price-input" 
							v-model="formData.price" 
							type="digit" 
							placeholder="0.00" 
							@input="validatePriceInput('price')"
						/>
					</view>
				</view>
				
				<view class="price-item">
					<text class="label">原价</text>
					<view class="price-input-wrapper">
						<text class="price-symbol">¥</text>
						<input 
							class="price-input" 
							v-model="formData.originalPrice" 
							type="digit" 
							placeholder="0.00（选填）" 
							@input="validatePriceInput('originalPrice')"
						/>
					</view>
				</view>
			</view>
			
			<!-- 分类区域 -->
			<view class="form-section">
				<view class="select-item" @tap="showLocationPicker">
					<text class="label">交易地点</text>
					<view class="select-value">
						<text class="value-text">{{formData.location || '请选择交易地点'}}</text>
						<uni-icons type="right" size="16" color="#999999"></uni-icons>
					</view>
				</view>
				
				<view class="select-item">
					<text class="label">交易方式</text>
					<view class="tag-group">
						<view class="tag-item active">
							<text class="tag-text">线下交易</text>
						</view>
					</view>
				</view>
			</view>
			
			<!-- 底部安全区域 -->
			<view class="safe-area-bottom"></view>
		</scroll-view>
		
		<!-- 底部操作栏 -->
		<view class="bottom-bar">
			<view class="draft-btn" @tap="saveDraft">
				<text class="draft-text">存草稿</text>
			</view>
			<view class="publish-btn" @tap="publishItem">
				<text class="publish-text">立即发布</text>
			</view>
		</view>
	</view>
</template>

<script>
	import deviceMixin from '@/mixins/device.js'
	import { createIdle, getInfoCategories } from '@/apis/content.js'
	import { API_BASE_URL } from '@/utils/constants.js'
	import { isLoggedIn } from '@/utils/auth.js'
	import { mapState, mapGetters } from 'vuex'
	import { saveIdleDraft, getIdleDraft } from '@/utils/storage.js'
	import { category, publish } from '@/apis/index.js'
	import deviceInfo from '@/utils/device-info.js'
	
	export default {
		mixins: [deviceMixin],
		data() {
			return {
				statusBarHeight: 0,
				navBarHeight: 44,
				images: [],
				loading: false,
				formData: {
					type: '', // 默认为空，等待接口返回后设置
					title: '',
					description: '',
					price: '',
					originalPrice: '',
					location: '同城交易',
					latitude: undefined, // 位置纬度
					longitude: undefined, // 位置经度
					contact: {
						name: '',
						phone: '',
						wechat: ''
					}
				},
				idleTypes: []
			}
		},
		computed: {
			...mapState('region', ['regionList']),
		},
		onLoad(options) {
			// 检查登录状态
			if (!isLoggedIn()) {
				uni.showToast({
					title: '请先登录',
					icon: 'none',
					duration: 2000
				})
				
				// 延迟跳转到登录页面
				setTimeout(() => {
					uni.navigateTo({
						url: '/pages/login/index'
					})
				}, 1500)
				return
			}
			
			// 获取状态栏高度
			this.statusBarHeight = deviceInfo.getStatusBarHeight();
			
			// 确保地区数据已加载
			this.ensureRegionDataLoaded()
			
			// 加载闲置类型数据
			this.loadIdleCategories()
			
			// 如果有草稿ID参数，加载草稿数据
			if (options.draftId) {
				this.loadDraftData(options.draftId)
			}
		},
		methods: {
			handleBack() {
				uni.navigateBack()
			},
			// 加载闲置类型分类数据
			async loadIdleCategories() {
				if (this.loading) return
				this.loading = true
				
				try {
					uni.showLoading({
						title: '加载中...',
						mask: true
					})
					
					// 获取闲置物品类型分类，type=2表示闲置物品分类
					const res = await getInfoCategories(2)
					
					console.log('闲置分类数据:', res)
					
					if (res.code === 0 && res.data && res.data.list) {
						// 转换为页面可用的格式
						this.idleTypes = res.data.list.map(item => {
							return {
								label: item.name,
								value: item.id.toString(),
								iconUrl: item.icon, // 使用API返回的图标URL
								id: item.id
							}
						})
						
						// 设置默认选中第一个分类
						if (this.idleTypes.length > 0) {
							this.formData.type = this.idleTypes[0].value
						}
					} else {
						throw new Error('获取分类失败')
					}
				} catch (error) {
					console.error('加载闲置分类失败:', error)
					// 使用默认分类数据
					this.useDefaultCategories()
					uni.showToast({
						title: '加载分类失败，使用默认分类',
						icon: 'none'
					})
				} finally {
					this.loading = false
					uni.hideLoading() // 确保在所有情况下都隐藏loading
				}
			},
			// 使用默认闲置分类数据（网络错误时的备选方案）
			useDefaultCategories() {
				// 使用空数组而不是默认数据
				this.idleTypes = [];
				
				// 重置类型选择
				this.formData.type = '';
				
				// 显示更明确的错误提示
				uni.showToast({
					title: '分类加载失败，请刷新重试',
					icon: 'none',
					duration: 2000
				});
			},
			async chooseImage() {
				try {
					// 选择图片
					const res = await new Promise((resolve, reject) => {
						uni.chooseImage({
							count: 9 - this.images.length,
							sizeType: ['compressed'],
							sourceType: ['album', 'camera'],
							success: resolve,
							fail: reject
						})
					})
					
					if (res.tempFilePaths && res.tempFilePaths.length > 0) {
						uni.showLoading({
							title: '上传中...',
							mask: true
						})
						
						try {
							// 上传所有选中的图片
							const uploadPromises = res.tempFilePaths.map(async (filePath) => {
								try {
									// 调用上传接口
									return await this.uploadFile(filePath)
								} catch (error) {
									console.error('上传图片失败:', error)
									return null
								}
							})
							
							// 等待所有图片上传完成
							const uploadResults = await Promise.all(uploadPromises)
							
							// 过滤掉上传失败的图片
							const successUrls = uploadResults.filter(url => url !== null)
							
							// 更新图片列表
							this.images = [...this.images, ...successUrls]
						} catch (error) {
							console.error('处理上传图片时出错:', error)
							uni.showToast({
								title: '上传图片失败',
								icon: 'none'
							})
						} finally {
							uni.hideLoading() // 确保在所有情况下都隐藏loading
						}
					}
				} catch (error) {
					console.error('选择图片失败:', error)
					uni.showToast({
						title: '选择图片失败',
						icon: 'none'
					})
					// 如果在chooseImage阶段可能显示了loading，需要确保隐藏
					if (this.loading) {
						uni.hideLoading()
					}
				}
			},
			
			// 上传文件到服务器
			async uploadFile(filePath) {
				return new Promise((resolve, reject) => {
					// 获取token
					const token = uni.getStorageSync('token')
					
					if (!token) {
						uni.showToast({
							title: '请先登录',
							icon: 'none'
						})
						return reject(new Error('请先登录'))
					}
					
					uni.uploadFile({
						url: API_BASE_URL + '/wx/upload/image',
						filePath: filePath,
						name: 'file',
						header: {
							'Authorization': 'Bearer ' + token
						},
						success: (uploadRes) => {
							try {
								const data = JSON.parse(uploadRes.data)
								console.log('上传响应:', data)
								if (data.code === 0 && data.data) {
									// 返回图片的URL
									resolve(data.data.url)
								} else {
									reject(new Error(data.message || '上传失败'))
								}
							} catch (error) {
								console.error('解析上传响应失败:', error)
								reject(error)
							}
						},
						fail: (err) => {
							console.error('上传请求失败:', err)
							reject(err)
						}
					})
				})
			},
			deleteImage(index) {
				this.images.splice(index, 1)
			},
			showLocationPicker() {
				// 调用微信choosePoi接口选择位置
				uni.choosePoi({
					success: (res) => {
						console.log('选择位置成功:', res);
						// 更新位置信息，根据类型决定使用城市还是具体地址
						if (res.type === 1) {
							// 选择城市，直接使用城市名称（不带省份）
							// 如果城市名称中包含"省"字，则只取省后面的部分
							let cityName = res.city || '';
							
							// 去除省份名称
							if (cityName.indexOf('省') > 0) {
								cityName = cityName.split('省')[1];
							} else if (cityName.indexOf('自治区') > 0) {
								cityName = cityName.split('自治区')[1];
							}
							
							this.formData.location = cityName;
							// 城市级别没有精确经纬度，将其设为空
							this.formData.latitude = undefined;
							this.formData.longitude = undefined;
						} else if (res.type === 2) {
							// 选择具体位置，只显示市和区县级别
							let address = res.address || '';
							let name = res.name || '';
							
							// 从地址中移除省份
							let formattedAddress = '';
							
							// 如果地址包含"省"字，则只取省后面的部分
							if (address.indexOf('省') > 0) {
								address = address.split('省')[1];
							} else if (address.indexOf('自治区') > 0) {
								address = address.split('自治区')[1];
							}
							
							// 如果还包含"市"，则取市和后面的部分
							if (address.indexOf('市') > 0) {
								// 直接使用市+区县+地点名的格式
								formattedAddress = address + ' ' + name;
							} else {
								// 如果没有市，直接使用净化后的地址+地点名
								formattedAddress = address + ' ' + name;
							}
							
							this.formData.location = formattedAddress;
							// 保存经纬度信息
							this.formData.latitude = res.latitude;
							this.formData.longitude = res.longitude;
						}
					},
					fail: (err) => {
						console.error('选择位置失败:', err);
						// 位置选择失败时的处理
						if (err.errMsg && err.errMsg.indexOf('cancel') !== -1) {
							// 用户取消选择，不做处理
							return;
						}
						
						// 其他错误显示提示
						uni.showToast({
							title: '位置选择失败，请重试',
							icon: 'none'
						});
					}
				});
			},
			selectIdleType(type) {
				this.formData.type = type
			},
			saveDraft() {
				if (!this.formData.title && !this.formData.description) {
					uni.showToast({
						title: '请填写标题或描述',
						icon: 'none'
					})
					return
				}
				
				try {
					// 构建草稿数据
					const draftData = {
						...this.formData,
						images: this.images
					}
					
					// 保存草稿
					const draftId = saveIdleDraft(draftData)
					
					if (draftId) {
						uni.showToast({
							title: '草稿保存成功',
							icon: 'success'
						})
						
						// 延迟返回上一页
						setTimeout(() => {
							uni.navigateBack()
						}, 1500)
					} else {
						throw new Error('保存失败')
					}
				} catch (e) {
					console.error('保存草稿失败', e)
					uni.showToast({
						title: '保存草稿失败',
						icon: 'none'
					})
				}
			},
			async publishItem() {
				// 检查登录状态
				if (!isLoggedIn()) {
					uni.showToast({
						title: '请先登录',
						icon: 'none'
					})
					return
				}
				
				// 表单验证
				if (!this.validateForm(true)) return
				
				// 使用变量跟踪loading状态
				let loadingShown = false;
				
				try {
					// 显示加载提示
					uni.showLoading({
						title: '正在发布...',
						mask: true
					});
					loadingShown = true;
					
					// 构建富文本内容
					let content = '';
					
					// 添加描述内容，分段落处理
					const textParagraphs = this.formData.description.split('\n').filter(p => p.trim());
					textParagraphs.forEach(paragraph => {
						content += `<p>${paragraph}</p>`;
					});
					
					// 添加图片内容
					if (this.images.length > 0) {
						this.images.forEach(imgUrl => {
							content += `<p><img src="${imgUrl}"></p>`;
						});
					}
					
					// 准备提交的数据
					const submitData = {
						categoryId: this.formData.type ? parseInt(this.formData.type) : 0,
						title: this.formData.title,
						price: parseFloat(this.formData.price) || 0,
						originalPrice: parseFloat(this.formData.originalPrice) || 0,
						content: content,
						tradeMethod: '线下交易',
						tradePlace: this.formData.location || '同城交易',
						images: this.images,
						// 添加地区ID
						regionId: this.getRegionId()
					}
					
					console.log('提交数据:', submitData);
					
					// 调用接口发布闲置信息
					const res = await createIdle(submitData);
					
					// 请求完成后隐藏loading
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					if (res.code === 0) {
						uni.showToast({
							title: '发布成功',
							icon: 'success'
						})
						
						// 延迟返回上一页
						setTimeout(() => {
							uni.navigateBack()
						}, 1500)
					} else {
						throw new Error(res.message || '发布失败')
					}
				} catch (error) {
					// 确保发生异常时也会隐藏loading
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					console.error('发布闲置失败:', error)
					uni.showToast({
						title: error.message || '发布失败，请重试',
						icon: 'none'
					})
				} finally {
					// 最后一道防线，确保无论如何都会隐藏loading
					if (loadingShown) {
						uni.hideLoading();
					}
				}
			},
			validatePriceInput(field) {
				// 价格输入验证：只允许输入数字和小数点
				const value = this.formData[field];
				
				// 如果为空，不处理
				if (!value) return;
				
				// 使用正则表达式验证，只允许数字和一个小数点
				// 允许：数字，一个小数点，小数点后最多两位
				const regex = /^\d+(\.\d{0,2})?$/;
				
				if (!regex.test(value)) {
					// 非法输入检测
					const hasForbiddenChars = /[^\d.]/.test(value);
					
					// 如果不符合正则，则去除非法字符
					const cleaned = value
						.replace(/[^\d.]/g, '') // 去除数字和点以外的字符
						.replace(/\.{2,}/g, '.') // 多个点替换为一个点
						.replace(/^(\d+\.\d{0,2}).*$/, '$1'); // 只保留一个小数点后两位
					
					// 更新值
					this.formData[field] = cleaned;
					
					// 如果存在非法字符，显示提示
					if (hasForbiddenChars) {
						uni.showToast({
							title: '请只输入数字和小数点',
							icon: 'none',
							duration: 1500
						});
					}
				}
			},
			validateForm(isPublish) {
				// 类型验证
				if (!this.formData.type) {
					uni.showToast({
						title: '请选择闲置类型',
						icon: 'none'
					})
					return false
				}
				
				// 图片验证
				if (this.images.length === 0) {
					uni.showToast({
						title: '请至少上传一张图片',
						icon: 'none'
					})
					return false
				}
				
				// 标题验证
				if (!this.formData.title.trim()) {
					uni.showToast({
						title: '请填写商品标题',
						icon: 'none'
					})
					return false
				}
				
				// 发布时的额外验证
				if (isPublish) {
					// 价格验证
					if (!this.formData.price) {
						uni.showToast({
							title: '请填写商品价格',
							icon: 'none'
						})
						return false
					}
					
					// 价格格式验证
					const priceRegex = /^\d+(\.\d{1,2})?$/;
					if (!priceRegex.test(this.formData.price)) {
						uni.showToast({
							title: '请输入正确的价格格式',
							icon: 'none'
						})
						return false
					}
					
					// 原价格式验证（如果有值但不是空字符串）
					if (this.formData.originalPrice !== undefined && this.formData.originalPrice !== '') {
						if (!priceRegex.test(this.formData.originalPrice)) {
							uni.showToast({
								title: '请输入正确的原价格式',
								icon: 'none'
							})
							return false
						}
						
						// 确保原价字段中没有中文或其他非法字符
						const hasForbiddenChars = /[^\d.]/.test(this.formData.originalPrice);
						if (hasForbiddenChars) {
							uni.showToast({
								title: '原价只能输入数字',
								icon: 'none'
							})
							return false
						}
					}
					
					// 交易地点验证
					if (!this.formData.location) {
						uni.showToast({
							title: '请选择交易地点',
							icon: 'none'
						})
						return false
					}
				}
				
				return true
			},
			// 确保区域数据已加载
			ensureRegionDataLoaded() {
				// 检查vuex store中是否已有区域数据
				if (this.$store && 
					this.$store.state.region && 
					(!this.$store.state.region.regionList || this.$store.state.region.regionList.length === 0)) {
					// 如果没有区域数据，则触发加载
					this.$store.dispatch('region/getRegionList').catch(error => {
						console.error('加载地区数据失败:', error)
					})
				}
			},
			getRegionId() {
				// 获取当前选择的地区ID
				let regionId = 0;
				const currentLocation = uni.getStorageSync('currentLocation');
				
				// 从Store中查找当前位置对应的地区ID
				if (this.regionList && this.regionList.length > 0) {
					const currentRegion = this.regionList.find(region => region.name === currentLocation);
					if (currentRegion) {
						regionId = currentRegion.id;
						console.log('找到当前地区ID:', regionId, '地区名称:', currentLocation);
					} else {
						console.log('未找到当前地区ID，使用默认ID 0, 当前地区名称:', currentLocation);
					}
				}
				
				return regionId;
			},
			// 加载草稿数据
			loadDraftData(draftId) {
				try {
					const draftData = getIdleDraft(Number(draftId))
					
					if (!draftData) {
						uni.showToast({
							title: '草稿不存在',
							icon: 'none'
						})
						return
					}
					
					// 填充表单数据
					if (draftData.type) {
						this.formData.type = draftData.type
					}
					
					if (draftData.title) {
						this.formData.title = draftData.title
					}
					
					if (draftData.description) {
						this.formData.description = draftData.description
					}
					
					if (draftData.price) {
						this.formData.price = draftData.price
					}
					
					if (draftData.originalPrice) {
						this.formData.originalPrice = draftData.originalPrice
					}
					
					if (draftData.location) {
						this.formData.location = draftData.location
					}
					
					if (draftData.latitude) {
						this.formData.latitude = draftData.latitude
					}
					
					if (draftData.longitude) {
						this.formData.longitude = draftData.longitude
					}
					
					if (draftData.contact) {
						this.formData.contact = draftData.contact
					}
					
					// 设置图片
					if (draftData.images && draftData.images.length > 0) {
						this.images = draftData.images
					}
				} catch (e) {
					console.error('加载草稿失败', e)
					uni.showToast({
						title: '加载草稿失败',
						icon: 'none'
					})
				}
			}
		}
	}
</script>

<style>
	.idle-publish-container {
		display: flex;
		flex-direction: column;
		height: 100vh;
		background-color: #f5f5f5;
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
	}
	
	.nav-title {
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
		font-size: 34rpx;
		font-weight: 500;
		color: #333333;
	}
	
	.placeholder {
		width: 100%;
	}
	
	.form-scroll {
		flex: 1;
		overflow: hidden;
	}
	
	.form-section {
		background-color: #ffffff;
		margin-bottom: 20rpx;
		padding: 20rpx 30rpx;
	}
	
	.section-title {
		font-size: 30rpx;
		font-weight: 500;
		color: #333333;
		margin-bottom: 20rpx;
	}
	
	.type-grid {
		display: flex;
		flex-wrap: wrap;
		margin: 0 -10rpx;
	}
	
	.type-item {
		width: calc(25% - 20rpx);
		height: 120rpx;
		margin: 10rpx;
		background-color: #f8f8f8;
		border-radius: 8rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}
	
	.type-item.active {
		background: linear-gradient(135deg, #1e90ff, #4169e1);
	}
	
	.type-text {
		font-size: 24rpx;
		color: #666666;
		margin-top: 10rpx;
	}
	
	.type-item.active .type-text {
		color: #ffffff;
	}
	
	.upload-area {
		width: 100%;
	}
	
	.upload-grid {
		display: flex;
		flex-wrap: wrap;
		margin: 0 -10rpx;
	}
	
	.image-item, .upload-btn {
		width: calc(33.33% - 20rpx);
		height: 200rpx;
		margin: 10rpx;
		border-radius: 8rpx;
		overflow: hidden;
		position: relative;
	}
	
	.preview-image {
		width: 100%;
		height: 100%;
	}
	
	.delete-btn {
		position: absolute;
		top: 10rpx;
		right: 10rpx;
		width: 40rpx;
		height: 40rpx;
		background-color: rgba(0, 0, 0, 0.5);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.upload-btn {
		background-color: #f8f8f8;
		border: 1rpx dashed #dddddd;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}
	
	.upload-text {
		font-size: 24rpx;
		color: #999999;
		margin-top: 10rpx;
	}
	
	.tip-text {
		font-size: 24rpx;
		color: #999999;
		margin-top: 20rpx;
	}
	
	.input-item {
		margin-bottom: 20rpx;
		position: relative;
	}
	
	.input-item:last-child {
		margin-bottom: 0;
	}
	
	.title-input {
		width: 100%;
		height: 80rpx;
		font-size: 30rpx;
		color: #333333;
	}
	
	.desc-input {
		width: 100%;
		min-height: 200rpx;
		font-size: 28rpx;
		color: #333333;
		line-height: 1.5;
	}
	
	.count-text {
		position: absolute;
		right: 0;
		bottom: 0;
		font-size: 24rpx;
		color: #999999;
	}
	
	.price-item {
		display: flex;
		align-items: center;
		margin-bottom: 20rpx;
		min-height: 80rpx;
	}
	
	.price-item:last-child {
		margin-bottom: 0;
	}
	
	.label {
		width: 160rpx;
		font-size: 30rpx;
		color: #333333;
	}
	
	.price-input-wrapper {
		flex: 1;
		display: flex;
		align-items: center;
	}
	
	.price-symbol {
		font-size: 30rpx;
		color: #333333;
		margin-right: 10rpx;
	}
	
	.price-input {
		flex: 1;
		height: 80rpx;
		font-size: 30rpx;
		color: #333333;
	}
	
	.select-item {
		display: flex;
		align-items: center;
		margin-bottom: 20rpx;
		min-height: 80rpx;
	}
	
	.select-item:last-child {
		margin-bottom: 0;
	}
	
	.select-value {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	
	.value-text {
		font-size: 30rpx;
		color: #333333;
	}
	
	.tag-group {
		flex: 1;
		display: flex;
	}
	
	.tag-item {
		padding: 10rpx 20rpx;
		background-color: #f5f5f5;
		border-radius: 8rpx;
		margin-right: 20rpx;
	}
	
	.tag-item.active {
		background-color: #fff2f0;
		border: 1rpx solid #fc3e2b;
	}
	
	.tag-text {
		font-size: 26rpx;
		color: #666666;
	}
	
	.tag-item.active .tag-text {
		color: #fc3e2b;
	}
	
	.safe-area-bottom {
		height: 180rpx;
	}
	
	.bottom-bar {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		height: 100rpx;
		background-color: #ffffff;
		display: flex;
		align-items: center;
		padding-bottom: env(safe-area-inset-bottom);
		box-shadow: 0 -2rpx 8rpx rgba(0, 0, 0, 0.04);
		z-index: 99;
	}
	
	.draft-btn, .publish-btn {
		flex: 1;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.draft-text {
		font-size: 30rpx;
		color: #666666;
	}
	
	.publish-btn {
		background: linear-gradient(to right, #ff8c00, #ff4500);
	}
	
	.publish-text {
		font-size: 30rpx;
		color: #ffffff;
		font-weight: 500;
	}
	
	.type-icon {
		width: 48rpx;
		height: 48rpx;
		margin-bottom: 8rpx;
	}
</style> 