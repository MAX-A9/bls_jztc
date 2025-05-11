"use strict";
const common_vendor = require("../../common/vendor.js");
const apis_index = require("../../apis/index.js");
const utils_deviceInfo = require("../../utils/device-info.js");
const _sfc_main = {
  data() {
    return {
      statusBarHeight: 0,
      navBarHeight: 44,
      title: "协议详情",
      content: "",
      type: "",
      loading: true,
      errorMsg: ""
    };
  },
  onLoad(options) {
    this.type = options.type || "privacy";
    this.statusBarHeight = utils_deviceInfo.deviceInfo.getStatusBarHeight();
    this.loadAgreementContent();
  },
  methods: {
    // 处理返回按钮点击
    handleBack() {
      common_vendor.index.navigateBack();
    },
    // 加载协议内容
    async loadAgreementContent() {
      this.loading = true;
      this.errorMsg = "";
      try {
        const res = await apis_index.agreement.getAgreement({ type: this.type });
        if (res.code === 0 && res.data) {
          this.content = res.data.content || `暂无${this.title}内容`;
          this.loading = false;
        } else {
          throw new Error(res.message || "获取协议内容失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/agreement/index.vue:88", "获取协议内容失败:", error);
        this.errorMsg = error.message || "获取协议内容失败，请重试";
        this.loading = false;
      }
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  _easycom_uni_icons2();
}
const _easycom_uni_icons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: common_vendor.t($data.title),
    d: $data.navBarHeight + "px",
    e: $data.statusBarHeight + "px",
    f: !$data.loading
  }, !$data.loading ? {
    g: $data.content
  } : {}, {
    h: $data.loading
  }, $data.loading ? {
    i: common_vendor.p({
      type: "spinner-cycle",
      size: "30",
      color: "#fc3e2b"
    })
  } : {}, {
    j: $data.errorMsg
  }, $data.errorMsg ? {
    k: common_vendor.p({
      type: "info",
      size: "50",
      color: "#cccccc"
    }),
    l: common_vendor.t($data.errorMsg),
    m: common_vendor.o((...args) => $options.loadAgreementContent && $options.loadAgreementContent(...args))
  } : {}, {
    n: `calc(${$data.statusBarHeight}px + ${$data.navBarHeight}px)`
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/agreement/index.js.map
