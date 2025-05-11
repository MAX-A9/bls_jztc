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
      currentTab: 0,
      isRefreshing: false,
      isLoading: false,
      hasMore: true,
      page: 1,
      pageSize: 10,
      total: 0,
      tabList: [
        { name: "全部", status: "all" },
        { name: "进行中", status: "process" },
        { name: "待支付", status: "unpaid" },
        { name: "已完成", status: "completed" }
      ],
      orderList: []
    };
  },
  onLoad(options) {
    if (options.tab) {
      this.currentTab = parseInt(options.tab);
    }
    this.loadOrderList();
  },
  methods: {
    // 获取订单列表
    async loadOrderList(isRefresh = false) {
      if (this.isLoading)
        return;
      let loadingShown = false;
      try {
        this.isLoading = true;
        if (!isRefresh) {
          common_vendor.index.showLoading({
            title: "加载中...",
            mask: true
          });
          loadingShown = true;
        }
        if (isRefresh) {
          this.page = 1;
          this.orderList = [];
        }
        const status = this.tabList[this.currentTab].status;
        const res = await apis_index.user.getOrderList({
          page: this.page,
          pageSize: this.pageSize,
          status
        });
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        if (res.code === 0 && res.data) {
          if (this.page === 1) {
            this.orderList = res.data.list || [];
          } else {
            this.orderList = [...this.orderList, ...res.data.list || []];
          }
          this.total = res.data.total || 0;
          this.hasMore = this.orderList.length < this.total;
        } else {
          throw new Error(res.message || "获取订单列表失败");
        }
      } catch (error) {
        if (loadingShown) {
          common_vendor.index.hideLoading();
          loadingShown = false;
        }
        common_vendor.index.__f__("error", "at pages/my/orders/index.vue:203", "获取订单列表失败:", error);
        common_vendor.index.showToast({
          title: "获取订单失败",
          icon: "none"
        });
      } finally {
        this.isLoading = false;
        if (this.isRefreshing) {
          this.isRefreshing = false;
        }
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
          return "processing";
        case 2:
          return "completed";
        case 3:
          return "cancelled";
        case 4:
          return "refunded";
        default:
          return "";
      }
    },
    // 返回上一页
    handleBack() {
      const pages = getCurrentPages();
      if (pages.length > 1) {
        common_vendor.index.navigateBack();
      } else {
        common_vendor.index.switchTab({
          url: "/pages/my/index"
        });
      }
    },
    // 切换选项卡
    switchTab(index) {
      if (this.currentTab === index)
        return;
      this.currentTab = index;
      this.page = 1;
      this.hasMore = true;
      this.loadOrderList(true);
    },
    // 下拉刷新
    onRefresh() {
      this.isRefreshing = true;
      this.page = 1;
      this.hasMore = true;
      this.loadOrderList(true);
    },
    // 加载更多
    loadMore() {
      if (!this.hasMore || this.isLoading)
        return;
      this.page++;
      this.loadOrderList();
    },
    // 点击订单
    handleOrderClick(order) {
      common_vendor.index.navigateTo({
        url: `/pages/my/orders/detail?orderNo=${order.orderNo}`
      });
    },
    // 联系客服
    handleContactService(order) {
      common_vendor.index.makePhoneCall({
        phoneNumber: "10086",
        fail: () => {
          common_vendor.index.showToast({
            title: "拨打电话失败",
            icon: "none"
          });
        }
      });
    },
    // 取消订单
    handleCancelOrder(order) {
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
                orderNo: order.orderNo,
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
              this.onRefresh();
            } catch (error) {
              if (loadingShown) {
                common_vendor.index.hideLoading();
                loadingShown = false;
              }
              common_vendor.index.__f__("error", "at pages/my/orders/index.vue:338", "取消订单失败:", error);
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
    async handlePayOrder(order) {
      let loadingShown = false;
      try {
        common_vendor.index.showLoading({
          title: "正在支付...",
          mask: true
          // 使用蒙层防止用户点击
        });
        loadingShown = true;
        const payResult = await utils_pay.requestOrderPay({
          orderNo: order.orderNo
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
          this.onRefresh();
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
        common_vendor.index.__f__("error", "at pages/my/orders/index.vue:399", "支付过程出错:", error);
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
    // 再次购买
    handleRebuyOrder(order) {
      common_vendor.index.navigateTo({
        url: "/pages/publish/info/index"
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
  return common_vendor.e({
    a: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: _ctx.navBarHeight + "px",
    d: common_vendor.f($data.tabList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.name),
        b: index,
        c: $data.currentTab === index ? 1 : "",
        d: common_vendor.o(($event) => $options.switchTab(index), index)
      };
    }),
    e: _ctx.statusBarHeight + "px",
    f: common_vendor.f($data.orderList, (item, k0, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.createdAt),
        b: common_vendor.t(item.statusText),
        c: common_vendor.n($options.getStatusClass(item.status)),
        d: "e6418f06-1-" + i0,
        e: common_vendor.t(item.productName),
        f: common_vendor.t(item.orderNo),
        g: common_vendor.t(item.amount.toFixed(2)),
        h: item.status === 0
      }, item.status === 0 ? {
        i: common_vendor.o(($event) => $options.handleCancelOrder(item), item.id),
        j: common_vendor.o(($event) => $options.handlePayOrder(item), item.id)
      } : item.status === 1 ? {
        l: common_vendor.o(($event) => $options.handleContactService(item), item.id)
      } : item.status === 2 || item.status === 3 ? {
        n: common_vendor.o(($event) => $options.handleRebuyOrder(item), item.id)
      } : {}, {
        k: item.status === 1,
        m: item.status === 2 || item.status === 3,
        o: item.id,
        p: common_vendor.o(($event) => $options.handleOrderClick(item), item.id)
      });
    }),
    g: common_vendor.p({
      type: "shop",
      size: "30",
      color: "#fc3e2b"
    }),
    h: $data.orderList.length > 0 && $data.hasMore
  }, $data.orderList.length > 0 && $data.hasMore ? {} : {}, {
    i: $data.orderList.length > 0 && !$data.hasMore
  }, $data.orderList.length > 0 && !$data.hasMore ? {} : {}, {
    j: $data.orderList.length === 0 && !$data.isLoading
  }, $data.orderList.length === 0 && !$data.isLoading ? {
    k: common_assets._imports_0$1
  } : {}, {
    l: `calc(${_ctx.statusBarHeight}px + ${_ctx.navBarHeight}px + 80rpx)`,
    m: $data.isRefreshing,
    n: common_vendor.o((...args) => $options.onRefresh && $options.onRefresh(...args)),
    o: common_vendor.o((...args) => $options.loadMore && $options.loadMore(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/my/orders/index.js.map
