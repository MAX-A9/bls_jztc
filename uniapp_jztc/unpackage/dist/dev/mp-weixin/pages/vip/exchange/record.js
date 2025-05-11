"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const apis_index = require("../../../apis/index.js");
const common_assets = require("../../../common/assets.js");
const uniIcons = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
const uniLoadMore = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-load-more/uni-load-more.js";
const _sfc_main = {
  components: {
    uniIcons,
    uniLoadMore
  },
  mixins: [mixins_device.deviceMixin],
  data() {
    return {
      recordList: [],
      page: 1,
      size: 10,
      totalPages: 0,
      isLoading: false,
      loadMoreStatus: "more"
      // more-加载前 loading-加载中 noMore-没有更多
    };
  },
  onLoad() {
    this.getExchangeRecords();
  },
  // 下拉刷新
  onPullDownRefresh() {
    this.page = 1;
    this.recordList = [];
    this.getExchangeRecords().then(() => {
      common_vendor.index.stopPullDownRefresh();
    });
  },
  // 上拉加载更多
  onReachBottom() {
    if (this.page < this.totalPages && this.loadMoreStatus !== "loading") {
      this.page++;
      this.loadMoreStatus = "loading";
      this.getExchangeRecords();
    }
  },
  methods: {
    handleBack() {
      common_vendor.index.navigateBack();
    },
    // 格式化状态为中文
    formatStatus(status) {
      if (!status)
        return "未知状态";
      const statusLower = typeof status === "string" ? status.toLowerCase() : status;
      if (statusLower === "processing" || statusLower === "处理中") {
        return "处理中";
      } else if (statusLower === "completed" || statusLower === "已完成") {
        return "已完成";
      } else if (statusLower === "failed" || statusLower === "失败") {
        return "失败";
      } else {
        return status;
      }
    },
    // 获取状态样式类
    getStatusClass(status) {
      const statusLower = typeof status === "string" ? status.toLowerCase() : status;
      if (statusLower === "processing" || statusLower === "处理中") {
        return "status-processing";
      } else if (statusLower === "completed" || statusLower === "已完成") {
        return "status-success";
      } else if (statusLower === "failed" || statusLower === "失败") {
        return "status-failed";
      } else {
        return "status-processing";
      }
    },
    // 获取兑换记录
    async getExchangeRecords() {
      if (this.isLoading)
        return;
      try {
        this.isLoading = true;
        const params = {
          page: this.page,
          size: this.size
        };
        const res = await apis_index.vip.getExchangeRecords(params);
        if (res.code === 0 && res.data) {
          if (this.page === 1) {
            this.recordList = res.data.list || [];
          } else {
            this.recordList = [...this.recordList, ...res.data.list || []];
          }
          this.totalPages = res.data.pages || 0;
          if (this.page >= this.totalPages) {
            this.loadMoreStatus = "noMore";
          } else {
            this.loadMoreStatus = "more";
          }
        } else {
          common_vendor.index.showToast({
            title: res.message || "获取兑换记录失败",
            icon: "none"
          });
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/exchange/record.vue:165", "获取兑换记录失败:", error);
        common_vendor.index.showToast({
          title: "获取记录失败，请重试",
          icon: "none"
        });
      } finally {
        this.isLoading = false;
      }
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _easycom_uni_load_more2 = common_vendor.resolveComponent("uni-load-more");
  (_easycom_uni_icons2 + _easycom_uni_load_more2)();
}
const _easycom_uni_icons = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
const _easycom_uni_load_more = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-load-more/uni-load-more.js";
if (!Math) {
  (_easycom_uni_icons + _easycom_uni_load_more)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: _ctx.statusBarHeight + "px",
    b: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    c: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    d: common_vendor.f($data.recordList, (item, k0, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.exchangeTime),
        b: common_vendor.t($options.formatStatus(item.status)),
        c: common_vendor.n($options.getStatusClass(item.status)),
        d: common_vendor.t(item.productName),
        e: common_vendor.t(item.duration),
        f: common_vendor.t(item.rechargeAccount),
        g: item.remark
      }, item.remark ? {
        h: common_vendor.t(item.remark)
      } : {}, {
        i: item.id
      });
    }),
    e: $data.recordList.length > 0
  }, $data.recordList.length > 0 ? {
    f: common_vendor.p({
      status: $data.loadMoreStatus
    })
  } : {}, {
    g: $data.recordList.length === 0 && !$data.isLoading
  }, $data.recordList.length === 0 && !$data.isLoading ? {
    h: common_assets._imports_0$4
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/vip/exchange/record.js.map
