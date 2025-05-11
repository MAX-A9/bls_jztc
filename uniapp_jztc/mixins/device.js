import { ref } from 'vue'
import deviceInfo from '@/utils/device-info.js'

export default {
	data() {
		return {
			// 设备相关
			isIOS: false,
			isAndroid: false,
			statusBarHeight: 20,
			navBarHeight: 44,
			// 屏幕信息
			screenWidth: 375,
			screenHeight: 667,
			// 安全区域
			safeAreaInsets: {
				top: 0,
				right: 0,
				bottom: 0,
				left: 0
			}
		}
	},
	created() {
		try {
			// 使用新的设备信息工具获取系统信息
			const windowInfo = deviceInfo.getWindowInfo();
			const appBaseInfo = deviceInfo.getAppBaseInfo();
			const deviceData = deviceInfo.getDeviceInfo();
			
			// 判断系统类型
			this.isIOS = (deviceData.platform === 'ios') || 
				(/ios/i.test(deviceData.system || ''));
			
			this.isAndroid = (deviceData.platform === 'android') || 
				(/android/i.test(deviceData.system || ''));
			
			// 获取状态栏高度
			this.statusBarHeight = windowInfo.statusBarHeight || 20;
			
			// 设置导航栏高度
			this.navBarHeight = this.isIOS ? 44 : 48;
			
			// 获取屏幕信息
			this.screenWidth = windowInfo.screenWidth || 375;
			this.screenHeight = windowInfo.screenHeight || 667;
			
			// 获取安全区域
			const safeAreaInsets = deviceInfo.getSafeAreaInsets();
			if (safeAreaInsets) {
				this.safeAreaInsets = safeAreaInsets;
			}
		} catch (error) {
			console.error('初始化设备信息失败:', error);
			// 使用默认值继续
		}
	},
	computed: {
		// 状态栏+导航栏的总高度
		navigationBarHeight() {
			return this.statusBarHeight + this.navBarHeight
		},
		// 获取安全区域高度
		safeAreaHeight() {
			return this.screenHeight - this.safeAreaInsets.top - this.safeAreaInsets.bottom
		},
		// 获取内容区域高度(去除导航栏和安全区域)
		contentHeight() {
			return this.safeAreaHeight - this.navigationBarHeight
		}
	},
	methods: {
		// rpx 转 px
		rpx2px(rpx) {
			return deviceInfo.rpx2px(rpx);
		},
		// px 转 rpx
		px2rpx(px) {
			return deviceInfo.px2rpx(px);
		},
		// 判断是否为 iPhone X 系列
		isIphoneX() {
			return deviceInfo.isIphoneX();
		},
		// 获取底部安全距离
		getBottomSafeDistance() {
			return this.isIphoneX() ? this.safeAreaInsets.bottom : 0
		},
		// 获取适配后的安全区域
		getSafeAreaInsets() {
			return {
				top: `${this.safeAreaInsets.top}px`,
				right: `${this.safeAreaInsets.right}px`,
				bottom: `${this.safeAreaInsets.bottom}px`,
				left: `${this.safeAreaInsets.left}px`
			}
		}
	}
} 