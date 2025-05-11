"use strict";
const common_vendor = require("../../common/vendor.js");
const utils_deviceInfo = require("../../utils/device-info.js");
const TabBar = () => "../../components/tab-bar/index.js";
const _sfc_main = {
  components: {
    TabBar
  },
  data() {
    return {
      tabIndex: -1,
      statusBarHeight: 0,
      navBarHeight: 44,
      categoryList: [
        { name: "闲置发布", type: "idle" },
        { name: "信息发布", type: "info" }
      ]
    };
  },
  onLoad() {
    this.statusBarHeight = utils_deviceInfo.deviceInfo.getStatusBarHeight();
  },
  onShow() {
    this.tabIndex = 2;
  },
  methods: {
    handleCategoryClick(category) {
      common_vendor.index.navigateTo({
        url: `/pages/publish/${category.type}/index`
      });
    },
    handleDraftClick() {
      common_vendor.index.navigateTo({
        url: "/pages/publish/draft/index"
      });
    }
  }
};
if (!Array) {
  const _component_tab_bar = common_vendor.resolveComponent("tab-bar");
  _component_tab_bar();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.o((...args) => $options.handleDraftClick && $options.handleDraftClick(...args)),
    b: $data.navBarHeight + "px",
    c: $data.statusBarHeight + "px",
    d: `calc(${$data.statusBarHeight}px + ${$data.navBarHeight}px)`,
    e: common_vendor.t($data.categoryList[0].name),
    f: common_vendor.o(($event) => $options.handleCategoryClick($data.categoryList[0])),
    g: common_vendor.t($data.categoryList[1].name),
    h: common_vendor.o(($event) => $options.handleCategoryClick($data.categoryList[1])),
    i: common_vendor.p({
      ["current-tab"]: $data.tabIndex
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/publish/index.js.map
