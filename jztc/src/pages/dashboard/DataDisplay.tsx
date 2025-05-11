import React from 'react';
import { Typography, Row, Col, Table, Tag, Button } from 'antd';
import {
  CheckOutlined,
  ClockCircleOutlined,
  CloseOutlined,
  SearchOutlined,
  GiftOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const { Title } = Typography;

// 定义订单数据类型
interface OrderData {
  id: number;
  orderNo: string;
  customer: string;
  amount: number;
  status: 'pending' | 'paid' | 'shipped' | 'completed' | 'cancelled';
  paymentMethod: string;
  createdAt: string;
  product: string;
}

// 定义兑换数据类型
interface ExchangeData {
  id: number;
  exchangeNo: string;
  userName: string;
  productName: string;
  points: number;
  status: 'pending' | 'processing' | 'completed' | 'cancelled';
  createdAt: string;
}

interface DataDisplayProps {
  orders: OrderData[];
  exchanges: ExchangeData[];
  loading: boolean;
}

const DataDisplay: React.FC<DataDisplayProps> = ({ orders, exchanges, loading }) => {
  const navigate = useNavigate();

  // 跳转到订单管理页面
  const navigateToOrders = () => {
    navigate('/orders');
  };

  // 跳转到兑换管理页面
  const navigateToExchanges = () => {
    navigate('/exchange/list');
  };

  // 获取状态标签
  const getStatusTag = (status: string) => {
    const statusConfig: Record<string, {color: string, text: string, icon: React.ReactNode}> = {
      'pending': { color: 'gold', text: '待处理', icon: <ClockCircleOutlined /> },
      'paid': { color: 'blue', text: '已支付', icon: <CheckOutlined /> },
      'shipped': { color: 'cyan', text: '已发货', icon: <CheckOutlined /> },
      'completed': { color: 'green', text: '已完成', icon: <CheckOutlined /> },
      'cancelled': { color: 'red', text: '已取消', icon: <CloseOutlined /> },
      'processing': { color: 'processing', text: '处理中', icon: <ClockCircleOutlined /> }
    };
    
    const config = statusConfig[status] || { color: 'default', text: status, icon: null };
    return <Tag color={config.color} icon={config.icon}>{config.text}</Tag>;
  };

  // 订单表格列定义
  const orderColumns = [
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (text: string) => (
        <span title={text}>{text.split(' ')[0]}</span>
      ),
    },
    {
      title: '客户',
      dataIndex: 'customer',
      key: 'customer',
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => `¥${amount.toFixed(2)}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
      filters: [
        { text: '待处理', value: 'pending' },
        { text: '已支付', value: 'paid' },
        { text: '已发货', value: 'shipped' },
        { text: '已完成', value: 'completed' },
        { text: '已取消', value: 'cancelled' },
      ],
      onFilter: (value: any, record: OrderData) => record.status === value,
    }
  ];

  // 兑换表格列定义
  const exchangeColumns = [
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (text: string) => (
        <span title={text}>{text.split(' ')[0]}</span>
      ),
    },
    {
      title: '用户',
      dataIndex: 'userName',
      key: 'userName',
    },
    {
      title: '商品',
      dataIndex: 'productName',
      key: 'productName',
      ellipsis: true,
    },
    {
      title: '时长',
      dataIndex: 'points',
      key: 'points',
      render: (points: number) => `${points}天`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
      filters: [
        { text: '待处理', value: 'pending' },
        { text: '处理中', value: 'processing' },
        { text: '已完成', value: 'completed' },
        { text: '已取消', value: 'cancelled' },
      ],
      onFilter: (value: any, record: ExchangeData) => record.status === value,
    }
  ];

  return (
    <>
      <div className="dashboard-title data-display-title">
        <Title level={4}>数据展示</Title>
      </div>
      
      {/* 平分布局的订单和兑换列表 */}
      <Row gutter={[16, 16]} className="dashboard-section">
        {/* 订单列表 */}
        <Col xs={24} lg={12}>
          <div className="list-subtitle-wrapper">
            <div className="list-subtitle">
              <span className="subtitle-icon"><SearchOutlined /></span>
              <span className="subtitle-text">订单列表</span>
            </div>
            <Button 
              type="link" 
              onClick={navigateToOrders}
            >
              查看更多
            </Button>
          </div>
          <Table 
            columns={orderColumns} 
            dataSource={orders.map(item => ({ ...item, key: item.id }))} 
            loading={loading}
            pagination={{ pageSize: 5, size: 'small' }}
            className="order-table"
            size="small"
            bordered={false}
          />
        </Col>
        
        {/* 兑换列表 */}
        <Col xs={24} lg={12}>
          <div className="list-subtitle-wrapper">
            <div className="list-subtitle">
              <span className="subtitle-icon"><GiftOutlined /></span>
              <span className="subtitle-text">兑换列表</span>
            </div>
            <Button 
              type="link" 
              onClick={navigateToExchanges}
            >
              查看更多
            </Button>
          </div>
          <Table 
            columns={exchangeColumns} 
            dataSource={exchanges.map(item => ({ ...item, key: item.id }))} 
            loading={loading}
            pagination={{ pageSize: 5, size: 'small' }}
            className="exchange-table"
            size="small"
            bordered={false}
          />
        </Col>
      </Row>
    </>
  );
};

export default DataDisplay; 