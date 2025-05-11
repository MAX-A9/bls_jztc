import request from '@/utils/request';

// 通用响应结构
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 内容列表项接口
export interface ContentItem {
  id: number;
  title: string;
  category: string;
  author: string;
  content: string;
  status: string;
  views: number;
  likes: number;
  comments: number;
  isRecommended: boolean;
  publishedAt: string;
  expiresAt: string | null;
  topUntil: string | null;
  createdAt: string;
  updatedAt: string;
}

// 内容列表请求参数
export interface ContentListReq {
  page?: number;
  pageSize?: number;
  title?: string;
  category?: string;
  status?: string;
  author?: string;
}

// 内容列表响应
export interface ContentListRes {
  code: number;
  message: string;
  data: {
    list: ContentItem[];
    total: number;
    page: number;
  };
}

// 内容详情响应
export interface ContentDetailRes {
  code: number;
  message: string;
  data: ContentItem;
}

// 创建内容请求参数
export interface ContentCreateReq {
  title: string;
  category: string;
  author: string;
  content: string;
  status: string;
  isRecommended: boolean;
  topUntil?: string | null;
}

// 创建内容响应
export interface ContentCreateRes {
  code: number;
  message: string;
  data: {
    id: number;
  };
}

// 更新内容请求参数
export interface ContentUpdateReq {
  id: number;
  title: string;
  category: string;
  author: string;
  content: string;
  status: string;
  isRecommended: boolean;
  topUntil?: string | null;
}

// 更新内容响应
export interface ContentUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 删除内容响应
export interface ContentDeleteRes {
  code: number;
  message: string;
  data: null;
}

// 更新内容状态请求参数
export interface ContentStatusUpdateReq {
  id: number;
  status: string;
}

// 更新内容状态响应
export interface ContentStatusUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 更新内容推荐状态请求参数
export interface ContentRecommendUpdateReq {
  id: number;
  isRecommended: boolean;
  topUntil?: string | null;
}

// 更新内容推荐状态响应
export interface ContentRecommendUpdateRes {
  code: number;
  message: string;
  data: null;
}

// 内页轮播图相关接口
// 内页轮播图条目接口
export interface InnerBannerItem {
  id: number;
  image: string;
  linkType: string;
  linkUrl: string;
  isEnabled: boolean;
  order: number;
  bannerType: string;
}

// 更新内页轮播图接口参数
export interface InnerBannerUpdateReq {
  id: number;
  bannerType: string;
  image: string;
  linkType: string;
  linkUrl: string;
  isEnabled: boolean;
  order: number;
}

// 创建内页轮播图请求参数
export interface InnerBannerCreateReq {
  bannerType: string; // 'home' 或 'idle'
  image: string;
  isEnabled: boolean;
  linkType: string;
  linkUrl: string;
  order: number;
}

// 创建内页轮播图
export const createInnerBanner = (data: InnerBannerCreateReq) => {
  return request.post<any, ApiResponse<{ id: number }>>('/content/inner-banner/create', data);
};

// 获取内页轮播图列表
export const getInnerBannerList = (bannerType: string) => {
  return request.get<any, ApiResponse<{ list: InnerBannerItem[], isGlobalEnabled: boolean }>>('/content/inner-banner/list', {
    params: { bannerType }
  });
};

// 更新内页轮播图
export const updateInnerBanner = (data: InnerBannerUpdateReq) => {
  return request.put<any, ApiResponse<null>>('/content/inner-banner/update', data);
};

// 删除内页轮播图
export const deleteInnerBanner = (id: number) => {
  return request.delete<any, ApiResponse<null>>('/content/inner-banner/delete', {
    data: { id }
  });
};

// 更新内页轮播图状态
export const updateInnerBannerStatus = (id: number, isEnabled: boolean) => {
  return request.put<any, ApiResponse<null>>('/content/inner-banner/status/update', {
    id,
    isEnabled
  });
};

// 更新首页轮播图全局开关状态
export const updateHomeInnerBannerGlobalStatus = (isEnabled: boolean) => {
  return request.put<any, ApiResponse<null>>('/content/inner-banner/home/global-status/update', {
    isEnabled
  });
};

// 更新闲置轮播图全局开关状态
export const updateIdleInnerBannerGlobalStatus = (isEnabled: boolean) => {
  return request.put<any, ApiResponse<null>>('/content/inner-banner/idle/global-status/update', {
    isEnabled
  });
};

// 底部Tab相关接口

// 底部Tab项接口定义
export interface BottomTabItem {
  id: number;
  name: string;
  icon: string;
  selectedIcon: string;
  path: string;
  order: number;
  isEnabled: boolean;
}

/**
 * 获取底部Tab列表
 */
export const getBottomTabList = async () => {
  return await request.get<any, ApiResponse<{
    list: BottomTabItem[];
  }>>('/content/bottom-tab/list');
};

/**
 * 更新底部Tab状态
 * @param id Tab ID
 * @param isEnabled 是否启用
 */
export const updateBottomTabStatus = async (id: number, isEnabled: boolean) => {
  return await request.put<any, ApiResponse<null>>('/content/bottom-tab/status/update', {
    id,
    isEnabled
  });
};

/**
 * 更新底部Tab信息
 * @param data Tab数据
 */
export const updateBottomTab = async (data: {
  id: number;
  name: string;
  path: string;
  icon: string;
  selectedIcon: string;
  order: number;
  isEnabled: boolean;
}) => {
  return await request.put<any, ApiResponse<null>>('/content/bottom-tab/update', data);
};

/**
 * 获取内容列表
 * @param params 查询参数
 * @returns 内容列表
 */
export const getContentList = (params: ContentListReq): Promise<ContentListRes> => {
  return request.get('/list', { params });
};

/**
 * 获取内容详情
 * @param id 内容ID
 * @returns 内容详情
 */
export const getContentDetail = (id: number): Promise<ContentDetailRes> => {
  return request.get('/detail', { params: { id } });
};

/**
 * 创建内容
 * @param data 内容数据
 * @returns 创建结果
 */
export const createContent = (data: ContentCreateReq): Promise<ContentCreateRes> => {
  return request({
    url: '/create',
    method: 'post',
    data,
  });
};

/**
 * 更新内容
 * @param data 内容数据
 * @returns 更新结果
 */
export const updateContent = (data: ContentUpdateReq): Promise<ContentUpdateRes> => {
  return request({
    url: '/update',
    method: 'put',
    data,
  });
};

/**
 * 删除内容
 * @param id 内容ID
 * @returns 删除结果
 */
export const deleteContent = (id: number): Promise<ContentDeleteRes> => {
  return request.delete('/delete', { params: { id } });
};

/**
 * 更新内容状态
 * @param data 状态数据
 * @returns 更新结果
 */
export const updateContentStatus = (data: ContentStatusUpdateReq): Promise<ContentStatusUpdateRes> => {
  return request({
    url: '/status/update',
    method: 'put',
    data,
  });
};

/**
 * 更新内容推荐状态
 * @param data 推荐状态数据
 * @returns 更新结果
 */
export const updateContentRecommend = (data: ContentRecommendUpdateReq): Promise<ContentRecommendUpdateRes> => {
  return request({
    url: '/recommend/update',
    method: 'put',
    data,
  });
}; 