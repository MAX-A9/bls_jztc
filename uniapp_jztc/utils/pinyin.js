/**
 * 汉字转拼音工具
 */
import { pinyin } from 'pinyin-pro';

/**
 * 获取汉字的拼音首字母
 * @param {String} char 单个汉字
 * @returns {String} 返回拼音首字母（大写）
 */
export function getFirstLetter(char) {
  if (!char) return '';
  
  // 如果是英文字母或数字，直接返回大写
  if (/[a-zA-Z0-9]/.test(char)) {
    return char.toUpperCase();
  }
  
  // 使用pinyin-pro获取拼音首字母
  const result = pinyin(char, { pattern: 'first', toneType: 'none' });
  if (result) {
    return result.toUpperCase();
  }
  
  // 无法转换的字符返回#
  return '#';
}

/**
 * 获取字符串的拼音首字母
 * @param {String} str 字符串
 * @returns {String} 返回首字母（大写）
 */
export function getStringFirstLetter(str) {
  if (!str) return '';
  return getFirstLetter(str.charAt(0));
}

/**
 * 获取完整的拼音
 * @param {String} str 字符串
 * @param {Object} options 配置选项
 * @returns {String} 返回拼音
 */
export function getPinyin(str, options = {}) {
  if (!str) return '';
  
  // 默认配置
  const defaultOptions = {
    toneType: 'none', // 不带声调
    type: 'normal',   // 默认格式
    separator: ' '    // 间隔符
  };
  
  // 合并选项
  const finalOptions = { ...defaultOptions, ...options };
  
  return pinyin(str, finalOptions);
}

export default {
  getFirstLetter,
  getStringFirstLetter,
  getPinyin
}; 