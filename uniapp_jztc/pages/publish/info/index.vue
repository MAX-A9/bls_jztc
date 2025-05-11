<template>
	<view class="info-publish-container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">信息发布</text>
			</view>
		</view>
		
		<!-- 占位元素 -->
		<view class="placeholder" :style="{ height: `calc(${statusBarHeight}px + ${navBarHeight}px)` }"></view>
		
		<!-- 表单内容 -->
		<scroll-view class="form-scroll" scroll-y>
			<!-- 信息类型选择 -->
			<view class="form-section">
				<view class="section-title">信息类型</view>
				<view class="type-grid">
					<view 
						v-for="(type, index) in infoTypes" 
						:key="index"
						class="type-item"
						:class="{ active: formData.type === type.id.toString() }"
						@tap="selectInfoType(type.id.toString())"
					>
						<image v-if="type.icon" class="type-icon" :src="type.icon" mode="aspectFit"></image>
						<uni-icons v-else type="info" size="24" :color="formData.type === type.id.toString() ? '#ffffff' : '#666666'"></uni-icons>
						<text class="type-text">{{ type.name }}</text>
					</view>
				</view>
			</view>
			
			<!-- 标题描述区域 -->
			<view class="form-section">
				<view class="input-item">
					<input 
						class="title-input" 
						v-model="formData.title" 
						placeholder="请输入信息标题" 
						maxlength="30"
					/>
					<text class="count-text">{{formData.title.length}}/30</text>
				</view>
				
				<view class="input-item">
					<textarea 
						class="desc-input" 
						v-model="formData.description" 
						placeholder="请详细描述您要发布的信息内容" 
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
				</view>
			</view>
			
			<!-- 发布设置 -->
			<view class="form-section" v-if="publishPackages.length > 0">
				<view class="section-title">选择展示套餐</view>
				
				<view class="option-cards">
					<view 
						class="option-card" 
						:class="{ active: formData.publishPackageId === item.id }"
						v-for="item in publishPackages"
						:key="'publish-'+item.id"
						@tap="selectPublishPackage(item.id)"
					>
						<view class="option-content">
							<text class="option-title">{{ item.title }}</text>
						</view>
						<view class="option-price">{{ item.price === 0 ? '免费' : item.price + '元' }}</view>
						<view class="check-icon" v-if="formData.publishPackageId === item.id"></view>
					</view>
				</view>
			</view>
			
			<view class="form-section" v-if="topPackages.length > 0">
				<view class="section-title">选择置顶套餐</view>
				
				<view class="option-cards">
					<view 
						class="option-card" 
						:class="{ active: formData.topPackageId === item.id }"
						v-for="item in topPackages"
						:key="'top-'+item.id"
						@tap="selectTopPackage(item.id)"
					>
						<view class="option-content">
							<text class="option-title">{{ item.title }}</text>
						</view>
						<view class="option-price">{{ item.price === 0 ? '免费' : item.price + '元' }}</view>
						<view class="check-icon" v-if="formData.topPackageId === item.id"></view>
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
	import { createInfo, getInfoCategories, getPackageList } from '@/apis/content.js'
	import { API_BASE_URL } from '@/utils/constants.js'
	import { isLoggedIn } from '@/utils/auth.js'
	import { mapState, mapGetters } from 'vuex'
	import { requestWxPay } from '@/utils/pay.js'
	import { saveInfoDraft, getInfoDraft } from '@/utils/storage.js'
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
					type: '',
					title: '',
					description: '',
					publishPackageId: null, // 发布套餐ID
					topPackageId: null // 置顶套餐ID，null表示不使用置顶
				},
				infoTypes: [],
				publishPackages: [], // 发布套餐列表
				topPackages: [] // 置顶套餐列表
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
			
			// 加载分类数据
			this.loadCategories()
			
			// 加载套餐数据
			this.loadPackages()
			
			// 确保地区数据已加载
			this.ensureRegionDataLoaded()
			
			// 如果有草稿ID参数，加载草稿数据
			if (options.draftId) {
				this.loadDraftData(options.draftId)
			}
		},
		methods: {
			handleBack() {
				uni.navigateBack()
			},
			// 加载分类数据
			async loadCategories() {
				if (this.loading) return
				this.loading = true
				
				try {
					uni.showLoading({
						title: '加载中...',
						mask: true
					})
					
					const res = await getInfoCategories(1)
					
					console.log('分类数据:', res)
					
					if (res.code === 0 && res.data && res.data.list) {
						this.infoTypes = res.data.list
						
						// 设置默认选中第一个分类
						if (this.infoTypes.length > 0) {
							this.formData.type = this.infoTypes[0].id.toString()
						}
					} else {
						throw new Error('获取分类失败')
					}
				} catch (error) {
					console.error('加载分类失败:', error)
					// 使用默认分类数据
					this.useDefaultCategories()
					uni.showToast({
						title: '加载分类失败，请重试',
						icon: 'none'
					})
				} finally {
					this.loading = false
					uni.hideLoading()
				}
			},
			// 使用默认分类数据
			useDefaultCategories() {
				// 使用空数组而不是默认数据
				this.infoTypes = [];
				
				// 重置类型选择
				this.formData.type = '';
				
				// 显示更明确的错误提示
				uni.showToast({
					title: '分类加载失败，请刷新重试',
					icon: 'none',
					duration: 2000
				});
			},
			handlePreview() {
				uni.showToast({
					title: '预览功能开发中',
					icon: 'none'
				})
			},
			selectInfoType(type) {
				this.formData.type = type
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
							uni.hideLoading()
						}
					}
				} catch (error) {
					console.error('选择图片失败:', error)
					uni.showToast({
						title: '选择图片失败',
						icon: 'none'
					})
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
			selectPublishPackage(packageId) {
				this.formData.publishPackageId = packageId
			},
			selectTopPackage(packageId) {
				this.formData.topPackageId = packageId
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
					const draftId = saveInfoDraft(draftData)
					
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
					
					setTimeout(() => {
						uni.navigateTo({
							url: '/pages/login/index'
						})
					}, 1500)
					return
				}
				
				// 表单验证
				if (!this.validateForm()) return
				
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
					
					// 添加文本内容，分段落处理
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
					
					// 获取当前选择的地区ID
					let regionId = this.getRegionId();
					
					// 准备提交的数据
					const submitData = {
						categoryId: this.formData.type ? parseInt(this.formData.type) : 0,
						title: this.formData.title,
						content: content,
						publishPackageId: this.formData.publishPackageId,
						topPackageId: this.formData.topPackageId,
						isTopRequest: this.formData.topPackageId !== null,
						topDays: 0, // 现在通过topPackageId来处理
						regionId: regionId,
						images: this.images
					}
					
					console.log('提交数据:', submitData);
					
					// 调用接口发布信息
					const res = await createInfo(submitData);
					
					// 请求完成后隐藏loading
					if (loadingShown) {
						uni.hideLoading();
						loadingShown = false;
					}
					
					if (res.code === 0) {
						// 获取订单号
						const orderNo = res.data.orderNo;
						
						// 检查是否需要支付
						if (orderNo) {
							// 获取套餐价格
							const publishPackage = this.publishPackages.find(p => p.id === this.formData.publishPackageId);
							const topPackage = this.formData.topPackageId ? this.topPackages.find(p => p.id === this.formData.topPackageId) : null;
							
							// 计算总价
							const publishPrice = publishPackage ? publishPackage.price : 0;
							const topPrice = topPackage ? topPackage.price : 0;
							const totalPrice = publishPrice + topPrice;
							
							// 只有当总价大于0时才调用支付
							if (totalPrice > 0) {
								try {
									// 支付商品描述
									const body = topPackage 
										? `${publishPackage.title}+${topPackage.title}`
										: publishPackage.title;
									
									// 调用支付接口
									const payResult = await requestWxPay({
										body: body,
										orderNo: orderNo,
										totalFee: totalPrice
									});
									
									if (payResult.success) {
										uni.showToast({
											title: '支付成功',
											icon: 'success'
										});
									} else {
										uni.showToast({
											title: payResult.message || '支付已取消',
											icon: 'none'
										});
									}
								} catch (payError) {
									console.error('支付过程出错:', payError);
									uni.showToast({
										title: payError.message || '支付失败',
										icon: 'none'
									});
								}
							}
						}
						
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
					
					console.error('发布信息失败:', error)
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
			validateForm() {
				// 类型验证
				if (!this.formData.type) {
					uni.showToast({
						title: '请选择信息类型',
						icon: 'none'
					})
					return false
				}
				
				// 标题验证
				if (!this.formData.title.trim()) {
					uni.showToast({
						title: '请填写信息标题',
						icon: 'none'
					})
					return false
				}
				
				// 描述验证
				if (!this.formData.description.trim()) {
					uni.showToast({
						title: '请填写信息描述',
						icon: 'none'
					})
					return false
				}
				
				// 发布套餐验证（只有在有发布套餐时才验证）
				if (this.publishPackages.length > 0 && !this.formData.publishPackageId) {
					uni.showToast({
						title: '请选择发布套餐',
						icon: 'none'
					})
					return false
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
			// 加载套餐数据
			async loadPackages() {
				try {
					uni.showLoading({
						title: '加载中...',
						mask: true
					})
					
					const res = await getPackageList()
					
					if (res.code === 0 && res.data) {
						// 获取套餐数据并按sortOrder排序
						let topPackages = res.data.topPackages || []
						let publishPackages = res.data.publishPackages || []
						
						// 获取套餐总开关状态
						const topEnabled = res.data.topEnabled !== undefined ? res.data.topEnabled : true
						const publishEnabled = res.data.publishEnabled !== undefined ? res.data.publishEnabled : true
						
						// 对置顶套餐按sortOrder排序
						if (topPackages.length > 0) {
							topPackages = topPackages.sort((a, b) => {
								// 如果没有sortOrder字段或相同，则按id排序
								if (a.sortOrder === b.sortOrder || a.sortOrder === undefined || b.sortOrder === undefined) {
									return a.id - b.id
								}
								return a.sortOrder - b.sortOrder
							})
						}
						
						// 对发布套餐按sortOrder排序
						if (publishPackages.length > 0) {
							publishPackages = publishPackages.sort((a, b) => {
								// 如果没有sortOrder字段或相同，则按id排序
								if (a.sortOrder === b.sortOrder || a.sortOrder === undefined || b.sortOrder === undefined) {
									return a.id - b.id
								}
								return a.sortOrder - b.sortOrder
							})
						}
						
						// 根据总开关状态决定是否显示套餐
						this.topPackages = topEnabled ? topPackages : []
						this.publishPackages = publishEnabled ? publishPackages : []
						
						// 如果有发布套餐，默认选中第一个
						if (this.publishPackages.length > 0) {
							this.formData.publishPackageId = this.publishPackages[0].id
						} else {
							// 如果发布套餐被禁用，需要重置 publishPackageId
							this.formData.publishPackageId = null
						}
						
						// 如果置顶套餐被禁用，重置 topPackageId
						if (!topEnabled) {
							this.formData.topPackageId = null
						}
					} else {
						throw new Error('获取套餐数据失败')
					}
				} catch (error) {
					console.error('加载套餐数据失败:', error)
					uni.showToast({
						title: '加载套餐失败，请重试',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			},
			// 加载草稿数据
			loadDraftData(draftId) {
				try {
					const draftData = getInfoDraft(Number(draftId))
					
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
					
					if (draftData.publishPackageId) {
						this.formData.publishPackageId = draftData.publishPackageId
					}
					
					if (draftData.topPackageId) {
						this.formData.topPackageId = draftData.topPackageId
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
	.info-publish-container {
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
	
	.type-icon {
		width: 48rpx;
		height: 48rpx;
	}
	
	.type-text {
		font-size: 24rpx;
		color: #666666;
		margin-top: 10rpx;
	}
	
	.type-item.active .type-text {
		color: #ffffff;
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
	
	.option-cards {
		display: flex;
		flex-direction: column;
		gap: 20rpx;
	}
	
	.option-card {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 30rpx;
		background-color: #f8f8f8;
		border-radius: 16rpx;
		border: 2rpx solid transparent;
		position: relative;
	}
	
	.option-card.active {
		border-color: #1e90ff;
		background-color: #f0f8ff;
	}
	
	.check-icon {
		position: absolute;
		bottom: 0;
		right: 0;
		width: 40rpx;
		height: 40rpx;
		background-color: #1e90ff;
		clip-path: polygon(0 100%, 100% 100%, 100% 0);
	}
	
	.check-icon::after {
		content: '✓';
		position: absolute;
		bottom: 2rpx;
		right: 6rpx;
		color: #ffffff;
		font-size: 20rpx;
		font-weight: bold;
	}
	
	.option-content {
		flex: 1;
	}
	
	.option-title {
		font-size: 28rpx;
		color: #333333;
		line-height: 1.4;
	}
	
	.option-price {
		font-size: 30rpx;
		color: #1e90ff;
		font-weight: 500;
		margin-right: 10rpx;
	}
	
	.hot-icon {
		color: #ff4500;
		font-size: 28rpx;
	}
	
	.ad-price {
		color: #ff8c00;
		font-weight: 500;
		font-size: 28rpx;
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
		background: linear-gradient(to right, #1e90ff, #4169e1);
	}
	
	.publish-text {
		font-size: 30rpx;
		color: #ffffff;
		font-weight: 500;
	}
</style> 