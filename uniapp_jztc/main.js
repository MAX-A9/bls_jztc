import App from './App'

// #ifndef VUE3
import Vue from 'vue'
import './uni.promisify.adaptor'
Vue.config.productionTip = false
App.mpType = 'app'

// 添加全局错误处理
Vue.config.errorHandler = function(err, vm, info) {
  console.error('Vue错误:', err);
  console.error('错误信息:', info);
};

// 添加Promise未捕获异常处理
uni.addInterceptor({
  returnValue(res) {
    if (!(res instanceof Promise)) {
      return res;
    }
    return new Promise((resolve, reject) => {
      res.then(res => {
        resolve(res);
      }).catch(error => {
        console.error('Promise未捕获异常:', error);
        reject(error);
      });
    });
  }
});

const app = new Vue({
  ...App
})
app.$mount()
// #endif

// #ifdef VUE3
import { createSSRApp } from 'vue'
import store from './store'

export function createApp() {
  const app = createSSRApp(App)
  
  // 挂载Vuex
  app.use(store)
  
  // 添加全局错误处理
  app.config.errorHandler = function(err, instance, info) {
    console.error('Vue错误:', err);
    console.error('错误信息:', info);
  };
  
  return {
    app
  }
}
// #endif