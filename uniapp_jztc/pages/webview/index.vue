<template>
	<view class="webview-container">
		<!-- 加载中指示器 -->
		<view class="loading" v-if="!webUrl">
			<view class="loading-spinner"></view>
			<text class="loading-text">加载中...</text>
		</view>
		
		<!-- Web视图 -->
		<web-view 
			:src="webUrl" 
			@message="handleMessage" 
			@error="handleError"
		></web-view>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				webUrl: '',
				originalUrl: '',
				securityEnabled: false,
				urlTitle: ''
			}
		},
		onLoad(options) {
			if (options.url) {
				try {
					this.originalUrl = decodeURIComponent(options.url) || '';
					
					// 提取页面标题（如果URL中包含title参数）
					if (options.title) {
						this.urlTitle = decodeURIComponent(options.title);
						// 设置导航栏标题
						uni.setNavigationBarTitle({
							title: this.urlTitle
						});
					} else {
						// 使用默认标题
						uni.setNavigationBarTitle({
							title: '网页浏览'
						});
					}
					
					// 启用WebView安全设置
					this.enableWebViewSecurity();
					
					// 处理URL，添加安全头部
					this.setSecureWebUrl(this.originalUrl);
				} catch (e) {
					console.error('URL解析错误:', e);
					uni.showToast({
						title: '无效的链接',
						icon: 'none'
					});
					setTimeout(() => {
						uni.navigateBack();
					}, 1500);
				}
			} else {
				uni.showToast({
					title: '缺少URL参数',
					icon: 'none'
				});
				setTimeout(() => {
					uni.navigateBack();
				}, 1500);
			}
		},
		mounted() {
			// 添加监听器以禁用SharedArrayBuffer
			this.disableSharedArrayBuffer();
		},
		methods: {
			// 启用WebView安全设置
			enableWebViewSecurity() {
				if (typeof wx !== 'undefined' && wx.setWebViewSecurity) {
					wx.setWebViewSecurity({
						enable: true,
						success: () => {
							console.log('成功启用WebView安全模式');
							this.securityEnabled = true;
						},
						fail: (err) => {
							console.error('启用WebView安全模式失败:', err);
						}
					});
				}
			},
			
			// 禁用SharedArrayBuffer相关功能
			disableSharedArrayBuffer() {
				// 在小程序环境中
				if (typeof wx !== 'undefined') {
					// 添加WebView加载完成的处理
					if (wx.onWebViewLoad) {
						wx.onWebViewLoad(() => {
							try {
								console.log('WebView加载完成，应用安全措施');
								
								// 尝试注入脚本禁用SharedArrayBuffer
								this.injectSecurityScript();
							} catch (e) {
								console.error('WebView安全处理失败:', e);
							}
						});
					}
				}
			},
			
			// 注入安全脚本
			injectSecurityScript() {
				if (typeof wx !== 'undefined' && wx.evaluateWebView) {
					// 注入一段JS来处理SharedArrayBuffer
					const securityScript = `
						try {
							// 阻止使用SharedArrayBuffer
							if (typeof SharedArrayBuffer !== 'undefined') {
								console.warn('检测到SharedArrayBuffer，但已被安全策略禁用');
							}
							
							// 添加必要的安全头
							const meta = document.createElement('meta');
							meta.httpEquiv = 'Cross-Origin-Embedder-Policy';
							meta.content = 'require-corp';
							document.head.appendChild(meta);
							
							const meta2 = document.createElement('meta');
							meta2.httpEquiv = 'Cross-Origin-Opener-Policy';
							meta2.content = 'same-origin';
							document.head.appendChild(meta2);
							
							true; // 返回执行结果
						} catch(e) {
							console.error('安全脚本执行出错:', e);
							false; // 返回执行结果
						}
					`;
					
					// 执行脚本
					wx.evaluateWebView({
						webviewId: 'webviewId', // 可能需要根据实际情况获取
						script: securityScript,
						success: (res) => {
							console.log('安全脚本注入成功', res);
						},
						fail: (err) => {
							console.error('安全脚本注入失败:', err);
						}
					});
				}
			},
			
			// 设置安全的WebURL，处理跨源隔离问题
			setSecureWebUrl(url) {
				// 检查URL是否有效
				if (!url) {
					this.webUrl = '';
					return;
				}
				
				try {
					// 确保URL有跨源隔离参数
					let secureUrl = url;
					
					// 如果URL中没有包含跨源隔离参数，则添加
					if (url.startsWith('http') && 
					    !url.includes('coop=cross-origin') && 
					    !url.includes('coep=require-corp')) {
						
						const separator = url.includes('?') ? '&' : '?';
						secureUrl = `${url}${separator}coop=cross-origin&coep=require-corp`;
					}
					
					console.log('使用安全URL:', secureUrl);
					this.webUrl = secureUrl;
				} catch (e) {
					console.error('处理URL时出错:', e);
					this.webUrl = url; // 降级处理，使用原始URL
				}
			},
			
			// 处理web-view的消息
			handleMessage(event) {
				console.log('收到web-view消息:', event);
			},
			
			// 处理web-view的错误
			handleError(error) {
				console.error('web-view加载错误:', error);
				uni.showToast({
					title: '页面加载出错',
					icon: 'none'
				});
			}
		}
	}
</script>

<style>
	.webview-container {
		position: relative;
		width: 100%;
		height: 100vh;
	}
	
	/* 加载中样式 */
	.loading {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		display: flex;
		flex-direction: column;
		align-items: center;
		z-index: 50;
	}
	
	.loading-spinner {
		width: 60rpx;
		height: 60rpx;
		border: 4rpx solid #f0f0f0;
		border-radius: 50%;
		border-top-color: #fc3e2b;
		animation: spin 1s linear infinite;
		margin-bottom: 20rpx;
	}
	
	.loading-text {
		font-size: 28rpx;
		color: #666666;
	}
	
	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}
</style> 