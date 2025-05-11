import { post } from '@/utils/request'

/**
 * 微信支付统一下单
 * @param {Object} data 支付参数
 * @param {String} data.body 商品描述
 * @param {String} data.orderNo 订单号
 * @param {Number} data.totalFee 订单金额
 * @returns {Promise} 返回Promise对象
 */
export function unifiedOrder(data) {
  return post('/wx/pay/unified-order', data)
}

/**
 * 订单支付接口
 * @param {Object} data 支付参数
 * @param {String} data.orderNo 订单号
 * @returns {Promise} 返回Promise对象
 */
export function orderPay(data) {
  return post('/wx/client/order/pay', data)
} 