import React, { useState } from 'react';
import { message, Button, Input, Form, Card } from 'antd';
import { verifyLicense, LicenseVerifyReq } from '../api/license';
import { useNavigate } from 'react-router-dom';
import './License.css';

const License: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const handleSubmit = async (values: { license_code: string }) => {
    try {
      setLoading(true);
      // 获取当前域名
      const domain = window.location.hostname;
      // 组装请求参数
      const params: LicenseVerifyReq = {
        domain,
        license_code: values.license_code
      };

      // 请求接口
      const response = await verifyLicense(params);
      
      // 使用类型断言
      const responseData = response as any;
      
      // 响应拦截器已经处理了返回的数据，如果能执行到这里，说明code=0
      if (responseData && responseData.data && responseData.data.valid) {
        message.success('授权验证成功');
        // 授权成功后跳转到登录页
        navigate('/login');
      } else {
        // 这部分代码可能不会执行，因为如果code不为0，响应拦截器会抛出异常
        message.error('授权验证失败，请检查授权码是否正确');
      }
    } catch (error) {
      console.error('授权验证出错:', error);
      message.error('授权验证出错，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="license-container">
      <div className="license-content">
        <Card className="license-card" title="系统授权">
          <div className="license-description">
            <p>请输入您的授权码以继续使用系统。</p>
            <p>如果您没有授权码，请联系：zl93339。</p>
          </div>
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            autoComplete="off"
          >
            <Form.Item
              name="license_code"
              label="授权码"
              rules={[{ required: true, message: '请输入授权码' }]}
            >
              <Input.TextArea
                placeholder="请输入您的授权码"
                rows={4}
                className="license-input"
              />
            </Form.Item>
            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                className="license-button"
              >
                验证授权
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>
    </div>
  );
};

export default License; 