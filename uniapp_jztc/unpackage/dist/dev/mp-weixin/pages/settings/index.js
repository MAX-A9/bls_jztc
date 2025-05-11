"use strict";
const common_vendor = require("../../common/vendor.js");
const apis_index = require("../../apis/index.js");
const utils_storage = require("../../utils/storage.js");
const utils_deviceInfo = require("../../utils/device-info.js");
const _sfc_main = {
  data() {
    return {
      statusBarHeight: 0,
      navBarHeight: 44,
      userInfo: {
        avatar: "/static/demo/0.png",
        nickname: "橘子用户",
        phone: ""
      },
      notifications: {
        message: true,
        comment: true
      },
      cacheSize: "8.5MB",
      updateData: {},
      // 临时存储要更新的数据
      appInfo: {
        name: "橘子同城",
        description: "为您提供便捷的本地生活服务"
      },
      showAboutPopup: false
      // 控制关于我们弹窗显示
    };
  },
  computed: {
    ...common_vendor.mapGetters("user", ["isLoggedIn"])
  },
  onLoad() {
    this.statusBarHeight = utils_deviceInfo.deviceInfo.getStatusBarHeight();
    this.loadUserInfo();
    this.updateCacheSize();
    this.fetchAppBaseInfo();
  },
  methods: {
    ...common_vendor.mapActions("user", ["getUserInfo"]),
    // 加载用户信息
    async loadUserInfo() {
      if (this.isLoggedIn) {
        try {
          const userInfo = await this.getUserInfo();
          const localUserInfo = utils_storage.getUserInfo();
          this.userInfo = {
            avatar: localUserInfo.avatarUrl || "/static/demo/0.png",
            nickname: localUserInfo.realName || "橘子用户",
            phone: localUserInfo.phone || ""
          };
        } catch (err) {
          common_vendor.index.__f__("error", "at pages/settings/index.vue:197", "加载用户信息失败:", err);
          common_vendor.index.showToast({
            title: "加载用户信息失败",
            icon: "none"
          });
        }
      }
    },
    // 处理返回按钮点击
    handleBack() {
      common_vendor.index.navigateBack();
    },
    // 头像选择处理
    async onChooseAvatar(e) {
      try {
        const { avatarUrl } = e.detail;
        if (!avatarUrl) {
          return;
        }
        common_vendor.index.showLoading({ title: "处理中..." });
        this.userInfo.avatar = avatarUrl;
        try {
          const result = await apis_index.user.uploadAvatar(avatarUrl);
          await this.getUserInfo();
          common_vendor.index.hideLoading();
          common_vendor.index.showToast({
            title: "头像已更新",
            icon: "success"
          });
        } catch (err) {
          common_vendor.index.hideLoading();
          common_vendor.index.showToast({
            title: "头像更新失败: " + (err.message || "未知错误"),
            icon: "none",
            duration: 3e3
          });
        }
      } catch (err) {
        common_vendor.index.hideLoading();
        common_vendor.index.showToast({
          title: "更新头像失败",
          icon: "none"
        });
      }
    },
    // 昵称输入处理
    async onNicknameInput(e) {
      try {
        const nickname = e.detail.value;
        if (!nickname || nickname === this.userInfo.nickname) {
          return;
        }
        this.userInfo.nickname = nickname;
        this.updateData.realName = nickname;
        await this.saveUserInfo();
        common_vendor.index.showToast({
          title: "昵称已更新",
          icon: "success"
        });
      } catch (err) {
        common_vendor.index.__f__("error", "at pages/settings/index.vue:276", "更新昵称失败:", err);
        common_vendor.index.showToast({
          title: "更新昵称失败",
          icon: "none"
        });
      }
    },
    // 修改手机号
    handleEditPhone() {
      common_vendor.index.showModal({
        title: "修改手机号",
        editable: true,
        placeholderText: "请输入11位手机号",
        content: this.userInfo.phone || "",
        success: async (res) => {
          if (res.confirm && res.content) {
            const phone = res.content;
            if (!/^1\d{10}$/.test(phone)) {
              common_vendor.index.showToast({
                title: "请输入有效的手机号",
                icon: "none"
              });
              return;
            }
            if (phone === this.userInfo.phone) {
              return;
            }
            try {
              this.userInfo.phone = phone;
              this.updateData.phone = phone;
              common_vendor.index.showLoading({ title: "更新中..." });
              await this.saveUserInfo();
              common_vendor.index.hideLoading();
              common_vendor.index.showToast({
                title: "手机号已更新",
                icon: "success"
              });
            } catch (err) {
              common_vendor.index.hideLoading();
              common_vendor.index.__f__("error", "at pages/settings/index.vue:324", "更新手机号失败:", err);
              common_vendor.index.showToast({
                title: "更新手机号失败",
                icon: "none"
              });
            }
          }
        }
      });
    },
    // 保存用户信息到后端
    async saveUserInfo() {
      if (Object.keys(this.updateData).length === 0) {
        return;
      }
      try {
        await apis_index.user.updateClientProfile(this.updateData);
        this.getUserInfo();
        this.updateData = {};
        return true;
      } catch (err) {
        common_vendor.index.__f__("error", "at pages/settings/index.vue:357", "保存用户信息失败:", err);
        throw err;
      }
    },
    // 处理开关变化
    handleSwitchChange(type, event) {
      this.notifications[type] = event.detail.value;
    },
    // 处理清除缓存
    handleClearCache() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要清除缓存吗？",
        success: (res) => {
          if (res.confirm) {
            try {
              common_vendor.index.showLoading({
                title: "清除中...",
                mask: true
              });
              const keysToPreserve = ["token", "USER_INFO", "currentLocation"];
              common_vendor.index.getStorageInfo({
                success: (res2) => {
                  const allKeys = res2.keys || [];
                  const keysToRemove = allKeys.filter((key) => !keysToPreserve.includes(key));
                  keysToRemove.forEach((key) => {
                    common_vendor.index.removeStorageSync(key);
                  });
                  common_vendor.index.getSavedFileList({
                    success: (res3) => {
                      const fileList = res3.fileList || [];
                      fileList.forEach((file) => {
                        common_vendor.index.removeSavedFile({
                          filePath: file.filePath,
                          fail: () => {
                          }
                        });
                      });
                    },
                    complete: () => {
                    }
                  });
                  setTimeout(() => {
                    this.updateCacheSize();
                    common_vendor.index.hideLoading();
                    common_vendor.index.showToast({
                      title: "缓存已清除",
                      icon: "success"
                    });
                  }, 500);
                },
                fail: () => {
                  common_vendor.index.hideLoading();
                  common_vendor.index.showToast({
                    title: "清除缓存失败",
                    icon: "none"
                  });
                }
              });
            } catch (error) {
              common_vendor.index.__f__("error", "at pages/settings/index.vue:432", "清除缓存出错:", error);
              common_vendor.index.hideLoading();
              common_vendor.index.showToast({
                title: "清除缓存失败",
                icon: "none"
              });
            }
          }
        }
      });
    },
    // 更新缓存大小
    updateCacheSize() {
      try {
        common_vendor.index.getStorageInfo({
          success: (res) => {
            const sizeInKB = res.currentSize || 0;
            if (sizeInKB < 1024) {
              this.cacheSize = sizeInKB + "KB";
            } else {
              const sizeInMB = (sizeInKB / 1024).toFixed(1);
              this.cacheSize = sizeInMB + "MB";
            }
          },
          fail: () => {
            this.cacheSize = "0KB";
          }
        });
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/settings/index.vue:465", "获取缓存大小失败:", error);
        this.cacheSize = "未知";
      }
    },
    // 处理关于我们
    handleAbout() {
      this.showAboutPopup = true;
    },
    // 获取应用基础信息
    async fetchAppBaseInfo() {
      try {
        const res = await apis_index.settings.getBaseSettings();
        if (res.code === 0 && res.data) {
          if (res.data.name) {
            this.appInfo.name = res.data.name;
          }
          if (res.data.description) {
            this.appInfo.description = res.data.description;
          }
        } else {
          common_vendor.index.__f__("error", "at pages/settings/index.vue:488", "获取应用基础信息失败:", res.message || "未知错误");
        }
      } catch (err) {
        common_vendor.index.__f__("error", "at pages/settings/index.vue:491", "请求应用基础信息接口出错:", err);
      }
    },
    // 处理协议点击
    handleAgreement(type) {
      const apiType = type === "service" ? "user" : "privacy";
      common_vendor.index.navigateTo({
        url: `/pages/agreement/index?type=${apiType}`
      });
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  _easycom_uni_icons2();
}
const _easycom_uni_icons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
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
    c: $data.navBarHeight + "px",
    d: $data.statusBarHeight + "px",
    e: $data.userInfo.avatar,
    f: common_vendor.o((...args) => $options.onChooseAvatar && $options.onChooseAvatar(...args)),
    g: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    h: $data.userInfo.nickname || "请输入昵称",
    i: common_vendor.o((...args) => $options.onNicknameInput && $options.onNicknameInput(...args)),
    j: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    k: common_vendor.t($data.userInfo.phone || "未设置"),
    l: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    m: common_vendor.o((...args) => $options.handleEditPhone && $options.handleEditPhone(...args)),
    n: $data.notifications.message,
    o: common_vendor.o(($event) => $options.handleSwitchChange("message", $event)),
    p: $data.notifications.comment,
    q: common_vendor.o(($event) => $options.handleSwitchChange("comment", $event)),
    r: common_vendor.t($data.cacheSize),
    s: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    t: common_vendor.o((...args) => $options.handleClearCache && $options.handleClearCache(...args)),
    v: common_vendor.p({
      type: "right",
      size: "16",
      color: "#CCCCCC"
    }),
    w: common_vendor.o((...args) => $options.handleAbout && $options.handleAbout(...args)),
    x: common_vendor.o(($event) => $options.handleAgreement("service")),
    y: common_vendor.o(($event) => $options.handleAgreement("privacy")),
    z: `calc(${$data.statusBarHeight}px + ${$data.navBarHeight}px)`,
    A: $data.showAboutPopup
  }, $data.showAboutPopup ? {
    B: common_vendor.t($data.appInfo.name),
    C: common_vendor.t($data.appInfo.description),
    D: common_vendor.o(($event) => $data.showAboutPopup = false)
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/settings/index.js.map
