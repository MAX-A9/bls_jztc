/**
 * 支付工具函数
 */
import { unifiedOrder, orderPay } from '@/apis/pay'

/**
 * 调用微信支付
 * @param {Object} data 支付参数
 * @param {String} data.body 商品描述
 * @param {String} data.orderNo 订单号
 * @param {Number} data.totalFee 订单金额
 * @returns {Promise} 支付结果Promise
 */
export function requestWxPay(data) {
  return new Promise(async (resolve, reject) => {
    try {
      // 调用统一下单接口
      const res = await unifiedOrder(data)
      
      if (res.code !== 0 || !res.data) {
        throw new Error(res.message || '支付下单失败')
      }
      
      // 获取支付参数
      const payParams = res.data
      
      // 调用微信支付
      uni.requestPayment({
        provider: 'wxpay',
        timeStamp: payParams.timeStamp,
        nonceStr: payParams.nonceStr,
        package: payParams.package,
        signType: payParams.signType,
        paySign: payParams.paySign,
        success: function(res) {
          resolve({
            success: true,
            message: '支付成功',
            data: res
          })
        },
        fail: function(err) {
          // 判断用户取消
          if (err.errMsg === 'requestPayment:fail cancel') {
            resolve({
              success: false,
              message: '用户取消支付',
              data: err
            })
          } else {
            reject(new Error(err.errMsg || '支付失败'))
          }
        }
      })
    } catch (error) {
      reject(error)
    }
  })
}

/**
 * 调用订单支付
 * @param {Object} data 支付参数
 * @param {String} data.orderNo 订单号
 * @returns {Promise} 支付结果Promise
 */
export function requestOrderPay(data) {
  return new Promise(async (resolve, reject) => {
    try {
      // 调用订单支付接口
      const res = await orderPay(data)
      
      if (res.code !== 0 || !res.data) {
        throw new Error(res.message || '支付下单失败')
      }
      
      // 获取支付参数
      const payParams = res.data
      
      // 调用微信支付
      uni.requestPayment({
        provider: 'wxpay',
        timeStamp: payParams.timeStamp,
        nonceStr: payParams.nonceStr,
        package: payParams.package,
        signType: payParams.signType,
        paySign: payParams.paySign,
        success: function(res) {
          resolve({
            success: true,
            message: '支付成功',
            data: res
          })
        },
        fail: function(err) {
          // 判断用户取消
          if (err.errMsg === 'requestPayment:fail cancel') {
            resolve({
              success: false,
              message: '用户取消支付',
              data: err
            })
          } else {
            reject(new Error(err.errMsg || '支付失败'))
          }
        }
      })
    } catch (error) {
      reject(error)
    }
  })
} 