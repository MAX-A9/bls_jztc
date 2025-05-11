import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Form,
  Input,
  Button,
  Space,
  Select,
  Switch,
  message,
  Row,
  Col,
  Alert,
} from 'antd';
import {
  SaveOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import './StorageSettings.css';
import { getStorageConfig, saveStorageConfig, SaveStorageConfigParams } from '../../api/storage';

const { Title } = Typography;
const { Option } = Select;

const StorageSettings: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  // 获取存储配置
  const fetchStorageConfig = async () => {
    setLoading(true);
    try {
      const response = await getStorageConfig();
      // 处理嵌套的数据
      const configData = response.data || response;
      
      // 设置表单初始值
      form.setFieldsValue({
        accessKeyId: configData.accessKeyId || '',
        accessKeySecret: configData.accessKeySecret || '',
        endpoint: configData.endpoint || '',
        bucket: configData.bucket || '',
        region: configData.region || '',
        directory: configData.directory || '',
        publicAccess: configData.publicAccess || false,
      });
    } catch (error) {
      console.error('获取存储配置失败:', error);
      message.error('获取存储配置失败，请刷新重试');
      
      // 设置一些默认值
      form.setFieldsValue({
        accessKeyId: '',
        accessKeySecret: '',
        endpoint: '',
        bucket: '',
        region: '',
        directory: 'uploads',
        publicAccess: false,
      });
    } finally {
      setLoading(false);
    }
  };

  // 初始化数据
  useEffect(() => {
    fetchStorageConfig();
  }, []);

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      
      // 构建请求参数
      const params: SaveStorageConfigParams = {
        accessKeyId: values.accessKeyId,
        accessKeySecret: values.accessKeySecret,
        endpoint: values.endpoint,
        bucket: values.bucket,
        region: values.region,
        directory: values.directory,
        publicAccess: values.publicAccess,
      };
      
      // 调用保存接口
      const response = await saveStorageConfig(params);
      const responseData = response.data || response;
      
      if (responseData.success) {
        message.success(responseData.responseMessage || '存储设置已更新');
      } else {
        message.error(responseData.responseMessage || '保存失败，请检查配置');
      }
    } catch (error) {
      console.error('保存存储配置失败:', error);
      message.error('保存失败，请检查表单信息');
    } finally {
      setLoading(false);
    }
  };

  // 刷新配置
  const handleRefresh = () => {
    fetchStorageConfig();
  };

  return (
    <div className="storage-settings-container">
      <Card className="storage-card">
        <div className="storage-header">
          <Title level={4}>存储设置</Title>
          <Space>
            <Button
              icon={<ReloadOutlined />}
              onClick={handleRefresh}
              loading={loading}
            >
              刷新
            </Button>
            <Button
              type="primary"
              onClick={handleSubmit}
              loading={loading}
              icon={<SaveOutlined />}
            >
              保存配置
            </Button>
          </Space>
        </div>

        <Form
          form={form}
          layout="vertical"
          className="storage-form"
        >
          <Alert
            message="阿里云OSS配置"
            description="设置阿里云对象存储服务(OSS)的相关配置。AccessKey信息具有高权限，请妥善保管。"
            type="info"
            showIcon
            className="storage-info"
          />
          <Row gutter={[16, 16]}>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Item
                name="accessKeyId"
                label="AccessKey ID"
                rules={[{ required: true, message: '请输入AccessKey ID' }]}
              >
                <Input placeholder="请输入AccessKey ID" />
              </Form.Item>
            </Col>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Item
                name="accessKeySecret"
                label="AccessKey Secret"
                rules={[{ required: true, message: '请输入AccessKey Secret' }]}
              >
                <Input.Password placeholder="请输入AccessKey Secret" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Item
                name="region"
                label="区域"
                rules={[{ required: true, message: '请选择区域' }]}
              >
                <Select placeholder="请选择OSS区域">
                  <Option value="oss-cn-hangzhou">华东1（杭州）</Option>
                  <Option value="oss-cn-shanghai">华东2（上海）</Option>
                  <Option value="oss-cn-beijing">华北2（北京）</Option>
                  <Option value="oss-cn-shenzhen">华南1（深圳）</Option>
                  <Option value="oss-cn-hongkong">香港</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Item
                name="endpoint"
                label="Endpoint"
                rules={[{ required: true, message: '请输入Endpoint' }]}
              >
                <Input placeholder="例如: oss-cn-beijing.aliyuncs.com" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Item
                name="bucket"
                label="Bucket"
                rules={[{ required: true, message: '请输入Bucket名称' }]}
              >
                <Input placeholder="请输入Bucket名称" />
              </Form.Item>
            </Col>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Item
                name="directory"
                label="存储目录"
                tooltip="OSS中的目录前缀，可留空"
              >
                <Input placeholder="例如: uploads/（可选）" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item
            name="publicAccess"
            label="公开访问"
            valuePropName="checked"
            tooltip="启用后，上传的文件将可以被公开访问"
          >
            <Switch checkedChildren="开启" unCheckedChildren="关闭" />
          </Form.Item>
          
          <Alert
            message="温馨提示"
            description="修改存储设置后，已上传的文件可能无法访问。请确保已备份重要数据。"
            type="warning"
            showIcon
            className="storage-warning"
          />
        </Form>
      </Card>
    </div>
  );
};

export default StorageSettings; 