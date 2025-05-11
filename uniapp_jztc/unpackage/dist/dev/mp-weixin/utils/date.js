"use strict";
const common_vendor = require("../common/vendor.js");
function getTimestamp(time) {
  if (!time)
    return Date.now();
  if (typeof time === "number")
    return time;
  if (time instanceof Date)
    return time.getTime();
  if (typeof time === "string") {
    time = time.trim();
    if (/^\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d{1,2}:\d{1,2}$/.test(time)) {
      const slashStr = time.replace(/-/g, "/");
      const date1 = new Date(slashStr);
      if (!isNaN(date1.getTime())) {
        return date1.getTime();
      }
      const isoStr = time.replace(" ", "T");
      const date2 = new Date(isoStr);
      if (!isNaN(date2.getTime())) {
        return date2.getTime();
      }
    }
    const directDate = new Date(time);
    if (!isNaN(directDate.getTime())) {
      return directDate.getTime();
    }
    const regex = /^(\d{4})[/-](\d{1,2})[/-](\d{1,2})(?:[ T](\d{1,2}):(\d{1,2})(?::(\d{1,2}))?)?$/;
    if (regex.test(time)) {
      const parts = time.match(regex);
      const year = parseInt(parts[1], 10);
      const month = parseInt(parts[2], 10) - 1;
      const day = parseInt(parts[3], 10);
      const hour = parts[4] ? parseInt(parts[4], 10) : 0;
      const minute = parts[5] ? parseInt(parts[5], 10) : 0;
      const second = parts[6] ? parseInt(parts[6], 10) : 0;
      const manualDate = new Date(year, month, day, hour, minute, second);
      return manualDate.getTime();
    }
  }
  common_vendor.index.__f__("warn", "at utils/date.js:63", "无法解析日期格式:", time);
  return Date.now();
}
function formatTimeAgo(time, lang = "zh") {
  if (!time)
    return "";
  const timestamp = getTimestamp(time);
  const now = Date.now();
  const diff = (now - timestamp) / 1e3;
  const units = {
    zh: {
      second: "秒",
      minute: "分钟",
      hour: "小时",
      day: "天",
      week: "周",
      month: "个月",
      year: "年"
    },
    en: {
      second: "second",
      minute: "minute",
      hour: "hour",
      day: "day",
      week: "week",
      month: "month",
      year: "year"
    }
  };
  const unit = units[lang] || units.zh;
  const pluralSuffix = lang === "en" ? "s" : "";
  if (diff < 5) {
    return lang === "zh" ? "刚刚" : "just now";
  } else if (diff < 60) {
    return `${Math.floor(diff)}${unit.second}${lang === "en" && Math.floor(diff) > 1 ? pluralSuffix : ""}${lang === "zh" ? "前" : " ago"}`;
  } else if (diff < 3600) {
    const minutes = Math.floor(diff / 60);
    return `${minutes}${unit.minute}${lang === "en" && minutes > 1 ? pluralSuffix : ""}${lang === "zh" ? "前" : " ago"}`;
  } else if (diff < 86400) {
    const hours = Math.floor(diff / 3600);
    return `${hours}${unit.hour}${lang === "en" && hours > 1 ? pluralSuffix : ""}${lang === "zh" ? "前" : " ago"}`;
  } else if (diff < 604800) {
    const days = Math.floor(diff / 86400);
    return `${days}${unit.day}${lang === "en" && days > 1 ? pluralSuffix : ""}${lang === "zh" ? "前" : " ago"}`;
  } else if (diff < 2592e3) {
    const weeks = Math.floor(diff / 604800);
    return `${weeks}${unit.week}${lang === "en" && weeks > 1 ? pluralSuffix : ""}${lang === "zh" ? "前" : " ago"}`;
  } else if (diff < 31536e3) {
    const months = Math.floor(diff / 2592e3);
    return `${months}${unit.month}${lang === "en" && months > 1 ? pluralSuffix : ""}${lang === "zh" ? "前" : " ago"}`;
  } else {
    const years = Math.floor(diff / 31536e3);
    return `${years}${unit.year}${lang === "en" && years > 1 ? pluralSuffix : ""}${lang === "zh" ? "前" : " ago"}`;
  }
}
exports.formatTimeAgo = formatTimeAgo;
exports.getTimestamp = getTimestamp;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/date.js.map
