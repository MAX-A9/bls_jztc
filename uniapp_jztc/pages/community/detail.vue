<template>
	<view class="container" :style="pageStyle">
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
			:show-scrollbar="false"
			:enhanced="true"
			:bounces="true"
		>
			<!-- 加载中 -->
			<view class="loading-container" v-if="loading">
				<view class="loading-indicator"></view>
				<text class="loading-text">内容加载中...</text>
			</view>
			
			<!-- 内容主体 -->
			<view class="content-wrapper" v-if="!loading">
				<!-- 商品图片轮播 -->
				<view class="swiper-container">
					<swiper 
						class="swiper" 
						:autoplay="false" 
						:duration="500"
						:circular="true"
						@change="handleSwiperChange"
					>
						<swiper-item v-for="(img, index) in goodsData.images" :key="index">
							<image :src="img" mode="aspectFill" @tap="previewImage(index)"></image>
						</swiper-item>
					</swiper>
					
					<!-- 自定义数字指示器 -->
					<view class="swiper-indicator" v-if="goodsData.images && goodsData.images.length > 0">
						<text>{{currentSwiperIndex + 1}}/{{goodsData.images.length}}</text>
					</view>
				</view>
				
				<!-- 商品信息 -->
				<view class="goods-info">
					<view class="price-row">
						<view class="price-left">
							<text class="price">¥{{goodsData.price}}</text>
							<text class="original-price" v-if="goodsData.originalPrice">¥{{goodsData.originalPrice}}</text>
						</view>
						<view class="price-right">
							<text class="stat-item">{{goodsData.likes || 0}}人想要</text>
							<text class="stat-separator">|</text>
							<text class="stat-item">{{goodsData.viewCount}}浏览</text>
						</view>
					</view>
					<view class="title">{{goodsData.title}}</view>
				</view>
				
				<!-- 商品描述 -->
				<view class="goods-detail">
					<view class="section-title">商品描述</view>
					<rich-text :nodes="processedDescription" class="detail-content"></rich-text>
				</view>
				
				<!-- 交易地点 -->	
				<view class="goods-detail">
					<view class="section-title">交易地点</view>
					<view class="location-content">
						<text class="detail-content location-text">{{goodsData.detailLocation || defaultLocation}}</text>
						<view class="copy-btn" @tap="copyLocation">
							<text class="copy-text">复制</text>
						</view>
					</view>
				</view>

				<!-- 交易方式 -->
				<view class="goods-detail">
					<view class="section-title">交易方式</view>
					<view class="tag-item active">
						<text class="tag-text">{{goodsData.tradeMethod}}</text>
					</view>
				</view>

				<!-- 发布人信息 -->
				<view class="publisher-card">
					<view class="publisher-info">
						<view class="publisher-left">
							<image class="avatar" :src="goodsData.publisher_avatar" mode="aspectFill"></image>
							<view class="info-content">
								<text class="name">{{goodsData.publisher_name}}</text>
								<text class="publish-count">已发布 {{goodsData.publish_count}} 件</text>
							</view>
						</view>
						<view class="follow-btn" @tap="handleFollow">
							<text>{{isFollowed ? '已关注' : '关注'}}</text>
						</view>
					</view>
				</view>
				
				<!-- Banner -->
				<view class="banner-container" v-if="bannerList.length > 0">
					<swiper 
						class="banner-swiper" 
						:circular="true" 
						:autoplay="true" 
						:interval="3000" 
						:duration="500"
						@change="handleBannerSwiperChange"
					>
						<swiper-item v-for="(item, index) in bannerList" :key="index" @tap="handleBannerClick(item)">
							<image :src="item.image" mode="aspectFill" class="banner-image"></image>
						</swiper-item>
					</swiper>
					
					<!-- 指示点 -->
					<view class="dots" v-if="bannerList.length > 1">
						<view 
							class="dot" 
							v-for="(item, index) in bannerList" 
							:key="index"
							:class="{ active: currentBannerIndex === index }"
						></view>
					</view>
				</view>
				
				<!-- 留言区域 -->
				<view class="comment-card">
					<view class="comment-header">
						<text class="comment-title">留言 {{goodsData.commentCount || 0}}</text>
						<view class="comment-btn" @tap="handleComment">
							<uni-icons type="chat" size="14" color="#007AFF"></uni-icons>
							<text>我要留言</text>
						</view>
					</view>
					
					<!-- 留言列表 -->
					<view class="comment-list" v-if="goodsData.comments && goodsData.comments.length">
						<view class="comment-item" v-for="(comment, index) in goodsData.comments" :key="index">
							<image class="comment-avatar" :src="comment.avatar" mode="aspectFill"></image>
							<view class="comment-content">
								<text class="comment-name">{{comment.name}}</text>
								<text class="comment-text">{{comment.content}}</text>
								<text class="comment-time">{{comment.time}}</text>
							</view>
						</view>
					</view>
					
					<!-- 无留言提示 -->
					<view class="no-comment" v-else>
						<text>暂无留言</text>
					</view>
				</view>
			</view>
		</scroll-view>
		
		<!-- 底部操作栏 -->
		<action-bar 
			:is-collected="isCollected"
			:publisher="{id: goodsData.publisher_id, name: goodsData.publisher_name}"
			@comment="handleComment"
			@collect="handleCollect"
			@share="handleShare"
			@message="handleMessage"
		/>
	</view>
