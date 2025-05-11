/**
 * 协议相关API
 */
import { get } from '../utils/request.js';

/**
 * 获取协议内容
 * @param {Object} params 参数对象
 * @param {String} params.type 协议类型 user-用户协议 privacy-隐私政策
 * @returns {Promise} 返回Promise对象
 */
export function getAgreement(params) {
  return get('/wx/agreement/get', params);
} 