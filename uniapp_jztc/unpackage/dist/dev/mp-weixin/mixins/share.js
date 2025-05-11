"use strict";
const common_vendor = require("../common/vendor.js");
const utils_share = require("../utils/share.js");
const shareMixin = {
  data() {
    return {
      // 分享相关数据
      shareData: {
        title: "",
        path: "",
        imageUrl: ""
      }
    };
  },
  // 页面生命周期钩子
  onLoad() {
    this.initShareData();
  },
  // 页面方法
  methods: {
    /**
     * 初始化分享数据
     */
    async initShareData() {
      try {
        if (this.isContentPage) {
          const contentData = this.contentData || this.detail || {};
          const shareOptions = await utils_share.getContentShareOptions(contentData);
          this.shareData = { ...shareOptions };
        } else if (this.isHomePage) {
          const shareOptions = await utils_share.getHomeShareOptions();
          this.shareData = { ...shareOptions };
        } else {
          const pages = getCurrentPages();
          const currentPage = pages[pages.length - 1];
          const path = `/${currentPage.route}${this._getPageQuery(currentPage)}`;
          const shareOptions = await utils_share.getDefaultShareOptions(path);
          this.shareData = { ...shareOptions };
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at mixins/share.js:58", "初始化分享数据失败:", error);
      }
    },
    /**
     * 获取页面查询参数
     * @param {Object} page 页面对象
     * @returns {String} 查询参数字符串
     */
    _getPageQuery(page) {
      if (!page.options || Object.keys(page.options).length === 0) {
        return "";
      }
      const query = Object.keys(page.options).map((key) => `${key}=${encodeURIComponent(page.options[key])}`).join("&");
      return `?${query}`;
    }
  },
  // 分享给朋友
  onShareAppMessage() {
    return {
      title: this.shareData.title,
      path: this.shareData.path,
      imageUrl: this.shareData.imageUrl
    };
  },
  // 分享到朋友圈
  onShareTimeline() {
    return {
      title: this.shareData.title,
      query: this.shareData.path.indexOf("?") > -1 ? this.shareData.path.split("?")[1] : "",
      imageUrl: this.shareData.imageUrl
    };
  }
};
exports.shareMixin = shareMixin;
//# sourceMappingURL=../../.sourcemap/mp-weixin/mixins/share.js.map
