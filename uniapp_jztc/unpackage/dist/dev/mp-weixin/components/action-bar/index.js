"use strict";
const common_vendor = require("../../common/vendor.js");
const _sfc_main = {
  name: "ActionBar",
  props: {
    // 是否已收藏
    isCollected: {
      type: Boolean,
      default: false
    },
    // 发布者信息
    publisher: {
      type: Object,
      default: () => ({
        id: "",
        name: ""
      })
    },
    // 是否显示留言按钮
    showComment: {
      type: Boolean,
      default: true
    },
    // 按钮文本
    messageText: {
      type: String,
      default: "私信"
    }
  },
  data() {
    return {};
  },
  methods: {
    // 留言
    handleComment() {
      this.$emit("comment");
    },
    // 收藏
    handleCollect() {
      this.$emit("collect");
    },
    // 私信
    handleMessage() {
      if (!this.publisher.id) {
        common_vendor.index.showToast({
          title: "无法获取发布者信息",
          icon: "none"
        });
        return;
      }
      this.$emit("message", this.publisher);
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
  return {
    a: common_vendor.p({
      type: "chat",
      size: "20",
      color: "#666666"
    }),
    b: common_vendor.o((...args) => $options.handleComment && $options.handleComment(...args)),
    c: common_vendor.p({
      type: $props.isCollected ? "star-filled" : "star",
      size: "20",
      color: $props.isCollected ? "#fc3e2b" : "#666666"
    }),
    d: $props.isCollected ? 1 : "",
    e: common_vendor.o((...args) => $options.handleCollect && $options.handleCollect(...args)),
    f: common_vendor.p({
      type: "redo",
      size: "20",
      color: "#666666"
    }),
    g: common_vendor.p({
      type: "chat",
      size: "20",
      color: "#FFFFFF"
    }),
    h: common_vendor.o((...args) => $options.handleMessage && $options.handleMessage(...args))
  };
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/components/action-bar/index.js.map
