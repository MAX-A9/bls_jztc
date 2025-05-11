import React, { useState, useEffect } from 'react';
import { Button, Card, Form, Input, Typography, message, Row, Col, Spin } from 'antd';
import { UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { login, getCaptcha, LoginParams } from '../api/auth';
import useUserStore from '../store/userStore';
import './Login.css';

const { Title } = Typography;

const Login: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [captchaLoading, setCaptchaLoading] = useState(false);
  const [captchaImage, setCaptchaImage] = useState('');
  const [captchaId, setCaptchaId] = useState('');
  const { setToken, setUserInfo } = useUserStore();

  // 获取验证码
  const fetchCaptcha = async () => {
    try {
      setCaptchaLoading(true);
      const res = await getCaptcha();
      // 接口返回的数据结构是 { code: 0, message: 'OK', data: { id, base64, expiredAt } }
      const captchaData = res.data || res;
      setCaptchaId(captchaData.id || '');
      setCaptchaImage(captchaData.base64 || '');
    } catch (error) {
      message.error('获取验证码失败，请刷新重试');
      console.error('获取验证码失败:', error);
    } finally {
      setCaptchaLoading(false);
    }
  };

  // 组件挂载时加载验证码
  useEffect(() => {
    fetchCaptcha();
  }, []);

  // 处理登录提交
  const handleSubmit = async (values: LoginParams) => {
    try {
      setLoading(true);
      values.captchaId = captchaId;
      
      const res = await login(values);
      const loginData = res.data || res;
      
      // 存储用户信息和token
      setToken(loginData.token || '');
      setUserInfo({
        id: loginData.userId || 0,
        username: values.username,
        realName: loginData.nickname || '',
        role: 'user',
        permissions: []
      });
      
      // 保存用户名到本地存储
      localStorage.setItem('username', values.username);
      
      message.success('登录成功');
      navigate('/dashboard');
    } catch (error: any) {
      message.error(error.message || '登录失败，请检查用户名和密码');
      // 刷新验证码
      fetchCaptcha();
      // 清空验证码输入框
      form.setFieldsValue({ captchaCode: '' });
    } finally {
      setLoading(false);
    }
  };

  // 从本地存储获取保存的用户名
  const savedUsername = localStorage.getItem('username') || '';

  return (
    <div className="login-container">
      <div className="login-content">
        <Card className="login-card">
          <div className="login-header">
            <Title level={2} className="login-title">系统登录</Title>
          </div>
          
          <Form
            form={form}
            name="login"
            initialValues={{ 
              username: savedUsername 
            }}
            onFinish={handleSubmit}
            autoComplete="off"
            size="large"
          >
            <Form.Item
              name="username"
              rules={[{ required: true, message: '请输入用户名' }]}
            >
              <Input prefix={<UserOutlined />} placeholder="用户名" />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[{ required: true, message: '请输入密码' }]}
            >
              <Input.Password prefix={<LockOutlined />} placeholder="密码" />
            </Form.Item>

            <Form.Item
              name="captchaCode"
              rules={[{ required: true, message: '请输入验证码' }]}
            >
              <Row gutter={8}>
                <Col span={16}>
                  <Input prefix={<SafetyOutlined />} placeholder="验证码" />
                </Col>
                <Col span={8}>
                  <div 
                    className="captcha-container" 
                    onClick={fetchCaptcha}
                    title="点击刷新验证码"
                    style={{ border: '1px solid #d9d9d9', height: '38px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}
                  >
                    {captchaLoading ? (
                      <div className="captcha-loading">
                        <Spin size="small" />
                      </div>
                    ) : captchaImage ? (
                      <div style={{ 
                        width: '100%', 
                        height: '100%', 
                        display: 'flex', 
                        alignItems: 'center', 
                        justifyContent: 'center',
                        backgroundColor: '#e8f0fe',
                        backgroundImage: 'linear-gradient(45deg, transparent 45%, #d0e1ff 45%, #d0e1ff 55%, transparent 55%)',
                        backgroundSize: '6px 6px',
                        borderRadius: '2px',
                        position: 'relative',
                        overflow: 'hidden',
                        userSelect: 'none',
                        WebkitUserSelect: 'none',
                        msUserSelect: 'none',
                        MozUserSelect: 'none',
                      }}>
                        {/* 添加干扰线 */}
                        <div style={{ 
                          position: 'absolute',
                          width: '120%',
                          height: '1px',
                          background: '#7597dd',
                          top: '50%',
                          left: '-10%',
                          transform: 'rotate(-5deg)',
                        }}></div>
                        <div style={{ 
                          position: 'absolute',
                          width: '120%',
                          height: '1px',
                          background: '#7597dd',
                          top: '30%',
                          left: '-10%',
                          transform: 'rotate(3deg)',
                        }}></div>
                        
                        {/* 显示验证码图片 */}
                        <img 
                          src={captchaImage} 
                          alt="验证码" 
                          className="captcha-image" 
                          style={{ 
                            maxWidth: '100%', 
                            maxHeight: '100%', 
                            display: 'block', 
                            pointerEvents: 'none',
                            objectFit: 'contain',
                            zIndex: 5,
                            position: 'relative'
                          }} 
                        />
                        
                        {/* 添加干扰点 */}
                        {Array.from({ length: 20 }).map((_, i) => (
                          <div 
                            key={`dot-${i}`}
                            style={{
                              position: 'absolute',
                              width: '2px',
                              height: '2px',
                              background: '#7597dd',
                              top: `${Math.random() * 100}%`,
                              left: `${Math.random() * 100}%`,
                              zIndex: 1
                            }}
                          ></div>
                        ))}
                      </div>
                    ) : (
                      <div className="captcha-placeholder">获取验证码</div>
                    )}
                  </div>
                </Col>
              </Row>
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                className="login-form-button"
                block
                loading={loading}
              >
                登录
              </Button>
            </Form.Item>
          </Form>
          
          <div className="login-footer">
            
          </div>
        </Card>
      </div>
      
      <div className="login-copyright">
        © {new Date().getFullYear()} 企业管理系统 - 版权所有
      </div>
    </div>
  );
};

export default Login; 