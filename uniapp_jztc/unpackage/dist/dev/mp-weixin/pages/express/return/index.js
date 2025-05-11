"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_deviceAdapter = require("../../../mixins/device-adapter.js");
const _sfc_main = {
  mixins: [mixins_deviceAdapter.deviceAdapter],
  data() {
    return {
      pageTitle: "网购退货",
      senderInfo: {
        name: "",
        phone: "",
        address: ""
      },
      receiverInfo: {
        name: "",
        phone: "",
        address: ""
      },
      currentExpress: -1,
      expressList: [
        {
          logo: "/static/express/sf.png",
          name: "顺丰快递",
          originalPrice: "18.0",
          currentPrice: "15.0",
          weight: "5"
        },
        {
          logo: "/static/express/zt.png",
          name: "中通快递",
          originalPrice: "15.0",
          currentPrice: "12.0",
          weight: "4"
        },
        {
          logo: "/static/express/yd.png",
          name: "韵达快递",
          originalPrice: "15.0",
          currentPrice: "12.0",
          weight: "4"
        }
        // ... 可以添加更多快递公司
      ],
      hasInsurance: false,
      // 是否选择保价服务
      isAgreed: false
      // 是否同意协议
    };
  },
  computed: {
    // 使用统一的布局尺寸
    pageStyle() {
      return {
        "--nav-height": `${this.layoutSize.navHeight}px`,
        "--settlement-height": `${this.layoutSize.settlementHeight}rpx`,
        "--content-padding": `${this.layoutSize.contentPadding}rpx`,
        "--card-gap": `${this.layoutSize.cardGap}rpx`,
        "--content-bottom": `${this.layoutSize.contentBottom}rpx`
      };
    }
  },
  methods: {
    handleBack() {
      common_vendor.index.navigateBack({
        delta: 1,
        fail: () => {
          common_vendor.index.switchTab({
            url: "/pages/index/index"
          });
        }
      });
    },
    handleSelectAddress(type) {
      common_vendor.index.navigateTo({
        url: `/pages/address/list?type=${type}`
      });
    },
    handleSelectGoods() {
      common_vendor.index.navigateTo({
        url: "/pages/express/goods/select"
      });
    },
    handleSelectPayment() {
      common_vendor.index.navigateTo({
        url: "/pages/express/payment/select"
      });
    },
    selectExpress(index) {
      this.currentExpress = index;
    },
    handleSelectTime() {
      common_vendor.index.navigateTo({
        url: "/pages/express/time/select"
      });
    },
    handleInsurance() {
      common_vendor.index.navigateTo({
        url: "/pages/express/insurance/set"
      });
    },
    handleShowDetail() {
      common_vendor.index.showModal({
        title: "费用明细",
        content: "快递费：¥15.0\n保价费：¥0.0\n总计：¥15.0",
        showCancel: false
      });
    },
    handleReadAgreement() {
      common_vendor.index.navigateTo({
        url: "/pages/agreement/express"
      });
    },
    toggleAgreement() {
      this.isAgreed = !this.isAgreed;
    },
    handleSubmitOrder() {
      if (!this.isAgreed) {
        common_vendor.index.showToast({
          title: "请先阅读并同意服务协议",
          icon: "none"
        });
        return;
      }
      common_vendor.index.showLoading({
        title: "提交中..."
      });
    },
    handleUploadImage() {
      common_vendor.index.chooseImage({
        count: 1,
        success: (res) => {
          common_vendor.index.__f__("log", "at pages/express/return/index.vue:349", "上传退货截图:", res.tempFilePaths[0]);
        }
      });
    },
    handlePasteAddress() {
      common_vendor.index.getClipboardData({
        success: (res) => {
          if (res.data) {
            common_vendor.index.__f__("log", "at pages/express/return/index.vue:359", "剪贴板内容:", res.data);
            common_vendor.index.showToast({
              title: "已识别地址",
              icon: "success"
            });
          } else {
            common_vendor.index.showToast({
              title: "剪贴板为空",
              icon: "none"
            });
          }
        }
      });
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
    a: _ctx.navigationBarHeight + "px",
    b: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    c: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    d: common_vendor.t($data.pageTitle),
    e: _ctx.statusBarHeight + "px",
    f: common_vendor.p({
      type: "right",
      size: "14",
      color: "#CCCCCC"
    }),
    g: common_vendor.p({
      type: "right",
      size: "14",
      color: "#CCCCCC"
    }),
    h: common_vendor.o(($event) => $options.handleSelectAddress("send")),
    i: common_vendor.p({
      type: "camera-filled",
      size: "16",
      color: "#007AFF"
    }),
    j: common_vendor.o((...args) => $options.handleUploadImage && $options.handleUploadImage(...args)),
    k: common_vendor.o(($event) => $options.handleSelectAddress("receive")),
    l: common_vendor.o((...args) => $options.handlePasteAddress && $options.handlePasteAddress(...args)),
    m: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    n: common_vendor.o((...args) => $options.handleSelectGoods && $options.handleSelectGoods(...args)),
    o: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    p: common_vendor.o((...args) => $options.handleSelectPayment && $options.handleSelectPayment(...args)),
    q: common_vendor.f($data.expressList, (item, index, i0) => {
      return {
        a: item.logo,
        b: common_vendor.t(item.name),
        c: common_vendor.t(item.currentPrice),
        d: common_vendor.t(item.originalPrice),
        e: common_vendor.t(item.weight),
        f: index,
        g: $data.currentExpress === index ? 1 : "",
        h: common_vendor.o(($event) => $options.selectExpress(index), index)
      };
    }),
    r: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    s: common_vendor.o((...args) => $options.handleSelectTime && $options.handleSelectTime(...args)),
    t: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    v: common_vendor.o((...args) => $options.handleInsurance && $options.handleInsurance(...args)),
    w: _ctx.navigationBarHeight + "px",
    x: common_vendor.p({
      type: "bottom",
      size: "12",
      color: "#666666"
    }),
    y: common_vendor.o((...args) => $options.handleShowDetail && $options.handleShowDetail(...args)),
    z: common_vendor.p({
      type: $data.isAgreed ? "checkbox-filled" : "circle",
      size: "18",
      color: $data.isAgreed ? "#007AFF" : "#CCCCCC"
    }),
    A: common_vendor.o((...args) => $options.toggleAgreement && $options.toggleAgreement(...args)),
    B: common_vendor.o((...args) => $options.handleReadAgreement && $options.handleReadAgreement(...args)),
    C: common_vendor.o((...args) => $options.handleSubmitOrder && $options.handleSubmitOrder(...args)),
    D: common_vendor.s($options.pageStyle)
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/express/return/index.js.map
