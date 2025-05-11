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
  Form,
  Modal,
  Select,
  DatePicker,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  SearchOutlined,
  ReloadOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import './ClientManagement.css'; // 引入客户管理页面的CSS
import './Orders.css'; // 引入订单管理页面的CSS
import { 
  getOrderList, 
  getOrderDetail, 
  cancelOrder, 
  deleteOrder, 
  updateOrderStatus,
  OrderListItem, 
  OrderListParams, 
  OrderDetailResponse,
  OrderUpdateStatusParams
} from '../api/order';

const { Title } = Typography;
const { Search } = Input;
const { RangePicker } = DatePicker;

const Orders: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [orders, setOrders] = useState<OrderListItem[]>([]);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [isPaging, setIsPaging] = useState(false); // 标记是否翻页操作
  
  // 搜索参数 - 初始不包含分页参数
  const [searchParams, setSearchParams] = useState<OrderListParams>({});
  
  // 模态框状态
  const [viewModalVisible, setViewModalVisible] = useState(false);
  const [currentOrder, setCurrentOrder] = useState<OrderDetailResponse | null>(null);
  
  // 更新订单状态模态框状态
  const [updateStatusModalVisible, setUpdateStatusModalVisible] = useState(false);
  const [updatingOrderNo, setUpdatingOrderNo] = useState('');
  
  // 表单实例
  const [form] = Form.useForm();
  const [updateStatusForm] = Form.useForm();

  useEffect(() => {
    fetchOrders();
  }, [current, pageSize, searchParams]);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      // 准备请求参数
      const params: OrderListParams = {
        ...searchParams
      };
      
      // 只有翻页操作时才添加分页参数
      if (isPaging) {
        params.page = current;
        params.pageSize = pageSize;
      }
      
      // 调用API获取订单列表
      const response = await getOrderList(params);
      
      // 处理嵌套在data中的数据
      const responseData = response.data || response;
      
      setOrders(responseData.list || []);
      setTotal(responseData.total || 0);
    } catch (error) {
      console.error('获取订单列表失败:', error);
      message.error('获取订单列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    // 搜索操作不是翻页
    setIsPaging(false);
    // 根据搜索条件更新搜索参数
    setSearchParams(prev => ({
      ...prev,
      orderNo: value,
    }));
    setCurrent(1); // 重置当前页码
  };

  const handleDateRangeChange = (dates: any) => {
    // 日期筛选不是翻页
    setIsPaging(false);
    if (dates && dates.length === 2) {
      setSearchParams(prev => ({
        ...prev,
        startTime: dates[0].format('YYYY-MM-DD'),
        endTime: dates[1].format('YYYY-MM-DD'),
      }));
    } else {
      setSearchParams(prev => ({
        ...prev,
        startTime: undefined,
        endTime: undefined,
      }));
    }
    setCurrent(1); // 重置当前页码
  };

  const handleStatusChange = (value: string) => {
    // 状态筛选不是翻页
    setIsPaging(false);
    setSearchParams(prev => ({
      ...prev,
      status: value,
    }));
    setCurrent(1); // 重置当前页码
  };

  const handleReset = () => {
    // 重置操作不是翻页
    setIsPaging(false);
    setSearchParams({});
    setCurrent(1);
    form.resetFields();
  };

  const getStatusTag = (status: number, statusText: string) => {
    switch (status) {
      case 0:
        return <Tag color="blue">{statusText}</Tag>;
      case 1:
        return <Tag color="green">{statusText}</Tag>;
      case 2:
        return <Tag color="red">{statusText}</Tag>;
      case 3:
        return <Tag color="orange">{statusText}</Tag>;
      case 4:
        return <Tag color="purple">{statusText}</Tag>;
      case 5:
        return <Tag color="cyan">{statusText}</Tag>;
      default:
        return <Tag>{statusText}</Tag>;
    }
  };

  const handleViewOrder = async (orderNo: string) => {
    try {
      setLoading(true);
      // 调用API获取订单详情
      const response = await getOrderDetail(orderNo);
      const orderDetail = response.data || response;
      
      setCurrentOrder(orderDetail);
      setViewModalVisible(true);
    } catch (error) {
      console.error('获取订单详情失败:', error);
      message.error('获取订单详情失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCancelOrder = async (orderNo: string) => {
    Modal.confirm({
      title: '确认取消订单',
      content: `确定要取消订单 ${orderNo} 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          setLoading(true);
          // 调用API取消订单
          await cancelOrder(orderNo);
          message.success('订单已取消');
          // 重新加载订单列表
          fetchOrders();
        } catch (error) {
          console.error('取消订单失败:', error);
          message.error('取消订单失败');
        } finally {
          setLoading(false);
        }
      }
    });
  };

  const handleDeleteOrder = async (orderNo: string) => {
    Modal.confirm({
      title: '确认删除订单',
      content: `确定要删除订单 ${orderNo} 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      okButtonProps: { danger: true },
      onOk: async () => {
        try {
          setLoading(true);
          // 调用API删除订单
          await deleteOrder(orderNo);
          message.success('订单已删除');
          // 重新加载订单列表
          fetchOrders();
        } catch (error) {
          console.error('删除订单失败:', error);
          message.error('删除订单失败');
        } finally {
          setLoading(false);
        }
      }
    });
  };

  // 打开更新订单状态模态框
  const handleOpenUpdateStatus = (record: OrderListItem) => {
    setUpdatingOrderNo(record.orderNo);
    // 设置表单初始值
    updateStatusForm.setFieldsValue({
      status: record.status,
      paymentMethod: record.paymentMethod,
      transactionId: '',
      remark: ''
    });
    setUpdateStatusModalVisible(true);
  };

  // 提交更新订单状态
  const handleUpdateStatus = async () => {
    try {
      const values = await updateStatusForm.validateFields();
      setLoading(true);
      
      const params: OrderUpdateStatusParams = {
        orderNo: updatingOrderNo,
        status: values.status,
        paymentMethod: values.paymentMethod,
        transactionId: values.transactionId,
        remark: values.remark
      };
      
      // 调用API更新订单状态
      await updateOrderStatus(params);
      message.success('订单状态已更新');
      
      // 关闭模态框并重新加载订单列表
      setUpdateStatusModalVisible(false);
      fetchOrders();
    } catch (error) {
      if (error instanceof Error) {
        console.error('更新订单状态失败:', error);
        message.error('更新订单状态失败: ' + error.message);
      } else {
        message.error('更新订单状态失败');
      }
    } finally {
      setLoading(false);
    }
  };

  const columns: ColumnsType<OrderListItem> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
      sorter: (a, b) => a.id - b.id,
    },
    {
      title: '订单号',
      dataIndex: 'orderNo',
      key: 'orderNo',
    },
    {
      title: '客户',
      dataIndex: 'clientName',
      key: 'clientName',
    },
    {
      title: '内容ID',
      dataIndex: 'contentId',
      key: 'contentId',
      width: 100,
    },
    {
      title: '商品',
      dataIndex: 'productName',
      key: 'productName',
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => `¥${amount.toFixed(2)}`,
      sorter: (a, b) => a.amount - b.amount,
    },
    {
      title: '支付方式',
      dataIndex: 'paymentMethod',
      key: 'paymentMethod',
    },
    {
      title: '订单状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number, record: OrderListItem) => getStatusTag(status, record.statusText),
      filters: [
        { text: '待支付', value: '0' },
        { text: '已支付', value: '1' },
        { text: '已取消', value: '2' },
        { text: '已退款', value: '3' },
        { text: '进行中', value: '4' },
        { text: '已完成', value: '5' },
      ],
      onFilter: (value, record) => record.status.toString() === value.toString(),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      sorter: (a, b) => {
        return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
      },
    },
    {
      title: '支付时间',
      dataIndex: 'payTime',
      key: 'payTime',
      render: (payTime: string) => payTime || '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: OrderListItem) => (
        <Space size="small">
          <Button
            type="primary"
            icon={<EyeOutlined />}
            size="small"
            onClick={() => handleViewOrder(record.orderNo)}
          >
            查看
          </Button>
          {record.status === 0 && (
            <Button
              danger
              size="small"
              onClick={() => handleCancelOrder(record.orderNo)}
            >
              取消
            </Button>
          )}
          <Button
            type="dashed"
            size="small"
            onClick={() => handleOpenUpdateStatus(record)}
          >
            更新状态
          </Button>
          <Button
            danger
            type="primary"
            size="small"
            onClick={() => handleDeleteOrder(record.orderNo)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  // 处理表格分页变化
  const handleTableChange = (page: number, pageSize?: number) => {
    setIsPaging(true); // 标记为翻页操作
    setCurrent(page);
    if (pageSize) {
      setPageSize(pageSize);
    }
  };

  return (
    <div className="client-management-container">
      <Card className="client-card">
        <div className="client-header">
          <Title level={4}>订单管理</Title>
          <Space size="large">
            <Search
              placeholder="输入订单号搜索"
              allowClear
              enterButton={<><SearchOutlined />搜索</>}
              onSearch={handleSearch}
              style={{ width: 250 }}
            />
            <Button 
              icon={<ReloadOutlined />} 
              onClick={handleReset}
              loading={loading}
            >
              刷新
            </Button>
          </Space>
        </div>
        
        {/* 筛选区域 */}
        <div className="filter-section" style={{ marginBottom: '16px' }}>
          <Row gutter={16} align="middle">
            <Col span={8}>
              <RangePicker 
                placeholder={['开始日期', '结束日期']}
                onChange={handleDateRangeChange}
                style={{ width: '100%' }}
              />
            </Col>
            <Col span={4}>
              <Select
                placeholder="订单状态"
                style={{ width: '100%' }}
                allowClear
                onChange={handleStatusChange}
              >
                <Select.Option value="0">待支付</Select.Option>
                <Select.Option value="1">已支付</Select.Option>
                <Select.Option value="2">已取消</Select.Option>
                <Select.Option value="3">已退款</Select.Option>
                <Select.Option value="4">进行中</Select.Option>
                <Select.Option value="5">已完成</Select.Option>
              </Select>
            </Col>
          </Row>
        </div>
        
        <Table
          dataSource={orders}
          columns={columns}
          rowKey="id"
          loading={loading}
          pagination={{
            current,
            pageSize,
            total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
            onChange: handleTableChange,
          }}
        />
      </Card>
      
      {/* 查看订单详情模态框 */}
      <Modal
        title="订单详情"
        open={viewModalVisible}
        onCancel={() => setViewModalVisible(false)}
        footer={[
          <Button key="back" onClick={() => setViewModalVisible(false)}>
            关闭
          </Button>,
        ]}
        width={700}
      >
        {currentOrder && (
          <div className="order-details">
            <Row gutter={[16, 16]}>
              <Col span={12}>
                <p><strong>订单号:</strong> {currentOrder.orderNo}</p>
                <p><strong>客户名称:</strong> {currentOrder.clientName}</p>
                <p><strong>内容ID:</strong> {currentOrder.contentId || '-'}</p>
                <p><strong>商品:</strong> {currentOrder.productName}</p>
                <p><strong>金额:</strong> ¥{currentOrder.amount?.toFixed(2)}</p>
                <p><strong>交易流水号:</strong> {currentOrder.transactionId || '-'}</p>
              </Col>
              <Col span={12}>
                <p><strong>支付方式:</strong> {currentOrder.paymentMethod || '-'}</p>
                <p><strong>订单状态:</strong> {currentOrder.status !== undefined && currentOrder.statusText && 
                  getStatusTag(currentOrder.status, currentOrder.statusText)}</p>
                <p><strong>创建时间:</strong> {currentOrder.createdAt}</p>
                <p><strong>支付时间:</strong> {currentOrder.payTime || '-'}</p>
                <p><strong>备注:</strong> {currentOrder.remark || '-'}</p>
              </Col>
            </Row>
          </div>
        )}
      </Modal>
      
      {/* 更新订单状态模态框 */}
      <Modal
        title="更新订单状态"
        open={updateStatusModalVisible}
        onCancel={() => setUpdateStatusModalVisible(false)}
        onOk={handleUpdateStatus}
        confirmLoading={loading}
        okText="确认"
        cancelText="取消"
      >
        <Form
          form={updateStatusForm}
          layout="vertical"
        >
          <Form.Item
            name="status"
            label="订单状态"
            rules={[{ required: true, message: '请选择订单状态' }]}
          >
            <Select>
              <Select.Option value={0}>待支付</Select.Option>
              <Select.Option value={1}>已支付</Select.Option>
              <Select.Option value={2}>已取消</Select.Option>
              <Select.Option value={3}>已退款</Select.Option>
              <Select.Option value={4}>进行中</Select.Option>
              <Select.Option value={5}>已完成</Select.Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="paymentMethod"
            label="支付方式"
          >
            <Select allowClear>
              <Select.Option value="wechat">微信</Select.Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="transactionId"
            label="交易流水号"
          >
            <Input placeholder="请输入交易流水号" />
          </Form.Item>
          
          <Form.Item
            name="remark"
            label="备注"
          >
            <Input.TextArea rows={4} placeholder="请输入备注信息" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Orders; 