import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Table,
  Button,
  Space,
  Input,
  message,
  Modal,
  Form,
  Popconfirm,
  Row,
  Col,
  Select,
} from 'antd';
import {
  ReloadOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getClientList,
  getClientDurationList,
  createClientDuration,
  updateClientDuration,
  deleteClientDuration,
  type ClientDurationItem,
  type ClientListItem,
  type ClientDurationCreateReq,
  type ClientDurationUpdateReq,
} from '@/api/client';
import '../ClientManagement.css'; // 引入客户管理页面的CSS
import './CustomerDuration.css'; // 引入独立的CSS文件

const { Title } = Typography;
const { Option } = Select;

// 分页数据
interface PaginationData {
  current: number;
  pageSize: number;
  total: number;
}

const CustomerDuration: React.FC = () => {
  const [form] = Form.useForm();
  const [durations, setDurations] = useState<ClientDurationItem[]>([]);
  const [clients, setClients] = useState<ClientListItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [modalVisible, setModalVisible] = useState<boolean>(false);
  const [modalTitle, setModalTitle] = useState<string>('');
  const [editingRecord, setEditingRecord] = useState<ClientDurationItem | null>(null);
  const [pagination, setPagination] = useState<PaginationData>({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [selectedClientId, setSelectedClientId] = useState<number | undefined>(undefined);

  // 页面加载时获取数据
  useEffect(() => {
    fetchData();
  }, [pagination.current, pagination.pageSize, selectedClientId]);

  // 首次加载获取客户列表
  useEffect(() => {
    fetchClients();
  }, []);

  // 获取客户列表，用于选择客户
  const fetchClients = async () => {
    try {
      const res = await getClientList({});
      if (res.code === 0) {
        setClients(res.data.list);
      } else {
        message.error(res.message || '获取客户列表失败');
      }
    } catch (error) {
      console.error('获取客户列表失败:', error);
      message.error('获取客户列表失败');
    }
  };

  // 获取客户时长数据
  const fetchData = async () => {
    setLoading(true);
    try {
      const params: any = {
        page: pagination.current,
        pageSize: pagination.pageSize,
      };
      
      if (selectedClientId) {
        params.clientId = selectedClientId;
      }

      const res = await getClientDurationList(params);
      if (res.code === 0) {
        setDurations(res.data.list);
        setPagination({
          ...pagination,
          total: res.data.total,
        });
      } else {
        message.error(res.message || '获取客户时长列表失败');
      }
    } catch (error) {
      console.error('获取客户时长列表失败:', error);
      message.error('获取客户时长列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 表格列定义
  const columns: ColumnsType<ClientDurationItem> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '客户ID',
      dataIndex: 'clientId',
      key: 'clientId',
      width: 120,
      className: 'customer-id-column',
    },
    {
      title: '客户名称',
      dataIndex: 'clientName',
      key: 'clientName',
      width: '15%',
    },
    {
      title: '总时长',
      dataIndex: 'totalDuration',
      key: 'totalDuration',
      width: '15%',
    },
    {
      title: '已使用时长',
      dataIndex: 'usedDuration',
      key: 'usedDuration',
      width: '15%',
    },
    {
      title: '剩余时长',
      dataIndex: 'remainingDuration',
      key: 'remainingDuration',
      width: '15%',
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: '15%',
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      key: 'updatedAt',
      width: '15%',
      render: (updatedAt) => updatedAt || '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      fixed: 'right',
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
            title="确定要删除吗?"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
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

  // 处理表格分页变化
  const handleTableChange = (pagination: any) => {
    setPagination({
      ...pagination,
      current: pagination.current,
      pageSize: pagination.pageSize,
    });
  };

  // 处理客户筛选
  const handleClientChange = (value: number | undefined) => {
    setSelectedClientId(value);
    setPagination({ ...pagination, current: 1 });
  };

  // 打开添加客户时长记录的模态框
  const handleAdd = () => {
    setModalTitle('添加客户时长');
    setEditingRecord(null);
    form.resetFields();
    setModalVisible(true);
  };

  // 打开编辑客户时长记录的模态框
  const handleEdit = (record: ClientDurationItem) => {
    setModalTitle('编辑客户时长');
    setEditingRecord(record);
    form.setFieldsValue({
      clientId: record.clientId,
      totalDuration: record.totalDuration,
      remainingDuration: record.remainingDuration,
      usedDuration: record.usedDuration,
    });
    setModalVisible(true);
  };

  // 处理删除客户时长记录
  const handleDelete = async (id: number) => {
    try {
      const res = await deleteClientDuration(id);
      if (res.code === 0) {
        message.success('删除成功');
        fetchData();
      } else {
        message.error(res.message || '删除失败');
      }
    } catch (error) {
      console.error('删除客户时长失败:', error);
      message.error('删除失败');
    }
  };

  // 处理刷新
  const handleRefresh = () => {
    fetchData();
  };

  // 处理模态框的确认
  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      
      // 获取选中的客户名称
      const selectedClient = clients.find(client => client.id === values.clientId);
      const clientName = selectedClient ? selectedClient.realName : '';
      
      if (editingRecord) {
        // 更新客户时长
        const updateData: ClientDurationUpdateReq = {
          id: editingRecord.id,
          clientId: values.clientId,
          clientName,
          totalDuration: values.totalDuration,
          remainingDuration: values.remainingDuration,
          usedDuration: values.usedDuration,
        };
        
        const res = await updateClientDuration(updateData);
        if (res.code === 0) {
          message.success('编辑成功');
          setModalVisible(false);
          fetchData();
        } else {
          message.error(res.message || '编辑失败');
        }
      } else {
        // 创建客户时长
        const createData: ClientDurationCreateReq = {
          clientId: values.clientId,
          clientName,
          totalDuration: values.totalDuration,
          remainingDuration: values.remainingDuration,
          usedDuration: values.usedDuration,
        };
        
        const res = await createClientDuration(createData);
        if (res.code === 0) {
          message.success('添加成功');
          setModalVisible(false);
          fetchData();
        } else {
          message.error(res.message || '添加失败');
        }
      }
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  // 渲染页面
  return (
    <div className="client-management-container">
      <Card className="client-card">
        <div className="client-header">
          <Title level={4}>客户时长</Title>
          <Space size="large">
            <Select
              placeholder="选择客户"
              style={{ width: 200 }}
              allowClear
              onChange={handleClientChange}
            >
              {clients.map(client => (
                <Option key={client.id} value={client.id}>
                  {client.realName}
                </Option>
              ))}
            </Select>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleAdd}
            >
              添加时长
            </Button>
            <Button
              icon={<ReloadOutlined />}
              onClick={handleRefresh}
              loading={loading}
            >
              刷新
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={durations}
          rowKey="id"
          pagination={{
            ...pagination,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: total => `共 ${total} 条记录`,
          }}
          scroll={{ x: true }}
          loading={loading}
          onChange={handleTableChange}
          className="customer-duration-table"
        />

        <Modal
          title={modalTitle}
          open={modalVisible}
          onOk={handleModalOk}
          onCancel={() => setModalVisible(false)}
          okText="保存"
          cancelText="取消"
          width={700}
        >
          <Form
            form={form}
            layout="vertical"
          >
            <Row gutter={16}>
              <Col span={24}>
                <Form.Item
                  name="clientId"
                  label="客户"
                  rules={[{ required: true, message: '请选择客户' }]}
                >
                  <Select placeholder="请选择客户">
                    {clients.map(client => (
                      <Option key={client.id} value={client.id}>
                        {client.realName}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>
              </Col>
            </Row>
            
            <Row gutter={16}>
              <Col span={8}>
                <Form.Item
                  name="totalDuration"
                  label="总时长"
                  rules={[{ required: true, message: '请输入总时长' }]}
                >
                  <Input prefix={<ClockCircleOutlined />} placeholder="例如：5天10小时30分钟" />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  name="usedDuration"
                  label="已使用时长"
                  rules={[{ required: true, message: '请输入已使用时长' }]}
                >
                  <Input prefix={<ClockCircleOutlined />} placeholder="例如：1天15小时48分钟" />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  name="remainingDuration"
                  label="剩余时长"
                  rules={[{ required: true, message: '请输入剩余时长' }]}
                >
                  <Input prefix={<ClockCircleOutlined />} placeholder="例如：3天18小时42分钟" />
                </Form.Item>
              </Col>
            </Row>
          </Form>
        </Modal>
      </Card>
    </div>
  );
};

export default CustomerDuration; 