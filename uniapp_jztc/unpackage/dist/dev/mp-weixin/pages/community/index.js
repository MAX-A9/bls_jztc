"use strict";
const common_vendor = require("../../common/vendor.js");
const mixins_device = require("../../mixins/device.js");
const mixins_share = require("../../mixins/share.js");
const utils_request = require("../../utils/request.js");
const utils_share = require("../../utils/share.js");
const TabBar = () => "../../components/tab-bar/index.js";
const CommunityNews = () => "../../components/community-news/index.js";
const _sfc_main = {
  components: {
    TabBar,
    CommunityNews
  },
  mixins: [mixins_device.deviceMixin, mixins_share.shareMixin],
  data() {
    return {
      tabIndex: 1,
      isRefreshing: false,
      currentMenu: 0,
      searchBoxWidth: 200,
      rightPadding: 0,
      menuList: [],
      menuLoading: false,
      isHomePage: true,
      // 标记为首页类型，用于分享功能
      shareData: null
    };
  },
  onLoad() {
    this.fetchCategories();
  },
  mounted() {
    const menuButtonInfo = common_vendor.index.getMenuButtonBoundingClientRect();
    const windowInfo = common_vendor.index.getWindowInfo();
    this.searchBoxWidth = menuButtonInfo.left - 100;
    this.rightPadding = windowInfo.windowWidth - menuButtonInfo.left;
  },
  onShow() {
    this.tabIndex = 1;
    this.updateShareData();
  },
  methods: {
    // 获取闲置分类数据
    fetchCategories() {
      this.menuLoading = true;
      const defaultMenu = {
        id: 0,
        name: "全部"
      };
      const fetchTimeout = setTimeout(() => {
        if (this.menuLoading) {
          this.useDefaultMenu();
          this.menuLoading = false;
          common_vendor.index.__f__("log", "at pages/community/index.vue:130", "获取闲置分类超时，使用默认值");
        }
      }, 3e3);
      utils_request.get("/wx/client/content/categories", { type: 2 }).then((res) => {
        clearTimeout(fetchTimeout);
        if (res.code === 0) {
          const categoryList = res.data.list || [];
          this.menuList = [defaultMenu, ...categoryList];
          common_vendor.index.__f__("log", "at pages/community/index.vue:145", "成功获取闲置分类:", this.menuList);
          this.updateShareData();
        } else {
          common_vendor.index.__f__("error", "at pages/community/index.vue:150", "获取闲置分类失败:", res.message || "未知错误");
          this.useDefaultMenu();
        }
      }).catch((err) => {
        clearTimeout(fetchTimeout);
        common_vendor.index.__f__("error", "at pages/community/index.vue:157", "请求闲置分类接口出错:", err);
        this.useDefaultMenu();
      }).finally(() => {
        clearTimeout(fetchTimeout);
        this.menuLoading = false;
      });
    },
    // 接口失败时使用默认菜单
    useDefaultMenu() {
      this.menuList = [
        { id: 0, name: "全部" }
      ];
    },
    handleSearch() {
      common_vendor.index.navigateTo({
        url: "/pages/search/index"
      });
    },
    switchMenu(index) {
      this.currentMenu = index;
      const selectedCategory = this.menuList[index];
      const categoryId = selectedCategory ? selectedCategory.id : 0;
      common_vendor.index.$emit("idleCategoryChanged", { categoryId });
      common_vendor.index.__f__("log", "at pages/community/index.vue:190", "切换到分类:", selectedCategory.name, "ID:", categoryId);
    },
    async onRefresh() {
      if (this.isRefreshing)
        return;
      this.isRefreshing = true;
      try {
        this.fetchCategories();
        if (this.$refs.newsComponent) {
          this.$refs.newsComponent.resetAndLoad();
        }
        await this.refreshData();
        this.updateShareData();
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/index.vue:210", "刷新失败:", error);
      } finally {
        this.isRefreshing = false;
      }
    },
    async refreshData() {
      await new Promise((resolve) => setTimeout(resolve, 1e3));
    },
    // 处理上拉加载更多
    handleScrollToLower() {
      common_vendor.index.__f__("log", "at pages/community/index.vue:220", "触发上拉加载更多");
      if (this.$refs.newsComponent) {
        this.$refs.newsComponent.loadMore();
      }
    },
    // 更新分享数据
    async updateShareData() {
      try {
        const settings = await utils_share.getShareSettings();
        if (settings) {
          const shareTitle = settings.content_share_text || "查看闲置社区";
          const shareImage = settings.content_share_image || "";
          this.shareData = {
            title: shareTitle,
            imageUrl: shareImage,
            path: "/pages/community/index"
          };
          common_vendor.index.__f__("log", "at pages/community/index.vue:244", "闲置社区分享数据已更新:", this.shareData);
        } else {
          this.initShareData();
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/index.vue:250", "更新分享数据失败:", error);
        this.initShareData();
      }
    },
    // 显示分享菜单
    handleShare() {
      this.updateShareData();
      common_vendor.index.showShareMenu({
        withShareTicket: true,
        menus: ["shareAppMessage", "shareTimeline"]
      });
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _component_community_news = common_vendor.resolveComponent("community-news");
  const _component_tab_bar = common_vendor.resolveComponent("tab-bar");
  (_easycom_uni_icons2 + _component_community_news + _component_tab_bar)();
}
const _easycom_uni_icons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.p({
      type: "search",
      size: "16",
      color: "#666666"
    }),
    b: common_vendor.o((...args) => $options.handleSearch && $options.handleSearch(...args)),
    c: $data.searchBoxWidth + "px",
    d: _ctx.navBarHeight + "px",
    e: $data.rightPadding + "px",
    f: _ctx.statusBarHeight + "px",
    g: common_vendor.f($data.menuList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.name),
        b: index,
        c: $data.currentMenu === index ? 1 : "",
        d: common_vendor.o(($event) => $options.switchMenu(index), index)
      };
    }),
    h: common_vendor.sr("newsComponent", "78e76304-1"),
    i: $data.isRefreshing,
    j: common_vendor.o((...args) => $options.onRefresh && $options.onRefresh(...args)),
    k: common_vendor.o((...args) => _ctx.onRestore && _ctx.onRestore(...args)),
    l: common_vendor.o((...args) => $options.handleScrollToLower && $options.handleScrollToLower(...args)),
    m: _ctx.navigationBarHeight + "px",
    n: common_vendor.p({
      ["current-tab"]: $data.tabIndex
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/community/index.js.map
