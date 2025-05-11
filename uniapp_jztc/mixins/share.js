/**
 * 分享功能混入
 */
import { getDefaultShareOptions, getHomeShareOptions, getContentShareOptions } from '../utils/share.js';

/**
 * 分享混入
 * 使用方法：
 * 1. 在页面中导入: import shareMixin from '@/mixins/share.js'
 * 2. 在组件中混入: mixins: [shareMixin]
 * 3. 如果需要自定义分享内容，可以覆盖onShareAppMessage方法
 */
export default {
  data() {
    return {
      // 分享相关数据
      shareData: {
        title: '',
        path: '',
        imageUrl: ''
      }
    };
  },
  
  // 页面生命周期钩子
  onLoad() {
    // 初始化分享数据
    this.initShareData();
  },
  
  // 页面方法
  methods: {
    /**
     * 初始化分享数据
     */
    async initShareData() {
      try {
        // 根据页面类型设置默认分享数据
        if (this.isContentPage) {
          // 内容页面
          const contentData = this.contentData || this.detail || {};
          const shareOptions = await getContentShareOptions(contentData);
          this.shareData = { ...shareOptions };
        } else if (this.isHomePage) {
          // 首页
          const shareOptions = await getHomeShareOptions();
          this.shareData = { ...shareOptions };
        } else {
          // 其他页面
          const pages = getCurrentPages();
          const currentPage = pages[pages.length - 1];
          const path = `/${currentPage.route}${this._getPageQuery(currentPage)}`;
          
          const shareOptions = await getDefaultShareOptions(path);
          this.shareData = { ...shareOptions };
        }
      } catch (error) {
        console.error('初始化分享数据失败:', error);
      }
    },
    
    /**
     * 获取页面查询参数
     * @param {Object} page 页面对象
     * @returns {String} 查询参数字符串
     */
    _getPageQuery(page) {
      if (!page.options || Object.keys(page.options).length === 0) {
        return '';
      }
      
      const query = Object.keys(page.options)
        .map(key => `${key}=${encodeURIComponent(page.options[key])}`)
        .join('&');
        
      return `?${query}`;
    }
  },
  
  // 分享给朋友
  onShareAppMessage() {
    return {
      title: this.shareData.title,
      path: this.shareData.path,
      imageUrl: this.shareData.imageUrl
    };
  },
  
  // 分享到朋友圈
  onShareTimeline() {
    return {
      title: this.shareData.title,
      query: this.shareData.path.indexOf('?') > -1 ? this.shareData.path.split('?')[1] : '',
      imageUrl: this.shareData.imageUrl
    };
  }
}; 