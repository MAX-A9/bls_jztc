/**
 * 消息轮询服务
 * 用于统一管理应用内的消息轮询，避免多个组件重复创建轮询
 */
import { getUnreadCount, getMessageList, getConversationList } from '@/apis/message.js';

class MessagePollingService {
  constructor() {
    this.pollingTimer = null;
    this.pollingInterval = 5000; // 默认5秒
    this.chatPollingInterval = 3000; // 聊天页面3秒
    this.isPolling = false;
    this.listeners = {
      unreadCount: [], // 未读消息数监听器
      chatMessages: [], // 聊天消息监听器
      conversationList: [] // 会话列表监听器
    };
    this.chatParams = null; // 聊天轮询参数
    this.conversationParams = null; // 会话列表轮询参数
    this.globalPaused = false; // 全局暂停状态
    this.previousPollingState = false; // 记录暂停前的轮询状态
    this.whitelistPages = ['/pages/chat/detail']; // 白名单页面不受全局暂停影响
    this.currentPage = ''; // 当前页面路径
  }

  /**
   * 设置当前页面路径
   * @param {String} pagePath - 页面路径
   */
  setCurrentPage(pagePath) {
    this.currentPage = pagePath;
    console.log('当前页面路径:', pagePath);
    
    // 如果当前页面在白名单中，且全局已暂停，则恢复轮询
    if (this.isPageInWhitelist(pagePath) && this.globalPaused) {
      console.log('当前页面在白名单中，恢复轮询');
      this.resumePollingForWhitelist();
    }
  }
  
  /**
   * 检查页面是否在白名单中
   * @param {String} pagePath - 页面路径
   * @returns {Boolean} 是否在白名单中
   */
  isPageInWhitelist(pagePath) {
    return this.whitelistPages.some(whitePath => pagePath && pagePath.indexOf(whitePath) !== -1);
  }
  
