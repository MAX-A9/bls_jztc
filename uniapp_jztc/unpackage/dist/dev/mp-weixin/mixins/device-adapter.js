"use strict";
const common_vendor = require("../common/vendor.js");
const utils_deviceInfo = require("../utils/device-info.js");
const deviceAdapter = {
  data() {
    return {
      // 缓存系统信息，避免重复获取
      _windowInfo: null,
      _deviceInfo: null,
      _appBaseInfo: null,
      _systemSetting: null
    };
  },
  created() {
    try {
      this._windowInfo = utils_deviceInfo.deviceInfo.getWindowInfo();
      this._deviceInfo = utils_deviceInfo.deviceInfo.getDeviceInfo();
      this._appBaseInfo = utils_deviceInfo.deviceInfo.getAppBaseInfo();
      this._systemSetting = utils_deviceInfo.deviceInfo.getSystemSetting();
    } catch (e) {
      common_vendor.index.__f__("error", "at mixins/device-adapter.js:22", "获取设备信息失败:", e);
      this._windowInfo = {
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
        }
      };
      this._deviceInfo = {
        brand: "",
        model: "",
        system: "",
        platform: ""
      };
      this._appBaseInfo = {
        SDKVersion: "",
        appName: "",
        appVersion: "",
        appLanguage: "zh-Hans"
      };
    }
  },
  computed: {
    // 获取状态栏高度
    statusBarHeight() {
      var _a;
      return ((_a = this._windowInfo) == null ? void 0 : _a.statusBarHeight) || 0;
    },
    // 导航栏总高度
    navigationBarHeight() {
      var _a;
      const isIOS = (((_a = this._deviceInfo) == null ? void 0 : _a.system) || "").toLowerCase().includes("ios");
      const navHeight = isIOS ? 44 : 48;
      return this.statusBarHeight + navHeight;
    },
    // 是否是小屏幕设备
    isSmallScreen() {
      var _a;
      return (((_a = this._windowInfo) == null ? void 0 : _a.windowWidth) || 0) <= 320;
    },
    // 是否是大屏幕设备
    isLargeScreen() {
      var _a;
      return (((_a = this._windowInfo) == null ? void 0 : _a.windowWidth) || 0) >= 768;
    },
    // 安全区域
    safeAreaInsets() {
      var _a;
      const safeArea = ((_a = this._windowInfo) == null ? void 0 : _a.safeArea) || {};
      return {
        top: safeArea.top || 0,
        bottom: safeArea.bottom || 0,
        left: safeArea.left || 0,
        right: safeArea.right || 0
      };
    },
    // 屏幕尺寸
    screenSize() {
      var _a, _b;
      return {
        width: ((_a = this._windowInfo) == null ? void 0 : _a.windowWidth) || 0,
        height: ((_b = this._windowInfo) == null ? void 0 : _b.windowHeight) || 0
      };
    },
    // 统一布局尺寸
    layoutSize() {
      return {
        // 底部结算栏高度(不含安全区域)
        settlementHeight: 120,
        // 固定120rpx
        // 内容区域底部间距
        contentBottom: 160,
        // 固定160rpx
        // 导航栏高度
        navHeight: 44,
        // 固定44px
        // 卡片间距
        cardGap: 20,
        // 固定20rpx
        // 内容区域水平内边距
        contentPadding: 20
        // 固定20rpx
      };
    }
  },
  methods: {
    // 获取元素布局信息
    async getElementRect(selector) {
      return new Promise((resolve) => {
        common_vendor.index.createSelectorQuery().select(selector).boundingClientRect((data) => {
          resolve(data);
        }).exec();
      });
    },
    // rpx 转 px
    rpxToPx(rpx) {
      var _a;
      const screenWidth = ((_a = this._windowInfo) == null ? void 0 : _a.windowWidth) || 375;
      return rpx / 750 * screenWidth;
    },
    // px 转 rpx 
    pxToRpx(px) {
      var _a;
      const screenWidth = ((_a = this._windowInfo) == null ? void 0 : _a.windowWidth) || 375;
      return px * 750 / screenWidth;
    }
  }
};
exports.deviceAdapter = deviceAdapter;
//# sourceMappingURL=../../.sourcemap/mp-weixin/mixins/device-adapter.js.map
