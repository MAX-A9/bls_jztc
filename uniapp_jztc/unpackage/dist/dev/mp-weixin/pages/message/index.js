"use strict";
const common_vendor = require("../../common/vendor.js");
const mixins_deviceAdapter = require("../../mixins/device-adapter.js");
const mixins_share = require("../../mixins/share.js");
const apis_index = require("../../apis/index.js");
const utils_messagePolling = require("../../utils/message-polling.js");
const common_assets = require("../../common/assets.js");
const TabBar = () => "../../components/tab-bar/index.js";
const _sfc_main = {
  components: {
    TabBar
  },
  mixins: [mixins_deviceAdapter.deviceAdapter, mixins_share.shareMixin],
  data() {
    return {
      tabIndex: 2,
      isRefreshing: false,
      totalUnreadCount: 0,
      // 系统消息 - 保留但不显示
      systemMessage: {
        lastMessage: "暂无系统通知",
        time: ""
      },
      chatList: [],
      pagination: {
        page: 1,
        size: 20,
        totalCount: 0,
        totalPage: 0
      },
      swipeOptions: [
        {
          text: "已读",
          style: {
            backgroundColor: "#8F8F94"
          }
        },
        {
          text: "删除",
          style: {
            backgroundColor: "#dd524d"
          }
        }
      ],
      conversationUnsubscribe: null,
      // 用于存储取消订阅函数
      _readSessionIds: {}
      // 用于记录已标记为已读的会话ID
    };
  },
  onLoad() {
    this.fetchSessionList();
    this.startConversationPolling();
  },
  onShow() {
    this.tabIndex = 3;
    this.startConversationPolling();
  },
  onHide() {
    this.stopConversationPolling();
  },
  onUnload() {
    this.stopConversationPolling();
  },
  methods: {
    // 开始会话列表轮询
    startConversationPolling() {
      this.stopConversationPolling();
      utils_messagePolling.messagePollingService.setPollingInterval(2e3);
      this.conversationUnsubscribe = utils_messagePolling.messagePollingService.setConversationParams(
        {
          page: this.pagination.page,
          size: this.pagination.size
        },
        this.handleConversationUpdate
      );
    },
    // 停止会话列表轮询
    stopConversationPolling() {
      if (this.conversationUnsubscribe) {
        this.conversationUnsubscribe();
        this.conversationUnsubscribe = null;
        utils_messagePolling.messagePollingService.setPollingInterval(5e3);
      }
    },
    // 处理会话列表更新
    handleConversationUpdate(data) {
      if (!data)
        return;
      const { list, totalCount, totalPage, currentPage, size } = data;
      this.pagination = {
        page: currentPage,
        size,
        totalCount,
        totalPage
      };
      if (list && list.length > 0) {
        this.chatList = list.map((item) => {
          return {
            ...item,
            time: this.formatTime(item.lastTime)
          };
        });
      } else {
        this.chatList = [];
      }
      this.chatList.forEach((item) => {
        const readItem = this._readSessionIds && this._readSessionIds[item.id];
        if (readItem) {
          item.unreadCount = 0;
        }
      });
    },
    // 处理滑动操作点击
    async handleSwipeClick(e, item) {
      if (e.index === 0) {
        await this.markAsRead(item);
      } else if (e.index === 1) {
        await this.deleteSession(item);
      }
    },
    // 标记会话为已读
    async markAsRead(item) {
      if (!item.unreadCount || item.unreadCount <= 0) {
        common_vendor.index.showToast({
          title: "没有未读消息",
          icon: "none"
        });
        return;
      }
      try {
        common_vendor.index.showLoading({
          title: "处理中...",
          mask: true
        });
        const response = await apis_index.message.markMessageRead({
          targetId: item.targetId
        });
        if (response && response.code === 0) {
          const index = this.chatList.findIndex((chat) => chat.id === item.id);
          if (index !== -1) {
            const previousUnread = this.chatList[index].unreadCount || 0;
            this.chatList[index].unreadCount = 0;
            this._readSessionIds[item.id] = true;
            if (previousUnread > 0) {
              this.totalUnreadCount = Math.max(0, this.totalUnreadCount - previousUnread);
            }
          }
          common_vendor.index.showToast({
            title: "已标记为已读",
            icon: "success"
          });
        } else {
          throw new Error((response == null ? void 0 : response.message) || "标记为已读失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/message/index.vue:282", "标记为已读失败:", error);
        common_vendor.index.showToast({
          title: "标记为已读失败",
          icon: "none"
        });
      } finally {
        common_vendor.index.hideLoading();
      }
    },
    // 删除会话
    async deleteSession(item) {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要删除此会话吗？",
        success: async (res) => {
          if (res.confirm) {
            try {
              common_vendor.index.showLoading({
                title: "处理中...",
                mask: true
              });
              const response = await apis_index.message.deleteConversation(item.id);
              common_vendor.index.hideLoading();
              if (response && response.code === 0 && response.data && response.data.success) {
                const index = this.chatList.findIndex((chat) => chat.id === item.id);
                if (index !== -1) {
                  const previousUnread = this.chatList[index].unreadCount || 0;
                  if (previousUnread > 0) {
                    this.totalUnreadCount = Math.max(0, this.totalUnreadCount - previousUnread);
                  }
                  this.chatList.splice(index, 1);
                }
                common_vendor.index.showToast({
                  title: "已删除会话",
                  icon: "success"
                });
              } else {
                throw new Error((response == null ? void 0 : response.message) || "删除会话失败");
              }
            } catch (error) {
              common_vendor.index.hideLoading();
              common_vendor.index.__f__("error", "at pages/message/index.vue:332", "删除会话失败:", error);
              common_vendor.index.showToast({
                title: "删除会话失败",
                icon: "none"
              });
            }
          }
        }
      });
    },
    // 获取聊天会话列表
    async fetchSessionList() {
      try {
        common_vendor.index.showLoading({
          title: "加载中...",
          mask: true
        });
        const res = await apis_index.message.getConversationList({
          page: this.pagination.page,
          size: this.pagination.size
        });
        if (res && res.code === 0) {
          const { list, totalCount, totalPage, currentPage, size } = res.data;
          this.pagination = {
            page: currentPage,
            size,
            totalCount,
            totalPage
          };
          utils_messagePolling.messagePollingService.updateConversationParams({
            page: currentPage,
            size
          });
          if (list && list.length > 0) {
            this.chatList = list.map((item) => {
              return {
                ...item,
                time: this.formatTime(item.lastTime)
              };
            });
          } else {
            this.chatList = [];
          }
        } else {
          throw new Error((res == null ? void 0 : res.message) || "获取会话列表失败");
        }
        common_vendor.index.hideLoading();
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/message/index.vue:391", "获取聊天会话列表失败:", error);
        common_vendor.index.showToast({
          title: "获取会话列表失败",
          icon: "none"
        });
        common_vendor.index.hideLoading();
      }
    },
    // 清除所有已读消息
    handleClearRead() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要清除所有已读消息吗？",
        success: async (res) => {
          if (res.confirm) {
            common_vendor.index.showLoading({
              title: "处理中...",
              mask: true
            });
            try {
              const hasReadItems = this.chatList.filter((item) => !item.unreadCount || item.unreadCount === 0);
              if (hasReadItems.length === 0) {
                common_vendor.index.hideLoading();
                common_vendor.index.showToast({
                  title: "没有可清除的已读消息",
                  icon: "none"
                });
                return;
              }
              let successCount = 0;
              for (const item of hasReadItems) {
                try {
                  const response = await apis_index.message.deleteConversation(item.id);
                  if (response && response.code === 0 && response.data && response.data.success) {
                    successCount++;
                  }
                } catch (err) {
                  common_vendor.index.__f__("error", "at pages/message/index.vue:435", `删除会话 ${item.id} 失败:`, err);
                }
              }
              this.chatList = this.chatList.filter((item) => item.unreadCount && item.unreadCount > 0);
              common_vendor.index.hideLoading();
              common_vendor.index.showToast({
                title: `已清除${successCount}个已读会话`,
                icon: "success"
              });
            } catch (error) {
              common_vendor.index.hideLoading();
              common_vendor.index.__f__("error", "at pages/message/index.vue:449", "清除已读消息失败:", error);
              common_vendor.index.showToast({
                title: "清除已读消息失败",
                icon: "none"
              });
            }
          }
        }
      });
    },
    // 聊天项点击处理
    async handleChatClick(chat) {
      common_vendor.index.navigateTo({
        url: `/pages/chat/detail?id=${chat.targetId}&nickName=${encodeURIComponent(chat.targetName || "用户")}`,
        success: async () => {
          if (chat.unreadCount && chat.unreadCount > 0) {
            try {
              const response = await apis_index.message.clearSessionUnread(chat.id);
              if (response && response.code === 0) {
                const index = this.chatList.findIndex((item) => item.id === chat.id);
                if (index !== -1) {
                  const previousUnread = this.chatList[index].unreadCount || 0;
                  this.chatList[index].unreadCount = 0;
                  this._readSessionIds[chat.id] = true;
                  if (previousUnread > 0) {
                    this.totalUnreadCount = Math.max(0, this.totalUnreadCount - previousUnread);
                  }
                }
              }
            } catch (error) {
              common_vendor.index.__f__("error", "at pages/message/index.vue:491", "清除会话未读消息失败:", error);
            }
          }
        }
      });
    },
    // 处理系统消息点击 - 已注释
    /*
    handleSystemClick() {
    	uni.navigateTo({
    		url: '/pages/notification/index'
    	})
    },
    */
    // 下拉刷新处理
    async onRefresh() {
      if (this.isRefreshing)
        return;
      this.isRefreshing = true;
      try {
        this.pagination.page = 1;
        await this.fetchSessionList();
        common_vendor.index.showToast({
          title: "刷新成功",
          icon: "success"
        });
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/message/index.vue:524", "刷新失败:", error);
        common_vendor.index.showToast({
          title: "刷新失败",
          icon: "none"
        });
      } finally {
        this.isRefreshing = false;
      }
    },
    // 刷新复位
    onRestore() {
      common_vendor.index.__f__("log", "at pages/message/index.vue:536", "刷新复位");
    },
    // 格式化时间显示
    formatTime(timestamp) {
      if (!timestamp)
        return "";
      let msgDate;
      if (typeof timestamp === "string") {
        const formattedTimestamp = timestamp.replace(/-/g, "/");
        msgDate = new Date(formattedTimestamp);
      } else {
        msgDate = new Date(timestamp);
      }
      if (isNaN(msgDate.getTime())) {
        common_vendor.index.__f__("error", "at pages/message/index.vue:555", "Invalid date format:", timestamp);
        return "";
      }
      const now = /* @__PURE__ */ new Date();
      const diffDays = Math.floor((now - msgDate) / (24 * 60 * 60 * 1e3));
      const hours = msgDate.getHours().toString().padStart(2, "0");
      const minutes = msgDate.getMinutes().toString().padStart(2, "0");
      const timeStr = `${hours}:${minutes}`;
      if (diffDays === 0) {
        return timeStr;
      }
      if (diffDays === 1) {
        return "昨天";
      }
      if (diffDays < 7) {
        const weekdays = ["周日", "周一", "周二", "周三", "周四", "周五", "周六"];
        return weekdays[msgDate.getDay()];
      }
      const month = (msgDate.getMonth() + 1).toString().padStart(2, "0");
      const day = msgDate.getDate().toString().padStart(2, "0");
      return `${month}-${day}`;
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _easycom_uni_swipe_action_item2 = common_vendor.resolveComponent("uni-swipe-action-item");
  const _easycom_uni_swipe_action2 = common_vendor.resolveComponent("uni-swipe-action");
  const _component_tab_bar = common_vendor.resolveComponent("tab-bar");
  (_easycom_uni_icons2 + _easycom_uni_swipe_action_item2 + _easycom_uni_swipe_action2 + _component_tab_bar)();
}
const _easycom_uni_icons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
const _easycom_uni_swipe_action_item = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-swipe-action-item/uni-swipe-action-item.js";
const _easycom_uni_swipe_action = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-swipe-action/uni-swipe-action.js";
if (!Math) {
  (_easycom_uni_icons + _easycom_uni_swipe_action_item + _easycom_uni_swipe_action)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.p({
      type: "trash",
      size: "16",
      color: "#666666"
    }),
    b: common_vendor.o((...args) => $options.handleClearRead && $options.handleClearRead(...args)),
    c: common_vendor.f($data.chatList, (item, index, i0) => {
      return common_vendor.e({
        a: item.targetAvatar || "/static/demo/0.png",
        b: common_vendor.t(item.targetName || "用户"),
        c: common_vendor.t(item.time || ""),
        d: common_vendor.t(item.lastMessage || ""),
        e: item.unreadCount && item.unreadCount > 0
      }, item.unreadCount && item.unreadCount > 0 ? {
        f: common_vendor.t(item.unreadCount)
      } : {}, {
        g: common_vendor.o(($event) => $options.handleChatClick(item), item.id || index),
        h: common_vendor.o((e) => $options.handleSwipeClick(e, item), item.id || index),
        i: "4af1779c-2-" + i0 + "," + ("4af1779c-1-" + i0),
        j: "4af1779c-1-" + i0,
        k: item.id || index
      });
    }),
    d: common_vendor.p({
      ["right-options"]: $data.swipeOptions
    }),
    e: $data.chatList.length === 0
  }, $data.chatList.length === 0 ? {
    f: common_assets._imports_0$1
  } : {}, {
    g: $data.isRefreshing,
    h: common_vendor.o((...args) => $options.onRefresh && $options.onRefresh(...args)),
    i: common_vendor.o((...args) => $options.onRestore && $options.onRestore(...args)),
    j: common_vendor.p({
      ["current-tab"]: $data.tabIndex
    })
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/message/index.js.map
