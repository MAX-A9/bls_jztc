"use strict";
const API_BASE_URL = "https://xcx.552i.com";
const API_TIMEOUT = 5e3;
const API_PATHS = {
  WXAPP_LOGIN: "/wx/wxapp-login",
  CLIENT_INFO: "/wx/client/info",
  CLIENT_UPDATE_PROFILE: "/wx/client/update-profile"
};
const STORAGE_KEYS = {
  USER_INFO: "USER_INFO",
  TOKEN: "token",
  // 确保和getStorageSync('token')使用的键名一致
  LOCATION: "currentLocation",
  IDLE_DRAFT: "idle_draft_list",
  // 闲置物品草稿列表
  INFO_DRAFT: "info_draft_list"
  // 信息发布草稿列表
};
const API_CODE = {
  SUCCESS: 0,
  ERROR: 1,
  UNAUTHORIZED: 401
};
exports.API_BASE_URL = API_BASE_URL;
exports.API_CODE = API_CODE;
exports.API_PATHS = API_PATHS;
exports.API_TIMEOUT = API_TIMEOUT;
exports.STORAGE_KEYS = STORAGE_KEYS;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/constants.js.map
