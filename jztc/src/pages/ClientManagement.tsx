import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Input,
  Modal,
  Form,
  Select,
  message,
  Popconfirm,
  Tag,
  Typography,
  Row,
  Col,
} from 'antd';
import {
  UserOutlined,
  SearchOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getClientList,
  createClient,
  updateClient,
  deleteClient,
  type ClientListItem,
  type ClientCreateReq,
  type ClientUpdateReq,
} from '@/api/client';
import './ClientManagement.css';

const { Title } = Typography;
const { Option } = Select;
const { Search } = Input;

const ClientManagement: React.FC = () => {
  const [clients, setClients] = useState<ClientListItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');
  
  // 模态框状态
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalTitle, setModalTitle] = useState('添加客户');
  const [editingClient, setEditingClient] = useState<ClientListItem | null>(null);
  
  // 表单引用
  const [form] = Form.useForm();

  // 获取客户列表数据
  const fetchClients = async () => {
    setLoading(true);
    try {
      // 只有在有搜索条件或不是第一页时才传递参数
      const params: any = {};
      if (searchText) {
        params.username = searchText; // 只提交username参数
      }
      if (current > 1) {
        params.page = current;
        params.pageSize = pageSize;
      }

      const res = await getClientList(params);
      if (res.code === 0) {
        setClients(res.data.list);
        setTotal(res.data.total);
      } else {
        message.error(res.message || '获取客户列表失败');
      }
    } catch (error) {
      console.error('获取客户列表失败:', error);
      message.error('获取客户列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 首次加载和依赖项变化时获取数据
  useEffect(() => {
    fetchClients();
  }, [current, pageSize, searchText]);

  // 处理搜索
  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrent(1); // 重置到第一页
  };

  // 处理添加客户
  const handleAdd = () => {
    setModalTitle('添加客户');
    setEditingClient(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  // 处理编辑客户
  const handleEdit = (record: ClientListItem) => {
    setModalTitle('编辑客户');
    setEditingClient(record);
    form.setFieldsValue({
      ...record,
      status: record.status === 1 ? '正常' : '禁用',
    });
    setIsModalVisible(true);
  };

  // 处理删除客户
  const handleDelete = async (id: number) => {
    try {
      await deleteClient(id);
      message.success('删除客户成功');
      fetchClients(); // 重新加载列表
    } catch (error) {
      console.error('删除客户失败:', error);
      message.error('删除客户失败');
    }
  };

  // 处理表单提交
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      
      if (editingClient) {
        // 更新客户
        const updateData: ClientUpdateReq = {
          id: editingClient.id,
          username: values.username,
          realName: values.realName,
          phone: values.phone,
          status: values.status === '正常' ? 1 : 0,
        };
        await updateClient(updateData);
        message.success('更新客户成功');
      } else {
        // 创建客户
        const createData: ClientCreateReq = {
          username: values.username,
          password: values.password,
          realName: values.realName,
          phone: values.phone,
          status: values.status === '正常' ? 1 : 0,
          identifier: values.identifier === '小程序' ? 'wxapp' : 'unknown',
        };
        await createClient(createData);
        message.success('添加客户成功');
      }
      
      setIsModalVisible(false);
      fetchClients(); // 重新加载列表
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  // 处理模态框取消
  const handleCancel = () => {
    setIsModalVisible(false);
  };

  // 渲染客户状态标签
  const renderStatusTag = (status: number) => {
    return status === 1 ? (
      <Tag color="green">正常</Tag>
    ) : (
      <Tag color="red">禁用</Tag>
    );
  };

  // 渲染标识标签
  const renderIdentifierTag = (identifier: string) => {
    let color = identifier === 'wxapp' ? 'cyan' : 'red';
    return <Tag color={color}>{identifier === 'wxapp' ? '小程序' : '未知'}</Tag>;
  };

  // 表格列定义
  const columns: ColumnsType<ClientListItem> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '姓名',
      dataIndex: 'realName',
      key: 'realName',
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: renderStatusTag,
    },
    {
      title: '来源',
      dataIndex: 'identifier',
      key: 'identifier',
      render: renderIdentifierTag,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="primary"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个客户吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
            icon={<ExclamationCircleOutlined style={{ color: 'red' }} />}
          >
            <Button 
              danger 
              size="small" 
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="client-management-container">
      <Card className="client-card">
        <div className="client-header">
          <Title level={4}>客户管理</Title>
          <Space size="large">
            <Search
              placeholder="输入用户名"
              onSearch={handleSearch}
              style={{ width: 250 }}
              allowClear
              enterButton={<><SearchOutlined />搜索</>}
            />
            <Button 
              type="primary" 
              icon={<PlusOutlined />} 
              onClick={handleAdd}
            >
              添加客户
            </Button>
            <Button 
              icon={<ReloadOutlined />} 
              onClick={fetchClients}
              loading={loading}
            >
              刷新
            </Button>
          </Space>
        </div>

        <Table
          dataSource={clients}
          columns={columns}
          rowKey="id"
          loading={loading}
          pagination={{
            defaultPageSize: 10,
            showSizeChanger: true,
            showTotal: (total: number) => `共 ${total} 条记录`,
            current,
            pageSize,
            total,
            onChange: (page, size) => {
              setCurrent(page);
              setPageSize(size);
            },
          }}
        />
      </Card>

      <Modal
        title={modalTitle}
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={handleCancel}
        maskClosable={false}
        destroyOnClose
        okText="确定"
        cancelText="取消"
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{ status: '正常', identifier: '小程序' }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="username"
                label="用户名"
                rules={[
                  { required: true, message: '请输入用户名' },
                  { min: 3, max: 20, message: '用户名长度在3-20个字符之间' }
                ]}
              >
                <Input prefix={<UserOutlined />} placeholder="请输入用户名" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="realName"
                label="姓名"
                rules={[{ required: true, message: '请输入姓名' }]}
              >
                <Input placeholder="请输入姓名" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="phone"
                label="手机号"
                rules={[
                  { required: true, message: '请输入手机号' },
                  {
                    pattern: /^1[3-9]\d{9}$/,
                    message: '请输入有效的手机号',
                  },
                ]}
              >
                <Input placeholder="请输入手机号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select placeholder="请选择状态">
                  <Option value="正常">正常</Option>
                  <Option value="禁用">禁用</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="identifier" label="来源">
                <Input disabled defaultValue="小程序" />
              </Form.Item>
            </Col>
            {!editingClient && (
              <Col span={12}>
                <Form.Item
                  name="password"
                  label="密码"
                  rules={[
                    { required: true, message: '请输入密码' },
                    { min: 6, max: 20, message: '密码长度在6-20个字符之间' }
                  ]}
                >
                  <Input.Password placeholder="请输入密码" />
                </Form.Item>
              </Col>
            )}
          </Row>
        </Form>
      </Modal>
    </div>
  );
};

export default ClientManagement; 