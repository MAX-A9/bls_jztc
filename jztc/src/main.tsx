import React from 'react';
import ReactDOM from 'react-dom/client';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import App from './App';
import './index.css';
import { ENV } from './utils/env';

// 在开发环境下输出环境变量信息
if (ENV.IS_DEV) {
  console.log('===== 环境变量信息 =====');
  console.log('应用名称：', ENV.APP_NAME);
  console.log('当前环境：', ENV.IS_DEV ? '开发环境' : '生产环境');
  console.log('API服务器：', ENV.API_SERVER);
  console.log('=====================');
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ConfigProvider locale={zhCN} theme={{ 
      token: { 
        colorPrimary: '#1890ff',
      },
    }}>
      <App />
    </ConfigProvider>
  </React.StrictMode>,
); 