import { create } from 'zustand';

// 用户信息接口
interface UserInfo {
  id: number;
  username: string;
  realName: string;
  avatar?: string;
  role: string;
  permissions: string[];
}

// 用户状态接口
interface UserState {
  userInfo: UserInfo | null;
  token: string | null;
  isLoggedIn: boolean;
  setUserInfo: (userInfo: UserInfo | null) => void;
  setToken: (token: string | null) => void;
  logout: () => void;
}

// 创建用户状态管理
const useUserStore = create<UserState>((set) => ({
  // 初始状态
  userInfo: null,
  token: localStorage.getItem('token'),
  isLoggedIn: !!localStorage.getItem('token'),
  
  // 设置用户信息
  setUserInfo: (userInfo) => set({ userInfo }),
  
  // 设置token
  setToken: (token) => {
    if (token) {
      localStorage.setItem('token', token);
      set({ token, isLoggedIn: true });
    } else {
      localStorage.removeItem('token');
      set({ token: null, isLoggedIn: false });
    }
  },
  
  // 退出登录
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    set({ userInfo: null, token: null, isLoggedIn: false });
  },
}));

export default useUserStore; 