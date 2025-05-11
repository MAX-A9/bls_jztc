"use strict";
const utils_request = require("../utils/request.js");
function getAgreement(params) {
  return utils_request.get("/wx/agreement/get", params);
}
const agreementApi = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  getAgreement
}, Symbol.toStringTag, { value: "Module" }));
exports.agreementApi = agreementApi;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/agreement.js.map
