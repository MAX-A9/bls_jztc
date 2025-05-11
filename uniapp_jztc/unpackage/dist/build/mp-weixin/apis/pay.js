"use strict";const r=require("../utils/request.js");exports.orderPay=function(e){return r.post("/wx/client/order/pay",e)},exports.unifiedOrder=function(e){return r.post("/wx/pay/unified-order",e)};
