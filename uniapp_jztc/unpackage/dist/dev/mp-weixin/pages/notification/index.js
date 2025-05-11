"use strict";
const common_vendor = require("../../common/vendor.js");
const common_assets = require("../../common/assets.js");
const _sfc_main = {
  data() {
    return {
      isRefreshing: false,
      notificationList: [
        {
          id: 1,
          type: "system",
          title: "系统更新通知",
          content: "亲爱的用户，我们的应用已更新到最新版本，新增了多项实用功能，欢迎体验！",
          time: "今天 12:30",
          isRead: false
        },
        {
          id: 2,
          type: "comment",
          title: "评论通知",
          content: '用户"张先生"评论了您的闲置物品："这个价格很合理，我很感兴趣"',
          time: "今天 10:15",
          isRead: true
        },
        {
          id: 3,
          type: "like",
          title: "点赞通知",
          content: '您的闲置物品"iPhone 12 128G 蓝色"收到了3个新的点赞',
          time: "昨天 15:40",
          isRead: false
        },
        {
          id: 4,
          type: "order",
          title: "订单通知",
          content: "您的订单#20230615001已发货，预计3天内送达",
          time: "昨天 09:20",
          isRead: true
        },
        {
          id: 5,
          type: "system",
          title: "账号安全提醒",
          content: "您的账号于6月10日在新设备上登录，如非本人操作，请及时修改密码",
          time: "3天前",
          isRead: true
        },
        {
          id: 6,
          type: "comment",
          title: "评论通知",
          content: '用户"李女士"回复了您的评论："好的，那就这个价格成交"',
          time: "4天前",
          isRead: true
        },
        {
          id: 7,
          type: "like",
          title: "点赞通知",
          content: '您的闲置物品"九成新小米平板5 Pro"收到了5个新的点赞',
          time: "5天前",
          isRead: true
        },
        {
          id: 8,
          type: "order",
          title: "订单通知",
          content: "您的订单#20230608002已完成，感谢您的使用",
          time: "一周前",
          isRead: true
        }
      ]
    };
  },
  methods: {
    getIconByType(type) {
      const iconMap = {
        "system": "notification-filled",
        "comment": "chat-filled",
        "like": "heart-filled",
        "order": "cart-filled"
      };
      return iconMap[type] || "notification-filled";
    },
    handleClearAll() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要清空所有通知吗？",
        success: (res) => {
          if (res.confirm) {
            this.notificationList = [];
            common_vendor.index.showToast({
              title: "已清空所有通知",
              icon: "success"
            });
          }
        }
      });
    },
    handleNotificationClick(notification) {
      const index = this.notificationList.findIndex((item) => item.id === notification.id);
      if (index !== -1 && !this.notificationList[index].isRead) {
        this.notificationList[index].isRead = true;
      }
      switch (notification.type) {
        case "comment":
          common_vendor.index.navigateTo({
            url: `/pages/content/detail?id=${notification.id}&type=comment`
          });
          break;
        case "like":
          common_vendor.index.navigateTo({
            url: `/pages/my/publish/index?highlight=${notification.id}`
          });
          break;
        case "order":
          common_vendor.index.navigateTo({
            url: `/pages/my/orders/index?id=${notification.id}`
          });
          break;
        case "system":
        default:
          common_vendor.index.showModal({
            title: notification.title,
            content: notification.content,
            showCancel: false,
            confirmText: "知道了"
          });
          break;
      }
    },
    async onRefresh() {
      if (this.isRefreshing)
        return;
      this.isRefreshing = true;
      try {
        await new Promise((resolve) => setTimeout(resolve, 1e3));
        this.notificationList = [
          {
            id: 9,
            type: "system",
            title: "新活动通知",
            content: "618年中大促即将开始，多款商品低至5折，先到先得！",
            time: "刚刚",
            isRead: false
          },
          ...this.notificationList
        ];
        common_vendor.index.showToast({
          title: "刷新成功",
          icon: "success"
        });
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/notification/index.vue:226", "刷新失败:", error);
        common_vendor.index.showToast({
          title: "刷新失败",
          icon: "error"
        });
      } finally {
        this.isRefreshing = false;
      }
    },
    onRestore() {
      common_vendor.index.__f__("log", "at pages/notification/index.vue:237", "刷新复位");
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
      type: "trash",
      size: "16",
      color: "#666666"
    }),
    b: common_vendor.o((...args) => $options.handleClearAll && $options.handleClearAll(...args)),
    c: $data.notificationList.length === 0
  }, $data.notificationList.length === 0 ? {
    d: common_assets._imports_0$7
  } : {
    e: common_vendor.f($data.notificationList, (item, index, i0) => {
      return common_vendor.e({
        a: "f06be568-1-" + i0,
        b: common_vendor.p({
          type: $options.getIconByType(item.type),
          size: "24",
          color: "#ffffff"
        }),
        c: common_vendor.n(item.type),
        d: common_vendor.t(item.title),
        e: common_vendor.t(item.time),
        f: common_vendor.t(item.content),
        g: !item.isRead
      }, !item.isRead ? {} : {}, {
        h: index,
        i: common_vendor.o(($event) => $options.handleNotificationClick(item), index)
      });
    })
  }, {
    f: $data.isRefreshing,
    g: common_vendor.o((...args) => $options.onRefresh && $options.onRefresh(...args)),
    h: common_vendor.o((...args) => $options.onRestore && $options.onRestore(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/notification/index.js.map
