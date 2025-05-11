import React, { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  UserOutlined,
  FileOutlined,
  SettingOutlined,
  LogoutOutlined,
  LockOutlined,
  GlobalOutlined,
  FileTextOutlined,
  ShoppingCartOutlined,
  GiftOutlined,
} from '@ant-design/icons';
import {
  Layout,
  Menu,
  Button,
  theme,
  Dropdown,
  Avatar,
  message,
} from 'antd';
import type { MenuProps } from 'antd';
import './MainLayout.css';
import { getUserInfo, logout, UserInfoResponse } from '../../api/auth';

const { Header, Sider, Content } = Layout;

// 定义菜单项类型
interface MenuItem {
  key: string;
  icon: React.ReactNode | null;
  label: string;
  path: string;
  children?: MenuItem[];
  disabled?: boolean;
  className?: string;
}

// 菜单项配置
const menuItems: MenuItem[] = [
  {
    key: 'dashboard',
    icon: <DashboardOutlined />,
    label: '仪表盘',
    path: '/dashboard',
  },
  {
    key: 'content',
    icon: <FileTextOutlined />,
    label: '内容管理',
    path: '/content',
    children: [
      {
        key: 'all-content',
        icon: null,
        label: '所有内容',
        path: '/content/all-content',
      },
      {
        key: 'all-comments',
        icon: null,
        label: '所有评论',
        path: '/content/all-comments',
      },
      {
        key: 'home-category',
        icon: null,
        label: '首页分类',
        path: '/content/home-category',
      },
      {
        key: 'idle-category',
        icon: null,
        label: '闲置分类',
        path: '/content/idle-category',
      },
    ],
  },
  {
    key: 'clients',
    icon: <UserOutlined />,
    label: '客户管理',
    path: '/clients',
  },
  {
    key: 'orders',
    icon: <ShoppingCartOutlined />,
    label: '订单管理',
    path: '/orders',
  },
  {
    key: 'exchange',
    icon: <GiftOutlined />,
    label: '兑换管理',
    path: '/exchange',
    children: [
      {
        key: 'exchange-list',
        icon: null,
        label: '兑换列表',
        path: '/exchange/list',
      },
      {
        key: 'mall',
        icon: null,
        label: '商城管理',
        path: '/exchange/mall',
      },
      {
        key: 'mall-category',
        icon: null,
        label: '商城分类',
        path: '/exchange/mall-category',
      },
      {
        key: 'customer-duration',
        icon: null,
        label: '客户时长',
        path: '/exchange/customer-duration',
      },
    ],
  },
  {
    key: 'files',
    icon: <FileOutlined />,
    label: '文件管理',
    path: '/files',
  },
  {
    key: 'auth',
    icon: <LockOutlined />,
    label: '权限管理',
    path: '/auth',
    children: [
      {
        key: 'user-management',
        icon: null,
        label: '用户管理',
        path: '/auth/user-management',
      },
      {
        key: 'operation-logs',
        icon: null,
        label: '操作日志',
        path: '/auth/operation-logs',
      },
    ],
  },
  {
    key: 'settings',
    icon: <SettingOutlined />,
    label: '常规管理',
    path: '/settings',
    children: [
      {
        key: 'sms-settings',
        icon: null,
        label: '短信设置',
        path: '/settings/sms',
      },
      {
        key: 'storage-settings',
        icon: null,
        label: '存储设置',
        path: '/settings/storage',
      },
      {
        key: 'payment-settings',
        icon: null,
        label: '支付设置',
        path: '/settings/payment',
      },
      {
        key: 'wxapp-settings',
        icon: null,
        label: '微信配置',
        path: '/settings/wxapp',
      },
    ],
  },
  {
    key: 'city-system',
    icon: <GlobalOutlined />,
    label: '同城系统',
    path: '/city-system',
    children: [
      {
        key: 'basic-settings',
        icon: null,
        label: '首页布局',
        path: '/city-system/basic-settings',
      },
      {
        key: 'page-settings',
        icon: null,
        label: '内页设置',
        path: '/city-system/page-settings',
      },
      {
        key: 'tab-settings',
        icon: null,
        label: '底部设置',
        path: '/city-system/tab-settings',
      },
      {
        key: 'mini-program',
        icon: null,
        label: '基础设置',
        path: '/city-system/mini-program',
      },
      {
        key: 'region-settings',
        icon: null,
        label: '地区设置',
        path: '/city-system/region-settings',
      },
      {
        key: 'publish-settings',
        icon: null,
        label: '套餐设置',
        path: '/city-system/publish-settings',
      },
    ],
  },
  // 可以在此添加更多菜单项
];

const MainLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [userInfo, setUserInfo] = useState<UserInfoResponse['data'] | null>(null);
  const navigate = useNavigate();
  const location = useLocation();
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  // 获取用户信息
  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const response = await getUserInfo();
        // 确保userData是正确的类型
        if (response.data) {
          setUserInfo(response.data);
        } else if (response.id !== undefined) {
          // 如果API直接返回数据而不是嵌套在data中
          setUserInfo({
            id: response.id,
            username: response.username || '',
            nickname: response.nickname || '',
            role: response.role || '',
            status: response.status || 0,
            lastLogin: response.lastLogin || ''
          });
        }
      } catch (error) {
        console.error('获取用户信息失败:', error);
        message.error('获取用户信息失败，请重新登录');
        // 如果获取用户信息失败，可能是token过期，跳转到登录页
        handleLogout();
      }
    };

    fetchUserInfo();
  }, []);

  // 处理菜单点击
  const handleMenuClick = (key: string) => {
    // 查找点击的菜单项
    const findMenuItemRecursive = (items: MenuItem[]): MenuItem | null => {
      for (const item of items) {
        if (item.key === key) return item;
        if (item.children) {
          const found = findMenuItemRecursive(item.children);
          if (found) return found;
        }
      }
      return null;
    };

    const menuItem = findMenuItemRecursive(menuItems);
    if (menuItem && menuItem.path) {
      navigate(menuItem.path);
    }
  };

  // 处理退出登录
  const handleLogout = async () => {
    try {
      message.loading('正在退出登录...', 0.5);
      await logout();
      message.success('已成功退出登录');
      navigate('/login');
    } catch (error) {
      console.error('退出登录失败:', error);
      message.error('退出登录失败，请重试');
      // 即使API调用失败，也尝试重定向到登录页
      navigate('/login');
    }
  };

  // 用户下拉菜单
  const userDropdownItems: MenuProps['items'] = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: () => handleLogout(),
    },
  ];

  // 获取当前选中的菜单项和展开的子菜单
  const getSelectedAndOpenKeys = () => {
    const pathname = location.pathname;
    let selectedKey = '';
    let openKey = '';

    // 递归检查路径匹配
    const checkPath = (items: MenuItem[]): boolean => {
      for (const item of items) {
        if (item.path && pathname.startsWith(item.path)) {
          if (item.children) {
            openKey = item.key;
            // 检查子菜单项
            for (const child of item.children) {
              // 使用严格匹配或前缀匹配
              if ((child.path === pathname) || 
                  (pathname.startsWith(child.path) && 
                   (pathname === child.path || 
                    pathname.substring(child.path.length, child.path.length + 1) === '/'))) {
                selectedKey = child.key;
                return true;
              }
            }
            // 如果没有找到精确匹配的子菜单，但路径以父菜单路径开头，则保持父菜单展开
            if (pathname.startsWith(item.path)) {
              return true;
            }
          } else {
            selectedKey = item.key;
            return true;
          }
        }
        if (item.children && checkPath(item.children)) {
          return true;
        }
      }
      return false;
    };

    checkPath(menuItems);
    return { selectedKeys: selectedKey ? [selectedKey] : [], openKeys: openKey ? [openKey] : [] };
  };

  const { selectedKeys, openKeys } = getSelectedAndOpenKeys();
  const [openedKeys, setOpenedKeys] = useState<string[]>(openKeys);
  
  // 当路由变化时更新选中的菜单项和展开的子菜单
  useEffect(() => {
    const { openKeys: newOpenKeys } = getSelectedAndOpenKeys();
    // 更新展开的子菜单，但不覆盖用户手动展开的子菜单
    setOpenedKeys(prev => {
      const mergedKeys = [...prev];
      newOpenKeys.forEach(key => {
        if (!mergedKeys.includes(key)) {
          mergedKeys.push(key);
        }
      });
      return mergedKeys;
    });
  }, [location.pathname]);

  // 将菜单项转换为Ant Design菜单格式
  const convertMenuItems = (items: MenuItem[]): MenuProps['items'] => {
    return items.map(item => {
      const menuItem: any = {
        key: item.key,
        icon: item.icon,
        label: item.label,
        disabled: item.disabled,
        className: item.className,
      };

      if (item.children && item.children.length > 0) {
        menuItem.children = convertMenuItems(item.children);
      }

      return menuItem;
    });
  };

  return (
    <Layout className={`main-layout ${collapsed ? 'collapsed' : ''}`}>
      <Sider 
        trigger={null} 
        collapsible 
        collapsed={collapsed}
        theme="light"
        className="main-sider"
        width={200}
      >
        <div className="logo">{collapsed ? '后台' : '管理系统'}</div>
        <Menu
          theme="light"
          mode="inline"
          selectedKeys={selectedKeys}
          openKeys={openedKeys}
          onOpenChange={(keys) => setOpenedKeys(keys as string[])}
          onClick={({ key }) => handleMenuClick(key)}
          items={convertMenuItems(menuItems)}
        />
      </Sider>
      <Layout>
        <Header className="main-header" style={{ background: colorBgContainer }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            className="trigger-button"
            title={collapsed ? "展开菜单" : "收起菜单"}
          />
          <div className="header-right">
            <Dropdown menu={{ items: userDropdownItems }} placement="bottomRight">
              <div className="user-dropdown">
                <Avatar 
                  style={{ backgroundColor: '#1890ff' }} 
                  icon={<UserOutlined />}
                />
                {userInfo ? (
                  <span className="username">{userInfo.nickname}</span>
                ) : (
                  <span className="username">加载中...</span>
                )}
              </div>
            </Dropdown>
          </div>
        </Header>
        <Content className="main-content">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout; 