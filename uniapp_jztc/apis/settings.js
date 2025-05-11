/**
 * 小程序设置相关API
 */
import { get } from '../utils/request.js';

/**
 * 获取小程序基础设置
 * @returns {Promise} 返回Promise对象
 */
export function getBaseSettings() {
  return get('/wx/mini-program/base/settings');
} 