"use strict";
const common_vendor = require("../../../common/vendor.js");
const mixins_device = require("../../../mixins/device.js");
const utils_storage = require("../../../utils/storage.js");
const utils_deviceInfo = require("../../../utils/device-info.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  mixins: [mixins_device.deviceMixin],
  data() {
    return {
      statusBarHeight: 0,
      navBarHeight: 44,
      draftList: []
    };
  },
  onLoad() {
    this.statusBarHeight = utils_deviceInfo.deviceInfo.getStatusBarHeight();
    this.loadDraftList();
  },
  methods: {
    // 返回上一页
    handleBack() {
      common_vendor.index.navigateBack();
    },
    // 加载草稿列表
    loadDraftList() {
      this.draftList = utils_storage.getAllDraftList();
    },
    // 格式化时间显示
    formatTime(timeStr) {
      if (!timeStr)
        return "未知时间";
      const date = new Date(timeStr);
      const now = /* @__PURE__ */ new Date();
      const diff = now - date;
      if (diff < 36e5) {
        const minutes = Math.floor(diff / 6e4);
        return `${minutes === 0 ? "刚刚" : minutes + "分钟前"}`;
      }
      if (diff < 864e5) {
        const hours = Math.floor(diff / 36e5);
        return `${hours}小时前`;
      }
      if (diff < 6048e5) {
        const days = Math.floor(diff / 864e5);
        return `${days}天前`;
      }
      return `${date.getMonth() + 1}月${date.getDate()}日`;
    },
    // 获取标题预览（如果没有标题则显示描述的一部分）
    getTitlePreview(draft) {
      if (draft.title && draft.title.trim()) {
        return draft.title;
      }
      if (draft.description && draft.description.trim()) {
        const preview = draft.description.slice(0, 20);
        return preview.length >= 20 ? preview + "..." : preview;
      }
      return draft.type === "idle" ? "闲置物品草稿" : "信息发布草稿";
    },
    // 处理草稿点击，根据类型跳转到不同页面
    handleDraftClick(draft) {
      if (draft.type === "idle") {
        common_vendor.index.navigateTo({
          url: `/pages/publish/idle/index?draftId=${draft.id}`
        });
      } else {
        common_vendor.index.navigateTo({
          url: `/pages/publish/info/index?draftId=${draft.id}`
        });
      }
    },
    // 处理删除草稿
    handleDeleteDraft(draft) {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要删除该草稿吗？",
        success: (res) => {
          if (res.confirm) {
            let success = false;
            if (draft.type === "idle") {
              success = utils_storage.deleteIdleDraft(draft.id);
            } else {
              success = utils_storage.deleteInfoDraft(draft.id);
            }
            if (success) {
              common_vendor.index.showToast({
                title: "删除成功",
                icon: "success"
              });
              this.loadDraftList();
            } else {
              common_vendor.index.showToast({
                title: "删除失败",
                icon: "none"
              });
            }
          }
        }
      });
    },
    // 清空所有草稿
    handleClearAll() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要清空所有草稿吗？",
        success: (res) => {
          if (res.confirm) {
            this.draftList.forEach((draft) => {
              if (draft.type === "idle") {
                utils_storage.deleteIdleDraft(draft.id);
              } else {
                utils_storage.deleteInfoDraft(draft.id);
              }
            });
            this.draftList = [];
            common_vendor.index.showToast({
              title: "清空成功",
              icon: "success"
            });
          }
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
    c: $data.draftList.length > 0
  }, $data.draftList.length > 0 ? {
    d: common_vendor.o((...args) => $options.handleClearAll && $options.handleClearAll(...args))
  } : {}, {
    e: $data.navBarHeight + "px",
    f: $data.statusBarHeight + "px",
    g: `calc(${$data.statusBarHeight}px + ${$data.navBarHeight}px)`,
    h: $data.draftList.length > 0
  }, $data.draftList.length > 0 ? {
    i: common_vendor.f($data.draftList, (item, k0, i0) => {
      return {
        a: common_vendor.t($options.getTitlePreview(item)),
        b: common_vendor.t(item.type === "idle" ? "闲置" : "信息"),
        c: common_vendor.t($options.formatTime(item.updateTime)),
        d: "481b4036-1-" + i0,
        e: common_vendor.o(($event) => $options.handleDeleteDraft(item), item.id),
        f: item.id,
        g: common_vendor.o(($event) => $options.handleDraftClick(item), item.id)
      };
    }),
    j: common_vendor.p({
      type: "trash",
      size: "18",
      color: "#ff4500"
    })
  } : {
    k: common_assets._imports_0$2
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/publish/draft/index.js.map
