"use strict";
const common_vendor = require("../common/vendor.js");
const apis_pay = require("../apis/pay.js");
function requestWxPay(data) {
  return new Promise(async (resolve, reject) => {
    try {
      const res = await apis_pay.unifiedOrder(data);
      if (res.code !== 0 || !res.data) {
        throw new Error(res.message || "支付下单失败");
      }
      const payParams = res.data;
      common_vendor.index.requestPayment({
        provider: "wxpay",
        timeStamp: payParams.timeStamp,
        nonceStr: payParams.nonceStr,
        package: payParams.package,
        signType: payParams.signType,
        paySign: payParams.paySign,
        success: function(res2) {
          resolve({
            success: true,
            message: "支付成功",
            data: res2
          });
        },
        fail: function(err) {
          if (err.errMsg === "requestPayment:fail cancel") {
            resolve({
              success: false,
              message: "用户取消支付",
              data: err
            });
          } else {
            reject(new Error(err.errMsg || "支付失败"));
          }
        }
      });
    } catch (error) {
      reject(error);
    }
  });
}
function requestOrderPay(data) {
  return new Promise(async (resolve, reject) => {
    try {
      const res = await apis_pay.orderPay(data);
      if (res.code !== 0 || !res.data) {
        throw new Error(res.message || "支付下单失败");
      }
      const payParams = res.data;
      common_vendor.index.requestPayment({
        provider: "wxpay",
        timeStamp: payParams.timeStamp,
        nonceStr: payParams.nonceStr,
        package: payParams.package,
        signType: payParams.signType,
        paySign: payParams.paySign,
        success: function(res2) {
          resolve({
            success: true,
            message: "支付成功",
            data: res2
          });
        },
        fail: function(err) {
          if (err.errMsg === "requestPayment:fail cancel") {
            resolve({
              success: false,
              message: "用户取消支付",
              data: err
            });
          } else {
            reject(new Error(err.errMsg || "支付失败"));
          }
        }
      });
    } catch (error) {
      reject(error);
    }
  });
}
exports.requestOrderPay = requestOrderPay;
exports.requestWxPay = requestWxPay;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/pay.js.map
