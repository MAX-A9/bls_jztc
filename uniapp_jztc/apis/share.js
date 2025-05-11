/**
 * 分享相关API
 */
import { get } from '../utils/request.js';

/**
 * 获取分享设置
 * @returns {Promise} 返回包含分享设置的Promise对象
 */
export function getShareSettings() {
  return get('/wx/share/settings');
} 