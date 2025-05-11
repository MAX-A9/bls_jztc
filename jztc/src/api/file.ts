import request from '@/utils/request';

// 通用响应结构
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 文件信息接口
export interface FileInfo {
  id: number;
  name: string;
  path: string;
  url: string;
  size: number;
  sizeFormat: string;
  type: string;
  extension: string;
  contentType: string;
  isPublic: boolean;
  userId: number;
  username: string;
  createdAt: string;
}

// 文件列表请求参数
export interface FileListReq {
  Page?: number;
  PageSize?: number;
  Keyword?: string;
  Type?: string;
  IsPublic?: boolean;
}

// 文件列表响应
export interface FileListRes {
  code: number;
  message: string;
  data: {
    list: FileInfo[];
    total: number;
    page: number;
  };
}

// 文件上传请求参数
export interface FileUploadReq {
  File: File;
  Directory?: string;
  IsPublic?: boolean;
}

// 文件上传响应
export interface FileUploadRes {
  id: number;
  name: string;
  extension: string;
  size: number;
  type: string;
  url: string;
  isPublic: boolean;
  uploadTime: string;
}

// 文件更新公开状态请求参数
export interface FileUpdatePublicReq {
  Id: number;
  IsPublic: boolean;
}

// 文件更新公开状态响应
export interface FileUpdatePublicRes {
  success: boolean;
  message: string;
  url: string;
}

// 文件删除响应
export interface FileDeleteRes {
  success: boolean;
  message: string;
}

// 文件批量删除响应
export interface FileBatchDeleteRes {
  success: boolean;
  message: string;
  count: number;
}

/**
 * 获取文件列表
 * @param params 查询参数
 * @returns 文件列表
 */
export const getFileList = (params: FileListReq): Promise<FileListRes> => {
  return request.get('/file/list', { params });
};

/**
 * 获取文件详情
 * @param id 文件ID
 * @returns 文件详情
 */
export const getFileDetail = (id: number): Promise<ApiResponse<FileInfo>> => {
  return request.get('/file/detail', { params: { Id: id } });
};

/**
 * 上传文件
 * @param data 上传文件数据
 * @returns 上传结果
 */
export const uploadFile = (data: FormData): Promise<ApiResponse<FileUploadRes>> => {
  return request({
    url: '/file/upload',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

/**
 * 更新文件公开状态
 * @param data 请求参数
 * @returns 更新结果
 */
export const updateFilePublic = (data: FileUpdatePublicReq): Promise<ApiResponse<FileUpdatePublicRes>> => {
  return request({
    url: '/file/update-public',
    method: 'put',
    data,
  });
};

/**
 * 删除文件
 * @param id 文件ID
 * @returns 删除结果
 */
export const deleteFile = (id: number): Promise<ApiResponse<FileDeleteRes>> => {
  return request({
    url: '/file/delete',
    method: 'delete',
    params: { Id: id },
  });
};

/**
 * 批量删除文件
 * @param ids 文件ID列表
 * @returns 批量删除结果
 */
export const batchDeleteFiles = (ids: number[]): Promise<ApiResponse<FileBatchDeleteRes>> => {
  return request({
    url: '/file/batch-delete',
    method: 'delete',
    params: { Ids: ids },
  });
}; 