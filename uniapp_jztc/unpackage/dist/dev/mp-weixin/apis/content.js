"use strict";
const utils_request = require("../utils/request.js");
function ensureValidPageNumber(params) {
  if (params && params.page !== void 0) {
    params.page = Math.max(1, parseInt(params.page) || 1);
  }
  return params;
}
function getInfoCategories(type = 1) {
  return utils_request.get("/wx/client/content/categories", { type }, true);
}
function getPublicContentDetail(id) {
  return utils_request.get("/wx/client/content/public/detail", { id }, true);
}
function createInfo(data) {
  return utils_request.post("/wx/client/content/info/create", data);
}
const followUser = (data) => {
  return utils_request.post("/wx/publisher/follow", {
    publisher_id: data.publisher_id
  });
};
const unfollowUser = (data) => {
  return utils_request.post("/wx/publisher/unfollow", {
    publisher_id: data.publisher_id
  });
};
function createIdle(data) {
  return utils_request.post("/wx/client/content/idle/create", data);
}
function getRegionContentList(params) {
  return utils_request.get("/wx/client/content/region/list", ensureValidPageNumber(params), true);
}
function createComment(data) {
  return utils_request.post("/wx/client/comment/create", data);
}
function getCommentList(contentId, params) {
  return utils_request.get("/wx/client/comment/list", { contentId, ...params }, true);
}
function addFavorite(contentId) {
  return utils_request.post("/wx/favorite/add", { contentId });
}
function cancelFavorite(contentId) {
  return utils_request.post("/wx/favorite/cancel", { contentId });
}
function getFavoriteStatus(contentId) {
  return utils_request.get("/wx/favorite/status", { contentId });
}
function getFavoriteList(params) {
  return utils_request.get("/wx/favorite/list", params);
}
function getFavoriteCount() {
  return utils_request.get("/wx/favorite/count");
}
function addBrowseRecord(contentId, contentType = "article") {
  return utils_request.post("/wx/client/browse-history/add", { contentId, contentType });
}
function getBrowseHistoryList(params) {
  return utils_request.get("/wx/client/browse-history/list", params);
}
function clearBrowseHistory(timeType = "all") {
  return utils_request.post("/wx/client/browse-history/clear", { timeType });
}
function getBrowseHistoryCount() {
  return utils_request.get("/wx/client/browse-history/count");
}
function getPublisherInfo(publisherId) {
  return utils_request.get("/wx/publisher/info", {
    publisher_id: publisherId
  }, true);
}
function getPublisherFollowStatus(publisherId) {
  return utils_request.get("/wx/publisher/follow/status", {
    publisher_id: publisherId
  });
}
function getFollowingList(params) {
  return utils_request.get("/wx/publisher/following/list", {
    page: params.page || 1,
    size: params.size || 10
  });
}
function getFollowingCount() {
  return utils_request.get("/wx/publisher/following/count");
}
function getRegionIdleList(params) {
  return utils_request.get("/wx/client/content/region/idle/list", ensureValidPageNumber(params), true);
}
function getPackageList() {
  return utils_request.get("/wx/client/package/list", {}, true);
}
function getMiniProgramList() {
  return utils_request.get("/wx/mini-program/list", {}, true);
}
exports.addBrowseRecord = addBrowseRecord;
exports.addFavorite = addFavorite;
exports.cancelFavorite = cancelFavorite;
exports.clearBrowseHistory = clearBrowseHistory;
exports.createComment = createComment;
exports.createIdle = createIdle;
exports.createInfo = createInfo;
exports.followUser = followUser;
exports.getBrowseHistoryCount = getBrowseHistoryCount;
exports.getBrowseHistoryList = getBrowseHistoryList;
exports.getCommentList = getCommentList;
exports.getFavoriteCount = getFavoriteCount;
exports.getFavoriteList = getFavoriteList;
exports.getFavoriteStatus = getFavoriteStatus;
exports.getFollowingCount = getFollowingCount;
exports.getFollowingList = getFollowingList;
exports.getInfoCategories = getInfoCategories;
exports.getMiniProgramList = getMiniProgramList;
exports.getPackageList = getPackageList;
exports.getPublicContentDetail = getPublicContentDetail;
exports.getPublisherFollowStatus = getPublisherFollowStatus;
exports.getPublisherInfo = getPublisherInfo;
exports.getRegionContentList = getRegionContentList;
exports.getRegionIdleList = getRegionIdleList;
exports.unfollowUser = unfollowUser;
//# sourceMappingURL=../../.sourcemap/mp-weixin/apis/content.js.map
