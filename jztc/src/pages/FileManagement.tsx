import React, { useState, useEffect, useRef, useCallback } from 'react';
import './FileManagement.css';
import {
  Card,
  Table,
  Button,
  Space,
  Input,
  Modal,
  Form,
  Upload,
  message,
  Popconfirm,
  Tag,
  Typography,
  Tooltip,
  Divider,
  Select,
  Radio,
  Checkbox,
  Row,
  Col,
  Empty,
  Spin,
} from 'antd';
import {
  FileOutlined,
  UploadOutlined,
  SearchOutlined,
  DeleteOutlined,
  DownloadOutlined,
  EyeOutlined,
  AppstoreOutlined,
  PictureOutlined,
  FileTextOutlined,
  FileExcelOutlined,
  FilePdfOutlined,
  FileZipOutlined,
  BarsOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { UploadProps } from 'antd';
import type { TablePaginationConfig } from 'antd/es/table';
import { 
  getFileList, 
  getFileDetail, 
  uploadFile, 
  updateFilePublic, 
  deleteFile, 
  batchDeleteFiles,
  type FileInfo,
  type FileListReq
} from '@/api/file';

const { Title } = Typography;
const { Dragger } = Upload;
const { Option } = Select;

type ViewMode = 'grid' | 'list';

const FileManagement: React.FC = () => {
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(24);
  const [searchText, setSearchText] = useState('');
  const [fileType, setFileType] = useState<string>('');
  const [isPublic, setIsPublic] = useState<boolean | undefined>(undefined);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isPreviewVisible, setIsPreviewVisible] = useState(false);
  const [previewFile, setPreviewFile] = useState<FileInfo | null>(null);
  const [viewMode, setViewMode] = useState<ViewMode>('grid');
  const [selectedFiles, setSelectedFiles] = useState<number[]>([]);
  const [form] = Form.useForm();
  const [pagination, setPagination] = useState<TablePaginationConfig>({
    current,
    pageSize,
    total,
    onChange: (page: number, pageSize: number) => {
      setCurrent(page);
      if (pageSize) setPageSize(pageSize);
    },
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number) => `共 ${total} 条`,
  });
  const [hasMore, setHasMore] = useState(true);
  const [gridFiles, setGridFiles] = useState<FileInfo[]>([]);
  const observer = useRef<IntersectionObserver | null>(null);
  const loadingRef = useCallback((node: HTMLDivElement | null) => {
    if (loading || !node) return;
    if (observer.current) observer.current.disconnect();
    
    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasMore && viewMode === 'grid' && !loading) {
        setCurrent(prev => prev + 1);
      }
    }, {
      rootMargin: '0px 0px 100px 0px',
      threshold: 0.1
    });
    
    node && observer.current.observe(node);
  }, [loading, hasMore, viewMode]);
  const fileGridRef = useRef<HTMLDivElement>(null);

  // 获取文件列表
  const fetchFiles = async () => {
    setLoading(true);
    try {
      const params: FileListReq = {
        Page: current,
        PageSize: pageSize
      };
      
      if (searchText) {
        params.Keyword = searchText;
      }
      
      if (fileType) {
        params.Type = fileType;
      }
      
      if (isPublic !== undefined) {
        params.IsPublic = isPublic;
      }

      const res = await getFileList(params);
      if (res.code === 0) {
        if (viewMode === 'grid') {
          if (current === 1) {
            setGridFiles(res.data.list);
          } else {
            setGridFiles(prev => [...prev, ...res.data.list]);
          }
          
          const newLength = current === 1 
            ? res.data.list.length 
            : gridFiles.length + res.data.list.length;
          
          setHasMore(newLength < res.data.total);
        } else {
          setFiles(res.data.list);
        }
        
        setTotal(res.data.total);
      } else {
        message.error(res.message || '获取文件列表失败');
      }
    } catch (error) {
      message.error('获取文件列表失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 初始加载数据
  useEffect(() => {
    fetchFiles();
  }, [current, pageSize, searchText, fileType, isPublic, viewMode]);
  
  // 视图切换时重置滚动位置
  useEffect(() => {
    if (viewMode === 'grid') {
      // 重置状态
      setCurrent(1);
      setGridFiles([]);
      setHasMore(true);
      
      // 重置滚动位置
      setTimeout(() => {
        if (fileGridRef.current) {
          fileGridRef.current.scrollTop = 0;
          
          // 重置父容器和页面滚动位置
          const parentElement = document.querySelector('.full-content-layout');
          if (parentElement) {
            parentElement.scrollTop = 0;
          }
          
          window.scrollTo(0, 0);
        }
      }, 0);
    }
  }, [viewMode]);
  
  // 添加观察器清理
  useEffect(() => {
    return () => {
      if (observer.current) {
        observer.current.disconnect();
      }
    };
  }, []);
  
  // 处理搜索
  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrent(1);
  };
  
  // 处理预览
  const handlePreview = async (id: number) => {
    try {
      const res = await getFileDetail(id);
      if (res.code === 0) {
        setPreviewFile(res.data);
        setIsPreviewVisible(true);
      } else {
        message.error(res.message || '获取文件详情失败');
      }
    } catch (error) {
      message.error('获取文件详情失败');
    }
  };
  
  // 处理下载
  const handleDownload = (file: FileInfo) => {
    message.success(`开始下载文件：${file.name}`);
    const link = document.createElement('a');
    link.href = file.url;
    link.target = '_blank';
    link.download = file.name;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };
  
  // 处理删除
  const handleDelete = async (id: number) => {
    try {
      const res = await deleteFile(id);
      if (res.code === 0 && res.data.success) {
        setSelectedFiles(selectedFiles.filter(fileId => fileId !== id));
        message.success(res.data.message || '文件删除成功');
        fetchFiles(); // 重新加载列表
      } else {
        message.error(res.data.message || '文件删除失败');
      }
    } catch (error) {
      message.error('删除文件失败');
    }
  };
  
  // 批量删除
  const handleBatchDelete = async () => {
    if (selectedFiles.length === 0) {
      message.warning('请选择要删除的文件');
      return;
    }
    
    try {
      const res = await batchDeleteFiles(selectedFiles);
      if (res.code === 0 && res.data.success) {
        setSelectedFiles([]);
        message.success(`成功删除${res.data.count}个文件`);
        fetchFiles(); // 重新加载列表
      } else {
        message.error(res.data.message || '批量删除文件失败');
      }
    } catch (error) {
      message.error('批量删除文件失败');
    }
  };
  
  // 更新文件公开状态
  const handleUpdatePublic = async (id: number, isPublic: boolean) => {
    try {
      const res = await updateFilePublic({
        Id: id,
        IsPublic: isPublic
      });
      
      if (res.code === 0 && res.data.success) {
        message.success(res.data.message || '更新文件状态成功');
        fetchFiles(); // 重新加载列表
      } else {
        message.error(res.data.message || '更新文件状态失败');
      }
    } catch (error) {
      message.error('更新文件状态失败');
    }
  };
  
  // 处理文件上传
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const formData = new FormData();
      
      // 获取上传的文件
      const { uploadFiles } = values;
      if (uploadFiles && uploadFiles.fileList && uploadFiles.fileList.length > 0) {
        formData.append('File', uploadFiles.fileList[0].originFileObj);
        
        const res = await uploadFile(formData);
        if (res.code === 0) {
          message.success('文件上传成功');
          setIsModalVisible(false);
          form.resetFields();
          fetchFiles(); // 重新加载列表
        } else {
          message.error(res.message || '文件上传失败');
        }
      }
    } catch (error) {
      message.error('上传文件失败');
    }
  };
  
  // 处理模态框取消
  const handleCancel = () => {
    if (isPreviewVisible) {
      setIsPreviewVisible(false);
      setPreviewFile(null);
    } else {
      setIsModalVisible(false);
    }
    form.resetFields();
  };
  
  // 处理文件选择
  const toggleFileSelection = (fileId: number, event?: React.MouseEvent) => {
    if (event) {
      event.stopPropagation();
    }
    
    setSelectedFiles(prevSelected => {
      if (prevSelected.includes(fileId)) {
        return prevSelected.filter(id => id !== fileId);
      } else {
        return [...prevSelected, fileId];
      }
    });
  };
  
  // 全选/取消全选
  const toggleSelectAll = () => {
    if (selectedFiles.length === files.length) {
      setSelectedFiles([]);
    } else {
      setSelectedFiles(files.map(file => file.id));
    }
  };
  
  // 处理行点击
  const handleRowClick = (record: FileInfo) => {
    return {
      onClick: () => {
        toggleFileSelection(record.id);
      },
    };
  };

  // 文件类型标签渲染
  const renderFileTypeTag = (type: string) => {
    let color = '';
    let text = type;
    
    if (type.includes('image')) {
      color = 'blue';
      text = '图片';
    } else if (type.includes('pdf')) {
      color = 'red';
      text = 'PDF';
    } else if (type.includes('excel') || type.includes('spreadsheet')) {
      color = 'green';
      text = '表格';
    } else if (type.includes('word') || type.includes('document')) {
      color = 'purple';
      text = '文档';
    } else if (type.includes('zip') || type.includes('compressed')) {
      color = 'orange';
      text = '压缩包';
    }
    
    return <Tag color={color}>{text}</Tag>;
  };
  
  // 根据文件类型获取图标
  const getFileIcon = (type: string, size: 'small' | 'large' = 'small') => {
    const fontSize = size === 'large' ? 48 : 24;
    
    if (type.includes('image')) {
      return <PictureOutlined style={{ fontSize }} />;
    } else if (type.includes('pdf')) {
      return <FilePdfOutlined style={{ fontSize }} />;
    } else if (type.includes('excel') || type.includes('spreadsheet')) {
      return <FileExcelOutlined style={{ fontSize }} />;
    } else if (type.includes('zip') || type.includes('compressed')) {
      return <FileZipOutlined style={{ fontSize }} />;
    } else if (type.includes('word') || type.includes('document')) {
      return <FileTextOutlined style={{ fontSize }} />;
    } else {
      return <FileOutlined style={{ fontSize }} />;
    }
  };
  
  // 表格列定义
  const columns: ColumnsType<FileInfo> = [
    {
      title: (
        <Checkbox
          checked={selectedFiles.length === files.length && files.length > 0}
          indeterminate={selectedFiles.length > 0 && selectedFiles.length < files.length}
          onChange={toggleSelectAll}
        />
      ),
      dataIndex: 'id',
      key: 'selection',
      width: 50,
      render: (id) => (
        <Checkbox
          checked={selectedFiles.includes(id)}
          onChange={(e) => {
            e.stopPropagation();
            toggleFileSelection(id);
          }}
          onClick={(e) => e.stopPropagation()}
        />
      ),
    },
    {
      title: '预览',
      dataIndex: 'type',
      key: 'preview',
      width: 70,
      render: (type, record) => (
        <div 
          className="file-icon" 
          onClick={(e) => {
            e.stopPropagation();
            handlePreview(record.id);
          }}
        >
          {getFileIcon(type)}
        </div>
      ),
    },
    {
      title: '文件名',
      dataIndex: 'name',
      key: 'fileName',
      render: (text) => <span className="file-name">{text}</span>,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'fileType',
      width: 90,
      render: (type) => renderFileTypeTag(type),
    },
    {
      title: '大小',
      dataIndex: 'sizeFormat',
      key: 'fileSize',
      width: 100,
    },
    {
      title: '上传人',
      dataIndex: 'username',
      key: 'uploader',
      width: 100,
    },
    {
      title: '上传时间',
      dataIndex: 'createdAt',
      key: 'uploadTime',
      width: 170,
    },
    {
      title: '状态',
      dataIndex: 'isPublic',
      key: 'status',
      width: 90,
      render: (isPublic) => (
        <Tag color={isPublic ? 'green' : 'blue'}>
          {isPublic ? '公开' : '私有'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          <Button 
            type="primary" 
            size="small" 
            icon={<EyeOutlined />} 
            onClick={(e) => {
              e.stopPropagation();
              handlePreview(record.id);
            }}
          >
            查看
          </Button>
          <Popconfirm
            title="确定删除此文件吗？"
            onConfirm={(e) => {
              if (e) e.stopPropagation();
              handleDelete(record.id);
            }}
            onCancel={(e) => e?.stopPropagation()}
            okText="确定"
            cancelText="取消"
          >
            <Button 
              danger 
              size="small" 
              icon={<DeleteOutlined />} 
              onClick={(e) => e.stopPropagation()}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];
  
  // 上传组件属性配置
  const uploadProps: UploadProps = {
    name: 'file',
    multiple: false,
    beforeUpload: (file) => {
      const isLt10M = file.size / 1024 / 1024 < 10;
      if (!isLt10M) {
        message.error('文件大小不能超过10MB!');
      }
      return false;
    },
    onChange(info) {
      const { status } = info.file;
      if (status === 'removed') {
        form.setFieldsValue({ uploadFiles: info });
      }
    },
  };
  
  // 渲染网格视图项
  const renderGridItem = (file: FileInfo) => {
    const isSelected = selectedFiles.includes(file.id);
    const isImage = file.type.includes('image');
    
    return (
      <div 
        className={`file-grid-item ${isSelected ? 'file-grid-item-selected' : ''}`}
        onClick={() => toggleFileSelection(file.id)}
      >
        <div className="file-checkbox">
          <Checkbox 
            checked={isSelected}
            onChange={(e) => {
              e.stopPropagation();
              toggleFileSelection(file.id);
            }}
            onClick={(e) => e.stopPropagation()}
          />
        </div>
        
        <div className="file-preview-area" onClick={(e) => {
          e.stopPropagation();
          handlePreview(file.id);
        }}>
          {isImage ? (
            <div className="file-thumbnail">
              <img src={file.url} alt={file.name} />
            </div>
          ) : (
            <div className="file-icon-large">
              {getFileIcon(file.type, 'large')}
            </div>
          )}
          <div className="file-delete-hover">
            <Popconfirm
              title="确定删除此文件吗？"
              onConfirm={(e) => {
                if (e) e.stopPropagation();
                handleDelete(file.id);
              }}
              onCancel={(e) => e?.stopPropagation()}
              okText="确定"
              cancelText="取消"
            >
              <Button 
                type="primary" 
                danger 
                size="small" 
                shape="circle"
                icon={<DeleteOutlined />} 
                onClick={(e) => e.stopPropagation()}
              />
            </Popconfirm>
          </div>
        </div>
        
        <div className="file-info">
          <Tooltip title={file.name}>
            <div className="file-name">{file.name}</div>
          </Tooltip>
          <div className="file-meta">
            <span>{renderFileTypeTag(file.type)}</span>
            <span>{file.sizeFormat}</span>
          </div>
        </div>
      </div>
    );
  };
  
  return (
    <div className="file-management-container">
      <Card className="main-dashboard-card">
        <div className="user-title">
          <Title level={4}>文件管理</Title>
          <div className="toolbar">
            <Radio.Group
              value={viewMode}
              onChange={e => setViewMode(e.target.value)}
              optionType="button"
              buttonStyle="solid"
            >
              <Radio.Button value="grid"><AppstoreOutlined /> 网格视图</Radio.Button>
              <Radio.Button value="list"><BarsOutlined /> 列表视图</Radio.Button>
            </Radio.Group>

            <Input.Search
              placeholder="搜索文件名或上传者"
              allowClear
              enterButton={<><SearchOutlined />搜索</>}
              size="middle"
              onSearch={handleSearch}
              style={{ width: 320 }}
            />
            
            <Button
              onClick={toggleSelectAll}
            >
              {selectedFiles.length === files.length && files.length > 0 ? '取消全选' : '全选'}
            </Button>
            
            {selectedFiles.length > 0 && (
              <Popconfirm
                title={`确定删除选中的 ${selectedFiles.length} 个文件吗？`}
                onConfirm={handleBatchDelete}
                okText="确定"
                cancelText="取消"
              >
                <Button
                  danger
                  icon={<DeleteOutlined />}
                >
                  批量删除 ({selectedFiles.length})
                </Button>
              </Popconfirm>
            )}
            
            <Button
              type="primary"
              icon={<UploadOutlined />}
              onClick={() => setIsModalVisible(true)}
            >
              上传文件
            </Button>
          </div>
        </div>
        
        <div className="full-content-layout">
          {/* 类型和公开状态筛选区域 */}
          <div className="filter-area">
            <Space>
              <Select
                placeholder="文件类型"
                style={{ width: 120 }}
                allowClear
                onChange={(value) => {
                  setFileType(value);
                  setCurrent(1);
                }}
              >
                <Option value="image">图片</Option>
                <Option value="document">文档</Option>
                <Option value="spreadsheet">表格</Option>
                <Option value="pdf">PDF</Option>
                <Option value="compressed">压缩包</Option>
              </Select>
              
              <Select
                placeholder="共享状态"
                style={{ width: 120 }}
                allowClear
                onChange={(value) => {
                  setIsPublic(value === undefined ? undefined : value === true);
                  setCurrent(1);
                }}
              >
                <Option value={true}>公开</Option>
                <Option value={false}>私有</Option>
              </Select>
            </Space>
          </div>
          
          {/* 网格/列表视图 */}
          <div className="file-list-container">
            {viewMode === 'grid' ? (
              <Spin spinning={loading} tip="加载中...">
                <div className="file-grid" ref={fileGridRef}>
                  {gridFiles.length > 0 ? (
                    <>
                      <Row gutter={[16, 16]}>
                        {gridFiles.map(file => (
                          <Col xs={24} sm={12} md={8} lg={6} xl={4} xxl={3} key={file.id}>
                            {renderGridItem(file)}
                          </Col>
                        ))}
                      </Row>
                      {hasMore && (
                        <div ref={loadingRef} className="loading-more-container">
                          <div className="loading-more">
                            {loading ? '加载中...' : '向下滚动加载更多'}
                          </div>
                        </div>
                      )}
                    </>
                  ) : (
                    <Empty description="暂无文件" />
                  )}
                </div>
              </Spin>
            ) : (
              <Table
                columns={columns}
                dataSource={files}
                rowKey="id"
                pagination={pagination}
                loading={loading}
                onChange={(newPagination) => setPagination(newPagination)}
                onRow={handleRowClick}
                expandable={{
                  expandedRowRender: record => (
                    <p className="file-description">{record.path}</p>
                  ),
                }}
              />
            )}
          </div>
        </div>
        
        {/* 上传文件对话框 */}
        <Modal
          title="上传文件"
          open={isModalVisible}
          onOk={handleSubmit}
          onCancel={handleCancel}
          okText="上传"
          cancelText="取消"
          width={640}
        >
          <Form form={form} layout="vertical">
            <Form.Item
              name="uploadFiles"
              rules={[{ required: true, message: '请选择要上传的文件' }]}
            >
              <Dragger {...uploadProps}>
                <p className="ant-upload-drag-icon">
                  <UploadOutlined />
                </p>
                <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
                <p className="ant-upload-hint">
                  支持单个或批量上传，文件大小不超过10MB
                </p>
              </Dragger>
            </Form.Item>
          </Form>
        </Modal>
        
        {/* 文件预览对话框 */}
        <Modal
          title="文件预览"
          open={isPreviewVisible}
          onCancel={handleCancel}
          footer={[
            <Button key="download" icon={<DownloadOutlined />} onClick={() => previewFile && handleDownload(previewFile)}>
              下载
            </Button>,
            <Button 
              key="togglePublic" 
              onClick={() => previewFile && handleUpdatePublic(previewFile.id, !previewFile.isPublic)}
            >
              {previewFile?.isPublic ? '设为私有' : '设为公开'}
            </Button>,
            <Button key="delete" danger icon={<DeleteOutlined />} onClick={() => previewFile && handleDelete(previewFile.id)}>
              删除
            </Button>,
            <Button key="close" onClick={handleCancel}>
              关闭
            </Button>,
          ]}
          width={800}
          className="preview-modal"
        >
          {previewFile && (
            <div className="file-preview">
              <div className="file-preview-header">
                <Space align="center">
                  {getFileIcon(previewFile.type, 'large')}
                  <div>
                    <h3>{previewFile.name}</h3>
                    <p>
                      <span>类型：{previewFile.type}</span>
                      <span style={{ marginLeft: 12 }}>大小：{previewFile.sizeFormat}</span>
                      <span style={{ marginLeft: 12 }}>上传人：{previewFile.username}</span>
                      <span style={{ marginLeft: 12 }}>上传时间：{previewFile.createdAt}</span>
                    </p>
                  </div>
                </Space>
              </div>
              
              <div className="file-preview-content">
                {previewFile.type.includes('image') ? (
                  <div className="image-preview">
                    <img src={previewFile.url} alt={previewFile.name} />
                  </div>
                ) : (
                  <div className="file-preview-placeholder">
                    {getFileIcon(previewFile.type, 'large')}
                    <p>暂不支持此文件类型的预览</p>
                  </div>
                )}
              </div>
              
              <div className="file-preview-info">
                <Divider />
                <h4>文件信息</h4>
                <p>文件路径: {previewFile.path}</p>
                <p>内容类型: {previewFile.contentType}</p>
                <p>访问状态: {previewFile.isPublic ? '公开' : '私有'}</p>
                <p>访问链接: <a href={previewFile.url} target="_blank" rel="noopener noreferrer">{previewFile.url}</a></p>
              </div>
            </div>
          )}
        </Modal>
      </Card>
    </div>
  );
};

export default FileManagement; 