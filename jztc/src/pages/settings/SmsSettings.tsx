import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Form,
  Input,
  Button,
  Switch,
  Select,
  Space,
  Table,
  Tag,
  message,
  Modal,
  Row,
  Col,
  Tabs
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined
} from '@ant-design/icons';
import './SmsSettings.css';

const { Title, Text } = Typography;
const { Option } = Select;
const { TabPane } = Tabs;
const { TextArea } = Input;

// SMS渠道数据接口
interface SmsChannel {
  id: number;
  name: string;
  appId: string;
  appKey: string;
  signName: string;
  status: boolean;
  isDefault: boolean;
}

// 短信模板数据接口
interface SmsTemplate {
  id: number;
  name: string;
  code: string;
  content: string;
  channelId: number;
  type: string;
  status: boolean;
}

const SmsSettings: React.FC = () => {
  // 状态定义
  const [loading, setLoading] = useState(false);
  const [channels, setChannels] = useState<SmsChannel[]>([]);
  const [templates, setTemplates] = useState<SmsTemplate[]>([]);
  const [filteredChannels, setFilteredChannels] = useState<SmsChannel[]>([]);
  const [filteredTemplates, setFilteredTemplates] = useState<SmsTemplate[]>([]);
  const [channelModalVisible, setChannelModalVisible] = useState(false);
  const [templateModalVisible, setTemplateModalVisible] = useState(false);
  const [currentChannel, setCurrentChannel] = useState<SmsChannel | null>(null);
  const [currentTemplate, setCurrentTemplate] = useState<SmsTemplate | null>(null);
  const [channelForm] = Form.useForm();
  const [templateForm] = Form.useForm();
  const [activeTab, setActiveTab] = useState('1');

  // 初始化加载数据
  useEffect(() => {
    fetchData();
  }, []);

  // 更新过滤后数据
  useEffect(() => {
    setFilteredChannels(channels);
  }, [channels]);

  useEffect(() => {
    setFilteredTemplates(templates);
  }, [templates]);

  // 获取数据
  const fetchData = async () => {
    setLoading(true);
    try {
      // 模拟API请求
      await new Promise(resolve => setTimeout(resolve, 800));
      
      // 模拟渠道数据
      const mockChannels: SmsChannel[] = [
        {
          id: 1,
          name: '阿里云短信',
          appId: 'LTAI5txxxxxxxxxxxxxxx',
          appKey: '********',
          signName: '测试签名',
          status: true,
          isDefault: true,
        },
        {
          id: 2,
          name: '腾讯云短信',
          appId: 'SDK_APPIDxxxxxxxx',
          appKey: '********',
          signName: '测试签名2',
          status: false,
          isDefault: false,
        },
      ];
      
      // 模拟模板数据
      const mockTemplates: SmsTemplate[] = [
        {
          id: 1,
          name: '注册验证码',
          code: 'SMS_123456789',
          content: '您的验证码为：${code}，有效期10分钟，请勿泄露给他人。',
          channelId: 1,
          type: '验证码',
          status: true,
        },
        {
          id: 2,
          name: '找回密码验证码',
          code: 'SMS_987654321',
          content: '您正在找回密码，验证码为：${code}，有效期10分钟，请勿泄露给他人。',
          channelId: 1,
          type: '验证码',
          status: true,
        },
        {
          id: 3,
          name: '订单支付成功通知',
          code: 'SMS_567891234',
          content: '您的订单${orderNo}已支付成功，感谢您的购买！',
          channelId: 2,
          type: '通知',
          status: true,
        },
      ];
      
      setChannels(mockChannels);
      setTemplates(mockTemplates);
    } catch (error) {
      console.error('获取数据失败:', error);
      message.error('获取数据失败');
    } finally {
      setLoading(false);
    }
  };

  // 添加短信渠道
  const showAddChannelModal = () => {
    setCurrentChannel(null);
    channelForm.resetFields();
    setChannelModalVisible(true);
  };

  // 编辑短信渠道
  const showEditChannelModal = (channel: SmsChannel) => {
    setCurrentChannel(channel);
    channelForm.setFieldsValue(channel);
    setChannelModalVisible(true);
  };

  // 处理渠道表单提交
  const handleChannelSubmit = async () => {
    try {
      const values = await channelForm.validateFields();
      
      if (currentChannel) {
        // 更新渠道
        setChannels(channels.map(item => 
          item.id === currentChannel.id ? { ...item, ...values } : item
        ));
        message.success('渠道更新成功');
      } else {
        // 添加渠道
        const newChannel = {
          id: channels.length > 0 ? Math.max(...channels.map(item => item.id)) + 1 : 1,
          ...values,
        };
        setChannels([...channels, newChannel]);
        message.success('渠道添加成功');
      }
      
      setChannelModalVisible(false);
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  // 删除短信渠道
  const handleDeleteChannel = (id: number) => {
    // 检查是否是默认渠道
    const channel = channels.find(item => item.id === id);
    if (channel?.isDefault) {
      message.error('不能删除默认渠道');
      return;
    }
    
    setChannels(channels.filter(item => item.id !== id));
    message.success('渠道删除成功');
  };

  // 设置默认渠道
  const handleSetDefaultChannel = (id: number) => {
    setChannels(channels.map(item => ({
      ...item,
      isDefault: item.id === id,
    })));
    message.success('默认渠道设置成功');
  };

  // 添加短信模板
  const showAddTemplateModal = () => {
    setCurrentTemplate(null);
    templateForm.resetFields();
    setTemplateModalVisible(true);
  };

  // 编辑短信模板
  const showEditTemplateModal = (template: SmsTemplate) => {
    setCurrentTemplate(template);
    templateForm.setFieldsValue(template);
    setTemplateModalVisible(true);
  };

  // 处理模板表单提交
  const handleTemplateSubmit = async () => {
    try {
      const values = await templateForm.validateFields();
      
      if (currentTemplate) {
        // 更新模板
        setTemplates(templates.map(item => 
          item.id === currentTemplate.id ? { ...item, ...values } : item
        ));
        message.success('模板更新成功');
      } else {
        // 添加模板
        const newTemplate = {
          id: templates.length > 0 ? Math.max(...templates.map(item => item.id)) + 1 : 1,
          ...values,
        };
        setTemplates([...templates, newTemplate]);
        message.success('模板添加成功');
      }
      
      setTemplateModalVisible(false);
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  // 删除短信模板
  const handleDeleteTemplate = (id: number) => {
    setTemplates(templates.filter(item => item.id !== id));
    message.success('模板删除成功');
  };

  // 渠道表格列定义
  const channelColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '渠道名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => <Text>{text}</Text>,
    },
    {
      title: 'AppID',
      dataIndex: 'appId',
      key: 'appId',
    },
    {
      title: '签名',
      dataIndex: 'signName',
      key: 'signName',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: boolean) => (
        status ? <Tag color="success">启用</Tag> : <Tag color="error">禁用</Tag>
      ),
      filters: [
        { text: '启用', value: true },
        { text: '禁用', value: false },
      ],
      onFilter: (value: any, record: SmsChannel) => record.status === value,
    },
    {
      title: '默认渠道',
      dataIndex: 'isDefault',
      key: 'isDefault',
      render: (isDefault: boolean) => (
        isDefault ? <Tag color="blue">是</Tag> : <Tag color="default">否</Tag>
      ),
      filters: [
        { text: '是', value: true },
        { text: '否', value: false },
      ],
      onFilter: (value: any, record: SmsChannel) => record.isDefault === value,
    },
    {
      title: '操作',
      key: 'action',
      width: 220,
      render: (_: any, record: SmsChannel) => (
        <Space size="small">
          <Button 
            type="primary" 
            size="small" 
            icon={<EditOutlined />}
            onClick={() => showEditChannelModal(record)}
          >
            编辑
          </Button>
          {!record.isDefault && (
            <>
              <Button 
                type="primary" 
                size="small" 
                ghost
                onClick={() => handleSetDefaultChannel(record.id)}
              >
                设为默认
              </Button>
              <Button 
                danger 
                size="small" 
                icon={<DeleteOutlined />}
                onClick={() => handleDeleteChannel(record.id)}
              >
                删除
              </Button>
            </>
          )}
        </Space>
      ),
    },
  ];

  // 模板表格列定义
  const templateColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '模板名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => <Text>{text}</Text>,
    },
    {
      title: '模板CODE',
      dataIndex: 'code',
      key: 'code',
    },
    {
      title: '模板类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => {
        if (type === '验证码') {
          return <Tag color="blue">验证码</Tag>;
        } else if (type === '通知') {
          return <Tag color="green">通知</Tag>;
        } else {
          return <Tag color="orange">{type}</Tag>;
        }
      },
      filters: [
        { text: '验证码', value: '验证码' },
        { text: '通知', value: '通知' },
        { text: '推广', value: '推广' },
      ],
      onFilter: (value: any, record: SmsTemplate) => record.type === value,
    },
    {
      title: '所属渠道',
      dataIndex: 'channelId',
      key: 'channelId',
      render: (channelId: number) => {
        const channel = channels.find(item => item.id === channelId);
        return channel ? channel.name : '-';
      },
      filters: channels.map(channel => ({ text: channel.name, value: channel.id })),
      onFilter: (value: any, record: SmsTemplate) => record.channelId === value,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: boolean) => (
        status ? <Tag color="success">启用</Tag> : <Tag color="error">禁用</Tag>
      ),
      filters: [
        { text: '启用', value: true },
        { text: '禁用', value: false },
      ],
      onFilter: (value: any, record: SmsTemplate) => record.status === value,
    },
    {
      title: '操作',
      key: 'action',
      width: 160,
      render: (_: any, record: SmsTemplate) => (
        <Space size="small">
          <Button 
            type="primary" 
            size="small" 
            icon={<EditOutlined />}
            onClick={() => showEditTemplateModal(record)}
          >
            编辑
          </Button>
          <Button 
            danger 
            size="small" 
            icon={<DeleteOutlined />}
            onClick={() => handleDeleteTemplate(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="sms-settings-container">
      <Card className="sms-card">
        <div className="sms-header">
          <Title level={4}>短信设置</Title>
          <Space>
            {activeTab === '1' ? (
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={showAddChannelModal}
              >
                添加渠道
              </Button>
            ) : (
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={showAddTemplateModal}
              >
                添加模板
              </Button>
            )}
            <Button
              icon={<ReloadOutlined />}
              onClick={fetchData}
            >
              刷新
            </Button>
          </Space>
        </div>

        <Tabs 
          activeKey={activeTab} 
          onChange={setActiveTab}
          className="sms-tabs"
        >
          <TabPane tab="短信渠道" key="1">
            <Table 
              className="sms-table"
              dataSource={filteredChannels.length > 0 ? filteredChannels : channels} 
              columns={channelColumns} 
              rowKey="id"
              loading={loading}
              pagination={{
                defaultPageSize: 10,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 个渠道`,
              }}
            />
          </TabPane>
          
          <TabPane tab="短信模板" key="2">
            <Table 
              className="sms-table"
              dataSource={filteredTemplates.length > 0 ? filteredTemplates : templates} 
              columns={templateColumns} 
              rowKey="id"
              loading={loading}
              pagination={{
                defaultPageSize: 10,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 个模板`,
              }}
            />
          </TabPane>
        </Tabs>
      </Card>

      {/* 渠道配置模态框 */}
      <Modal
        title={currentChannel ? '编辑渠道' : '添加渠道'}
        open={channelModalVisible}
        onOk={handleChannelSubmit}
        onCancel={() => setChannelModalVisible(false)}
        maskClosable={false}
        width={600}
      >
        <Form
          form={channelForm}
          layout="vertical"
          className="sms-form"
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="name"
                label="渠道名称"
                rules={[{ required: true, message: '请输入渠道名称' }]}
              >
                <Input placeholder="请输入渠道名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="signName"
                label="短信签名"
                rules={[{ required: true, message: '请输入短信签名' }]}
              >
                <Input placeholder="请输入短信签名" />
              </Form.Item>
            </Col>
          </Row>
          
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="appId"
                label="AppID"
                rules={[{ required: true, message: '请输入AppID' }]}
              >
                <Input placeholder="请输入AppID" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="appKey"
                label="AppKey"
                rules={[{ required: true, message: '请输入AppKey' }]}
              >
                <Input.Password placeholder="请输入AppKey" />
              </Form.Item>
            </Col>
          </Row>
          
          <Form.Item
            name="status"
            label="状态"
            valuePropName="checked"
            initialValue={true}
          >
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
          
          <Form.Item
            name="isDefault"
            label="默认渠道"
            valuePropName="checked"
            initialValue={false}
          >
            <Switch checkedChildren="是" unCheckedChildren="否" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 模板配置模态框 */}
      <Modal
        title={currentTemplate ? '编辑模板' : '添加模板'}
        open={templateModalVisible}
        onOk={handleTemplateSubmit}
        onCancel={() => setTemplateModalVisible(false)}
        maskClosable={false}
        width={600}
      >
        <Form
          form={templateForm}
          layout="vertical"
          className="sms-form"
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="name"
                label="模板名称"
                rules={[{ required: true, message: '请输入模板名称' }]}
              >
                <Input placeholder="请输入模板名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="code"
                label="模板CODE"
                rules={[{ required: true, message: '请输入模板CODE' }]}
              >
                <Input placeholder="请输入模板CODE" />
              </Form.Item>
            </Col>
          </Row>
          
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="type"
                label="模板类型"
                rules={[{ required: true, message: '请选择模板类型' }]}
              >
                <Select placeholder="请选择模板类型">
                  <Option value="验证码">验证码</Option>
                  <Option value="通知">通知</Option>
                  <Option value="推广">推广</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="channelId"
                label="所属渠道"
                rules={[{ required: true, message: '请选择所属渠道' }]}
              >
                <Select placeholder="请选择所属渠道">
                  {channels.map(channel => (
                    <Option key={channel.id} value={channel.id}>{channel.name}</Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>
          
          <Form.Item
            name="content"
            label="模板内容"
            rules={[{ required: true, message: '请输入模板内容' }]}
          >
            <TextArea rows={4} placeholder="请输入模板内容，变量请使用${变量名}格式" />
          </Form.Item>
          
          <Form.Item
            name="status"
            label="状态"
            valuePropName="checked"
            initialValue={true}
          >
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default SmsSettings; 