"use strict";
const utils_request = require("../utils/request.js");
function getConversationList(params = {}) {
  return utils_request.request({
    url: "/wx/client/conversation/list",
    method: "GET",
    data: {
      page: params.page || 1,
      size: params.size || 20,
      ...params
    }
  });
}
function clearAllUnread() {
  return utils_request.request({
    url: "/wx/client/conversation/clear-unread",
    method: "POST"
  });
}
function clearSessionUnread(sessionId) {
  return utils_request.request({
    url: "/wx/client/conversation/clear-session-unread",
    method: "POST",
    data: { sessionId }
  });
}
function sendMessage(params) {
  return utils_request.request({
    url: "/wx/client/message/send",
    method: "POST",
    data: params
  });
}
function getMessageList(params = {}) {
  return utils_request.request({
    url: "/wx/client/message/list",
    method: "GET",
    data: {
      targetId: params.targetId,
      page: params.page || 1,
      size: params.size || 20
    }
  });
}
function createConversation(params) {
  return utils_request.request({
    url: "/wx/client/conversation/create",
    method: "POST",
    data: params
  });
}
function markMessageRead(params) {
  return utils_request.request({
    url: "/wx/client/message/read",
    method: "POST",
    data: params
  });
}
function deleteConversation(conversationId) {
  return utils_request.request({
    url: "/wx/client/conversation/delete",
    method: "POST",
    data: { id: conversationId }
  });
}
function getUnreadCount() {
  return utils_request.request({
    url: "/wx/client/message/unread/count",
    method: "GET"
  });
}
const messageApi = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  clearAllUnread,
  clearSessionUnread,
  createConversation,
  deleteConversation,
  getConversationList,
  getMessageList,
  getUnreadCount,
  markMessageRead,
  sendMessage
}, Symbol.toStringTag, { value: "Module" }));
exports.createConversation = createConversation;
exports.getConversationList = getConversationList;
exports.getMessageList = getMessageList;
exports.getUnreadCount = getUnreadCount;
exports.markMessageRead = markMessageRead;
exports.messageApi = messageApi;
exports.sendMessage = sendMessage;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/message.js.map
