"use strict";
const common_vendor = require("../../common/vendor.js");
const utils_storage = require("../../utils/storage.js");
const utils_auth = require("../../utils/auth.js");
common_vendor.index.$on("showLoginModal", () => {
  const store = getApp().$vm.$store;
  if (store) {
    store.dispatch("user/showLoginModal");
  }
});
const state = {
  // 用户信息
  userInfo: utils_storage.getUserInfo() || {},
  // token
  token: utils_storage.getToken() || "",
  // 是否已登录
  isLogin: utils_auth.isLoggedIn(),
  // 登录加载状态
  loginLoading: false,
  // 是否显示登录窗口
  showLoginModal: false
};
const mutations = {
  // 设置用户信息
  SET_USER_INFO(state2, userInfo) {
    state2.userInfo = userInfo;
    utils_storage.setUserInfo(userInfo);
  },
  // 设置token
  SET_TOKEN(state2, token) {
    state2.token = token;
    utils_storage.setToken(token);
    state2.isLogin = utils_auth.isLoggedIn();
  },
  // 清除登录状态
  CLEAR_LOGIN_STATE(state2) {
    state2.userInfo = {};
    state2.token = "";
    state2.isLogin = false;
    utils_storage.clearUserLoginState();
  },
  // 设置登录加载状态
  SET_LOGIN_LOADING(state2, status) {
    state2.loginLoading = status;
  },
  // 设置登录窗口显示状态
  SET_LOGIN_MODAL(state2, status) {
    state2.showLoginModal = status;
  }
};
const actions = {
  // 静默登录
  async silentLogin({ commit }) {
    commit("SET_LOGIN_LOADING", true);
    try {
      const result = await utils_auth.silentLogin();
      if (result) {
        if (result.token) {
          commit("SET_TOKEN", result.token);
        }
        if (result.userInfo) {
          commit("SET_USER_INFO", result.userInfo);
        } else if (utils_auth.isLoggedIn()) {
          try {
            const userInfo = await utils_auth.fetchAndSaveUserInfo();
            if (userInfo) {
              commit("SET_USER_INFO", userInfo);
            }
          } catch (infoError) {
          }
        }
      }
      commit("SET_LOGIN_LOADING", false);
      return result;
    } catch (error) {
      commit("SET_LOGIN_LOADING", false);
      return Promise.reject(error);
    }
  },
  // 完整登录流程
  async login({ commit }) {
    commit("SET_LOGIN_LOADING", true);
    try {
      const result = await utils_auth.fullLogin();
      if (result) {
        if (result.token) {
          commit("SET_TOKEN", result.token);
        }
        if (result.userInfo) {
          commit("SET_USER_INFO", result.userInfo);
        }
        commit("SET_LOGIN_MODAL", false);
      }
      commit("SET_LOGIN_LOADING", false);
      return result;
    } catch (error) {
      commit("SET_LOGIN_LOADING", false);
      return Promise.reject(error);
    }
  },
  // 获取用户信息
  async getUserInfo({ commit, state: state2 }) {
    try {
      const userInfo = await utils_auth.fetchAndSaveUserInfo();
      if (userInfo) {
        if (userInfo.id && !userInfo.clientId) {
          userInfo.clientId = userInfo.id;
        }
        commit("SET_USER_INFO", userInfo);
        if (!state2.isLogin && (userInfo.clientId || userInfo.id)) {
          const token = utils_storage.getToken();
          if (token) {
            commit("SET_TOKEN", token);
          }
        }
      }
      return userInfo;
    } catch (error) {
      if (error.code === 401) {
        commit("CLEAR_LOGIN_STATE");
      }
      return Promise.reject(error);
    }
  },
  // 退出登录
  logout({ commit }) {
    utils_auth.logout();
    commit("CLEAR_LOGIN_STATE");
  },
  // 显示登录窗口
  showLoginModal({ commit }) {
    commit("CLEAR_LOGIN_STATE");
    commit("SET_LOGIN_MODAL", true);
  },
  // 关闭登录窗口
  hideLoginModal({ commit }) {
    commit("SET_LOGIN_MODAL", false);
  }
};
const getters = {
  // 是否已登录
  isLoggedIn: (state2) => state2.isLogin,
  // 用户信息
  userInfo: (state2) => state2.userInfo,
  // token
  token: (state2) => state2.token,
  // 是否显示登录窗口
  showLoginModal: (state2) => state2.showLoginModal
};
const user = {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
};
exports.user = user;
//# sourceMappingURL=../../../.sourcemap/mp-weixin/store/modules/user.js.map
