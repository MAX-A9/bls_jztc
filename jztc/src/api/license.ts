import request from '../utils/request';

// 授权验证请求参数
export interface LicenseVerifyReq {
  domain: string;
  license_code: string;
}

// 授权验证响应数据结构
export interface LicenseVerifyData {
  valid: boolean;
  expire_time?: string;
}

// 授权验证响应
export interface LicenseVerifyRes {
  code: number;
  message: string;
  data: LicenseVerifyData;
}

// 授权验证
export function verifyLicense(data: LicenseVerifyReq) {
  return request<LicenseVerifyRes>({
    url: '/license/verify',
    method: 'POST',
    data,
  });
}

// 检查授权状态
export function checkLicenseStatus(domain: string) {
  return request<LicenseVerifyRes>({
    url: '/license/status',
    method: 'GET',
    params: { domain },
  });
} 