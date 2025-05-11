<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<view class="nav-bar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="nav-content" :style="{ height: navBarHeight + 'px' }">
				<view class="back-btn" @tap="handleBack">
					<uni-icons type="left" size="20" color="#333333"></uni-icons>
				</view>
				<text class="nav-title">师傅入驻</text>
			</view>
		</view>
		
		<!-- 表单内容 -->
		<scroll-view 
			class="content-scroll" 
			scroll-y
			:style="{ top: `calc(${statusBarHeight}px + ${navBarHeight}px)` }"
		>
			<view class="form-wrapper">
				<view class="form-header">
					<image class="header-image" src="/static/demo/0.png" mode="aspectFill"></image>
					<view class="header-text">
						<text class="header-title">成为橘子师傅</text>
						<text class="header-desc">加入我们，获得更多工作机会</text>
					</view>
				</view>
				
				<view class="form-section">
					<view class="section-title">基本信息</view>
					
					<view class="form-item">
						<text class="form-label">姓名</text>
						<input class="form-input" type="text" v-model="formData.name" placeholder="请输入真实姓名" />
					</view>
					
					<view class="form-item">
						<text class="form-label">手机号</text>
						<input class="form-input" type="number" v-model="formData.phone" placeholder="请输入手机号码" maxlength="11" />
					</view>
					
					<view class="form-item">
						<text class="form-label">身份证号</text>
						<input class="form-input" type="idcard" v-model="formData.idCard" placeholder="请输入身份证号码" maxlength="18" />
					</view>
					
					<view class="form-item">
						<text class="form-label">所在城市</text>
						<picker 
							class="form-picker" 
							mode="region" 
							@change="handleRegionChange"
							:value="formData.region"
						>
							<view class="picker-text">{{ formData.regionText || '请选择所在城市' }}</view>
						</picker>
						<uni-icons class="picker-arrow" type="right" size="16" color="#999999"></uni-icons>
					</view>
				</view>
				
				<view class="form-section">
					<view class="section-title">专业技能</view>
					
					<view class="skill-list">
						<view 
							class="skill-item" 
							v-for="(item, index) in skillList" 
							:key="index"
							:class="{ active: formData.skills.includes(item.value) }"
							@tap="toggleSkill(item.value)"
						>
							<text class="skill-text">{{ item.label }}</text>
						</view>
					</view>
				</view>
				
				<view class="form-section">
					<view class="section-title">工作经验</view>
					
					<view class="form-item">
						<text class="form-label">工作年限</text>
						<picker 
							class="form-picker" 
							mode="selector" 
							:range="experienceOptions"
							@change="handleExperienceChange"
							:value="formData.experienceIndex"
						>
							<view class="picker-text">{{ formData.experience || '请选择工作年限' }}</view>
						</picker>
						<uni-icons class="picker-arrow" type="right" size="16" color="#999999"></uni-icons>
					</view>
					
					<view class="form-item">
						<text class="form-label">个人简介</text>
						<textarea 
							class="form-textarea" 
							v-model="formData.introduction" 
							placeholder="请简要介绍您的工作经历和专业技能"
							maxlength="200"
						></textarea>
						<text class="textarea-counter">{{ formData.introduction.length }}/200</text>
					</view>
				</view>
				
				<view class="form-section">
					<view class="section-title">上传资质</view>
					
					<view class="upload-list">
						<view class="upload-item" @tap="uploadImage('idCardFront')">
							<image 
								class="upload-image" 
								:src="formData.idCardFront || '/static/upload-id-front.png'" 
								mode="aspectFill"
							></image>
							<text class="upload-text">身份证正面</text>
						</view>
						
						<view class="upload-item" @tap="uploadImage('idCardBack')">
							<image 
								class="upload-image" 
								:src="formData.idCardBack || '/static/upload-id-back.png'" 
								mode="aspectFill"
							></image>
							<text class="upload-text">身份证反面</text>
						</view>
						
						<view class="upload-item" @tap="uploadImage('qualification')">
							<image 
								class="upload-image" 
								:src="formData.qualification || '/static/upload-qualification.png'" 
								mode="aspectFill"
							></image>
							<text class="upload-text">资格证书</text>
						</view>
					</view>
				</view>
				
				<view class="agreement-section">
					<view class="agreement-checkbox" @tap="toggleAgreement">
						<uni-icons 
							:type="formData.agreement ? 'checkbox-filled' : 'circle'" 
							size="20" 
							:color="formData.agreement ? '#fc3e2b' : '#999999'"
						></uni-icons>
					</view>
					<text class="agreement-text">我已阅读并同意</text>
					<text class="agreement-link" @tap.stop="showAgreement">《师傅入驻协议》</text>
				</view>
				
				<button class="submit-btn" :disabled="!formData.agreement" @tap="handleSubmit">提交申请</button>
			</view>
		</scroll-view>
	</view>
