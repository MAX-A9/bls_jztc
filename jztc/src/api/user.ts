import request from '../utils/request';

// 用户列表查询参数
export interface UserListParams {
  page: number;
  pageSize: number;
  username?: string;
  nickname?: string;
  status?: number;
}

// 用户列表项
export interface UserListItem {
  id: number;
  username: string;
  nickname: string;
  status: number;
  statusText: string;
  lastLoginIp: string;
  lastLoginTime: string;
}

// 用户列表响应
export interface UserListResponse {
  code?: number;
  message?: string;
  data?: {
    list: UserListItem[];
    total: number;
    page: number;
  };
  list?: UserListItem[];
  total?: number;
  page?: number;
}

// 创建用户参数
export interface UserCreateParams {
  username: string;
  password: string;
  nickname: string;
}

// 创建用户响应
export interface UserCreateResponse {
  code?: number;
  message?: string;
  data?: {
    id: number;
  };
  id?: number;
}

// 更新用户参数
export interface UserUpdateParams {
  id: number;
  username: string;
  nickname: string;
  status: number;
}

// 删除用户参数
export interface UserDeleteParams {
  id: number;
}

/**
 * 获取用户列表
 * @param params 查询参数
 * @returns 用户列表
 */
export const getUserList = (params: UserListParams): Promise<UserListResponse> => {
  return request.get('/user/list', { params });
};

/**
 * 创建用户
 * @param params 用户信息
 * @returns 创建结果
 */
export const createUser = (params: UserCreateParams): Promise<UserCreateResponse> => {
  return request.post('/user/create', params);
};

/**
 * 更新用户
 * @param params 用户信息
 * @returns 更新结果
 */
export const updateUser = (params: UserUpdateParams): Promise<void> => {
  return request.put('/user/update', params);
};

/**
 * 删除用户
 * @param id 用户ID
 * @returns 删除结果
 */
export const deleteUser = (id: number): Promise<void> => {
  return request.delete('/user/delete', { data: { id } });
}; 