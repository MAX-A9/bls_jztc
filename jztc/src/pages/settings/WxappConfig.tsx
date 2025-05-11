import React, { useState, useEffect } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  Switch,
  message,
  Typography,
  Spin,
  Alert,
  Space,
  Row,
  Col,
} from 'antd';
import { SaveOutlined, ReloadOutlined } from '@ant-design/icons';
import { getWxappConfig, saveWxappConfig, WxappConfig as WxappConfigType } from '@/api/wxapp';
import './WxappConfig.css';

const { Title } = Typography;

const WxappConfig: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saveLoading, setSaveLoading] = useState(false);

  // 获取微信配置
  const fetchWxappConfig = async () => {
    setLoading(true);
    try {
      const res = await getWxappConfig();
      if (res.code === 0) {
        form.setFieldsValue(res.data);
      } else {
        message.error(res.message || '获取微信配置失败');
      }
    } catch (error) {
      console.error('获取微信配置失败:', error);
      message.error('获取微信配置失败，请检查网络');
    } finally {
      setLoading(false);
    }
  };

  // 保存微信配置
  const handleSave = async (values: WxappConfigType) => {
    setSaveLoading(true);
    try {
      const res = await saveWxappConfig(values);
      if (res.code === 0) {
        message.success('保存微信配置成功');
      } else {
        message.error(res.message || '保存微信配置失败');
      }
    } catch (error) {
      console.error('保存微信配置失败:', error);
      message.error('保存微信配置失败，请检查网络');
    } finally {
      setSaveLoading(false);
    }
  };

  // 初次加载获取配置
  useEffect(() => {
    fetchWxappConfig();
  }, []);

  return (
    <div className="wxapp-config-container">
      <Card className="wxapp-card">
        <div className="wxapp-header">
          <Title level={4}>微信小程序配置</Title>
          <Space>
            <Button
              icon={<ReloadOutlined />}
              onClick={fetchWxappConfig}
              loading={loading}
            >
              刷新
            </Button>
            <Button
              type="primary"
              onClick={() => form.submit()}
              loading={saveLoading}
              icon={<SaveOutlined />}
            >
              保存配置
            </Button>
          </Space>
        </div>

        <div className="wxapp-form-container" style={{ minHeight: '300px' }}>
          <Spin spinning={loading} tip="加载中..." size="large">
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSave}
              initialValues={{ enabled: false }}
              className="wxapp-form"
            >
              <Alert
                message="微信小程序配置"
                description="配置微信小程序的AppID和AppSecret，启用微信小程序登录功能。"
                type="info"
                showIcon
                className="wxapp-info"
              />
              
              <Row gutter={[16, 16]}>
                <Col xs={24} sm={24} md={12} lg={12} xl={12}>
                  <Form.Item
                    name="appId"
                    label="小程序AppID"
                    rules={[{ required: true, message: '请输入小程序AppID' }]}
                  >
                    <Input placeholder="请输入小程序AppID" />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={24} md={12} lg={12} xl={12}>
                  <Form.Item
                    name="appSecret"
                    label="小程序AppSecret"
                    rules={[{ required: true, message: '请输入小程序AppSecret' }]}
                  >
                    <Input.Password placeholder="请输入小程序AppSecret" />
                  </Form.Item>
                </Col>
              </Row>
              
              <Form.Item
                name="enabled"
                label="启用状态"
                valuePropName="checked"
                tooltip="启用后，用户可以通过微信小程序进行登录"
              >
                <Switch checkedChildren="已启用" unCheckedChildren="已禁用" />
              </Form.Item>

              <Alert
                message="安全提示"
                description="请妥善保管您的AppSecret，不要泄露给他人，以防被恶意使用。"
                type="warning"
                showIcon
                className="wxapp-warning"
              />
            </Form>
          </Spin>
        </div>
      </Card>
    </div>
  );
};

export default WxappConfig; 