/**
 * 将日期字符串转换为时间戳，兼容iOS
 * @param {String|Number|Date} time 需要转换的时间
 * @return {Number} 时间戳
 */
export function getTimestamp(time) {
  // 若参数为空，返回当前时间戳
  if (!time) return Date.now();
  
  // 如果已经是时间戳，直接返回
  if (typeof time === 'number') return time;
  
  // 如果是Date对象，返回其时间戳
  if (time instanceof Date) return time.getTime();
  
  // 字符串类型需要特殊处理
  if (typeof time === 'string') {
    // 标准化日期格式，移除额外空格
    time = time.trim();
    
    // 处理 yyyy-MM-dd HH:mm:ss 格式
    if (/^\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d{1,2}:\d{1,2}$/.test(time)) {
      // iOS兼容方案1: 转换为 yyyy/MM/dd HH:mm:ss
      const slashStr = time.replace(/-/g, '/');
      const date1 = new Date(slashStr);
      if (!isNaN(date1.getTime())) {
        return date1.getTime();
      }
      
      // iOS兼容方案2: 转换为 yyyy-MM-ddTHH:mm:ss
      const isoStr = time.replace(' ', 'T');
      const date2 = new Date(isoStr);
      if (!isNaN(date2.getTime())) {
        return date2.getTime();
      }
    }
    
    // 尝试直接解析，某些格式在某些设备上可能有效
    const directDate = new Date(time);
    if (!isNaN(directDate.getTime())) {
      return directDate.getTime();
    }
    
    // 手动解析日期时间
    // 匹配 yyyy-MM-dd 或 yyyy/MM/dd 格式，可选时间部分
    const regex = /^(\d{4})[/-](\d{1,2})[/-](\d{1,2})(?:[ T](\d{1,2}):(\d{1,2})(?::(\d{1,2}))?)?$/;
    if (regex.test(time)) {
      const parts = time.match(regex);
      const year = parseInt(parts[1], 10);
      const month = parseInt(parts[2], 10) - 1; // 月份从0开始
      const day = parseInt(parts[3], 10);
      const hour = parts[4] ? parseInt(parts[4], 10) : 0;
      const minute = parts[5] ? parseInt(parts[5], 10) : 0;
      const second = parts[6] ? parseInt(parts[6], 10) : 0;
      
      // 使用本地时间创建Date对象
      const manualDate = new Date(year, month, day, hour, minute, second);
      return manualDate.getTime();
    }
  }
  
  // 所有方法都失败，使用当前时间并记录警告
  console.warn('无法解析日期格式:', time);
  return Date.now();
}

/**
 * 将时间格式化为多久前
 * @param {String|Number|Date} time 需要格式化的时间
 * @param {String} lang 语言，默认zh为中文，en为英文
 * @return {String} 格式化后的时间
 */
export function formatTimeAgo(time, lang = 'zh') {
  if (!time) return '';
  
  // 获取时间戳
  const timestamp = getTimestamp(time);
  
  const now = Date.now();
  const diff = (now - timestamp) / 1000; // 转为秒
  
  // 定义时间单位
  const units = {
    zh: {
      second: '秒',
      minute: '分钟',
      hour: '小时',
      day: '天',
      week: '周',
      month: '个月',
      year: '年'
    },
    en: {
      second: 'second',
      minute: 'minute',
      hour: 'hour',
      day: 'day',
      week: 'week',
      month: 'month',
      year: 'year'
    }
  };
  
  // 获取对应语言的单位
  const unit = units[lang] || units.zh;
  const pluralSuffix = lang === 'en' ? 's' : ''; // 英文复数形式
  
  // 定义不同时间段的显示格式
  if (diff < 5) {
    return lang === 'zh' ? '刚刚' : 'just now';
  } else if (diff < 60) {
    return `${Math.floor(diff)}${unit.second}${lang === 'en' && Math.floor(diff) > 1 ? pluralSuffix : ''}${lang === 'zh' ? '前' : ' ago'}`;
  } else if (diff < 3600) {
    const minutes = Math.floor(diff / 60);
    return `${minutes}${unit.minute}${lang === 'en' && minutes > 1 ? pluralSuffix : ''}${lang === 'zh' ? '前' : ' ago'}`;
  } else if (diff < 86400) {
    const hours = Math.floor(diff / 3600);
    return `${hours}${unit.hour}${lang === 'en' && hours > 1 ? pluralSuffix : ''}${lang === 'zh' ? '前' : ' ago'}`;
  } else if (diff < 604800) {
    const days = Math.floor(diff / 86400);
    return `${days}${unit.day}${lang === 'en' && days > 1 ? pluralSuffix : ''}${lang === 'zh' ? '前' : ' ago'}`;
  } else if (diff < 2592000) {
    const weeks = Math.floor(diff / 604800);
    return `${weeks}${unit.week}${lang === 'en' && weeks > 1 ? pluralSuffix : ''}${lang === 'zh' ? '前' : ' ago'}`;
  } else if (diff < 31536000) {
    const months = Math.floor(diff / 2592000);
    return `${months}${unit.month}${lang === 'en' && months > 1 ? pluralSuffix : ''}${lang === 'zh' ? '前' : ' ago'}`;
  } else {
    const years = Math.floor(diff / 31536000);
    return `${years}${unit.year}${lang === 'en' && years > 1 ? pluralSuffix : ''}${lang === 'zh' ? '前' : ' ago'}`;
  }
}

/**
 * 格式化日期为指定格式
 * @param {String|Number|Date} time 需要格式化的时间
 * @param {String} format 格式化的格式，默认为 YYYY-MM-DD
 * @return {String} 格式化后的时间
 */
export function formatDate(time, format = 'YYYY-MM-DD') {
  if (!time) return '';
  
  // 获取时间戳并转换为Date对象
  const timestamp = getTimestamp(time);
  const date = new Date(timestamp);
  
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const hour = date.getHours();
  const minute = date.getMinutes();
  const second = date.getSeconds();
  
  // 补零函数
  const pad = (n) => n < 10 ? '0' + n : n;
  
  return format
    .replace(/YYYY/g, year)
    .replace(/YY/g, String(year).slice(2))
    .replace(/MM/g, pad(month))
    .replace(/M/g, month)
    .replace(/DD/g, pad(day))
    .replace(/D/g, day)
    .replace(/HH/g, pad(hour))
    .replace(/H/g, hour)
    .replace(/hh/g, pad(hour % 12 || 12))
    .replace(/h/g, hour % 12 || 12)
    .replace(/mm/g, pad(minute))
    .replace(/m/g, minute)
    .replace(/ss/g, pad(second))
    .replace(/s/g, second);
} 