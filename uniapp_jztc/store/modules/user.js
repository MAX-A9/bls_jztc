/**
 * 用户状态管理模块
 */
import { getUserInfo, setUserInfo, getToken, setToken, clearUserLoginState } from '../../utils/storage.js';
import { isLoggedIn, silentLogin, fullLogin, logout, fetchAndSaveUserInfo } from '../../utils/auth.js';

// 添加监听登录窗口事件
uni.$on('showLoginModal', () => {
  const store = getApp().$vm.$store;
  if (store) {
    // 触发登录窗口显示的action
    store.dispatch('user/showLoginModal');
  }
});

// 状态
export const state = {
  // 用户信息
  userInfo: getUserInfo() || {},
  // token
  token: getToken() || '',
  // 是否已登录
  isLogin: isLoggedIn(),
  // 登录加载状态
  loginLoading: false,
  // 是否显示登录窗口
  showLoginModal: false
};

// 修改状态的同步方法
export const mutations = {
  // 设置用户信息
  SET_USER_INFO(state, userInfo) {
    state.userInfo = userInfo;
    setUserInfo(userInfo);
  },
  // 设置token
  SET_TOKEN(state, token) {
    state.token = token;
    setToken(token);
    state.isLogin = isLoggedIn();
  },
  // 清除登录状态
  CLEAR_LOGIN_STATE(state) {
    state.userInfo = {};
    state.token = '';
    state.isLogin = false;
    clearUserLoginState();
  },
  // 设置登录加载状态
  SET_LOGIN_LOADING(state, status) {
    state.loginLoading = status;
  },
  // 设置登录窗口显示状态
  SET_LOGIN_MODAL(state, status) {
    state.showLoginModal = status;
  }
};

// 包含异步操作的方法
export const actions = {
  // 静默登录
  async silentLogin({ commit }) {
    commit('SET_LOGIN_LOADING', true);
    try {
      const result = await silentLogin();
      
      if (result) {
        // 确保先设置token，再设置用户信息，顺序很重要
        if (result.token) {
          commit('SET_TOKEN', result.token);
        }
        
        // 如果有用户信息，直接使用
        if (result.userInfo) {
          commit('SET_USER_INFO', result.userInfo);
        } 
        // 否则尝试获取用户信息
        else if (isLoggedIn()) {
          try {
            const userInfo = await fetchAndSaveUserInfo();
            if (userInfo) {
              commit('SET_USER_INFO', userInfo);
            }
          } catch (infoError) {
            // 获取用户信息失败
          }
        }
      }
      
      commit('SET_LOGIN_LOADING', false);
      return result;
    } catch (error) {
      commit('SET_LOGIN_LOADING', false);
      return Promise.reject(error);
    }
  },
  
  // 完整登录流程
  async login({ commit }) {
    commit('SET_LOGIN_LOADING', true);
    try {
      const result = await fullLogin();
      if (result) {
        if (result.token) {
          commit('SET_TOKEN', result.token);
        }
        if (result.userInfo) {
          commit('SET_USER_INFO', result.userInfo);
        }
        // 登录成功后关闭登录窗口
        commit('SET_LOGIN_MODAL', false);
      }
      commit('SET_LOGIN_LOADING', false);
      return result;
    } catch (error) {
      commit('SET_LOGIN_LOADING', false);
      return Promise.reject(error);
    }
  },
  
  // 获取用户信息
  async getUserInfo({ commit, state }) {
    // 始终获取最新的用户信息，确保数据准确
    try {
      const userInfo = await fetchAndSaveUserInfo();
      if (userInfo) {
        // 确保兼容性：将id字段复制到clientId
        if (userInfo.id && !userInfo.clientId) {
          userInfo.clientId = userInfo.id;
        }
        
        commit('SET_USER_INFO', userInfo);
        
        // 确保登录状态一致
        if (!state.isLogin && (userInfo.clientId || userInfo.id)) {
          // 有用户ID但登录状态为false时，更新登录状态
          const token = getToken();
          if (token) {
            commit('SET_TOKEN', token);
          }
        }
      }
      return userInfo;
    } catch (error) {
      // 如果是401错误，清除登录状态
      if (error.code === 401) {
        commit('CLEAR_LOGIN_STATE');
      }
      
      return Promise.reject(error);
    }
  },
  
  // 退出登录
  logout({ commit }) {
    logout();
    commit('CLEAR_LOGIN_STATE');
  },
  
  // 显示登录窗口
  showLoginModal({ commit }) {
    commit('CLEAR_LOGIN_STATE'); // 先清除登录状态
    commit('SET_LOGIN_MODAL', true); // 再显示登录窗口
  },
  
  // 关闭登录窗口
  hideLoginModal({ commit }) {
    commit('SET_LOGIN_MODAL', false);
  }
};

// 获取状态的计算属性
export const getters = {
  // 是否已登录
  isLoggedIn: state => state.isLogin,
  // 用户信息
  userInfo: state => state.userInfo,
  // token
  token: state => state.token,
  // 是否显示登录窗口
  showLoginModal: state => state.showLoginModal
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}; 