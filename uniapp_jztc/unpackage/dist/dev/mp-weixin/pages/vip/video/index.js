"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  mixins: [mixins_device.deviceMixin],
  data() {
    return {
      platforms: [
        { name: "爱奇艺", icon: "/static/platform/iqy.png" },
        { name: "腾讯视频", icon: "/static/platform/tx.png" },
        { name: "优酷", icon: "/static/platform/yk.png" },
        { name: "芒果TV", icon: "/static/platform/mg.png" },
        { name: "哔哩哔哩", icon: "/static/platform/bl.png" },
        { name: "搜狐视频", icon: "/static/platform/sh.png" }
      ],
      packages: [
        {
          tag: "特惠专享",
          name: "爱奇艺VIP",
          price: "0.01",
          days: "30"
        },
        {
          tag: "全网低价",
          name: "腾讯视频VIP",
          price: "9.66",
          days: "1"
        },
        {
          tag: "6.1折",
          name: "优酷视频VIP",
          price: "11",
          days: "30"
        }
      ]
    };
  },
  methods: {
    handleBack() {
      common_vendor.index.navigateBack();
    },
    handleWatch() {
      common_vendor.index.showLoading({
        title: "加载中..."
      });
      setTimeout(() => {
        common_vendor.index.hideLoading();
        common_vendor.index.showToast({
          title: "领取成功",
          icon: "success"
        });
      }, 1e3);
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  _easycom_uni_icons2();
}
const _easycom_uni_icons = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.p({
      type: "left",
      size: "20",
      color: "#ffecd8"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: _ctx.statusBarHeight + "px",
    d: common_assets._imports_0$5,
    e: common_vendor.p({
      type: "right",
      size: "14",
      color: "#862c13"
    }),
    f: common_vendor.p({
      type: "right",
      size: "14",
      color: "#862c13"
    }),
    g: common_vendor.f($data.platforms, (item, index, i0) => {
      return {
        a: item.icon,
        b: common_vendor.t(item.name),
        c: index
      };
    }),
    h: common_vendor.f($data.packages, (item, index, i0) => {
      return {
        a: common_vendor.t(item.tag),
        b: common_vendor.t(item.name),
        c: common_vendor.t(item.price),
        d: common_vendor.t(item.days),
        e: index
      };
    }),
    i: common_vendor.o((...args) => $options.handleWatch && $options.handleWatch(...args)),
    j: common_vendor.p({
      type: "info-filled",
      size: "20",
      color: "#ffecd8"
    }),
    k: common_assets._imports_1
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/vip/video/index.js.map
