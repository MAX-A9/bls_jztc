import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Table,
  Space,
  Button,
  Tag,
  Input,
  Row,
  Col,
  message,
  Popconfirm,
  Modal,
  Form,
  Select,
  DatePicker,
  InputNumber,
} from 'antd';
import {
  SearchOutlined,
  ReloadOutlined,
  EyeOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  SyncOutlined,
  PlusOutlined,
  UserOutlined,
} from '@ant-design/icons';
import '../ClientManagement.css';
import { 
  getExchangeRecordList, 
  getExchangeRecord, 
  deleteExchangeRecord, 
  updateExchangeRecordStatus,
  createExchangeRecord,
  type ExchangeRecord,
  type ExchangeRecordListParams,
  type ExchangeRecordStatus,
  type CreateExchangeRecordReq
} from '../../api/exchange';
import { getClientList } from '../../api/client';
import dayjs from 'dayjs';

const { Title } = Typography;
const { Option } = Select;

// 前端UI使用的兑换记录接口（与API接口有区别，保留原有属性以兼容UI）
interface UIExchangeRecord {
  id: number;
  userId: string; // 客户ID
  userName: string; // 客户名称
  productId?: number;
  productName: string;
  quantity?: number;
  phone: string; // 充值账号
  exchangeTime: string;
  status: ExchangeRecordStatus;
  address?: string;
  remark?: string;
  duration?: number; // 兑换所需时长
}

