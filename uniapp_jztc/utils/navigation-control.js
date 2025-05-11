/**
 * 导航控制工具
 * 用于管理页面导航和相关资源控制
 */
import messagePollingService from './message-polling.js';

// 底部导航页面路径列表
const tabBarPages = [
  '/pages/index/index',
  '/pages/community/index',
  '/pages/publish/index',
  '/pages/message/index',
  '/pages/my/index'
];

/**
 * 检查页面是否为底部导航页面
 * @param {String} pagePath - 页面路径
 * @returns {Boolean} 是否为底部导航页面
 */
function isTabBarPage(pagePath) {
  // 去掉可能存在的参数部分
  const path = pagePath.split('?')[0];
  return tabBarPages.includes(path);
}

/**
 * 页面跳转前处理
 * @param {String} url - 目标页面URL
 */
function beforeNavigateTo(url) {
  // 提取页面路径
  const pagePath = url.split('?')[0];
  
  // 设置当前页面路径
  messagePollingService.setCurrentPage(pagePath);
  
  // 如果将要跳转到非底部导航页面，暂停轮询
  if (!isTabBarPage(pagePath)) {
    messagePollingService.pausePollingGlobally();
    console.log('导航到非底部导航页面，暂停轮询:', pagePath);
  }
}

/**
 * 页面返回前处理
 * @param {String} currentPage - 当前页面路径
 * @param {String} targetPage - 返回目标页面路径
 */
function beforeNavigateBack(currentPage, targetPage) {
  // 设置当前页面将要变为的目标页面
  messagePollingService.setCurrentPage(targetPage);
  
  // 如果当前是非底部导航页面，要返回到底部导航页面，恢复轮询
  if (!isTabBarPage(currentPage) && isTabBarPage(targetPage)) {
    messagePollingService.resumePollingGlobally();
    console.log('返回到底部导航页面，恢复轮询:', targetPage);
  }
}

/**
 * 初始化导航拦截器
 * 拦截uni.navigateTo等方法，添加资源控制
 */
function initNavigationInterceptors() {
  // 保存原始方法
  const originalNavigateTo = uni.navigateTo;
  const originalRedirectTo = uni.redirectTo;
  const originalReLaunch = uni.reLaunch;
  const originalSwitchTab = uni.switchTab;
  const originalNavigateBack = uni.navigateBack;
  
  // 重写navigateTo
  uni.navigateTo = function(options) {
    const url = options.url;
    beforeNavigateTo(url);
    return originalNavigateTo.call(this, options);
  };
  
  // 重写redirectTo
  uni.redirectTo = function(options) {
    const url = options.url;
    beforeNavigateTo(url);
    return originalRedirectTo.call(this, options);
  };
  
  // 重写reLaunch
  uni.reLaunch = function(options) {
    const url = options.url;
    // reLaunch会关闭所有页面，无论目标是什么，都优先暂停轮询
    messagePollingService.pausePollingGlobally();
    
    // 如果目标是底部导航页面，则在页面加载后恢复轮询
    if (isTabBarPage(url.split('?')[0])) {
      setTimeout(() => {
        messagePollingService.resumePollingGlobally();
      }, 100);
    }
    
    return originalReLaunch.call(this, options);
  };
  
  // 重写switchTab
  uni.switchTab = function(options) {
    const url = options.url;
    // switchTab只能切换到底部导航页面，总是恢复轮询
    setTimeout(() => {
      messagePollingService.resumePollingGlobally();
    }, 100);
    
    return originalSwitchTab.call(this, options);
  };
  
  // 重写navigateBack
  uni.navigateBack = function(options) {
    // navigateBack比较特殊，需要获取当前页面和返回的目标页面
    // 由于无法直接知道返回到哪个页面，我们在App.vue中通过getCurrentPages()处理
    return originalNavigateBack.call(this, options);
  };
  
  console.log('导航拦截器初始化完成');
}

export default {
  isTabBarPage,
  beforeNavigateTo,
  beforeNavigateBack,
  initNavigationInterceptors
}; 