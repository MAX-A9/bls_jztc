"use strict";
const common_vendor = require("../../../../../common/vendor.js");
const _sfc_main = {
  name: "uniSwipeAction",
  data() {
    return {};
  },
  created() {
    this.children = [];
  },
  methods: {
    // 公开给用户使用，重制组件样式
    resize() {
    },
    // 公开给用户使用，关闭全部 已经打开的组件
    closeAll() {
      this.children.forEach((vm) => {
        vm.is_show = "none";
      });
    },
    closeOther(vm) {
      if (this.openItem && this.openItem !== vm) {
        this.openItem.is_show = "none";
      }
      this.openItem = vm;
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {};
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../../../../.sourcemap/mp-weixin/node-modules/@dcloudio/uni-ui/lib/uni-swipe-action/uni-swipe-action.js.map
