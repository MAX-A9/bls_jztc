import React, { useState, useEffect, useRef } from 'react';
import { 
  Modal, 
  Tabs, 
  Upload, 
  Button, 
  message, 
  Spin, 
  Empty,
  Form
} from 'antd';
import type { TabsProps, UploadProps, UploadFile } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import { getFileList, uploadFile, type FileInfo } from '@/api/file';
import './styles.css';

const { Dragger } = Upload;

interface FileSelectorProps {
  visible: boolean;
  onCancel: () => void;
  onSelect: (url: string) => void;
  title?: string;
  accept?: string;
}

const FileSelector: React.FC<FileSelectorProps> = ({
  visible,
  onCancel,
  onSelect,
  title = '选择文件',
  accept = 'image/*'
}) => {
  // 文件列表
  const [fileList, setFileList] = useState<FileInfo[]>([]);
  // 加载状态
  const [loading, setLoading] = useState<boolean>(false);
  // 上传状态
  const [uploading, setUploading] = useState<boolean>(false);
  // 激活的标签页
  const [activeTab, setActiveTab] = useState<string>('media');
  // 上传的文件
  const [uploadFileList, setUploadFileList] = useState<UploadFile[]>([]);
  // 分页参数
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const pageSize = 20; // 每页数量
  
  // 媒体库容器引用
  const mediaLibraryRef = useRef<HTMLDivElement>(null);
  
  // 是否已经初始化
  const initialized = useRef(false);
  
  // 是否禁用自动加载
  const disableAutoLoad = useRef(true);

  // 获取文件列表
  const fetchFiles = async (pageNum = 1, append = false) => {
    if (!hasMore && pageNum > 1) return;
    
    setLoading(true);
    try {
      const res = await getFileList({
        Page: pageNum,
        PageSize: pageSize,
        Type: accept.includes('image') ? 'image' : undefined
      });
      
      if (res.code === 0) {
        const newData = res.data.list || [];
        if (append) {
          setFileList(prev => [...prev, ...newData]);
        } else {
          setFileList(newData);
        }
        
        // 判断是否还有更多数据
        const totalCount = res.data.total || 0;
        const loadedCount = append ? fileList.length + newData.length : newData.length;
        const hasMoreData = newData.length === pageSize && loadedCount < totalCount;
        
        setHasMore(hasMoreData);
        setPage(pageNum);
        
        // 第一页加载完成后启用滚动加载
        if (!append) {
          initialized.current = true;
          // 立即启用自动加载，不再使用延迟
          disableAutoLoad.current = false;
        }
      } else {
        message.error(res.message || '获取文件列表失败');
      }
    } catch (error) {
      console.error('获取文件列表失败:', error);
      message.error('获取文件列表失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 处理媒体库滚动
  const handleScroll = () => {
    if (!mediaLibraryRef.current || loading || !hasMore || disableAutoLoad.current) return;
    
    const { scrollTop, scrollHeight, clientHeight } = mediaLibraryRef.current;
    // 当滚动到底部附近时加载更多
    if (scrollHeight - scrollTop - clientHeight < 100) {
      fetchFiles(page + 1, true);
    }
  };

  // 监听滚动事件
  useEffect(() => {
    const mediaLibrary = mediaLibraryRef.current;
    if (mediaLibrary && visible && activeTab === 'media') {
      mediaLibrary.addEventListener('scroll', handleScroll);
    }
    
    return () => {
      if (mediaLibrary) {
        mediaLibrary.removeEventListener('scroll', handleScroll);
      }
    };
  }, [visible, activeTab, page, hasMore, loading]);

  // 当弹窗打开时获取文件列表
  useEffect(() => {
    if (visible) {
      // 重置状态
      setPage(1);
      setHasMore(true);
      setFileList([]);
      initialized.current = false;
      disableAutoLoad.current = true;
      
      // 重置滚动位置
      if (mediaLibraryRef.current) {
        mediaLibraryRef.current.scrollTop = 0;
      }
      
      // 获取第一页数据
      fetchFiles(1, false);
    }
  }, [visible]);
  
  // 标签页切换
  useEffect(() => {
    if (activeTab === 'media' && visible) {
      // 切换到媒体标签时，立即启用自动加载
      disableAutoLoad.current = false;
    }
  }, [activeTab]);

  // 关闭弹窗时重置状态
  const handleCancel = () => {
    setUploadFileList([]);
    onCancel();
  };

  // 选择文件
  const handleSelectFile = (file: FileInfo) => {
    onSelect(file.url);
    handleCancel();
  };

  // 手动加载更多
  const handleLoadMore = () => {
    if (!loading && hasMore) {
      fetchFiles(page + 1, true);
    }
  };

  // 上传文件
  const handleUpload = async () => {
    if (uploadFileList.length === 0) {
      message.warning('请先选择要上传的文件');
      return;
    }

    const file = uploadFileList[0].originFileObj;
    if (!file) {
      message.warning('文件无效');
      return;
    }

    setUploading(true);
    try {
      const formData = new FormData();
      formData.append('File', file);
      
      const res = await uploadFile(formData);
      if (res.code === 0) {
        message.success('文件上传成功');
        onSelect(res.data.url);
        handleCancel();
      } else {
        message.error(res.message || '文件上传失败');
      }
    } catch (error) {
      console.error('上传文件失败:', error);
      message.error('上传文件失败');
    } finally {
      setUploading(false);
    }
  };

  // 上传属性设置
  const uploadProps: UploadProps = {
    name: 'file',
    multiple: false,
    accept,
    maxCount: 1,
    fileList: uploadFileList,
    beforeUpload: (file) => {
      // 检查文件类型
      const acceptTypes = accept.split(',').map(type => type.trim());
      const isAcceptable = acceptTypes.some(type => {
        if (type === 'image/*') return file.type.startsWith('image/');
        return file.type === type;
      });
      
      if (!isAcceptable) {
        message.error(`只能上传${accept.includes('image') ? '图片' : '指定类型的'}文件!`);
        return Upload.LIST_IGNORE;
      }
      
      // 检查文件大小 (10MB)
      const isLt10M = file.size / 1024 / 1024 < 10;
      if (!isLt10M) {
        message.error('文件大小不能超过10MB!');
        return Upload.LIST_IGNORE;
      }
      
      // 更新上传列表但不自动上传
      setUploadFileList([{ uid: '-1', name: file.name, status: 'done', size: file.size, type: file.type, originFileObj: file }]);
      return false;
    },
    onRemove: () => {
      setUploadFileList([]);
      return true;
    }
  };

  // 标签页设置
  const tabItems: TabsProps['items'] = [
    {
      key: 'media',
      label: '媒体库',
      children: (
        <div className="media-library" ref={mediaLibraryRef}>
          {fileList.length > 0 ? (
            <div className="file-selector-grid">
              {fileList.map(file => (
                <div 
                  key={file.id} 
                  onClick={() => handleSelectFile(file)}
                  className="media-item"
                >
                  <div className="media-item-thumbnail">
                    <img 
                      src={file.url} 
                      alt={file.name} 
                      loading="lazy"
                    />
                  </div>
                  <div className="media-item-name">
                    {file.name}
                  </div>
                </div>
              ))}
            </div>
          ) : !loading ? (
            <Empty 
              description="暂无图片文件" 
              className="empty-placeholder"
            />
          ) : null}
          
          {/* 加载状态和提示 */}
          <div className="load-more-indicator">
            {loading && (
              <div className="loading-more">
                <Spin size="small" />
                <span>加载中...</span>
              </div>
            )}
            {!loading && hasMore && (
              <div className="scroll-tip">
                <Button type="link" onClick={handleLoadMore}>加载更多</Button>
              </div>
            )}
            {!hasMore && fileList.length > 0 && (
              <div className="no-more-data">
                没有更多文件了
              </div>
            )}
          </div>
        </div>
      ),
    },
    {
      key: 'upload',
      label: '上传文件',
      children: (
        <div>
          <Form layout="vertical">
            <Form.Item
              name="uploadFiles"
              rules={[{ required: true, message: '请选择要上传的图片' }]}
            >
              <Dragger {...uploadProps}>
                <p className="ant-upload-drag-icon">
                  <UploadOutlined />
                </p>
                <p className="ant-upload-text">点击或拖拽图片到此区域上传</p>
                <p className="ant-upload-hint">
                  支持单个图片上传，文件大小不超过10MB
                </p>
              </Dragger>
            </Form.Item>
            <Form.Item className="upload-button-container">
              <Button 
                type="primary" 
                onClick={handleUpload}
                disabled={uploading || uploadFileList.length === 0}
                loading={uploading}
              >
                {uploading ? '上传中...' : '开始上传'}
              </Button>
            </Form.Item>
          </Form>
        </div>
      ),
    },
  ];

  return (
    <Modal
      title={title}
      open={visible}
      onCancel={handleCancel}
      footer={null}
      width={900}
      styles={{ body: { padding: 0 } }}
      destroyOnClose
    >
      <Tabs 
        items={tabItems} 
        activeKey={activeTab}
        onChange={setActiveTab}
        className="file-selector-tabs"
      />
    </Modal>
  );
};

export default FileSelector; 