</template>

<script>
	import deviceAdapter from '@/mixins/device-adapter.js'
	import shareMixin from '@/mixins/share.js'
	import { getPublicContentDetail, getCommentList, addFavorite, cancelFavorite, getFavoriteStatus } from '@/apis/content.js'
	import { createComment, getPublisherInfo, followUser, unfollowUser, getPublisherFollowStatus } from '@/apis/content.js'
	import { formatTimeAgo } from '@/utils/date.js'
	import ActionBar from '@/components/action-bar/index.vue'
	import { get } from '@/utils/request.js'
	import { getShareSettings } from '@/utils/share.js'
	
	export default {
		mixins: [deviceAdapter, shareMixin],
		components: {
			ActionBar
		},
		onLoad(options) {
			// 获取传递的商品ID
			const id = options.id
			console.log('商品ID:', id)
			// 根据ID加载商品数据
			if (id) {
				this.goodsId = id
				this.loadGoodsData(id)
				// 加载轮播图数据
				this.loadBannerData()
			} else {
				uni.showToast({
					title: '缺少商品ID',
					icon: 'none'
				})
			}
		},
		data() {
			return {
				goodsId: null,
				goodsData: {
					id: '',
					title: '',
					price: '0',
					originalPrice: '0',
					condition: '',
					location: '',
					detailLocation: '',
					tradeMethod: '线下交易',
					description: '',
					images: [],
					publisher_avatar: '',
					publisher_name: '',
					publisher_id: '',
					publish_count: 0,
					likes: 0,
					viewCount: 0,
					comments: [],
					cover: '' // 用于分享图片
				},
				loading: false,
				loadError: false,
				errorMsg: '',
				isFollowed: false,
				isCollected: false,
				currentSwiperIndex: 0,
				commentLoading: false,
				defaultLocation: '江西南昌-青山湖区-高新大道1888号',
				bannerList: [],
				currentBannerIndex: 0,
				isContentPage: true, // 标记为内容页面，用于分享功能
				shareData: null
			}
		},
		computed: {
			pageTitle() {
				return '闲置详情'
			},
			pageStyle() {
				return {
					'--nav-height': `${this.layoutSize.navHeight}px`,
					'--content-padding': `${this.layoutSize.contentPadding}rpx`
				}
			},
			// 处理HTML内容，添加类名以适配小程序
			processedDescription() {
				if (!this.goodsData.description) return '';
				
				// 为段落添加类名
				let content = this.goodsData.description
					.replace(/<p/g, '<p class="rich-paragraph"')
					.replace(/<span/g, '<span class="rich-text"')
					.replace(/<div/g, '<div class="rich-paragraph"');
					
				return content;
			}
		},
		methods: {
			handleBack() {
				uni.navigateBack()
			},
			previewImage(index) {
				uni.previewImage({
					current: index,
					urls: this.goodsData.images
				})
			},
			async handleFollow() {
				try {
					// 安全检查：确保publisher_id非空
					const publisherId = this.goodsData.publisher_id;
					if (!publisherId) {
						uni.showToast({
							title: '无法获取发布者信息',
							icon: 'none'
						});
						return;
					}
					
					// 显示加载提示
					uni.showLoading({
						title: this.isFollowed ? '取消关注中...' : '关注中...',
						mask: true
					});
					
					// 根据当前关注状态调用不同的API
					let result;
					
					if (this.isFollowed) {
						// 已关注，执行取消关注
						result = await unfollowUser({
							publisher_id: publisherId
						});
					} else {
						// 未关注，执行关注
						result = await followUser({
							publisher_id: publisherId
						});
					}
					
					// 关闭加载提示
					uni.hideLoading();
					
					if (result && result.code === 0) {
						// 操作成功后重新加载关注状态，保证与服务器同步
						await this.loadFollowStatus(String(publisherId));
						
						// 显示操作成功提示
						uni.showToast({
							title: this.isFollowed ? '已关注' : '已取消关注',
							icon: 'success'
						});
					} else {
						throw new Error(result?.message || '操作失败');
					}
				} catch (error) {
					// 确保加载提示已关闭
					uni.hideLoading();
					
					console.error('关注操作失败:', error);
					uni.showToast({
						title: error.message || '操作失败，请重试',
						icon: 'none'
					});
				}
			},
			async handleCollect() {
				try {
					// 显示加载提示
					uni.showLoading({
						title: this.isCollected ? '取消收藏中...' : '收藏中...',
						mask: true
					})
					
					// 调用接口
					let result
					if (this.isCollected) {
						// 已收藏，执行取消收藏
						result = await cancelFavorite(this.goodsId)
					} else {
						// 未收藏，执行添加收藏
						result = await addFavorite(this.goodsId)
					}
					
					// 关闭加载提示
					uni.hideLoading()
					
					if (result.code === 0) {
						// 操作成功，更新状态
						this.isCollected = !this.isCollected
						uni.showToast({
							title: this.isCollected ? '已收藏' : '已取消收藏',
							icon: 'success'
						})
					} else {
						throw new Error(result.message || '操作失败')
					}
				} catch (error) {
					uni.hideLoading()
					console.error('收藏操作失败:', error)
					uni.showToast({
						title: error.message || '操作失败，请重试',
						icon: 'none'
					})
				}
			},
			handleMessage() {
				// 获取发布者ID
				const publisherId = this.goodsData.publisher_id || '';
				if (!publisherId) {
					uni.showToast({
						title: '无法获取发布者信息',
						icon: 'none'
					});
					return;
				}
				
				// 跳转到聊天页面
				uni.navigateTo({
					url: `/pages/chat/detail?userId=${publisherId}&userName=${this.goodsData.publisher_name}`,
					success: () => {
						console.log('跳转到私信页面成功')
					},
					fail: (err) => {
						console.error('跳转到私信页面失败:', err)
						uni.showToast({
							title: '无法打开私信',
							icon: 'none'
						})
					}
				})
			},
			handleShare() {
				// 更新分享数据
				this.updateShareData();
				// 显示分享菜单，包括分享朋友圈
				uni.showShareMenu({
					withShareTicket: true,
					menus: ['shareAppMessage', 'shareTimeline']
				});
			},
			// 过滤HTML内容中的图片标签
			filterImgTags(htmlContent) {
				if (!htmlContent) return '';
				// 只移除img标签，保留其他HTML标签
				return htmlContent.replace(/<img[^>]*>/g, '');
			},
			
			// 格式化发布时间为"多久前"
			formatPublishTime(time) {
				if (!time) return '';
				return formatTimeAgo(time);
			},
			
			// 使用API加载商品数据
			async loadGoodsData(id) {
				this.loading = true
				
				try {
					// 调用接口获取商品详情
					const result = await getPublicContentDetail(id)
					
					if (result.code === 0 && result.data) {
						const data = result.data
						
						// 过滤content中的img标签（因为images数组中已经有图片）
						const filteredContent = this.filterImgTags(data.content)
						// 格式化发布时间
						const formattedTime = this.formatPublishTime(data.publishTime)
						
						// 更新商品数据，保持原有结构
						this.goodsData = {
							id: data.id,
							title: data.title || '',
							price: data.price || 0,
							originalPrice: data.originalPrice || 0,
							description: filteredContent || '',
							images: data.images || [],
							condition: data.condition || '',
							location: data.location || '',
							detailLocation: data.tradePlace || this.defaultLocation,
							tradeMethod: data.tradeMethod || '线下交易',
							publisher_avatar: data.publisherAvatar || '',
							publisher_name: data.publisher || '匿名用户',
							publisher_id: data.publisherId || data.publisher_id || '',
							publish_count: data.publishCount || 0,
							likes: data.likes || 0,
							viewCount: data.views || 0,
							comments: [],
							cover: data.cover || '' // 更新cover字段
						}
						
						// 加载评论列表
						this.loadComments()
						
						// 加载收藏状态
						this.loadFavoriteStatus()
						
						// 加载发布者信息（确保有publisher_id才调用）
						if (this.goodsData.publisher_id) {
							this.loadPublisherInfo()
						}
						
						// 更新分享数据
						this.updateShareData();
					} else {
						throw new Error(result.message || '获取商品详情失败')
					}
				} catch (error) {
					console.error('加载商品详情失败:', error)
					uni.showToast({
						title: '加载商品详情失败',
						icon: 'none'
					})
				} finally {
					this.loading = false
				}
			},
			
			// 加载收藏状态
			async loadFavoriteStatus() {
				if (!this.goodsId) return
				
				try {
					const result = await getFavoriteStatus(this.goodsId)
					if (result.code === 0 && result.data) {
						this.isCollected = result.data.isFavorite
					}
				} catch (error) {
					console.error('获取收藏状态失败:', error)
				}
			},
			
			// 加载评论列表
			async loadComments() {
				if (!this.goodsData.id) return
				
				this.commentLoading = true
				
				try {
					const result = await getCommentList(this.goodsData.id, {
						page: 1,
						pageSize: 10
					})
					
					if (result.code === 0 && result.data) {
						// 格式化评论数据，适配实际API返回的字段名
						this.goodsData.comments = (result.data.list || []).map(item => ({
							id: item.id,
							avatar: item.avatarUrl || '',
							name: item.realName || '匿名用户',
							content: item.comment || '',
							time: this.formatPublishTime(item.createdAt)
						}))
						
						// 更新评论总数
						if (result.data.total !== undefined) {
							this.goodsData.commentCount = result.data.total
						}
					} else {
						console.error('获取评论列表失败:', result.message)
					}
				} catch (error) {
					console.error('加载评论失败:', error)
				} finally {
					this.commentLoading = false
				}
			},
			
			handleComment() {
				uni.showModal({
					title: '留言',
					editable: true,
					placeholderText: '请输入留言内容',
					success: async (res) => {
						if (res.confirm && res.content) {
							// 先定义一个变量标记是否已显示loading
							let loadingShown = false;
							
							try {
								// 调用创建评论接口
								const commentData = {
									comment: res.content,
									contentId: this.goodsData.id
								};
								
								// 显示加载提示
								uni.showLoading({
									title: '提交中...'
								});
								loadingShown = true;
								
								const result = await createComment(commentData);
								
								// 操作完成，隐藏加载提示
								uni.hideLoading();
								loadingShown = false;
								
								if (result.code === 0) {
									// 评论成功
									uni.showToast({
										title: '留言成功',
										icon: 'success'
									});
									
									// 刷新评论列表
									this.loadComments();
								} else {
									throw new Error(result.message || '评论失败');
								}
							} catch (error) {
								// 确保发生异常时也会隐藏加载提示
								if (loadingShown) {
									uni.hideLoading();
									loadingShown = false;
								}
								
								console.error('提交评论失败:', error);
								uni.showToast({
									title: error.message || '留言失败，请重试',
									icon: 'none'
								});
							} finally {
								// 最后一道防线，确保无论如何都会隐藏加载提示
								if (loadingShown) {
									uni.hideLoading();
								}
							}
						}
					}
				});
			},
			handleSwiperChange(e) {
				this.currentSwiperIndex = e.detail.current
			},
			handleBannerClick(item) {
				// 处理Banner点击事件
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
			},
			handleBannerSwiperChange(e) {
				this.currentBannerIndex = e.detail.current
			},
			copyLocation() {
				uni.setClipboardData({
					data: this.goodsData.detailLocation || this.defaultLocation,
					success: () => {
						uni.showToast({
							title: '地址已复制',
							icon: 'success'
						})
					}
				})
			},
			// 加载发布者信息
			async loadPublisherInfo() {
				try {
					// 安全检查：确保publisher_id存在且有效
					const publisherId = this.goodsData.publisher_id;
					if (!publisherId) {
						console.error('无法加载发布者信息: publisher_id不存在');
						return;
					}
					
					// 将publisherId转换为字符串
					const publisherIdStr = String(publisherId);
					
					// 调用API获取发布者信息
					const result = await getPublisherInfo(publisherIdStr);
					
					// 检查API返回结果
					if (result && result.code === 0 && result.data) {
						// 获取数据
						const publisherData = result.data;
						
						// 安全地更新UI数据，使用API返回的字段
						if (publisherData.real_name) {
							this.goodsData.publisher_name = publisherData.real_name;
						}
						
						if (publisherData.avatar_url) {
							this.goodsData.publisher_avatar = publisherData.avatar_url;
						}
						
						if (publisherData.publish_count !== undefined) {
							this.goodsData.publish_count = publisherData.publish_count;
						}
						
						// 加载关注状态
						await this.loadFollowStatus(publisherIdStr);
					} else {
						console.error('获取发布者信息失败:', result?.message || '未知错误');
					}
				} catch (error) {
					console.error('加载发布者信息异常:', error);
				}
			},
			// 加载关注状态
			async loadFollowStatus(publisherId) {
				try {
					if (!publisherId) {
						console.error('无法获取关注状态: publisher_id不存在');
						return;
					}
					
					const result = await getPublisherFollowStatus(publisherId);
					
					if (result && result.code === 0 && result.data) {
						// 更新关注状态
						this.isFollowed = !!result.data.is_followed;
					} else {
						console.error('获取关注状态失败:', result?.message || '未知错误');
					}
				} catch (error) {
					console.error('加载关注状态异常:', error);
				}
			},
			// 加载轮播图数据
			async loadBannerData() {
				try {
					const result = await get('/wx/inner-banner/list', {
						bannerType: 'idle'
					})
					
					if (result && result.code === 0 && result.data) {
						// 根据返回结果，更新轮播图数据
						if (result.data.isGlobalEnabled && result.data.list && result.data.list.length > 0) {
							this.bannerList = result.data.list
								.filter(item => item.isEnabled)
								.sort((a, b) => a.order - b.order)
						}
					}
				} catch (error) {
					console.error('获取轮播图数据失败', error)
				}
			},
			// 更新分享数据
			async updateShareData() {
				try {
					// 获取分享配置
					const settings = await getShareSettings();
					
					if (settings && this.goodsData) {
						// 分享标题处理
						// 1. 如果有商品标题，优先使用商品标题
						// 2. 否则使用后端配置的content_share_text
						const shareTitle = this.goodsData.title || settings.content_share_text || '查看闲置商品';
						
						// 分享图片处理：
						// 1. 优先使用第一张商品图片作为封面
						// 2. 如果没有商品图片，使用后端配置的content_share_image
						let shareImage = '';
						if (this.goodsData.images && this.goodsData.images.length > 0) {
							shareImage = this.goodsData.images[0];
						} else if (settings.content_share_image) {
							shareImage = settings.content_share_image;
						}
						
						// 更新商品的封面图
						this.goodsData.cover = shareImage;
						
						// 直接设置分享数据
						this.shareData = {
							title: shareTitle,
							imageUrl: shareImage,
							path: `/pages/community/detail?id=${this.goodsId}`
						};
						
						console.log('闲置商品分享数据已更新:', this.shareData);
					} else {
						// 配置获取失败，调用混入的默认初始化
						this.initShareData();
					}
				} catch (error) {
					console.error('更新分享数据失败:', error);
					// 配置获取失败，调用混入的默认初始化
					this.initShareData();
				}
			},
			// 添加onShow生命周期
			onShow() {
				// 确保分享功能启用，使用最新配置
				this.updateShareData();
			}
		}
	}
