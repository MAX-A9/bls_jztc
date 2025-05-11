"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const mixins_share = require("../../../mixins/share.js");
const utils_request = require("../../../utils/request.js");
const apis_index = require("../../../apis/index.js");
const utils_share = require("../../../utils/share.js");
const exchangePopup = () => "../../../components/exchange-popup/index.js";
const uniIcons = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
const _sfc_main = {
  components: {
    uniIcons,
    exchangePopup
  },
  mixins: [mixins_device.deviceMixin, mixins_share.shareMixin],
  data() {
    return {
      recentExchanges: [],
      // 最近兑换记录
      currentTab: 0,
      tabList: ["全部"],
      // 初始只有全部选项，后面从接口获取
      memberDays: "0",
      memberHours: "0",
      memberMinutes: "0",
      categories: [],
      // 保存完整的分类数据，包括ID
      products: [],
      // 保存产品列表数据
      page: 1,
      size: 10,
      currentCategoryId: null,
      // 当前选中的分类ID
      currentProduct: null,
      // 当前选中的产品
      shareImageUrl: ""
      // 保存分享图片地址
    };
  },
  onLoad() {
    this.getMemberDuration();
    this.getCategoryList();
    this.getRecentExchangeRecords();
    this.getShareImage();
  },
  // 自定义分享给好友
  onShareAppMessage() {
    return {
      title: "兑换视频会员权益",
      // 使用页面导航栏标题
      path: "/pages/vip/exchange/index",
      imageUrl: this.shareImageUrl
    };
  },
  // 自定义分享到朋友圈
  onShareTimeline() {
    return {
      title: "兑换视频会员权益",
      // 使用页面导航栏标题
      query: "",
      imageUrl: this.shareImageUrl
    };
  },
  methods: {
    // 获取分享图片
    async getShareImage() {
      try {
        const settings = await utils_share.getShareSettings();
        if (settings && settings.default_share_image) {
          this.shareImageUrl = settings.default_share_image;
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:179", "获取分享图片失败:", error);
      }
    },
    handleBack() {
      common_vendor.index.navigateBack();
    },
    handleExchangeRecord() {
      common_vendor.index.navigateTo({
        url: "/pages/vip/exchange/record"
      });
    },
    // 获取最近兑换记录
    async getRecentExchangeRecords() {
      try {
        const res = await apis_index.vip.getRecentExchangeRecords();
        if (res.code === 0 && res.data) {
          this.recentExchanges = res.data.list || [];
        } else {
          common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:198", "获取最近兑换记录失败:", res == null ? void 0 : res.message);
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:201", "获取最近兑换记录异常:", error);
      }
    },
    switchTab(index) {
      this.currentTab = index;
      if (index === 0) {
        this.currentCategoryId = null;
      } else if (this.categories.length > 0 && index - 1 < this.categories.length) {
        this.currentCategoryId = this.categories[index - 1].id;
      }
      this.getProductList();
    },
    handleExchangeVip(item) {
      this.currentProduct = item;
      this.$nextTick(() => {
        if (this.$refs.exchangePopup) {
          this.$refs.exchangePopup.open(item);
        } else {
          common_vendor.index.showToast({
            title: "组件加载中，请稍后再试",
            icon: "none"
          });
        }
      });
    },
    // 处理确认兑换
    handleConfirmExchange(data) {
      this.exchangeProduct(data);
    },
    async exchangeProduct(data) {
      if (!data)
        return;
      try {
        common_vendor.index.showLoading({ title: "兑换中..." });
        const res = await apis_index.vip.createExchangeRecord(data);
        common_vendor.index.hideLoading();
        if (res.code === 0) {
          common_vendor.index.showToast({
            title: "兑换成功",
            icon: "success",
            duration: 2e3
          });
          setTimeout(() => {
            this.getMemberDuration();
            this.getRecentExchangeRecords();
          }, 1e3);
        } else {
          common_vendor.index.showToast({
            title: res.msg || "兑换失败",
            icon: "none"
          });
        }
      } catch (error) {
        common_vendor.index.hideLoading();
        common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:271", "兑换失败:", error);
        common_vendor.index.showToast({
          title: "兑换失败，请重试",
          icon: "none"
        });
      }
    },
    // 获取会员时长
    async getMemberDuration() {
      try {
        const res = await utils_request.get("/wx/client/duration/remaining");
        if (res.code === 0 && res.data) {
          const remainingDuration = res.data.remainingDuration || "0天0小时0分钟";
          this.memberDays = remainingDuration.match(/(\d+)天/) ? remainingDuration.match(/(\d+)天/)[1] : "0";
          this.memberHours = remainingDuration.match(/(\d+)小时/) ? remainingDuration.match(/(\d+)小时/)[1] : "0";
          this.memberMinutes = remainingDuration.match(/(\d+)分钟/) ? remainingDuration.match(/(\d+)分钟/)[1] : "0";
        } else {
          common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:291", "获取会员时长失败:", res == null ? void 0 : res.message);
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:294", "获取会员时长异常:", error);
      }
    },
    // 获取平台分类列表
    async getCategoryList() {
      try {
        const res = await utils_request.get("/wx/client/shop-category/list");
        if (res.code === 0 && res.data && res.data.list && res.data.list.length > 0) {
          this.categories = res.data.list;
          this.tabList = ["全部", ...res.data.list.map((item) => item.name)];
          this.getProductList();
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:312", "获取平台分类列表失败:", error);
      }
    },
    // 获取产品列表
    async getProductList() {
      try {
        common_vendor.index.showLoading({ title: "加载中..." });
        const params = {
          page: this.page,
          size: this.size
        };
        if (this.currentCategoryId) {
          params.categoryId = this.currentCategoryId;
        }
        const res = await utils_request.get("/wx/product/list", params);
        if (res.code === 0 && res.data) {
          this.products = res.data.list || [];
        } else {
          common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:335", "获取产品列表失败:", res == null ? void 0 : res.message);
          this.products = [];
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/exchange/index.vue:339", "获取产品列表异常:", error);
        this.products = [];
      } finally {
        common_vendor.index.hideLoading();
      }
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _component_exchange_popup = common_vendor.resolveComponent("exchange-popup");
  (_easycom_uni_icons2 + _component_exchange_popup)();
}
const _easycom_uni_icons = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: _ctx.statusBarHeight + "px",
    b: common_vendor.p({
      type: "left",
      size: "20",
      color: "#ffecd8"
    }),
    c: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    d: $data.recentExchanges.length > 0
  }, $data.recentExchanges.length > 0 ? {
    e: common_vendor.f($data.recentExchanges, (item, index, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.clientName),
        b: common_vendor.t(item.productName),
        c: index < $data.recentExchanges.length - 1
      }, index < $data.recentExchanges.length - 1 ? {} : {}, {
        d: index
      });
    })
  } : {}, {
    f: $data.recentExchanges.length > 0
  }, $data.recentExchanges.length > 0 ? {
    g: common_vendor.f($data.recentExchanges, (item, index, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.clientName),
        b: common_vendor.t(item.productName),
        c: index < $data.recentExchanges.length - 1
      }, index < $data.recentExchanges.length - 1 ? {} : {}, {
        d: index
      });
    })
  } : {}, {
    h: common_vendor.o((...args) => $options.handleExchangeRecord && $options.handleExchangeRecord(...args)),
    i: common_vendor.t($data.memberDays),
    j: common_vendor.t($data.memberHours),
    k: common_vendor.t($data.memberMinutes),
    l: common_vendor.f($data.tabList, (item, index, i0) => {
      return {
        a: common_vendor.t(item),
        b: index,
        c: $data.currentTab === index ? 1 : "",
        d: common_vendor.o(($event) => $options.switchTab(index), index)
      };
    }),
    m: $data.products.length > 0
  }, $data.products.length > 0 ? {
    n: common_vendor.f($data.products, (item, k0, i0) => {
      return common_vendor.e({
        a: item.tags
      }, item.tags ? {
        b: common_vendor.t(item.tags)
      } : {}, {
        c: item.thumbnail,
        d: common_vendor.t(item.name),
        e: common_vendor.t(item.description),
        f: common_vendor.t(item.duration),
        g: common_vendor.o(($event) => $options.handleExchangeVip(item), item.id),
        h: item.id
      });
    })
  } : {}, {
    o: $data.products.length === 0
  }, $data.products.length === 0 ? {} : {}, {
    p: common_vendor.sr("exchangePopup", "0f245fb6-1"),
    q: common_vendor.o($options.handleConfirmExchange)
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
_sfc_main.__runtimeHooks = 6;
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/vip/exchange/index.js.map
