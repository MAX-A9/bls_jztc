/**
 * 环境变量工具类
 * 提供对环境变量的统一访问方法
 */

// 环境变量前缀，Vite中默认以VITE_开头的变量才会暴露给客户端代码
const ENV_PREFIX = 'VITE_';

/**
 * 获取环境变量
 * @param key 环境变量名称，不需要包含前缀
 * @param defaultValue 默认值，当环境变量不存在时返回
 * @returns 环境变量值
 */
export function getEnvVariable<T = string>(key: string, defaultValue?: T): T {
  const envKey = ENV_PREFIX + key;
  const envValue = import.meta.env[envKey] as unknown;
  
  if (envValue === undefined) {
    return defaultValue as T;
  }
  
  // 处理布尔值
  if (envValue === 'true') {
    return true as unknown as T;
  }
  
  if (envValue === 'false') {
    return false as unknown as T;
  }
  
  // 处理数字
  if (!Number.isNaN(Number(envValue)) && typeof envValue === 'string' && envValue.trim() !== '') {
    return Number(envValue) as unknown as T;
  }
  
  return envValue as T;
}

/**
 * 判断当前是否为开发环境
 * @returns 是否为开发环境
 */
export function isDevelopment(): boolean {
  return import.meta.env.DEV;
}

/**
 * 判断当前是否为生产环境
 * @returns 是否为生产环境
 */
export function isProduction(): boolean {
  return import.meta.env.PROD;
}

/**
 * 获取API基础URL
 * @returns API基础URL
 */
export function getApiBaseUrl(): string {
  return getEnvVariable('API_BASE_URL', '/api');
}

/**
 * 获取应用名称
 * @returns 应用名称
 */
export function getAppName(): string {
  return getEnvVariable('APP_NAME', '管理系统');
}

/**
 * 获取API服务器地址
 * @returns API服务器地址
 */
export function getApiServer(): string {
  return getEnvVariable('API_SERVER', '');
}

// 导出环境变量对象，方便直接访问
export const ENV = {
  APP_NAME: getAppName(),
  API_BASE_URL: getApiBaseUrl(),
  API_SERVER: getApiServer(),
  IS_DEV: isDevelopment(),
  IS_PROD: isProduction(),
}; 