</template>

<script>
	import deviceInfo from '@/utils/device-info.js'

	export default {
		data() {
			return {
				statusBarHeight: 0,
				navBarHeight: 44,
				formData: {
					name: '',
					phone: '',
					idCard: '',
					region: ['', '', ''],
					regionText: '',
					skills: [],
					experienceIndex: 0,
					experience: '',
					introduction: '',
					idCardFront: '',
					idCardBack: '',
					qualification: '',
					agreement: false
				},
				skillList: [
					{ label: '家电维修', value: 'appliance' },
					{ label: '水电维修', value: 'plumbing' },
					{ label: '空调安装', value: 'ac' },
					{ label: '家具安装', value: 'furniture' },
					{ label: '管道疏通', value: 'pipe' },
					{ label: '电器安装', value: 'electronics' },
					{ label: '防水补漏', value: 'waterproof' },
					{ label: '开锁换锁', value: 'lock' }
				],
				experienceOptions: ['1年以下', '1-3年', '3-5年', '5-10年', '10年以上']
			}
		},
		onLoad() {
			// 获取状态栏高度
			this.statusBarHeight = deviceInfo.getStatusBarHeight();
		},
		methods: {
			handleBack() {
				uni.navigateBack()
			},
			handleRegionChange(e) {
				this.formData.region = e.detail.value
				this.formData.regionText = e.detail.value.join(' ')
			},
			toggleSkill(value) {
				const index = this.formData.skills.indexOf(value)
				if (index === -1) {
					this.formData.skills.push(value)
				} else {
					this.formData.skills.splice(index, 1)
				}
			},
			handleExperienceChange(e) {
				this.formData.experienceIndex = e.detail.value
				this.formData.experience = this.experienceOptions[e.detail.value]
			},
			uploadImage(field) {
				uni.chooseImage({
					count: 1,
					sizeType: ['compressed'],
					sourceType: ['album', 'camera'],
					success: (res) => {
						this.formData[field] = res.tempFilePaths[0]
					}
				})
			},
			toggleAgreement() {
				this.formData.agreement = !this.formData.agreement
			},
			showAgreement() {
				uni.showModal({
					title: '师傅入驻协议',
					content: '这里是师傅入驻协议的详细内容...',
					showCancel: false
				})
			},
			handleSubmit() {
				// 表单验证
				if (!this.formData.name) {
					return uni.showToast({ title: '请输入姓名', icon: 'none' })
				}
				if (!this.formData.phone || !/^1\d{10}$/.test(this.formData.phone)) {
					return uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
				}
				if (!this.formData.idCard) {
					return uni.showToast({ title: '请输入身份证号', icon: 'none' })
				}
				if (!this.formData.regionText) {
					return uni.showToast({ title: '请选择所在城市', icon: 'none' })
				}
				if (this.formData.skills.length === 0) {
					return uni.showToast({ title: '请选择至少一项专业技能', icon: 'none' })
				}
				if (!this.formData.experience) {
					return uni.showToast({ title: '请选择工作年限', icon: 'none' })
				}
				if (!this.formData.idCardFront) {
					return uni.showToast({ title: '请上传身份证正面照片', icon: 'none' })
				}
				if (!this.formData.idCardBack) {
					return uni.showToast({ title: '请上传身份证反面照片', icon: 'none' })
				}
				
				// 提交表单
				uni.showLoading({ title: '提交中...' })
				
				setTimeout(() => {
					uni.hideLoading()
					uni.showModal({
						title: '提交成功',
						content: '您的入驻申请已提交，我们将在1-3个工作日内审核，请留意短信通知',
						showCancel: false,
						success: () => {
							uni.navigateBack()
						}
					})
				}, 1500)
			}
		}
	}
</script>

