import request from '@/utils/request';

// 导航小程序项接口
export interface MiniProgramItem {
  id: number;
  name: string;
  appId: string;
  logo: string;
  isEnabled: boolean;
  order: number;
}

// 轮播图项接口
export interface BannerItem {
  id: number;
  image: string;
  linkType: string;
  linkUrl: string;
  isEnabled: boolean;
  order: number;
}

// 活动区域项接口
export interface ActivityAreaItem {
  title: string;
  description: string;
  linkType: string;
  linkUrl: string;
}

// 活动区域数据响应
export interface ActivityAreaGetRes {
  code: number;
  message: string;
  data: {
    topLeft: ActivityAreaItem;
    bottomLeft: ActivityAreaItem;
    right: ActivityAreaItem;
    isGlobalEnabled: boolean;
  };
}

// 活动区域保存请求
export interface ActivityAreaSaveReq {
  topLeft: ActivityAreaItem;
  bottomLeft: ActivityAreaItem;
  right: ActivityAreaItem;
}

// 活动区域保存响应
export interface ActivityAreaSaveRes {
  code: number;
  message: string;
  data: null;
}

// 活动区域配置信息
export interface ActivityAreaInfo {
  leftTopTitle: string;
  leftTopDescription: string;
  leftTopLinkType: string;
  leftTopLinkUrl: string;
  leftBottomTitle: string;
  leftBottomDescription: string;
  leftBottomLinkType: string;
  leftBottomLinkUrl: string;
  rightTitle: string;
  rightDescription: string;
  rightLinkType: string;
  rightLinkUrl: string;
}

// 获取导航小程序列表请求
export interface MiniProgramListReq {
  // 无需参数
}

// 获取导航小程序列表响应
export interface MiniProgramListRes {
  code: number;
  message: string;
  data: {
    list: MiniProgramItem[];
    isGlobalEnabled: boolean;
  };
}

// 创建导航小程序请求
export interface MiniProgramCreateReq {
  name: string;
  appId: string;
  logo: string;
  isEnabled: boolean;
  order: number;
}

// 创建导航小程序响应
export interface MiniProgramCreateRes {
  code: number;
  message: string;
  data: {
    id: number;
  };
}

// 更新导航小程序请求
export interface MiniProgramUpdateReq {
  id: number;
  name: string;
  appId: string;
  logo: string;
  isEnabled: boolean;
  order: number;
}

// 更新导航小程序响应
export interface MiniProgramUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 更新导航小程序状态请求
export interface MiniProgramStatusUpdateReq {
  id: number;
  isEnabled: boolean;
}

// 更新导航小程序状态响应
export interface MiniProgramStatusUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 删除导航小程序响应
export interface MiniProgramDeleteRes {
  code: number;
  message: string;
  data: null;
}

/**
 * 获取导航小程序列表
 * @returns 导航小程序列表响应
 */
export const getMiniProgramList = (): Promise<MiniProgramListRes> => {
  return request.get('/content/mini-program/list');
};

/**
 * 创建导航小程序
 * @param data 导航小程序数据
 * @returns 创建结果响应
 */
export const createMiniProgram = (data: MiniProgramCreateReq): Promise<MiniProgramCreateRes> => {
  return request({
    url: '/content/mini-program/create',
    method: 'post',
    data,
  });
};

/**
 * 更新导航小程序
 * @param data 导航小程序数据
 * @returns 更新结果响应
 */
export const updateMiniProgram = (data: MiniProgramUpdateReq): Promise<MiniProgramUpdateRes> => {
  return request({
    url: '/content/mini-program/update',
    method: 'put',
    data,
  });
};

/**
 * 删除导航小程序
 * @param id 导航小程序ID
 * @returns 删除结果响应
 */
export const deleteMiniProgram = (id: number): Promise<MiniProgramDeleteRes> => {
  return request.delete('/content/mini-program/delete', { params: { id } });
};

/**
 * 更新导航小程序状态
 * @param data 状态数据
 * @returns 更新结果响应
 */
export const updateMiniProgramStatus = (data: MiniProgramStatusUpdateReq): Promise<MiniProgramStatusUpdateRes> => {
  return request({
    url: '/content/mini-program/status/update',
    method: 'put',
    data,
  });
};

/**
 * 获取活动区域数据
 * @returns 活动区域数据响应
 */
export const getActivityArea = (): Promise<ActivityAreaGetRes> => {
  return request.get('/content/activity-area/get');
};

/**
 * 保存活动区域数据
 * @param data 活动区域数据
 * @returns 保存结果响应
 */
export const saveActivityArea = (data: ActivityAreaSaveReq): Promise<ActivityAreaSaveRes> => {
  return request({
    url: '/content/activity-area/save',
    method: 'post',
    data,
  });
};

// 获取轮播图列表响应
export interface BannerListRes {
  code: number;
  message: string;
  data: {
    list: BannerItem[];
    isGlobalEnabled: boolean;
  };
}

