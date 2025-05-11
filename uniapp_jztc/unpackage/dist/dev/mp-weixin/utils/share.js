"use strict";
const common_vendor = require("../common/vendor.js");
const apis_index = require("../apis/index.js");
let shareSettings = null;
async function getShareSettings(forceRefresh = false) {
  if (shareSettings && !forceRefresh) {
    return shareSettings;
  }
  try {
    const result = await apis_index.share.getShareSettings();
    if (result && result.code === 0 && result.data) {
      shareSettings = result.data;
      return shareSettings;
    } else {
      common_vendor.index.__f__("error", "at utils/share.js:28", "获取分享设置失败:", result);
      return null;
    }
  } catch (error) {
    common_vendor.index.__f__("error", "at utils/share.js:32", "获取分享设置出错:", error);
    return null;
  }
}
async function getContentShareOptions(content) {
  const settings = await getShareSettings();
  if (!settings) {
    return {
      title: "分享内容",
      path: "/pages/index/index"
    };
  }
  return {
    title: (content == null ? void 0 : content.title) || settings.content_share_text,
    imageUrl: (content == null ? void 0 : content.cover) || settings.content_share_image,
    path: (content == null ? void 0 : content.id) ? `/pages/content/detail?id=${content.id}` : "/pages/index/index"
  };
}
async function getHomeShareOptions() {
  const settings = await getShareSettings();
  if (!settings) {
    return {
      title: "欢迎访问",
      path: "/pages/index/index"
    };
  }
  return {
    title: settings.home_share_text,
    imageUrl: settings.home_share_image,
    path: "/pages/index/index"
  };
}
async function getDefaultShareOptions(path = "/pages/index/index") {
  const settings = await getShareSettings();
  if (!settings) {
    return {
      title: "欢迎访问",
      path
    };
  }
  return {
    title: settings.default_share_text,
    imageUrl: settings.default_share_image,
    path
  };
}
exports.getContentShareOptions = getContentShareOptions;
exports.getDefaultShareOptions = getDefaultShareOptions;
exports.getHomeShareOptions = getHomeShareOptions;
exports.getShareSettings = getShareSettings;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/share.js.map
