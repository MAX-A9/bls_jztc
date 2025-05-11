import request from '@/utils/request';

// 通用响应结构
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 地区列表项接口
export interface RegionListItem {
  id: number;
  name: string;
  location: string;
  level: string;
  status: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

// 地区列表请求参数
export interface RegionListReq {
  page?: number;
  pageSize?: number;
  name?: string;
  location?: string;
  level?: string;
  status?: string;
}

// 地区列表响应
export interface RegionListRes {
  code: number;
  message: string;
  data: {
    list: RegionListItem[];
    total: number;
    page: number;
  };
}

// 地区详情请求参数
export interface RegionDetailReq {
  id: number;
}

// 地区详情响应
export interface RegionDetailRes {
  code: number;
  message: string;
  data: RegionListItem;
}

// 创建地区请求参数
export interface RegionCreateReq {
  name: string;
  location: string;
  level: string;
  status: string;
  sortOrder?: number;
}

// 创建地区响应
export interface RegionCreateRes {
  code: number;
  message: string;
  data: {
    id: number;
  };
}

// 更新地区请求参数
export interface RegionUpdateReq {
  id: number;
  name: string;
  location: string;
  level: string;
  status: string;
  sortOrder?: number;
}

// 更新地区响应
export interface RegionUpdateRes {
  code: number;
  message: string;
  data: any;
}

// 删除地区请求参数
export interface RegionDeleteReq {
  id: number;
}

// 删除地区响应
export interface RegionDeleteRes {
  code: number;
  message: string;
  data: any;
}

/**
 * 获取地区列表
 * @param params 查询参数
 * @returns 地区列表
 */
export const getRegionList = (params: RegionListReq = {}): Promise<RegionListRes> => {
  return request.get('/region/list', { params });
};

/**
 * 获取地区详情
 * @param id 地区ID
 * @returns 地区详情
 */
export const getRegionDetail = (id: number): Promise<RegionDetailRes> => {
  return request.get('/region/detail', { params: { id } });
};

/**
 * 创建地区
 * @param data 地区信息
 * @returns 创建结果
 */
export const createRegion = (data: RegionCreateReq): Promise<RegionCreateRes> => {
  return request.post('/region/create', data);
};

/**
 * 更新地区
 * @param data 地区信息
 * @returns 更新结果
 */
export const updateRegion = (data: RegionUpdateReq): Promise<RegionUpdateRes> => {
  return request.put('/region/update', data);
};

/**
 * 删除地区
 * @param id 地区ID
 * @returns 删除结果
 */
export const deleteRegion = (id: number): Promise<RegionDeleteRes> => {
  return request.delete('/region/delete', { params: { id } });
}; 