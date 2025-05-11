/**
 * 本地存储相关方法
 */
import { STORAGE_KEYS } from './constants.js';

/**
 * 存储数据
 * @param {String} key 存储键名
 * @param {any} data 存储的数据
 */
export function setStorage(key, data) {
  try {
    uni.setStorageSync(key, data);
  } catch (e) {
    console.error(`数据存储失败 [${key}]`);
  }
}

/**
 * 获取存储的数据
 * @param {String} key 存储键名
 * @returns {any} 存储的数据
 */
export function getStorage(key) {
  try {
    const data = uni.getStorageSync(key);
    return data;
  } catch (e) {
    console.error(`数据读取失败 [${key}]`);
    return null;
  }
}

/**
 * 移除存储的数据
 * @param {String} key 存储键名
 */
export function removeStorage(key) {
  try {
    uni.removeStorageSync(key);
  } catch (e) {
    console.error(`数据删除失败 [${key}]`);
  }
}

/**
 * 清除所有存储数据
 */
export function clearStorage() {
  try {
    uni.clearStorageSync();
  } catch (e) {
    console.error('清空数据失败');
  }
}

/**
 * 存储用户信息
 * @param {Object} userInfo 用户信息
 */
export function setUserInfo(userInfo) {
  setStorage(STORAGE_KEYS.USER_INFO, userInfo);
}

/**
 * 获取用户信息
 * @returns {Object} 用户信息
 */
export function getUserInfo() {
  const info = getStorage(STORAGE_KEYS.USER_INFO) || {};
  
  // 确保兼容性：如果只有id字段，复制到clientId
  if (info && info.id && !info.clientId) {
    info.clientId = info.id;
  }
  
  return info;
}

/**
 * 存储token
 * @param {String} token 登录token
 */
export function setToken(token) {
  setStorage(STORAGE_KEYS.TOKEN, token);
}

/**
 * 获取token
 * @returns {String} 登录token
 */
export function getToken() {
  const token = getStorage(STORAGE_KEYS.TOKEN) || '';
  return token;
}

/**
 * 清除用户登录态
 */
export function clearUserLoginState() {
  removeStorage(STORAGE_KEYS.USER_INFO);
  removeStorage(STORAGE_KEYS.TOKEN);
}

/**
 * 保存闲置物品草稿
 * @param {Object} draftData 草稿数据
 * @returns {Number} 草稿ID
 */
export function saveIdleDraft(draftData) {
  try {
    // 获取现有的草稿列表
    const draftList = getIdleDraftList() || [];
    
    // 生成唯一ID
    const draftId = Date.now();
    
    // 添加更新时间和ID
    const newDraft = {
      ...draftData,
      id: draftId,
      updateTime: new Date().toISOString(),
      type: 'idle' // 标记为闲置物品类型
    };
    
    // 添加到列表开头（最新的在前面）
    draftList.unshift(newDraft);
    
    // 保存回存储
    setStorage(STORAGE_KEYS.IDLE_DRAFT, draftList);
    
    return draftId;
  } catch (e) {
    console.error('保存闲置物品草稿失败', e);
    return 0;
  }
}

/**
 * 获取闲置物品草稿列表
 * @returns {Array} 草稿列表
 */
export function getIdleDraftList() {
  return getStorage(STORAGE_KEYS.IDLE_DRAFT) || [];
}

/**
 * 获取单个闲置物品草稿
 * @param {Number} draftId 草稿ID
 * @returns {Object|null} 草稿数据
 */
export function getIdleDraft(draftId) {
  const draftList = getIdleDraftList();
  return draftList.find(item => item.id === draftId) || null;
}

/**
 * 删除闲置物品草稿
 * @param {Number} draftId 草稿ID
 * @returns {Boolean} 是否删除成功
 */
export function deleteIdleDraft(draftId) {
  try {
    let draftList = getIdleDraftList();
    draftList = draftList.filter(item => item.id !== draftId);
    setStorage(STORAGE_KEYS.IDLE_DRAFT, draftList);
    return true;
  } catch (e) {
    console.error('删除闲置物品草稿失败', e);
    return false;
  }
}

/**
 * 保存信息发布草稿
 * @param {Object} draftData 草稿数据
 * @returns {Number} 草稿ID
 */
export function saveInfoDraft(draftData) {
  try {
    // 获取现有的草稿列表
    const draftList = getInfoDraftList() || [];
    
    // 生成唯一ID
    const draftId = Date.now();
    
    // 添加更新时间和ID
    const newDraft = {
      ...draftData,
      id: draftId,
      updateTime: new Date().toISOString(),
      type: 'info' // 标记为信息发布类型
    };
    
    // 添加到列表开头（最新的在前面）
    draftList.unshift(newDraft);
    
    // 保存回存储
    setStorage(STORAGE_KEYS.INFO_DRAFT, draftList);
    
    return draftId;
  } catch (e) {
    console.error('保存信息发布草稿失败', e);
    return 0;
  }
}

/**
 * 获取信息发布草稿列表
 * @returns {Array} 草稿列表
 */
export function getInfoDraftList() {
  return getStorage(STORAGE_KEYS.INFO_DRAFT) || [];
}

/**
 * 获取单个信息发布草稿
 * @param {Number} draftId 草稿ID
 * @returns {Object|null} 草稿数据
 */
export function getInfoDraft(draftId) {
  const draftList = getInfoDraftList();
  return draftList.find(item => item.id === draftId) || null;
}

/**
 * 删除信息发布草稿
 * @param {Number} draftId 草稿ID
 * @returns {Boolean} 是否删除成功
 */
export function deleteInfoDraft(draftId) {
  try {
    let draftList = getInfoDraftList();
    draftList = draftList.filter(item => item.id !== draftId);
    setStorage(STORAGE_KEYS.INFO_DRAFT, draftList);
    return true;
  } catch (e) {
    console.error('删除信息发布草稿失败', e);
    return false;
  }
}

/**
 * 获取所有草稿列表（包括闲置和信息）
 * @returns {Array} 草稿列表
 */
export function getAllDraftList() {
  const idleDrafts = getIdleDraftList();
  const infoDrafts = getInfoDraftList();
  
  // 合并两个列表并按时间排序
  return [...idleDrafts, ...infoDrafts].sort((a, b) => {
    return new Date(b.updateTime) - new Date(a.updateTime);
  });
}

/**
 * 获取闲置物品草稿列表（别名，为兼容性）
 * @returns {Array} 草稿列表
 */
export function getIdleDrafts() {
  return getIdleDraftList();
}

/**
 * 获取信息发布草稿列表（别名，为兼容性）
 * @returns {Array} 草稿列表
 */
export function getInfoDrafts() {
  return getInfoDraftList();
} 