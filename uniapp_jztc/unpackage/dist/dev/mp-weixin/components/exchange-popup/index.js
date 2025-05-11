"use strict";
const common_vendor = require("../../common/vendor.js");
const uniPopup = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-popup/uni-popup.js";
const uniIcons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
const _sfc_main = {
  name: "exchange-popup",
  components: {
    uniPopup,
    uniIcons
  },
  data() {
    return {
      phone: "",
      confirmPhone: "",
      productInfo: null
    };
  },
  mounted() {
    common_vendor.index.__f__("log", "at components/exchange-popup/index.vue:69", "Exchange popup component mounted");
  },
  methods: {
    // 打开弹窗
    open(product) {
      common_vendor.index.__f__("log", "at components/exchange-popup/index.vue:74", "Try to open popup with product:", product);
      this.productInfo = product;
      this.phone = "";
      this.confirmPhone = "";
      this.$nextTick(() => {
        if (this.$refs.popup) {
          this.$refs.popup.open("bottom");
        } else {
          common_vendor.index.__f__("error", "at components/exchange-popup/index.vue:84", "Popup ref not found");
        }
      });
    },
    // 关闭弹窗
    closePopup() {
      this.$refs.popup.close();
    },
    // 确认兑换
    handleConfirm() {
      if (!this.phone) {
        common_vendor.index.showToast({
          title: "请输入手机号",
          icon: "none"
        });
        return;
      }
      if (!/^1\d{10}$/.test(this.phone)) {
        common_vendor.index.showToast({
          title: "请输入正确的手机号",
          icon: "none"
        });
        return;
      }
      if (this.phone !== this.confirmPhone) {
        common_vendor.index.showToast({
          title: "两次输入的手机号不一致",
          icon: "none"
        });
        return;
      }
      const data = {
        duration: this.productInfo.duration || 0,
        productName: this.productInfo.name || "",
        rechargeAccount: this.phone
      };
      this.$emit("confirm", data);
      this.closePopup();
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _easycom_uni_popup2 = common_vendor.resolveComponent("uni-popup");
  (_easycom_uni_icons2 + _easycom_uni_popup2)();
}
const _easycom_uni_icons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
const _easycom_uni_popup = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-popup/uni-popup.js";
if (!Math) {
  (_easycom_uni_icons + _easycom_uni_popup)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.p({
      type: "close",
      size: "20",
      color: "#333"
    }),
    b: common_vendor.o((...args) => $options.closePopup && $options.closePopup(...args)),
    c: $data.productInfo
  }, $data.productInfo ? {
    d: $data.productInfo.thumbnail,
    e: common_vendor.t($data.productInfo.name),
    f: common_vendor.t($data.productInfo.duration)
  } : {}, {
    g: $data.phone,
    h: common_vendor.o(($event) => $data.phone = $event.detail.value),
    i: $data.confirmPhone,
    j: common_vendor.o(($event) => $data.confirmPhone = $event.detail.value),
    k: common_vendor.o((...args) => $options.handleConfirm && $options.handleConfirm(...args)),
    l: common_vendor.sr("popup", "2792bd8a-0"),
    m: common_vendor.p({
      type: "bottom"
    })
  });
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/components/exchange-popup/index.js.map