</script>

<style>
	.container {
		width: 100%;
		min-height: 100vh;
		position: relative;
		background-color: #FFFFFF;
	}
	
	/* 导航栏样式 */
	.nav-bar {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 100;
		background-color: #FFFFFF;
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
		font-size: 32rpx;
		color: #262626;
		font-weight: 600;
	}
	
	/* 轮播图容器样式 */
	.swiper-container {
		position: relative;
		width: 100%;
		height: 750rpx;
		background-color: #FFFFFF;
	}
	
	/* 轮播图样式 */
	.swiper {
		width: 100%;
		height: 750rpx;
		background-color: #FFFFFF;
	}
	
	.swiper image {
		width: 100%;
		height: 100%;
	}
	
	/* 自定义数字指示器样式 */
	.swiper-indicator {
		position: absolute;
		right: 2rpx;
		bottom: 2rpx;
		background-color: rgba(0, 0, 0, 0.3);
		padding: 2rpx 12rpx;
		height: 32rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 100rpx;
		z-index: 99;
		pointer-events: none;
	}
	
	.swiper-indicator text {
		color: #FFFFFF;
		font-size: 24rpx;
		line-height: 1;
	}
	
	/* 商品信息样式 */
	.goods-info {
		background-color: #FFFFFF;
		padding: 30rpx;
		padding-left: var(--content-padding);
		padding-right: var(--content-padding);
	}
	
	.price-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16rpx;
		margin-bottom: 20rpx;
	}
	
	.price-left {
		display: flex;
		align-items: baseline;
		gap: 16rpx;
	}
	
	.price {
		font-size: 48rpx;
		color: #fc3e2b;
		font-weight: 600;
	}
	
	.original-price {
		font-size: 28rpx;
		color: #999999;
		text-decoration: line-through;
	}
	
	.price-right {
		display: flex;
		align-items: center;
		gap: 12rpx;
	}
	
	.stat-item {
		font-size: 24rpx;
		color: #999999;
	}
	
	.stat-separator {
		font-size: 24rpx;
		color: #EEEEEE;
	}
	
	.title {
		font-size: 32rpx;
		color: #262626;
		font-weight: 500;
		margin-bottom: 16rpx;
		line-height: 1.4;
	}
	
	/* 商品描述样式 */
	.goods-detail {
		background-color: #FFFFFF;
		padding: 30rpx;
		padding-left: var(--content-padding);
		padding-right: var(--content-padding);
	}
	
	.section-title {
		font-size: 30rpx;
		color: #262626;
		font-weight: 600;
		margin-bottom: 20rpx;
	}
	
	/* 确保rich-text中段落的样式 */
	.detail-content {
		font-size: 28rpx;
		color: #333333;
		line-height: 1.6;
	}
	
	/* 小程序不支持标签选择器，使用类选择器替代 */
	.rich-paragraph {
		margin-bottom: 10rpx;
	}
	
	.rich-text {
		font-size: 28rpx;
		color: #333333;
		line-height: 1.6;
	}
	
	/* 发布人信息样式 */
	.publisher-card {
		background-color: #FFFFFF;
		padding: 30rpx;
		padding-left: var(--content-padding);
		padding-right: var(--content-padding);
	}
	
	.publisher-info {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.publisher-left {
		display: flex;
		align-items: center;
		gap: 20rpx;
	}
	
	.avatar {
		width: 80rpx;
		height: 80rpx;
		border-radius: 50%;
		background-color: #f5f5f5;
	}
	
	.info-content {
		display: flex;
		flex-direction: column;
		gap: 8rpx;
	}
	
	.name {
		font-size: 28rpx;
		color: #262626;
		font-weight: 500;
	}
	
	.publish-count {
		font-size: 24rpx;
		color: #999999;
	}
	
	.follow-btn {
		padding: 12rpx 32rpx;
		border-radius: 32rpx;
		background-color: #fc3e2b;
	}
	
	.follow-btn text {
		font-size: 26rpx;
		color: #FFFFFF;
	}
	
	/* 内容区域样式 */
	.content-scroll {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 100rpx;
		z-index: 1;
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
	}
	
	.content-wrapper {
		position: relative;
		padding-bottom: calc(env(safe-area-inset-bottom) + 40rpx);
	}
	
	/* Banner样式 */
	.banner-container {
		width: 100%;
		height: 200rpx;
		display: flex;
		align-items: center;
		position: relative;
	}
	
	.banner-swiper {
		width: 100%;
		height: 100%;
	}
	
	.banner-image {
		width: 100%;
		height: calc(100% - 10rpx);
	}
	
	.dots {
		position: absolute;
		bottom: 20rpx;
		right: 20rpx;
		display: flex;
		justify-content: flex-end;
		gap: 8rpx;
	}
	
	.dot {
		width: 12rpx;
		height: 12rpx;
		border-radius: 50%;
		background-color: rgba(255, 255, 255, 0.6);
	}
	
	.dot.active {
		background-color: #fc3e2b;
		width: 24rpx;
		border-radius: 6rpx;
	}
	
	.location-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	
	.location-text {
		flex: 1;
		padding-right: 20rpx;
	}
	
	.copy-btn {
		padding: 6rpx 16rpx;
		background-color: #f5f5f5;
		border-radius: 30rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.copy-text {
		font-size: 24rpx;
		color: #666666;
	}
	
	/* 交易方式样式 */
	.trade-method-content {
		padding: 16rpx 24rpx;
		border: 1px solid #e5e5e5;
		border-radius: 8rpx;
		background-color: #ffffff;
	}
	
	.tag-item {
		display: inline-block;
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
	
	/* 加载状态样式 */
	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 300rpx;
		margin-top: 100rpx;
	}
	
	.loading-indicator {
		width: 40rpx;
		height: 40rpx;
		border: 4rpx solid #f0f0f0;
		border-radius: 50%;
		border-top-color: #fc3e2b;
		animation: spin 1s linear infinite;
		margin-bottom: 20rpx;
	}
	
	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}
	
	.loading-text {
		font-size: 26rpx;
		color: #999999;
	}
	
	/* 留言区域样式 */
	.comment-card {
		background-color: #FFFFFF;
		padding: 30rpx;
		padding-left: var(--content-padding);
		padding-right: var(--content-padding);
		margin-bottom: 40rpx;
	}
	
	.comment-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 30rpx;
	}
	
	.comment-title {
		font-size: 30rpx;
		color: #262626;
		font-weight: 600;
	}
	
	.comment-btn {
		display: flex;
		align-items: center;
		gap: 8rpx;
	}
	
	.comment-btn text {
		font-size: 26rpx;
		color: #007AFF;
	}
	
	.comment-list {
		display: flex;
		flex-direction: column;
		gap: 30rpx;
	}
	
	.comment-item {
		display: flex;
		gap: 20rpx;
	}
	
	.comment-avatar {
		width: 64rpx;
		height: 64rpx;
		border-radius: 50%;
		background-color: #f5f5f5;
	}
	
	.comment-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 8rpx;
	}
	
	.comment-name {
		font-size: 26rpx;
		color: #666666;
	}
	
	.comment-text {
		font-size: 28rpx;
		color: #262626;
		line-height: 1.5;
	}
	
	.comment-time {
		font-size: 24rpx;
		color: #999999;
	}
	
	.no-comment {
		text-align: center;
		padding: 40rpx 0;
	}
	
	.no-comment text {
		font-size: 26rpx;
		color: #999999;
	}
</style> 