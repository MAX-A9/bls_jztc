import deviceInfo from '@/utils/device-info.js'

export default {
  data() {
    return {
      // 缓存系统信息，避免重复获取
      _windowInfo: null,
      _deviceInfo: null,
      _appBaseInfo: null,
      _systemSetting: null,
    }
  },
  created() {
    // 在组件创建时获取一次信息，避免重复调用
    try {
      // 使用新的工具类获取设备信息
      this._windowInfo = deviceInfo.getWindowInfo()
      this._deviceInfo = deviceInfo.getDeviceInfo()
      this._appBaseInfo = deviceInfo.getAppBaseInfo()
      this._systemSetting = deviceInfo.getSystemSetting()
    } catch (e) {
      console.error('获取设备信息失败:', e)
      // 使用默认值
      this._windowInfo = {
        windowWidth: 375,
        windowHeight: 667,
        screenWidth: 375,
        screenHeight: 667,
        statusBarHeight: 20,
        safeArea: {
          top: 0,
          right: 375,
          bottom: 667,
          left: 0
        }
      }
      this._deviceInfo = {
        brand: '',
        model: '',
        system: '',
        platform: ''
      }
      this._appBaseInfo = {
        SDKVersion: '',
        appName: '',
        appVersion: '',
        appLanguage: 'zh-Hans'
      }
    }
  },
  computed: {
    // 获取状态栏高度
    statusBarHeight() {
      return this._windowInfo?.statusBarHeight || 0
    },
    
    // 导航栏总高度
    navigationBarHeight() {
      const isIOS = (this._deviceInfo?.system || '').toLowerCase().includes('ios')
      // iOS 和 Android 导航栏高度不同
      const navHeight = isIOS ? 44 : 48
      return this.statusBarHeight + navHeight
    },

    // 是否是小屏幕设备
    isSmallScreen() {
      return (this._windowInfo?.windowWidth || 0) <= 320
    },

    // 是否是大屏幕设备
    isLargeScreen() {
      return (this._windowInfo?.windowWidth || 0) >= 768
    },

    // 安全区域
    safeAreaInsets() {
      const safeArea = this._windowInfo?.safeArea || {}
      return {
        top: safeArea.top || 0,
        bottom: safeArea.bottom || 0,
        left: safeArea.left || 0,
        right: safeArea.right || 0
      }
    },

    // 屏幕尺寸
    screenSize() {
      return {
        width: this._windowInfo?.windowWidth || 0,
        height: this._windowInfo?.windowHeight || 0
      }
    },

    // 统一布局尺寸
    layoutSize() {
      return {
        // 底部结算栏高度(不含安全区域)
        settlementHeight: 120,  // 固定120rpx
        // 内容区域底部间距
        contentBottom: 160,     // 固定160rpx
        // 导航栏高度
        navHeight: 44,         // 固定44px
        // 卡片间距
        cardGap: 20,           // 固定20rpx
        // 内容区域水平内边距
        contentPadding: 20     // 固定20rpx
      }
    }
  },

  methods: {
    // 获取元素布局信息
    async getElementRect(selector) {
      return new Promise((resolve) => {
        uni.createSelectorQuery()
          .select(selector)
          .boundingClientRect(data => {
            resolve(data)
          })
          .exec()
      })
    },

    // rpx 转 px
    rpxToPx(rpx) {
      const screenWidth = this._windowInfo?.windowWidth || 375
      return (rpx / 750) * screenWidth
    },

    // px 转 rpx 
    pxToRpx(px) {
      const screenWidth = this._windowInfo?.windowWidth || 375
      return (px * 750) / screenWidth
    }
  }
} 