import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Space,
  Input,
  DatePicker,
  Button,
  Typography,
  Tag,
  Select,
  Form,
  Row,
  Col,
  message,
  Popconfirm,
} from 'antd';
import {
  SearchOutlined,
  FileSearchOutlined,
  ReloadOutlined,
  ExportOutlined,
  UpOutlined,
  DownOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { RangePickerProps } from 'antd/es/date-picker';
import dayjs from 'dayjs';
import './OperationLogs.css';
import { getOperationLogs, deleteOperationLog, exportOperationLogs, LogListParams } from '../api/logs';

const { Title } = Typography;
const { Option } = Select;
const { RangePicker } = DatePicker;

// 日志类型定义
interface LogData {
  id: number;
  operatorName: string;
  ip: string;
  action: string;
  status: '成功' | '失败';
  time: string;
  details?: string;
}

// 日志类型选项
const logActionTypes = [
  { value: '登录', color: 'blue' },
  { value: '登出', color: 'purple' },
  { value: '添加', color: 'green' },
  { value: '修改', color: 'orange' },
  { value: '删除', color: 'red' },
  { value: '导出', color: 'cyan' },
  { value: '导入', color: 'magenta' },
  { value: '查询', color: 'gold' },
];

const OperationLogs: React.FC = () => {
  const [logs, setLogs] = useState<LogData[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  
  // 搜索条件
  const [form] = Form.useForm();
  const [searchParams, setSearchParams] = useState({
    keyword: '',
    dateRange: null as [dayjs.Dayjs, dayjs.Dayjs] | null,
    action: '',
    result: '',
  });

  // 获取日志数据
  const fetchLogs = async () => {
    setLoading(true);
    try {
      // 构建API参数
      const params: LogListParams = {
        page: current,
        pageSize,
        keyword: searchParams.keyword || undefined,
        action: searchParams.action || undefined,
        result: searchParams.result || undefined,
      };

      // 添加日期范围
      if (searchParams.dateRange && searchParams.dateRange.length === 2) {
        params.startTime = searchParams.dateRange[0].format('YYYY-MM-DD HH:mm:ss');
        params.endTime = searchParams.dateRange[1].format('YYYY-MM-DD HH:mm:ss');
      }

      // 发起API请求
      const response = await getOperationLogs(params);
      // 处理嵌套的数据
      const responseData = response.data || response;
      
      // 转换数据格式
      const formattedLogs: LogData[] = (responseData.list || []).map(item => ({
        id: item.id,
        operatorName: item.username,
        ip: item.operationIp,
        action: item.action,
        status: item.operationResult === 1 ? '成功' : '失败',
        time: item.operationTime,
        details: item.details,
      }));
      
      setLogs(formattedLogs);
      setTotal(responseData.total || 0);
    } catch (error) {
      console.error('获取日志数据失败:', error);
      message.error('获取日志数据失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 初始加载和条件变化时重新获取数据
  useEffect(() => {
    fetchLogs();
  }, [current, pageSize, searchParams]);

  // 处理搜索
  const handleSearch = () => {
    const values = form.getFieldsValue();
    setSearchParams({
      ...searchParams,
      dateRange: values.dateRange,
      action: values.action || '',
      result: values.result || '',
    });
    setCurrent(1); // 重置到第一页
    setIsFilterVisible(false); // 隐藏高级筛选
  };

  // 重置搜索
  const handleReset = () => {
    form.resetFields();
    setSearchParams({
      keyword: '',
      dateRange: null,
      action: '',
      result: '',
    });
    setCurrent(1);
  };

  // 导出日志数据
  const handleExport = async () => {
    try {
      setLoading(true);
      // 构建导出参数
      const params: LogListParams = {
        keyword: searchParams.keyword || undefined,
        action: searchParams.action || undefined,
        result: searchParams.result || undefined,
      };

      // 添加日期范围
      if (searchParams.dateRange && searchParams.dateRange.length === 2) {
        params.startTime = searchParams.dateRange[0].format('YYYY-MM-DD HH:mm:ss');
        params.endTime = searchParams.dateRange[1].format('YYYY-MM-DD HH:mm:ss');
      }

      const response = await exportOperationLogs(params);
      // 处理嵌套的数据
      const responseData = response.data || response;
      
      // 创建一个临时链接并点击它来下载文件
      if (responseData && responseData.url) {
        const link = document.createElement('a');
        link.href = responseData.url;
        link.target = '_blank';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        message.success('日志数据导出成功');
      } else {
        message.error('导出失败，未获取到下载链接');
      }
    } catch (error) {
      console.error('导出日志数据失败:', error);
      message.error('导出日志数据失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 自定义日期范围选择限制（只能选择最近30天）
  const disabledDate: RangePickerProps['disabledDate'] = (current) => {
    return current && current > dayjs().endOf('day') || current < dayjs().subtract(30, 'days').startOf('day');
  };

  // 渲染状态标签
  const renderStatusTag = (status: string) => {
    return status === '成功' ? (
      <div style={{ textAlign: 'center' }}>
        <Tag color="success">成功</Tag>
      </div>
    ) : (
      <div style={{ textAlign: 'center' }}>
        <Tag color="error">失败</Tag>
      </div>
    );
  };

  // 渲染操作类型标签
  const renderActionTag = (action: string) => {
    const actionType = logActionTypes.find(type => type.value === action);
    return <Tag color={actionType?.color || 'default'}>{action}</Tag>;
  };

  // 表格列定义
  const columns: ColumnsType<LogData> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '操作人',
      dataIndex: 'operatorName',
      key: 'operatorName',
      width: 100,
    },
    {
      title: 'IP地址',
      dataIndex: 'ip',
      key: 'ip',
      width: 120,
    },
    {
      title: '操作类型',
      dataIndex: 'action',
      key: 'action',
      width: 100,
      render: renderActionTag,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: renderStatusTag,
      align: 'center',
    },
    {
      title: '操作时间',
      dataIndex: 'time',
      key: 'time',
      width: 180,
      sorter: (a, b) => dayjs(a.time).unix() - dayjs(b.time).unix(),
    },
    {
      title: '操作',
      key: 'operation',
      width: 100,
      render: (_, record) => (
        <Popconfirm
          title="确定删除此日志记录吗？"
          description="此操作不可恢复，请谨慎操作。"
          onConfirm={() => handleDelete(record.id)}
          okText="确定"
          cancelText="取消"
        >
          <Button type="link" danger>
            删除
          </Button>
        </Popconfirm>
      ),
      align: 'center',
    },
  ];

  // 高级筛选显示状态
  const [isFilterVisible, setIsFilterVisible] = useState(false);

  // 删除日志
  const handleDelete = async (id: number) => {
    try {
      setLoading(true);
      const response = await deleteOperationLog(id);
      message.success('日志删除成功');
      
      // 如果当前页没有数据了，且不是第一页，则回到上一页
      if (logs.length === 1 && current > 1) {
        setCurrent(current - 1);
      } else {
        // 重新加载数据
        fetchLogs();
      }
    } catch (error) {
      console.error('删除日志失败:', error);
      message.error('删除日志失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="operation-logs-container">
      <Card className="main-dashboard-card">
        <div className="operation-logs-header">
          <div className="header-left">
            <Title level={4}>操作日志</Title>
          </div>
          <div className="header-right">
            <Space>
              <Input.Search
                placeholder="搜索操作人/IP"
                allowClear
                style={{ width: 240 }}
                onSearch={(value) => {
                  setSearchParams({ ...searchParams, keyword: value });
                  setCurrent(1);
                }}
              />
              <Button
                type="text"
                icon={isFilterVisible ? <UpOutlined /> : <DownOutlined />}
                onClick={() => setIsFilterVisible(!isFilterVisible)}
              >
                高级筛选
              </Button>
              <Button
                type="primary"
                ghost
                icon={<ReloadOutlined />}
                onClick={() => fetchLogs()}
              >
                刷新
              </Button>
              <Button
                type="primary"
                icon={<ExportOutlined />}
                onClick={handleExport}
              >
                导出
              </Button>
            </Space>
          </div>
        </div>

        {isFilterVisible && (
          <div className="operation-logs-filter">
            <Form form={form} layout="horizontal" onFinish={handleSearch}>
              <Row gutter={16}>
                <Col xs={24} sm={24} md={8} lg={8} xl={6}>
                  <Form.Item name="dateRange" label="时间范围">
                    <RangePicker
                      style={{ width: '100%' }}
                      disabledDate={disabledDate}
                      ranges={{
                        '今天': [dayjs().startOf('day'), dayjs().endOf('day')],
                        '昨天': [dayjs().subtract(1, 'day').startOf('day'), dayjs().subtract(1, 'day').endOf('day')],
                        '最近3天': [dayjs().subtract(2, 'day').startOf('day'), dayjs().endOf('day')],
                        '最近7天': [dayjs().subtract(6, 'day').startOf('day'), dayjs().endOf('day')],
                      }}
                    />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12} md={6} lg={6} xl={4}>
                  <Form.Item name="action" label="操作类型">
                    <Select placeholder="请选择" allowClear>
                      {logActionTypes.map(type => (
                        <Option key={type.value} value={type.value}>
                          {type.value}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12} md={6} lg={6} xl={4}>
                  <Form.Item name="result" label="状态">
                    <Select placeholder="请选择" allowClear>
                      <Option value="1">成功</Option>
                      <Option value="0">失败</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col xs={24} sm={24} md={4} lg={4} xl={10}>
                  <div style={{ textAlign: 'right' }}>
                    <Space>
                      <Button 
                        type="primary" 
                        htmlType="submit" 
                        icon={<SearchOutlined />}
                      >
                        搜索
                      </Button>
                      <Button 
                        onClick={handleReset} 
                        icon={<FileSearchOutlined />}
                      >
                        重置
                      </Button>
                    </Space>
                  </div>
                </Col>
              </Row>
            </Form>
          </div>
        )}

        <Table
          columns={columns}
          dataSource={logs}
          rowKey="id"
          expandable={{
            expandedRowRender: (record) => record.details ? (
              <div style={{ padding: '12px 0' }}>
                <div style={{ fontSize: '12px', color: '#666' }}>{record.details}</div>
              </div>
            ) : null,
            expandRowByClick: true,
            showExpandColumn: false,
            rowExpandable: (record) => !!record.details,
          }}
          pagination={{
            current,
            pageSize,
            total,
            showTotal: (total) => `共 ${total} 条记录`,
            showSizeChanger: true,
            showQuickJumper: true,
            onChange: (page, pageSize) => {
              setCurrent(page);
              setPageSize(pageSize);
            },
          }}
          loading={loading}
          scroll={{ x: 'max-content' }}
        />
      </Card>
    </div>
  );
};

export default OperationLogs; 