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
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  LockOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import './UserManagement.css';
import { 
  getUserList, 
  createUser, 
  updateUser, 
  deleteUser,
  UserListItem,
  UserListParams,
  UserCreateParams,
  UserUpdateParams
} from '../api/user';

const { Title } = Typography;
const { Option } = Select;

const UserManagement: React.FC = () => {
  const [users, setUsers] = useState<UserListItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');
  
  // 模态框状态
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalTitle, setModalTitle] = useState('添加用户');
  const [editingUser, setEditingUser] = useState<UserListItem | null>(null);
  
  // 表单引用
  const [form] = Form.useForm();

  // 获取用户列表数据
  const fetchUsers = async () => {
    setLoading(true);
    try {
      // 准备API请求参数
      const params: UserListParams = {
        page: current,
        pageSize: pageSize
      };
      
      // 添加搜索条件
      if (searchText) {
        // 简单判断搜索内容类型
        if (searchText.includes('@')) {
          params.username = searchText; // 可能是邮箱
        } else {
          // 可能是用户名或昵称
          params.username = searchText;
          params.nickname = searchText;
        }
      }
      
      // 调用API获取数据
      const response = await getUserList(params);
      // 处理嵌套在data中的数据
      const responseData = response.data || response;
      
      setUsers(responseData.list || []);
      setTotal(responseData.total || 0);
    } catch (error) {
      console.error('获取用户列表失败:', error);
      message.error('获取用户列表失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 初始加载和搜索、分页变化时重新获取数据
  useEffect(() => {
    fetchUsers();
  }, [current, pageSize, searchText]);

  // 搜索处理
  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrent(1); // 重置到第一页
  };

  // 添加用户
  const handleAdd = () => {
    setModalTitle('添加用户');
    setEditingUser(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  // 编辑用户
  const handleEdit = (record: UserListItem) => {
    setModalTitle('编辑用户');
    setEditingUser(record);
    form.setFieldsValue({
      username: record.username,
      nickname: record.nickname,
      status: record.status
    });
    setIsModalVisible(true);
  };

  // 删除用户
  const handleDelete = async (id: number) => {
    try {
      setLoading(true);
      // 调用API删除用户
      const response = await deleteUser(id);
      
      message.success('删除成功');
      
      // 如果当前页没有数据了，且不是第一页，则回到上一页
      if (users.length === 1 && current > 1) {
        setCurrent(current - 1);
      } else {
        // 重新加载数据
        fetchUsers();
      }
    } catch (error) {
      console.error('删除用户失败:', error);
      message.error('删除用户失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      
      if (editingUser) {
        // 编辑现有用户
        const updateParams: UserUpdateParams = {
          id: editingUser.id,
          username: values.username,
          nickname: values.nickname,
          status: values.status
        };
        
        // 调用API更新用户
        const response = await updateUser(updateParams);
        message.success('用户信息更新成功');
      } else {
        // 添加新用户
        const createParams: UserCreateParams = {
          username: values.username,
          password: values.password,
          nickname: values.nickname
        };
        
        // 调用API创建用户
        const response = await createUser(createParams);
        message.success('用户添加成功');
      }
      
      // 关闭模态框并重置状态
      setIsModalVisible(false);
      setEditingUser(null);
      form.resetFields();
      
      // 重新加载数据
      fetchUsers();
    } catch (error) {
      console.error('表单提交失败:', error);
      message.error('操作失败，请检查表单并重试');
    } finally {
      setLoading(false);
    }
  };

  // 取消表单
  const handleCancel = () => {
    setIsModalVisible(false);
    setEditingUser(null);
    form.resetFields();
  };

  // 状态标签渲染
  const renderStatusTag = (status: number, statusText: string) => {
    return status === 1 ? (
      <Tag color="success">{statusText}</Tag>
    ) : (
      <Tag color="error">{statusText}</Tag>
    );
  };

  // 表格列定义
  const columns: ColumnsType<UserListItem> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number, record) => renderStatusTag(status, record.statusText),
      filters: [
        { text: '正常', value: 1 },
        { text: '禁用', value: 0 },
      ],
      onFilter: (value, record) => record.status === value,
    },
    {
      title: '最后登录IP',
      dataIndex: 'lastLoginIp',
      key: 'lastLoginIp',
    },
    {
      title: '最后登录时间',
      dataIndex: 'lastLoginTime',
      key: 'lastLoginTime',
    },
    {
      title: '操作',
      key: 'action',
      fixed: 'right',
      width: 200,
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
            title="确定要删除此用户吗？"
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
    <div className="user-management-container">
      <Card className="main-dashboard-card">
        <div className="user-header">
          <Title level={4}>用户管理</Title>
          <Space>
            <Input.Search
              placeholder="搜索用户名、昵称"
              allowClear
              enterButton
              onSearch={handleSearch}
              style={{ width: 300 }}
            />
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleAdd}
            >
              添加用户
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={users}
          rowKey="id"
          pagination={{
            current,
            pageSize,
            total,
            onChange: (page, pageSize) => {
              setCurrent(page);
              setPageSize(pageSize);
            },
            showTotal: (total) => `共 ${total} 条记录`,
            showSizeChanger: true,
            showQuickJumper: true,
          }}
          loading={loading}
          scroll={{ x: 1200 }}
        />
      </Card>

      <Modal
        title={modalTitle}
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={handleCancel}
        width={700}
        maskClosable={false}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{ status: 1 }}
        >
          <Row gutter={16}>
            <Col span={12}>
            <Form.Item
              name="username"
              label="用户名"
              rules={[
                { required: true, message: '请输入用户名' },
                { min: 3, max: 20, message: '用户名长度应为3-20个字符' },
                {
                  pattern: /^[a-zA-Z0-9_-]+$/,
                  message: '用户名只能包含字母、数字、下划线和连字符',
                },
              ]}
            >
              <Input prefix={<UserOutlined />} placeholder="请输入用户名" />
            </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="nickname"
                label="昵称"
                rules={[{ required: true, message: '请输入昵称' }]}
              >
                <Input placeholder="请输入昵称" />
              </Form.Item>
            </Col>
          </Row>
          
          {editingUser && (
            <Form.Item
              name="status"
              label="账户状态"
              rules={[{ required: true, message: '请选择账户状态' }]}
            >
              <Select>
                <Option value={1}>正常</Option>
                <Option value={0}>禁用</Option>
              </Select>
            </Form.Item>
          )}

          {!editingUser && (
            <Form.Item
              name="password"
              label="初始密码"
              rules={[
                { required: true, message: '请输入初始密码' },
                { min: 6, message: '密码长度不能少于6个字符' },
                {
                  validator(_, value) {
                    if (!value) return Promise.resolve();
                    const hasUpper = /[A-Z]/.test(value);
                    const hasLower = /[a-z]/.test(value);
                    const hasNumber = /[0-9]/.test(value);
                    
                    if (hasUpper && hasLower && hasNumber) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error('密码必须包含大小写字母和数字'));
                  },
                },
              ]}
            >
              <Input.Password prefix={<LockOutlined />} placeholder="请输入初始密码" />
            </Form.Item>
          )}
        </Form>
      </Modal>
    </div>
  );
};

export default UserManagement; 