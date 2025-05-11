"use strict";
const utils_request = require("../utils/request.js");
function getExchangeRecords(params = {}) {
  return utils_request.get("/wx/client/exchange-record/page", params);
}
function getRecentExchangeRecords() {
  return utils_request.get("/wx/exchange-record/list");
}
function createExchangeRecord(data) {
  return utils_request.post("/wx/client/exchange-record/create", data);
}
const vipApi = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  createExchangeRecord,
  getExchangeRecords,
  getRecentExchangeRecords
}, Symbol.toStringTag, { value: "Module" }));
exports.vipApi = vipApi;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/vip.js.map