const ExchangeList: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [records, setRecords] = useState<UIExchangeRecord[]>([]);
  const [filteredRecords, setFilteredRecords] = useState<UIExchangeRecord[]>([]);
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentRecord, setCurrentRecord] = useState<UIExchangeRecord | null>(null);
  const [searchText, setSearchText] = useState('');
  // 分页相关状态
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  // 筛选状态
  const [statusFilter, setStatusFilter] = useState<ExchangeRecordStatus | ''>('');
  // 创建记录相关状态
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [createForm] = Form.useForm();
  // 客户选择器相关状态
  const [clientSelectVisible, setClientSelectVisible] = useState(false);
  const [clientList, setClientList] = useState<Array<{id: number, realName: string, phone: string}>>([]);
  const [clientLoading, setClientLoading] = useState(false);
  
  // 将API响应数据转换为UI数据
  const convertToUIRecord = (apiRecord: ExchangeRecord): UIExchangeRecord => {
    return {
      id: apiRecord.id,
      userId: String(apiRecord.clientId),
      userName: apiRecord.clientName,
      productName: apiRecord.productName,
      phone: apiRecord.rechargeAccount,
      exchangeTime: apiRecord.exchangeTime,
      status: apiRecord.status,
      remark: apiRecord.remark,
      duration: apiRecord.duration
    };
  };

  // 初始化获取数据
  useEffect(() => {
    fetchRecords();
  }, [currentPage, pageSize, statusFilter]);

  // 监听搜索条件变化
  useEffect(() => {
    if (searchText) {
      const filtered = records.filter(record =>
        record.userName.toLowerCase().includes(searchText.toLowerCase()) ||
        record.productName.toLowerCase().includes(searchText.toLowerCase()) ||
        record.userId.toLowerCase().includes(searchText.toLowerCase()) ||
        record.phone.includes(searchText)
      );
      setFilteredRecords(filtered);
    } else {
      setFilteredRecords(records);
    }
  }, [searchText, records]);

  // 获取兑换记录
  const fetchRecords = async () => {
    setLoading(true);
    try {
      const params: ExchangeRecordListParams = {
        page: currentPage,
        size: pageSize
      };
      
      // 如果有状态筛选，则添加状态参数
      if (statusFilter) {
        params.status = statusFilter;
      }
      
      const result = await getExchangeRecordList(params);
      
      if (result.code === 0) {
        // 将API返回的数据转换为UI使用的格式
        const uiRecords = result.data.list.map(convertToUIRecord);
        setRecords(uiRecords);
        setFilteredRecords(uiRecords);
        setTotal(result.data.total);
      } else {
        message.error(result.message || '获取兑换记录失败');
        setRecords([]);
        setFilteredRecords([]);
        setTotal(0);
      }
    } catch (error) {
      console.error('获取兑换记录失败:', error);
      message.error('获取兑换记录失败');
      setRecords([]);
      setFilteredRecords([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // 查看详情
  const showDetail = async (record: UIExchangeRecord) => {
    setLoading(true);
    try {
      const result = await getExchangeRecord(record.id);
      
      if (result.code === 0) {
        const detailRecord = convertToUIRecord(result.data);
        setCurrentRecord(detailRecord);
        setDetailVisible(true);
      } else {
        message.error(result.message || '获取兑换记录详情失败');
      }
    } catch (error) {
      console.error('获取兑换记录详情失败:', error);
      message.error('获取兑换记录详情失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 删除兑换记录
  const handleDelete = async (id: number) => {
    try {
      const result = await deleteExchangeRecord(id);
      
      if (result.code === 0) {
        message.success('删除成功');
        fetchRecords(); // 刷新列表
      } else {
        message.error(result.message || '删除失败');
      }
    } catch (error) {
      console.error('删除兑换记录失败:', error);
      message.error('删除失败');
    }
  };

  // 更新状态
  const updateStatus = async (id: number, status: ExchangeRecordStatus) => {
    try {
      const result = await updateExchangeRecordStatus({ id, status });
      
      if (result.code === 0) {
        message.success('状态已更新');
        fetchRecords(); // 刷新列表
      } else {
        message.error(result.message || '更新状态失败');
      }
    } catch (error) {
      console.error('更新状态失败:', error);
      message.error('更新状态失败');
    }
  };

  // 处理搜索
  const handleSearch = (value: string) => {
    setSearchText(value);
    // 重置为第一页
    setCurrentPage(1);
  };
  
  // 表格分页变化
  const handleTableChange = (pagination: any, filters: any) => {
    setCurrentPage(pagination.current);
    setPageSize(pagination.pageSize);
    
    // 处理状态筛选
    if (filters.status && filters.status.length > 0) {
      setStatusFilter(filters.status[0] as ExchangeRecordStatus);
    } else {
      setStatusFilter('');
    }
  };

  // 显示创建记录模态框
  const showCreateModal = () => {
    createForm.resetFields();
    // 设置初始值
    createForm.setFieldsValue({
      status: 'processing', // 默认为处理中
      exchangeTime: dayjs(), // 当前时间
    });
    setCreateModalVisible(true);
  };

  // 显示客户选择器
  const showClientSelector = () => {
    // 获取客户列表
    setClientLoading(true);
    getClientList({
      page: 1,
      pageSize: 100
    }).then(result => {
      if (result.code === 0 && result.data && result.data.list) {
        setClientList(result.data.list);
      } else {
        message.error(result.message || '获取客户列表失败');
        setClientList([]);
      }
    }).catch(error => {
      console.error('获取客户列表失败:', error);
      message.error('获取客户列表失败');
      setClientList([]);
    }).finally(() => {
      setClientLoading(false);
    });
    
    setClientSelectVisible(true);
  };
  
  // 选择客户
  const handleSelectClient = (client: {id: number, realName: string}) => {
    createForm.setFieldsValue({
      clientId: client.id,
      clientName: client.realName
    });
    
    setClientSelectVisible(false);
  };

  // 创建兑换记录
  const handleCreateRecord = async () => {
    try {
      const values = await createForm.validateFields();
      
      const recordData: CreateExchangeRecordReq = {
        clientId: parseInt(values.clientId),
        clientName: values.clientName,
        duration: values.duration,
        exchangeTime: values.exchangeTime.format('YYYY-MM-DD HH:mm:ss'),
        productName: values.productName,
        rechargeAccount: values.rechargeAccount,
        remark: values.remark || '',
        status: values.status,
      };
      
      const result = await createExchangeRecord(recordData);
      
      if (result.code === 0) {
        message.success('创建兑换记录成功');
        setCreateModalVisible(false);
        fetchRecords(); // 刷新列表
      } else {
        message.error(result.message || '创建兑换记录失败');
      }
    } catch (error) {
      console.error('创建兑换记录失败:', error);
      message.error('创建兑换记录失败');
    }
  };

  // 表格列定义
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 70,
    },
    {
      title: '客户',
      dataIndex: 'userName',
      key: 'userName',
      render: (text: string, record: UIExchangeRecord) => (
        <span>{text} (ID: {record.userId})</span>
      ),
    },
    {
      title: '充值账号',
      dataIndex: 'phone',
      key: 'phone',
      width: 130,
    },
    {
      title: '商品',
      dataIndex: 'productName',
      key: 'productName',
    },
    {
      title: '所需时长',
      dataIndex: 'duration',
      key: 'duration',
      width: 100,
      render: (duration?: number) => duration ? `${duration}天` : '-',
    },
    {
      title: '兑换时间',
      dataIndex: 'exchangeTime',
      key: 'exchangeTime',
      width: 180,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 120,
      render: (status: ExchangeRecordStatus) => {
        switch (status) {
          case 'processing':
            return <Tag icon={<SyncOutlined spin />} color="processing">处理中</Tag>;
          case 'completed':
            return <Tag icon={<CheckCircleOutlined />} color="success">已完成</Tag>;
          case 'failed':
            return <Tag icon={<CloseCircleOutlined />} color="error">失败</Tag>;
          default:
            return <Tag color="default">{status}</Tag>;
        }
      },
      filters: [
        { text: '处理中', value: 'processing' },
        { text: '已完成', value: 'completed' },
        { text: '失败', value: 'failed' },
      ],
      filteredValue: statusFilter ? [statusFilter] : null,
    },
    {
      title: '操作',
      key: 'action',
      width: 250,
      render: (_: any, record: UIExchangeRecord) => (
        <Space size="small">
          <Button 
            type="primary" 
            size="small" 
            icon={<EyeOutlined />}
            onClick={() => showDetail(record)}
          >
            详情
          </Button>
          {record.status === 'processing' && (
            <>
              <Button 
                size="small" 
                type="primary"
                onClick={() => updateStatus(record.id, 'completed')}
              >
                标记完成
              </Button>
              <Button 
                size="small" 
                danger
                onClick={() => updateStatus(record.id, 'failed')}
              >
                标记失败
              </Button>
            </>
          )}
          <Popconfirm
            title="确定要删除吗?"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button size="small" danger>
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
          <Title level={4}>兑换列表</Title>
          <Space size="large">
            <Input.Search
              placeholder="搜索客户/充值账号/商品"
              onSearch={handleSearch}
              style={{ width: 250 }}
              allowClear
              enterButton={<><SearchOutlined />搜索</>}
            />
            <Button 
              type="primary"
              icon={<PlusOutlined />}
              onClick={showCreateModal}
            >
              新增兑换
            </Button>
            <Button 
              icon={<ReloadOutlined />} 
              onClick={fetchRecords}
              loading={loading}
            >
              刷新
            </Button>
          </Space>
        </div>
        
        <Table
          dataSource={searchText ? filteredRecords : records}
          columns={columns}
          rowKey="id"
          loading={loading}
          onChange={handleTableChange}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条记录`,
          }}
        />
      </Card>
      
      {/* 详情模态框 */}
      <Modal
        title="兑换详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailVisible(false)}>
            关闭
          </Button>,
        ]}
      >
        {currentRecord && (
          <div>
            <p><strong>记录编号:</strong> {currentRecord.id}</p>
            <p><strong>客户:</strong> {currentRecord.userName} (ID: {currentRecord.userId})</p>
            <p><strong>充值账号:</strong> {currentRecord.phone}</p>
            <p><strong>商品:</strong> {currentRecord.productName}</p>
            {currentRecord.duration !== undefined && <p><strong>所需时长:</strong> {currentRecord.duration}天</p>}
            <p><strong>兑换时间:</strong> {currentRecord.exchangeTime}</p>
            <p><strong>状态:</strong> {
              currentRecord.status === 'processing' ? '处理中' : 
              currentRecord.status === 'completed' ? '已完成' : '失败'
            }</p>
            {currentRecord.address && <p><strong>收货地址:</strong> {currentRecord.address}</p>}
            {currentRecord.remark && <p><strong>备注:</strong> {currentRecord.remark}</p>}
          </div>
        )}
      </Modal>
      
      {/* 创建记录模态框 */}
      <Modal
        title="创建兑换记录"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onOk={handleCreateRecord}
        okText="保存"
        cancelText="取消"
        width={600}
      >
        <Form
          form={createForm}
          layout="vertical"
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="clientId"
                label="客户ID"
                rules={[{ required: true, message: '请输入客户ID' }]}
              >
                <Input 
                  placeholder="请输入客户ID" 
                  addonAfter={
                    <Button 
                      type="link" 
                      size="small" 
                      icon={<UserOutlined />} 
                      onClick={showClientSelector}
                      style={{ marginRight: -7, marginLeft: -7 }}
                    >
                      选择
                    </Button>
                  }
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="clientName"
                label="客户名称"
                rules={[{ required: true, message: '请输入客户名称' }]}
              >
                <Input placeholder="请输入客户名称" />
              </Form.Item>
            </Col>
          </Row>
          
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="productName"
                label="商品名称"
                rules={[{ required: true, message: '请输入商品名称' }]}
              >
                <Input placeholder="请输入商品名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="rechargeAccount"
                label="充值账号"
                rules={[{ required: true, message: '请输入充值账号' }]}
              >
                <Input placeholder="请输入充值账号" />
              </Form.Item>
            </Col>
          </Row>
          
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="duration"
                label="所需时长(天)"
                rules={[{ required: true, message: '请输入所需时长' }]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  step={0.5}
                  placeholder="请输入所需时长"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="exchangeTime"
                label="兑换时间"
                rules={[{ required: true, message: '请选择兑换时间' }]}
              >
                <DatePicker
                  style={{ width: '100%' }}
                  showTime
                  format="YYYY-MM-DD HH:mm:ss"
                />
              </Form.Item>
            </Col>
          </Row>
          
          <Form.Item
            name="status"
            label="状态"
            rules={[{ required: true, message: '请选择状态' }]}
          >
            <Select placeholder="请选择状态">
              <Option value="processing">处理中</Option>
              <Option value="completed">已完成</Option>
              <Option value="failed">失败</Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="remark"
            label="备注"
          >
            <Input.TextArea rows={3} placeholder="请输入备注" />
          </Form.Item>
        </Form>
      </Modal>
      
      {/* 客户选择器对话框 */}
      <Modal
        title="选择客户"
        open={clientSelectVisible}
        onCancel={() => setClientSelectVisible(false)}
        footer={null}
        width={700}
      >
        <Input.Search
          placeholder="搜索客户名称或手机号"
          allowClear
          enterButton="搜索"
          onSearch={(value) => {
            setClientLoading(true);
            getClientList({
              page: 1,
              pageSize: 100,
              ...(value ? (/^\d+$/.test(value) ? { phone: value } : { realName: value }) : {})
            }).then(result => {
              if (result.code === 0 && result.data && result.data.list) {
                setClientList(result.data.list);
              } else {
                setClientList([]);
              }
            }).catch(() => {
              setClientList([]);
            }).finally(() => {
              setClientLoading(false);
            });
          }}
          style={{ marginBottom: 16 }}
        />
        
        <Table
          dataSource={clientList}
          rowKey="id"
          loading={clientLoading}
          columns={[
            {
              title: 'ID',
              dataIndex: 'id',
              key: 'id',
              width: 80,
            },
            {
              title: '客户名称',
              dataIndex: 'realName',
              key: 'realName',
            },
            {
              title: '手机号',
              dataIndex: 'phone',
              key: 'phone',
            },
            {
              title: '操作',
              key: 'action',
              width: 100,
              render: (_, record) => (
                <Button 
                  type="primary" 
                  size="small"
                  onClick={() => handleSelectClient(record)}
                >
                  选择
                </Button>
              ),
            },
          ]}
          pagination={false}
          scroll={{ y: 300 }}
        />
      </Modal>
    </div>
  );
};

export default ExchangeList; 