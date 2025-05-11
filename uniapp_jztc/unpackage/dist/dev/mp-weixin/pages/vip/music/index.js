"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const mixins_share = require("../../../mixins/share.js");
const utils_request = require("../../../utils/request.js");
const utils_storage = require("../../../utils/storage.js");
const utils_share = require("../../../utils/share.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  mixins: [mixins_device.deviceMixin, mixins_share.shareMixin],
  data() {
    return {
      platforms: [],
      packages: [],
      currentCategoryId: null,
      adInfo: {
        enableWxAd: false,
        rewardedVideoAdId: ""
      },
      videoAd: null,
      isLoading: false,
      userAvatar: "/static/images/default-avatar.png",
      memberDuration: {
        days: "0",
        hours: "0"
      },
      shareImageUrl: ""
    };
  },
  computed: {
    ...common_vendor.mapState("user", ["userInfo", "isLogin"])
  },
  onLoad() {
    this.getAdSettings();
    this.getPlatformList();
    this.getUserAvatar();
    this.getMemberDuration();
    this.getShareImage();
  },
  // 自定义分享给好友
  onShareAppMessage() {
    return {
      title: "天天领音乐会员",
      // 使用页面导航栏标题
      path: "/pages/vip/music/index",
      imageUrl: this.shareImageUrl
    };
  },
  // 自定义分享到朋友圈
  onShareTimeline() {
    return {
      title: "天天领音乐会员",
      // 使用页面导航栏标题
      query: "",
      imageUrl: this.shareImageUrl
    };
  },
  methods: {
    // 获取用户头像
    getUserAvatar() {
      try {
        if (this.isLogin && this.userInfo && this.userInfo.avatarUrl) {
          this.userAvatar = this.userInfo.avatarUrl;
          return;
        }
        const userInfo = utils_storage.getUserInfo();
        if (userInfo && userInfo.avatarUrl) {
          this.userAvatar = userInfo.avatarUrl;
        } else {
          this.userAvatar = "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0";
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:210", "获取用户头像失败", error);
        this.userAvatar = "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0";
      }
    },
    // 获取平台列表
    async getPlatformList() {
      try {
        const res = await utils_request.get("/wx/client/shop-category/list");
        if (res.code === 0 && res.data && res.data.list && res.data.list.length > 0) {
          this.platforms = res.data.list.map((item) => {
            return {
              id: item.id,
              name: item.name,
              icon: item.image
            };
          });
          if (this.platforms.length > 0) {
            this.currentCategoryId = this.platforms[0].id;
            this.getPackageList(this.currentCategoryId);
          }
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:235", "获取平台列表失败", error);
      }
    },
    // 获取会员套餐列表
    async getPackageList(categoryId) {
      try {
        common_vendor.index.__f__("log", "at pages/vip/music/index.vue:242", "开始获取会员套餐列表，分类ID:", categoryId);
        common_vendor.index.showLoading({ title: "加载中..." });
        const res = await utils_request.request({
          url: "/wx/product/list",
          method: "GET",
          data: { categoryId }
        });
        common_vendor.index.hideLoading();
        common_vendor.index.__f__("log", "at pages/vip/music/index.vue:254", "套餐接口响应:", JSON.stringify(res));
        this.packages = [];
        if (res.code === 0) {
          if (res.data && res.data.list && res.data.list.length > 0) {
            setTimeout(() => {
              this.packages = res.data.list.map((item) => {
                return {
                  id: item.id,
                  tag: item.tags,
                  name: item.name,
                  price: item.price.toFixed(2),
                  days: item.duration,
                  description: item.description
                };
              });
              common_vendor.index.__f__("log", "at pages/vip/music/index.vue:273", "处理后的套餐数据:", JSON.stringify(this.packages));
            }, 100);
          } else {
            common_vendor.index.__f__("log", "at pages/vip/music/index.vue:276", "当前平台没有可用套餐");
          }
        } else {
          common_vendor.index.__f__("log", "at pages/vip/music/index.vue:279", "接口请求失败:", res.message || "未知错误");
        }
      } catch (error) {
        common_vendor.index.hideLoading();
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:283", "获取会员套餐列表失败", error);
      }
    },
    // 选择平台
    selectPlatform(categoryId) {
      common_vendor.index.__f__("log", "at pages/vip/music/index.vue:289", "选择平台:", categoryId, "当前平台:", this.currentCategoryId);
      if (this.currentCategoryId !== categoryId) {
        this.currentCategoryId = categoryId;
        this.getPackageList(categoryId);
      }
    },
    // 获取广告设置
    async getAdSettings() {
      let loading = false;
      try {
        loading = true;
        common_vendor.index.showLoading({ title: "加载中..." });
        const res = await utils_request.request({
          url: "/wx/ad/settings",
          method: "GET"
        });
        if (res && res.code === 0) {
          this.adInfo = res.data;
          if (this.adInfo.enableWxAd && this.adInfo.rewardedVideoAdId) {
            this.createRewardedVideoAd();
          }
        } else {
          common_vendor.index.__f__("error", "at pages/vip/music/index.vue:316", "获取广告设置失败:", res == null ? void 0 : res.message);
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:319", "获取广告设置异常:", error);
      } finally {
        if (loading) {
          common_vendor.index.hideLoading();
        }
      }
    },
    // 创建激励视频广告
    createRewardedVideoAd() {
      if (!common_vendor.wx$1 || !common_vendor.wx$1.createRewardedVideoAd) {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:331", "当前环境不支持激励视频广告");
        return;
      }
      this.videoAd = common_vendor.wx$1.createRewardedVideoAd({
        adUnitId: this.adInfo.rewardedVideoAdId
      });
      this.videoAd.onLoad(() => {
        common_vendor.index.__f__("log", "at pages/vip/music/index.vue:342", "激励视频广告加载成功");
      });
      this.videoAd.onError((err) => {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:347", "激励视频广告出错:", err);
        common_vendor.index.showToast({
          title: "广告加载失败",
          icon: "none"
        });
      });
      this.videoAd.onClose((res) => {
        if (res && res.isEnded) {
          this.reportAdViewed();
        } else {
          common_vendor.index.showToast({
            title: "请完整观看广告才能获得奖励",
            icon: "none"
          });
        }
      });
    },
    // 显示激励视频广告
    showRewardedVideoAd() {
      if (!this.videoAd) {
        common_vendor.index.showToast({
          title: "广告加载中，请稍后再试",
          icon: "none"
        });
        return;
      }
      this.videoAd.show().catch(() => {
        this.videoAd.load().then(() => this.videoAd.show()).catch((err) => {
          common_vendor.index.__f__("error", "at pages/vip/music/index.vue:385", "激励视频广告显示失败:", err);
          common_vendor.index.showToast({
            title: "广告加载失败，请稍后再试",
            icon: "none"
          });
        });
      });
    },
    // 上报广告观看完成
    async reportAdViewed() {
      var _a;
      if (this.isLoading)
        return;
      this.isLoading = true;
      let loading = false;
      try {
        loading = true;
        common_vendor.index.showLoading({ title: "领取奖励中..." });
        const res = await utils_request.post("/wx/ad/reward/viewed");
        if (res && res.code === 0) {
          const message = ((_a = res.data) == null ? void 0 : _a.message) || "恭喜获得会员奖励！";
          if (loading) {
            common_vendor.index.hideLoading();
            loading = false;
          }
          common_vendor.index.showToast({
            title: message,
            icon: "success",
            duration: 3e3
          });
          setTimeout(() => {
            this.refreshUserInfo();
          }, 2e3);
        } else {
          if (loading) {
            common_vendor.index.hideLoading();
            loading = false;
          }
          common_vendor.index.showToast({
            title: (res == null ? void 0 : res.message) || "领取失败，请重试",
            icon: "none"
          });
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:440", "上报广告观看失败:", error);
        if (loading) {
          common_vendor.index.hideLoading();
          loading = false;
        }
        common_vendor.index.showToast({
          title: "网络异常，请重试",
          icon: "none"
        });
      } finally {
        if (loading) {
          common_vendor.index.hideLoading();
        }
        this.isLoading = false;
      }
    },
    // 刷新用户信息
    refreshUserInfo() {
      let loading = false;
      try {
        loading = true;
        common_vendor.index.showLoading({ title: "刷新数据..." });
        this.getMemberDuration();
        setTimeout(() => {
        }, 1e3);
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:478", "刷新用户信息失败:", error);
      } finally {
        setTimeout(() => {
          if (loading) {
            common_vendor.index.hideLoading();
          }
        }, 1e3);
      }
    },
    // 获取会员时长
    async getMemberDuration() {
      try {
        const res = await utils_request.get("/wx/client/duration/remaining");
        if (res.code === 0 && res.data) {
          const durationStr = res.data.remainingDuration || "0天0小时0分钟";
          const days = durationStr.match(/(\d+)天/) ? durationStr.match(/(\d+)天/)[1] : "0";
          const hours = durationStr.match(/(\d+)小时/) ? durationStr.match(/(\d+)小时/)[1] : "0";
          this.memberDuration = {
            days,
            hours
          };
        } else {
          common_vendor.index.__f__("error", "at pages/vip/music/index.vue:506", "获取会员时长失败:", res == null ? void 0 : res.message);
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:509", "获取会员时长异常:", error);
      }
    },
    handleBack() {
      common_vendor.index.navigateBack();
    },
    // 修改观看按钮点击事件
    handleWatch() {
      if (this.adInfo.enableWxAd && this.adInfo.rewardedVideoAdId) {
        this.showRewardedVideoAd();
      } else {
        common_vendor.index.showToast({
          title: "广告功能暂未开放",
          icon: "none"
        });
      }
    },
    goToExchange() {
      common_vendor.index.navigateTo({
        url: "/pages/vip/exchange/index"
      });
    },
    // 获取分享图片
    async getShareImage() {
      try {
        const settings = await utils_share.getShareSettings();
        if (settings && settings.default_share_image) {
          this.shareImageUrl = settings.default_share_image;
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/vip/music/index.vue:545", "获取分享图片失败:", error);
      }
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
      color: "#ffecd8"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: _ctx.statusBarHeight + "px",
    d: $data.userAvatar,
    e: common_vendor.t($data.memberDuration.days),
    f: common_vendor.t($data.memberDuration.hours),
    g: common_vendor.p({
      type: "right",
      size: "14",
      color: "#862c13"
    }),
    h: common_vendor.p({
      type: "right",
      size: "14",
      color: "#862c13"
    }),
    i: common_vendor.o((...args) => $options.goToExchange && $options.goToExchange(...args)),
    j: common_vendor.f($data.platforms, (item, index, i0) => {
      return {
        a: item.icon,
        b: common_vendor.t(item.name),
        c: index,
        d: $data.currentCategoryId === item.id ? 1 : "",
        e: common_vendor.o(($event) => $options.selectPlatform(item.id), index)
      };
    }),
    k: $data.packages.length > 0
  }, $data.packages.length > 0 ? {
    l: common_vendor.f($data.packages, (item, index, i0) => {
      return {
        a: common_vendor.t(item.tag),
        b: common_vendor.t(item.name),
        c: common_vendor.t(item.description),
        d: common_vendor.t(item.days),
        e: index
      };
    })
  } : {}, {
    m: $data.packages.length === 0
  }, $data.packages.length === 0 ? {} : {}, {
    n: common_vendor.o((...args) => $options.handleWatch && $options.handleWatch(...args)),
    o: common_vendor.p({
      type: "info-filled",
      size: "20",
      color: "#ffecd8"
    }),
    p: common_assets._imports_1
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
_sfc_main.__runtimeHooks = 6;
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/vip/music/index.js.map