// 创建轮播图请求
export interface BannerCreateReq {
  image: string;
  linkType: string;
  linkUrl: string;
  isEnabled: boolean;
  order: number;
}

// 创建轮播图响应
export interface BannerCreateRes {
  code: number;
  message: string;
  data: {
    id: number;
  };
}

// 更新轮播图请求
export interface BannerUpdateReq {
  id: number;
  image: string;
  linkType: string;
  linkUrl: string;
  isEnabled: boolean;
  order: number;
}

// 更新轮播图响应
export interface BannerUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 更新轮播图状态请求
export interface BannerStatusUpdateReq {
  id: number;
  isEnabled: boolean;
}

// 更新轮播图状态响应
export interface BannerStatusUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 删除轮播图响应
export interface BannerDeleteRes {
  code: number;
  message: string;
  data: null;
}

/**
 * 获取轮播图列表
 * @returns 轮播图列表响应
 */
export const getBannerList = (): Promise<BannerListRes> => {
  return request.get('/content/banner/list');
};

/**
 * 创建轮播图
 * @param data 轮播图数据
 * @returns 创建结果响应
 */
export const createBanner = (data: BannerCreateReq): Promise<BannerCreateRes> => {
  return request({
    url: '/content/banner/create',
    method: 'post',
    data,
  });
};

/**
 * 更新轮播图
 * @param data 轮播图数据
 * @returns 更新结果响应
 */
export const updateBanner = (data: BannerUpdateReq): Promise<BannerUpdateRes> => {
  return request({
    url: '/content/banner/update',
    method: 'put',
    data,
  });
};

/**
 * 删除轮播图
 * @param id 轮播图ID
 * @returns 删除结果响应
 */
export const deleteBanner = (id: number): Promise<BannerDeleteRes> => {
  return request.delete('/content/banner/delete', { params: { id } });
};

/**
 * 更新轮播图状态
 * @param data 状态数据
 * @returns 更新结果响应
 */
export const updateBannerStatus = (data: BannerStatusUpdateReq): Promise<BannerStatusUpdateRes> => {
  return request({
    url: '/content/banner/status/update',
    method: 'put',
    data,
  });
};

// 导航区域全局状态更新请求
export interface MiniProgramGlobalStatusUpdateReq {
  isEnabled: boolean;
}

// 导航区域全局状态更新响应
export interface MiniProgramGlobalStatusUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 活动区域全局状态更新请求
export interface ActivityAreaGlobalStatusUpdateReq {
  isEnabled: boolean;
}

// 活动区域全局状态更新响应
export interface ActivityAreaGlobalStatusUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 轮播图区域全局状态更新请求
export interface BannerGlobalStatusUpdateReq {
  isEnabled: boolean;
}

// 轮播图区域全局状态更新响应
export interface BannerGlobalStatusUpdateRes {
  code: number;
  message: string;
  data: null;
}

/**
 * 更新导航区域全局状态
 * @param data 全局状态数据
 * @returns 更新结果响应
 */
export const updateMiniProgramGlobalStatus = (data: MiniProgramGlobalStatusUpdateReq): Promise<MiniProgramGlobalStatusUpdateRes> => {
  return request({
    url: '/content/mini-program/global-status/update',
    method: 'put',
    data,
  });
};

/**
 * 更新活动区域全局状态
 * @param data 全局状态数据
 * @returns 更新结果响应
 */
export const updateActivityAreaGlobalStatus = (data: ActivityAreaGlobalStatusUpdateReq): Promise<ActivityAreaGlobalStatusUpdateRes> => {
  return request({
    url: '/content/activity-area/global-status/update',
    method: 'put',
    data,
  });
};

/**
 * 更新轮播图区域全局状态
 * @param data 全局状态数据
 * @returns 更新结果响应
 */
export const updateBannerGlobalStatus = (data: BannerGlobalStatusUpdateReq): Promise<BannerGlobalStatusUpdateRes> => {
  return request({
    url: '/content/banner/global-status/update',
    method: 'put',
    data,
  });
};

// 分享设置接口定义
export interface ShareSettings {
  default_share_text: string;
  default_share_image: string;
  content_share_text: string;
  content_share_image: string;
  home_share_text: string;
  home_share_image: string;
}

// 获取分享设置响应
export interface ShareSettingsRes {
  code: number;
  message: string;
  data: ShareSettings;
}

// 保存分享设置响应
export interface SaveShareSettingsRes {
  code: number;
  message: string;
  data: {
    success: boolean;
  };
}

/**
 * 获取分享设置
 * @returns 分享设置响应
 */
export const getShareSettings = (): Promise<ShareSettingsRes> => {
  return request.get('/settings/share/settings');
};

/**
 * 保存分享设置
 * @param data 分享设置数据
 * @returns 保存结果响应
 */
export const saveShareSettings = (data: ShareSettings): Promise<SaveShareSettingsRes> => {
  return request({
    url: '/settings/share/settings/save',
    method: 'post',
    data,
  });
}; 