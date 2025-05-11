import request from '../utils/request';

// 日志列表项类型定义
export interface LogItem {
  id: number;
  action: string;
  module: string;
  details?: string;
  operationIp: string;
  operationResult: number;
  operationTime: string;
  resultText: string;
  username: string;
  userId: number;
}

// 日志列表请求参数
export interface LogListParams {
  page?: number;
  pageSize?: number;
  username?: string;
  action?: string;
  result?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
}

// 日志列表响应
export interface LogListResponse {
  code?: number;
  message?: string;
  data?: {
    list: LogItem[];
    total: number;
    page: number;
  };
  list?: LogItem[];
  total?: number;
  page?: number;
}

// 导出日志响应
export interface LogExportResponse {
  code?: number;
  message?: string;
  data?: {
    url: string;
  };
  url?: string;
}

/**
 * 获取操作日志列表
 * @param params 查询参数
 * @returns 日志列表
 */
export const getOperationLogs = (params: LogListParams): Promise<LogListResponse> => {
  return request.get('/log/list', { params });
};

/**
 * 删除操作日志
 * @param id 日志ID
 * @returns 删除结果
 */
export const deleteOperationLog = (id: number): Promise<void> => {
  return request.delete('/log/delete', { params: { id } });
};

/**
 * 导出操作日志
 * @param params 查询参数
 * @returns 导出结果
 */
export const exportOperationLogs = (params: LogListParams): Promise<LogExportResponse> => {
  return request.get('/log/export', { params });
}; 