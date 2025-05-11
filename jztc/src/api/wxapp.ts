import request from '@/utils/request';

// 定义接口返回类型
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 微信小程序配置类型
export interface WxappConfig {
  appId: string;
  appSecret: string;
  enabled: boolean;
}

// 小程序基本设置类型
export interface MiniProgramBaseSettings {
  name: string;
  description: string;
  logo: string;
}

/**
 * 获取微信小程序配置
 * @returns 微信小程序配置信息
 */
export const getWxappConfig = () => {
  // request的响应拦截器会直接返回response.data
  return request.get<any, ApiResponse<WxappConfig>>('/wxapp/config');
};

/**
 * 保存微信小程序配置
 * @param data 配置信息
 * @returns 保存结果
 */
export const saveWxappConfig = (data: WxappConfig) => {
  // request的响应拦截器会直接返回response.data
  return request.post<any, ApiResponse<null>>('/wxapp/config', data);
};

/**
 * 获取小程序基本设置
 * @returns 小程序基本设置信息
 */
export const getMiniProgramBaseSettings = () => {
  return request.get<any, ApiResponse<MiniProgramBaseSettings>>('/settings/mini-program/base/settings');
};

/**
 * 保存小程序基本设置
 * @param data 设置信息
 * @returns 保存结果
 */
export const saveMiniProgramBaseSettings = (data: MiniProgramBaseSettings) => {
  return request.post<any, ApiResponse<null>>('/settings/mini-program/base/settings/save', data);
}; 