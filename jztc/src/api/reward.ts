import request from '@/utils/request';

// 奖励设置响应接口
export interface RewardSettingsRes {
  code: number;
  message: string;
  data: {
    enableReward: boolean;
    firstViewMinRewardMin: number;
    firstViewMaxRewardDay: number;
    singleAdMinRewardMin: number;
    singleAdMaxRewardDay: number;
    dailyRewardLimit: number;
    dailyMaxAccumulatedDay: number;
    rewardExpirationDays: number;
  };
}

// 保存奖励设置请求接口
export interface RewardSettingsSaveReq {
  enableReward: boolean;
  firstViewMinRewardMin: number;
  firstViewMaxRewardDay: number;
  singleAdMinRewardMin: number;
  singleAdMaxRewardDay: number;
  dailyRewardLimit: number;
  dailyMaxAccumulatedDay: number;
  rewardExpirationDays: number;
}

// 保存奖励设置响应接口
export interface RewardSettingsSaveRes {
  code: number;
  message: string;
  data: {
    isSuccess: boolean;
  };
}

/**
 * 获取奖励设置
 * @returns 奖励设置响应
 */
export const getRewardSettings = (): Promise<RewardSettingsRes> => {
  return request.get('/settings/reward/settings');
};

/**
 * 保存奖励设置
 * @param data 奖励设置数据
 * @returns 保存结果响应
 */
export const saveRewardSettings = (data: RewardSettingsSaveReq): Promise<RewardSettingsSaveRes> => {
  return request({
    url: '/settings/reward/settings/save',
    method: 'post',
    data,
  });
};

// 协议设置响应接口
export interface AgreementSettingsRes {
  code: number;
  message: string;
  data: {
    privacyPolicy: string;
    userAgreement: string;
  };
}

// 保存协议设置请求接口
export interface AgreementSettingsSaveReq {
  privacyPolicy: string;
  userAgreement: string;
}

// 保存协议设置响应接口
export interface AgreementSettingsSaveRes {
  code: number;
  message: string;
  data: {
    isSuccess: boolean;
  };
}

/**
 * 获取协议设置
 * @returns 协议设置响应
 */
export const getAgreementSettings = (): Promise<AgreementSettingsRes> => {
  return request.get('/settings/agreement/settings');
};

/**
 * 保存协议设置
 * @param data 协议设置数据
 * @returns 保存结果响应
 */
export const saveAgreementSettings = (data: AgreementSettingsSaveReq): Promise<AgreementSettingsSaveRes> => {
  return request({
    url: '/settings/agreement/settings/save',
    method: 'post',
    data,
  });
}; 