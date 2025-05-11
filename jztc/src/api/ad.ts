import request from '@/utils/request';

// 定义接口返回类型
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 广告设置类型
export interface AdSettings {
  enableWxAd: boolean;
  rewardedVideoAdId: string;
}

/**
 * 获取广告设置
 * @returns 广告设置信息
 */
export const getAdSettings = () => {
  return request.get<any, ApiResponse<AdSettings>>('/settings/ad/settings');
};

/**
 * 保存广告设置
 * @param data 广告设置
 * @returns 保存结果
 */
export const saveAdSettings = (data: AdSettings) => {
  return request.post<any, ApiResponse<null>>('/settings/ad/settings/save', data);
}; 