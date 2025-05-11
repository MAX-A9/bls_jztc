import request from '../utils/request';

/**
 * 统计数据返回接口
 */
export interface StatsData {
  clientCount: number;   // 客户数量
  exchangeCount: number; // 兑换次数
  publishCount: number;  // 发布数量
  revenueAmount: number; // 收益金额
}

/**
 * 趋势数据类型
 */
export interface TrendData {
  labels: string[];  // 日期标签
  values: number[];  // 当前选择的数据类型的值
  allValues: {       // 所有数据类型的值
    clients: number[];    // 客户数量
    exchanges: number[];  // 兑换数量
    publishes: number[];  // 发布数量
    revenue: (number|string)[];  // 收益金额
  };
  dataType: 'clients' | 'exchanges' | 'publishes' | 'revenue' | 'all';  // 数据类型
  period: 'week' | 'month' | 'year';  // 周期
}

/**
 * 获取仪表盘统计数据
 * @param periodType 周期类型：week-本周，month-本月，year-本年
 * @returns 统计数据
 */
export async function getStatsData(periodType: 'week' | 'month' | 'year') {
  return request.get<any, { data: StatsData }>('/stats/data', {
    params: { periodType }
  });
}

/**
 * 获取趋势分析数据
 * @param periodType 周期类型：week-本周，month-本月，year-本年
 * @param dataType 数据类型：clients-客户数量，exchanges-兑换数量，publishes-发布数量，revenue-收益金额，all-所有数据
 * @returns 趋势数据
 */
export async function getTrendData(
  periodType: 'week' | 'month' | 'year',
  dataType: 'clients' | 'exchanges' | 'publishes' | 'revenue' | 'all' = 'all'
) {
  return request.get<any, { data: TrendData }>('/stats/trend', {
    params: { periodType, dataType }
  });
} 