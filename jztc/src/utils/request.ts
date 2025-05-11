import axios from 'axios';
import { message } from 'antd';
import { ENV } from './env';

// 创建axios实例
const request = axios.create({
  baseURL: ENV.API_BASE_URL, // 使用环境变量中配置的API基础URL
  timeout: 10000, // 请求超时时间
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    
    // 如果有token则添加到请求头
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data;
    
    // 如果返回的状态码不是0，说明接口请求有问题，直接抛出错误
    if (res.code !== 0) {
      message.error(res.message || '请求失败');
      
      // 401: 未登录或token过期
      if (res.code === 401) {
        // 清除token
        localStorage.removeItem('token');
        // 重定向到登录页
        window.location.href = '/login';
      }
      
      return Promise.reject(new Error(res.message || '请求失败'));
    }
    
    // 返回完整的响应数据
    return res;
  },
  (error) => {
    console.error('请求错误', error);
    
    // 处理网络错误
    let errorMessage = '网络错误，请稍后重试';
    
    if (error.response) {
      const { status } = error.response;
      
      // 根据状态码显示不同的错误信息
      switch (status) {
        case 401:
          errorMessage = '未授权，请重新登录';
          // 清除token
          localStorage.removeItem('token');
          // 重定向到登录页
          window.location.href = '/login';
          break;
        case 403:
          errorMessage = '拒绝访问';
          break;
        case 404:
          errorMessage = '请求的资源不存在';
          break;
        case 500:
          errorMessage = '服务器错误';
          break;
        default:
          errorMessage = `请求错误 (${status})`;
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      errorMessage = '服务器无响应，请稍后重试';
    }
    
    message.error(errorMessage);
    return Promise.reject(error);
  }
);

export default request; 