<script>
	import { checkAndAutoLogin } from './utils/auth.js';
	import navigationControl from './utils/navigation-control.js';
	import { initUniCompatibility } from './utils/uni-compatibility.js';
	import LoginModal from './components/login-modal/index.vue';
	
	// 添加初始化状态标志
	let isInitialized = false;
	
	export default {
		components: {
			LoginModal
		},
		onLaunch: function() {
			// 只在首次加载时执行初始化
			if (!isInitialized) {
				// 初始化 uni-app 兼容性适配
				initUniCompatibility();
				
				// 初始化请求拦截器
				this.initRequestInterceptor();
				
				// 初始化导航拦截器
				this.initNavigationInterceptor();
				
				// 标记已初始化完成
				isInitialized = true;
			}
			
			// 设置延迟以确保Vuex初始化完成后再执行自动登录
			setTimeout(() => {
				// 添加超时处理，避免登录请求阻塞应用加载
				const loginTimeout = setTimeout(() => {
					console.log('自动登录超时，继续应用初始化');
				}, 3000); // 设置3秒超时
				
				this.autoLogin().finally(() => {
					clearTimeout(loginTimeout);
				});
			}, 100);
			
			// 禁用SharedArrayBuffer相关功能，避免安全警告
			this.disableSharedArrayBuffer();
		},
		onShow: function() {
			// App Show时不重复执行自动登录，避免重复请求
			// 仅在onLaunch中执行自动登录即可
			
			// 添加区域列表的初始化，确保每次打开应用时都加载最新的区域列表
			setTimeout(() => {
				// 获取store实例
				const store = this.$store;
				if (store && store.dispatch) {
					// 加载区域列表数据
					console.log('App显示时开始加载区域列表数据...');
					store.dispatch('region/getRegionList')
						.then((regionList) => {
							// 区域列表加载成功后，检查并更新当前位置
							if (regionList && regionList.length > 0) {
								// 检查当前是否已设置位置
								const currentLocationId = uni.getStorageSync('currentLocationId');
								
								// 如果没有设置位置ID，则使用第一个区域
								if (!currentLocationId) {
									const firstRegion = regionList[0];
									uni.setStorageSync('currentLocationId', firstRegion.id);
									uni.setStorageSync('currentLocation', firstRegion.name);
									
									console.log('App显示时更新位置信息:', firstRegion.name, firstRegion.id);
									
									// 通知所有监听此事件的组件位置已变更
									uni.$emit('locationChanged', { regionId: firstRegion.id });
								} else {
									// 如果已有位置ID，也发送locationChanged事件以触发内容刷新
									console.log('App显示时发送位置更新事件，使用已有位置ID:', currentLocationId);
									uni.$emit('locationChanged', { regionId: currentLocationId });
								}
							}
						})
						.catch(err => {
							console.error('初始化区域列表失败:', err);
						});
				}
			}, 200);
		},
		onHide: function() {
			// App隐藏时的处理
		},
		methods: {
			// 自动登录方法
			async autoLogin() {
				try {
					// 检查是否已登录
					const hasToken = this.$store.getters['user/token'];
					if (!hasToken) {
						const result = await checkAndAutoLogin();
						if (result) {
							// 先设置token，确保登录状态
							if (result.token) {
								this.$store.commit('user/SET_TOKEN', result.token);
							}
							// 再设置用户信息
							this.$store.commit('user/SET_USER_INFO', result.userInfo || result);
						}
					} else {
						// 有token但没有完整用户信息时，获取用户信息
						const userInfo = this.$store.getters['user/userInfo'];
						// 检查是否有用户ID (clientId或id)和基本信息
						const hasUserId = userInfo && (userInfo.clientId || userInfo.id);
						const hasBasicInfo = userInfo && userInfo.realName;
						
						if (!hasUserId || !hasBasicInfo) {
							try {
								await this.$store.dispatch('user/getUserInfo');
							} catch (infoError) {
								// 尝试静默登录
								await this.$store.dispatch('user/silentLogin');
							}
						}
					}
				} catch (error) {
					// 登录失败时不做特殊处理，用户可以在页面上手动登录
					console.log('自动登录失败:', error);
				}
				return Promise.resolve(); // 确保总是返回已解决的Promise
			},
			
			// 初始化请求拦截器
			initRequestInterceptor() {
				// 添加请求拦截器
				uni.addInterceptor('request', {
					fail(err) {
						return err;
					}
				});
				
				// 添加响应拦截器
				uni.addInterceptor('response', {
					response(res) {
						// 统一处理401认证失败
						if (res.statusCode === 401) {
							// 提示用户重新登录
							uni.showToast({
								title: '登录已过期，请重新登录',
								icon: 'none'
							});
							
							// 清除登录状态
							const store = getApp().$vm.$store;
							if (store) {
								store.dispatch('user/logout');
								// 显示登录弹窗
								store.dispatch('user/showLoginModal');
							}
						}
						return res;
					}
				});
			},
			
			// 初始化导航拦截器
			initNavigationInterceptor() {
				// 初始化基本的导航拦截
				navigationControl.initNavigationInterceptors();
				
				// 处理页面返回的特殊情况
				uni.addInterceptor('navigateBack', {
					invoke(args) {
						// 获取当前页面栈
						const pages = getCurrentPages();
						if (pages.length < 2) return args;
						
						// 当前页面和目标页面
						const currentPage = pages[pages.length - 1];
						const targetPage = pages[pages.length - 2];
						
						// 获取页面路径
						const currentPath = `/${currentPage.route}`;
						const targetPath = `/${targetPage.route}`;
						
						// 处理页面返回的轮询控制
						navigationControl.beforeNavigateBack(currentPath, targetPath);
						
						return args;
					}
				});
			},
			
			// 禁用SharedArrayBuffer相关功能
			disableSharedArrayBuffer() {
				// 在小程序环境中，通过设置自定义拦截器来处理网络请求头
				uni.addInterceptor('request', {
					invoke(args) {
						// 添加安全相关的请求头
						if (!args.header) {
							args.header = {};
						}
						// 防止使用SharedArrayBuffer
						args.header['Cross-Origin-Opener-Policy'] = 'same-origin';
						args.header['Cross-Origin-Embedder-Policy'] = 'require-corp';
						return args;
					}
				});
				
				// 拦截WebView组件的加载
				if (typeof wx !== 'undefined' && wx.onWebViewLoad) {
					wx.onWebViewLoad(function() {
						try {
							// 禁用SharedArrayBuffer
							if (typeof SharedArrayBuffer !== 'undefined') {
								// 在开发环境输出日志，不影响功能
								console.warn('已禁用SharedArrayBuffer以提高安全性');
							}
						} catch (e) {
							// 忽略错误
						}
					});
				}
			}
		}
	}
</script>

<style>
	/*每个页面公共css */
	/* 隐藏所有滚动条 */
	::-webkit-scrollbar {
		display: none;
		width: 0 !important;
		height: 0 !important;
	}
	
	/* 隐藏scroll-view的滚动条 - 修改为更符合微信小程序规范的写法 */
	.scroll-view-container ::-webkit-scrollbar {
		display: none;
	}
	
	/* 添加全局滚动条隐藏 */
	::-webkit-scrollbar {
		width: 0;
		height: 0;
		color: transparent;
		display: none;
	}
	
	page {
		height: 100%;
		overflow: hidden;
	}
</style>

<!-- 添加登录弹窗组件 -->
<template>
	<login-modal></login-modal>
</template>
