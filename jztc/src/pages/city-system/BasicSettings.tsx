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
  Tabs,
  Modal,
  Select
} from 'antd';
import { 
  ReloadOutlined, 
  PlusOutlined
} from '@ant-design/icons';
import type { TabsProps } from 'antd';
import './styles.css';
import { NavigationModule, ActivityModule, BannerModule } from './modules';
import FileSelector from '@/components/FileSelector';
import {
  getMiniProgramList,
  createMiniProgram,
  updateMiniProgram,
  deleteMiniProgram,
  updateMiniProgramStatus,
  updateMiniProgramGlobalStatus,
  getActivityArea,
  saveActivityArea,
  updateActivityAreaGlobalStatus,
  getBannerList,
  createBanner,
  updateBanner,
  deleteBanner,
  updateBannerStatus,
  updateBannerGlobalStatus,
  type MiniProgramItem,
  type BannerItem,
  type BannerCreateReq,
  type BannerUpdateReq,
  type BannerStatusUpdateReq,
  type MiniProgramGlobalStatusUpdateReq,
  type ActivityAreaGlobalStatusUpdateReq,
  type BannerGlobalStatusUpdateReq
} from '@/api/settings';

const { Title } = Typography;

const BasicSettings: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('function');
  
  // 轮播图列表数据
  const [bannerList, setBannerList] = useState<BannerItem[]>([]);
  
  // 添加/编辑轮播图
  const [currentBanner, setCurrentBanner] = useState<BannerItem | null>(null);
  const [bannerForm] = Form.useForm();
  
  // 添加轮播图状态
  const [bannerModalVisible, setBannerModalVisible] = useState(false);
  
  const handleAddBanner = () => {
    setCurrentBanner(null);
    bannerForm.resetFields();
    setBannerModalVisible(true);
  };
  
  const handleEditBanner = (banner: BannerItem) => {
    setCurrentBanner(banner);
    bannerForm.setFieldsValue({
      linkType: banner.linkType,
      linkUrl: banner.linkUrl,
      isEnabled: banner.isEnabled,
      order: banner.order || bannerList.findIndex(item => item.id === banner.id) + 1,
    });
    setBannerModalVisible(true);
  };
  
  // 获取轮播图列表
  const fetchBannerList = async () => {
    setLoading(true);
    try {
      const response = await getBannerList();
      if (response.code === 0) {
        setBannerList(response.data.list);
        // 从接口响应中获取轮播图区域全局状态
        setBannerEnabled(response.data.isGlobalEnabled);
      } else {
        message.error(response.message || '获取轮播图列表失败');
      }
    } catch (error) {
      console.error('获取轮播图列表失败:', error);
      message.error('获取轮播图列表失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 修改下面这部分代码替换原有函数
  // BannerModule内部调用此函数
  const handleSaveBanner = async () => {
    try {
      const values = await bannerForm.validateFields();
      setLoading(true);
      
      // 获取image URL
      let imageUrl = '';
      if (values.image?.fileList?.[0]?.url) {
        imageUrl = values.image.fileList[0].url;
      } else if (values.image?.fileList?.[0]?.response?.url) {
        imageUrl = values.image.fileList[0].response.url;
      } else if (currentBanner?.image) {
        imageUrl = currentBanner.image;
      }
      
      const formData = {
        linkType: values.linkType || 'page',
        linkUrl: values.linkUrl || '',
        image: imageUrl,
        isEnabled: values.isEnabled !== undefined ? values.isEnabled : true,
        order: values.order || (bannerList.length + 1),
      };
      
      if (currentBanner) {
        // 编辑现有Banner
        const updateData: BannerUpdateReq = {
          id: currentBanner.id,
          ...formData,
        };
        
        const response = await updateBanner(updateData);
        if (response.code === 0) {
          message.success('轮播图编辑成功');
          // 重新加载轮播图列表
          fetchBannerList();
          // 关闭弹窗
          setBannerModalVisible(false);
        } else {
          message.error(response.message || '更新轮播图失败');
        }
      } else {
        // 添加新Banner
        const createData: BannerCreateReq = formData;
        
        const response = await createBanner(createData);
        if (response.code === 0) {
          message.success('轮播图添加成功');
          // 重新加载轮播图列表
          fetchBannerList();
          // 关闭弹窗
          setBannerModalVisible(false);
        } else {
          message.error(response.message || '添加轮播图失败');
        }
      }
    } catch (error) {
      console.error('保存轮播图失败:', error);
      message.error('保存轮播图失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 删除轮播图
  const handleDeleteBanner = async (id: number) => {
    setLoading(true);
    try {
      const response = await deleteBanner(id);
      if (response.code === 0) {
        message.success('轮播图已删除');
        fetchBannerList();
      } else {
        message.error(response.message || '删除轮播图失败');
      }
    } catch (error) {
      console.error('删除轮播图失败:', error);
      message.error('删除失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 切换轮播图状态
  const toggleBannerStatus = async (id: number, isEnabled: boolean) => {
    setLoading(true);
    try {
      const data: BannerStatusUpdateReq = { id, isEnabled };
      const response = await updateBannerStatus(data);
      if (response.code === 0) {
        message.success(`${isEnabled ? '启用' : '禁用'}成功`);
        fetchBannerList();
      } else {
        message.error(response.message || '更新状态失败');
      }
    } catch (error) {
      console.error('更新轮播图状态失败:', error);
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 导航小程序列表
  const [miniPrograms, setMiniPrograms] = useState<MiniProgramItem[]>([]);
  
  // 导航区域全局状态
  const [navigationEnabled, setNavigationEnabled] = useState(true);
  
  // 活动区域全局状态
  const [activityAreaEnabled, setActivityAreaEnabled] = useState(true);
  
  // 轮播图区域全局状态
  const [bannerEnabled, setBannerEnabled] = useState(true);
  
  // 添加导航小程序
  const [miniProgramModalVisible, setMiniProgramModalVisible] = useState(false);
  const [miniProgramForm] = Form.useForm();
  const [currentMiniProgram, setCurrentMiniProgram] = useState<MiniProgramItem | null>(null);
  
  // 文件选择器状态
  const [fileSelectorVisible, setFileSelectorVisible] = useState(false);
  const [bannerSelectorVisible, setBannerSelectorVisible] = useState(false);
  
  // 处理文件选择 - 小程序图标
  const handleFileSelected = (url: string) => {
    // 创建一个模拟的文件对象
    const mockFile = {
      uid: '-1',
      name: url.split('/').pop() || '图片',
      status: 'done',
      url: url,
      thumbUrl: url,
    };
    
    // 设置miniProgramForm的logo字段值
    miniProgramForm.setFieldsValue({
      logo: { fileList: [mockFile] }
    });
    
    setFileSelectorVisible(false);
  };
  
  // 处理文件选择 - 轮播图
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
  
  const handleAddMiniProgram = () => {
    setCurrentMiniProgram(null);
    setMiniProgramModalVisible(true);
    miniProgramForm.resetFields();
  };
  
  const handleEditMiniProgram = (miniProgram: MiniProgramItem) => {
    setCurrentMiniProgram(miniProgram);
    miniProgramForm.setFieldsValue({
      name: miniProgram.name,
      appId: miniProgram.appId,
      order: miniProgram.order,
      isEnabled: miniProgram.isEnabled,
    });
    setMiniProgramModalVisible(true);
  };
  
  const handleSaveMiniProgram = async () => {
    try {
      const values = await miniProgramForm.validateFields();
      setLoading(true);
      
      // 获取logo URL
      let logoUrl = '';
      if (values.logo?.fileList?.[0]?.url) {
        logoUrl = values.logo.fileList[0].url;
      } else if (values.logo?.fileList?.[0]?.response?.url) {
        logoUrl = values.logo.fileList[0].response.url;
      } else if (currentMiniProgram?.logo) {
        logoUrl = currentMiniProgram.logo;
      }
      
      const formData = {
        name: values.name,
        appId: values.appId,
        logo: logoUrl,
        isEnabled: values.isEnabled !== undefined ? values.isEnabled : true,
        order: values.order || (miniPrograms.length + 1),
      };
      
      if (currentMiniProgram) {
        // 编辑现有小程序
        const updateData = {
          id: currentMiniProgram.id,
          ...formData,
        };
        
        const response = await updateMiniProgram(updateData);
        if (response.code === 0) {
          message.success('导航小程序编辑成功');
          // 重新加载小程序列表
          fetchMiniProgramList();
        } else {
          message.error(response.message || '更新导航小程序失败');
        }
      } else {
        // 添加新小程序
        const response = await createMiniProgram(formData);
        if (response.code === 0) {
          message.success('导航小程序添加成功');
          // 重新加载小程序列表
          fetchMiniProgramList();
        } else {
          message.error(response.message || '添加导航小程序失败');
        }
      }
      
      setMiniProgramModalVisible(false);
    } catch (error) {
      console.error('保存导航小程序失败:', error);
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 删除导航小程序
  const handleDeleteMiniProgram = async (id: number) => {
    try {
      setLoading(true);
      const response = await deleteMiniProgram(id);
      if (response.code === 0) {
        message.success('导航小程序已删除');
        fetchMiniProgramList();
      } else {
        message.error(response.message || '删除导航小程序失败');
      }
    } catch (error) {
      console.error('删除导航小程序失败:', error);
      message.error('删除失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 切换导航小程序状态
  const toggleMiniProgramStatus = async (id: number, isEnabled: boolean) => {
    try {
      // 添加调试信息
      console.log('更新导航小程序状态:', { id, isEnabled });
      
      setLoading(true);
      const response = await updateMiniProgramStatus({ id, isEnabled });
      if (response.code === 0) {
        message.success(`${isEnabled ? '启用' : '禁用'}成功`);
        fetchMiniProgramList();
      } else {
        message.error(response.message || '更新状态失败');
      }
    } catch (error) {
      console.error('更新导航小程序状态失败:', error);
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };

  // 切换导航区域全局状态
  const toggleNavigationGlobalStatus = async (isEnabled: boolean) => {
    try {
      setLoading(true);
      const data: MiniProgramGlobalStatusUpdateReq = { isEnabled };
      const response = await updateMiniProgramGlobalStatus(data);
      
      if (response.code === 0) {
        setNavigationEnabled(isEnabled);
        message.success(`导航区域${isEnabled ? '已启用' : '已禁用'}`);
        
        // 重新加载所有数据
        fetchMiniProgramList();
        fetchActivityArea();
        fetchBannerList();
      } else {
        message.error(response.message || '更新导航区域状态失败');
      }
    } catch (error) {
      console.error('更新导航区域全局状态失败:', error);
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };

  // 切换活动区域全局状态
  const toggleActivityAreaGlobalStatus = async (isEnabled: boolean) => {
    try {
      setLoading(true);
      const data: ActivityAreaGlobalStatusUpdateReq = { isEnabled };
      const response = await updateActivityAreaGlobalStatus(data);
      
      if (response.code === 0) {
        setActivityAreaEnabled(isEnabled);
        message.success(`活动区域${isEnabled ? '已启用' : '已禁用'}`);
        
        // 重新加载所有数据
        fetchMiniProgramList();
        fetchActivityArea();
        fetchBannerList();
      } else {
        message.error(response.message || '更新活动区域状态失败');
      }
    } catch (error) {
      console.error('更新活动区域全局状态失败:', error);
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };
  
  // 切换轮播图区域全局状态
  const toggleBannerGlobalStatus = async (isEnabled: boolean) => {
    try {
      setLoading(true);
      const data: BannerGlobalStatusUpdateReq = { isEnabled };
      const response = await updateBannerGlobalStatus(data);
      
      if (response.code === 0) {
        setBannerEnabled(isEnabled);
        message.success(`轮播图区域${isEnabled ? '已启用' : '已禁用'}`);
        
        // 重新加载所有数据
        fetchMiniProgramList();
        fetchActivityArea();
        fetchBannerList();
      } else {
        message.error(response.message || '更新轮播图区域状态失败');
      }
    } catch (error) {
      console.error('更新轮播图区域全局状态失败:', error);
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };

  // 初始化加载数据
  useEffect(() => {
    fetchMiniProgramList();
    fetchActivityArea();
    fetchBannerList();
  }, []);

  // 获取导航小程序列表
  const fetchMiniProgramList = async () => {
    try {
      const response = await getMiniProgramList();
      if (response.code === 0) {
        setMiniPrograms(response.data.list);
        // 从接口响应中获取导航区域全局状态
        setNavigationEnabled(response.data.isGlobalEnabled);
      } else {
        message.error(response.message || '获取导航小程序列表失败');
      }
    } catch (error) {
      console.error('获取导航小程序列表失败:', error);
      message.error('获取导航小程序列表失败');
    }
  };

  // 获取活动区域数据
  const fetchActivityArea = async () => {
    setLoading(true);
    try {
      const response = await getActivityArea();
      if (response.code === 0) {
        // 将API返回的数据格式转换为表单字段格式
        form.setFieldsValue({
          // 左上模块
          leftTopTitle: response.data.topLeft.title,
          leftTopDescription: response.data.topLeft.description,
          leftTopLinkType: response.data.topLeft.linkType,
          leftTopLinkUrl: response.data.topLeft.linkUrl,
          
          // 左下模块
          leftBottomTitle: response.data.bottomLeft.title,
          leftBottomDescription: response.data.bottomLeft.description,
          leftBottomLinkType: response.data.bottomLeft.linkType,
          leftBottomLinkUrl: response.data.bottomLeft.linkUrl,
          
          // 右侧模块
          rightTitle: response.data.right.title,
          rightDescription: response.data.right.description,
          rightLinkType: response.data.right.linkType,
          rightLinkUrl: response.data.right.linkUrl,
        });
        
        // 获取活动区域全局状态
        setActivityAreaEnabled(response.data.isGlobalEnabled);
      } else {
        message.error(response.message || '获取活动区域数据失败');
      }
    } catch (error) {
      console.error('获取活动区域数据失败:', error);
      message.error('获取活动区域数据失败');
    } finally {
      setLoading(false);
    }
  };

  // 保存配置
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      
      // 使用values进行基本验证
      console.log('表单数据已验证:', values);
      
      // 构建活动区域请求数据
      const activityAreaData = {
        topLeft: {
          title: values.leftTopTitle || '',
          description: values.leftTopDescription || '',
          linkType: values.leftTopLinkType || 'page',
          linkUrl: values.leftTopLinkUrl || '',
        },
        bottomLeft: {
          title: values.leftBottomTitle || '',
          description: values.leftBottomDescription || '',
          linkType: values.leftBottomLinkType || 'page',
          linkUrl: values.leftBottomLinkUrl || '',
        },
        right: {
          title: values.rightTitle || '',
          description: values.rightDescription || '',
          linkType: values.rightLinkType || 'page',
          linkUrl: values.rightLinkUrl || '',
        }
      };
      
      // 发送保存请求
      const response = await saveActivityArea(activityAreaData);
      if (response.code === 0) {
        // 保存成功提示
        message.success('活动区域配置已保存');
      } else {
        message.error(response.message || '保存活动区域失败');
      }
    } catch (error) {
      console.error('保存基础配置失败:', error);
      message.error('保存失败，请检查表单填写是否正确');
    } finally {
      setLoading(false);
    }
  };

  // 重置表单
  const handleReset = () => {
    setLoading(true);
    try {
      // 重置表单
      form.resetFields();
      
      // 重新加载数据
      fetchMiniProgramList();
      fetchActivityArea();
      fetchBannerList();
      
      // 更新提示信息
      message.success('页面已刷新');
    } catch (error) {
      console.error('刷新页面失败:', error);
      message.error('刷新失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };
  
  // 标签页配置
  const items: TabsProps['items'] = [
    {
      key: 'function',
      label: '导航区域',
      children: (
        <NavigationModule
          miniPrograms={miniPrograms}
          handleAddMiniProgram={handleAddMiniProgram}
          handleEditMiniProgram={handleEditMiniProgram}
          handleDeleteMiniProgram={handleDeleteMiniProgram}
          toggleMiniProgramStatus={toggleMiniProgramStatus}
          navigationEnabled={navigationEnabled}
          toggleNavigationGlobalStatus={toggleNavigationGlobalStatus}
        />
      ),
    },
    {
      key: 'activity',
      label: '活动区域',
      children: (
        <ActivityModule 
          onSave={handleSubmit}
          activityAreaEnabled={activityAreaEnabled}
          toggleActivityAreaGlobalStatus={toggleActivityAreaGlobalStatus}
        />
      ),
    },
    {
      key: 'banner',
      label: '轮播区域',
      children: (
        <BannerModule
          bannerList={bannerList}
          handleAddBanner={handleAddBanner}
          handleEditBanner={handleEditBanner}
          handleDeleteBanner={handleDeleteBanner}
          toggleBannerStatus={toggleBannerStatus}
          bannerEnabled={bannerEnabled}
          toggleBannerGlobalStatus={toggleBannerGlobalStatus}
        />
      ),
    },
  ];

  return (
    <>
      <div className="mini-program-container">
        <Card className="main-dashboard-card">
          <div className="category-header">
            <Title level={4}>首页布局</Title>
            <Space>
              <Button 
                icon={<ReloadOutlined />}
                onClick={handleReset}
              >
                刷新
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
      
      {/* 添加小程序弹窗 */}
      <Modal
        title={currentMiniProgram ? "编辑导航" : "添加导航"}
        open={miniProgramModalVisible}
        onOk={handleSaveMiniProgram}
        onCancel={() => {
          setMiniProgramModalVisible(false);
          setFileSelectorVisible(false);
        }}
        okText="保存"
        cancelText="取消"
      >
        <Form form={miniProgramForm} layout="vertical">
          <Form.Item
            name="name"
            label="小程序名称"
            rules={[{ required: true, message: '请输入小程序名称' }]}
          >
            <Input placeholder="请输入小程序名称" />
          </Form.Item>
          
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
            name="appId"
            label="小程序AppID"
            rules={[{ required: true, message: '请输入小程序AppID' }]}
          >
            <Input placeholder="请输入小程序AppID" />
          </Form.Item>
          
          <Form.Item
            name="logo"
            label="小程序图标"
          >
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
              {/* 已上传图片 */}
              {(miniProgramForm.getFieldValue('logo')?.fileList?.[0]?.url || currentMiniProgram?.logo) && (
                <div style={{ width: 104, height: 104, border: '1px dashed #d9d9d9', padding: 8, boxSizing: 'border-box' }}>
                  <img 
                    src={miniProgramForm.getFieldValue('logo')?.fileList?.[0]?.url || currentMiniProgram?.logo} 
                    alt="小程序图标" 
                    style={{ width: '100%', height: '100%', objectFit: 'contain' }} 
                  />
                </div>
              )}
              
              {/* 上传按钮 - 简洁样式 */}
              <div 
                style={{ 
                  width: 104, 
                  height: 104, 
                  border: '1px dashed #d9d9d9', 
                  display: 'flex', 
                  justifyContent: 'center', 
                  alignItems: 'center',
                  cursor: 'pointer',
                  flexDirection: 'column',
                  backgroundColor: '#fafafa'
                }}
                onClick={() => setFileSelectorVisible(true)}
              >
                <PlusOutlined style={{ fontSize: 20, color: '#999' }} />
                <div style={{ marginTop: 8, color: '#666' }}>上传</div>
              </div>
            </div>
          </Form.Item>
          
          <Form.Item
            name="isEnabled"
            label="是否启用"
            valuePropName="checked"
            initialValue={true}
          >
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
        </Form>
      </Modal>
      
      {/* 文件选择器 */}
      {fileSelectorVisible && (
        <FileSelector
          visible={fileSelectorVisible}
          onCancel={() => setFileSelectorVisible(false)}
          onSelect={handleFileSelected}
        />
      )}
      
      {/* 轮播图选择器 */}
      {bannerSelectorVisible && (
        <FileSelector
          visible={bannerSelectorVisible}
          onCancel={() => setBannerSelectorVisible(false)}
          onSelect={handleBannerImageSelected}
        />
      )}
      
      {/* 添加/编辑轮播图弹窗 */}
      <Modal
        title={currentBanner ? "编辑轮播图" : "添加轮播图"}
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
            label="轮播图片"
          >
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
              {/* 已上传图片 */}
              {(bannerForm.getFieldValue('image')?.fileList?.[0]?.url || currentBanner?.image) && (
                <div style={{ width: 200, height: 120, border: '1px dashed #d9d9d9', padding: 8, boxSizing: 'border-box' }}>
                  <img 
                    src={bannerForm.getFieldValue('image')?.fileList?.[0]?.url || currentBanner?.image} 
                    alt="轮播图片" 
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
            initialValue="page"
          >
            <Select placeholder="请选择跳转类型">
              <Select.Option value="page">小程序页面</Select.Option>
              <Select.Option value="miniprogram">其他小程序</Select.Option>
              <Select.Option value="webview">网页</Select.Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="linkUrl"
            label="跳转地址"
          >
            <Input placeholder="请输入跳转地址" />
          </Form.Item>
          
          <Form.Item
            name="isEnabled"
            label="是否启用"
            valuePropName="checked"
            initialValue={true}
          >
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default BasicSettings; 