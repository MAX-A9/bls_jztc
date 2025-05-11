"use strict";
const common_vendor = require("../common/vendor.js");
const utils_storage = require("./storage.js");
const utils_constants = require("./constants.js");
let isShowingLoginModal = false;
function request(options) {
  return new Promise((resolve, reject) => {
    const token = utils_storage.getToken();
    const url = /^(http|https):\/\//.test(options.url) ? options.url : utils_constants.API_BASE_URL + options.url;
    const header = {
      "Content-Type": options.contentType || "application/json",
      ...options.header
    };
    if (token) {
      header["Authorization"] = `Bearer ${token}`;
    }
    let timeoutTimer = null;
    const timeoutPromise = new Promise((_, timeoutReject) => {
      timeoutTimer = setTimeout(() => {
        timeoutReject({ message: "请求超时，请检查网络", code: "TIMEOUT" });
        task && task.abort();
      }, options.timeout || utils_constants.API_TIMEOUT);
    });
    const task = common_vendor.index.request({
      url,
      data: options.data,
      method: options.method || "GET",
      header,
      success: (res) => {
        clearTimeout(timeoutTimer);
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data);
        } else if (res.statusCode === 401) {
          if (!isShowingLoginModal) {
            isShowingLoginModal = true;
            common_vendor.index.removeStorageSync("token");
            common_vendor.index.removeStorageSync("USER_INFO");
            common_vendor.index.showModal({
              title: "登录已过期",
              content: "您的登录已过期，请重新登录",
              showCancel: false,
              success: () => {
                common_vendor.index.$emit("showLoginModal");
                isShowingLoginModal = false;
              }
            });
          }
          reject({ code: 401, message: "未授权或token已过期" });
        } else {
          reject(res.data || { message: `请求失败，状态码：${res.statusCode}` });
        }
      },
      fail: (err) => {
        clearTimeout(timeoutTimer);
        reject(err || { message: "网络请求失败" });
      },
      complete: () => {
        if (options.complete) {
          options.complete();
        }
      }
    });
    return Promise.race([task, timeoutPromise]);
  });
}
function get(url, data = {}, noAuth = false) {
  return new Promise((resolve, reject) => {
    common_vendor.index.request({
      url: utils_constants.API_BASE_URL + url,
      method: "GET",
      data,
      header: {
        "Content-Type": "application/json",
        ...!noAuth && { "Authorization": `Bearer ${utils_storage.getToken()}` }
      },
      success: (res) => {
        if (res.statusCode === 200) {
          resolve(res.data);
        } else if (res.statusCode === 401) {
          if (!isShowingLoginModal) {
            isShowingLoginModal = true;
            common_vendor.index.removeStorageSync("token");
            common_vendor.index.removeStorageSync("USER_INFO");
            common_vendor.index.showModal({
              title: "登录已过期",
              content: "您的登录已过期，请重新登录",
              showCancel: false,
              success: () => {
                common_vendor.index.$emit("showLoginModal");
                isShowingLoginModal = false;
              }
            });
          }
          reject({ code: 401, message: "未授权或token已过期" });
        } else {
          reject(res);
        }
      },
      fail: (err) => {
        reject(err);
      }
    });
  });
}
function post(url, data = {}, options = {}) {
  return request({
    url,
    data,
    method: "POST",
    ...options
  });
}
function put(url, data = {}, options = {}) {
  return request({
    url,
    data,
    method: "PUT",
    ...options
  });
}
function del(url, data = {}, options = {}) {
  return request({
    url,
    data,
    method: "DELETE",
    ...options
  });
}
exports.del = del;
exports.get = get;
exports.post = post;
exports.put = put;
exports.request = request;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/request.js.map
