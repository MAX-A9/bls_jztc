"use strict";
const common_vendor = require("../../common/vendor.js");
const apis_content = require("../../apis/content.js");
const _sfc_main = {
  name: "FunctionArea",
  data() {
    return {
      scrollLeft: 0,
      contentWidth: 0,
      // 内容总宽度
      viewWidth: 0,
      // 可视区域宽度
      currentPage: 0,
      // 当前页码
      isGlobalEnabled: true,
      // 默认显示，加载完成后根据接口返回值更新
      miniProgramList: []
      // 小程序列表
    };
  },
  created() {
  },
  mounted() {
    this.initScrollBar();
  },
  methods: {
    // 加载小程序列表
    async loadMiniProgramList() {
      try {
        const res = await apis_content.getMiniProgramList();
        if (res.code === 0 && res.data) {
          this.isGlobalEnabled = res.data.isGlobalEnabled;
          if (this.isGlobalEnabled && res.data.list) {
            this.miniProgramList = res.data.list.filter((item) => item.isEnabled).sort((a, b) => a.order - b.order);
          }
        } else {
          common_vendor.index.__f__("error", "at components/function-area/index.vue:70", "获取小程序列表失败:", res.message);
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at components/function-area/index.vue:73", "加载小程序列表异常:", error);
      }
    },
    // 跳转到小程序
    navigateToMiniProgram(miniProgram) {
      common_vendor.index.navigateToMiniProgram({
        appId: miniProgram.appId,
        success(res) {
          common_vendor.index.__f__("log", "at components/function-area/index.vue:81", "跳转成功", res);
        },
        fail(err) {
          common_vendor.index.__f__("error", "at components/function-area/index.vue:84", "跳转失败", err);
          common_vendor.index.showToast({
            title: "跳转失败",
            icon: "none"
          });
        }
      });
    },
    initScrollBar() {
      const query = common_vendor.index.createSelectorQuery().in(this);
      query.select(".scroll-row").boundingClientRect((data) => {
        if (data) {
          this.viewWidth = data.width;
        }
      }).exec();
    },
    handleScroll(e) {
      const { scrollLeft, scrollWidth } = e.detail;
      const maxScroll = scrollWidth - this.viewWidth;
      this.currentPage = maxScroll > 0 ? scrollLeft / maxScroll : 0;
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.isGlobalEnabled
  }, $data.isGlobalEnabled ? common_vendor.e({
    b: common_vendor.f($data.miniProgramList, (item, k0, i0) => {
      return {
        a: item.logo,
        b: common_vendor.t(item.name),
        c: "mini-" + item.id,
        d: common_vendor.o(($event) => $options.navigateToMiniProgram(item), "mini-" + item.id)
      };
    }),
    c: common_vendor.o((...args) => $options.handleScroll && $options.handleScroll(...args)),
    d: $data.miniProgramList.length > 0
  }, $data.miniProgramList.length > 0 ? {
    e: `translateX(${$data.currentPage * 100}%)`
  } : {}) : {});
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/components/function-area/index.js.map
