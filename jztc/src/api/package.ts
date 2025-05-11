import request from '@/utils/request';

// 通用响应结构
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 套餐类型
export type PackageType = 'top' | 'publish';

// 时长单位类型
export type DurationType = 'hour' | 'day' | 'month';

// 套餐信息结构
export interface Package {
  id: number;
  title: string;
  description: string;
  price: number;
  type: PackageType;
  duration: number;
  durationType: DurationType;
  sortOrder: number;
  is_enabled?: number; // 添加是否启用字段
  createdAt: string;
  updatedAt: string;
}

// 获取套餐列表请求参数
export interface GetPackageListParams {
  type?: PackageType;
  sort?: string; // 排序字段
  order?: 'asc' | 'desc'; // 排序方式
}

// 获取套餐列表响应
export interface GetPackageListRes {
  list: Package[];
  isGlobalEnabled?: boolean; // 添加全局启用状态字段
}

// 创建套餐请求参数
export interface CreatePackageParams {
  title: string;
  description: string;
  price: number;
  type: PackageType;
  duration: number;
  durationType: DurationType;
  sortOrder?: number;
  is_enabled?: number; // 添加是否启用字段
}

// 创建套餐响应
export interface CreatePackageRes {
  id: number;
}

// 更新套餐请求参数
export interface UpdatePackageParams {
  id: number;
  title: string;
  description: string;
  price: number;
  type: PackageType;
  duration: number;
  durationType: DurationType;
  sortOrder?: number;
  is_enabled?: number; // 添加是否启用字段
}

// 设置套餐类型启用状态请求参数
export interface SetPackageTypeEnabledParams {
  isEnabled: boolean; // true启用，false禁用
}

/**
 * 获取套餐列表
 * @param params 查询参数
 * @returns 套餐列表
 */
export const getPackageList = (params: GetPackageListParams): Promise<ApiResponse<GetPackageListRes>> => {
  return request.get('/package/list', { params });
};

/**
 * 获取套餐详情
 * @param id 套餐ID
 * @returns 套餐详情
 */
export const getPackageDetail = (id: number): Promise<ApiResponse<Package>> => {
  return request.get('/package/detail', { params: { id } });
};

/**
 * 创建套餐
 * @param data 套餐数据
 * @returns 创建结果
 */
export const createPackage = (data: CreatePackageParams): Promise<ApiResponse<CreatePackageRes>> => {
  return request.post('/package/create', data);
};

/**
 * 更新套餐
 * @param data 套餐数据
 * @returns 更新结果
 */
export const updatePackage = (data: UpdatePackageParams): Promise<ApiResponse<any>> => {
  return request.put('/package/update', data);
};

/**
 * 删除套餐
 * @param id 套餐ID
 * @returns 删除结果
 */
export const deletePackage = (id: number): Promise<ApiResponse<any>> => {
  return request.delete('/package/delete', { params: { id } });
};

/**
 * 设置套餐类型的启用状态
 * @param data 启用状态数据
 * @returns 设置结果
 */
export const setPackageTypeEnabled = (type: PackageType, data: SetPackageTypeEnabledParams): Promise<ApiResponse<any>> => {
  const url = type === 'top' 
    ? '/package/top/global-status/update' 
    : '/package/publish/global-status/update';
  return request.put(url, data);
}; 