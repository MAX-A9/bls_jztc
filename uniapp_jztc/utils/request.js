/**
 * 网络请求封装
 */
import { getToken } from './storage.js';
import { API_BASE_URL, API_TIMEOUT } from './constants.js';

// 添加事件总线，用于触发登录窗口
let isShowingLoginModal = false;

/**
 * 发起网络请求
 * @param {Object} options 请求配置
 * @returns {Promise} 返回Promise对象
 */
export function request(options) {
  return new Promise((resolve, reject) => {
    // 获取token
    const token = getToken();
    
    // 处理请求URL
    const url = /^(http|https):\/\//.test(options.url) 
      ? options.url 
      : API_BASE_URL + options.url;
    
    // 请求头
    const header = {
      'Content-Type': options.contentType || 'application/json',
      ...options.header
    };
    
    // 添加token到请求头
    if (token) {
      header['Authorization'] = `Bearer ${token}`;
    }
    
    // 请求超时定时器
    let timeoutTimer = null;
    
    // 创建超时Promise
    const timeoutPromise = new Promise((_, timeoutReject) => {
      timeoutTimer = setTimeout(() => {
        timeoutReject({ message: '请求超时，请检查网络', code: 'TIMEOUT' });
        task && task.abort(); // 超时后中断请求
      }, options.timeout || API_TIMEOUT);
    });
    
    // 发起请求
    const task = uni.request({
      url,
      data: options.data,
      method: options.method || 'GET',
      header,
      success: (res) => {
        clearTimeout(timeoutTimer);
        
        // 请求成功
        if (res.statusCode >= 200 && res.statusCode < 300) {
          // 直接返回接口数据，不再额外处理
          resolve(res.data);
        } 
        // 未授权
        else if (res.statusCode === 401) {
          // 弹出登录窗口，避免多次弹出
          if (!isShowingLoginModal) {
            isShowingLoginModal = true;

            // 清除本地登录信息
            uni.removeStorageSync('token');
            uni.removeStorageSync('USER_INFO');
            
            // 显示提示
            uni.showModal({
              title: '登录已过期',
              content: '您的登录已过期，请重新登录',
              showCancel: false,
              success: () => {
                // 触发全局登录事件
                uni.$emit('showLoginModal');
                isShowingLoginModal = false;
              }
            });
          }
          
          // 返回错误信息
          reject({ code: 401, message: '未授权或token已过期' });
        } 
        // 其他错误
        else {
          reject(res.data || { message: `请求失败，状态码：${res.statusCode}` });
        }
      },
      fail: (err) => {
        clearTimeout(timeoutTimer);
        reject(err || { message: '网络请求失败' });
      },
      complete: () => {
        // 请求完成的回调
        if (options.complete) {
          options.complete();
        }
      }
    });
    
    // 使用Promise.race竞争，哪个先完成就返回哪个结果
    return Promise.race([task, timeoutPromise]);
  });
}

/**
 * GET请求
 * @param {String} url 请求地址
 * @param {Object} data 请求参数
 * @param {Boolean} noAuth 是否不需要授权
 * @returns {Promise} 返回Promise对象
 */
export function get(url, data = {}, noAuth = false) {
  return new Promise((resolve, reject) => {
    uni.request({
      url: API_BASE_URL + url,
      method: 'GET',
      data: data,
      header: {
        'Content-Type': 'application/json',
        ...(!noAuth && { 'Authorization': `Bearer ${getToken()}` })
      },
      success: (res) => {
        if (res.statusCode === 200) {
          resolve(res.data)
        } else if (res.statusCode === 401) {
          // 处理401未授权错误
          if (!isShowingLoginModal) {
            isShowingLoginModal = true;
            
            // 清除本地登录信息
            uni.removeStorageSync('token');
            uni.removeStorageSync('USER_INFO');
            
            // 显示提示
            uni.showModal({
              title: '登录已过期',
              content: '您的登录已过期，请重新登录',
              showCancel: false,
              success: () => {
                // 触发全局登录事件
                uni.$emit('showLoginModal');
                isShowingLoginModal = false;
              }
            });
          }
          
          reject({ code: 401, message: '未授权或token已过期' });
        } else {
          reject(res)
        }
      },
      fail: (err) => {
        reject(err)
      }
    })
  })
}

/**
 * POST请求
 * @param {String} url 请求地址
 * @param {Object} data 请求参数
 * @param {Object} options 其他配置
 * @returns {Promise} 返回Promise对象
 */
export function post(url, data = {}, options = {}) {
  return request({
    url,
    data,
    method: 'POST',
    ...options
  });
}

/**
 * PUT请求
 * @param {String} url 请求地址
 * @param {Object} data 请求参数
 * @param {Object} options 其他配置
 * @returns {Promise} 返回Promise对象
 */
export function put(url, data = {}, options = {}) {
  return request({
    url,
    data,
    method: 'PUT',
    ...options
  });
}

/**
 * DELETE请求
 * @param {String} url 请求地址
 * @param {Object} data 请求参数
 * @param {Object} options 其他配置
 * @returns {Promise} 返回Promise对象
 */
export function del(url, data = {}, options = {}) {
  return request({
    url,
    data,
    method: 'DELETE',
    ...options
  });
}

export default {
  request,
  get,
  post,
  put,
  del
}; 