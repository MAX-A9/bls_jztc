/**
 * 用户相关API
 */
import { post, get, put } from '../utils/request.js';
import { API_PATHS } from '../utils/constants.js';

/**
 * 微信小程序登录
 * @param {Object} data 请求数据，包含code和userInfo
 * @returns {Promise} 返回Promise对象
 */
export function wxappLogin(data) {
  return post(API_PATHS.WXAPP_LOGIN, data)
    .catch(err => {
      throw err;
    });
}

/**
 * 获取客户信息
 * @returns {Promise} 返回Promise对象，包含客户信息
 */
export function getClientInfo() {
  return get(API_PATHS.CLIENT_INFO)
    .catch(err => {
      throw err;
    });
}

/**
 * 更新客户个人信息
 * @param {Object} data 请求数据，可包含phone、realName等字段
 * @returns {Promise} 返回Promise对象
 */
export function updateClientProfile(data) {
  return put(API_PATHS.CLIENT_UPDATE_PROFILE, data)
    .catch(err => {
      throw err;
    });
}

/**
 * 获取微信登录凭证
 * @returns {Promise} 返回Promise对象，包含登录凭证code
 */
export function getWxLoginCode() {
  return new Promise((resolve, reject) => {
    uni.login({
      provider: 'weixin',
      success: (res) => {
        if (res.code) {
          resolve(res.code);
        } else {
          const error = { message: '获取微信登录凭证失败', res };
          reject(error);
        }
      },
      fail: (err) => {
        reject(err || { message: '微信登录失败' });
      }
    });
  });
}

/**
 * 获取微信用户信息
 * @returns {Promise} 返回Promise对象，包含用户信息
 */
export function getWxUserInfo() {
  return new Promise((resolve, reject) => {
    uni.getUserProfile({
      desc: '用于完善用户资料',
      success: (res) => {
        if (res.userInfo) {
          resolve(res.userInfo);
        } else {
          reject({ message: '获取用户信息失败' });
        }
      },
      fail: (err) => {
        reject(err || { message: '用户拒绝授权' });
      }
    });
  });
}

/**
 * 获取图片base64编码
 * @param {String} imageUrl 图片URL
 * @returns {Promise<String>} 返回Promise对象，包含图片base64编码
 */
export function getImageBase64(imageUrl) {
  return new Promise((resolve, reject) => {
    uni.getImageInfo({
      src: imageUrl,
      success: (imgInfo) => {
        // 微信小程序环境下使用传统方式处理
        const canvasSize = 200; // 设置合适的尺寸
        const ctx = uni.createCanvasContext('avatarCanvas');
        
        // 清空画布
        ctx.clearRect(0, 0, canvasSize, canvasSize);
        
        ctx.drawImage(imgInfo.path, 0, 0, canvasSize, canvasSize);
        
        // 设置2秒超时，防止卡住
        let drawTimeout = setTimeout(() => {
          reject(new Error('Canvas绘制超时'));
        }, 2000);
        
        ctx.draw(false, () => {
          clearTimeout(drawTimeout);
          
          uni.canvasToTempFilePath({
            canvasId: 'avatarCanvas',
            fileType: 'jpg',
            quality: 0.8,
            success: (res) => {
              // 读取临时文件为base64
              uni.getFileSystemManager().readFile({
                filePath: res.tempFilePath,
                encoding: 'base64',
                success: (base64Res) => {
                  const base64Data = 'data:image/jpeg;base64,' + base64Res.data;
                  resolve(base64Data);
                },
                fail: (err) => {
                  reject(new Error('读取图片文件失败: ' + JSON.stringify(err)));
                }
              });
            },
            fail: (err) => {
              reject(new Error('导出图片失败: ' + JSON.stringify(err)));
            }
          });
        });
      },
      fail: (err) => {
        reject(new Error('获取图片信息失败: ' + JSON.stringify(err)));
      }
    });
  });
}

/**
 * 上传用户头像
 * @param {String} avatarUrl 微信返回的头像URL
 * @returns {Promise} 返回Promise对象
 */
export function uploadAvatar(avatarUrl) {
  return new Promise(async (resolve, reject) => {
    try {
      // 1. 将头像URL转为base64编码
      const base64Image = await getImageBase64(avatarUrl);
      
      // 2. 调用接口更新头像
      const result = await put(API_PATHS.CLIENT_UPDATE_PROFILE, { 
        avatarUrl: base64Image
      });
      
      resolve(result);
    } catch (err) {
      reject(err);
    }
  });
}

/**
 * 获取用户订单列表
 * @param {Object} params 查询参数
 * @param {Number} params.page 页码
 * @param {Number} params.pageSize 每页条数
 * @param {String} params.status 订单状态：all-全部 process-进行中 unpaid-待支付 completed-已完成 cancelled-已取消 refunded-已退款
 * @returns {Promise} 返回Promise对象
 */
export function getOrderList(params) {
  return get('/wx/client/order/list', params);
}

/**
 * 获取订单详情
 * @param {Object} params 查询参数
 * @param {String} params.orderNo 订单编号
 * @returns {Promise} 返回Promise对象
 */
export function getOrderDetail(params) {
  return get('/wx/client/order/detail', params);
}

/**
 * 取消订单
 * @param {Object} data 请求数据
 * @param {String} data.orderNo 订单编号
 * @param {String} data.reason 取消原因
 * @returns {Promise} 返回Promise对象
 */
export function cancelOrder(data) {
  return post('/wx/client/order/cancel', data);
}

// 获取管家二维码图片
export function getButlerImage() {
  return get('/wx/client/butler/image')
} 