import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Table, 
  Space, 
  Button, 
  Input, 
  Form, 
  Select, 
  Row, 
  Col, 
  Modal, 
  message, 
  Typography, 
  Tag, 
  Tooltip, 
  Popconfirm,
  Badge,
  Spin
} from 'antd';
import { 
  SearchOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  CheckCircleOutlined, 
  CloseCircleOutlined, 
  EyeOutlined, 
  UserOutlined, 
  ReloadOutlined,
  DownOutlined,
  UpOutlined
} from '@ant-design/icons';
import { getCommentList, updateCommentStatus, deleteComment, CommentItem, updateComment, getCommentDetail, CommentDetailItem } from '@/api/comment';
import dayjs from 'dayjs';
import './styles.css';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { TextArea } = Input;

const AllComments: React.FC = () => {
  // 状态定义
  const [form] = Form.useForm();
  const [editForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<CommentItem[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [selectedComment, setSelectedComment] = useState<CommentItem | null>(null);
  const [commentDetail, setCommentDetail] = useState<CommentDetailItem | null>(null);
  const [viewModalVisible, setViewModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [searchParams, setSearchParams] = useState({});
  const [detailLoading, setDetailLoading] = useState(false);
  const [isFilterVisible, setIsFilterVisible] = useState(false);

  // 组件加载完成后获取数据
  useEffect(() => {
    fetchComments();
  }, [pagination.current, pagination.pageSize, searchParams]);

  // 获取评论列表
  const fetchComments = async () => {
    setLoading(true);
    try {
      const params = {
        page: pagination.current,
        pageSize: pagination.pageSize,
        ...searchParams,
      };
      const res = await getCommentList(params);
      if (res.code === 0 && res.data) {
        setData(res.data.list);
        setPagination({
          ...pagination,
          total: res.data.total,
        });
      } else {
        message.error(res.message || '获取评论列表失败');
      }
    } catch (error) {
      console.error('获取评论列表出错:', error);
      message.error('获取评论列表出错');
    } finally {
      setLoading(false);
    }
  };

  // 搜索表单提交
  const onFinish = (values: any) => {
    // 过滤掉空值
    const params = Object.entries(values).reduce((acc: any, [key, value]) => {
      if (value !== undefined && value !== '') {
        acc[key] = value;
      }
      return acc;
    }, {});
    
    setPagination({ ...pagination, current: 1 });
    setSearchParams(params);
    setIsFilterVisible(false); // 隐藏高级筛选
  };

  // 重置搜索表单
  const handleReset = () => {
    form.resetFields();
    setPagination({ ...pagination, current: 1 });
    setSearchParams({});
  };

  // 表格分页、排序、筛选变化
  const handleTableChange = (newPagination: any) => {
    setPagination({
      ...pagination,
      current: newPagination.current,
      pageSize: newPagination.pageSize,
    });
  };

  // 查看评论详情
  const handleView = async (record: CommentItem) => {
    setSelectedComment(record);
    setViewModalVisible(true);
    setDetailLoading(true);
    
    try {
      const res = await getCommentDetail(record.id);
      if (res.code === 0 && res.data) {
        setCommentDetail(res.data);
      } else {
        message.error(res.message || '获取评论详情失败');
      }
    } catch (error) {
      console.error('获取评论详情失败:', error);
      message.error('获取评论详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  // 编辑评论
  const handleEdit = (record: CommentItem) => {
    setSelectedComment(record);
    // 先获取详情，再加载表单
    setDetailLoading(true);
    
    getCommentDetail(record.id)
      .then(res => {
        if (res.code === 0 && res.data) {
          setCommentDetail(res.data);
          editForm.setFieldsValue({
            comment: res.data.comment,
            status: res.data.status,
            realName: res.data.realName,
          });
          setEditModalVisible(true);
        } else {
          message.error(res.message || '获取评论详情失败');
        }
      })
      .catch(error => {
        console.error('获取评论详情失败:', error);
        message.error('获取评论详情失败');
      })
      .finally(() => {
        setDetailLoading(false);
      });
  };

  // 提交编辑
  const handleEditSubmit = async () => {
    try {
      const values = await editForm.validateFields();
      if (selectedComment) {
        setLoading(true);
        const response = await updateComment({
          id: selectedComment.id,
          ...values,
        });
        
        if (response.code === 0) {
          message.success('评论更新成功');
          setEditModalVisible(false);
          fetchComments();
        } else {
          message.error(response.message || '更新失败');
        }
      }
    } catch (error) {
      console.error('表单验证失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 更新评论状态
  const handleUpdateStatus = async (id: number, status: string, realName: string) => {
    try {
      setLoading(true);
      const response = await updateCommentStatus({ id, status, realName });
      
      if (response.code === 0) {
        message.success(`评论状态已更新为${status}`);
        fetchComments();
      } else {
        message.error(response.message || '更新状态失败');
      }
    } catch (error) {
      console.error('更新评论状态出错:', error);
      message.error('更新评论状态出错');
    } finally {
      setLoading(false);
    }
  };

  // 删除评论
  const handleDelete = async (id: number) => {
    try {
      setLoading(true);
      const response = await deleteComment(id);
      
      if (response.code === 0) {
        message.success('评论已删除');
        fetchComments();
      } else {
        message.error(response.message || '删除失败');
      }
    } catch (error) {
      console.error('删除评论出错:', error);
      message.error('删除评论出错');
    } finally {
      setLoading(false);
    }
  };

  // 渲染状态标签
  const renderStatusTag = (status: string) => {
    switch (status) {
      case '已审核':
        return <Tag color="success">已审核</Tag>;
      case '待审核':
        return <Tag color="warning">待审核</Tag>;
      case '已拒绝':
        return <Tag color="error">已拒绝</Tag>;
      default:
        return <Tag>{status}</Tag>;
    }
  };

  // 表格列定义
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '内容标题',
      dataIndex: 'contentTitle',
      key: 'contentTitle',
      width: 200,
      ellipsis: true,
      render: (text: string) => (
        <Tooltip title={text}>
          <span>{text}</span>
        </Tooltip>
      ),
    },
    {
      title: '评论用户',
      dataIndex: 'realName',
      key: 'realName',
      width: 120,
      render: (text: string) => (
        <Space>
          <UserOutlined />
          {text}
        </Space>
      ),
    },
    {
      title: '评论内容',
      dataIndex: 'comment',
      key: 'comment',
      ellipsis: true,
      render: (text: string) => (
        <Tooltip title={text}>
          <span>{text}</span>
        </Tooltip>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => renderStatusTag(status),
    },
    {
      title: '评论时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 170,
      render: (text: string) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      width: 240,
      render: (_: any, record: CommentItem) => (
        <Space size="small">
          <Button 
            type="primary" 
            icon={<EyeOutlined />} 
            size="small"
            onClick={() => handleView(record)}
          >
            查看
          </Button>
          
          <Button 
            type="primary" 
            icon={<EditOutlined />} 
            size="small"
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>

          {record.status !== '已审核' && (
            <Button
              type="primary"
              icon={<CheckCircleOutlined />}
              size="small"
              style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
              onClick={() => handleUpdateStatus(record.id, '已审核', record.realName)}
            >
              通过
            </Button>
          )}
          
          {record.status !== '已拒绝' && (
            <Button
              danger
              type="primary"
              icon={<CloseCircleOutlined />}
              size="small"
              onClick={() => handleUpdateStatus(record.id, '已拒绝', record.realName)}
            >
              拒绝
            </Button>
          )}
          
          <Popconfirm
            title="确定要删除此评论吗？"
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
          <div className="header-left">
            <Title level={4}>所有评论</Title>
          </div>
          <div className="header-right">
            <Space>
              <Input.Search
                placeholder="搜索评论内容/用户名"
                allowClear
                style={{ width: 240 }}
                onSearch={(value) => {
                  // 同时搜索评论内容和用户名
                  const params: any = {};
                  if (value) {
                    params.comment = value;
                    params.realName = value;
                  }
                  setSearchParams(params);
                  setPagination({ ...pagination, current: 1 });
                }}
              />
              <Button
                type="text"
                icon={isFilterVisible ? <UpOutlined /> : <DownOutlined />}
                onClick={() => setIsFilterVisible(!isFilterVisible)}
              >
                高级筛选
              </Button>
              <Badge count={data.filter(item => item.status === '待审核').length} overflowCount={999}>
                <Button
                  type="primary"
                  icon={<ReloadOutlined />}
                  onClick={() => fetchComments()}
                >
                  刷新数据
                </Button>
              </Badge>
            </Space>
          </div>
        </div>

        {isFilterVisible && (
          <div className="operation-logs-filter">
            <Form form={form} layout="horizontal" onFinish={onFinish}>
              <Row gutter={16}>
                <Col xs={24} sm={12} md={6} lg={6}>
                  <Form.Item name="contentId" label="内容ID">
                    <Input placeholder="输入内容ID" allowClear />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12} md={6} lg={6}>
                  <Form.Item name="realName" label="用户名">
                    <Input placeholder="输入用户名" allowClear />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12} md={6} lg={6}>
                  <Form.Item name="comment" label="评论内容">
                    <Input placeholder="输入评论内容" allowClear />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12} md={6} lg={6}>
                  <Form.Item name="status" label="状态">
                    <Select placeholder="请选择状态" allowClear>
                      <Option value="已审核">已审核</Option>
                      <Option value="待审核">待审核</Option>
                      <Option value="已拒绝">已拒绝</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col xs={24} sm={24} md={24} lg={24}>
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
                        icon={<ReloadOutlined />}
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
          dataSource={data}
          rowKey="id"
          pagination={{
            current: pagination.current,
            pageSize: pagination.pageSize,
            total: pagination.total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条`,
          }}
          loading={loading}
          onChange={handleTableChange}
          scroll={{ x: 'max-content' }}
          className="content-table"
        />
      </Card>

      {/* 查看评论详情的模态框 */}
      <Modal
        title="评论详情"
        open={viewModalVisible}
        onCancel={() => {
          setViewModalVisible(false);
          setCommentDetail(null);
        }}
        footer={[
          <Button key="close" onClick={() => {
            setViewModalVisible(false);
            setCommentDetail(null);
          }}>
            关闭
          </Button>,
        ]}
        width={700}
      >
        <Spin spinning={detailLoading}>
          {commentDetail ? (
            <div className="comment-detail">
              <Row gutter={[16, 16]}>
                <Col span={12}>
                  <Text strong>评论ID:</Text> {commentDetail.id}
                </Col>
                <Col span={12}>
                  <Text strong>状态:</Text> {renderStatusTag(commentDetail.status)}
                </Col>
                <Col span={12}>
                  <Text strong>用户:</Text> {commentDetail.realName} (ID: {commentDetail.clientId})
                </Col>
                <Col span={12}>
                  <Text strong>评论时间:</Text> {dayjs(commentDetail.createdAt).format('YYYY-MM-DD HH:mm:ss')}
                </Col>
                <Col span={24}>
                  <Text strong>内容标题:</Text>
                  <Paragraph ellipsis={{ rows: 1, expandable: true }}>{commentDetail.contentTitle}</Paragraph>
                </Col>
                <Col span={24}>
                  <Text strong>评论内容:</Text>
                  <div className="comment-content-box">
                    {commentDetail.comment}
                  </div>
                </Col>
                <Col span={24}>
                  <Text strong>更新时间:</Text> {dayjs(commentDetail.updatedAt).format('YYYY-MM-DD HH:mm:ss')}
                </Col>
              </Row>
            </div>
          ) : (
            <div className="comment-detail-empty">
              <p>加载评论详情中...</p>
            </div>
          )}
        </Spin>
      </Modal>

      {/* 编辑评论的模态框 */}
      <Modal
        title="编辑评论"
        open={editModalVisible}
        onCancel={() => {
          setEditModalVisible(false);
          setCommentDetail(null);
        }}
        onOk={handleEditSubmit}
        confirmLoading={loading}
      >
        <Spin spinning={detailLoading}>
          <Form
            form={editForm}
            layout="vertical"
          >
            <Form.Item
              name="comment"
              label="评论内容"
              rules={[{ required: true, message: '请输入评论内容' }]}
            >
              <TextArea rows={4} placeholder="请输入评论内容" />
            </Form.Item>
            <Form.Item
              name="realName"
              label="用户名"
              rules={[{ required: true, message: '请输入用户名' }]}
            >
              <Input placeholder="请输入用户名" />
            </Form.Item>
            <Form.Item
              name="status"
              label="评论状态"
              rules={[{ required: true, message: '请选择评论状态' }]}
            >
              <Select placeholder="请选择状态">
                <Option value="已审核">已审核</Option>
                <Option value="待审核">待审核</Option>
                <Option value="已拒绝">已拒绝</Option>
              </Select>
            </Form.Item>
          </Form>
        </Spin>
      </Modal>
    </div>
  );
};

export default AllComments; 