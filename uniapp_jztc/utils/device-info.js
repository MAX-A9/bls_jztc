/**
 * 设备信息工具
 * 统一封装设备信息相关API，解决兼容性问题
 */

// 缓存获取的系统信息，避免重复调用API
let cachedWindowInfo = null;
let cachedDeviceInfo = null;
let cachedAppBaseInfo = null;
let cachedSystemSetting = null;

/**
 * 获取窗口信息
 * @returns {Object} 窗口信息对象
 */
export function getWindowInfo() {
  if (cachedWindowInfo) return cachedWindowInfo;
  
  try {
    // 优先使用新API
    if (uni.canIUse('getWindowInfo')) {
      cachedWindowInfo = uni.getWindowInfo();
      return cachedWindowInfo;
    }
    
    // 回退到兼容方式
    const sysInfo = uni.getSystemInfo({ sync: true });
    cachedWindowInfo = {
      windowWidth: sysInfo.windowWidth,
      windowHeight: sysInfo.windowHeight,
      screenWidth: sysInfo.screenWidth,
      screenHeight: sysInfo.screenHeight,
      statusBarHeight: sysInfo.statusBarHeight,
      safeArea: sysInfo.safeArea || {
        top: 0,
        right: sysInfo.screenWidth,
        bottom: sysInfo.screenHeight,
        left: 0
      },
      pixelRatio: sysInfo.pixelRatio
    };
    return cachedWindowInfo;
  } catch (error) {
    console.error('获取窗口信息失败:', error);
    // 返回默认值
    return {
      windowWidth: 375,
      windowHeight: 667,
      screenWidth: 375,
      screenHeight: 667,
      statusBarHeight: 20,
      safeArea: {
        top: 0,
        right: 375,
        bottom: 667,
        left: 0
      },
      pixelRatio: 2
    };
  }
}

/**
 * 获取设备信息
 * @returns {Object} 设备信息对象
 */
export function getDeviceInfo() {
  if (cachedDeviceInfo) return cachedDeviceInfo;
  
  try {
    // 优先使用新API
    if (uni.canIUse('getDeviceInfo')) {
      cachedDeviceInfo = uni.getDeviceInfo();
      return cachedDeviceInfo;
    }
    
    // 回退到兼容方式
    const sysInfo = uni.getSystemInfo({ sync: true });
    cachedDeviceInfo = {
      brand: sysInfo.brand,
      model: sysInfo.model,
      system: sysInfo.system,
      platform: sysInfo.platform,
      deviceId: sysInfo.deviceId,
      devicePixelRatio: sysInfo.pixelRatio,
      deviceOrientation: sysInfo.deviceOrientation,
      deviceType: sysInfo.deviceType
    };
    return cachedDeviceInfo;
  } catch (error) {
    console.error('获取设备信息失败:', error);
    // 返回默认值
    return {
      brand: '',
      model: '',
      system: '',
      platform: '',
      deviceId: '',
      devicePixelRatio: 2,
      deviceOrientation: 'portrait',
      deviceType: 'phone'
    };
  }
}

/**
 * 获取应用基础信息
 * @returns {Object} 应用基础信息对象
 */
export function getAppBaseInfo() {
  if (cachedAppBaseInfo) return cachedAppBaseInfo;
  
  try {
    // 优先使用新API
    if (uni.canIUse('getAppBaseInfo')) {
      cachedAppBaseInfo = uni.getAppBaseInfo();
      return cachedAppBaseInfo;
    }
    
    // 回退到兼容方式
    const sysInfo = uni.getSystemInfo({ sync: true });
    cachedAppBaseInfo = {
      SDKVersion: sysInfo.SDKVersion,
      appName: sysInfo.appName,
      appVersion: sysInfo.version,
      appLanguage: sysInfo.language,
      theme: sysInfo.theme,
      host: sysInfo.host
    };
    return cachedAppBaseInfo;
  } catch (error) {
    console.error('获取应用基础信息失败:', error);
    // 返回默认值
    return {
      SDKVersion: '',
      appName: '',
      appVersion: '',
      appLanguage: 'zh-Hans',
      theme: 'light'
    };
  }
}

