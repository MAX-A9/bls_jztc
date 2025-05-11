"use strict";
const common_vendor = require("../common/vendor.js");
const utils_deviceInfo = require("../utils/device-info.js");
const deviceMixin = {
  data() {
    return {
      // 设备相关
      isIOS: false,
      isAndroid: false,
      statusBarHeight: 20,
      navBarHeight: 44,
      // 屏幕信息
      screenWidth: 375,
      screenHeight: 667,
      // 安全区域
      safeAreaInsets: {
        top: 0,
        right: 0,
        bottom: 0,
        left: 0
      }
    };
  },
  created() {
    try {
      const windowInfo = utils_deviceInfo.deviceInfo.getWindowInfo();
      const appBaseInfo = utils_deviceInfo.deviceInfo.getAppBaseInfo();
      const deviceData = utils_deviceInfo.deviceInfo.getDeviceInfo();
      this.isIOS = deviceData.platform === "ios" || /ios/i.test(deviceData.system || "");
      this.isAndroid = deviceData.platform === "android" || /android/i.test(deviceData.system || "");
      this.statusBarHeight = windowInfo.statusBarHeight || 20;
      this.navBarHeight = this.isIOS ? 44 : 48;
      this.screenWidth = windowInfo.screenWidth || 375;
      this.screenHeight = windowInfo.screenHeight || 667;
      const safeAreaInsets = utils_deviceInfo.deviceInfo.getSafeAreaInsets();
      if (safeAreaInsets) {
        this.safeAreaInsets = safeAreaInsets;
      }
    } catch (error) {
      common_vendor.index.__f__("error", "at mixins/device.js:54", "初始化设备信息失败:", error);
    }
  },
  computed: {
    // 状态栏+导航栏的总高度
    navigationBarHeight() {
      return this.statusBarHeight + this.navBarHeight;
    },
    // 获取安全区域高度
    safeAreaHeight() {
      return this.screenHeight - this.safeAreaInsets.top - this.safeAreaInsets.bottom;
    },
    // 获取内容区域高度(去除导航栏和安全区域)
    contentHeight() {
      return this.safeAreaHeight - this.navigationBarHeight;
    }
  },
  methods: {
    // rpx 转 px
    rpx2px(rpx) {
      return utils_deviceInfo.deviceInfo.rpx2px(rpx);
    },
    // px 转 rpx
    px2rpx(px) {
      return utils_deviceInfo.deviceInfo.px2rpx(px);
    },
    // 判断是否为 iPhone X 系列
    isIphoneX() {
      return utils_deviceInfo.deviceInfo.isIphoneX();
    },
    // 获取底部安全距离
    getBottomSafeDistance() {
      return this.isIphoneX() ? this.safeAreaInsets.bottom : 0;
    },
    // 获取适配后的安全区域
    getSafeAreaInsets() {
      return {
        top: `${this.safeAreaInsets.top}px`,
        right: `${this.safeAreaInsets.right}px`,
        bottom: `${this.safeAreaInsets.bottom}px`,
        left: `${this.safeAreaInsets.left}px`
      };
    }
  }
};
exports.deviceMixin = deviceMixin;
//# sourceMappingURL=../../.sourcemap/mp-weixin/mixins/device.js.map
