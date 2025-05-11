import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Popconfirm,
  message,
  Typography,
  Tag,
  Tooltip,
  Input,
  Select,
  Modal,
  DatePicker
} from 'antd';
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  SearchOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import './styles.css';
import { 
  getContentList, 
  deleteContent, 
  updateContentStatus, 
  updateContentRecommend,
  type ContentItem 
} from '@/api/content';
import { getCategoriesForSelect } from '@/api/category';
import dayjs from 'dayjs';

const { Title } = Typography;
const { Option } = Select;

const AllContent: React.FC = () => {
  const [contents, setContents] = useState<ContentItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [searchText, setSearchText] = useState<string>('');
  const [filterCategory, setFilterCategory] = useState<string>('');
  const [total, setTotal] = useState<number>(0);
  const [current, setCurrent] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [topUntilDate, setTopUntilDate] = useState<dayjs.Dayjs | null>(null);
  const [isTopModalVisible, setIsTopModalVisible] = useState<boolean>(false);
  const [currentTopItemId, setCurrentTopItemId] = useState<number | null>(null);
  const [categories, setCategories] = useState<{ label: string; value: string; type: string }[]>([]);
  const [categoryLoading, setCategoryLoading] = useState<boolean>(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchCategories();
    fetchContents();
  }, [current, pageSize, searchText, filterCategory]);

  const fetchContents = async () => {
    setLoading(true);
    try {
      const params: any = {
        page: current,
        pageSize: pageSize
      };
      
      if (searchText) {
        params.title = searchText;
      }
      
      if (filterCategory) {
        params.category = filterCategory;
      }
      
      const res = await getContentList(params);
      
      if (res.code === 0) {
        setContents(res.data.list);
        setTotal(res.data.total);
      } else {
        message.error(res.message || '获取内容列表失败');
      }
    } catch (error) {
      console.error('获取内容列表失败:', error);
      message.error('获取内容列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    setCategoryLoading(true);
    try {
      const options = await getCategoriesForSelect();
      setCategories(options);
    } catch (error) {
      console.error('获取分类列表失败:', error);
      message.error('获取分类列表失败');
    } finally {
      setCategoryLoading(false);
    }
  };

  const handleAdd = () => {
    navigate('/content/add-content');
  };

  const handleEdit = (record: ContentItem) => {
    navigate(`/content/edit-content/${record.id}`);
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await deleteContent(id);
      
      if (res.code === 0) {
        message.success('内容删除成功');
        fetchContents();
      } else {
        message.error(res.message || '删除内容失败');
      }
    } catch (error) {
      console.error('删除内容失败:', error);
      message.error('删除内容失败');
    }
  };
  
  // 更新内容状态
  const handleStatusChange = async (id: number, status: string) => {
    try {
      const res = await updateContentStatus({ id, status });
      
      if (res.code === 0) {
        message.success('状态更新成功');
        fetchContents();
      } else {
        message.error(res.message || '状态更新失败');
      }
    } catch (error) {
      console.error('状态更新失败:', error);
      message.error('状态更新失败');
    }
  };
  
  // 打开置顶设置模态框
  const showTopModal = (id: number) => {
    setCurrentTopItemId(id);
    setTopUntilDate(dayjs().add(1, 'day')); // 默认设置为明天
    setIsTopModalVisible(true);
  };
  
  // 关闭置顶设置模态框
  const handleTopModalCancel = () => {
    setIsTopModalVisible(false);
    setCurrentTopItemId(null);
    setTopUntilDate(null);
  };
  
  // 确认置顶设置
  const handleTopModalOk = async () => {
    if (!currentTopItemId || !topUntilDate) {
      message.error('请设置置顶截止时间');
      return;
    }
    
    try {
      const topUntil = topUntilDate.format('YYYY-MM-DD HH:mm:ss');
      
      const res = await updateContentRecommend({ 
        id: currentTopItemId, 
        isRecommended: true,
        topUntil 
      });
      
      if (res.code === 0) {
        message.success('设置置顶成功');
        fetchContents();
        handleTopModalCancel();
      } else {
        message.error(res.message || '设置置顶失败');
      }
    } catch (error) {
      console.error('设置置顶失败:', error);
      message.error('设置置顶失败');
    }
  };
  
  // 取消置顶
  const handleCancelTop = async (id: number) => {
    try {
      const res = await updateContentRecommend({ 
        id: id, 
        isRecommended: false 
      });
      
      if (res.code === 0) {
        message.success('取消置顶成功');
        fetchContents();
      } else {
        message.error(res.message || '取消置顶失败');
      }
    } catch (error) {
      console.error('取消置顶失败:', error);
      message.error('取消置顶失败');
    }
  };

  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrent(1); // 搜索时重置到第一页
  };

  const handleCategoryFilter = (value: string) => {
    setFilterCategory(value);
    setCurrent(1); // 筛选时重置到第一页
  };

  const getStatusTag = (status: string) => {
    switch (status) {
      case '待审核':
        return <Tag color="processing">待审核</Tag>;
      case '已发布':
        return <Tag color="success">已发布</Tag>;
      case '已下架':
        return <Tag color="error">已下架</Tag>;
      default:
        return <Tag color="default">{status}</Tag>;
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
      fixed: 'left' as const,
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: {
        showTitle: false,
      },
      render: (text: string) => (
        <Tooltip placement="topLeft" title={text}>
          <div className="content-title-cell">{text}</div>
        </Tooltip>
      ),
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      width: 100,
      filters: categories
        .filter(item => item.type !== 'group')
        .map(item => ({
          text: item.label,
          value: item.label
        })),
      onFilter: (value: any, record: ContentItem) => record.category === value,
    },
    {
      title: '作者',
      dataIndex: 'author',
      key: 'author',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string, record: ContentItem) => (
        <Tooltip title="点击修改状态">
          <div 
            onClick={() => {
              const newStatus = status === '已发布' ? '已下架' : (status === '已下架' ? '待审核' : '已发布');
              handleStatusChange(record.id, newStatus);
            }}
            style={{ cursor: 'pointer' }}
          >
            {getStatusTag(status)}
          </div>
        </Tooltip>
      ),
    },
    {
      title: '浏览量',
      dataIndex: 'views',
      key: 'views',
      width: 100,
      sorter: (a: ContentItem, b: ContentItem) => a.views - b.views,
    },
    {
      title: '想要',
      dataIndex: 'likes',
      key: 'likes',
      width: 100,
      sorter: (a: ContentItem, b: ContentItem) => a.likes - b.likes,
    },
    {
      title: '评论数',
      dataIndex: 'comments',
      key: 'comments',
      width: 100,
      sorter: (a: ContentItem, b: ContentItem) => a.comments - b.comments,
    },
    {
      title: '置顶',
      dataIndex: 'isRecommended',
      key: 'isRecommended',
      width: 80,
      render: (isRecommended: boolean, record: ContentItem) => {
        if (!isRecommended) {
          return (
            <Tooltip title="点击设置置顶">
              <Tag 
                color="default" 
                style={{ cursor: 'pointer' }}
                onClick={() => showTopModal(record.id)}
              >
                未置顶
              </Tag>
            </Tooltip>
          );
        }
        
        // 根据置顶时间判断是否过期
        const isTopExpired = record.topUntil && new Date(record.topUntil) < new Date();
        
        if (isTopExpired) {
          return <Tag color="default">已过期</Tag>;
        } else {
          return (
            <Tooltip title="点击取消置顶">
              <Tag 
                color="gold" 
                style={{ cursor: 'pointer' }}
                onClick={() => handleCancelTop(record.id)}
              >
                置顶
              </Tag>
            </Tooltip>
          );
        }
      },
    },
    {
      title: '发布时间',
      dataIndex: 'publishedAt',
      key: 'publishedAt',
      width: 180,
    },
    {
      title: '到期时间',
      dataIndex: 'expiresAt',
      key: 'expiresAt',
      width: 180,
      render: (expiresAt: string | null) => expiresAt ? expiresAt : '-',
    },
    {
      title: '置顶时间',
      dataIndex: 'topUntil',
      key: 'topUntil',
      width: 180,
      render: (topUntil: string | null) => topUntil ? topUntil : '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right' as const,
      render: (_: any, record: ContentItem) => (
        <Space size="small">
          <Button 
            type="primary" 
            icon={<EditOutlined />} 
            size="small"
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除此内容吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button 
              danger 
              icon={<DeleteOutlined />} 
              size="small"
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="category-container">
      <Card className="main-dashboard-card">
        <div className="category-header">
          <Title level={4}>所有内容</Title>
          <Space>
            <Input
              placeholder="搜索标题或作者"
              prefix={<SearchOutlined />}
              onChange={(e) => handleSearch(e.target.value)}
              style={{ width: 200 }}
              allowClear
            />
            <Select
              placeholder="分类筛选"
              style={{ width: 120 }}
              onChange={handleCategoryFilter}
              allowClear
              loading={categoryLoading}
            >
              {categories
                .filter(item => item.type !== 'group')
                .map(item => (
                  <Option key={item.value} value={item.label}>
                    {item.label}
                  </Option>
                ))}
            </Select>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
              添加内容
            </Button>
            <Button icon={<ReloadOutlined />} onClick={() => fetchContents()} loading={loading}>
              刷新
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={contents}
          rowKey="id"
          loading={loading}
          pagination={{
            current: current,
            pageSize: pageSize,
            total: total,
            onChange: (page, pageSize) => {
              setCurrent(page);
              setPageSize(pageSize || 10);
            },
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
          }}
          scroll={{ x: 'max-content' }}
          className="content-table"
        />
      </Card>
      
      {/* 置顶设置模态框 */}
      <Modal
        title="设置置顶"
        open={isTopModalVisible}
        onOk={handleTopModalOk}
        onCancel={handleTopModalCancel}
        okText="确定"
        cancelText="取消"
      >
        <div style={{ marginBottom: 16 }}>请选择置顶截止时间：</div>
        <DatePicker
          showTime
          placeholder="选择置顶截止时间"
          value={topUntilDate}
          onChange={(date) => setTopUntilDate(date)}
          style={{ width: '100%' }}
          format="YYYY-MM-DD HH:mm:ss"
          disabledDate={(current) => current && current < dayjs().startOf('day')}
        />
      </Modal>
    </div>
  );
};

export default AllContent; 