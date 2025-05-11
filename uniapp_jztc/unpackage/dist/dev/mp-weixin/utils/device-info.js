"use strict";
const common_vendor = require("../common/vendor.js");
let cachedWindowInfo = null;
let cachedDeviceInfo = null;
let cachedAppBaseInfo = null;
let cachedSystemSetting = null;
function getWindowInfo() {
  if (cachedWindowInfo)
    return cachedWindowInfo;
  try {
    if (common_vendor.index.canIUse("getWindowInfo")) {
      cachedWindowInfo = common_vendor.index.getWindowInfo();
      return cachedWindowInfo;
    }
    const sysInfo = common_vendor.index.getSystemInfo({ sync: true });
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
    common_vendor.index.__f__("error", "at utils/device-info.js:44", "获取窗口信息失败:", error);
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
function getDeviceInfo() {
  if (cachedDeviceInfo)
    return cachedDeviceInfo;
  try {
    if (common_vendor.index.canIUse("getDeviceInfo")) {
      cachedDeviceInfo = common_vendor.index.getDeviceInfo();
      return cachedDeviceInfo;
    }
    const sysInfo = common_vendor.index.getSystemInfo({ sync: true });
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
    common_vendor.index.__f__("error", "at utils/device-info.js:91", "获取设备信息失败:", error);
    return {
      brand: "",
      model: "",
      system: "",
      platform: "",
      deviceId: "",
      devicePixelRatio: 2,
      deviceOrientation: "portrait",
      deviceType: "phone"
    };
  }
}
function getAppBaseInfo() {
  if (cachedAppBaseInfo)
    return cachedAppBaseInfo;
  try {
    if (common_vendor.index.canIUse("getAppBaseInfo")) {
      cachedAppBaseInfo = common_vendor.index.getAppBaseInfo();
      return cachedAppBaseInfo;
    }
    const sysInfo = common_vendor.index.getSystemInfo({ sync: true });
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
    common_vendor.index.__f__("error", "at utils/device-info.js:132", "获取应用基础信息失败:", error);
    return {
      SDKVersion: "",
      appName: "",
      appVersion: "",
      appLanguage: "zh-Hans",
      theme: "light"
    };
  }
}
function getSystemSetting() {
  if (cachedSystemSetting)
    return cachedSystemSetting;
  try {
    if (common_vendor.index.canIUse("getSystemSetting")) {
      cachedSystemSetting = common_vendor.index.getSystemSetting();
      return cachedSystemSetting;
    }
    const sysInfo = common_vendor.index.getSystemInfo({ sync: true });
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
      bluetoothAuthorized: sysInfo.bluetoothAuthorized
    };
    return cachedSystemSetting;
  } catch (error) {
    common_vendor.index.__f__("error", "at utils/device-info.js:177", "获取系统设置信息失败:", error);
    return {
      bluetoothEnabled: false,
      locationEnabled: false,
      wifiEnabled: true,
      deviceOrientation: "portrait"
    };
  }
}
function getStatusBarHeight() {
  const windowInfo = getWindowInfo();
  return windowInfo.statusBarHeight || 20;
}
function getNavigationBarHeight() {
  const statusBarHeight = getStatusBarHeight();
  const deviceInfo2 = getDeviceInfo();
  const isIOS = deviceInfo2.platform === "ios" || /ios/i.test(deviceInfo2.system || "");
  const navBarHeight = isIOS ? 44 : 48;
  return statusBarHeight + navBarHeight;
}
function getSafeAreaInsets() {
  const windowInfo = getWindowInfo();
  const safeArea = windowInfo.safeArea || {};
  return {
    top: safeArea.top || 0,
    right: windowInfo.screenWidth - safeArea.right || 0,
    bottom: windowInfo.screenHeight - safeArea.bottom || 0,
    left: safeArea.left || 0
  };
}
function isIphoneX() {
  const deviceInfo2 = getDeviceInfo();
  const safeAreaInsets = getSafeAreaInsets();
  return (deviceInfo2.platform === "ios" || /ios/i.test(deviceInfo2.system || "")) && safeAreaInsets.bottom > 0;
}
function rpx2px(rpx) {
  const windowInfo = getWindowInfo();
  return rpx / 750 * windowInfo.windowWidth;
}
function px2rpx(px) {
  const windowInfo = getWindowInfo();
  return px * 750 / windowInfo.windowWidth;
}
const deviceInfo = {
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
exports.deviceInfo = deviceInfo;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/device-info.js.map
