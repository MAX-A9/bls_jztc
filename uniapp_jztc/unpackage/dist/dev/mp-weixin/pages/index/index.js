"use strict";
const common_vendor = require("../../common/vendor.js");
const mixins_device = require("../../mixins/device.js");
const mixins_share = require("../../mixins/share.js");
const utils_request = require("../../utils/request.js");
const apis_index = require("../../apis/index.js");
const FunctionArea = () => "../../components/function-area/index.js";
const SwiperBanner = () => "../../components/swiper-banner/index.js";
const ContentList = () => "../../components/content-list/index.js";
const ActivityArea = () => "../../components/activity-area/index.js";
const TabBar = () => "../../components/tab-bar/index.js";
const _sfc_main = {
  components: {
    FunctionArea,
    SwiperBanner,
    ContentList,
    ActivityArea,
    TabBar
  },
  mixins: [mixins_device.deviceMixin, mixins_share.shareMixin],
  data() {
    return {
      tabIndex: 0,
      currentLocation: "",
      isLocating: false,
      menuList: [],
      // 改为空数组，通过API获取
      menuLoading: false,
      // 添加加载状态
      currentMenu: 0,
      headerHeight: 0,
      showFixedMenu: false,
      // 是否显示固定菜单
      menuTop: 0,
      // 菜单导航的顶部位置
      contentTop: 0,
      // 内容区域顶部位置
      showBackToTop: false,
      // 是否显示返回顶部按钮
      scrollTop: 0,
      // 用于控制滚动位置
      oldScrollTop: 0,
      // 添加这个变量用于强制更新 scroll-top
      showLocationGuide: false,
      // 是否显示位置引导弹窗
      pageLoading: true,
      // 页面加载状态
      refreshing: false,
      // 添加下拉刷新状态
      logoUrl: "",
      // 默认logo地址
      contentListInitialized: false,
      // 标记内容列表是否已初始化
      isHomePage: true
      // 标记为首页，用于分享功能
    };
  },
  onLoad() {
    setTimeout(() => {
      if (this.pageLoading) {
        this.pageLoading = false;
        common_vendor.index.__f__("log", "at pages/index/index.vue:209", "页面加载超时，强制显示页面");
      }
    }, 3e3);
    this.initLocation();
    this.staggeredLoading();
  },
  mounted() {
    this.calculateHeaderHeight();
    setTimeout(() => {
      this.getMenuPosition();
    }, 300);
  },
  onReady() {
    this.getMenuPosition();
  },
  onResize() {
    this.calculateHeaderHeight();
    this.getMenuPosition();
  },
  onShow() {
    this.tabIndex = 0;
    const savedLocation = common_vendor.index.getStorageSync("currentLocation");
    if (savedLocation) {
      this.currentLocation = savedLocation;
      common_vendor.index.__f__("log", "at pages/index/index.vue:243", "页面显示时更新位置信息:", savedLocation);
    }
    this.$nextTick(() => {
      if (this.$refs.contentList && !this.contentListInitialized) {
        common_vendor.index.__f__("log", "at pages/index/index.vue:249", "页面显示时，初始化内容列表组件");
        this.$refs.contentList.init();
        this.contentListInitialized = true;
      } else if (this.$refs.contentList) {
        common_vendor.index.__f__("log", "at pages/index/index.vue:254", "内容列表已初始化，使用已有数据");
      }
    });
    const isFirstLaunch = !common_vendor.index.getStorageSync("notFirstLaunch");
    if (isFirstLaunch) {
      common_vendor.index.setStorageSync("notFirstLaunch", true);
      setTimeout(() => {
        this.showLocationGuide = true;
      }, 1e3);
    }
    this.showShareMenu();
  },
  onNavigationBarButtonTap(e) {
    if (e.index === 0) {
      this.navigateToLocationSelect();
    }
  },
  methods: {
    // 添加一个分阶段加载的方法
    staggeredLoading() {
      setTimeout(() => {
        this.fetchMenuCategories();
      }, 200);
      setTimeout(() => {
        this.fetchBaseSettings();
      }, 800);
      setTimeout(() => {
        this.initVisibleComponents();
      }, 1500);
    },
    // 初始化可见区域内的组件
    initVisibleComponents() {
      if (this.$refs.functionArea) {
        this.$refs.functionArea.loadMiniProgramList();
      }
      setTimeout(() => {
        if (this.$refs.swiperBanner) {
          this.$refs.swiperBanner.fetchBannerData();
        }
      }, 300);
      setTimeout(() => {
        if (this.$refs.activityArea) {
          this.$refs.activityArea.fetchActivityData();
        }
      }, 600);
    },
    // 获取菜单分类数据
    fetchMenuCategories() {
      this.menuLoading = true;
      const defaultMenu = {
        id: 0,
        name: "推荐"
      };
      const fetchTimeout = setTimeout(() => {
        if (this.menuLoading) {
          this.useDefaultMenu();
          this.menuLoading = false;
          this.pageLoading = false;
          common_vendor.index.__f__("log", "at pages/index/index.vue:335", "获取菜单分类超时，使用默认值");
        }
      }, 3e3);
      utils_request.get("/wx/client/content/categories", { type: 1 }).then((res) => {
        clearTimeout(fetchTimeout);
        if (res.code === 0) {
          const categoryList = res.data.list || [];
          this.menuList = [defaultMenu, ...categoryList];
          common_vendor.index.__f__("log", "at pages/index/index.vue:350", "成功获取菜单分类:", this.menuList);
        } else {
          common_vendor.index.__f__("error", "at pages/index/index.vue:352", "获取菜单分类失败:", res.message || "未知错误");
          this.useDefaultMenu();
        }
      }).catch((err) => {
        clearTimeout(fetchTimeout);
        common_vendor.index.__f__("error", "at pages/index/index.vue:359", "请求菜单分类接口出错:", err);
        this.useDefaultMenu();
      }).finally(() => {
        clearTimeout(fetchTimeout);
        this.menuLoading = false;
        this.pageLoading = false;
      });
    },
    // 接口失败时使用默认菜单
    useDefaultMenu() {
      this.menuList = [];
    },
    // 初始化位置信息
    initLocation() {
      const savedLocation = common_vendor.index.getStorageSync("currentLocation");
      if (savedLocation) {
        this.currentLocation = savedLocation;
        common_vendor.index.__f__("log", "at pages/index/index.vue:381", "使用本地存储的位置信息:", savedLocation);
        const locationId = common_vendor.index.getStorageSync("currentLocationId");
        if (locationId) {
          common_vendor.index.__f__("log", "at pages/index/index.vue:386", "使用已保存的地区ID:", locationId);
          common_vendor.index.$emit("locationChanged", { regionId: locationId });
          this.$nextTick(() => {
            if (this.$refs.contentList && !this.contentListInitialized) {
              common_vendor.index.__f__("log", "at pages/index/index.vue:393", "位置初始化后，调用内容列表初始化方法");
              this.$refs.contentList.init();
              this.contentListInitialized = true;
            }
          });
        }
        return;
      }
      this.useDefaultLocation();
    },
    // 使用默认位置
    useDefaultLocation() {
      const savedLocation = common_vendor.index.getStorageSync("currentLocation");
      if (savedLocation) {
        this.currentLocation = savedLocation;
        common_vendor.index.__f__("log", "at pages/index/index.vue:413", "使用已保存的位置:", savedLocation);
        const locationId = common_vendor.index.getStorageSync("currentLocationId");
        if (locationId) {
          common_vendor.index.__f__("log", "at pages/index/index.vue:418", "使用已保存的地区ID:", locationId);
          common_vendor.index.$emit("locationChanged", { regionId: locationId });
          this.$nextTick(() => {
            if (this.$refs.contentList && !this.contentListInitialized) {
              common_vendor.index.__f__("log", "at pages/index/index.vue:425", "使用已保存位置后，调用内容列表初始化方法");
              this.$refs.contentList.init();
              this.contentListInitialized = true;
            }
          });
        }
        return;
      }
      const store = this.$store;
      if (store && store.state.region && store.state.region.regionList.length > 0) {
        const firstRegion = store.state.region.regionList[0];
        this.currentLocation = firstRegion.name;
        common_vendor.index.__f__("log", "at pages/index/index.vue:443", "从store中设置默认地区ID:", firstRegion.id);
        common_vendor.index.setStorageSync("currentLocationId", firstRegion.id);
        common_vendor.index.setStorageSync("currentLocation", firstRegion.name);
        common_vendor.index.$emit("locationChanged", { regionId: firstRegion.id });
        this.$nextTick(() => {
          if (this.$refs.contentList && !this.contentListInitialized) {
            common_vendor.index.__f__("log", "at pages/index/index.vue:453", "设置默认位置后，调用内容列表初始化方法");
            this.$refs.contentList.init();
            this.contentListInitialized = true;
          }
        });
        common_vendor.index.__f__("log", "at pages/index/index.vue:459", "设置默认位置:", firstRegion.name, "来自", firstRegion.location);
      } else {
        this.currentLocation = "请选择地区";
        if (store && store.dispatch) {
          store.dispatch("region/getRegionList").then(() => {
            if (store.state.region.regionList.length > 0) {
              const firstRegion = store.state.region.regionList[0];
              this.currentLocation = firstRegion.name;
              common_vendor.index.setStorageSync("currentLocation", this.currentLocation);
              common_vendor.index.__f__("log", "at pages/index/index.vue:472", "异步获取后设置地区ID:", firstRegion.id);
              common_vendor.index.setStorageSync("currentLocationId", firstRegion.id);
              common_vendor.index.$emit("locationChanged", { regionId: firstRegion.id });
              this.$nextTick(() => {
                if (this.$refs.contentList && !this.contentListInitialized) {
                  common_vendor.index.__f__("log", "at pages/index/index.vue:481", "异步获取位置后，调用内容列表初始化方法");
                  this.$refs.contentList.init();
                  this.contentListInitialized = true;
                }
              });
              common_vendor.index.__f__("log", "at pages/index/index.vue:487", "异步更新默认位置:", firstRegion.name);
            }
          }).catch((err) => {
            common_vendor.index.__f__("error", "at pages/index/index.vue:490", "获取地区列表失败:", err);
          });
        }
      }
      if (!savedLocation) {
        common_vendor.index.setStorageSync("currentLocation", this.currentLocation);
      }
    },
    // 关闭位置引导弹窗
    closeLocationGuide() {
      this.showLocationGuide = false;
    },
    // 引导页面跳转到位置选择
    goToLocationSelectFromGuide() {
      this.showLocationGuide = false;
      this.navigateToLocationSelect();
    },
    // 跳转到位置选择页面
    navigateToLocationSelect() {
      common_vendor.index.navigateTo({
        url: "/pages/location/select/index",
        events: {
          // 选择位置后的回调
          locationSelected: (data) => {
            if (data && data.name) {
              this.currentLocation = data.name;
              common_vendor.index.setStorageSync("currentLocation", data.name);
              if (data.id) {
                common_vendor.index.__f__("log", "at pages/index/index.vue:528", "位置选择数据:", data);
                common_vendor.index.removeStorageSync("currentLocationId");
                common_vendor.index.setStorageSync("currentLocationId", data.id);
                if (this.$refs.contentList) {
                  common_vendor.index.__f__("log", "at pages/index/index.vue:537", "直接调用内容列表组件方法更新区域ID:", data.id);
                  this.$refs.contentList.regionId = data.id;
                  this.$refs.contentList.resetAndLoad();
                }
                common_vendor.index.__f__("log", "at pages/index/index.vue:543", "发送locationChanged事件，regionId:", data.id);
                common_vendor.index.$emit("locationChanged", { regionId: data.id });
              }
              common_vendor.index.__f__("log", "at pages/index/index.vue:547", "用户选择了新位置:", data.name, "ID:", data.id);
            }
          }
        }
      });
    },
    switchMenu(index) {
      if (this.menuLoading || !this.menuList.length) {
        return;
      }
      if (this.currentMenu !== index) {
        this.currentMenu = index;
        const selectedCategory = this.menuList[index];
        const categoryId = selectedCategory ? selectedCategory.id : "";
        if (this.$refs.contentList) {
          this.$refs.contentList.handleCategoryChange(categoryId);
        }
      }
    },
    calculateHeaderHeight() {
      const query = common_vendor.index.createSelectorQuery().in(this);
      query.select(".fixed-header").boundingClientRect((data) => {
        this.headerHeight = data.height;
      }).exec();
    },
    handleSearch(e) {
      const keyword = e.detail.value;
      common_vendor.index.__f__("log", "at pages/index/index.vue:581", "搜索关键词:", keyword);
    },
    // 获取菜单导航的位置
    getMenuPosition() {
      const query = common_vendor.index.createSelectorQuery().in(this);
      query.select(".fixed-header").boundingClientRect((data) => {
        if (data) {
          this.headerHeight = data.height;
        }
      }).exec();
      query.select("#category-menu").boundingClientRect((data) => {
        if (data) {
          this.menuTop = data.top;
        }
      }).exec();
      query.select("#content-area").boundingClientRect((data) => {
        if (data) {
          this.contentTop = data.top - this.headerHeight;
          common_vendor.index.__f__("log", "at pages/index/index.vue:607", "内容区域位置:", this.contentTop);
        }
      }).exec();
    },
    // 处理滚动事件
    handleScroll(e) {
      const scrollTop = e.detail.scrollTop;
      this.oldScrollTop = scrollTop;
      this.showBackToTop = scrollTop > 500;
      if (this.contentTop && this.contentTop > 0) {
        this.showFixedMenu = scrollTop >= this.contentTop;
      } else {
        this.showFixedMenu = scrollTop >= this.menuTop;
      }
    },
    // 返回顶部
    scrollToTop() {
      this.scrollTop = 1;
      this.$nextTick(() => {
        this.scrollTop = 0;
      });
    },
    // 处理底部导航切换
    handleTabChange(index) {
      common_vendor.index.__f__("log", "at pages/index/index.vue:641", "切换到标签:", index);
    },
    // 处理发布按钮点击
    handlePublish() {
      common_vendor.index.__f__("log", "at pages/index/index.vue:646", "点击发布按钮");
    },
    // 处理下拉刷新
    handleRefresh() {
      this.refreshing = true;
      this.fetchMenuCategories();
      if (this.$refs.functionArea) {
        this.$refs.functionArea.loadMiniProgramList();
      }
      if (this.$refs.activityArea) {
        this.$refs.activityArea.fetchActivityData();
      }
      if (this.$refs.swiperBanner) {
        this.$refs.swiperBanner.fetchBannerData();
      }
      if (this.$refs.contentList) {
        this.$refs.contentList.resetAndLoad();
      }
      this.initShareData();
      setTimeout(() => {
        if (this.refreshing) {
          this.refreshing = false;
          common_vendor.index.__f__("log", "at pages/index/index.vue:682", "刷新状态超时自动关闭");
        }
      }, 3e3);
    },
    // 处理上拉加载更多
    handleScrollToLower() {
      common_vendor.index.__f__("log", "at pages/index/index.vue:688", "触发上拉加载更多");
      if (this.$refs.contentList) {
        this.$refs.contentList.loadMore();
      }
    },
    // 获取小程序基础设置
    fetchBaseSettings() {
      apis_index.settings.getBaseSettings().then((res) => {
        if (res.code === 0 && res.data) {
          if (res.data.logo) {
            this.logoUrl = res.data.logo;
            common_vendor.index.__f__("log", "at pages/index/index.vue:702", "成功获取小程序Logo:", this.logoUrl);
          }
        } else {
          common_vendor.index.__f__("error", "at pages/index/index.vue:705", "获取小程序基础设置失败:", res.message || "未知错误");
        }
      }).catch((err) => {
        common_vendor.index.__f__("error", "at pages/index/index.vue:709", "请求小程序基础设置接口出错:", err);
      });
    },
    // 显示分享菜单，包括分享朋友圈
    showShareMenu() {
      this.initShareData();
      common_vendor.index.showShareMenu({
        withShareTicket: true,
        menus: ["shareAppMessage", "shareTimeline"]
      });
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _component_function_area = common_vendor.resolveComponent("function-area");
  const _component_activity_area = common_vendor.resolveComponent("activity-area");
  const _component_swiper_banner = common_vendor.resolveComponent("swiper-banner");
  const _component_content_list = common_vendor.resolveComponent("content-list");
  const _component_tab_bar = common_vendor.resolveComponent("tab-bar");
  (_easycom_uni_icons2 + _component_function_area + _component_activity_area + _component_swiper_banner + _component_content_list + _component_tab_bar)();
}
const _easycom_uni_icons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.pageLoading
  }, $data.pageLoading ? {
    b: common_vendor.f(4, (i, k0, i0) => {
      return {
        a: i
      };
    })
  } : {}, {
    c: $data.logoUrl,
    d: _ctx.isIOS ? "28px" : "32px",
    e: common_vendor.t($data.currentLocation || "定位中..."),
    f: common_vendor.p({
      type: "bottom",
      size: "12",
      color: "#ffffff"
    }),
    g: common_vendor.o((...args) => $options.navigateToLocationSelect && $options.navigateToLocationSelect(...args)),
    h: _ctx.navBarHeight + "px",
    i: _ctx.statusBarHeight + "px",
    j: common_vendor.p({
      type: "search",
      size: "16",
      color: "#666666"
    }),
    k: common_vendor.o((...args) => $options.handleSearch && $options.handleSearch(...args)),
    l: !$data.showFixedMenu,
    m: common_vendor.f($data.menuList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.name),
        b: index,
        c: $data.currentMenu === index ? 1 : "",
        d: common_vendor.o(($event) => $options.switchMenu(index), index)
      };
    }),
    n: $data.menuLoading
  }, $data.menuLoading ? {} : {}, {
    o: $data.showFixedMenu ? 1 : "",
    p: $data.showFixedMenu,
    q: common_vendor.sr("functionArea", "6d276527-2"),
    r: common_vendor.sr("activityArea", "6d276527-3"),
    s: common_vendor.sr("swiperBanner", "6d276527-4"),
    t: common_vendor.f($data.menuList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.name),
        b: index,
        c: $data.currentMenu === index ? 1 : "",
        d: common_vendor.o(($event) => $options.switchMenu(index), index)
      };
    }),
    v: $data.menuLoading
  }, $data.menuLoading ? {} : {}, {
    w: common_vendor.sr("contentList", "6d276527-5"),
    x: $data.headerHeight + "px",
    y: common_vendor.o((...args) => $options.handleScroll && $options.handleScroll(...args)),
    z: $data.scrollTop,
    A: common_vendor.o((...args) => $options.handleScrollToLower && $options.handleScrollToLower(...args)),
    B: common_vendor.o((...args) => $options.handleRefresh && $options.handleRefresh(...args)),
    C: $data.refreshing,
    D: common_vendor.p({
      type: "top",
      size: "24",
      color: "#ffffff"
    }),
    E: $data.showBackToTop,
    F: common_vendor.o((...args) => $options.scrollToTop && $options.scrollToTop(...args)),
    G: common_vendor.o($options.handleTabChange),
    H: common_vendor.o($options.handlePublish),
    I: $data.showLocationGuide
  }, $data.showLocationGuide ? {
    J: common_vendor.t($data.currentLocation),
    K: common_vendor.o((...args) => $options.goToLocationSelectFromGuide && $options.goToLocationSelectFromGuide(...args)),
    L: common_vendor.o(() => {
    }),
    M: common_vendor.o((...args) => $options.closeLocationGuide && $options.closeLocationGuide(...args))
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/index/index.js.map
