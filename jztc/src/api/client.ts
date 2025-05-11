import request from '@/utils/request';

// 通用响应结构
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 客户列表项接口
export interface ClientListItem {
  id: number;
  username: string;
  realName: string;
  phone: string;
  status: number;
  statusText: string;
  identifier: string;
  identifierText: string;
  createdAt: string;
  lastLoginAt: string;
  lastLoginIp: string;
}

// 客户列表请求参数
export interface ClientListReq {
  page?: number;
  pageSize?: number;
  username?: string;
  realName?: string;
  phone?: string;
  status?: number;
}

// 客户列表响应
export interface ClientListRes {
  code: number;
  message: string;
  data: {
    list: ClientListItem[];
    total: number;
    page: number;
  };
}

// 创建客户请求参数
export interface ClientCreateReq {
  username: string;
  password: string;
  realName: string;
  phone: string;
  status: number;
  identifier: 'wxapp' | 'unknown';
}

// 创建客户响应
export interface ClientCreateRes {
  id: number;
}

// 更新客户请求参数
export interface ClientUpdateReq {
  id: number;
  username: string;
  realName: string;
  phone: string;
  status: number;
}

// 客户时长列表项接口
export interface ClientDurationItem {
  id: number;
  clientId: number;
  clientName: string;
  remainingDuration: string;
  totalDuration: string;
  usedDuration: string;
  createdAt: string;
  updatedAt: string;
}

// 客户时长列表请求参数
export interface ClientDurationListReq {
  page?: number;
  pageSize?: number;
  clientId?: number;
}

// 客户时长列表响应
export interface ClientDurationListRes {
  code: number;
  message: string;
  data: {
    list: ClientDurationItem[];
    total: number;
    page: number;
  };
}

// 创建客户时长请求参数
export interface ClientDurationCreateReq {
  clientId: number;
  clientName: string;
  totalDuration: string;
  remainingDuration: string;
  usedDuration: string;
}

// 更新客户时长请求参数
export interface ClientDurationUpdateReq {
  id: number;
  clientId: number;
  clientName: string;
  totalDuration: string;
  remainingDuration: string;
  usedDuration: string;
}

/**
 * 获取客户列表
 * @param params 查询参数
 * @returns 客户列表
 */
export const getClientList = (params: ClientListReq): Promise<ClientListRes> => {
  return request.get('/client/list', { params });
};

// 创建客户
export const createClient = (data: ClientCreateReq) => {
  return request<ApiResponse<ClientCreateRes>>({
    url: '/client/create',
    method: 'post',
    data,
  });
};

// 更新客户
export const updateClient = (data: ClientUpdateReq) => {
  return request<ApiResponse<null>>({
    url: '/client/update',
    method: 'put',
    data,
  });
};

// 删除客户
export const deleteClient = (id: number) => {
  return request<ApiResponse<null>>({
    url: '/client/delete',
    method: 'delete',
    params: { id },
  });
};

/**
 * 获取客户时长列表
 * @param params 查询参数
 * @returns 客户时长列表
 */
export const getClientDurationList = (params: ClientDurationListReq): Promise<ClientDurationListRes> => {
  return request.get('/client/duration/list', { params });
};

/**
 * 获取客户时长详情
 * @param id 客户时长ID
 * @returns 客户时长详情
 */
export const getClientDurationDetail = (id: number): Promise<ApiResponse<ClientDurationItem>> => {
  return request.get('/client/duration/detail', { params: { id } });
};

/**
 * 创建客户时长
 * @param data 客户时长信息
 * @returns 创建结果
 */
export const createClientDuration = (data: ClientDurationCreateReq): Promise<ApiResponse<{ id: number }>> => {
  return request({
    url: '/client/duration/create',
    method: 'post',
    data,
  });
};

/**
 * 更新客户时长
 * @param data 客户时长信息
 * @returns 更新结果
 */
export const updateClientDuration = (data: ClientDurationUpdateReq): Promise<ApiResponse<null>> => {
  return request({
    url: '/client/duration/update',
    method: 'put',
    data,
  });
};

/**
 * 删除客户时长
 * @param id 客户时长ID
 * @returns 删除结果
 */
export const deleteClientDuration = (id: number): Promise<ApiResponse<null>> => {
  return request({
    url: '/client/duration/delete',
    method: 'delete',
    params: { id },
  });
}; 