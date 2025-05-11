import request from '../utils/request';

// 验证码返回数据类型
export interface CaptchaResponse {
  code?: number;
  message?: string;
  data?: {
    id: string;
    base64: string;
    expiredAt: string;
  };
  id?: string;
  base64?: string;
  expiredAt?: string;
}

// 登录返回数据类型
export interface LoginResponse {
  code?: number;
  message?: string;
  data?: {
    token: string;
    userId: number;
    nickname: string;
    expireIn: number;
  };
  token?: string;
  userId?: number;
  nickname?: string;
  expireIn?: number;
}

// 登录请求参数类型
export interface LoginParams {
  username: string;
  password: string;
  captchaId: string;
  captchaCode: string;
}

// 用户信息返回数据类型
export interface UserInfoResponse {
  code?: number;
  message?: string;
  data?: {
    id: number;
    username: string;
    nickname: string;
    role: string;
    status: number;
    lastLogin: string;
  };
  id?: number;
  username?: string;
  nickname?: string;
  role?: string;
  status?: number;
  lastLogin?: string;
}

// 退出登录响应类型
export interface LogoutResponse {
  code?: number;
  message?: string;
}

/**
 * 获取验证码
 * @returns 验证码信息
 */
export const getCaptcha = (): Promise<CaptchaResponse> => {
  return request.get('/captcha');
};

/**
 * 用户登录
 * @param params 登录参数
 * @returns 登录结果
 */
export const login = (params: LoginParams): Promise<LoginResponse> => {
  return request.post('/login', params);
};

/**
 * 获取当前登录用户信息
 * @returns 用户信息
 */
export const getUserInfo = (): Promise<UserInfoResponse> => {
  return request.get('/info');
};

/**
 * 退出登录
 * @returns 退出结果
 */
export const logout = async (): Promise<void> => {
  try {
    // 调用后端退出登录接口
    await request.post('/logout');
  } finally {
    // 无论API调用成功与否，都清除本地存储的token
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    localStorage.removeItem('username');
    localStorage.removeItem('nickname');
  }
}; 