import request from '../utils/request';

// 存储配置接口响应
export interface StorageConfigResponse {
  code?: number;
  message?: string;
  data?: {
    accessKeyId: string;
    accessKeySecret: string;
    endpoint: string;
    bucket: string;
    region: string;
    directory: string;
    publicAccess: boolean;
  };
  accessKeyId?: string;
  accessKeySecret?: string;
  endpoint?: string;
  bucket?: string;
  region?: string;
  directory?: string;
  publicAccess?: boolean;
}

// 保存配置请求参数
export interface SaveStorageConfigParams {
  accessKeyId: string;
  accessKeySecret: string;
  endpoint: string;
  bucket: string;
  region: string;
  directory?: string;
  publicAccess?: boolean;
}

// 保存配置响应
export interface SaveStorageConfigResponse {
  code?: number;
  message?: string;
  data?: {
    success: boolean;
    responseMessage: string;
  };
  success?: boolean;
  responseMessage?: string;
}

/**
 * 获取存储配置
 * @returns 存储配置信息
 */
export const getStorageConfig = (): Promise<StorageConfigResponse> => {
  return request.get('/storage/config');
};

/**
 * 保存存储配置
 * @param params 存储配置参数
 * @returns 保存结果
 */
export const saveStorageConfig = (params: SaveStorageConfigParams): Promise<SaveStorageConfigResponse> => {
  return request.post('/storage/config', params);
}; 