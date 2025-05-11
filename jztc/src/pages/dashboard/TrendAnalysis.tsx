import React, { useState, useEffect, useRef } from 'react';
import { Typography, Card, Row, Col, Segmented } from 'antd';
import { RiseOutlined, FallOutlined } from '@ant-design/icons';
import * as echarts from 'echarts';

const { Title } = Typography;

// 定义统计数据类型
interface StatisticsData {
  registration: number;
  exchange: number;
  publish: number;
  income: number;
}

// 定义趋势数据类型
interface TrendData {
  date: string;
  value: number;
  category: string;
}

interface TrendAnalysisProps {
  statistics: StatisticsData;
  trendData: TrendData[];
}

const TrendAnalysis: React.FC<TrendAnalysisProps> = ({ trendData }) => {
  // 趋势图类型，默认设置为'all'（全部）
  const [trendType, setTrendType] = useState<'all' | 'registration' | 'exchange' | 'publish' | 'income'>('all');
  
  // 图表实例ref
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  // 初始化并更新图表
  useEffect(() => {
    // 如果图表实例存在，销毁
    if (chartInstance.current) {
      chartInstance.current.dispose();
    }

    // 确保DOM已经渲染
    if (chartRef.current) {
      // 初始化图表
      chartInstance.current = echarts.init(chartRef.current);
      // 设置选项
      chartInstance.current.setOption(getEChartsOption());
    }

    // 组件卸载时清理
    return () => {
      if (chartInstance.current) {
        chartInstance.current.dispose();
        chartInstance.current = null;
      }
    };
  }, [trendData, trendType]);

  // 处理窗口大小变化
  useEffect(() => {
    const handleResize = () => {
      if (chartInstance.current) {
        chartInstance.current.resize();
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  // 处理趋势类型变化
  const handleTrendTypeChange = (value: string) => {
    setTrendType(value as 'all' | 'registration' | 'exchange' | 'publish' | 'income');
  };

  // 计算趋势变化百分比
  const calculateTrend = (type: 'registration' | 'exchange' | 'publish' | 'income'): number => {
    // 筛选当前类型的数据
    const currentData = trendData.filter(item => {
      return (type === 'registration' && item.category === '注册用户') || 
             (type === 'exchange' && item.category === '兑换次数') || 
             (type === 'publish' && item.category === '发布内容') || 
             (type === 'income' && item.category === '收益金额');
    });
    
    // 如果数据少于2个点，无法计算趋势
    if (currentData.length < 2) {
      return 0;
    }
    
    // 计算最早期和最近期的数据
    const oldest = currentData[0].value;
    const newest = currentData[currentData.length - 1].value;
    
    // 防止除以0的情况
    if (oldest === 0) {
      return newest > 0 ? 100 : 0;
    }
    
    // 计算百分比变化
    return Math.round(((newest - oldest) / oldest) * 100);
  };

  // 获取趋势图配置
  const getEChartsOption = () => {
    const seriesData: any[] = [];
    const categories: string[] = [];
    let dates: string[] = [];

    // 提取日期
    if (trendData.length > 0) {
      const firstCategory = trendData[0].category;
      dates = trendData
        .filter(item => item.category === firstCategory)
        .map(item => item.date);
    }

    // 按类别分组数据
    interface CategoryGroups {
      [key: string]: TrendData[];
    }
    
    const categoryGroups: CategoryGroups = {
      '注册用户': [],
      '兑换次数': [],
      '发布内容': [],
      '收益金额': []
    };

    trendData.forEach(item => {
      if (categoryGroups[item.category]) {
        categoryGroups[item.category].push(item);
      }
    });

    // 设置类别对应颜色
    const categoryColors = {
      '注册用户': '#1890ff',
      '兑换次数': '#722ed1',
      '发布内容': '#13c2c2',
      '收益金额': '#fa8c16'
    };
    
    // 按趋势类型过滤和处理数据
    Object.keys(categoryGroups).forEach(category => {
      if (
        trendType === 'all' || 
        (trendType === 'registration' && category === '注册用户') ||
        (trendType === 'exchange' && category === '兑换次数') ||
        (trendType === 'publish' && category === '发布内容') ||
        (trendType === 'income' && category === '收益金额')
      ) {
        const data = categoryGroups[category].map(item => item.value);
        if (data.length > 0) {
          categories.push(category);
          const color = categoryColors[category as keyof typeof categoryColors];
          
          seriesData.push({
            name: category,
            type: 'line',
            data: data,
            smooth: true,
            symbol: 'circle',
            symbolSize: 6,
            lineStyle: {
              width: 3,
              shadowColor: 'rgba(0,0,0,0.2)',
              shadowBlur: 10,
              shadowOffsetY: 10
            },
            itemStyle: {
              color: color,
              borderWidth: 2
            },
            emphasis: {
              itemStyle: {
                borderWidth: 4,
                shadowColor: color,
                shadowBlur: 10
              }
            },
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [{
                  offset: 0, 
                  color: color.replace('rgb', 'rgba').replace(')', ', 0.3)')
                }, {
                  offset: 1, 
                  color: color.replace('rgb', 'rgba').replace(')', ', 0.05)')
                }]
              }
            }
          });
        }
      }
    });

    return {
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        borderWidth: 0,
        padding: [16, 20],
        textStyle: {
          color: '#666'
        },
        shadowColor: 'rgba(0, 0, 0, 0.1)',
        shadowBlur: 10,
        shadowOffsetX: 2,
        shadowOffsetY: 2,
        borderRadius: 8,
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#ddd'
          }
        },
        formatter: (params: any[]) => {
          let result = `<div style="font-weight: 500; margin-bottom: 8px;">${params[0].axisValue}</div>`;
          params.forEach(param => {
            result += `<div style="display: flex; align-items: center; margin: 5px 0;">
              <span style="display:inline-block; width:8px; height:8px; border-radius:50%; background-color:${param.color}; margin-right: 8px;"></span>
              <span style="margin-right: 12px;">${param.seriesName}:</span>
              <span style="font-weight: 500;">${param.value}</span>
            </div>`;
          });
          return result;
        }
      },
      legend: {
        data: categories,
        top: 0,
        right: 80,
        textStyle: {
          color: '#666'
        },
        icon: 'circle',
        itemWidth: 10,
        itemHeight: 10,
        itemGap: 20
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        top: 50,
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: dates,
        axisLine: {
          lineStyle: {
            color: '#eee'
          }
        },
        axisTick: {
          show: false
        },
        axisLabel: {
          color: '#666',
          fontSize: 12
        }
      },
      yAxis: {
        type: 'value',
        axisLine: {
          show: false
        },
        axisTick: {
          show: false
        },
        axisLabel: {
          color: '#666',
          fontSize: 12,
          margin: 10
        },
        splitLine: {
          lineStyle: {
            color: '#f0f0f0',
            type: 'dashed'
          }
        }
      },
      series: seriesData
    };
  };

  // 获取趋势百分比显示
  const getTrendValue = () => {
    if (trendType === 'all') return null;
    return calculateTrend(trendType);
  };

  const trendValue = getTrendValue();

  return (
    <>
      <div className="dashboard-title trend-analysis-title">
        <Title level={4}>趋势分析</Title>
        <Segmented
          options={[
            { label: '全部', value: 'all' },
            { label: '注册用户', value: 'registration' },
            { label: '兑换次数', value: 'exchange' },
            { label: '发布内容', value: 'publish' },
            { label: '收益金额', value: 'income' }
          ]}
          value={trendType}
          onChange={value => handleTrendTypeChange(value as string)}
        />
      </div>
      
      <Row className="trend-chart-row">
        <Col span={24}>
          <Card className="trend-chart-card">
            <div className="trend-indicator">
              {trendType !== 'all' && trendValue !== null && (
                <div className={`trend-badge ${trendValue > 0 ? 'trend-up' : trendValue < 0 ? 'trend-down' : 'trend-flat'}`}>
                  {trendValue > 0 ? <RiseOutlined /> : trendValue < 0 ? <FallOutlined /> : null}
                  <span>{Math.abs(trendValue)}%</span>
                </div>
              )}
            </div>
            <div className="trend-chart">
              <div ref={chartRef} style={{ height: '100%', width: '100%' }}></div>
            </div>
          </Card>
        </Col>
      </Row>
    </>
  );
};

export default TrendAnalysis; 