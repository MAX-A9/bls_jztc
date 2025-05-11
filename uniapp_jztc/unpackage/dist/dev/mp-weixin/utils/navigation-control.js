"use strict";
const common_vendor = require("../common/vendor.js");
const utils_messagePolling = require("./message-polling.js");
const tabBarPages = [
  "/pages/index/index",
  "/pages/community/index",
  "/pages/publish/index",
  "/pages/message/index",
  "/pages/my/index"
];
function isTabBarPage(pagePath) {
  const path = pagePath.split("?")[0];
  return tabBarPages.includes(path);
}
function beforeNavigateTo(url) {
  const pagePath = url.split("?")[0];
  utils_messagePolling.messagePollingService.setCurrentPage(pagePath);
  if (!isTabBarPage(pagePath)) {
    utils_messagePolling.messagePollingService.pausePollingGlobally();
    common_vendor.index.__f__("log", "at utils/navigation-control.js:41", "导航到非底部导航页面，暂停轮询:", pagePath);
  }
}
function beforeNavigateBack(currentPage, targetPage) {
  utils_messagePolling.messagePollingService.setCurrentPage(targetPage);
  if (!isTabBarPage(currentPage) && isTabBarPage(targetPage)) {
    utils_messagePolling.messagePollingService.resumePollingGlobally();
    common_vendor.index.__f__("log", "at utils/navigation-control.js:57", "返回到底部导航页面，恢复轮询:", targetPage);
  }
}
function initNavigationInterceptors() {
  const originalNavigateTo = common_vendor.index.navigateTo;
  const originalRedirectTo = common_vendor.index.redirectTo;
  const originalReLaunch = common_vendor.index.reLaunch;
  const originalSwitchTab = common_vendor.index.switchTab;
  const originalNavigateBack = common_vendor.index.navigateBack;
  common_vendor.index.navigateTo = function(options) {
    const url = options.url;
    beforeNavigateTo(url);
    return originalNavigateTo.call(this, options);
  };
  common_vendor.index.redirectTo = function(options) {
    const url = options.url;
    beforeNavigateTo(url);
    return originalRedirectTo.call(this, options);
  };
  common_vendor.index.reLaunch = function(options) {
    const url = options.url;
    utils_messagePolling.messagePollingService.pausePollingGlobally();
    if (isTabBarPage(url.split("?")[0])) {
      setTimeout(() => {
        utils_messagePolling.messagePollingService.resumePollingGlobally();
      }, 100);
    }
    return originalReLaunch.call(this, options);
  };
  common_vendor.index.switchTab = function(options) {
    options.url;
    setTimeout(() => {
      utils_messagePolling.messagePollingService.resumePollingGlobally();
    }, 100);
    return originalSwitchTab.call(this, options);
  };
  common_vendor.index.navigateBack = function(options) {
    return originalNavigateBack.call(this, options);
  };
  common_vendor.index.__f__("log", "at utils/navigation-control.js:121", "导航拦截器初始化完成");
}
const navigationControl = {
  isTabBarPage,
  beforeNavigateTo,
  beforeNavigateBack,
  initNavigationInterceptors
};
exports.navigationControl = navigationControl;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/navigation-control.js.map
