"use strict";
const common_vendor = require("../common/vendor.js");
function getFirstLetter(char) {
  if (!char)
    return "";
  if (/[a-zA-Z0-9]/.test(char)) {
    return char.toUpperCase();
  }
  const result = common_vendor.pinyin(char, { pattern: "first", toneType: "none" });
  if (result) {
    return result.toUpperCase();
  }
  return "#";
}
function getStringFirstLetter(str) {
  if (!str)
    return "";
  return getFirstLetter(str.charAt(0));
}
function getPinyin(str, options = {}) {
  if (!str)
    return "";
  const defaultOptions = {
    toneType: "none",
    // 不带声调
    type: "normal",
    // 默认格式
    separator: " "
    // 间隔符
  };
  const finalOptions = { ...defaultOptions, ...options };
  return common_vendor.pinyin(str, finalOptions);
}
exports.getPinyin = getPinyin;
exports.getStringFirstLetter = getStringFirstLetter;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/pinyin.js.map
