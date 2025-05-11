"use strict";
const common_vendor = require("../../common/vendor.js");
const mixins_deviceAdapter = require("../../mixins/device-adapter.js");
const _sfc_main = {
  mixins: [mixins_deviceAdapter.deviceAdapter],
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
  data() {
    return {
      pageTitle: "安装维修",
      isAgreed: false,
      // 是否同意协议
      currentService: "",
      // 当前选中的服务类型
      demandText: "",
      // 安装需求输入框的值
      keywords: ["柜子", "空调", "洗衣机", "电视", "床", "沙发", "衣柜", "冰箱", "热水器", "油烟机"],
      // 关键词列表
      selectedAddress: null,
      // 选中的地址
      selectedTime: "",
      // 选中的预约时间
      imageList: [],
      // 物品图片列表
      noteText: ""
      // 备注信息输入框的值
    };
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
    handleShowDetail() {
      common_vendor.index.showModal({
        title: "费用明细",
        content: "服务费：¥15.0\n总计：¥15.0",
        showCancel: false
      });
    },
    handleReadAgreement() {
      common_vendor.index.navigateTo({
        url: "/pages/agreement/service"
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
    selectService(type) {
      this.currentService = type;
    },
    selectKeyword(keyword) {
      if (this.demandText) {
        this.demandText = this.demandText.trim() + "，" + keyword;
      } else {
        this.demandText = keyword;
      }
    },
    handleSelectAddress() {
      common_vendor.index.navigateTo({
        url: "/pages/address/list?type=service"
      });
    },
    handleSelectTime() {
      common_vendor.index.navigateTo({
        url: "/pages/service/time"
      });
    },
    previewImage(index) {
      common_vendor.index.previewImage({
        urls: this.imageList,
        current: index
      });
    },
    deleteImage(index) {
      this.imageList.splice(index, 1);
    },
    chooseImage() {
      common_vendor.index.chooseImage({
        count: 4 - this.imageList.length,
        sizeType: ["compressed"],
        sourceType: ["album", "camera"],
        success: (res) => {
          this.imageList = [...this.imageList, ...res.tempFilePaths];
        }
      });
    },
    handleCheckPrice() {
      common_vendor.index.navigateTo({
        url: "/pages/service/nearby-price"
      });
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
      type: "tools",
      size: "32",
      color: $data.currentService === "install" ? "#007AFF" : "#666666"
    }),
    g: $data.currentService === "install" ? 1 : "",
    h: common_vendor.o(($event) => $options.selectService("install")),
    i: common_vendor.p({
      type: "refresh",
      size: "32",
      color: $data.currentService === "repair" ? "#007AFF" : "#666666"
    }),
    j: $data.currentService === "repair" ? 1 : "",
    k: common_vendor.o(($event) => $options.selectService("repair")),
    l: common_vendor.p({
      type: "clear",
      size: "32",
      color: $data.currentService === "clean" ? "#007AFF" : "#666666"
    }),
    m: $data.currentService === "clean" ? 1 : "",
    n: common_vendor.o(($event) => $options.selectService("clean")),
    o: $data.demandText,
    p: common_vendor.o(($event) => $data.demandText = $event.detail.value),
    q: common_vendor.t($data.demandText.length),
    r: common_vendor.f($data.keywords, (item, index, i0) => {
      return {
        a: common_vendor.t(item),
        b: index,
        c: common_vendor.o(($event) => $options.selectKeyword(item), index)
      };
    }),
    s: !$data.selectedAddress
  }, !$data.selectedAddress ? {
    t: common_vendor.p({
      type: "location",
      size: "20",
      color: "#666666"
    })
  } : {
    v: common_vendor.t($data.selectedAddress.name),
    w: common_vendor.t($data.selectedAddress.phone),
    x: common_vendor.t($data.selectedAddress.address)
  }, {
    y: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    z: common_vendor.o((...args) => $options.handleSelectAddress && $options.handleSelectAddress(...args)),
    A: !$data.selectedTime
  }, !$data.selectedTime ? {
    B: common_vendor.p({
      type: "calendar",
      size: "20",
      color: "#666666"
    })
  } : {
    C: common_vendor.t($data.selectedTime)
  }, {
    D: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    E: common_vendor.o((...args) => $options.handleSelectTime && $options.handleSelectTime(...args)),
    F: common_vendor.f($data.imageList, (item, index, i0) => {
      return {
        a: item,
        b: "6007866c-8-" + i0,
        c: common_vendor.o(($event) => $options.deleteImage(index), index),
        d: index,
        e: common_vendor.o(($event) => $options.previewImage(index), index)
      };
    }),
    G: common_vendor.p({
      type: "clear",
      size: "12",
      color: "#FFFFFF"
    }),
    H: $data.imageList.length < 4
  }, $data.imageList.length < 4 ? {
    I: common_vendor.p({
      type: "camera",
      size: "24",
      color: "#666666"
    }),
    J: common_vendor.o((...args) => $options.chooseImage && $options.chooseImage(...args))
  } : {}, {
    K: $data.noteText,
    L: common_vendor.o(($event) => $data.noteText = $event.detail.value),
    M: common_vendor.t($data.noteText.length),
    N: _ctx.navigationBarHeight + "px",
    O: common_vendor.p({
      type: "right",
      size: "12",
      color: "#E6F0FF"
    }),
    P: common_vendor.o((...args) => $options.handleSubmitOrder && $options.handleSubmitOrder(...args)),
    Q: common_vendor.s($options.pageStyle)
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/service/index.js.map