<style>
	.container {
		min-height: 100vh;
		background-color: #f5f5f5;
		position: relative;
	}
	
	.nav-bar {
		background-color: #ffffff;
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 99;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
	}
	
	.nav-content {
		display: flex;
		align-items: center;
		position: relative;
	}
	
	.back-btn {
		width: 88rpx;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-left: 7px;
	}
	
	.nav-title {
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
		font-size: 34rpx;
		font-weight: 500;
		color: #333333;
		max-width: 350rpx;
		text-align: center;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.content-scroll {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: #f5f5f5;
		z-index: 1;
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
	}
	
	.form-wrapper {
		padding: 20rpx 30rpx;
		padding-bottom: calc(env(safe-area-inset-bottom) + 40rpx);
	}
	
	.form-header {
		display: flex;
		align-items: center;
		background-color: #ffffff;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
	}
	
	.header-image {
		width: 120rpx;
		height: 120rpx;
		border-radius: 60rpx;
		margin-right: 30rpx;
	}
	
	.header-text {
		flex: 1;
	}
	
	.header-title {
		font-size: 36rpx;
		font-weight: bold;
		color: #333333;
		margin-bottom: 10rpx;
	}
	
	.header-desc {
		font-size: 28rpx;
		color: #666666;
	}
	
	.form-section {
		background-color: #ffffff;
		border-radius: 16rpx;
		padding: 30rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
	}
	
	.section-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333333;
		margin-bottom: 30rpx;
		position: relative;
		padding-left: 20rpx;
	}
	
	.section-title::before {
		content: '';
		position: absolute;
		left: 0;
		top: 6rpx;
		width: 8rpx;
		height: 32rpx;
		background-color: #fc3e2b;
		border-radius: 4rpx;
	}
	
	.form-item {
		margin-bottom: 30rpx;
		position: relative;
	}
	
	.form-item:last-child {
		margin-bottom: 0;
	}
	
	.form-label {
		font-size: 28rpx;
		color: #333333;
		margin-bottom: 16rpx;
		display: block;
	}
	
	.form-input {
		width: 100%;
		height: 88rpx;
		background-color: #f8f8f8;
		border-radius: 12rpx;
		padding: 0 24rpx;
		font-size: 28rpx;
		color: #333333;
		box-sizing: border-box;
	}
	
	.form-picker {
		width: 100%;
		height: 88rpx;
		background-color: #f8f8f8;
		border-radius: 12rpx;
		padding: 0 24rpx;
		font-size: 28rpx;
		color: #333333;
		display: flex;
		align-items: center;
		box-sizing: border-box;
	}
	
	.picker-text {
		flex: 1;
		color: #333333;
	}
	
	.picker-arrow {
		position: absolute;
		right: 24rpx;
		top: 50%;
		transform: translateY(-50%);
	}
	
	.form-textarea {
		width: 100%;
		height: 200rpx;
		background-color: #f8f8f8;
		border-radius: 12rpx;
		padding: 24rpx;
		font-size: 28rpx;
		color: #333333;
		box-sizing: border-box;
	}
	
	.textarea-counter {
		position: absolute;
		right: 24rpx;
		bottom: 24rpx;
		font-size: 24rpx;
		color: #999999;
	}
	
	.skill-list {
		display: flex;
		flex-wrap: wrap;
		margin: 0 -10rpx;
	}
	
	.skill-item {
		width: calc(33.33% - 20rpx);
		height: 80rpx;
		background-color: #f8f8f8;
		border-radius: 12rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0 10rpx 20rpx;
	}
	
	.skill-item.active {
		background-color: #fff2f0;
		border: 1rpx solid #fc3e2b;
	}
	
	.skill-text {
		font-size: 28rpx;
		color: #333333;
	}
	
	.skill-item.active .skill-text {
		color: #fc3e2b;
	}
	
	.upload-list {
		display: flex;
		justify-content: space-between;
	}
	
	.upload-item {
		width: 200rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
	.upload-image {
		width: 200rpx;
		height: 140rpx;
		background-color: #f8f8f8;
		border-radius: 12rpx;
		margin-bottom: 16rpx;
	}
	
	.upload-text {
		font-size: 26rpx;
		color: #666666;
	}
	
	.agreement-section {
		display: flex;
		align-items: center;
		margin: 40rpx 0;
	}
	
	.agreement-checkbox {
		margin-right: 12rpx;
	}
	
	.agreement-text {
		font-size: 26rpx;
		color: #666666;
	}
	
	.agreement-link {
		font-size: 26rpx;
		color: #fc3e2b;
	}
	
	.submit-btn {
		width: 100%;
		height: 88rpx;
		background-color: #fc3e2b;
		color: #ffffff;
		font-size: 32rpx;
		font-weight: 500;
		border-radius: 44rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.submit-btn[disabled] {
		background-color: #cccccc;
		color: #ffffff;
	}
</style> 