"use strict";
const utils_request = require("../utils/request.js");
function getShareSettings() {
  return utils_request.get("/wx/share/settings");
}
const shareApi = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  getShareSettings
}, Symbol.toStringTag, { value: "Module" }));
exports.shareApi = shareApi;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/share.js.map
