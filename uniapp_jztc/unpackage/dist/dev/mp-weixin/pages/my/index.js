"use strict";
const common_vendor = require("../../common/vendor.js");
const mixins_device = require("../../mixins/device.js");
const apis_content = require("../../apis/content.js");
const apis_user = require("../../apis/user.js");
const apis_index = require("../../apis/index.js");
const utils_deviceInfo = require("../../utils/device-info.js");
const TabBar = () => "../../components/tab-bar/index.js";
const _sfc_main = {
  components: {
    TabBar
  },
  mixins: [mixins_device.deviceMixin],
  data() {
    return {
      tabIndex: 4,
      statusBarHeight: 0,
      navBarHeight: 44,
      showQrCodePopup: false,
      favoriteCount: 0,
      // 收藏数量
      browseHistoryCount: 0,
      // 浏览记录数量
      followingCount: 0,
      // 关注数量
      publishCount: 0,
      // 发布数量
      butlerImageUrl: ""
      // 管家二维码图片
    };
  },
  computed: {
    // 从Vuex中获取用户信息和登录状态
    userInfo() {
      return this.$store.getters["user/userInfo"];
    },
    hasLogin() {
      return this.$store.getters["user/isLoggedIn"];
    },
    loginLoading() {
      return this.$store.state.user.loginLoading;
    }
  },
  onLoad() {
    this.statusBarHeight = utils_deviceInfo.deviceInfo.getStatusBarHeight();
    this.loadFavoriteData();
  },
  onShow() {
    this.tabIndex = 4;
    this.checkAndRefreshUserInfo();
    this.getFavoriteCount();
    this.getBrowseHistoryCount();
    this.getFollowingCount();
    this.getPublishCount();
  },
  methods: {
    // 处理登录方法
    async handleLogin() {
      try {
        common_vendor.index.showLoading({
          title: "登录中",
          mask: true
        });
        const code = await apis_user.getWxLoginCode();
        const loginData = { code };
        const loginResult = await apis_user.wxappLogin(loginData);
        if (loginResult && loginResult.code === 0 && loginResult.data && loginResult.data.token) {
          this.$store.commit("user/SET_TOKEN", loginResult.data.token);
          await this.$store.dispatch("user/getUserInfo");
          common_vendor.index.hideLoading();
          common_vendor.index.showToast({
            title: "登录成功",
            icon: "success"
          });
          this.getFavoriteCount();
          this.getBrowseHistoryCount();
          this.getFollowingCount();
        } else {
          common_vendor.index.hideLoading();
          common_vendor.index.showToast({
            title: loginResult.message || "登录失败",
            icon: "none",
            duration: 2e3
          });
        }
      } catch (error) {
        common_vendor.index.hideLoading();
        common_vendor.index.__f__("error", "at pages/my/index.vue:262", "登录错误:", error);
        common_vendor.index.showToast({
          title: error.message || "登录失败，请重试",
          icon: "none",
          duration: 2e3
        });
      }
    },
    // 处理用户信息区域点击
    handleUserInfoClick() {
      if (!this.hasLogin) {
        this.handleLogin();
      } else if (!this.userInfo.phone) {
        common_vendor.index.showToast({
          title: "该功能暂未开放",
          icon: "none"
        });
      }
    },
    // 检查并刷新用户信息
    async checkAndRefreshUserInfo() {
      const isLoggedIn = this.$store.getters["user/isLoggedIn"];
      if (!isLoggedIn) {
        return;
      }
      const userInfo = this.$store.getters["user/userInfo"];
      const hasUserId = userInfo && (userInfo.clientId || userInfo.id);
      if (hasUserId && userInfo.realName) {
        return;
      }
      try {
        await this.$store.dispatch("user/getUserInfo");
      } catch (error) {
        if (error.code === 401) {
          try {
            await this.$store.dispatch("user/silentLogin");
          } catch (loginError) {
          }
        }
      }
    },
    handleSetting() {
      common_vendor.index.navigateTo({
        url: "/pages/settings/index"
      });
    },
    handleFeatureClick(type) {
      if (!this.hasLogin) {
        common_vendor.index.showToast({
          title: "请先登录",
          icon: "none"
        });
        return;
      }
      const pages = {
        publish: "/pages/my/publish/index",
        follow: "/pages/my/follow/index",
        favorite: "/pages/my/favorite/index",
        history: "/pages/my/history/index"
      };
      if (pages[type]) {
        common_vendor.index.navigateTo({
          url: pages[type]
        });
      } else {
        common_vendor.index.showToast({
          title: "功能开发中",
          icon: "none"
        });
      }
    },
    handleSignIn() {
      if (!this.hasLogin) {
        common_vendor.index.showToast({
          title: "请先登录",
          icon: "none"
        });
        return;
      }
      common_vendor.index.navigateTo({
        url: "/pages/vip/music/index"
      });
    },
    handleViewAllOrders() {
      if (!this.hasLogin) {
        common_vendor.index.showToast({
          title: "请先登录",
          icon: "none"
        });
        return;
      }
      common_vendor.index.navigateTo({
        url: "/pages/my/orders/index"
      });
    },
    handleOrderType(type) {
      if (!this.hasLogin) {
        common_vendor.index.showToast({
          title: "请先登录",
          icon: "none"
        });
        return;
      }
      let tabIndex = 0;
      switch (type) {
        case "processing":
          tabIndex = 1;
          break;
        case "unpaid":
          tabIndex = 2;
          break;
        case "completed":
          tabIndex = 3;
          break;
        default:
          tabIndex = 0;
      }
      common_vendor.index.navigateTo({
        url: `/pages/my/orders/index?tab=${tabIndex}`
      });
    },
    showQrCode() {
      this.showQrCodePopup = true;
      this.getButlerImage();
    },
    closeQrCode() {
      this.showQrCodePopup = false;
    },
    // 获取管家二维码图片
    async getButlerImage() {
      try {
        const res = await apis_user.getButlerImage();
        if (res.code === 0 && res.data && res.data.imageUrl) {
          this.butlerImageUrl = res.data.imageUrl;
        } else {
          common_vendor.index.__f__("error", "at pages/my/index.vue:438", "获取管家二维码失败", res);
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/my/index.vue:441", "获取管家二维码出错", error);
      }
    },
    handleMasterJoin() {
      if (!this.hasLogin) {
        common_vendor.index.showToast({
          title: "请先登录",
          icon: "none"
        });
        return;
      }
      common_vendor.index.navigateTo({
        url: "/pages/master/join/index"
      });
    },
    // 处理联系客服
    handleContact() {
      if (common_vendor.wx$1.openCustomerServiceChat) {
        common_vendor.wx$1.openCustomerServiceChat({
          extInfo: { url: "https://work.weixin.qq.com/" },
          // 根据实际情况修改
          corpId: "ww5823288888ed1111",
          // 需要替换为企业微信的企业ID
          success(res) {
            common_vendor.index.__f__("log", "at pages/my/index.vue:468", "客服会话打开成功", res);
          },
          fail(err) {
            common_vendor.index.__f__("error", "at pages/my/index.vue:471", "客服会话打开失败", err);
            common_vendor.index.showModal({
              title: "在线客服",
              content: "无法连接客服，您可以通过右上角胶囊按钮 -> 关于 -> 客服来联系我们",
              confirmText: "我知道了",
              showCancel: false
            });
          }
        });
      } else {
        common_vendor.index.showModal({
          title: "在线客服",
          content: "您可以通过右上角胶囊按钮 -> 关于 -> 客服来联系我们",
          confirmText: "我知道了",
          showCancel: false
        });
      }
    },
    // 处理意见反馈
    handleFeedback() {
      if (typeof common_vendor.wx$1.showFeedback === "function") {
        common_vendor.wx$1.showFeedback();
      } else {
        common_vendor.index.showModal({
          title: "意见反馈",
          content: '您可以通过"右上角胶囊按钮 -> 关于 -> 反馈与投诉"来提交反馈',
          confirmText: "我知道了",
          showCancel: false
        });
      }
    },
    async getFavoriteCount() {
      if (!this.hasLogin) {
        this.favoriteCount = 0;
        return;
      }
      try {
        const res = await apis_content.getFavoriteCount();
        if (res.code === 0 && res.data) {
          this.favoriteCount = res.data.total || 0;
        } else {
          this.favoriteCount = 0;
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/my/index.vue:540", "获取收藏数量失败", error);
        this.favoriteCount = 0;
      }
    },
    async getBrowseHistoryCount() {
      if (!this.hasLogin) {
        this.browseHistoryCount = 0;
        return;
      }
      try {
        const res = await apis_content.getBrowseHistoryCount();
        if (res.code === 0 && res.data) {
          this.browseHistoryCount = res.data.count || 0;
        } else {
          this.browseHistoryCount = 0;
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/my/index.vue:560", "获取浏览记录数量失败", error);
        this.browseHistoryCount = 0;
      }
    },
    async getFollowingCount() {
      if (!this.hasLogin) {
        this.followingCount = 0;
        return;
      }
      try {
        const res = await apis_content.getFollowingCount();
        if (res.code === 0 && res.data) {
          this.followingCount = res.data.count || 0;
        } else {
          this.followingCount = 0;
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/my/index.vue:580", "获取关注数量失败", error);
        this.followingCount = 0;
      }
    },
    async getPublishCount() {
      if (!this.hasLogin) {
        this.publishCount = 0;
        return;
      }
      try {
        const res = await apis_index.publish.getPublishCount();
        if (res.code === 0 && res.data) {
          this.publishCount = res.data.total || 0;
        } else {
          this.publishCount = 0;
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/my/index.vue:600", "获取发布数量失败", error);
        this.publishCount = 0;
      }
    },
    loadFavoriteData() {
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _component_tab_bar = common_vendor.resolveComponent("tab-bar");
  (_easycom_uni_icons2 + _component_tab_bar)();
}
const _easycom_uni_icons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.statusBarHeight + "px",
    b: common_vendor.p({
      type: "gear",
      size: "24",
      color: "#ffffff"
    }),
    c: common_vendor.o((...args) => $options.handleSetting && $options.handleSetting(...args)),
    d: $data.navBarHeight + "px",
    e: $options.userInfo.avatarUrl || "/static/demo/0.png",
    f: common_vendor.t($options.userInfo.realName || $options.userInfo.nickName || "未登录"),
    g: common_vendor.t($options.userInfo.phone || ($options.hasLogin ? "绑定手机号" : "点击登录")),
    h: common_vendor.o((...args) => $options.handleUserInfoClick && $options.handleUserInfoClick(...args)),
    i: common_vendor.p({
      type: $options.hasLogin ? "gift" : "person",
      size: "16",
      color: "#ffffff"
    }),
    j: common_vendor.t($options.hasLogin ? "签到领福利" : "立即登录"),
    k: common_vendor.o(($event) => $options.hasLogin ? $options.handleSignIn() : $options.handleLogin()),
    l: common_vendor.t($data.publishCount),
    m: common_vendor.o(($event) => $options.handleFeatureClick("publish")),
    n: common_vendor.t($data.followingCount),
    o: common_vendor.o(($event) => $options.handleFeatureClick("follow")),
    p: common_vendor.t($data.favoriteCount),
    q: common_vendor.o(($event) => $options.handleFeatureClick("favorite")),
    r: common_vendor.t($data.browseHistoryCount),
    s: common_vendor.o(($event) => $options.handleFeatureClick("history")),
    t: common_vendor.p({
      type: "right",
      size: "14",
      color: "#999999"
    }),
    v: common_vendor.o((...args) => $options.handleViewAllOrders && $options.handleViewAllOrders(...args)),
    w: common_vendor.p({
      type: "refresh",
      size: "28",
      color: "#fc3e2b"
    }),
    x: common_vendor.o(($event) => $options.handleOrderType("processing")),
    y: common_vendor.p({
      type: "wallet",
      size: "28",
      color: "#fc3e2b"
    }),
    z: common_vendor.o(($event) => $options.handleOrderType("unpaid")),
    A: common_vendor.p({
      type: "checkbox-filled",
      size: "28",
      color: "#fc3e2b"
    }),
    B: common_vendor.o(($event) => $options.handleOrderType("completed")),
    C: common_vendor.p({
      type: "list",
      size: "28",
      color: "#fc3e2b"
    }),
    D: common_vendor.o(($event) => $options.handleOrderType("all")),
    E: common_vendor.p({
      type: "right",
      size: "16",
      color: "#999999"
    }),
    F: common_vendor.o((...args) => $options.showQrCode && $options.showQrCode(...args)),
    G: common_vendor.p({
      type: "shop",
      size: "28",
      color: "#fc3e2b"
    }),
    H: common_vendor.p({
      type: "headphones",
      size: "28",
      color: "#fc3e2b"
    }),
    I: common_vendor.o((...args) => $options.handleContact && $options.handleContact(...args)),
    J: common_vendor.p({
      type: "help",
      size: "28",
      color: "#fc3e2b"
    }),
    K: common_vendor.p({
      type: "chat",
      size: "28",
      color: "#fc3e2b"
    }),
    L: common_vendor.o((...args) => $options.handleFeedback && $options.handleFeedback(...args)),
    M: common_vendor.p({
      ["current-tab"]: $data.tabIndex
    }),
    N: $data.showQrCodePopup
  }, $data.showQrCodePopup ? {
    O: $data.butlerImageUrl || "/static/images/qrcode.png",
    P: common_vendor.p({
      type: "closeempty",
      size: "24",
      color: "#666666"
    }),
    Q: common_vendor.o((...args) => $options.closeQrCode && $options.closeQrCode(...args)),
    R: common_vendor.o(() => {
    }),
    S: common_vendor.o((...args) => $options.closeQrCode && $options.closeQrCode(...args))
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/my/index.js.map
