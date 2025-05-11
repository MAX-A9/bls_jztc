import request from '../utils/request';

// 兑换记录状态类型
export type ExchangeRecordStatus = 'processing' | 'completed' | 'failed';

// 兑换记录数据接口
export interface ExchangeRecord {
  id: number;
  clientId: number;
  clientName: string;
  productName: string;
  duration: number;
  exchangeTime: string;
  rechargeAccount: string;
  remark: string;
  status: ExchangeRecordStatus;
  createdAt?: string;
  updatedAt?: string;
}

// 兑换记录列表查询参数
export interface ExchangeRecordListParams {
  page: number;
  size: number;
  id?: number;
  clientId?: number;
  status?: ExchangeRecordStatus;
}

// 兑换记录列表响应
export interface ExchangeRecordListRes {
  code: number;
  message: string;
  data: {
    list: ExchangeRecord[];
    total: number;
    page: number;
    size: number;
    pages: number;
  };
}

// 创建兑换记录请求参数
export interface CreateExchangeRecordReq {
  clientId: number;
  clientName: string;
  duration: number;
  exchangeTime: string;
  productName: string;
  rechargeAccount: string;
  remark: string;
  status: ExchangeRecordStatus;
}

// 更新兑换记录请求参数
export interface UpdateExchangeRecordReq {
  id: number;
  clientId: number;
  clientName: string;
  duration: number;
  exchangeTime: string;
  productName: string;
  rechargeAccount: string;
  remark: string;
  status: ExchangeRecordStatus;
}

// 更新兑换记录状态请求参数
export interface UpdateExchangeRecordStatusReq {
  id: number;
  status: ExchangeRecordStatus;
}

// 通用响应接口
export interface CommonResponse {
  code: number;
  message: string;
  data: any;
}

/**
 * 获取兑换记录列表
 * @param params 查询参数
 * @returns 兑换记录列表
 */
export const getExchangeRecordList = (params: ExchangeRecordListParams): Promise<ExchangeRecordListRes> => {
  return request.get('/exchange-record/list', { params });
};

/**
 * 获取兑换记录详情
 * @param id 记录ID
 * @returns 兑换记录详情
 */
export const getExchangeRecord = (id: number): Promise<CommonResponse> => {
  return request.get('/exchange-record/get', { params: { id } });
};

/**
 * 创建兑换记录
 * @param data 兑换记录数据
 * @returns 创建结果
 */
export const createExchangeRecord = (data: CreateExchangeRecordReq): Promise<CommonResponse> => {
  return request.post('/exchange-record/create', data);
};

/**
 * 更新兑换记录
 * @param data 兑换记录数据
 * @returns 更新结果
 */
export const updateExchangeRecord = (data: UpdateExchangeRecordReq): Promise<CommonResponse> => {
  return request.put('/exchange-record/update', data);
};

/**
 * 删除兑换记录
 * @param id 记录ID
 * @returns 删除结果
 */
export const deleteExchangeRecord = (id: number): Promise<CommonResponse> => {
  return request.delete('/exchange-record/delete', { params: { id } });
};

/**
 * 更新兑换记录状态
 * @param data 状态更新数据
 * @returns 更新结果
 */
export const updateExchangeRecordStatus = (data: UpdateExchangeRecordStatusReq): Promise<CommonResponse> => {
  return request.put('/exchange-record/status/update', data);
}; 