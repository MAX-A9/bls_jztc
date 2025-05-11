"use strict";
const common_vendor = require("../common/vendor.js");
let cachedSystemInfo = null;
function getCompatibleSystemInfo() {
  if (cachedSystemInfo)
    return cachedSystemInfo;
  try {
    const windowInfo = common_vendor.index.getWindowInfo ? common_vendor.index.getWindowInfo() : {};
    const deviceInfo = common_vendor.index.getDeviceInfo ? common_vendor.index.getDeviceInfo() : {};
    const appBaseInfo = common_vendor.index.getAppBaseInfo ? common_vendor.index.getAppBaseInfo() : {};
    const systemSetting = common_vendor.index.getSystemSetting ? common_vendor.index.getSystemSetting() : {};
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
    common_vendor.index.__f__("error", "at utils/uni-compatibility.js:60", "获取兼容系统信息失败:", error);
    try {
      const legacyInfo = common_vendor.index.getSystemInfo({ sync: true });
      cachedSystemInfo = legacyInfo;
      return legacyInfo;
    } catch (fallbackError) {
      common_vendor.index.__f__("error", "at utils/uni-compatibility.js:68", "获取系统信息兼容处理失败:", fallbackError);
      return {
        windowWidth: 375,
        windowHeight: 667,
        screenWidth: 375,
        screenHeight: 667,
        statusBarHeight: 20,
        pixelRatio: 2,
        platform: "android",
        language: "zh-Hans"
      };
    }
  }
}
function initUniCompatibility() {
  common_vendor.index.getSystemInfoSync;
  common_vendor.index.getSystemInfoSync = function() {
    common_vendor.index.__f__("warn", "at utils/uni-compatibility.js:94", "wx.getSystemInfoSync 已废弃，请使用 wx.getWindowInfo/wx.getDeviceInfo/wx.getAppBaseInfo 等代替");
    return getCompatibleSystemInfo();
  };
  common_vendor.index.getSystemInfo;
  common_vendor.index.getSystemInfo = function(options) {
    if (options && options.sync === true) {
      if (typeof options.success === "function") {
        options.success(getCompatibleSystemInfo());
      }
      return getCompatibleSystemInfo();
    }
    return new Promise((resolve, reject) => {
      try {
        const info = getCompatibleSystemInfo();
        if (typeof (options == null ? void 0 : options.success) === "function") {
          options.success(info);
        }
        resolve(info);
      } catch (error) {
        if (typeof (options == null ? void 0 : options.fail) === "function") {
          options.fail(error);
        }
        reject(error);
      } finally {
        if (typeof (options == null ? void 0 : options.complete) === "function") {
          options.complete();
        }
      }
    });
  };
  common_vendor.index.__f__("log", "at utils/uni-compatibility.js:130", "uni-app 兼容性适配已完成");
}
exports.initUniCompatibility = initUniCompatibility;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/uni-compatibility.js.map
