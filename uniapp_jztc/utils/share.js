/**
 * 分享工具类
 */
import { share as shareApi } from '../apis/index.js';

// 分享设置缓存
let shareSettings = null;

/**
 * 获取分享设置
 * @param {Boolean} forceRefresh 是否强制刷新
 * @returns {Promise} 返回分享设置Promise
 */
export async function getShareSettings(forceRefresh = false) {
  // 如果有缓存且不强制刷新，则直接返回缓存
  if (shareSettings && !forceRefresh) {
    return shareSettings;
  }
  
  try {
    // 调用API获取分享设置
    const result = await shareApi.getShareSettings();
    if (result && result.code === 0 && result.data) {
      // 更新缓存
      shareSettings = result.data;
      return shareSettings;
    } else {
      console.error('获取分享设置失败:', result);
      return null;
    }
  } catch (error) {
    console.error('获取分享设置出错:', error);
    return null;
  }
}

/**
 * 获取内容分享参数
 * @param {Object} content 内容对象
 * @returns {Promise<Object>} 分享参数
 */
export async function getContentShareOptions(content) {
  const settings = await getShareSettings();
  if (!settings) {
    return {
      title: '分享内容',
      path: '/pages/index/index'
    };
  }
  
  // 使用内容信息构建分享参数
  return {
    title: content?.title || settings.content_share_text,
    imageUrl: content?.cover || settings.content_share_image,
    path: content?.id ? `/pages/content/detail?id=${content.id}` : '/pages/index/index'
  };
}

/**
 * 获取首页分享参数
 * @returns {Promise<Object>} 分享参数
 */
export async function getHomeShareOptions() {
  const settings = await getShareSettings();
  if (!settings) {
    return {
      title: '欢迎访问',
      path: '/pages/index/index'
    };
  }
  
  return {
    title: settings.home_share_text,
    imageUrl: settings.home_share_image,
    path: '/pages/index/index'
  };
}

/**
 * 获取默认分享参数
 * @param {String} path 分享路径
 * @returns {Promise<Object>} 分享参数
 */
export async function getDefaultShareOptions(path = '/pages/index/index') {
  const settings = await getShareSettings();
  if (!settings) {
    return {
      title: '欢迎访问',
      path: path
    };
  }
  
  return {
    title: settings.default_share_text,
    imageUrl: settings.default_share_image,
    path: path
  };
}

export default {
  getShareSettings,
  getContentShareOptions,
  getHomeShareOptions,
  getDefaultShareOptions
}; 