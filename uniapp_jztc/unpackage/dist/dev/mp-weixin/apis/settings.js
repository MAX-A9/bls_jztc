"use strict";
const utils_request = require("../utils/request.js");
function getBaseSettings() {
  return utils_request.get("/wx/mini-program/base/settings");
}
const settingsApi = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  getBaseSettings
}, Symbol.toStringTag, { value: "Module" }));
exports.settingsApi = settingsApi;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/settings.js.map
