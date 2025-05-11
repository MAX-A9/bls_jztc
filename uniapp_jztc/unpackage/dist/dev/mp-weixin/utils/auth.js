"use strict";
const apis_index = require("../apis/index.js");
const utils_storage = require("./storage.js");
const utils_constants = require("./constants.js");
function isLoggedIn() {
  return !!utils_storage.getToken();
}
async function fetchAndSaveUserInfo() {
  try {
    const result = await apis_index.user.getClientInfo();
    if (result && result.code === utils_constants.API_CODE.SUCCESS && result.data) {
      const userData = result.data;
      if (userData.id && !userData.clientId) {
        userData.clientId = userData.id;
      }
      utils_storage.setUserInfo(userData);
      return userData;
    } else {
      const error = new Error(result.message || "获取客户信息失败");
      throw error;
    }
  } catch (error) {
    return Promise.reject(error);
  }
}
async function silentLogin() {
  try {
    const code = await apis_index.user.getWxLoginCode();
    const loginData = { code };
    const loginResult = await apis_index.user.wxappLogin(loginData);
    if (loginResult && loginResult.code === utils_constants.API_CODE.SUCCESS && loginResult.data && loginResult.data.token) {
      utils_storage.setToken(loginResult.data.token);
      const userInfo = await fetchAndSaveUserInfo();
      return {
        ...loginResult.data,
        userInfo
      };
    } else {
      const error = new Error(loginResult.message || "登录失败，返回数据不完整");
      throw error;
    }
  } catch (error) {
    return Promise.reject(error);
  }
}
async function fullLogin() {
  try {
    const code = await apis_index.user.getWxLoginCode();
    const wxUserInfo = await apis_index.user.getWxUserInfo();
    const loginData = {
      code,
      userInfo: wxUserInfo
    };
    const loginResult = await apis_index.user.wxappLogin(loginData);
    if (loginResult && loginResult.code === utils_constants.API_CODE.SUCCESS && loginResult.data && loginResult.data.token) {
      utils_storage.setToken(loginResult.data.token);
      const userInfo = await fetchAndSaveUserInfo();
      return {
        ...loginResult.data,
        userInfo
      };
    } else {
      const error = new Error(loginResult.message || "登录失败，返回数据不完整");
      throw error;
    }
  } catch (error) {
    return Promise.reject(error);
  }
}
function logout() {
  utils_storage.clearUserLoginState();
}
async function checkAndAutoLogin() {
  if (isLoggedIn()) {
    const cachedUserInfo = utils_storage.getUserInfo();
    const hasUserId = cachedUserInfo && (cachedUserInfo.clientId || cachedUserInfo.id);
    const hasBasicInfo = cachedUserInfo && cachedUserInfo.realName;
    if (hasUserId && hasBasicInfo) {
      return {
        token: utils_storage.getToken(),
        userInfo: cachedUserInfo
      };
    }
    try {
      const userInfo = await fetchAndSaveUserInfo();
      return {
        token: utils_storage.getToken(),
        userInfo
      };
    } catch (error) {
      utils_storage.clearUserLoginState();
    }
  }
  try {
    const loginResult = await silentLogin();
    if (!loginResult.userInfo && loginResult.token) {
      try {
        const userInfo = await fetchAndSaveUserInfo();
        return {
          ...loginResult,
          userInfo
        };
      } catch (infoError) {
        return loginResult;
      }
    }
    return loginResult;
  } catch (loginError) {
    return Promise.reject(loginError);
  }
}
exports.checkAndAutoLogin = checkAndAutoLogin;
exports.fetchAndSaveUserInfo = fetchAndSaveUserInfo;
exports.fullLogin = fullLogin;
exports.isLoggedIn = isLoggedIn;
exports.logout = logout;
exports.silentLogin = silentLogin;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/auth.js.map
