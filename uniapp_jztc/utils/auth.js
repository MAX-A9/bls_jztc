/**
 * 登录授权相关方法
 */
import { user } from '../apis/index.js';
import { setUserInfo, setToken, getUserInfo, getToken, clearUserLoginState } from './storage.js';
import { API_CODE } from './constants.js';

/**
 * 检查是否已登录
 * @returns {Boolean} 是否已登录
 */
export function isLoggedIn() {
  return !!getToken();
}

// 重新导出getUserInfo函数，以便其他文件可以从auth.js中导入
export { getUserInfo };

/**
 * 获取客户信息并保存
 * @returns {Promise} 返回Promise对象
 */
export async function fetchAndSaveUserInfo() {
  try {
    const result = await user.getClientInfo();
    
    if (result && result.code === API_CODE.SUCCESS && result.data) {
      // 后端返回的是id字段，而非clientId字段，需要转换
      const userData = result.data;
      
      // 确保兼容性：如果后端返回的是id字段，将其复制到clientId字段
      if (userData.id && !userData.clientId) {
        userData.clientId = userData.id;
      }
      
      // 保存用户信息到本地
      setUserInfo(userData);
      return userData;
    } else {
      const error = new Error(result.message || '获取客户信息失败');
      throw error;
    }
  } catch (error) {
    return Promise.reject(error);
  }
}

/**
 * 静默登录
 * @returns {Promise} 返回Promise对象
 */
export async function silentLogin() {
  try {
    // 1. 获取微信登录凭证(code)
    const code = await user.getWxLoginCode();
    
    // 2. 使用code调用后端登录接口
    const loginData = { code };
    
    const loginResult = await user.wxappLogin(loginData);
    
    // 3. 处理后端返回结果，适配实际接口结构
    if (loginResult && loginResult.code === API_CODE.SUCCESS && loginResult.data && loginResult.data.token) {
      // 保存token
      setToken(loginResult.data.token);
      
      // 4. 获取并保存用户信息
      const userInfo = await fetchAndSaveUserInfo();
      
      return {
        ...loginResult.data,
        userInfo
      };
    } else {
      const error = new Error(loginResult.message || '登录失败，返回数据不完整');
      throw error;
    }
  } catch (error) {
    return Promise.reject(error);
  }
}

/**
 * 完整登录流程（包含获取用户信息）
 * @returns {Promise} 返回Promise对象
 */
export async function fullLogin() {
  try {
    // 1. 获取微信登录凭证(code)
    const code = await user.getWxLoginCode();
    
    // 2. 获取用户信息
    const wxUserInfo = await user.getWxUserInfo();
    
    // 3. 使用code和用户信息调用后端登录接口
    const loginData = {
      code,
      userInfo: wxUserInfo
    };
    
    const loginResult = await user.wxappLogin(loginData);
    
    // 4. 处理后端返回结果，适配实际接口结构
    if (loginResult && loginResult.code === API_CODE.SUCCESS && loginResult.data && loginResult.data.token) {
      // 保存token
      setToken(loginResult.data.token);
      
      // 5. 获取并保存用户信息
      const userInfo = await fetchAndSaveUserInfo();
      
      // 返回完整的结果
      return {
        ...loginResult.data,
        userInfo
      };
    } else {
      const error = new Error(loginResult.message || '登录失败，返回数据不完整');
      throw error;
    }
  } catch (error) {
    return Promise.reject(error);
  }
}

/**
 * 退出登录
 */
export function logout() {
  clearUserLoginState();
}

/**
 * 检查并自动登录（无感登录）
 * @returns {Promise} 返回Promise对象
 */
export async function checkAndAutoLogin() {
  // 如果已登录，先尝试使用已有token获取用户信息
  if (isLoggedIn()) {
    const cachedUserInfo = getUserInfo();
    
    // 缓存中有完整用户信息，直接返回
    // 确保兼容id和clientId字段
    const hasUserId = cachedUserInfo && (cachedUserInfo.clientId || cachedUserInfo.id);
    const hasBasicInfo = cachedUserInfo && cachedUserInfo.realName;
    
    if (hasUserId && hasBasicInfo) {
      return {
        token: getToken(),
        userInfo: cachedUserInfo
      };
    }
    
    // 缓存中没有完整用户信息，但有token，尝试获取用户信息
    try {
      const userInfo = await fetchAndSaveUserInfo();
      return {
        token: getToken(),
        userInfo: userInfo
      };
    } catch (error) {
      // 获取用户信息失败，可能是token已失效，清除登录状态
      clearUserLoginState();
    }
  }
  
  // 没有token或token已失效，尝试静默登录
  try {
    const loginResult = await silentLogin();
    
    // 确保登录结果包含必要的信息
    if (!loginResult.userInfo && loginResult.token) {
      try {
        const userInfo = await fetchAndSaveUserInfo();
        return {
          ...loginResult,
          userInfo
        };
      } catch (infoError) {
        return loginResult;
      }
    }
    
    return loginResult;
  } catch (loginError) {
    return Promise.reject(loginError);
  }
} 