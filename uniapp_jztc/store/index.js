import { createStore } from 'vuex';
import user from './modules/user.js';
import region from './modules/region.js';

// 创建store实例
const store = createStore({
  modules: {
    user,
    region
  }
});

export default store; 