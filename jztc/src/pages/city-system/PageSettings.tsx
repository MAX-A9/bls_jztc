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
  InputNumber,
  Alert,
  Tabs,
  Select,
  Table,
  Avatar,
  Popconfirm,
  Modal,
  Image,
  Tag,
  Tooltip
} from 'antd';
import { 
  ReloadOutlined, 
  PlusOutlined, 
  DeleteOutlined,
  EditOutlined
} from '@ant-design/icons';
import type { TabsProps } from 'antd';
import './styles.css';
import FileSelector from '../../components/FileSelector';
import { 
  createInnerBanner, 
  updateInnerBanner, 
  getInnerBannerList, 
  deleteInnerBanner, 
  updateInnerBannerStatus, 
  updateHomeInnerBannerGlobalStatus,
  updateIdleInnerBannerGlobalStatus,
  InnerBannerItem 
} from '../../api/content';

const { Title } = Typography;
const { Option } = Select;

const PageSettings: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('home');
  
  // 首页banner数据
  const [homeBannerList, setHomeBannerList] = useState<InnerBannerItem[]>([]);
  
  // 闲置banner数据
  const [idleBannerList, setIdleBannerList] = useState<InnerBannerItem[]>([]);
  
  // 添加/编辑banner相关状态
  const [bannerModalVisible, setBannerModalVisible] = useState(false);
  const [currentBanner, setCurrentBanner] = useState<InnerBannerItem | null>(null);
  const [currentBannerType, setCurrentBannerType] = useState<'home' | 'idle'>('home');
  const [bannerForm] = Form.useForm();
  
  // 文件选择器状态
  const [bannerSelectorVisible, setBannerSelectorVisible] = useState(false);
  
  // 添加banner
  const handleAddBanner = (type: 'home' | 'idle') => {
    setCurrentBanner(null);
    setCurrentBannerType(type);
    bannerForm.resetFields();
    setBannerModalVisible(true);
  };
  
  // 编辑banner
  const handleEditBanner = (banner: InnerBannerItem, type: 'home' | 'idle') => {
    setCurrentBanner(banner);
    setCurrentBannerType(type);
    
    // 创建一个模拟的文件对象用于图片显示
    const mockFile = banner.image ? {
      uid: '-1',
      name: banner.image.split('/').pop() || '图片',
      status: 'done',
      url: banner.image,
      thumbUrl: banner.image,
    } : null;
    
    bannerForm.setFieldsValue({
      linkType: banner.linkType,
      linkUrl: banner.linkUrl,
      isEnabled: banner.isEnabled,
      order: banner.order,
      image: mockFile ? { fileList: [mockFile] } : undefined
    });
    
    setBannerModalVisible(true);
  };
  
  // 保存banner
  const handleSaveBanner = async () => {
    try {
      const values = await bannerForm.validateFields();
      setLoading(true);
      
      const imageUrl = values.image?.fileList?.[0]?.url || (currentBanner?.image || '');
      if (!imageUrl) {
        message.error('请上传Banner图片');
        setLoading(false);
        return;
      }
      
      if (currentBanner) {
        // 编辑现有Banner
        const updateData = {
          id: currentBanner.id,
          bannerType: currentBannerType,
          image: imageUrl,
          linkType: values.linkType,
          linkUrl: values.linkUrl,
          isEnabled: !!values.isEnabled,
          order: values.order || 0
        };
        
        try {
          const response = await updateInnerBanner(updateData);
          if (response.code === 0) {
            message.success(`编辑${currentBannerType === 'home' ? '首页' : '闲置'}Banner成功`);
            // 刷新列表
            fetchPageSettings();
          } else {
            message.error(response.message || '更新Banner失败');
          }
        } catch (error) {
          console.error('更新Banner失败:', error);
          message.error('更新Banner失败，请重试');
        }
      } else {
        // 添加新Banner
        const createData = {
          bannerType: currentBannerType,
          image: imageUrl,
          linkType: values.linkType,
          linkUrl: values.linkUrl,
          isEnabled: !!values.isEnabled,
          order: values.order || 0
        };
        
        try {
          const response = await createInnerBanner(createData);
          if (response.code === 0) {
            message.success(`添加${currentBannerType === 'home' ? '首页' : '闲置'}Banner成功`);
            // 刷新列表
            fetchPageSettings();
          } else {
            message.error(response.message || '添加Banner失败');
          }
        } catch (error) {
          console.error('添加Banner失败:', error);
          message.error('添加Banner失败，请重试');
        }
      }
      
      setBannerModalVisible(false);
    } catch (error) {
      console.error('保存Banner失败:', error);
      message.error('表单验证失败，请检查输入');
    } finally {
      setLoading(false);
    }
  };
  
  // 删除banner
  const handleDeleteBanner = async (id: number, type: 'home' | 'idle') => {
    setLoading(true);
    try {
      const response = await deleteInnerBanner(id);
      if (response.code === 0) {
        message.success(`删除${type === 'home' ? '首页' : '闲置'}Banner成功`);
        // 刷新列表
        fetchPageSettings();
      } else {
        message.error(response.message || '删除Banner失败');
      }
    } catch (error) {
      console.error('删除Banner失败:', error);
      message.error('删除操作失败，请重试');
    } finally {
      setLoading(false);
    }
  };
  
  // 切换banner状态
  const toggleBannerStatus = async (id: number, isEnabled: boolean, type: 'home' | 'idle') => {
    setLoading(true);
    try {
      const response = await updateInnerBannerStatus(id, isEnabled);
      if (response.code === 0) {
        message.success(`${isEnabled ? '启用' : '禁用'}${type === 'home' ? '首页' : '闲置'}Banner成功`);
        // 刷新列表
        fetchPageSettings();
      } else {
        message.error(response.message || '更新状态失败');
      }
    } catch (error) {
      console.error('更新Banner状态失败:', error);
      message.error('操作失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 处理切换Banner总开关
  const handleGlobalSwitchChange = async (type: 'home' | 'idle', checked: boolean) => {
    setLoading(true);
    try {
      // 根据类型调用不同的API
      const response = type === 'home' 
        ? await updateHomeInnerBannerGlobalStatus(checked)
        : await updateIdleInnerBannerGlobalStatus(checked);
        
      if (response.code === 0) {
        message.success(`${checked ? '启用' : '禁用'}${type === 'home' ? '首页' : '闲置'}Banner成功`);
        // 更新表单状态
        form.setFieldsValue({
          [type === 'home' ? 'enableHomeBanner' : 'enableIdleBanner']: checked
        });
      } else {
        message.error(response.message || '更新总开关状态失败');
        // 更新失败时回滚UI状态
        form.setFieldsValue({
          [type === 'home' ? 'enableHomeBanner' : 'enableIdleBanner']: !checked
        });
      }
    } catch (error) {
      console.error('更新Banner总开关状态失败:', error);
      message.error('操作失败，请重试');
      // 更新失败时回滚UI状态
      form.setFieldsValue({
        [type === 'home' ? 'enableHomeBanner' : 'enableIdleBanner']: !checked
      });
    } finally {
      setLoading(false);
    }
  };

  // 加载现有配置
  useEffect(() => {
    fetchPageSettings();
  }, []);

  // 获取页面设置
  const fetchPageSettings = async () => {
    setLoading(true);
    try {
      // 获取首页Banner列表
      const homeResponse = await getInnerBannerList('home');
      if (homeResponse.code === 0) {
        setHomeBannerList(homeResponse.data.list || []);
        
        // 设置首页Banner总开关状态
        form.setFieldsValue({
          enableHomeBanner: homeResponse.data.isGlobalEnabled
        });
      } else {
        message.error('获取首页Banner列表失败：' + homeResponse.message);
      }
      
      // 获取闲置Banner列表
      const idleResponse = await getInnerBannerList('idle');
      if (idleResponse.code === 0) {
        setIdleBannerList(idleResponse.data.list || []);
        
        // 设置闲置Banner总开关状态
        form.setFieldsValue({
          enableIdleBanner: idleResponse.data.isGlobalEnabled
        });
      } else {
        message.error('获取闲置Banner列表失败：' + idleResponse.message);
      }
    } catch (error) {
      console.error('获取内页设置失败:', error);
      message.error('获取内页设置失败');
    } finally {
      setLoading(false);
    }
  };

  // 重置表单
  const handleReset = () => {
    form.resetFields();
    fetchPageSettings();
    message.info('表单已重置');
  };
  
  // 渲染Banner表格
  const renderBannerTable = (type: 'home' | 'idle') => {
    const dataSource = type === 'home' ? homeBannerList : idleBannerList;
    
    return (
      <Table
        dataSource={dataSource}
        rowKey="id"
        pagination={false}
        columns={[
          {
            title: '排序',
            dataIndex: 'order',
            key: 'order',
            width: 80,
            render: (order, _, index) => (
              <div style={{ textAlign: 'center' }}>{order || index + 1}</div>
            ),
          },
          {
            title: '预览图',
            dataIndex: 'image',
            key: 'image',
            width: 100,
            render: image => (
              <div style={{ textAlign: 'center' }}>
                {image ? (
                  <Image 
                    src={image} 
                    alt="Banner图" 
                    width={80}
                    height={45}
                    style={{ objectFit: 'cover' }}
                  />
                ) : (
                  <Avatar size={40} shape="square" style={{ backgroundColor: '#87d068' }}>
                    图
                  </Avatar>
                )}
              </div>
            ),
          },
          {
            title: '跳转类型',
            dataIndex: 'linkType',
            key: 'linkType',
            render: (linkType) => {
              let tagColor = 'blue';
              let linkText = '小程序页面';
              
              if (linkType === 'miniprogram') {
                tagColor = 'green';
                linkText = '其他小程序';
              } else if (linkType === 'webview') {
                tagColor = 'orange';
                linkText = '网页';
              }
              
              return <Tag color={tagColor}>{linkText}</Tag>;
            }
          },
          {
            title: '跳转地址',
            dataIndex: 'linkUrl',
            key: 'linkUrl',
            ellipsis: true,
            render: (linkUrl) => (
              <Tooltip title={linkUrl} placement="topLeft">
                <span>{linkUrl || '-'}</span>
              </Tooltip>
            )
          },
          {
            title: '状态',
            dataIndex: 'isEnabled',
            key: 'isEnabled',
            render: (isEnabled, record) => (
              <Switch
                checked={isEnabled}
                onChange={(checked) => toggleBannerStatus(record.id, checked, type)}
                checkedChildren="启用"
                unCheckedChildren="禁用"
              />
            ),
          },
          {
            title: '操作',
            key: 'action',
            width: 180,
            render: (_, record) => (
              <Space size="small">
                <Button 
                  type="primary" 
                  size="small"
                  icon={<EditOutlined />} 
                  onClick={() => handleEditBanner(record, type)}
                >
                  编辑
                </Button>
                <Popconfirm
                  title={`确定要删除这个${type === 'home' ? '首页' : '闲置'}Banner吗？`}
                  onConfirm={() => handleDeleteBanner(record.id, type)}
                  okText="是"
                  cancelText="否"
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
        ]}
      />
    );
  };

  // 处理文件选择 - Banner图片
  const handleBannerImageSelected = (url: string) => {
    // 创建一个模拟的文件对象
    const mockFile = {
      uid: '-1',
      name: url.split('/').pop() || '图片',
      status: 'done',
      url: url,
      thumbUrl: url,
    };
    
    // 设置bannerForm的image字段值
    bannerForm.setFieldsValue({
      image: { fileList: [mockFile] }
    });
    
    setBannerSelectorVisible(false);
  };

  // 标签页配置
  const items: TabsProps['items'] = [
    {
      key: 'home',
      label: '首页Banner设置',
      children: (
        <>
          <Alert
            message="首页Banner设置"
            description="配置首页轮播图，可以添加多张Banner图，设置跳转链接。"
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Form.Item
              name="enableHomeBanner"
              label="启用首页Banner"
              valuePropName="checked"
              style={{ marginBottom: 0 }}
            >
              <Switch 
                checkedChildren="已启用" 
                unCheckedChildren="已禁用" 
                onChange={(checked) => handleGlobalSwitchChange('home', checked)}
              />
            </Form.Item>
            
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => handleAddBanner('home')}
            >
              添加Banner
            </Button>
          </div>
          
          {renderBannerTable('home')}
        </>
      ),
    },
    {
      key: 'idle',
      label: '闲置Banner设置',
      children: (
        <>
          <Alert
            message="闲置Banner设置"
            description="配置闲置轮播图，可以添加多张Banner图，设置跳转链接。"
            type="info"
            showIcon
            style={{ marginBottom: 24 }}
          />
          
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Form.Item
              name="enableIdleBanner"
              label="启用闲置Banner"
              valuePropName="checked"
              style={{ marginBottom: 0 }}
            >
              <Switch 
                checkedChildren="已启用" 
                unCheckedChildren="已禁用" 
                onChange={(checked) => handleGlobalSwitchChange('idle', checked)}
              />
            </Form.Item>
            
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => handleAddBanner('idle')}
            >
              添加Banner
            </Button>
          </div>
          
          {renderBannerTable('idle')}
        </>
      ),
    },
  ];

  return (
    <>
      <div className="mini-program-container">
        <Card className="main-dashboard-card">
          <div className="category-header">
            <Title level={4}>内页设置</Title>
            <Space>
              <Button 
                icon={<ReloadOutlined />}
                onClick={handleReset}
              >
                重置
              </Button>
            </Space>
          </div>
          
          <Form
            form={form}
            layout="vertical"
            disabled={loading}
          >
            <Tabs
              activeKey={activeTab}
              onChange={setActiveTab}
              items={items}
              className="mp-settings-tabs"
            />
          </Form>
        </Card>
      </div>
      
      {/* Banner编辑弹窗 */}
      <Modal
        title={currentBanner 
          ? `编辑${currentBannerType === 'home' ? '首页' : '闲置'}Banner` 
          : `添加${currentBannerType === 'home' ? '首页' : '闲置'}Banner`}
        open={bannerModalVisible}
        onOk={handleSaveBanner}
        onCancel={() => {
          setBannerModalVisible(false);
          setBannerSelectorVisible(false);
        }}
        okText="保存"
        cancelText="取消"
      >
        <Form form={bannerForm} layout="vertical">
          <Form.Item
            name="order"
            label="排序"
            rules={[
              { required: true, message: '请输入排序数字' },
              { type: 'number', message: '请输入有效的数字' }
            ]}
            tooltip="数字越小排序越靠前"
          >
            <InputNumber 
              style={{ width: '100%' }} 
              min={1} 
              precision={0}
              placeholder="请输入排序值，数字越小排序越靠前"
            />
          </Form.Item>
          
          <Form.Item
            name="image"
            label="Banner图片"
          >
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
              {/* 已上传图片 */}
              {(bannerForm.getFieldValue('image')?.fileList?.[0]?.url || currentBanner?.image) && (
                <div style={{ width: 200, height: 120, border: '1px dashed #d9d9d9', padding: 8, boxSizing: 'border-box' }}>
                  <img 
                    src={bannerForm.getFieldValue('image')?.fileList?.[0]?.url || currentBanner?.image} 
                    alt="Banner图片" 
                    style={{ width: '100%', height: '100%', objectFit: 'contain' }} 
                  />
                </div>
              )}
              
              {/* 上传按钮 - 简洁样式 */}
              <div 
                style={{ 
                  width: 200, 
                  height: 120, 
                  border: '1px dashed #d9d9d9', 
                  display: 'flex', 
                  justifyContent: 'center', 
                  alignItems: 'center',
                  cursor: 'pointer',
                  flexDirection: 'column',
                  backgroundColor: '#fafafa'
                }}
                onClick={() => setBannerSelectorVisible(true)}
              >
                <PlusOutlined style={{ fontSize: 20, color: '#999' }} />
                <div style={{ marginTop: 8, color: '#666' }}>上传</div>
              </div>
            </div>
          </Form.Item>
          
          <Form.Item
            name="linkType"
            label="跳转类型"
            rules={[{ required: true, message: '请选择跳转类型' }]}
          >
            <Select placeholder="请选择跳转类型">
              <Option value="page">小程序页面</Option>
              <Option value="miniprogram">其他小程序</Option>
              <Option value="webview">网页</Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="linkUrl"
            label="跳转地址"
            rules={[{ required: true, message: '请输入跳转地址' }]}
          >
            <Input placeholder="请输入跳转地址" />
          </Form.Item>
          
          <Form.Item
            name="isEnabled"
            label="是否启用"
            valuePropName="checked"
          >
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
        </Form>
      </Modal>
      
      {/* 轮播图选择器 */}
      {bannerSelectorVisible && (
        <FileSelector
          visible={bannerSelectorVisible}
          onCancel={() => setBannerSelectorVisible(false)}
          onSelect={handleBannerImageSelected}
        />
      )}
    </>
  );
};

export default PageSettings; 