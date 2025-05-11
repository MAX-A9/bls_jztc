import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Form,
  Input,
  Switch,
  Button,
  Row,
  Col,
  Tabs,
  Space,
  InputNumber,
  Alert,
  Image,
  App,
} from 'antd';
import { 
  SaveOutlined, 
  ReloadOutlined, 
  UploadOutlined, 
  PlusOutlined,
} from '@ant-design/icons';
import type { TabsProps } from 'antd';
import './styles.css';
import { getMiniProgramBaseSettings, saveMiniProgramBaseSettings } from '@/api/wxapp';
import { getAdSettings, saveAdSettings } from '@/api/ad';
import FileSelector from '@/components/FileSelector';
import { getRewardSettings, saveRewardSettings, getAgreementSettings, saveAgreementSettings } from '@/api/reward';
import request from '@/utils/request';
import { 
  getShareSettings,
  saveShareSettings,
} from '@/api/settings';

const { Title } = Typography;
const { TextArea } = Input;

// 主MiniProgram组件定义
const MiniProgramContent: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saveLoading, setSaveLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('basic');
  const [fileSelectVisible, setFileSelectVisible] = useState(false);
  const [logoUrl, setLogoUrl] = useState<string>('');
  const [qrcodeSelectVisible, setQrcodeSelectVisible] = useState(false);
  const [qrcodeUrl, setQrcodeUrl] = useState<string>('');
  const [contentImageSelectVisible, setContentImageSelectVisible] = useState(false);
  const [contentImageUrl, setContentImageUrl] = useState<string>('');
  const [homeImageSelectVisible, setHomeImageSelectVisible] = useState(false);
  const [homeImageUrl, setHomeImageUrl] = useState<string>('');
  const [defaultImageSelectVisible, setDefaultImageSelectVisible] = useState(false);
  const [defaultImageUrl, setDefaultImageUrl] = useState<string>('');
  
  // 使用App上下文中的message
  const { message } = App.useApp();

  // 加载现有配置
  useEffect(() => {
    fetchTabData(activeTab);
  }, []);

  // 获取当前选项卡数据
  const fetchTabData = async (tabKey: string) => {
    setLoading(true);
    try {
      // 基本设置数据从API获取
      if (tabKey === 'basic') {
        const res = await getMiniProgramBaseSettings();
        if (res.code === 0) {
          // 只更新基本设置相关字段
          form.setFieldsValue({
            name: res.data.name,
            description: res.data.description,
          });
          
          // 设置Logo URL
          if (res.data.logo) {
            setLogoUrl(res.data.logo);
          }
        } else {
          message.error(res.message || '获取小程序基本设置失败');
        }
      } 
      // 广告设置数据从API获取
      else if (tabKey === 'ads') {
        const res = await getAdSettings();
        if (res.code === 0) {
          form.setFieldsValue({
            enableWxAd: res.data.enableWxAd,
            rewardedVideoAdId: res.data.rewardedVideoAdId
          });
        } else {
          message.error(res.message || '获取广告设置失败');
        }
      }
      else if (tabKey === 'rewards') {
        try {
          // 从API获取奖励设置数据
          const res = await getRewardSettings();
          if (res.code === 0) {
            form.setFieldsValue({
              enableRewards: res.data.enableReward,
              firstTimeMinMinutes: res.data.firstViewMinRewardMin,
              firstTimeMaxValue: res.data.firstViewMaxRewardDay,
              rewardMinMinutes: res.data.singleAdMinRewardMin,
              rewardMaxDays: res.data.singleAdMaxRewardDay,
              dailyRewardLimit: res.data.dailyRewardLimit,
              maxDailyRewardTime: res.data.dailyMaxAccumulatedDay,
              rewardExpireDays: res.data.rewardExpirationDays
            });
          } else {
            message.error(res.message || '获取奖励设置失败');
          }
        } catch (error) {
          console.error('获取奖励设置失败:', error);
          message.error('获取奖励设置失败，请检查网络连接');
        }
      }
      else if (tabKey === 'agreement') {
        try {
          // 从API获取协议设置数据
          const res = await getAgreementSettings();
          if (res.code === 0) {
            form.setFieldsValue({
              privacyPolicy: res.data.privacyPolicy,
              userAgreement: res.data.userAgreement
            });
          } else {
            message.error(res.message || '获取协议设置失败');
          }
        } catch (error) {
          console.error('获取协议设置失败:', error);
          message.error('获取协议设置失败，请检查网络连接');
        }
      }
      else if (tabKey === 'exclusiveManager') {
        try {
          // 从API获取专属管家图片
          const response = await request({
            url: '/butler',
            method: 'get'
          });
          
          console.log('专属管家返回数据:', response);
          
          // 判断response是否直接就是我们需要的数据格式
          if (response && typeof response === 'object') {
            let resData;
            
            // 判断返回数据结构
            if ('data' in response && 'code' in response) {
              // response本身就是最终数据
              resData = response;
            } else if ('data' in response && typeof response.data === 'object' && 'code' in response.data) {
              // response是axios包装的，data才是最终数据
              resData = response.data;
            } else {
              // 无法识别的数据结构
              console.error('无法识别的数据结构:', response);
              message.error('获取专属管家设置失败: 数据格式错误');
              return;
            }
            
            if (resData.code === 0 && resData.data) {
              // 设置专属管家图片URL
              if (resData.data.imageUrl) {
                setQrcodeUrl(resData.data.imageUrl);
                form.setFieldsValue({ qrcode: resData.data.imageUrl });
              }
            } else {
              message.error(resData.message || '获取专属管家设置失败');
            }
          } else {
            message.error('获取专属管家设置失败: 返回数据无效');
          }
        } catch (error) {
          console.error('获取专属管家设置失败:', error);
          message.error('获取专属管家设置失败');
        }
      }
      else if (tabKey === 'shareSettings') {
        try {
          // 从API获取分享设置数据
          const response = await getShareSettings();
          
          console.log('分享设置返回数据:', response);
          
          if (response && response.code === 0 && response.data) {
            // 设置表单数据和状态
            const shareData = response.data;
            
            setContentImageUrl(shareData.content_share_image || '');
            setHomeImageUrl(shareData.home_share_image || '');
            setDefaultImageUrl(shareData.default_share_image || '');
            
            form.setFieldsValue({
              defaultShareTitle: shareData.default_share_text || '',
              defaultShareImage: shareData.default_share_image || '',
              contentShareTitle: shareData.content_share_text || '',
              homeShareTitle: shareData.home_share_text || '',
              contentDefaultImage: shareData.content_share_image || '',
              homeDefaultImage: shareData.home_share_image || ''
            });
          } else {
            message.error(response.message || '获取分享设置失败');
          }
        } catch (error) {
          console.error('获取分享设置失败:', error);
          message.error('获取分享设置失败');
        }
      }
    } catch (error) {
      console.error(`获取${tabKey}选项卡数据失败:`, error);
      message.error(`获取${tabKey}选项卡数据失败`);
    } finally {
      setLoading(false);
    }
  };

  // 保存配置
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSaveLoading(true);
      
      // 基本设置保存到API
      if (activeTab === 'basic') {
        try {
          const res = await saveMiniProgramBaseSettings({
            name: values.name,
            description: values.description,
            logo: logoUrl
          });
          
          if (res.code === 0) {
            message.success('小程序基本设置已保存');
          } else {
            message.error(res.message || '保存小程序基本设置失败');
          }
        } catch (error) {
          console.error('保存小程序基本设置失败:', error);
          message.error('保存小程序基本设置失败，请检查网络连接');
        }
      } 
      // 广告设置保存到API
      else if (activeTab === 'ads') {
        try {
          const res = await saveAdSettings({
            enableWxAd: values.enableWxAd,
            rewardedVideoAdId: values.rewardedVideoAdId
          });
          
          if (res.code === 0) {
            message.success('广告设置已保存');
          } else {
            message.error(res.message || '保存广告设置失败');
          }
        } catch (error) {
          console.error('保存广告设置失败:', error);
          message.error('保存广告设置失败，请检查网络连接');
        }
      }
      // 奖励设置保存到API
      else if (activeTab === 'rewards') {
        try {
          const res = await saveRewardSettings({
            enableReward: values.enableRewards,
            firstViewMinRewardMin: values.firstTimeMinMinutes,
            firstViewMaxRewardDay: values.firstTimeMaxValue,
            singleAdMinRewardMin: values.rewardMinMinutes,
            singleAdMaxRewardDay: values.rewardMaxDays,
            dailyRewardLimit: values.dailyRewardLimit,
            dailyMaxAccumulatedDay: values.maxDailyRewardTime,
            rewardExpirationDays: values.rewardExpireDays
          });
          
          if (res.code === 0) {
            message.success('奖励设置已保存');
          } else {
            message.error(res.message || '保存奖励设置失败');
          }
        } catch (error) {
          console.error('保存奖励设置失败:', error);
          message.error('保存奖励设置失败，请检查网络连接');
        }
      }
      // 协议设置保存到API
      else if (activeTab === 'agreement') {
        try {
          const res = await saveAgreementSettings({
            privacyPolicy: values.privacyPolicy,
            userAgreement: values.userAgreement
          });
          
          if (res.code === 0) {
            message.success('协议设置已保存');
          } else {
            message.error(res.message || '保存协议设置失败');
          }
        } catch (error) {
          console.error('保存协议设置失败:', error);
          message.error('保存协议设置失败，请检查网络连接');
        }
      }
      // 专属管家设置保存到API
      else if (activeTab === 'exclusiveManager') {
        try {
          const response = await request({
            url: '/butler/image/save',
            method: 'post',
            data: {
              imageUrl: qrcodeUrl,
              status: 1
            }
          });
          
          console.log('保存专属管家图片返回:', response);
          
          // 判断response是否直接就是我们需要的数据格式
          if (response && typeof response === 'object') {
            let resData;
            
            // 判断返回数据结构
            if ('data' in response && 'code' in response) {
              // response本身就是最终数据
              resData = response;
            } else if ('data' in response && typeof response.data === 'object' && 'code' in response.data) {
              // response是axios包装的，data才是最终数据
              resData = response.data;
            } else {
              // 无法识别的数据结构
              console.error('无法识别的数据结构:', response);
              message.error('保存专属管家设置失败: 数据格式错误');
              return;
            }
            
            if (resData.code === 0) {
              message.success('专属管家设置已保存');
            } else {
              message.error(resData.message || '保存专属管家设置失败');
            }
          } else {
            message.error('保存专属管家设置失败: 返回数据无效');
          }
        } catch (error) {
          console.error('保存专属管家设置失败:', error);
          message.error('保存专属管家设置失败');
        }
      }
      // 分享设置保存到API
      else if (activeTab === 'shareSettings') {
        try {
          // 从表单和状态收集数据
          const shareSettings = {
            default_share_text: values.defaultShareTitle,
            default_share_image: defaultImageUrl,
            content_share_text: values.contentShareTitle,
            content_share_image: contentImageUrl,
            home_share_text: values.homeShareTitle,
            home_share_image: homeImageUrl
          };
          
          // 调用API保存分享设置
          const response = await saveShareSettings(shareSettings);
          
          console.log('保存分享设置返回:', response);
          
          if (response && response.code === 0) {
            message.success('分享设置已保存');
          } else {
            message.error(response.message || '保存分享设置失败');
          }
        } catch (error) {
          console.error('保存分享设置失败:', error);
          message.error('保存分享设置失败');
        }
      }
      else {
        // 准备提交数据（其他选项卡的模拟保存）
        const submitData = {
          ...values,
          logo: logoUrl,
        };
        
        console.log('提交的小程序配置:', submitData);
        
        // 模拟API请求延迟
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        message.success('小程序配置已保存');
      }
    } catch (error) {
      console.error('保存小程序配置失败:', error);
      message.error('保存失败，请检查表单填写是否正确');
    } finally {
      setSaveLoading(false);
    }
  };

  // 重置表单改为刷新数据
  const handleReset = () => {
    // 不再重置表单，只是重新获取数据
    fetchTabData(activeTab);
    message.info('数据已刷新');
  };

  // 处理文件选择
  const handleFileSelect = (url: string) => {
    setLogoUrl(url);
    form.setFieldsValue({ logo: url });
  };

  // 处理二维码选择
  const handleQrcodeSelect = (url: string) => {
    setQrcodeUrl(url);
    form.setFieldsValue({ qrcode: url });
  };

  // 处理内容默认图片选择
  const handleContentImageSelect = (url: string) => {
    setContentImageUrl(url);
    form.setFieldsValue({ contentDefaultImage: url });
  };
  
  // 处理首页默认图片选择
  const handleHomeImageSelect = (url: string) => {
    setHomeImageUrl(url);
    form.setFieldsValue({ homeDefaultImage: url });
  };

  // 处理默认分享图片选择
  const handleDefaultImageSelect = (url: string) => {
    setDefaultImageUrl(url);
    form.setFieldsValue({ defaultShareImage: url });
  };

  // 切换标签页时重新加载数据
  const handleTabChange = (key: string) => {
    setActiveTab(key);
    fetchTabData(key);
  };

  // 标签页配置
  const items: TabsProps['items'] = [
    {
      key: 'basic',
      label: '基本设置',
      children: (
        <>
          <Form.Item
            name="name"
            label="小程序名称"
            rules={[{ required: true, message: '请输入小程序名称' }]}
          >
            <Input placeholder="请输入小程序显示名称" />
          </Form.Item>

          <Form.Item
            name="description"
            label="小程序描述"
            rules={[{ required: true, message: '请输入小程序描述' }]}
          >
            <TextArea rows={4} placeholder="请输入小程序的简要描述" maxLength={500} showCount />
          </Form.Item>

          <Form.Item
            label="小程序Logo"
            tooltip="建议上传尺寸为300x300像素的图片，支持JPG、PNG格式"
          >
            <div className="logo-upload-container">
              {logoUrl ? (
                <div className="logo-preview">
                  <Image
                    src={logoUrl}
                    alt="小程序Logo"
                    width={128}
                    height={128}
                    style={{ objectFit: 'contain' }}
                  />
                  <div className="logo-actions">
                    <Button 
                      type="primary" 
                      onClick={() => setFileSelectVisible(true)}
                      icon={<UploadOutlined />}
                      size="small"
                    >
                      更换Logo
                    </Button>
                  </div>
                </div>
              ) : (
                <Button 
                  icon={<PlusOutlined />} 
                  onClick={() => setFileSelectVisible(true)}
                  className="upload-button"
                >
                  上传Logo
                </Button>
              )}
            </div>
            <Form.Item name="logo" hidden>
              <Input />
            </Form.Item>
          </Form.Item>

          {/* 文件选择器组件 */}
          <FileSelector
            visible={fileSelectVisible}
            onCancel={() => setFileSelectVisible(false)}
            onSelect={handleFileSelect}
            title="选择小程序Logo"
            accept="image/*"
          />
        </>
      ),
    },
    {
      key: 'ads',
      label: '广告设置',
      children: (
        <>
          <Alert
            message="智能接入广告"
            description="微信小程序现已支持免开发智能接入广告，无需填写大部分广告位ID和参数配置。只需在小程序后台开启广告功能即可。激励广告仍需手动配置。"
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <Form.Item
            name="enableWxAd"
            label="启用微信广告"
            valuePropName="checked"
          >
            <Switch checkedChildren="已启用" unCheckedChildren="已禁用" />
          </Form.Item>
          
          <Form.Item
            name="rewardedVideoAdId"
            label="激励视频广告位ID"
            rules={[{ required: true, message: '请输入激励视频广告位ID' }]}
          >
            <Input placeholder="请输入微信激励视频广告位ID" />
          </Form.Item>
          
          <Alert
            message="广告设置说明"
            description="激励视频广告需要手动配置广告位ID。其他类型广告请前往微信小程序管理后台 - 流量主 - 广告位管理，开通并配置智能接入广告。"
            type="warning"
            showIcon
            style={{ marginTop: 24 }}
          />
        </>
      ),
    },
    {
      key: 'rewards',
      label: '奖励设置',
      children: (
        <>
          <Alert
            message="广告奖励时长设置"
            description="用户观看广告后可获得的奖励时长设置，累积的时长可在商城中兑换商品或服务。"
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <Form.Item
            name="enableRewards"
            label="启用奖励功能"
            valuePropName="checked"
          >
            <Switch checkedChildren="已启用" unCheckedChildren="已禁用" />
          </Form.Item>
          
          <Row gutter={24}>
            <Col span={12}>
              <Form.Item
                label="首次观看广告奖励范围"
                tooltip="用户首次观看广告可获得的随机奖励范围"
              >
                <Space.Compact block>
                  <Form.Item
                    name="firstTimeMinMinutes"
                    noStyle
                    rules={[{ required: true, message: '请输入最小值' }]}
                  >
                    <InputNumber
                      style={{ width: 150 }}
                      min={1}
                      max={1440} // 最多24小时(1440分钟)
                      placeholder="最小值"
                      addonBefore="最小"
                      addonAfter="分钟"
                    />
                  </Form.Item>
                  <Input
                    style={{ width: 30, textAlign: 'center', pointerEvents: 'none', backgroundColor: '#fff' }}
                    placeholder="~"
                    disabled
                  />
                  <Form.Item
                    name="firstTimeMaxValue"
                    noStyle
                    rules={[{ required: true, message: '请输入最大值' }]}
                  >
                    <InputNumber
                      style={{ width: 150 }}
                      min={0.1}
                      max={365}
                      step={0.1}
                      precision={1}
                      placeholder="最大值"
                      addonBefore="最大"
                      addonAfter="天"
                    />
                  </Form.Item>
                </Space.Compact>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="单次广告随机奖励范围"
                tooltip="系统将在设定的最小值和最大值之间随机选择一个奖励时长"
              >
                <Space.Compact block>
                  <Form.Item
                    name="rewardMinMinutes"
                    noStyle
                    rules={[{ required: true, message: '请输入最小值' }]}
                  >
                    <InputNumber
                      style={{ width: 150 }}
                      min={1}
                      max={60}
                      placeholder="最小值"
                      addonBefore="最小"
                      addonAfter="分钟"
                    />
                  </Form.Item>
                  <Input
                    style={{ width: 30, textAlign: 'center', pointerEvents: 'none', backgroundColor: '#fff' }}
                    placeholder="~"
                    disabled
                  />
                  <Form.Item
                    name="rewardMaxDays"
                    noStyle
                    rules={[{ required: true, message: '请输入最大值' }]}
                  >
                    <InputNumber
                      style={{ width: 150 }}
                      min={0.1}
                      max={365}
                      step={0.1}
                      precision={1}
                      placeholder="最大值"
                      addonBefore="最大"
                      addonAfter="天"
                    />
                  </Form.Item>
                </Space.Compact>
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col span={12}>
              <Form.Item
                name="dailyRewardLimit"
                label="每日奖励次数上限"
                rules={[{ required: true, message: '请设置每日奖励次数上限' }]}
              >
                <InputNumber
                  min={1}
                  max={100}
                  style={{ width: '100%' }}
                  placeholder="用户每天最多可获得奖励的次数"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="maxDailyRewardTime"
                label="每日最大累计奖励"
                tooltip="限制用户每天最多能获得多少天的奖励时长"
                rules={[{ required: true, message: '请设置每日最大累计奖励' }]}
              >
                <InputNumber
                  min={0.1}
                  max={30}
                  step={0.1}
                  precision={1}
                  style={{ width: '100%' }}
                  placeholder="用户每天最多可获得的总时长"
                  addonAfter="天"
                />
              </Form.Item>
            </Col>
          </Row>
          
          <Form.Item
            name="rewardExpireDays"
            label="奖励过期天数"
            tooltip="设置为0表示永不过期"
            rules={[{ required: true, message: '请设置奖励过期天数' }]}
          >
            <InputNumber
              min={0}
              style={{ width: '100%' }}
              placeholder="奖励时长在多少天后过期，0表示永不过期"
              addonAfter="天"
            />
          </Form.Item>
        </>
      ),
    },
    {
      key: 'agreement',
      label: '协议设置',
      children: (
        <>
          <Form.Item
            name="privacyPolicy"
            label="隐私政策"
            rules={[{ required: true, message: '请输入隐私政策' }]}
          >
            <TextArea rows={8} placeholder="请输入隐私政策内容，支持Markdown格式" />
          </Form.Item>

          <Form.Item
            name="userAgreement"
            label="用户协议"
            rules={[{ required: true, message: '请输入用户协议' }]}
          >
            <TextArea rows={8} placeholder="请输入用户协议内容，支持Markdown格式" />
          </Form.Item>
        </>
      ),
    },
    {
      key: 'exclusiveManager',
      label: '专属管家',
      children: (
        <>
          <Alert
            message="专属管家设置"
            description="上传专属管家图片，方便用户获取专属客服联系方式。"
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <Form.Item
            label="专属管家图片"
            tooltip="建议上传清晰的图片，支持JPG、PNG格式"
          >
            <div className="logo-upload-container">
              {qrcodeUrl ? (
                <div className="logo-preview">
                  <Image
                    src={qrcodeUrl}
                    alt="专属管家图片"
                    width={128}
                    height={128}
                    style={{ objectFit: 'contain' }}
                  />
                  <div className="logo-actions">
                    <Button 
                      type="primary" 
                      onClick={() => setQrcodeSelectVisible(true)}
                      icon={<UploadOutlined />}
                      size="small"
                    >
                      更换图片
                    </Button>
                  </div>
                </div>
              ) : (
                <Button 
                  icon={<PlusOutlined />} 
                  onClick={() => setQrcodeSelectVisible(true)}
                  className="upload-button"
                >
                  上传图片
                </Button>
              )}
            </div>
            <Form.Item name="qrcode" hidden>
              <Input />
            </Form.Item>
          </Form.Item>

          {/* 文件选择器组件 */}
          <FileSelector
            visible={qrcodeSelectVisible}
            onCancel={() => setQrcodeSelectVisible(false)}
            onSelect={handleQrcodeSelect}
            title="选择专属管家图片"
            accept="image/*"
          />
        </>
      ),
    },
    {
      key: 'shareSettings',
      label: '分享设置',
      children: (
        <>
          <Alert
            message="分享设置"
            description="设置小程序分享时的默认图片和分享语，包括内容页分享和首页分享。"
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <Form.Item
            name="defaultShareTitle"
            label="默认分享语"
            tooltip="当没有指定分享语时使用的默认文案"
            rules={[{ required: true, message: '请输入默认分享语' }]}
          >
            <Input placeholder="请输入默认分享语" maxLength={30} showCount />
          </Form.Item>
          
          <Form.Item
            label="默认分享图片"
            tooltip="当没有指定分享图片时使用的默认图片"
            rules={[{ required: true, message: '请上传默认分享图片' }]}
          >
            <div className="logo-upload-container">
              {defaultImageUrl ? (
                <div className="logo-preview">
                  <Image
                    src={defaultImageUrl}
                    alt="默认分享图片"
                    width={128}
                    height={128}
                    style={{ objectFit: 'contain' }}
                  />
                  <div className="logo-actions">
                    <Button 
                      type="primary" 
                      onClick={() => setDefaultImageSelectVisible(true)}
                      icon={<UploadOutlined />}
                      size="small"
                    >
                      更换图片
                    </Button>
                  </div>
                </div>
              ) : (
                <Button 
                  icon={<PlusOutlined />} 
                  onClick={() => setDefaultImageSelectVisible(true)}
                  className="upload-button"
                >
                  上传图片
                </Button>
              )}
            </div>
            <Form.Item name="defaultShareImage" hidden>
              <Input />
            </Form.Item>
          </Form.Item>
          
          <Card title="内容页分享设置" variant="borderless" style={{ marginBottom: 24 }}>
            <Form.Item
              name="contentShareTitle"
              label="内容页分享语"
              tooltip="分享内容页时使用的文案，支持{title}变量，会替换为实际内容标题"
              rules={[{ required: true, message: '请输入内容页分享语' }]}
            >
              <Input placeholder="例如：推荐一篇好文章《{title}》，快来看看吧！" maxLength={40} showCount />
            </Form.Item>

            <Form.Item
              label="内容默认图片"
              tooltip="当分享内容页面时，如内容没有图片则使用此图片"
            >
              <div className="logo-upload-container">
                {contentImageUrl ? (
                  <div className="logo-preview">
                    <Image
                      src={contentImageUrl}
                      alt="内容默认图片"
                      width={128}
                      height={128}
                      style={{ objectFit: 'contain' }}
                    />
                    <div className="logo-actions">
                      <Button 
                        type="primary" 
                        onClick={() => setContentImageSelectVisible(true)}
                        icon={<UploadOutlined />}
                        size="small"
                      >
                        更换图片
                      </Button>
                    </div>
                  </div>
                ) : (
                  <Button 
                    icon={<PlusOutlined />} 
                    onClick={() => setContentImageSelectVisible(true)}
                    className="upload-button"
                  >
                    上传图片
                  </Button>
                )}
              </div>
              <Form.Item name="contentDefaultImage" hidden>
                <Input />
              </Form.Item>
            </Form.Item>
          </Card>
          
          <Card title="首页分享设置" variant="borderless" style={{ marginBottom: 24 }}>
            <Form.Item
              name="homeShareTitle"
              label="首页分享语"
              tooltip="分享小程序首页时使用的文案"
              rules={[{ required: true, message: '请输入首页分享语' }]}
            >
              <Input placeholder="请输入首页分享语" maxLength={30} showCount />
            </Form.Item>

            <Form.Item
              label="首页默认图片"
              tooltip="分享小程序首页时使用的默认图片"
            >
              <div className="logo-upload-container">
                {homeImageUrl ? (
                  <div className="logo-preview">
                    <Image
                      src={homeImageUrl}
                      alt="首页默认图片"
                      width={128}
                      height={128}
                      style={{ objectFit: 'contain' }}
                    />
                    <div className="logo-actions">
                      <Button 
                        type="primary" 
                        onClick={() => setHomeImageSelectVisible(true)}
                        icon={<UploadOutlined />}
                        size="small"
                      >
                        更换图片
                      </Button>
                    </div>
                  </div>
                ) : (
                  <Button 
                    icon={<PlusOutlined />} 
                    onClick={() => setHomeImageSelectVisible(true)}
                    className="upload-button"
                  >
                    上传图片
                  </Button>
                )}
              </div>
              <Form.Item name="homeDefaultImage" hidden>
                <Input />
              </Form.Item>
            </Form.Item>
          </Card>

          {/* 文件选择器组件 */}
          <FileSelector
            visible={contentImageSelectVisible}
            onCancel={() => setContentImageSelectVisible(false)}
            onSelect={handleContentImageSelect}
            title="选择内容默认图片"
            accept="image/*"
          />
          <FileSelector
            visible={homeImageSelectVisible}
            onCancel={() => setHomeImageSelectVisible(false)}
            onSelect={handleHomeImageSelect}
            title="选择首页默认图片"
            accept="image/*"
          />
          <FileSelector
            visible={defaultImageSelectVisible}
            onCancel={() => setDefaultImageSelectVisible(false)}
            onSelect={handleDefaultImageSelect}
            title="选择默认分享图片"
            accept="image/*"
          />
        </>
      ),
    },
  ];

  return (
    <div className="mini-program-container">
      <Card className="main-dashboard-card">
        <div className="category-header">
          <Title level={4}>基础设置</Title>
          <Space>
            <Button 
              type="primary" 
              icon={<SaveOutlined />} 
              onClick={handleSubmit}
              loading={saveLoading}
            >
              保存配置
            </Button>
            <Button 
              icon={<ReloadOutlined />}
              onClick={handleReset}
            >
              刷新
            </Button>
          </Space>
        </div>
        
        <Form
          form={form}
          layout="vertical"
          disabled={loading}
        >
          <Tabs
            activeKey={activeTab}
            onChange={handleTabChange}
            items={items}
            className="mp-settings-tabs"
          />
        </Form>
      </Card>
    </div>
  );
};

// 包装组件，添加AntdApp上下文
const MiniProgram: React.FC = () => {
  return (
    <App>
      <MiniProgramContent />
    </App>
  );
};

export default MiniProgram; 