  /**
   * 为白名单页面恢复轮询，但不改变全局暂停状态
   */
  resumePollingForWhitelist() {
    if (this.globalPaused && !this.isPolling && this.previousPollingState) {
      console.log('白名单页面特殊处理：恢复轮询但保持全局暂停状态');
      this.isPolling = true;
      this.poll();
      
      // 根据当前活跃的轮询类型选择合适的间隔
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
    // 如果正在轮询，需要重新启动以应用新间隔
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
    
    // 如果正在轮询且有聊天参数，需要重新启动以应用新间隔
    if (this.isPolling && this.chatParams) {
      this.stopPolling();
      this.startPolling();
    }
  }

  /**
   * 启动轮询
   */
  startPolling() {
    if (this.isPolling || this.globalPaused) return;

    this.isPolling = true;
    this.poll();

    // 根据当前活跃的轮询类型选择合适的间隔
    const interval = this.chatParams ? this.chatPollingInterval : this.pollingInterval;
    
    this.pollingTimer = setInterval(() => {
      this.poll();
    }, interval);

    console.log(`消息轮询已启动，间隔 ${interval} ms`);
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
    console.log('消息轮询已停止');
  }

  /**
   * 全局暂停轮询（适用于进入非底部导航页面时）
   * 此方法会记住当前的轮询状态，以便恢复时能够还原
   */
  pausePollingGlobally() {
    if (this.globalPaused) return;
    
    // 如果当前页面在白名单中，不执行暂停
    if (this.isPageInWhitelist(this.currentPage)) {
      console.log('当前页面在白名单中，不执行全局暂停');
      return;
    }
    
    this.previousPollingState = this.isPolling;
    this.globalPaused = true;
    
    if (this.isPolling) {
      this.stopPolling();
      console.log('消息轮询已全局暂停');
    }
  }
  
  /**
   * 全局恢复轮询（适用于返回到底部导航页面时）
   * 如果之前处于轮询状态，则会恢复轮询
   */
  resumePollingGlobally() {
    if (!this.globalPaused) return;
    
    this.globalPaused = false;
    
    if (this.previousPollingState) {
      this.startPolling();
      console.log('消息轮询已全局恢复');
    }
  }

  /**
   * 执行轮询
   */
  async poll() {
    // 如果全局暂停了，且当前页面不在白名单中，直接返回
    if (this.globalPaused && !this.isPageInWhitelist(this.currentPage)) return;
    
    // 如果有注册的未读消息监听器，获取未读消息数
    if (this.listeners.unreadCount.length > 0) {
      this.pollUnreadCount();
    }

    // 如果有聊天参数和聊天消息监听器，获取聊天消息
    if (this.chatParams && this.listeners.chatMessages.length > 0) {
      this.pollChatMessages();
    }
    
    // 如果有会话列表监听器，获取会话列表
    if (this.conversationParams && this.listeners.conversationList.length > 0) {
      this.pollConversationList();
    }
  }

  /**
   * 获取未读消息数
   */
  async pollUnreadCount() {
    try {
      const res = await getUnreadCount();
      if (res && res.code === 0) {
        const unreadCount = res.data.unreadCount || 0;
        // 通知所有监听器
        this.listeners.unreadCount.forEach(listener => {
          listener(unreadCount);
        });
      }
    } catch (error) {
      console.error('获取未读消息数失败:', error);
    }
  }

  /**
   * 获取聊天消息
   */
  async pollChatMessages() {
    if (!this.chatParams || !this.chatParams.targetId) return;

    try {
      const res = await getMessageList({
        targetId: this.chatParams.targetId,
        page: 1,
        size: 20,
        lastId: this.chatParams.lastId || 0
      });

      if (res && res.code === 0) {
        // 通知所有聊天消息监听器
        this.listeners.chatMessages.forEach(listener => {
          listener(res.data);
        });
      }
    } catch (error) {
      console.error('获取聊天消息失败:', error);
    }
  }
  
  /**
   * 获取会话列表
   */
  async pollConversationList() {
    if (!this.conversationParams) return;
    
    try {
      // 确保每次获取的是最新数据，不受缓存影响
      const timestamp = new Date().getTime();
      const res = await getConversationList({
        page: this.conversationParams.page || 1,
        size: this.conversationParams.size || 20,
        _t: timestamp // 添加时间戳防止缓存
      });
      
      if (res && res.code === 0) {
        // 通知所有会话列表监听器
        this.listeners.conversationList.forEach(listener => {
          listener(res.data);
        });
      }
    } catch (error) {
      console.error('获取会话列表失败:', error);
    }
  }

  /**
   * 添加未读消息数监听器
   * @param {Function} callback - 监听器回调函数
   * @returns {Function} 移除监听器的函数
   */
  addUnreadCountListener(callback) {
    if (typeof callback !== 'function') return () => {};
    
    this.listeners.unreadCount.push(callback);
    
    // 启动轮询（如果尚未启动）
    if (!this.isPolling) {
      this.startPolling();
    } else {
      // 已经在轮询中，立即获取一次数据
      this.pollUnreadCount();
    }
    
    // 返回移除监听器的函数
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
    
    // 如果所有监听器都已移除，停止轮询
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
    
    if (typeof callback === 'function') {
      this.listeners.chatMessages.push(callback);
      
      // 调整为聊天轮询的间隔
      if (this.isPolling) {
        this.stopPolling();
      }
      this.startPolling();
      
      // 返回移除监听器的函数
      return () => {
        this.removeChatMessageListener(callback);
      };
    }
    
    return () => {};
  }

  /**
   * 更新聊天参数
   * @param {Object} params - 聊天参数
   */
  updateChatParams(params) {
    if (this.chatParams) {
      this.chatParams = {...this.chatParams, ...params};
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
    
    // 如果所有聊天监听器被移除，清除聊天参数
    if (this.listeners.chatMessages.length === 0) {
      this.chatParams = null;
      
      // 根据当前活跃的监听器状态调整轮询
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
    
    if (typeof callback === 'function') {
      this.listeners.conversationList.push(callback);
      
      // 启动轮询（如果尚未启动）
      if (!this.isPolling) {
        this.startPolling();
      } else {
        // 已经在轮询中，立即获取一次数据
        this.pollConversationList();
      }
      
      // 返回移除监听器的函数
      return () => {
        this.removeConversationListListener(callback);
      };
    }
    
    return () => {};
  }
  
  /**
   * 更新会话列表参数
   * @param {Object} params - 会话列表参数
   */
  updateConversationParams(params) {
    if (this.conversationParams) {
      this.conversationParams = {...this.conversationParams, ...params};
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
    
    // 如果所有会话列表监听器被移除，清除会话列表参数
    if (this.listeners.conversationList.length === 0) {
      this.conversationParams = null;
      
      // 根据当前活跃的监听器状态调整轮询
      this.adjustPollingState();
    }
  }
  
  /**
   * 根据当前活跃的监听器调整轮询状态
   */
  adjustPollingState() {
    // 如果还有任何类型的监听器，保持轮询但可能调整间隔
    const hasAnyListeners = 
      this.listeners.unreadCount.length > 0 || 
      this.listeners.chatMessages.length > 0 ||
      this.listeners.conversationList.length > 0;
    
    if (hasAnyListeners) {
      // 如果有聊天监听器，使用聊天轮询间隔，否则使用默认间隔
      const shouldUseChatInterval = this.listeners.chatMessages.length > 0;
      
      // 如果轮询间隔需要变更，重新启动轮询
      if (this.isPolling) {
        this.stopPolling();
      }
      this.startPolling();
    } else {
      // 如果没有任何监听器，停止轮询
      this.stopPolling();
    }
  }
  
  /**
   * 检查并停止轮询（如果没有活跃的监听器）
   */
  checkAndStopPolling() {
    const hasAnyListeners = 
      this.listeners.unreadCount.length > 0 || 
      this.listeners.chatMessages.length > 0 ||
      this.listeners.conversationList.length > 0;
    
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

// 导出单例
export default new MessagePollingService(); 