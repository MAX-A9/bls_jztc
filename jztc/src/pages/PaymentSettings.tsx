import React, { useState } from 'react';
import {
  Card,
  Typography,
  Form,
  Input,
  Button,
  Switch,
  Space,
  Divider,
  App,
} from 'antd';
import {
  WechatOutlined,
  SaveOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import './PaymentSettings.css';
import { getPaymentConfig, savePaymentConfig, SavePaymentConfigParams } from '../api/payment';

const { Title } = Typography;

const PaymentSettings: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const { message } = App.useApp();

  // 获取当前域名
  const getCurrentDomain = () => {
    return window.location.hostname;
  };

  // 构建默认回调地址
  const getDefaultNotifyUrl = () => {
    const domain = getCurrentDomain();
    return `https://${domain}/wx/pay/notify`;
  };

  // 获取微信支付配置
  const fetchWechatPayConfig = async () => {
    setLoading(true);
    try {
      const response = await getPaymentConfig();
      // 处理嵌套的数据
      const configData = response.data || response;
      
      // 设置表单初始值
      form.setFieldsValue({
        appId: configData.appId || '',
        mchId: configData.mchId || '',
        apiKey: configData.apiKey || '',
        notifyUrl: configData.notifyUrl || getDefaultNotifyUrl(),
        isEnabled: configData.isEnabled || false,
      });
    } catch (error) {
      console.error('获取支付配置失败:', error);
      message.error('获取支付配置失败，请刷新重试');
      
      // 设置默认值
      form.setFieldsValue({
        appId: '',
        mchId: '',
        apiKey: '',
        notifyUrl: getDefaultNotifyUrl(),
        isEnabled: false,
      });
    } finally {
      setLoading(false);
    }
  };

  // 保存微信支付配置
  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      
      // 构建请求参数
      const params: SavePaymentConfigParams = {
        appId: values.appId,
        mchId: values.mchId,
        apiKey: values.apiKey,
        notifyUrl: values.notifyUrl,
        isEnabled: values.isEnabled,
      };
      
      // 调用保存接口
      const response = await savePaymentConfig(params);
      const responseData = response.data || response;
      
      if (responseData.success) {
        message.success(responseData.responseMessage || '支付设置已更新');
      } else {
        message.error(responseData.responseMessage || '保存失败，请检查配置');
      }
    } catch (error) {
      console.error('保存支付配置失败:', error);
      message.error('保存失败，请检查表单填写是否正确');
    } finally {
      setLoading(false);
    }
  };

  // 重置为默认回调地址
  const resetToDefaultNotifyUrl = () => {
    form.setFieldValue('notifyUrl', getDefaultNotifyUrl());
    message.success('已重置为默认回调地址');
  };

  // 组件加载时获取配置
  React.useEffect(() => {
    fetchWechatPayConfig();
  }, []);

  return (
    <div className="payment-settings-container">
      <Card className="payment-settings-card">
        <div className="payment-settings-header">
          <Title level={4}>支付设置</Title>
          <Space>
            <Button
              icon={<ReloadOutlined />}
              onClick={fetchWechatPayConfig}
              loading={loading}
            >
              刷新
            </Button>
            <Button
              type="primary"
              icon={<SaveOutlined />}
              onClick={handleSave}
              loading={loading}
            >
              保存
            </Button>
          </Space>
        </div>

        <Form
          form={form}
          layout="vertical"
        >
          <div className="payment-section">
            <div className="section-header">
              <div className="section-title">
                <WechatOutlined style={{ fontSize: '20px', color: '#07C160', marginRight: '8px' }} />
                <Title level={5} style={{ margin: 0 }}>微信支付配置</Title>
              </div>
              <Form.Item
                name="isEnabled"
                valuePropName="checked"
                style={{ marginBottom: 0 }}
              >
                <Switch checkedChildren="启用" unCheckedChildren="禁用" />
              </Form.Item>
            </div>
            <Divider />
            
            <Form.Item
              name="appId"
              label="AppID"
              rules={[{ required: true, message: '请输入AppID' }]}
            >
              <Input placeholder="请输入小程序AppID" />
            </Form.Item>

            <Form.Item
              name="mchId"
              label="商户号"
              rules={[{ required: true, message: '请输入商户号' }]}
            >
              <Input placeholder="请输入微信支付商户号" />
            </Form.Item>

            <Form.Item
              name="apiKey"
              label="API密钥"
              rules={[{ required: true, message: '请输入API密钥' }]}
            >
              <Input.Password placeholder="请输入微信支付API密钥" />
            </Form.Item>

            <Form.Item
              name="notifyUrl"
              label="回调通知地址"
              rules={[{ required: true, message: '请输入回调通知地址' }]}
            >
              <Input 
                placeholder="请输入支付回调通知地址" 
                addonAfter={
                  <Button 
                    type="link" 
                    style={{ padding: 0, fontSize: '12px' }} 
                    onClick={resetToDefaultNotifyUrl}
                  >
                    使用默认地址
                  </Button>
                }
              />
            </Form.Item>
          </div>
        </Form>
      </Card>
    </div>
  );
};

export default PaymentSettings; 