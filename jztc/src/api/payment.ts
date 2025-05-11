import request from '../utils/request';

// 支付配置接口响应
export interface PaymentConfigResponse {
  code?: number;
  message?: string;
  data?: {
    appId: string;
    mchId: string;
    apiKey: string;
    notifyUrl: string;
    isEnabled: boolean;
  };
  appId?: string;
  mchId?: string;
  apiKey?: string;
  notifyUrl?: string;
  isEnabled?: boolean;
}

// 保存配置请求参数
export interface SavePaymentConfigParams {
  appId: string;
  mchId: string;
  apiKey: string;
  notifyUrl: string;
  isEnabled?: boolean;
}

// 保存配置响应
export interface SavePaymentConfigResponse {
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
 * 获取支付配置
 * @returns 支付配置信息
 */
export const getPaymentConfig = (): Promise<PaymentConfigResponse> => {
  return request.get('/payment/config');
};

/**
 * 保存支付配置
 * @param params 支付配置参数
 * @returns 保存结果
 */
export const savePaymentConfig = (params: SavePaymentConfigParams): Promise<SavePaymentConfigResponse> => {
  return request.post('/payment/config', params);
}; 