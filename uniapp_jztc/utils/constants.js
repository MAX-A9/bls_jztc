/**
 * 常量配置
 */

// API基础URL
export const API_BASE_URL = 'http://localhost:8000'; // 替换为实际的API地址

// 图片上传URL (使用基础URL加上路径)
export const UPLOAD_URL = API_BASE_URL + '/wx/upload/image'; 

// API超时设置（毫秒）
export const API_TIMEOUT = 5000;

// API路径
export const API_PATHS = {
  WXAPP_LOGIN: '/wx/wxapp-login',
  CLIENT_INFO: '/wx/client/info',
  CLIENT_UPDATE_PROFILE: '/wx/client/update-profile'
};

// 存储键名
export const STORAGE_KEYS = {
  USER_INFO: 'USER_INFO',
  TOKEN: 'token', // 确保和getStorageSync('token')使用的键名一致
  LOCATION: 'currentLocation',
  IDLE_DRAFT: 'idle_draft_list', // 闲置物品草稿列表
  INFO_DRAFT: 'info_draft_list' // 信息发布草稿列表
};

// 接口状态码
export const API_CODE = {
  SUCCESS: 0,
  ERROR: 1,
  UNAUTHORIZED: 401
}; 