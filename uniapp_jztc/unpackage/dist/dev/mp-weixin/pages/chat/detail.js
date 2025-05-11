"use strict";
const common_vendor = require("../../common/vendor.js");
const apis_message = require("../../apis/message.js");
const utils_messagePolling = require("../../utils/message-polling.js");
const _sfc_main = {
  data() {
    return {
      chatId: null,
      chatInfo: {
        id: 0,
        name: "",
        avatar: ""
      },
      messageList: [],
      inputMessage: "",
      showMorePanel: false,
      inputFocus: false,
      currentDate: this.formatDate(/* @__PURE__ */ new Date()),
      keyboardVisible: false,
      userAvatar: "",
      pagination: {
        page: 1,
        size: 20,
        totalCount: 0,
        totalPage: 0
      },
      isLoading: false,
      nickName: "",
      // 聊天对象昵称
      _markingAsRead: false,
      chatMessageUnsubscribe: null
      // 用于存储取消订阅函数
    };
  },
  onLoad(options) {
    if (options.id || options.sessionId) {
      this.chatId = options.id || options.sessionId;
      if (options.nickName) {
        this.nickName = decodeURIComponent(options.nickName);
      }
      this.loadChatInfo();
      this.loadChatHistory();
      const userInfo = common_vendor.index.getStorageSync("USER_INFO") || {};
      this.userAvatar = userInfo.avatarUrl || "/static/demo/0.png";
      utils_messagePolling.messagePollingService.setCurrentPage("/pages/chat/detail");
    }
  },
  onReady() {
    if (this.chatInfo.name) {
      common_vendor.index.setNavigationBarTitle({
        title: this.chatInfo.name
      });
    }
  },
  // 监听键盘高度变化
  onKeyboardHeightChange(e) {
    this.keyboardVisible = e.height > 0;
    if (this.keyboardVisible) {
      this.showMorePanel = false;
    }
  },
  // 页面显示时标记消息为已读
  onShow() {
    if (this.chatId) {
      common_vendor.index.__f__("log", "at pages/chat/detail.vue:208", "聊天详情页面显示，chatId:", this.chatId);
      utils_messagePolling.messagePollingService.setCurrentPage("/pages/chat/detail");
      if (this.messageList.length > 0) {
        this.markMessagesAsRead();
      }
      this.startChatPolling();
    }
  },
  // 显示导航栏右侧按钮
  onNavigationBarButtonTap(e) {
    if (e.index === 0) {
      this.showMoreOptions();
    }
  },
  onHide() {
    common_vendor.index.__f__("log", "at pages/chat/detail.vue:230", "聊天详情页面隐藏");
  },
  onUnload() {
    common_vendor.index.__f__("log", "at pages/chat/detail.vue:234", "聊天详情页面卸载");
    this.stopChatPolling();
    utils_messagePolling.messagePollingService.setCurrentPage("");
  },
  methods: {
    // 加载聊天信息
    loadChatInfo() {
      if (this.nickName) {
        this.chatInfo = {
          id: this.chatId,
          name: this.nickName,
          avatar: "/static/demo/0.png"
          // 使用默认头像，后续会从消息中获取并更新
        };
        common_vendor.index.setNavigationBarTitle({
          title: this.chatInfo.name
        });
      } else {
        this.chatInfo = {
          id: this.chatId,
          name: "用户" + this.chatId,
          avatar: "/static/demo/0.png"
          // 使用默认头像，后续会从消息中获取并更新
        };
        common_vendor.index.setNavigationBarTitle({
          title: this.chatInfo.name
        });
      }
    },
    // 加载聊天历史
    async loadChatHistory() {
      if (!this.chatId || this.isLoading)
        return;
      this.isLoading = true;
      try {
        common_vendor.index.showLoading({
          title: "加载中...",
          mask: true
        });
        const res = await apis_message.getMessageList({
          targetId: this.chatId,
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
          if (list && list.length > 0) {
            const formattedMessages = list.map((item) => {
              const formattedDate = item.createdAt ? item.createdAt.replace(/-/g, "/") : "";
              const createdAtTime = formattedDate ? new Date(formattedDate).getTime() : 0;
              return {
                type: "text",
                // 默认都是文本消息
                content: item.content,
                time: this.formatMessageTime(item.createdAt),
                isSelf: item.isSelf,
                id: item.id,
                senderId: item.senderId,
                senderName: item.senderName,
                senderAvatar: item.senderAvatar || "/static/demo/0.png",
                receiverId: item.receiverId,
                receiverName: item.receiverName,
                receiverAvatar: item.receiverAvatar || "/static/demo/0.png",
                isRead: item.isRead,
                createdAt: createdAtTime
                // 用于排序
              };
            });
            formattedMessages.sort((a, b) => a.id - b.id);
            this.messageList = formattedMessages;
            this.markMessagesAsRead();
          } else {
            this.messageList = [];
          }
        } else {
          throw new Error((res == null ? void 0 : res.message) || "获取消息列表失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/chat/detail.vue:337", "加载聊天记录失败:", error);
        common_vendor.index.showToast({
          title: "加载聊天记录失败",
          icon: "none"
        });
      } finally {
        this.isLoading = false;
        common_vendor.index.hideLoading();
      }
    },
    // 格式化消息时间
    formatMessageTime(timestamp) {
      if (!timestamp)
        return "";
      let dateObj;
      if (typeof timestamp === "string") {
        const formattedTimestamp = timestamp.replace(/-/g, "/");
        dateObj = new Date(formattedTimestamp);
      } else {
        dateObj = new Date(timestamp);
      }
      if (isNaN(dateObj.getTime())) {
        common_vendor.index.__f__("error", "at pages/chat/detail.vue:364", "Invalid date format:", timestamp);
        return "";
      }
      return this.formatTime(dateObj);
    },
    // 显示更多选项
    showMoreOptions() {
      common_vendor.index.showActionSheet({
        itemList: ["清空聊天记录", "投诉", "加入黑名单"],
        success: (res) => {
          switch (res.tapIndex) {
            case 0:
              this.clearChatHistory();
              break;
            case 1:
              this.reportUser();
              break;
            case 2:
              this.blockUser();
              break;
          }
        }
      });
    },
    // 清空聊天记录
    clearChatHistory() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要清空聊天记录吗？",
        success: (res) => {
          if (res.confirm) {
            this.messageList = [];
            common_vendor.index.showToast({
              title: "已清空聊天记录",
              icon: "success"
            });
          }
        }
      });
    },
    // 投诉用户
    reportUser() {
      common_vendor.index.showToast({
        title: "已提交投诉",
        icon: "success"
      });
    },
    // 拉黑用户
    blockUser() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要将该用户加入黑名单吗？",
        success: (res) => {
          if (res.confirm) {
            common_vendor.index.showToast({
              title: "已加入黑名单",
              icon: "success"
            });
            setTimeout(() => {
              common_vendor.index.navigateBack();
            }, 1500);
          }
        }
      });
    },
    // 聊天消息处理函数
    handleChatMessages(data) {
      if (!data || !data.list)
        return;
      const { list } = data;
      if (list && list.length > 0) {
        const newMessages = list.filter((item) => {
          const exists = this.messageList.some((existing) => existing.id === item.id);
          return !exists;
        }).map((item) => {
          const formattedDate = item.createdAt ? item.createdAt.replace(/-/g, "/") : "";
          const createdAtTime = formattedDate ? new Date(formattedDate).getTime() : 0;
          return {
            type: "text",
            // 默认都是文本消息
            content: item.content,
            time: this.formatMessageTime(item.createdAt),
            isSelf: item.isSelf,
            id: item.id,
            senderId: item.senderId,
            senderName: item.senderName,
            senderAvatar: item.senderAvatar || "/static/demo/0.png",
            receiverId: item.receiverId,
            receiverName: item.receiverName,
            receiverAvatar: item.receiverAvatar || "/static/demo/0.png",
            isRead: item.isRead,
            createdAt: createdAtTime
            // 用于排序
          };
        });
        if (newMessages.length > 0) {
          newMessages.sort((a, b) => a.id - b.id);
          this.messageList = [...this.messageList, ...newMessages];
          this.markMessagesAsRead();
        }
      }
    },
    // 开始聊天轮询
    startChatPolling() {
      if (!this.chatId)
        return;
      let lastMsgId = 0;
      if (this.messageList.length > 0) {
        const lastMsg = this.messageList[this.messageList.length - 1];
        lastMsgId = lastMsg.id || 0;
      }
      this.stopChatPolling();
      this.chatMessageUnsubscribe = utils_messagePolling.messagePollingService.setChatParams(
        {
          targetId: this.chatId,
          lastId: lastMsgId
        },
        this.handleChatMessages
      );
    },
    // 停止聊天轮询
    stopChatPolling() {
      if (this.chatMessageUnsubscribe) {
        this.chatMessageUnsubscribe();
        this.chatMessageUnsubscribe = null;
      }
    },
    // 发送文本消息成功后更新lastId
    updateLastId(messageId) {
      if (messageId && this.chatId) {
        utils_messagePolling.messagePollingService.updateChatParams({ lastId: messageId });
      }
    },
    // 发送文本消息
    async sendTextMessage() {
      if (!this.inputMessage.trim())
        return;
      const messageContent = this.inputMessage.trim();
      this.inputMessage = "";
      try {
        const now = /* @__PURE__ */ new Date();
        const tempMessage = {
          type: "text",
          content: messageContent,
          time: this.formatTime(now),
          isSelf: true,
          sending: true,
          // 标记为发送中
          createdAt: now.getTime(),
          // 用于排序
          senderAvatar: this.userAvatar || "/static/demo/0.png"
        };
        this.messageList.push(tempMessage);
        this.showMorePanel = false;
        const response = await apis_message.sendMessage({
          content: messageContent,
          receiverId: parseInt(this.chatId)
        });
        if (response && response.code === 0) {
          const lastIndex = this.messageList.length - 1;
          if (lastIndex >= 0) {
            this.messageList[lastIndex].sending = false;
            if (response.data && response.data.id) {
              this.messageList[lastIndex].id = response.data.id;
              this.updateLastId(response.data.id);
            }
          }
        } else {
          throw new Error((response == null ? void 0 : response.message) || "发送消息失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/chat/detail.vue:579", "发送消息失败:", error);
        const lastIndex = this.messageList.length - 1;
        if (lastIndex >= 0 && this.messageList[lastIndex].sending) {
          this.messageList[lastIndex].sending = false;
          this.messageList[lastIndex].failed = true;
        }
        common_vendor.index.showToast({
          title: "发送失败，请重试",
          icon: "none"
        });
      }
    },
    // 切换更多功能面板
    toggleMorePanel() {
      this.showMorePanel = !this.showMorePanel;
      this.inputFocus = false;
    },
    // 隐藏更多功能面板
    hideMorePanel() {
      if (this.showMorePanel) {
        this.showMorePanel = false;
      }
    },
    // 聚焦输入框
    focusInput() {
      this.inputFocus = true;
      this.hideMorePanel();
    },
    // 输入框获得焦点
    onInputFocus(e) {
      this.keyboardVisible = true;
    },
    // 输入框失去焦点
    onInputBlur() {
      setTimeout(() => {
        this.keyboardVisible = false;
      }, 100);
    },
    // 选择图片
    chooseImage() {
      common_vendor.index.chooseImage({
        count: 1,
        success: (res) => {
          const tempFilePath = res.tempFilePaths[0];
          const now = /* @__PURE__ */ new Date();
          this.messageList.push({
            type: "image",
            content: tempFilePath,
            time: this.formatTime(now),
            isSelf: true,
            sending: true,
            createdAt: now.getTime()
            // 用于排序
          });
          this.showMorePanel = false;
          common_vendor.index.showToast({
            title: "图片功能开发中",
            icon: "none"
          });
        }
      });
    },
    // 选择位置
    chooseLocation() {
      common_vendor.index.showToast({
        title: "选择位置功能开发中",
        icon: "none"
      });
      this.showMorePanel = false;
    },
    // 分享商品
    shareProduct() {
      common_vendor.index.showToast({
        title: "分享商品功能开发中",
        icon: "none"
      });
      this.showMorePanel = false;
    },
    // 预览图片
    previewImage(url) {
      const imageUrls = this.messageList.filter((msg) => msg.type === "image").map((msg) => msg.content);
      common_vendor.index.previewImage({
        current: url,
        urls: imageUrls
      });
    },
    // 格式化日期
    formatDate(date) {
      const year = date.getFullYear();
      const month = date.getMonth() + 1;
      const day = date.getDate();
      return `${year}年${month}月${day}日`;
    },
    // 格式化时间
    formatTime(date) {
      const hours = date.getHours().toString().padStart(2, "0");
      const minutes = date.getMinutes().toString().padStart(2, "0");
      return `${hours}:${minutes}`;
    },
    // 标记消息为已读
    async markMessagesAsRead() {
      if (!this.chatId) {
        common_vendor.index.__f__("error", "at pages/chat/detail.vue:709", "无法标记消息为已读：chatId不存在");
        return;
      }
      if (this._markingAsRead)
        return;
      this._markingAsRead = true;
      try {
        common_vendor.index.__f__("log", "at pages/chat/detail.vue:718", "标记消息为已读，targetId:", this.chatId);
        let targetId = this.chatId;
        if (typeof targetId === "string") {
          targetId = parseInt(targetId);
        }
        const response = await apis_message.markMessageRead({
          targetId
        });
        if (response && response.code === 0) {
          common_vendor.index.__f__("log", "at pages/chat/detail.vue:732", "消息已成功标记为已读");
          this.messageList.forEach((msg) => {
            if (!msg.isSelf) {
              msg.isRead = true;
            }
          });
        } else {
          common_vendor.index.__f__("error", "at pages/chat/detail.vue:741", "标记消息为已读失败:", (response == null ? void 0 : response.message) || "未知错误");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/chat/detail.vue:744", "标记消息为已读请求失败:", error);
      } finally {
        this._markingAsRead = false;
      }
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
    a: common_vendor.t($data.currentDate),
    b: common_vendor.f($data.messageList, (message, index, i0) => {
      return common_vendor.e({
        a: !message.isSelf
      }, !message.isSelf ? {
        b: message.senderAvatar
      } : {}, {
        c: message.type === "text"
      }, message.type === "text" ? {
        d: common_vendor.t(message.content)
      } : message.type === "image" ? {
        f: message.content,
        g: common_vendor.o(($event) => $options.previewImage(message.content), index)
      } : message.type === "product" ? {
        i: message.content.image,
        j: common_vendor.t(message.content.title),
        k: common_vendor.t(message.content.price)
      } : {}, {
        e: message.type === "image",
        h: message.type === "product",
        l: message.isSelf
      }, message.isSelf ? {
        m: message.senderAvatar
      } : {}, {
        n: index,
        o: common_vendor.n(message.isSelf ? "self" : "other"),
        p: "msg-" + index
      });
    }),
    c: "msg-" + ($data.messageList.length - 1),
    d: common_vendor.o((...args) => $options.hideMorePanel && $options.hideMorePanel(...args)),
    e: $data.inputFocus,
    f: common_vendor.o((...args) => $options.sendTextMessage && $options.sendTextMessage(...args)),
    g: common_vendor.o((...args) => $options.onInputFocus && $options.onInputFocus(...args)),
    h: common_vendor.o((...args) => $options.onInputBlur && $options.onInputBlur(...args)),
    i: $data.inputMessage,
    j: common_vendor.o(($event) => $data.inputMessage = $event.detail.value),
    k: common_vendor.o((...args) => $options.focusInput && $options.focusInput(...args)),
    l: $data.inputMessage.trim()
  }, $data.inputMessage.trim() ? {} : {
    m: common_vendor.p({
      type: "plusempty",
      size: "24",
      color: "#333333"
    })
  }, {
    n: $data.inputMessage.trim() ? 1 : "",
    o: common_vendor.o(($event) => $data.inputMessage.trim() ? $options.sendTextMessage() : $options.toggleMorePanel()),
    p: $data.showMorePanel && !$data.keyboardVisible
  }, $data.showMorePanel && !$data.keyboardVisible ? {
    q: common_vendor.p({
      type: "image",
      size: "24",
      color: "#ffffff"
    }),
    r: common_vendor.o((...args) => $options.chooseImage && $options.chooseImage(...args)),
    s: common_vendor.p({
      type: "location",
      size: "24",
      color: "#ffffff"
    }),
    t: common_vendor.o((...args) => $options.chooseLocation && $options.chooseLocation(...args)),
    v: common_vendor.p({
      type: "shop",
      size: "24",
      color: "#ffffff"
    }),
    w: common_vendor.o((...args) => $options.shareProduct && $options.shareProduct(...args)),
    x: common_vendor.o(() => {
    })
  } : {}, {
    y: common_vendor.o((...args) => $options.hideMorePanel && $options.hideMorePanel(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/chat/detail.js.map
