"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const _sfc_main = {
  mixins: [mixins_device.deviceMixin],
  data() {
    return {
      searchText: "",
      historyList: [
        {
          number: "SF1234567890",
          company: "顺丰快递",
          status: "delivered",
          statusText: "已签收"
        },
        {
          number: "YT9876543210",
          company: "圆通快递",
          status: "shipping",
          statusText: "运输中"
        }
      ],
      companyList: [
        { name: "顺丰快递", logo: "/static/express/sf.png", phone: "95338" },
        { name: "中通快递", logo: "/static/express/zt.png", phone: "95311" },
        { name: "韵达快递", logo: "/static/express/yd.png", phone: "95546" },
        { name: "申通快递", logo: "/static/express/st.png", phone: "95543" },
        { name: "圆通快递", logo: "/static/express/yt.png", phone: "95554" },
        { name: "百世快递", logo: "/static/express/bs.png", phone: "95320" },
        { name: "京东快递", logo: "/static/express/jd.png", phone: "950616" },
        { name: "邮政快递", logo: "/static/express/yz.png", phone: "11183" }
      ]
    };
  },
  methods: {
    handleBack() {
      common_vendor.index.navigateBack();
    },
    handleSearch() {
      if (!this.searchText) {
        common_vendor.index.showToast({
          title: "请输入运单号",
          icon: "none"
        });
        return;
      }
      common_vendor.index.__f__("log", "at pages/express/search/index.vue:122", "搜索运单号:", this.searchText);
    },
    handleScan() {
      common_vendor.index.scanCode({
        success: (res) => {
          this.searchText = res.result;
          this.handleSearch();
        }
      });
    },
    handleClearHistory() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要清空历史记录吗？",
        success: (res) => {
          if (res.confirm) {
            this.historyList = [];
          }
        }
      });
    },
    handleHistoryClick(item) {
      common_vendor.index.navigateTo({
        url: `/pages/express/detail?number=${item.number}`
      });
    },
    handleCompanyClick(item) {
      common_vendor.index.makePhoneCall({
        phoneNumber: item.phone,
        fail(err) {
          common_vendor.index.__f__("log", "at pages/express/search/index.vue:152", "拨打电话失败:", err);
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
  return common_vendor.e({
    a: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: _ctx.statusBarHeight + "px",
    d: common_vendor.p({
      type: "search",
      size: "18",
      color: "#666666"
    }),
    e: common_vendor.o((...args) => $options.handleSearch && $options.handleSearch(...args)),
    f: $data.searchText,
    g: common_vendor.o(($event) => $data.searchText = $event.detail.value),
    h: common_vendor.p({
      type: "scan",
      size: "20",
      color: "#007AFF"
    }),
    i: common_vendor.o((...args) => $options.handleScan && $options.handleScan(...args)),
    j: $data.historyList.length > 0
  }, $data.historyList.length > 0 ? {
    k: common_vendor.p({
      type: "trash",
      size: "14",
      color: "#999999"
    }),
    l: common_vendor.o((...args) => $options.handleClearHistory && $options.handleClearHistory(...args)),
    m: common_vendor.f($data.historyList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.number),
        b: common_vendor.t(item.company),
        c: common_vendor.t(item.statusText),
        d: common_vendor.n(item.status),
        e: index,
        f: common_vendor.o(($event) => $options.handleHistoryClick(item), index)
      };
    })
  } : {}, {
    n: common_vendor.f($data.companyList, (item, index, i0) => {
      return {
        a: item.logo,
        b: common_vendor.t(item.name),
        c: common_vendor.t(item.phone),
        d: index,
        e: common_vendor.o(($event) => $options.handleCompanyClick(item), index)
      };
    })
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/express/search/index.js.map
