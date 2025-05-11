"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const apis_index = require("../../../apis/index.js");
const utils_pay = require("../../../utils/pay.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  mixins: [mixins_device.deviceMixin],
  data() {
    return {
      orderId: null,
      orderNo: null,
      loading: true,
      orderDetail: null
    };
  },
  onLoad(options) {
    if (options.id) {
      this.orderId = options.id;
      this.loadOrderDetail();
    } else if (options.orderNo) {
      this.orderNo = options.orderNo;
      this.loadOrderDetail();
    } else {
      this.loading = false;
      common_vendor.index.showToast({
        title: "订单参数不存在",
        icon: "none"
      });
    }
  },
  methods: {
    // 返回上一页
    handleBack() {
      const pages = getCurrentPages();
      if (pages.length > 1) {
        common_vendor.index.navigateBack();
      } else {
        common_vendor.index.redirectTo({
          url: "/pages/my/orders/index"
        });
      }
    },
    // 加载订单详情
    async loadOrderDetail() {
      let loadingShown = false;
      try {
        this.loading = true;
        common_vendor.index.showLoading({
          title: "加载中...",
          mask: true
        });
        loadingShown = true;
        const params = this.orderNo ? { orderNo: this.orderNo } : { id: this.orderId };
        const res = await apis_index.user.getOrderDetail(params);
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        if (res.code === 0 && res.data) {
          this.orderDetail = res.data;
          this.loading = false;
        } else {
          throw new Error(res.message || "获取订单详情失败");
        }
      } catch (error) {
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        common_vendor.index.__f__("error", "at pages/my/orders/detail.vue:191", "获取订单详情失败:", error);
        common_vendor.index.showToast({
          title: "获取订单失败",
          icon: "none"
        });
        this.loading = false;
      } finally {
        if (loadingShown) {
          common_vendor.index.hideLoading();
        }
      }
    },
    // 根据订单状态获取CSS类名
    getStatusClass(status) {
      switch (status) {
        case 0:
          return "unpaid";
        case 1:
          return "paid";
        case 2:
          return "cancelled";
        case 3:
          return "refunded";
        case 4:
          return "processing";
        case 5:
          return "completed";
        default:
          return "";
      }
    },
    // 根据订单状态获取描述文本
    getStatusDesc(status) {
      switch (status) {
        case 0:
          return "请及时完成支付";
        case 1:
          return "您的订单已支付";
        case 2:
          return "您的订单已取消";
        case 3:
          return "您的订单已退款";
        case 4:
          return "您的订单正在进行中";
        case 5:
          return "您的订单已完成";
        default:
          return "";
      }
    },
    // 获取支付方式文本
    getPaymentMethodText(method) {
      switch (method) {
        case "wechat":
          return "微信支付";
        case "alipay":
          return "支付宝";
        default:
          return "未知";
      }
    },
    // 复制订单号
    copyOrderNo() {
      if (!this.orderDetail || !this.orderDetail.orderNo)
        return;
      common_vendor.index.setClipboardData({
        data: this.orderDetail.orderNo,
        success: () => {
          common_vendor.index.showToast({
            title: "复制成功",
            icon: "success"
          });
        }
      });
    },
    // 取消订单
    handleCancelOrder() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要取消该订单吗？",
        success: async (res) => {
          if (res.confirm) {
            let loadingShown = false;
            try {
              common_vendor.index.showLoading({
                title: "处理中...",
                mask: true
              });
              loadingShown = true;
              await apis_index.user.cancelOrder({
                orderNo: this.orderDetail.orderNo,
                reason: "用户主动取消"
              });
              if (loadingShown) {
                common_vendor.index.hideLoading();
                loadingShown = false;
              }
              common_vendor.index.showToast({
                title: "已取消订单",
                icon: "success"
              });
              this.loadOrderDetail();
            } catch (error) {
              if (loadingShown) {
                common_vendor.index.hideLoading();
                loadingShown = false;
              }
              common_vendor.index.__f__("error", "at pages/my/orders/detail.vue:299", "取消订单失败:", error);
              common_vendor.index.showToast({
                title: error.message || "取消失败",
                icon: "none"
              });
            } finally {
              if (loadingShown) {
                common_vendor.index.hideLoading();
              }
            }
          }
        }
      });
    },
    // 支付订单
    async handlePayOrder() {
      let loadingShown = false;
      try {
        common_vendor.index.showLoading({
          title: "正在支付...",
          mask: true
        });
        loadingShown = true;
        const payResult = await utils_pay.requestOrderPay({
          orderNo: this.orderDetail.orderNo
        });
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        if (payResult.success) {
          common_vendor.index.showToast({
            title: "支付成功",
            icon: "success"
          });
          this.loadOrderDetail();
        } else {
          common_vendor.index.showToast({
            title: payResult.message || "支付已取消",
            icon: "none"
          });
        }
      } catch (error) {
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        common_vendor.index.__f__("error", "at pages/my/orders/detail.vue:360", "支付过程出错:", error);
        common_vendor.index.showToast({
          title: error.message || "支付失败",
          icon: "none"
        });
      } finally {
        if (loadingShown) {
          common_vendor.index.hideLoading();
        }
      }
    },
    // 联系客服
    handleContactService() {
      common_vendor.index.__f__("log", "at pages/my/orders/detail.vue:376", "联系客服");
    },
    // 再次购买
    handleRebuyOrder() {
      common_vendor.index.__f__("log", "at pages/my/orders/detail.vue:382", "再次购买");
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
  return common_vendor.e({
    a: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: _ctx.navBarHeight + "px",
    d: _ctx.statusBarHeight + "px",
    e: `calc(${_ctx.statusBarHeight}px + ${_ctx.navBarHeight}px)`,
    f: $data.loading
  }, $data.loading ? {
    g: common_vendor.p({
      type: "spinner-cycle",
      size: "30",
      color: "#999999"
    })
  } : $data.orderDetail ? common_vendor.e({
    i: common_vendor.t($data.orderDetail.statusText),
    j: common_vendor.n($options.getStatusClass($data.orderDetail.status)),
    k: common_vendor.t($options.getStatusDesc($data.orderDetail.status)),
    l: common_vendor.p({
      type: "shop",
      size: "30",
      color: "#fc3e2b"
    }),
    m: common_vendor.t($data.orderDetail.productName),
    n: common_vendor.t($data.orderDetail.amount.toFixed(2)),
    o: common_vendor.t($data.orderDetail.orderNo),
    p: common_vendor.o((...args) => $options.copyOrderNo && $options.copyOrderNo(...args)),
    q: common_vendor.t($data.orderDetail.createdAt),
    r: $data.orderDetail.payTime
  }, $data.orderDetail.payTime ? {
    s: common_vendor.t($data.orderDetail.payTime)
  } : {}, {
    t: common_vendor.t($options.getPaymentMethodText($data.orderDetail.paymentMethod)),
    v: $data.orderDetail.remark
  }, $data.orderDetail.remark ? {
    w: common_vendor.t($data.orderDetail.remark)
  } : {}, {
    x: common_vendor.t($data.orderDetail.amount.toFixed(2)),
    y: $data.orderDetail.status === 0
  }, $data.orderDetail.status === 0 ? {
    z: common_vendor.o((...args) => $options.handleCancelOrder && $options.handleCancelOrder(...args)),
    A: common_vendor.o((...args) => $options.handlePayOrder && $options.handlePayOrder(...args))
  } : $data.orderDetail.status === 1 || $data.orderDetail.status === 4 ? {
    C: common_vendor.o((...args) => $options.handleContactService && $options.handleContactService(...args))
  } : $data.orderDetail.status === 2 || $data.orderDetail.status === 3 || $data.orderDetail.status === 5 ? {
    E: common_vendor.o((...args) => $options.handleRebuyOrder && $options.handleRebuyOrder(...args))
  } : {}, {
    B: $data.orderDetail.status === 1 || $data.orderDetail.status === 4,
    D: $data.orderDetail.status === 2 || $data.orderDetail.status === 3 || $data.orderDetail.status === 5
  }) : {
    F: common_assets._imports_0$1,
    G: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args))
  }, {
    h: $data.orderDetail
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/my/orders/detail.js.map
