import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { App as AntdApp, ConfigProvider, Spin } from 'antd';
import zhCN from 'antd/lib/locale/zh_CN';
import MainLayout from './components/layouts/MainLayout';
import Login from './pages/Login';
import License from './pages/License';
import Dashboard from './pages/Dashboard';
import ClientManagement from './pages/ClientManagement';
import FileManagement from './pages/FileManagement';
import UserManagement from './pages/UserManagement';
import OperationLogs from './pages/OperationLogs';
import NotFound from './pages/NotFound';
import HomeCategory from './pages/content/HomeCategory';
import IdleCategory from './pages/content/IdleCategory';
import AllContent from './pages/content/AllContent';
import AllComments from './pages/content/AllComments';
import ContentForm from './pages/content/ContentForm';
import MiniProgram from './pages/city-system/MiniProgram';
import BasicSettings from './pages/city-system/BasicSettings';
import TabSettings from './pages/city-system/TabSettings';
import PageSettings from './pages/city-system/PageSettings';
import PublishSettings from './pages/city-system/PublishSettings';
import RegionSettings from './pages/city-system/RegionSettings';
import Orders from './pages/Orders';
import Mall from './pages/exchange/Mall';
import MallCategory from './pages/exchange/MallCategory';
import ExchangeList from './pages/exchange/ExchangeList';
import CustomerDuration from './pages/exchange/CustomerDuration';
import SmsSettings from './pages/settings/SmsSettings';
import StorageSettings from './pages/settings/StorageSettings';
import WxappConfig from './pages/settings/WxappConfig';
import PaymentSettings from './pages/PaymentSettings';
import { checkLicenseStatus } from './api/license';
import './App.css';

// 路由守卫组件，检查用户是否已登录
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  // 从localStorage中获取token判断用户是否已登录
  const isAuthenticated = localStorage.getItem('token');
  
  if (!isAuthenticated) {
    // 如果未登录，重定向到登录页面
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

// 主应用组件
const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <AntdApp>
        <AppContent />
      </AntdApp>
    </ConfigProvider>
  );
};

// 内容组件，使用AntdApp上下文
const AppContent: React.FC = () => {
  const [isLicensed, setIsLicensed] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(true);
  
  // 使用AntdApp消息API的引用，此处使用是正确的
  const { message } = AntdApp.useApp();

  useEffect(() => {
    // 检查系统是否已授权
    const checkLicense = async () => {
      try {
        // 获取当前域名
        const domain = window.location.hostname;
        const res = await checkLicenseStatus(domain);
        
        // 使用类型断言
        const responseData = res as any;
        if (responseData && responseData.data) {
          setIsLicensed(responseData.data.valid);
          
          // 显示授权状态
          if (responseData.data.valid) {
            message.success('系统已授权');
          }
        } else {
          setIsLicensed(false);
        }
      } catch (error) {
        console.error('授权验证出错:', error);
        setIsLicensed(false);
      } finally {
        setLoading(false);
      }
    };

    checkLicense();
  }, [message]);

  if (loading) {
    return (
      <div style={{ width: '100vw', height: '100vh', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
        <div style={{ textAlign: 'center' }}>
          <Spin size="large" />
          <div style={{ marginTop: 16 }}>正在检查系统授权状态...</div>
        </div>
      </div>
    );
  }

  return (
    <Router>
      <Routes>
        {/* 授权路由，如果未授权则显示License页面 */}
        {!isLicensed && <Route path="*" element={<License />} />}
        
        {/* 已授权后才能访问的路由 */}
        {isLicensed && (
          <>
            <Route path="/login" element={<Login />} />
            
            {/* 受保护的路由，使用MainLayout布局 */}
            <Route 
              path="/" 
              element={
                <ProtectedRoute>
                  <MainLayout />
                </ProtectedRoute>
              }
            >
              {/* 默认重定向到仪表盘 */}
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="dashboard" element={<Dashboard />} />
              
              {/* 内容管理相关路由 */}
              <Route path="content">
                <Route index element={<Navigate to="/content/all-content" replace />} />
                <Route path="all-content" element={<AllContent />} />
                <Route path="all-comments" element={<AllComments />} />
                <Route path="add-content" element={<ContentForm />} />
                <Route path="edit-content/:id" element={<ContentForm />} />
                <Route path="home-category" element={<HomeCategory />} />
                <Route path="idle-category" element={<IdleCategory />} />
              </Route>
              
              <Route path="clients" element={<ClientManagement />} />
              <Route path="orders" element={<Orders />} />
              
              {/* 兑换管理相关路由 */}
              <Route path="exchange">
                <Route index element={<Navigate to="/exchange/list" replace />} />
                <Route path="list" element={<ExchangeList />} />
                <Route path="mall" element={<Mall />} />
                <Route path="mall-category" element={<MallCategory />} />
                <Route path="customer-duration" element={<CustomerDuration />} />
              </Route>
              
              <Route path="files" element={<FileManagement />} />
              
              {/* 权限管理相关路由 */}
              <Route path="auth">
                <Route index element={<Navigate to="/auth/user-management" replace />} />
                <Route path="user-management" element={<UserManagement />} />
                <Route path="operation-logs" element={<OperationLogs />} />
              </Route>
              
              {/* 系统设置相关路由 */}
              <Route path="settings">
                <Route index element={<Navigate to="/settings/sms" replace />} />
                <Route path="sms" element={<SmsSettings />} />
                <Route path="storage" element={<StorageSettings />} />
                <Route path="payment" element={<PaymentSettings />} />
                <Route path="wxapp" element={<WxappConfig />} />
              </Route>
              
              {/* 同城系统相关路由 */}
              <Route path="city-system">
                <Route index element={<Navigate to="/city-system/basic-settings" replace />} />
                <Route path="basic-settings" element={<BasicSettings />} />
                <Route path="page-settings" element={<PageSettings />} />
                <Route path="tab-settings" element={<TabSettings />} />
                <Route path="mini-program" element={<MiniProgram />} />
                <Route path="region-settings" element={<RegionSettings />} />
                <Route path="publish-settings" element={<PublishSettings />} />
              </Route>
              
              {/* 404页面 */}
              <Route path="*" element={<NotFound />} />
            </Route>
          </>
        )}
      </Routes>
    </Router>
  );
};

export default App; 