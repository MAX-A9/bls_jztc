import React, { useState, useEffect } from 'react';
import { 
  Typography, 
  Card, 
  Row, 
  Col, 
  Segmented
} from 'antd';
import {
  UserAddOutlined,
  SwapOutlined,
  CloudUploadOutlined,
  MoneyCollectOutlined
} from '@ant-design/icons';
import { getStatsData, getTrendData } from '../api/stats';
import { getOrderList, OrderListItem } from '../api/order';
import { getExchangeRecordList, ExchangeRecord, ExchangeRecordStatus } from '../api/exchange';

// 使用更明确的路径
import TrendAnalysis from '../pages/dashboard/TrendAnalysis';
import DataDisplay from '../pages/dashboard/DataDisplay';
import './Dashboard.css';

const { Title } = Typography;

// 定义统计数据类型
interface StatisticsData {
  registration: number;
  exchange: number;
  publish: number;
  income: number;
}

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

// 定义趋势数据类型（组件内部使用）
interface TrendData {
  date: string;
  value: number;
  category: string;
}

// 将API返回的订单状态转换为组件中使用的状态
const mapOrderStatus = (status: number): 'pending' | 'paid' | 'shipped' | 'completed' | 'cancelled' => {
  switch (status) {
    case 0: return 'pending'; // 待支付
    case 1: return 'paid';    // 已支付
    case 4: return 'shipped'; // 进行中
    case 5: return 'completed'; // 已完成
    case 2: return 'cancelled'; // 已取消
    case 3: return 'cancelled'; // 已退款，映射为已取消
    default: return 'pending';
  }
};

// 将API返回的兑换状态转换为组件中使用的状态
const mapExchangeStatus = (status: ExchangeRecordStatus): 'pending' | 'processing' | 'completed' | 'cancelled' => {
  switch (status) {
    case 'processing': return 'processing';
    case 'completed': return 'completed';
    case 'failed': return 'cancelled';
    default: return 'pending';
  }
};

// 定义时间范围类型
type TimeRangeType = '本周' | '本月' | '本年';
// 对应API参数
const timeRangeApiMap: Record<TimeRangeType, 'week' | 'month' | 'year'> = {
  '本周': 'week',
  '本月': 'month',
  '本年': 'year'
};

