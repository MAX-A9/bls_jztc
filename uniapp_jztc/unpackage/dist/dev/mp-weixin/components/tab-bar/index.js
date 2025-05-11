"use strict";
const common_vendor = require("../../common/vendor.js");
const mixins_device = require("../../mixins/device.js");
const utils_messagePolling = require("../../utils/message-polling.js");
const common_assets = require("../../common/assets.js");
const _sfc_main = {
  name: "TabBar",
  mixins: [mixins_device.deviceMixin],
  props: {
    currentTab: {
      type: Number,
      default: 0
    }
  },
  data() {
    return {
      unreadCount: 0,
      unreadCountUnsubscribe: null
      // 用于存储取消订阅函数
    };
  },
  created() {
    this.unreadCountUnsubscribe = utils_messagePolling.messagePollingService.addUnreadCountListener(
      (unreadCount) => {
        this.unreadCount = unreadCount;
      }
    );
  },
  beforeDestroy() {
    if (this.unreadCountUnsubscribe) {
      this.unreadCountUnsubscribe();
      this.unreadCountUnsubscribe = null;
    }
  },
  methods: {
    switchTab(index) {
      const routes = [
        "/pages/index/index",
        "/pages/community/index",
        "/pages/publish/index",
        "/pages/message/index",
        "/pages/my/index"
      ];
      if (index >= 0 && index < routes.length) {
        common_vendor.index.switchTab({
          url: routes[index]
        });
      }
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: _ctx.safeAreaInsets.bottom + "px",
    b: $props.currentTab === 0 ? "/static/tabbar/sy-liang.png" : "/static/tabbar/sy.png",
    c: $props.currentTab === 0 ? 1 : "",
    d: common_vendor.o(($event) => $options.switchTab(0)),
    e: $props.currentTab === 1 ? "/static/tabbar/sq-liang.png" : "/static/tabbar/sq.png",
    f: $props.currentTab === 1 ? 1 : "",
    g: common_vendor.o(($event) => $options.switchTab(1)),
    h: common_assets._imports_0$8,
    i: common_vendor.o(($event) => $options.switchTab(2)),
    j: $props.currentTab === 3 ? "/static/tabbar/xiaoxi-liang.png" : "/static/tabbar/xiaoxi.png",
    k: $props.currentTab === 3 ? 1 : "",
    l: $data.unreadCount > 0
  }, $data.unreadCount > 0 ? {
    m: common_vendor.t($data.unreadCount > 99 ? "99+" : $data.unreadCount)
  } : {}, {
    n: common_vendor.o(($event) => $options.switchTab(3)),
    o: $props.currentTab === 4 ? "/static/tabbar/my-liang.png" : "/static/tabbar/my.png",
    p: $props.currentTab === 4 ? 1 : "",
    q: common_vendor.o(($event) => $options.switchTab(4)),
    r: _ctx.safeAreaInsets.bottom + "px"
  });
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/components/tab-bar/index.js.map
