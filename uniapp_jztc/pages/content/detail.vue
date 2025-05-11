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
		>
			<!-- 加载中 -->
			<view class="loading-container" v-if="loading">
				<view class="loading-indicator"></view>
				<text class="loading-text">内容加载中...</text>
			</view>
			
			<!-- 加载失败 -->
			<view class="error-container" v-if="loadError">
				<text class="error-text">{{errorMsg}}</text>
				<view class="retry-btn" @tap="loadDetail">重新加载</view>
			</view>
			
			<!-- 内容展示 -->
			<view class="content-wrapper" v-if="!loading && !loadError">
				<!-- 标签区域 -->
				<view class="tag-row">
					<view class="tag top-tag" v-if="contentData.isTop">置顶</view>
					<view class="tag category-tag">{{contentData.category}}</view>
				</view>
				
				<!-- 标题 -->
				<view class="content-title">{{contentData.title}}</view>
				
				<!-- 发布信息 -->
				<view class="publish-info">
					<text class="publisher">{{contentData.publisher}}</text>
					<text class="publish-time">{{contentData.publishTime}}</text>
				</view>
				
				<!-- 内容区域 -->
				<view class="content-body">
					<!-- 文字内容 -->
					<view class="content-text" v-if="contentData.content">
						<rich-text :nodes="contentData.content"></rich-text>
					</view>
					
					<!-- 图片内容 -->
					<view class="image-list" v-if="contentData.images && contentData.images.length">
						<image 
							v-for="(img, index) in contentData.images" 
							:key="index"
							:src="img"
							mode="widthFix"
							@tap="previewImage(index)"
						></image>
					</view>
				</view>
				
				<!-- 发布人信息 -->
				<view class="publisher-card">
					<view class="publisher-info">
						<view class="publisher-left">
							<image class="avatar" :src="contentData.publisher_avatar" mode="aspectFill"></image>
							<view class="info-content">
								<text class="name">{{contentData.publisher_name}}</text>
								<text class="publish-count">已发布 {{contentData.publish_count}} 条</text>
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
						@change="handleSwiperChange"
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
				<view class="comment-section">
					<view class="comment-header">
						<text class="comment-title">留言 {{contentData.comment_count}}</text>
						<view class="comment-btn" @tap="handleComment">
							<uni-icons type="chat" size="14" color="#007AFF"></uni-icons>
							<text>我要留言</text>
						</view>
					</view>
					<!-- 留言列表 -->
					<view class="comment-list" v-if="contentData.comments && contentData.comments.length">
						<view class="comment-item" v-for="(comment, index) in contentData.comments" :key="index">
							<image class="comment-avatar" :src="comment.avatar" mode="aspectFill"></image>
							<view class="comment-content">
								<text class="comment-name">{{comment.name}}</text>
								<text class="comment-text">{{comment.content}}</text>
								<text class="comment-time">{{comment.time}}</text>
							</view>
						</view>
					</view>
				</view>
			</view>
		</scroll-view>
		
		<!-- 底部固定操作栏 -->
		<action-bar 
			:is-collected="isCollected"
			:publisher="{id: contentData.publisher_id, name: contentData.publisher_name}"
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
	import { getPublicContentDetail, getCommentList, addFavorite, cancelFavorite, getFavoriteStatus, addBrowseRecord, getPublisherInfo, followUser, unfollowUser, getPublisherFollowStatus } from '@/apis/content.js'
	import { createComment } from '@/apis/content.js'
	import { formatTimeAgo } from '@/utils/date.js'
	import ActionBar from '@/components/action-bar/index.vue'
	import { get } from '@/utils/request.js'
	import { createConversation } from '@/apis/message.js'
	import { getShareSettings } from '@/utils/share.js'
	
	export default {
		mixins: [deviceAdapter, shareMixin],
		components: {
			ActionBar
		},
		data() {
			return {
				contentId: null,
				contentData: {
					id: '',
					title: '',
					content: '',
					category: '',
					publisher: '',
					publishTime: '',
					publisher_name: '',
					publisher_id: '',
					publisher_avatar: '',
					images: [],
					isTop: false,
					publish_count: 0,
					comment_count: 0,
					comments: [],
					cover: '', // 添加封面图属性，用于分享
				},
				loading: false,
				loadError: false,
				errorMsg: '',
				isFollowed: false,
				isCollected: false,
				commentLoading: false,
				bannerList: [],
				currentBannerIndex: 0,
				isContentPage: true, // 标记为内容页面，用于分享功能
			}
		},
		computed: {
			pageTitle() {
				return this.contentData?.title || '内容详情'
			},
			pageStyle() {
				return {
					'--nav-height': `${this.layoutSize.navHeight}px`,
					'--content-padding': `${this.layoutSize.contentPadding}rpx`
				}
			}
		},
		onLoad(options) {
			// 获取内容ID
			if (options.id) {
				this.contentId = options.id
				// 加载内容详情
				this.loadDetail()
				// 加载轮播图数据
				this.loadBannerData()
			} else {
				this.loadError = true
				this.errorMsg = '参数错误，缺少内容ID'
			}
		},
		onShow() {
			// 确保分享功能启用，使用最新配置
			this.updateShareData();
		},
		methods: {
			// 过滤HTML内容中的图片标签
			filterImgTags(htmlContent) {
				if (!htmlContent) return '';
				// 使用正则表达式移除img标签
				return htmlContent.replace(/<img[^>]*>/g, '');
			},
			
			// 格式化发布时间为"多久前"
			formatPublishTime(time) {
				if (!time) return '';
				return formatTimeAgo(time);
			},
			
			// 加载内容详情
			async loadDetail() {
				if (!this.contentId) {
					this.loadError = true
					this.errorMsg = '内容ID无效'
					return
				}
				
				this.loading = true
				
				try {
					// 调用接口获取详情
					const result = await getPublicContentDetail(this.contentId)
					
					if (result.code === 0 && result.data) {
						// 获取数据成功
						const data = result.data
						
						// 过滤content中的img标签
						const filteredContent = this.filterImgTags(data.content)
						
						// 格式化发布时间
						const formattedTime = this.formatPublishTime(data.publishTime)
						
						// 更新内容数据
						this.contentData = {
							id: data.id,
							title: data.title || '',
							content: filteredContent,
							category: data.category || '',
							publishTime: formattedTime,
							publisher: data.publisher || '',
							publisher_name: data.publisher || '',
							publisher_id: data.publisher_id || '',
							publisher_avatar: '',
							isTop: data.isTop || false,
							images: data.images || [],
							publish_count: 0,
							comment_count: 0,
							comments: [],
							cover: data.cover || (data.images && data.images.length > 0 ? data.images[0] : ''), // 优先使用封面图，其次使用第一张图片
						}
						
						// 加载评论列表
						this.loadComments()
						
						// 加载收藏状态
						this.loadFavoriteStatus()
						
						// 加载发布者信息（确保有publisher_id才调用）
						if (this.contentData.publisher_id) {
							this.loadPublisherInfo()
						}
						
						// 在内容成功加载后记录浏览历史
						this.recordBrowseHistory()
						
						// 更新分享数据
						this.updateShareData()
					} else {
						throw new Error(result.message || '获取内容详情失败')
					}
				} catch (error) {
					console.error('加载内容详情失败:', error)
					this.loadError = true
					this.errorMsg = error.message || '加载失败，请重试'
				} finally {
					this.loading = false
				}
			},
			
			// 加载收藏状态
			async loadFavoriteStatus() {
				try {
					const result = await getFavoriteStatus(this.contentId)
					if (result.code === 0 && result.data) {
						this.isCollected = result.data.isFavorite
					}
				} catch (error) {
					console.error('获取收藏状态失败:', error)
				}
			},
			
			// 加载评论列表
			async loadComments() {
				if (!this.contentData.id) return
				
				this.commentLoading = true
				
				try {
					const result = await getCommentList(this.contentData.id, {
						page: 1,
						pageSize: 10
					})
					
					if (result.code === 0 && result.data) {
						// 格式化评论数据，适配实际API返回的字段名
						this.contentData.comments = (result.data.list || []).map(item => ({
							id: item.id,
							avatar: item.avatarUrl || '',
							name: item.realName || '匿名用户',
							content: item.comment || '',
							time: this.formatPublishTime(item.createdAt)
						}))
						
						// 更新评论总数
						if (result.data.total !== undefined) {
							this.contentData.comment_count = result.data.total
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
			handleBack() {
				uni.navigateBack()
			},
			previewImage(index) {
				uni.previewImage({
					current: index,
					urls: this.contentData.images,
					success: () => {
						console.log('图片预览成功')
					},
					fail: (err) => {
						console.error('图片预览失败:', err)
					}
				})
			},
			async handleFollow() {
				try {
					// 安全检查：确保publisher_id非空
					const publisherId = this.contentData.publisher_id;
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
						// 切换关注状态
						this.isFollowed = !this.isFollowed;
						
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
									contentId: this.contentId
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
									
									// 刷新评论列表，而不是直接添加到本地
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
			handleShare() {
				// 更新分享数据
				this.updateShareData();
				// 显示分享菜单，包括分享朋友圈
				uni.showShareMenu({
					withShareTicket: true,
					menus: ['shareAppMessage', 'shareTimeline']
				});
			},
			async handleMessage() {
				// 获取发布者ID
				const publisherId = this.contentData.publisher_id || '';
				if (!publisherId) {
					uni.showToast({
						title: '无法获取发布者信息',
						icon: 'none'
					});
					return;
				}
				
				try {
					// 调用创建会话接口
					const result = await createConversation({
						targetId: parseInt(publisherId)
					});
					
					if (result && result.code === 0) {
						// 会话创建成功后，跳转到聊天页面
						uni.navigateTo({
							url: `/pages/chat/detail?id=${publisherId}&nickName=${encodeURIComponent(this.contentData.publisher_name || '用户')}`,
							fail: (err) => {
								console.error('跳转到聊天页面失败:', err);
								uni.showToast({
									title: '无法打开聊天页面',
									icon: 'none'
								});
							}
						});
					} else {
						throw new Error(result?.message || '创建会话失败');
					}
				} catch (error) {
					console.error('私信功能错误:', error);
					uni.showToast({
						title: '无法发起私信，请重试',
						icon: 'none'
					});
				}
			},
			// 处理收藏
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
						result = await cancelFavorite(this.contentId)
					} else {
						// 未收藏，执行添加收藏
						result = await addFavorite(this.contentId)
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
					// 检查URL是否有效
					if (!this.isValidUrl(item.linkUrl)) {
						uni.showToast({
							title: '无效的URL',
							icon: 'none'
						});
						return;
					}
					
					// 添加安全处理
					if (typeof wx !== 'undefined' && wx.setWebViewSecurity) {
						// 先设置WebView安全模式
						wx.setWebViewSecurity({
							enable: true,
							success: () => {
								console.log('成功设置WebView安全模式');
								this.navigateToWebView(item.linkUrl);
							},
							fail: (err) => {
								console.error('设置WebView安全模式失败:', err);
								// 失败也继续跳转，但可能会有警告
								this.navigateToWebView(item.linkUrl);
							}
						});
					} else {
						// 如果不支持setWebViewSecurity，直接跳转
						this.navigateToWebView(item.linkUrl);
					}
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
			// 检查URL是否有效
			isValidUrl(url) {
				if (!url) return false;
				try {
					new URL(url);
					return true;
				} catch (e) {
					return false;
				}
			},
			// 导航到WebView页面
			navigateToWebView(url) {
				// 添加跨源隔离参数
				const separator = url.includes('?') ? '&' : '?';
				const secureParams = `${separator}coop=cross-origin&coep=require-corp`;
				let secureUrl = url;
				
				// 只对http和https开头的URL添加参数
				if (url.startsWith('http://') || url.startsWith('https://')) {
					// 为了不修改原始URL太多，只添加必要的跨源隔离参数
					secureUrl = url + secureParams;
				}
				
				// 跳转到WebView页面
				uni.navigateTo({
					url: '/pages/webview/index?url=' + encodeURIComponent(secureUrl),
					fail: (err) => {
						console.error('网页跳转失败', url, err);
						uni.showToast({
							title: '网页打开失败',
							icon: 'none'
						});
					}
				});
			},
			// 记录浏览历史 - 简化错误处理
			async recordBrowseHistory() {
				try {
					// 使用服务器返回的内容ID，而不是路由参数中的ID
					const contentId = this.contentData.id || this.contentId;
					if (!contentId) {
						console.error('记录浏览历史失败: 缺少有效的内容ID');
						return;
					}
					
					// 使用article作为内容类型
					await addBrowseRecord(contentId, 'article');
				} catch (error) {
					console.error('记录浏览历史失败:', error);
				}
			},
			// 加载发布者信息
			async loadPublisherInfo() {
				try {
					// 安全检查：确保publisher_id存在且有效
					const publisherId = this.contentData.publisher_id;
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
							this.contentData.publisher_name = publisherData.real_name;
						}
						
						if (publisherData.avatar_url) {
							this.contentData.publisher_avatar = publisherData.avatar_url;
						}
						
						if (publisherData.publish_count !== undefined) {
							this.contentData.publish_count = publisherData.publish_count;
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
						bannerType: 'home'
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
			handleSwiperChange(e) {
				this.currentBannerIndex = e.detail.current;
			},
			// 更新分享数据
			async updateShareData() {
				try {
					// 获取分享配置
					const settings = await getShareSettings();
					
					if (settings && this.contentData) {
						// 分享标题处理
						// 1. 如果有内容标题，优先使用内容标题
						// 2. 否则使用后端配置的content_share_text
						const shareTitle = this.contentData.title || settings.content_share_text || '查看内容详情';
						
						// 分享图片处理：
						// 1. 优先使用内容页自带的封面图
						// 2. 如果没有封面图，使用内容的第一张图片
						// 3. 如果都没有，使用后端配置的content_share_image
						let shareImage = '';
						if (this.contentData.cover) {
							shareImage = this.contentData.cover;
						} else if (this.contentData.images && this.contentData.images.length > 0) {
							shareImage = this.contentData.images[0];
						} else if (settings.content_share_image) {
							shareImage = settings.content_share_image;
						}
						
						// 直接设置分享数据，避免再次调用initShareData从而获得更精确的控制
						this.shareData = {
							title: shareTitle,
							imageUrl: shareImage,
							path: `/pages/content/detail?id=${this.contentId}`
						};
						
						console.log('内容分享数据已更新:', this.shareData);
					} else {
						// 配置获取失败，调用混入的默认初始化
						this.initShareData();
					}
				} catch (error) {
					console.error('更新分享数据失败:', error);
					// 配置获取失败，调用混入的默认初始化
					this.initShareData();
				}
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
	
	/* 导航栏 */
	.nav-bar {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 2;
		background-color: #FFFFFF;
		border-bottom: 1rpx solid #EEEEEE;
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
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		padding: 0 88rpx;
	}
	
	/* 内容区域 */
	.content-scroll {
		flex: 1;
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
	
	.tag-row {
		display: flex;
		gap: 12rpx;
		margin-bottom: 20rpx;
		padding: 30rpx 30rpx 0;
	}
	
	.tag {
		padding: 4rpx 12rpx;
		border-radius: 4rpx;
		font-size: 22rpx;
	}
	
	.top-tag {
		background-color: #fc3e2b;
		color: #ffffff;
	}
	
	.category-tag {
		background-color: #f5f5f5;
		color: #666666;
	}
	
	.content-title {
		font-size: 36rpx;
		font-weight: 600;
		color: #262626;
		margin-bottom: 20rpx;
		padding: 0 30rpx;
	}
	
	.publish-info {
		display: flex;
		align-items: center;
		gap: 20rpx;
		margin-bottom: 30rpx;
		padding: 0 30rpx;
	}
	
	.publisher {
		font-size: 28rpx;
		color: #333333;
	}
	
	.publish-time {
		font-size: 24rpx;
		color: #999999;
	}
	
	.content-text {
		font-size: 28rpx;
		color: #333333;
		line-height: 1.6;
		white-space: pre-wrap;
		margin-bottom: 20rpx;
	}
	
	
	/* 内容区域样式 */
	.content-body {
		display: flex;
		flex-direction: column;
		gap: 20rpx;
		padding: 0 30rpx;
	}
	
	/* 图片列表样式 */
	.image-list {
		display: flex;
		flex-direction: column;
		gap: 20rpx;
	}
	
	.image-list image {
		width: 100%;
		border-radius: 12rpx;
		background-color: #f5f5f5;
	}
	
	/* 发布人信息卡片 */
	.publisher-card {
		background-color: #FFFFFF;
		padding: 30rpx;
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
	
	/* 留言区域样式 */
	.comment-section {
		background-color: #FFFFFF;
		padding: 30rpx;
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
		color: #007AFF;
		font-size: 26rpx;
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
	
	/* 加载中样式 */
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
	
	/* 错误状态样式 */
	.error-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 300rpx;
		margin-top: 100rpx;
	}
	
	.error-text {
		font-size: 28rpx;
		color: #999999;
		margin-bottom: 20rpx;
	}
	
	.retry-btn {
		padding: 16rpx 40rpx;
		background-color: #fc3e2b;
		color: #FFFFFF;
		font-size: 28rpx;
		border-radius: 40rpx;
	}
</style> 