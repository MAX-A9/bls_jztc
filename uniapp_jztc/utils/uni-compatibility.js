/**
 * uni-app API 兼容性适配
 * 用于处理废弃的API，确保第三方UI组件可以正常工作
 */

// 缓存获取的系统信息，避免重复调用API
let cachedSystemInfo = null;

/**
 * 获取系统信息的兼容处理函数
 * 用于替代废弃的 getSystemInfoSync
 * @returns {Object} 系统信息
 */
function getCompatibleSystemInfo() {
  if (cachedSystemInfo) return cachedSystemInfo;
  
  try {
    // 使用新的API
    const windowInfo = uni.getWindowInfo ? uni.getWindowInfo() : {};
    const deviceInfo = uni.getDeviceInfo ? uni.getDeviceInfo() : {};
    const appBaseInfo = uni.getAppBaseInfo ? uni.getAppBaseInfo() : {};
    const systemSetting = uni.getSystemSetting ? uni.getSystemSetting() : {};
    
    // 组合成兼容旧API的格式
    cachedSystemInfo = {
      // 窗口信息
      windowWidth: windowInfo.windowWidth,
      windowHeight: windowInfo.windowHeight, 
      screenWidth: windowInfo.screenWidth,
      screenHeight: windowInfo.screenHeight,
      statusBarHeight: windowInfo.statusBarHeight,
      safeArea: windowInfo.safeArea,
      pixelRatio: windowInfo.pixelRatio,
      
      // 设备信息
      brand: deviceInfo.brand,
      model: deviceInfo.model,
      system: deviceInfo.system,
      platform: deviceInfo.platform,
      
      // APP基础信息
      SDKVersion: appBaseInfo.SDKVersion,
      appName: appBaseInfo.appName,
      appVersion: appBaseInfo.appVersion,
      appLanguage: appBaseInfo.appLanguage,
      theme: appBaseInfo.theme,
      
      // 系统设置
      bluetoothEnabled: systemSetting.bluetoothEnabled,
      locationEnabled: systemSetting.locationEnabled,
      wifiEnabled: systemSetting.wifiEnabled,
      
      // 兼容字段
      language: appBaseInfo.appLanguage,
      version: appBaseInfo.appVersion
    };
    
    return cachedSystemInfo;
  } catch (error) {
    console.error('获取兼容系统信息失败:', error);
    
    // 尝试使用旧的API方式获取
    try {
      const legacyInfo = uni.getSystemInfo({ sync: true });
      cachedSystemInfo = legacyInfo;
      return legacyInfo;
    } catch (fallbackError) {
      console.error('获取系统信息兼容处理失败:', fallbackError);
      // 返回最小默认值
      return {
        windowWidth: 375,
        windowHeight: 667,
        screenWidth: 375,
        screenHeight: 667,
        statusBarHeight: 20,
        pixelRatio: 2,
        platform: 'android',
        language: 'zh-Hans'
      };
    }
  }
}

/**
 * 兼容性API适配器初始化
 * 用于全局替换废弃的API
 */
function initUniCompatibility() {
  // 保存原始方法
  const originalGetSystemInfoSync = uni.getSystemInfoSync;
  
  // 替换为兼容方法
  uni.getSystemInfoSync = function() {
    console.warn('wx.getSystemInfoSync 已废弃，请使用 wx.getWindowInfo/wx.getDeviceInfo/wx.getAppBaseInfo 等代替');
    return getCompatibleSystemInfo();
  };
  
  // 确保uni.getSystemInfo方法也使用新的API
  const originalGetSystemInfo = uni.getSystemInfo;
  uni.getSystemInfo = function(options) {
    // 如果是同步调用
    if (options && options.sync === true) {
      if (typeof options.success === 'function') {
        options.success(getCompatibleSystemInfo());
      }
      return getCompatibleSystemInfo();
    }
    
    // 如果是异步调用
    return new Promise((resolve, reject) => {
      try {
        const info = getCompatibleSystemInfo();
        if (typeof options?.success === 'function') {
          options.success(info);
        }
        resolve(info);
      } catch (error) {
        if (typeof options?.fail === 'function') {
          options.fail(error);
        }
        reject(error);
      } finally {
        if (typeof options?.complete === 'function') {
          options.complete();
        }
      }
    });
  };
  
  console.log('uni-app 兼容性适配已完成');
}

// 导出函数
export { getCompatibleSystemInfo, initUniCompatibility }

// 导出默认对象
export default {
  getCompatibleSystemInfo,
  initUniCompatibility
}; 