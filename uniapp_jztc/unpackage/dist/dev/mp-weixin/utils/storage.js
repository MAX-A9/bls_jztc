"use strict";
const common_vendor = require("../common/vendor.js");
const utils_constants = require("./constants.js");
function setStorage(key, data) {
  try {
    common_vendor.index.setStorageSync(key, data);
  } catch (e) {
    common_vendor.index.__f__("error", "at utils/storage.js:15", `数据存储失败 [${key}]`);
  }
}
function getStorage(key) {
  try {
    const data = common_vendor.index.getStorageSync(key);
    return data;
  } catch (e) {
    common_vendor.index.__f__("error", "at utils/storage.js:29", `数据读取失败 [${key}]`);
    return null;
  }
}
function removeStorage(key) {
  try {
    common_vendor.index.removeStorageSync(key);
  } catch (e) {
    common_vendor.index.__f__("error", "at utils/storage.js:42", `数据删除失败 [${key}]`);
  }
}
function setUserInfo(userInfo) {
  setStorage(utils_constants.STORAGE_KEYS.USER_INFO, userInfo);
}
function getUserInfo() {
  const info = getStorage(utils_constants.STORAGE_KEYS.USER_INFO) || {};
  if (info && info.id && !info.clientId) {
    info.clientId = info.id;
  }
  return info;
}
function setToken(token) {
  setStorage(utils_constants.STORAGE_KEYS.TOKEN, token);
}
function getToken() {
  const token = getStorage(utils_constants.STORAGE_KEYS.TOKEN) || "";
  return token;
}
function clearUserLoginState() {
  removeStorage(utils_constants.STORAGE_KEYS.USER_INFO);
  removeStorage(utils_constants.STORAGE_KEYS.TOKEN);
}
function saveIdleDraft(draftData) {
  try {
    const draftList = getIdleDraftList() || [];
    const draftId = Date.now();
    const newDraft = {
      ...draftData,
      id: draftId,
      updateTime: (/* @__PURE__ */ new Date()).toISOString(),
      type: "idle"
      // 标记为闲置物品类型
    };
    draftList.unshift(newDraft);
    setStorage(utils_constants.STORAGE_KEYS.IDLE_DRAFT, draftList);
    return draftId;
  } catch (e) {
    common_vendor.index.__f__("error", "at utils/storage.js:134", "保存闲置物品草稿失败", e);
    return 0;
  }
}
function getIdleDraftList() {
  return getStorage(utils_constants.STORAGE_KEYS.IDLE_DRAFT) || [];
}
function getIdleDraft(draftId) {
  const draftList = getIdleDraftList();
  return draftList.find((item) => item.id === draftId) || null;
}
function deleteIdleDraft(draftId) {
  try {
    let draftList = getIdleDraftList();
    draftList = draftList.filter((item) => item.id !== draftId);
    setStorage(utils_constants.STORAGE_KEYS.IDLE_DRAFT, draftList);
    return true;
  } catch (e) {
    common_vendor.index.__f__("error", "at utils/storage.js:169", "删除闲置物品草稿失败", e);
    return false;
  }
}
function saveInfoDraft(draftData) {
  try {
    const draftList = getInfoDraftList() || [];
    const draftId = Date.now();
    const newDraft = {
      ...draftData,
      id: draftId,
      updateTime: (/* @__PURE__ */ new Date()).toISOString(),
      type: "info"
      // 标记为信息发布类型
    };
    draftList.unshift(newDraft);
    setStorage(utils_constants.STORAGE_KEYS.INFO_DRAFT, draftList);
    return draftId;
  } catch (e) {
    common_vendor.index.__f__("error", "at utils/storage.js:203", "保存信息发布草稿失败", e);
    return 0;
  }
}
function getInfoDraftList() {
  return getStorage(utils_constants.STORAGE_KEYS.INFO_DRAFT) || [];
}
function getInfoDraft(draftId) {
  const draftList = getInfoDraftList();
  return draftList.find((item) => item.id === draftId) || null;
}
function deleteInfoDraft(draftId) {
  try {
    let draftList = getInfoDraftList();
    draftList = draftList.filter((item) => item.id !== draftId);
    setStorage(utils_constants.STORAGE_KEYS.INFO_DRAFT, draftList);
    return true;
  } catch (e) {
    common_vendor.index.__f__("error", "at utils/storage.js:238", "删除信息发布草稿失败", e);
    return false;
  }
}
function getAllDraftList() {
  const idleDrafts = getIdleDraftList();
  const infoDrafts = getInfoDraftList();
  return [...idleDrafts, ...infoDrafts].sort((a, b) => {
    return new Date(b.updateTime) - new Date(a.updateTime);
  });
}
exports.clearUserLoginState = clearUserLoginState;
exports.deleteIdleDraft = deleteIdleDraft;
exports.deleteInfoDraft = deleteInfoDraft;
exports.getAllDraftList = getAllDraftList;
exports.getIdleDraft = getIdleDraft;
exports.getInfoDraft = getInfoDraft;
exports.getToken = getToken;
exports.getUserInfo = getUserInfo;
exports.saveIdleDraft = saveIdleDraft;
exports.saveInfoDraft = saveInfoDraft;
exports.setToken = setToken;
exports.setUserInfo = setUserInfo;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/storage.js.map
