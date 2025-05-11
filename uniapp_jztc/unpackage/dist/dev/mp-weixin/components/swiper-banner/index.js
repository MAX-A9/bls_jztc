"use strict";
const common_vendor = require("../../common/vendor.js");
const utils_request = require("../../utils/request.js");
const _sfc_main = {
  name: "SwiperBanner",
  data() {
    return {
      current: 0,
      bannerList: [],
      isGlobalEnabled: false
    };
  },
  created() {
  },
  methods: {
    async fetchBannerData() {
      try {
        const result = await utils_request.get("/wx/banner/list");
        if (result && result.code === 0) {
          this.isGlobalEnabled = result.data.isGlobalEnabled;
          if (this.isGlobalEnabled && result.data.list) {
            this.bannerList = result.data.list.filter((item) => item.isEnabled).sort((a, b) => a.order - b.order);
          }
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at components/swiper-banner/index.vue:58", "获取轮播图数据失败", error);
      }
    },
    handleChange(e) {
      this.current = e.detail.current;
    },
    handleClick(item) {
      if (!item || !item.linkUrl)
        return;
      if (item.linkType === "page") {
        common_vendor.index.navigateTo({
          url: "/" + item.linkUrl,
          fail: (err) => {
            common_vendor.index.__f__("error", "at components/swiper-banner/index.vue:71", "页面跳转失败", item.linkUrl, err);
          }
        });
      } else if (item.linkType === "webview") {
        common_vendor.index.navigateTo({
          url: "/pages/webview/index?url=" + encodeURIComponent(item.linkUrl),
          fail: (err) => {
            common_vendor.index.__f__("error", "at components/swiper-banner/index.vue:78", "网页跳转失败", item.linkUrl, err);
          }
        });
      } else if (item.linkType === "miniprogram") {
        common_vendor.index.navigateToMiniProgram({
          appId: item.linkUrl,
          fail: (err) => {
            common_vendor.index.__f__("error", "at components/swiper-banner/index.vue:86", "小程序跳转失败", item.linkUrl, err);
            common_vendor.index.showToast({
              title: "跳转失败",
              icon: "none"
            });
          }
        });
      }
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.isGlobalEnabled
  }, $data.isGlobalEnabled ? {
    b: common_vendor.f($data.bannerList, (item, index, i0) => {
      return {
        a: item.image,
        b: index,
        c: common_vendor.o(($event) => $options.handleClick(item), index)
      };
    }),
    c: common_vendor.o((...args) => $options.handleChange && $options.handleChange(...args)),
    d: common_vendor.f($data.bannerList, (item, index, i0) => {
      return {
        a: index,
        b: $data.current === index ? 1 : ""
      };
    })
  } : {});
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/components/swiper-banner/index.js.map
