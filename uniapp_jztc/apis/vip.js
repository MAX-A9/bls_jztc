/**
 * VIP相关API
 */
import { get, post } from '@/utils/request.js';

/**
 * 获取兑换记录列表
 * @param {Object} params 查询参数 {page, size}
 * @returns {Promise} 返回Promise对象
 */
export function getExchangeRecords(params = {}) {
  return get('/wx/client/exchange-record/page', params);
}

/**
 * 获取最近兑换记录列表（用于通知展示）
 * @returns {Promise} 返回Promise对象
 */
export function getRecentExchangeRecords() {
  return get('/wx/exchange-record/list');
}

/**
 * 创建兑换记录
 * @param {Object} data 兑换信息
 * @returns {Promise} 返回Promise对象
 */
export function createExchangeRecord(data) {
  return post('/wx/client/exchange-record/create', data);
} 