const Dashboard: React.FC = () => {
  // 筛选类型状态
  const [timeRange, setTimeRange] = useState<TimeRangeType>('本周');
  
  // 统计数据状态
  const [statistics, setStatistics] = useState<StatisticsData>({
    registration: 0,
    exchange: 0,
    publish: 0,
    income: 0
  });

  // 订单数据状态
  const [orders, setOrders] = useState<OrderData[]>([]);
  // 兑换数据状态
  const [exchanges, setExchanges] = useState<ExchangeData[]>([]);
  // 加载状态
  const [loading, setLoading] = useState(false);

  // 趋势数据状态
  const [trendData, setTrendData] = useState<TrendData[]>([]);

  // 根据筛选类型更新数据
  useEffect(() => {
    // 加载对应时间范围的数据
    const loadData = async () => {
      setLoading(true);
      try {
        // 调用API获取统计数据
        const statsResponse = await getStatsData(timeRangeApiMap[timeRange]);
        
        // 将API返回的数据映射到统计数据中
        const statsData: StatisticsData = {
          registration: statsResponse.data.clientCount,
          exchange: statsResponse.data.exchangeCount,
          publish: statsResponse.data.publishCount,
          income: statsResponse.data.revenueAmount
        };
        
        setStatistics(statsData);
        
        // 获取趋势数据
        const trendResponse = await getTrendData(timeRangeApiMap[timeRange], 'all');
        
        // 将API返回的趋势数据转换为组件使用的格式
        const transformedTrendData: TrendData[] = [];
        
        // 处理所有趋势数据
        const categoryMap = {
          'clients': '注册用户',
          'exchanges': '兑换次数',
          'publishes': '发布内容',
          'revenue': '收益金额'
        };
        
        // 按类别处理趋势数据
        Object.entries(categoryMap).forEach(([key, category]) => {
          const dataKey = key as keyof typeof trendResponse.data.allValues;
          trendResponse.data.labels.forEach((label, index) => {
            let value = trendResponse.data.allValues[dataKey][index];
            if (typeof value === 'string') {
              value = parseFloat(value);
            }
            
            transformedTrendData.push({
              date: label,
              value: value as number,
              category: category
            });
          });
        });
        
        setTrendData(transformedTrendData);
        
        // 获取订单数据
        const orderListResponse = await getOrderList({
          page: 1,
          pageSize: 5
        });
        
        const orderResponseData = orderListResponse.data || orderListResponse;
        const orderList = orderResponseData.list || [];
        
        // 转换订单数据格式
        const transformedOrders: OrderData[] = orderList.map((order: OrderListItem) => ({
          id: order.id,
          orderNo: order.orderNo,
          customer: order.clientName,
          amount: order.amount,
          status: mapOrderStatus(order.status),
          paymentMethod: order.paymentMethod,
          createdAt: order.createdAt,
          product: order.productName
        }));
        
        setOrders(transformedOrders);
        
        // 获取兑换数据
        const exchangeResponse = await getExchangeRecordList({
          page: 1,
          size: 5
        });
        
        if (exchangeResponse.code === 0) {
          // 转换兑换数据格式
          const transformedExchanges: ExchangeData[] = exchangeResponse.data.list.map((exchange: ExchangeRecord) => ({
            id: exchange.id,
            exchangeNo: `EX${exchange.id.toString().padStart(10, '0')}`,
            userName: exchange.clientName,
            productName: exchange.productName,
            points: exchange.duration || 0,
            status: mapExchangeStatus(exchange.status),
            createdAt: exchange.exchangeTime
          }));
          
          setExchanges(transformedExchanges);
        }
        
      } catch (error) {
        console.error('获取数据失败', error);
      } finally {
        setLoading(false);
      }
    };
    
    loadData();
  }, [timeRange]);

  // 处理筛选变化
  const handleRangeChange = (value: string) => {
    setTimeRange(value as TimeRangeType);
  };

  return (
    <Card className="main-dashboard-card">
      <div className="dashboard-title">
        <Title level={4}>概览</Title>
        <Segmented
          options={['本周', '本月', '本年']}
          value={timeRange}
          onChange={value => handleRangeChange(value as string)}
        />
      </div>
      
      {/* 统计卡片区域 */}
      <Row gutter={[16, 16]} className="stats-row">
        <Col xs={24} sm={12} md={12} lg={6}>
          <Card variant="borderless" className="stat-card card-registration">
            <div className="stat-icon">
              <UserAddOutlined />
            </div>
            <div className="stat-content">
              <div className="stat-value">{statistics.registration}</div>
              <div className="stat-title">注册用户</div>
            </div>
          </Card>
        </Col>
        
        <Col xs={24} sm={12} md={12} lg={6}>
          <Card variant="borderless" className="stat-card card-exchange">
            <div className="stat-icon">
              <SwapOutlined />
            </div>
            <div className="stat-content">
              <div className="stat-value">{statistics.exchange}</div>
              <div className="stat-title">兑换次数</div>
            </div>
          </Card>
        </Col>
        
        <Col xs={24} sm={12} md={12} lg={6}>
          <Card variant="borderless" className="stat-card card-publish">
            <div className="stat-icon">
              <CloudUploadOutlined />
            </div>
            <div className="stat-content">
              <div className="stat-value">{statistics.publish}</div>
              <div className="stat-title">发布内容</div>
            </div>
          </Card>
        </Col>
        
        <Col xs={24} sm={12} md={12} lg={6}>
          <Card variant="borderless" className="stat-card card-income">
            <div className="stat-icon">
              <MoneyCollectOutlined />
            </div>
            <div className="stat-content">
              <div className="stat-value">{statistics.income}</div>
              <div className="stat-title">收益金额</div>
            </div>
          </Card>
        </Col>
      </Row>

      {/* 趋势分析区域 */}
      <TrendAnalysis 
        statistics={statistics}
        trendData={trendData} 
      />

      {/* 数据展示区域 */}
      <DataDisplay 
        orders={orders} 
        exchanges={exchanges} 
        loading={loading} 
      />
    </Card>
  );
};

export default Dashboard; 