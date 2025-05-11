"use strict";
const common_vendor = require("../common/vendor.js");
const apis_message = require("../apis/message.js");
class MessagePollingService {
  constructor() {
    this.pollingTimer = null;
    this.pollingInterval = 5e3;
    this.chatPollingInterval = 3e3;
    this.isPolling = false;
    this.listeners = {
      unreadCount: [],
      // 未读消息数监听器
      chatMessages: [],
      // 聊天消息监听器
      conversationList: []
      // 会话列表监听器
    };
    this.chatParams = null;
    this.conversationParams = null;
    this.globalPaused = false;
    this.previousPollingState = false;
    this.whitelistPages = ["/pages/chat/detail"];
    this.currentPage = "";
  }
  /**
   * 设置当前页面路径
   * @param {String} pagePath - 页面路径
   */
  setCurrentPage(pagePath) {
    this.currentPage = pagePath;
    common_vendor.index.__f__("log", "at utils/message-polling.js:32", "当前页面路径:", pagePath);
    if (this.isPageInWhitelist(pagePath) && this.globalPaused) {
      common_vendor.index.__f__("log", "at utils/message-polling.js:36", "当前页面在白名单中，恢复轮询");
      this.resumePollingForWhitelist();
    }
  }
  /**
   * 检查页面是否在白名单中
   * @param {String} pagePath - 页面路径
   * @returns {Boolean} 是否在白名单中
   */
  isPageInWhitelist(pagePath) {
    return this.whitelistPages.some((whitePath) => pagePath && pagePath.indexOf(whitePath) !== -1);
  }
  /**
   * 为白名单页面恢复轮询，但不改变全局暂停状态
   */
  resumePollingForWhitelist() {
    if (this.globalPaused && !this.isPolling && this.previousPollingState) {
      common_vendor.index.__f__("log", "at utils/message-polling.js:55", "白名单页面特殊处理：恢复轮询但保持全局暂停状态");
      this.isPolling = true;
      this.poll();
      const interval = this.chatParams ? this.chatPollingInterval : this.pollingInterval;
      this.pollingTimer = setInterval(() => {
        this.poll();
      }, interval);
    }
  }
  /**
   * 设置轮询间隔时间
   * @param {Number} interval - 轮询间隔(毫秒)
   */
  setPollingInterval(interval) {
    this.pollingInterval = interval;
    if (this.isPolling) {
      this.stopPolling();
      this.startPolling();
    }
  }
  /**
   * 设置聊天轮询间隔时间
   * @param {Number} interval - 轮询间隔(毫秒)
   */
  setChatPollingInterval(interval) {
    this.chatPollingInterval = interval;
    if (this.isPolling && this.chatParams) {
      this.stopPolling();
      this.startPolling();
    }
  }
  /**
   * 启动轮询
   */
  startPolling() {
    if (this.isPolling || this.globalPaused)
      return;
    this.isPolling = true;
    this.poll();
    const interval = this.chatParams ? this.chatPollingInterval : this.pollingInterval;
    this.pollingTimer = setInterval(() => {
      this.poll();
    }, interval);
    common_vendor.index.__f__("log", "at utils/message-polling.js:111", `消息轮询已启动，间隔 ${interval} ms`);
  }
  /**
   * 停止轮询
   */
  stopPolling() {
    if (this.pollingTimer) {
      clearInterval(this.pollingTimer);
      this.pollingTimer = null;
    }
    this.isPolling = false;
    common_vendor.index.__f__("log", "at utils/message-polling.js:123", "消息轮询已停止");
  }
  /**
   * 全局暂停轮询（适用于进入非底部导航页面时）
   * 此方法会记住当前的轮询状态，以便恢复时能够还原
   */
  pausePollingGlobally() {
    if (this.globalPaused)
      return;
    if (this.isPageInWhitelist(this.currentPage)) {
      common_vendor.index.__f__("log", "at utils/message-polling.js:135", "当前页面在白名单中，不执行全局暂停");
      return;
    }
    this.previousPollingState = this.isPolling;
    this.globalPaused = true;
    if (this.isPolling) {
      this.stopPolling();
      common_vendor.index.__f__("log", "at utils/message-polling.js:144", "消息轮询已全局暂停");
    }
  }
  /**
   * 全局恢复轮询（适用于返回到底部导航页面时）
   * 如果之前处于轮询状态，则会恢复轮询
   */
  resumePollingGlobally() {
    if (!this.globalPaused)
      return;
    this.globalPaused = false;
    if (this.previousPollingState) {
      this.startPolling();
      common_vendor.index.__f__("log", "at utils/message-polling.js:159", "消息轮询已全局恢复");
    }
  }
  /**
   * 执行轮询
   */
  async poll() {
    if (this.globalPaused && !this.isPageInWhitelist(this.currentPage))
      return;
    if (this.listeners.unreadCount.length > 0) {
      this.pollUnreadCount();
    }
    if (this.chatParams && this.listeners.chatMessages.length > 0) {
      this.pollChatMessages();
    }
    if (this.conversationParams && this.listeners.conversationList.length > 0) {
      this.pollConversationList();
    }
  }
  /**
   * 获取未读消息数
   */
  async pollUnreadCount() {
    try {
      const res = await apis_message.getUnreadCount();
      if (res && res.code === 0) {
        const unreadCount = res.data.unreadCount || 0;
        this.listeners.unreadCount.forEach((listener) => {
          listener(unreadCount);
        });
      }
    } catch (error) {
      common_vendor.index.__f__("error", "at utils/message-polling.js:200", "获取未读消息数失败:", error);
    }
  }
  /**
   * 获取聊天消息
   */
  async pollChatMessages() {
    if (!this.chatParams || !this.chatParams.targetId)
      return;
    try {
      const res = await apis_message.getMessageList({
        targetId: this.chatParams.targetId,
        page: 1,
        size: 20,
        lastId: this.chatParams.lastId || 0
      });
      if (res && res.code === 0) {
        this.listeners.chatMessages.forEach((listener) => {
          listener(res.data);
        });
      }
    } catch (error) {
      common_vendor.index.__f__("error", "at utils/message-polling.js:225", "获取聊天消息失败:", error);
    }
  }
  /**
   * 获取会话列表
   */
  async pollConversationList() {
    if (!this.conversationParams)
      return;
    try {
      const timestamp = (/* @__PURE__ */ new Date()).getTime();
      const res = await apis_message.getConversationList({
        page: this.conversationParams.page || 1,
        size: this.conversationParams.size || 20,
        _t: timestamp
        // 添加时间戳防止缓存
      });
      if (res && res.code === 0) {
        this.listeners.conversationList.forEach((listener) => {
          listener(res.data);
        });
      }
    } catch (error) {
      common_vendor.index.__f__("error", "at utils/message-polling.js:251", "获取会话列表失败:", error);
    }
  }
  /**
   * 添加未读消息数监听器
   * @param {Function} callback - 监听器回调函数
   * @returns {Function} 移除监听器的函数
   */
  addUnreadCountListener(callback) {
    if (typeof callback !== "function")
      return () => {
      };
    this.listeners.unreadCount.push(callback);
    if (!this.isPolling) {
      this.startPolling();
    } else {
      this.pollUnreadCount();
    }
    return () => {
      this.removeUnreadCountListener(callback);
    };
  }
  /**
   * 移除未读消息数监听器
   * @param {Function} callback - 要移除的监听器回调函数
   */
  removeUnreadCountListener(callback) {
    const index = this.listeners.unreadCount.indexOf(callback);
    if (index !== -1) {
      this.listeners.unreadCount.splice(index, 1);
    }
    this.checkAndStopPolling();
  }
  /**
   * 设置聊天参数并添加聊天消息监听器
   * @param {Object} params - 聊天参数
   * @param {Number} params.targetId - 聊天目标ID
   * @param {Number} params.lastId - 最后一条消息ID
   * @param {Function} callback - 监听器回调函数
   * @returns {Function} 移除监听器的函数
   */
  setChatParams(params, callback) {
    this.chatParams = params;
    if (typeof callback === "function") {
      this.listeners.chatMessages.push(callback);
      if (this.isPolling) {
        this.stopPolling();
      }
      this.startPolling();
      return () => {
        this.removeChatMessageListener(callback);
      };
    }
    return () => {
    };
  }
  /**
   * 更新聊天参数
   * @param {Object} params - 聊天参数
   */
  updateChatParams(params) {
    if (this.chatParams) {
      this.chatParams = { ...this.chatParams, ...params };
    } else {
      this.chatParams = params;
    }
  }
  /**
   * 移除聊天消息监听器
   * @param {Function} callback - 要移除的监听器回调函数
   */
  removeChatMessageListener(callback) {
    const index = this.listeners.chatMessages.indexOf(callback);
    if (index !== -1) {
      this.listeners.chatMessages.splice(index, 1);
    }
    if (this.listeners.chatMessages.length === 0) {
      this.chatParams = null;
      this.adjustPollingState();
    }
  }
  /**
   * 设置会话列表参数并添加会话列表监听器
   * @param {Object} params - 会话列表参数
   * @param {Number} params.page - 页码
   * @param {Number} params.size - 每页条数
   * @param {Function} callback - 监听器回调函数
   * @returns {Function} 移除监听器的函数
   */
  setConversationParams(params, callback) {
    this.conversationParams = params;
    if (typeof callback === "function") {
      this.listeners.conversationList.push(callback);
      if (!this.isPolling) {
        this.startPolling();
      } else {
        this.pollConversationList();
      }
      return () => {
        this.removeConversationListListener(callback);
      };
    }
    return () => {
    };
  }
  /**
   * 更新会话列表参数
   * @param {Object} params - 会话列表参数
   */
  updateConversationParams(params) {
    if (this.conversationParams) {
      this.conversationParams = { ...this.conversationParams, ...params };
    } else {
      this.conversationParams = params;
    }
  }
  /**
   * 移除会话列表监听器
   * @param {Function} callback - 要移除的监听器回调函数
   */
  removeConversationListListener(callback) {
    const index = this.listeners.conversationList.indexOf(callback);
    if (index !== -1) {
      this.listeners.conversationList.splice(index, 1);
    }
    if (this.listeners.conversationList.length === 0) {
      this.conversationParams = null;
      this.adjustPollingState();
    }
  }
  /**
   * 根据当前活跃的监听器调整轮询状态
   */
  adjustPollingState() {
    const hasAnyListeners = this.listeners.unreadCount.length > 0 || this.listeners.chatMessages.length > 0 || this.listeners.conversationList.length > 0;
    if (hasAnyListeners) {
      this.listeners.chatMessages.length > 0;
      if (this.isPolling) {
        this.stopPolling();
      }
      this.startPolling();
    } else {
      this.stopPolling();
    }
  }
  /**
   * 检查并停止轮询（如果没有活跃的监听器）
   */
  checkAndStopPolling() {
    const hasAnyListeners = this.listeners.unreadCount.length > 0 || this.listeners.chatMessages.length > 0 || this.listeners.conversationList.length > 0;
    if (!hasAnyListeners) {
      this.stopPolling();
    }
  }
  /**
   * 清除所有监听器
   */
  clearListeners() {
    this.listeners.unreadCount = [];
    this.listeners.chatMessages = [];
    this.listeners.conversationList = [];
    this.chatParams = null;
    this.conversationParams = null;
    this.stopPolling();
  }
}
const messagePollingService = new MessagePollingService();
exports.messagePollingService = messagePollingService;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/message-polling.js.map
