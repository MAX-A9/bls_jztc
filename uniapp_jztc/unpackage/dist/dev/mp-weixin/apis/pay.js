"use strict";
const utils_request = require("../utils/request.js");
function unifiedOrder(data) {
  return utils_request.post("/wx/pay/unified-order", data);
}
function orderPay(data) {
  return utils_request.post("/wx/client/order/pay", data);
}
exports.orderPay = orderPay;
exports.unifiedOrder = unifiedOrder;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/pay.js.map
