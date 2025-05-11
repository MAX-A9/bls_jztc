/**
 * API 入口文件，统一导出所有API
 */

// 导入模块
import * as userApi from './user.js';
import * as contentApi from './content.js';
import * as payApi from './pay.js';
import * as agreementApi from './agreement.js';
import * as settingsApi from './settings.js';
import * as messageApi from './message.js';
import * as categoryApi from './category.js';
import * as publishApi from './publish.js';
import * as vipApi from './vip.js';
import * as shareApi from './share.js';

// 导出模块
export const user = userApi;
export const content = contentApi;
export const pay = payApi;
export const agreement = agreementApi;
export const settings = settingsApi;
export const message = messageApi;
export const category = categoryApi;
export const publish = publishApi;
export const vip = vipApi;
export const share = shareApi;

// 默认导出所有API
export default {
  user,
  content,
  pay,
  agreement,
  settings,
  message,
  category,
  publish,
  vip,
  share
}; 