/**
 * 获取系统设置
 * @returns {Object} 系统设置对象
 */
export function getSystemSetting() {
  if (cachedSystemSetting) return cachedSystemSetting;
  
  try {
    // 优先使用新API
    if (uni.canIUse('getSystemSetting')) {
      cachedSystemSetting = uni.getSystemSetting();
      return cachedSystemSetting;
    }
    
    // 回退到兼容方式
    const sysInfo = uni.getSystemInfo({ sync: true });
    cachedSystemSetting = {
      bluetoothEnabled: sysInfo.bluetoothEnabled,
      locationEnabled: sysInfo.locationEnabled,
      wifiEnabled: sysInfo.wifiEnabled,
      deviceOrientation: sysInfo.deviceOrientation,
      // 兼容性处理
      locationAuthorized: sysInfo.locationAuthorized,
      microphoneAuthorized: sysInfo.microphoneAuthorized,
      cameraAuthorized: sysInfo.cameraAuthorized,
      notificationAuthorized: sysInfo.notificationAuthorized,
      notificationAlertAuthorized: sysInfo.notificationAlertAuthorized,
      notificationBadgeAuthorized: sysInfo.notificationBadgeAuthorized,
      notificationSoundAuthorized: sysInfo.notificationSoundAuthorized,
      bluetoothAuthorized: sysInfo.bluetoothAuthorized,
    };
    return cachedSystemSetting;
  } catch (error) {
    console.error('获取系统设置信息失败:', error);
    // 返回默认值
    return {
      bluetoothEnabled: false,
      locationEnabled: false,
      wifiEnabled: true,
      deviceOrientation: 'portrait'
    };
  }
}

/**
 * 获取状态栏高度
 * @returns {Number} 状态栏高度(px)
 */
export function getStatusBarHeight() {
  const windowInfo = getWindowInfo();
  return windowInfo.statusBarHeight || 20;
}

/**
 * 获取导航栏高度
 * @returns {Number} 导航栏高度(px)
 */
export function getNavigationBarHeight() {
  const statusBarHeight = getStatusBarHeight();
  const deviceInfo = getDeviceInfo();
  const isIOS = (deviceInfo.platform === 'ios') || 
    (/ios/i.test(deviceInfo.system || ''));
  
  // iOS 和 Android 导航栏高度不同
  const navBarHeight = isIOS ? 44 : 48;
  return statusBarHeight + navBarHeight;
}

/**
 * 获取安全区域
 * @returns {Object} 安全区域对象
 */
export function getSafeAreaInsets() {
  const windowInfo = getWindowInfo();
  const safeArea = windowInfo.safeArea || {};
  
  return {
    top: safeArea.top || 0,
    right: (windowInfo.screenWidth - safeArea.right) || 0,
    bottom: (windowInfo.screenHeight - safeArea.bottom) || 0,
    left: safeArea.left || 0
  };
}

/**
 * 是否为 iPhone X 系列（有底部安全区域的机型）
 * @returns {Boolean} 是否为 iPhone X 系列
 */
export function isIphoneX() {
  const deviceInfo = getDeviceInfo();
  const safeAreaInsets = getSafeAreaInsets();
  
  return (deviceInfo.platform === 'ios' || 
    /ios/i.test(deviceInfo.system || '')) && 
    safeAreaInsets.bottom > 0;
}

/**
 * rpx 转 px
 * @param {Number} rpx - 要转换的rpx值
 * @returns {Number} 转换后的px值
 */
export function rpx2px(rpx) {
  const windowInfo = getWindowInfo();
  return (rpx / 750) * windowInfo.windowWidth;
}

/**
 * px 转 rpx
 * @param {Number} px - 要转换的px值
 * @returns {Number} 转换后的rpx值
 */
export function px2rpx(px) {
  const windowInfo = getWindowInfo();
  return (px * 750) / windowInfo.windowWidth;
}

// 导出默认对象，方便一次性导入所有方法
export default {
  getWindowInfo,
  getDeviceInfo,
  getAppBaseInfo,
  getSystemSetting,
  getStatusBarHeight,
  getNavigationBarHeight,
  getSafeAreaInsets,
  isIphoneX,
  rpx2px,
  px2rpx
}; 