import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Form,
  Input,
  Switch,
  Button,
  message,
  Space,
  Alert,
  Table,
  Modal,
  InputNumber
} from 'antd';
import { 
  ReloadOutlined, 
  PlusOutlined,
  EditOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import './styles.css';
import FileSelector from '../../components/FileSelector';
import { getBottomTabList, updateBottomTabStatus, updateBottomTab, BottomTabItem } from '../../api/content';

const { Title } = Typography;

// 样式常量
const STYLES = {
  // 表格列宽
  columnWidths: {
    order: 80,
    icon: 110,
    selectedIcon: 100,
    action: 100
  },
  // 图标尺寸
  iconSize: {
    width: 24,
    height: 24
  },
  // 图片上传区域
  imageUploader: {
    container: {
      display: 'flex', 
      flexWrap: 'wrap' as const, 
      gap: '8px'
    },
    preview: {
      width: 104, 
      height: 104, 
      border: '1px dashed #d9d9d9', 
      padding: 8, 
      boxSizing: 'border-box' as const
    },
    uploadButton: {
      width: 104, 
      height: 104, 
      border: '1px dashed #d9d9d9', 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center',
      cursor: 'pointer',
      flexDirection: 'column' as const,
      backgroundColor: '#fafafa'
    },
    uploadIcon: {
      fontSize: 20, 
      color: '#999'
    },
    uploadText: {
      marginTop: 8, 
      color: '#666'
    }
  },
  // 其他样式
  alertMargin: {
    marginBottom: 24
  }
};

// Tab项目的最大数量限制
const MAX_TAB_COUNT = 5;

const TabSettings: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [tabList, setTabList] = useState<BottomTabItem[]>([]);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [currentTab, setCurrentTab] = useState<BottomTabItem | null>(null);
  const [editForm] = Form.useForm();
  
  // 文件选择器状态
  const [iconSelectorVisible, setIconSelectorVisible] = useState(false);
  const [selectedIconSelectorVisible, setSelectedIconSelectorVisible] = useState(false);

  // 加载现有配置
  useEffect(() => {
    fetchTabSettings();
  }, []);

  // 获取Tab设置
  const fetchTabSettings = async () => {
    setLoading(true);
    try {
      const response = await getBottomTabList();
      if (response.code === 0) {
        setTabList(response.data.list || []);
      } else {
        message.error(response.message || '获取底部Tab列表失败');
      }
    } catch (error) {
      console.error('获取Tab设置失败:', error);
      message.error('获取Tab设置失败');
    } finally {
      setLoading(false);
    }
  };

  // 重置表单
  const handleReset = () => {
    fetchTabSettings();
    message.info('Tab设置已重置');
  };
  
  // 创建文件对象
  const createFileObject = (url: string, fileType: string) => {
    if (!url) return null;
    
    return {
      uid: '-1',
      name: url.split('/').pop() || fileType,
      status: 'done',
      url,
      thumbUrl: url,
    };
  };
  
  // 编辑Tab
  const handleEditTab = (tab: BottomTabItem) => {
    setCurrentTab(tab);
    
    // 创建文件对象用于图标显示
    const mockFile = createFileObject(tab.icon, '图标');
    const mockSelectedFile = createFileObject(tab.selectedIcon, '选中图标');
    
    editForm.setFieldsValue({
      name: tab.name,
      isEnabled: tab.isEnabled,
      order: tab.order,
      path: tab.path || '',
      icon: mockFile ? { fileList: [mockFile] } : undefined,
      selectedIcon: mockSelectedFile ? { fileList: [mockSelectedFile] } : undefined
    });
    
    setEditModalVisible(true);
  };
  
  // 获取表单中的文件URL
  const getFileUrl = (fileList: any, defaultUrl: string = '') => {
    if (fileList?.fileList?.[0]?.url) {
      return fileList.fileList[0].url;
    }
    
    if (fileList?.fileList?.[0]?.response?.url) {
      return fileList.fileList[0].response.url;
    }
    
    return defaultUrl;
  };
  
  // 保存Tab编辑
  const handleSaveTabEdit = async () => {
    try {
      const values = await editForm.validateFields();
      setLoading(true);
      
      if (currentTab) {
        const updateData = {
          id: currentTab.id,
          name: values.name,
          path: values.path,
          icon: getFileUrl(values.icon, currentTab.icon),
          selectedIcon: getFileUrl(values.selectedIcon, currentTab.selectedIcon),
          order: values.order,
          isEnabled: values.isEnabled
        };
        
        try {
          const response = await updateBottomTab(updateData);
          if (response.code === 0) {
            message.success('Tab编辑成功');
            fetchTabSettings(); // 重新获取列表
          } else {
            message.error(response.message || '更新Tab失败');
          }
        } catch (error) {
          console.error('更新Tab失败:', error);
          message.error('编辑失败，请重试');
        }
        
        setEditModalVisible(false);
      }
    } catch (error) {
      message.error('请填写正确的Tab信息');
    } finally {
      setLoading(false);
    }
  };
  
  // 启用/禁用Tab
  const toggleTabActive = async (id: number, isEnabled: boolean) => {
    setLoading(true);
    try {
      const response = await updateBottomTabStatus(id, isEnabled);
      if (response.code === 0) {
        message.success(`Tab已${isEnabled ? '启用' : '禁用'}`);
        // 更新本地列表状态
        setTabList(prevList => 
          prevList.map(tab => 
            tab.id === id ? { ...tab, isEnabled } : tab
          )
        );
      } else {
        message.error(response.message || '更新状态失败');
      }
    } catch (error) {
      console.error('更新Tab状态失败:', error);
      message.error('操作失败，请重试');
    } finally {
      setLoading(false);
    }
  };
  
  // 处理文件选择
  const handleFileSelected = (url: string, fieldName: string, setVisible: React.Dispatch<React.SetStateAction<boolean>>) => {
    const mockFile = createFileObject(url, '图片');
    
    // 设置表单字段值
    editForm.setFieldsValue({
      [fieldName]: mockFile ? { fileList: [mockFile] } : undefined
    });
    
    setVisible(false);
  };
  
  // 处理未选中图标选择
  const handleIconSelected = (url: string) => {
    handleFileSelected(url, 'icon', setIconSelectorVisible);
  };
  
  // 处理选中图标选择
  const handleSelectedIconSelected = (url: string) => {
    handleFileSelected(url, 'selectedIcon', setSelectedIconSelectorVisible);
  };
  
  // 渲染图标
  const renderIcon = (iconUrl: string, altText: string) => (
    <div style={{ textAlign: 'center' }}>
      {iconUrl ? (
        <img 
          src={iconUrl} 
          alt={altText} 
          style={STYLES.iconSize} 
        />
      ) : (
        <span>无图标</span>
      )}
    </div>
  );
  
  // Tab列表列定义
  const columns: ColumnsType<BottomTabItem> = [
    {
      title: '排序',
      dataIndex: 'order',
      key: 'order',
      width: STYLES.columnWidths.order,
      sorter: (a, b) => a.order - b.order,
    },
    {
      title: 'Tab名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '未选中图标',
      dataIndex: 'icon',
      key: 'icon',
      width: STYLES.columnWidths.icon,
      render: (icon) => renderIcon(icon, '图标'),
    },
    {
      title: '选中图标',
      dataIndex: 'selectedIcon',
      key: 'selectedIcon',
      width: STYLES.columnWidths.selectedIcon,
      render: (selectedIcon) => renderIcon(selectedIcon, '选中图标'),
    },
    {
      title: '页面路径',
      dataIndex: 'path',
      key: 'path',
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      key: 'isEnabled',
      render: (isEnabled, record) => (
        <Switch 
          checked={isEnabled} 
          onChange={(checked) => toggleTabActive(record.id, checked)}
          checkedChildren="启用" 
          unCheckedChildren="禁用" 
        />
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: STYLES.columnWidths.action,
      render: (_, record) => (
        <Button
          type="primary"
          size="small"
          icon={<EditOutlined />}
          onClick={() => handleEditTab(record)}
        >
          编辑
        </Button>
      ),
    },
  ];

  // 渲染文件预览和上传按钮
  const renderFileUploader = (fieldName: string, label: string, setVisible: React.Dispatch<React.SetStateAction<boolean>>) => {
    const fileUrl = editForm.getFieldValue(fieldName)?.fileList?.[0]?.url || 
                   (currentTab && fieldName === 'icon' ? currentTab.icon : '') ||
                   (currentTab && fieldName === 'selectedIcon' ? currentTab.selectedIcon : '');
    
    return (
      <Form.Item
        name={fieldName}
        label={label}
        rules={[{ required: true, message: `请上传${label}` }]}
      >
        <div style={STYLES.imageUploader.container}>
          {/* 已上传图片 */}
          {fileUrl && (
            <div style={STYLES.imageUploader.preview}>
              <img 
                src={fileUrl} 
                alt={label}
                style={{ width: '100%', height: '100%', objectFit: 'contain' }} 
              />
            </div>
          )}
          
          {/* 上传按钮 */}
          <div 
            style={STYLES.imageUploader.uploadButton}
            onClick={() => setVisible(true)}
          >
            <PlusOutlined style={STYLES.imageUploader.uploadIcon} />
            <div style={STYLES.imageUploader.uploadText}>上传</div>
          </div>
        </div>
      </Form.Item>
    );
  };

  return (
    <div className="mini-program-container">
      <Card className="main-dashboard-card">
        <div className="category-header">
          <Title level={4}>微信小程序底部设置</Title>
          <Space>
            <Button 
              icon={<ReloadOutlined />}
              onClick={handleReset}
            >
              重置
            </Button>
          </Space>
        </div>
        
        <Alert
          message="微信小程序底部导航配置"
          description={`微信小程序底部最多支持${MAX_TAB_COUNT}个导航项，每个导航项包含图标和文字。请设置每个导航项的名称、未选中图标、选中图标和对应的页面路径。`}
          type="info"
          showIcon
          style={STYLES.alertMargin}
        />
        
        <Table
          columns={columns}
          dataSource={tabList}
          rowKey="id"
          pagination={false}
          bordered
          loading={loading}
        />
      </Card>
      
      {/* Tab编辑弹窗 */}
      <Modal
        title="编辑微信小程序Tab"
        open={editModalVisible}
        onOk={handleSaveTabEdit}
        onCancel={() => {
          setEditModalVisible(false);
          setIconSelectorVisible(false);
          setSelectedIconSelectorVisible(false);
        }}
        okText="保存"
        cancelText="取消"
      >
        <Form form={editForm} layout="vertical">
          <Form.Item
            name="name"
            label="Tab名称"
            rules={[{ required: true, message: '请输入Tab名称' }]}
          >
            <Input placeholder="请输入Tab名称" />
          </Form.Item>
          
          <Form.Item
            name="order"
            label="Tab排序"
            rules={[
              { required: true, message: '请输入Tab排序' },
              { type: 'number', message: '请输入有效的数字' }
            ]}
          >
            <InputNumber 
              style={{ width: '100%' }} 
              min={1}
              max={MAX_TAB_COUNT}
              placeholder="数字越小排序越靠前"
            />
          </Form.Item>
          
          <Form.Item
            name="path"
            label="页面路径"
            rules={[{ required: true, message: '请输入页面路径' }]}
          >
            <Input placeholder="请输入页面路径，例如：pages/index/index" />
          </Form.Item>
          
          {renderFileUploader('icon', '未选中图标', setIconSelectorVisible)}
          {renderFileUploader('selectedIcon', '选中图标', setSelectedIconSelectorVisible)}
          
          <Form.Item
            name="isEnabled"
            label="是否启用"
            valuePropName="checked"
          >
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
        </Form>
      </Modal>
      
      {/* 未选中图标选择器 */}
      {iconSelectorVisible && (
        <FileSelector
          visible={iconSelectorVisible}
          onCancel={() => setIconSelectorVisible(false)}
          onSelect={handleIconSelected}
          title="选择未选中Tab图标"
        />
      )}
      
      {/* 选中图标选择器 */}
      {selectedIconSelectorVisible && (
        <FileSelector
          visible={selectedIconSelectorVisible}
          onCancel={() => setSelectedIconSelectorVisible(false)}
          onSelect={handleSelectedIconSelected}
          title="选择选中Tab图标"
        />
      )}
    </div>
  );
};

export default TabSettings; 