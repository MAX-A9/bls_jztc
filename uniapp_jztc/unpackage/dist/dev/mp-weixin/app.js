"use strict";
Object.defineProperty(exports, Symbol.toStringTag, { value: "Module" });
const common_vendor = require("./common/vendor.js");
const utils_auth = require("./utils/auth.js");
const utils_navigationControl = require("./utils/navigation-control.js");
const utils_uniCompatibility = require("./utils/uni-compatibility.js");
const common_assets = require("./common/assets.js");
const store_index = require("./store/index.js");
if (!Math) {
  "./pages/index/index.js";
  "./pages/community/index.js";
  "./pages/message/index.js";
  "./pages/my/index.js";
  "./pages/publish/index.js";
  "./pages/publish/draft/index.js";
  "./pages/publish/idle/index.js";
  "./pages/publish/info/index.js";
  "./pages/my/publish/index.js";
  "./pages/my/follow/index.js";
  "./pages/my/favorite/index.js";
  "./pages/my/history/index.js";
  "./pages/my/orders/index.js";
  "./pages/my/orders/detail.js";
  "./pages/express/send/index.js";
  "./pages/express/search/index.js";
  "./pages/express/return/index.js";
  "./pages/service/index.js";
  "./pages/content/detail.js";
  "./pages/community/detail.js";
  "./pages/vip/music/index.js";
  "./pages/vip/exchange/index.js";
  "./pages/vip/exchange/record.js";
  "./pages/vip/video/index.js";
  "./pages/master/join/index.js";
  "./pages/settings/index.js";
  "./pages/location/select/index.js";
  "./pages/chat/detail.js";
  "./pages/notification/index.js";
  "./pages/agreement/index.js";
  "./pages/webview/index.js";
}
const _sfc_main$1 = {
  name: "LoginModal",
  data() {
    return {
      loading: false,
      isAgreed: true
    };
  },
  computed: {
    visible() {
      return this.$store.getters["user/showLoginModal"];
    }
  },
  methods: {
    // 处理登录
    async handleLogin() {
      if (!this.isAgreed) {
        common_vendor.index.showToast({
          title: "请先同意用户协议和隐私政策",
          icon: "none"
        });
        return;
      }
      if (this.loading)
        return;
      this.loading = true;
      try {
        await this.$store.dispatch("user/login");
        common_vendor.index.showToast({
          title: "登录成功",
          icon: "success"
        });
      } catch (error) {
        common_vendor.index.__f__("error", "at components/login-modal/index.vue:82", "登录失败:", error);
        common_vendor.index.showToast({
          title: "登录失败，请重试",
          icon: "none"
        });
      } finally {
        this.loading = false;
      }
    },
    // 切换协议勾选状态
    toggleAgreement() {
      this.isAgreed = !this.isAgreed;
    },
    // 跳转至协议页面
    goToAgreement(type) {
      common_vendor.index.navigateTo({
        url: `/pages/agreement/index?type=${type}`
      });
    },
    // 取消登录
    cancel() {
      this.$store.dispatch("user/hideLoginModal");
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  _easycom_uni_icons2();
}
const _easycom_uni_icons = () => "./node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render$1(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $options.visible
  }, $options.visible ? common_vendor.e({
    b: common_vendor.o((...args) => $options.cancel && $options.cancel(...args)),
    c: common_vendor.p({
      type: "close",
      size: "20",
      color: "#999999"
    }),
    d: common_vendor.o((...args) => $options.cancel && $options.cancel(...args)),
    e: common_assets._imports_0,
    f: !$data.loading
  }, !$data.loading ? {
    g: common_vendor.p({
      type: "weixin",
      size: "20",
      color: "#ffffff"
    })
  } : {}, {
    h: common_vendor.t($data.loading ? "登录中..." : "微信登录"),
    i: common_vendor.o((...args) => $options.handleLogin && $options.handleLogin(...args)),
    j: $data.loading,
    k: $data.isAgreed
  }, $data.isAgreed ? {} : {}, {
    l: common_vendor.n({
      checked: $data.isAgreed
    }),
    m: common_vendor.o((...args) => $options.toggleAgreement && $options.toggleAgreement(...args)),
    n: common_vendor.o(($event) => $options.goToAgreement("user")),
    o: common_vendor.o(($event) => $options.goToAgreement("privacy"))
  }) : {});
}
const LoginModal = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main$1, [["render", _sfc_render$1]]);
let isInitialized = false;
const _sfc_main = {
  components: {
    LoginModal
  },
  onLaunch: function() {
    if (!isInitialized) {
      utils_uniCompatibility.initUniCompatibility();
      this.initRequestInterceptor();
      this.initNavigationInterceptor();
      isInitialized = true;
    }
    setTimeout(() => {
      const loginTimeout = setTimeout(() => {
        common_vendor.index.__f__("log", "at App.vue:34", "自动登录超时，继续应用初始化");
      }, 3e3);
      this.autoLogin().finally(() => {
        clearTimeout(loginTimeout);
      });
    }, 100);
    this.disableSharedArrayBuffer();
  },
  onShow: function() {
    setTimeout(() => {
      const store = this.$store;
      if (store && store.dispatch) {
        common_vendor.index.__f__("log", "at App.vue:55", "App显示时开始加载区域列表数据...");
        store.dispatch("region/getRegionList").then((regionList) => {
          if (regionList && regionList.length > 0) {
            const currentLocationId = common_vendor.index.getStorageSync("currentLocationId");
            if (!currentLocationId) {
              const firstRegion = regionList[0];
              common_vendor.index.setStorageSync("currentLocationId", firstRegion.id);
              common_vendor.index.setStorageSync("currentLocation", firstRegion.name);
              common_vendor.index.__f__("log", "at App.vue:69", "App显示时更新位置信息:", firstRegion.name, firstRegion.id);
              common_vendor.index.$emit("locationChanged", { regionId: firstRegion.id });
            } else {
              common_vendor.index.__f__("log", "at App.vue:75", "App显示时发送位置更新事件，使用已有位置ID:", currentLocationId);
              common_vendor.index.$emit("locationChanged", { regionId: currentLocationId });
            }
          }
        }).catch((err) => {
          common_vendor.index.__f__("error", "at App.vue:81", "初始化区域列表失败:", err);
        });
      }
    }, 200);
  },
  onHide: function() {
  },
  methods: {
    // 自动登录方法
    async autoLogin() {
      try {
        const hasToken = this.$store.getters["user/token"];
        if (!hasToken) {
          const result = await utils_auth.checkAndAutoLogin();
          if (result) {
            if (result.token) {
              this.$store.commit("user/SET_TOKEN", result.token);
            }
            this.$store.commit("user/SET_USER_INFO", result.userInfo || result);
          }
        } else {
          const userInfo = this.$store.getters["user/userInfo"];
          const hasUserId = userInfo && (userInfo.clientId || userInfo.id);
          const hasBasicInfo = userInfo && userInfo.realName;
          if (!hasUserId || !hasBasicInfo) {
            try {
              await this.$store.dispatch("user/getUserInfo");
            } catch (infoError) {
              await this.$store.dispatch("user/silentLogin");
            }
          }
        }
      } catch (error) {
        common_vendor.index.__f__("log", "at App.vue:123", "自动登录失败:", error);
      }
      return Promise.resolve();
    },
    // 初始化请求拦截器
    initRequestInterceptor() {
      common_vendor.index.addInterceptor("request", {
        fail(err) {
          return err;
        }
      });
      common_vendor.index.addInterceptor("response", {
        response(res) {
          if (res.statusCode === 401) {
            common_vendor.index.showToast({
              title: "登录已过期，请重新登录",
              icon: "none"
            });
            const store = getApp().$vm.$store;
            if (store) {
              store.dispatch("user/logout");
              store.dispatch("user/showLoginModal");
            }
          }
          return res;
        }
      });
    },
    // 初始化导航拦截器
    initNavigationInterceptor() {
      utils_navigationControl.navigationControl.initNavigationInterceptors();
      common_vendor.index.addInterceptor("navigateBack", {
        invoke(args) {
          const pages = getCurrentPages();
          if (pages.length < 2)
            return args;
          const currentPage = pages[pages.length - 1];
          const targetPage = pages[pages.length - 2];
          const currentPath = `/${currentPage.route}`;
          const targetPath = `/${targetPage.route}`;
          utils_navigationControl.navigationControl.beforeNavigateBack(currentPath, targetPath);
          return args;
        }
      });
    },
    // 禁用SharedArrayBuffer相关功能
    disableSharedArrayBuffer() {
      common_vendor.index.addInterceptor("request", {
        invoke(args) {
          if (!args.header) {
            args.header = {};
          }
          args.header["Cross-Origin-Opener-Policy"] = "same-origin";
          args.header["Cross-Origin-Embedder-Policy"] = "require-corp";
          return args;
        }
      });
      if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.onWebViewLoad) {
        common_vendor.wx$1.onWebViewLoad(function() {
          try {
            if (typeof SharedArrayBuffer !== "undefined") {
              common_vendor.index.__f__("warn", "at App.vue:212", "已禁用SharedArrayBuffer以提高安全性");
            }
          } catch (e) {
          }
        });
      }
    }
  }
};
if (!Array) {
  const _component_login_modal = common_vendor.resolveComponent("login-modal");
  _component_login_modal();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {};
}
const App = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
function createApp() {
  const app = common_vendor.createSSRApp(App);
  app.use(store_index.store);
  app.config.errorHandler = function(err, instance, info) {
    common_vendor.index.__f__("error", "at main.js:50", "Vue错误:", err);
    common_vendor.index.__f__("error", "at main.js:51", "错误信息:", info);
  };
  return {
    app
  };
}
createApp().app.mount("#app");
exports.createApp = createApp;
//# sourceMappingURL=../.sourcemap/mp-weixin/app.js.map
