"use strict";
const common_vendor = require("../../common/vendor.js");
const utils_request = require("../../utils/request.js");
const mixins_deviceAdapter = require("../../mixins/device-adapter.js");
const _sfc_main = {
  name: "ActivityArea",
  mixins: [mixins_deviceAdapter.deviceAdapter],
  data() {
    return {
      activityList: [],
      isGlobalEnabled: false,
      leftItems: [],
      rightItem: null
    };
  },
  created() {
  },
  methods: {
    async fetchActivityData() {
      try {
        const result = await utils_request.get("/wx/activity-area/get");
        if (result && result.code === 0) {
          this.isGlobalEnabled = result.data.isGlobalEnabled;
          if (this.isGlobalEnabled && result.data.list) {
            this.activityList = result.data.list;
            this.processActivityData();
          }
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at components/activity-area/index.vue:70", "获取活动区域数据失败", error);
      }
    },
    processActivityData() {
      this.leftItems = this.activityList.filter(
        (item) => item.position === "topLeft" || item.position === "bottomLeft"
      );
      this.rightItem = this.activityList.find((item) => item.position === "right");
    },
    getPatternClass(position) {
      if (position === "topLeft")
        return "music";
      if (position === "bottomLeft")
        return "coupon";
      if (position === "right")
        return "video";
      return "";
    },
    handleItemClick(item) {
      if (item.linkType === "page" && item.linkUrl) {
        common_vendor.index.navigateTo({
          url: "/" + item.linkUrl
        });
      } else if (item.linkType === "webview" && item.linkUrl) {
        common_vendor.index.navigateTo({
          url: "/pages/webview/index?url=" + encodeURIComponent(item.linkUrl)
        });
      } else if (item.linkType === "miniprogram" && item.linkUrl) {
        common_vendor.index.navigateToMiniProgram({
          appId: item.linkUrl,
          fail: (err) => {
            common_vendor.index.__f__("error", "at components/activity-area/index.vue:102", "小程序跳转失败", item.linkUrl, err);
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
  }, $data.isGlobalEnabled ? common_vendor.e({
    b: common_vendor.f($data.leftItems, (item, index, i0) => {
      return {
        a: common_vendor.n($options.getPatternClass(item.position)),
        b: common_vendor.t(item.title),
        c: common_vendor.t(item.description),
        d: index,
        e: common_vendor.o(($event) => $options.handleItemClick(item), index)
      };
    }),
    c: _ctx.isSmallScreen ? 1 : "",
    d: _ctx.isSmallScreen ? 1 : "",
    e: $data.rightItem
  }, $data.rightItem ? {
    f: common_vendor.t($data.rightItem.title),
    g: _ctx.isSmallScreen ? 1 : "",
    h: common_vendor.t($data.rightItem.description),
    i: _ctx.isSmallScreen ? 1 : "",
    j: common_vendor.o(($event) => $options.handleItemClick($data.rightItem))
  } : {}, {
    k: _ctx.isSmallScreen ? 1 : "",
    l: _ctx.isLargeScreen ? 1 : ""
  }) : {});
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/components/activity-area/index.js.map
