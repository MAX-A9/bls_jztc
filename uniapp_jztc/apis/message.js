/**
 * 消息相关API
 */
import { request } from '@/utils/request.js';

/**
 * 获取会话列表
 * @param {Object} params - 请求参数
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.size - 每页条数，默认20
 * @returns {Promise} 返回Promise对象
 */
export function getConversationList(params = {}) {
  return request({
    url: '/wx/client/conversation/list',
    method: 'GET',
    data: {
      page: params.page || 1,
      size: params.size || 20,
      ...params
    }
  });
}

/**
 * 清除所有未读消息
 * @returns {Promise} 返回Promise对象
 */
export function clearAllUnread() {
  return request({
    url: '/wx/client/conversation/clear-unread',
    method: 'POST'
  });
}

/**
 * 清除指定会话的未读消息
 * @param {String} sessionId - 会话ID
 * @returns {Promise} 返回Promise对象
 */
export function clearSessionUnread(sessionId) {
  return request({
    url: '/wx/client/conversation/clear-session-unread',
    method: 'POST',
    data: { sessionId }
  });
}

/**
 * 发送私信
 * @param {Object} params - 请求参数
 * @param {String} params.content - 消息内容
 * @param {Number} params.receiverId - 接收者ID
 * @returns {Promise} 返回Promise对象
 */
export function sendMessage(params) {
  return request({
    url: '/wx/client/message/send',
    method: 'POST',
    data: params
  });
}

/**
 * 获取消息列表
 * @param {Object} params - 请求参数
 * @param {Number} params.targetId - 目标用户ID
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.size - 每页条数，默认20
 * @returns {Promise} 返回Promise对象
 */
export function getMessageList(params = {}) {
  return request({
    url: '/wx/client/message/list',
    method: 'GET',
    data: {
      targetId: params.targetId,
      page: params.page || 1,
      size: params.size || 20
    }
  });
}

/**
 * 创建会话
 * @param {Object} params - 请求参数
 * @param {Number} params.targetId - 目标用户ID
 * @returns {Promise} 返回Promise对象
 */
export function createConversation(params) {
  return request({
    url: '/wx/client/conversation/create',
    method: 'POST',
    data: params
  });
}

/**
 * 标记消息为已读
 * @param {Object} params - 请求参数
 * @param {Number} params.targetId - 目标用户ID
 * @returns {Promise} 返回Promise对象
 */
export function markMessageRead(params) {
  return request({
    url: '/wx/client/message/read',
    method: 'POST',
    data: params
  });
}

/**
 * 删除会话
 * @param {String} conversationId - 会话ID
 * @returns {Promise} 返回Promise对象
 */
export function deleteConversation(conversationId) {
  return request({
    url: '/wx/client/conversation/delete',
    method: 'POST',
    data: { id: conversationId }
  });
}

/**
 * 获取未读消息数量
 * @returns {Promise} 返回Promise对象，包含未读消息数量
 */
export function getUnreadCount() {
  return request({
    url: '/wx/client/message/unread/count',
    method: 'GET'
  });
} 