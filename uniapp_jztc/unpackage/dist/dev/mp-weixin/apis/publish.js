"use strict";
const utils_request = require("../utils/request.js");
function createInfo(data) {
  return utils_request.post("/wx/client/content/info/create", data);
}
function createIdle(data) {
  return utils_request.post("/wx/client/content/idle/create", data);
}
function getPackageList() {
  return utils_request.get("/wx/client/package/list", {}, true);
}
function updateContent(id, data) {
  return utils_request.put(`/update/${id}`, data);
}
function deleteContent(id) {
  return utils_request.del(`/delete/${id}`);
}
function getUserPublishList(params) {
  return utils_request.get("/wx/client/content/user/list", params);
}
function getMyPublishList(params) {
  return utils_request.get("/wx/client/publish/list", params);
}
function getPublishCount() {
  return utils_request.get("/wx/client/publish/count");
}
const publish = {
  createInfo,
  createIdle,
  getPackageList,
  updateContent,
  deleteContent,
  getUserPublishList,
  getMyPublishList,
  getPublishCount
};
const publishApi = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  createIdle,
  createInfo,
  default: publish,
  deleteContent,
  getMyPublishList,
  getPackageList,
  getPublishCount,
  getUserPublishList,
  updateContent
}, Symbol.toStringTag, { value: "Module" }));
exports.publishApi = publishApi;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/publish.js.map
