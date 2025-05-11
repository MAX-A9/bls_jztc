import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Tabs,
  Button,
  Table,
  Space,
  Modal,
  Form,
  Input,
  InputNumber,
  message,
  Popconfirm,
  Select,
  Switch,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import '../ClientManagement.css';
import './styles.css';
import {
  getPackageList,
  createPackage,
  updatePackage,
  deletePackage,
  setPackageTypeEnabled,
  Package,
  PackageType,
  DurationType
} from '../../api/package';

const { Title } = Typography;
const { TextArea } = Input;

const PublishSettings: React.FC = () => {
  // 状态定义
  const [topPackages, setTopPackages] = useState<Package[]>([]);
  const [publishPackages, setPublishPackages] = useState<Package[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [currentPackage, setCurrentPackage] = useState<Package | null>(null);
  const [packageType, setPackageType] = useState<PackageType>('top');
  const [activeTab, setActiveTab] = useState<string>('top');
  const [form] = Form.useForm();
  
  // 全局启用状态
  const [topPackageEnabled, setTopPackageEnabled] = useState<boolean>(true);
  const [publishPackageEnabled, setPublishPackageEnabled] = useState<boolean>(true);
  const [switchLoading, setSwitchLoading] = useState<boolean>(false);

  // 初始化加载数据
  useEffect(() => {
    fetchPackages();
  }, []);

  // 获取套餐数据
  const fetchPackages = async () => {
    setLoading(true);
    try {
      // 获取置顶套餐
      const topRes = await getPackageList({ type: 'top' });
      if (topRes.code === 0) {
        // 对数据按sortOrder从小到大排序
        const sortedTopPackages = [...(topRes.data.list || [])].sort((a, b) => a.sortOrder - b.sortOrder);
        setTopPackages(sortedTopPackages);
        
        // 使用接口返回的全局启用状态
        if (topRes.data.isGlobalEnabled !== undefined) {
          setTopPackageEnabled(topRes.data.isGlobalEnabled);
        }
      } else {
        message.error(`获取置顶套餐失败: ${topRes.message}`);
      }
      
      // 获取发布套餐
      const publishRes = await getPackageList({ type: 'publish' });
      if (publishRes.code === 0) {
        // 对数据按sortOrder从小到大排序
        const sortedPublishPackages = [...(publishRes.data.list || [])].sort((a, b) => a.sortOrder - b.sortOrder);
        setPublishPackages(sortedPublishPackages);
        
        // 使用接口返回的全局启用状态
        if (publishRes.data.isGlobalEnabled !== undefined) {
          setPublishPackageEnabled(publishRes.data.isGlobalEnabled);
        }
      } else {
        message.error(`获取发布套餐失败: ${publishRes.message}`);
      }
    } catch (error) {
      console.error('获取套餐数据失败:', error);
      message.error('获取套餐数据失败');
    } finally {
      setLoading(false);
    }
  };

  // 处理全局开关变化
  const handlePackageTypeEnabledChange = async (type: PackageType, checked: boolean) => {
    setSwitchLoading(true);
    try {
      const res = await setPackageTypeEnabled(type, {
        isEnabled: checked
      });
      
      if (res.code === 0) {
        message.success(`${type === 'top' ? '置顶' : '展示'}套餐已${checked ? '启用' : '禁用'}`);
        if (type === 'top') {
          setTopPackageEnabled(checked);
        } else {
          setPublishPackageEnabled(checked);
        }
      } else {
        message.error(`操作失败: ${res.message}`);
        // 恢复之前的状态
        if (type === 'top') {
          setTopPackageEnabled(!checked);
        } else {
          setPublishPackageEnabled(!checked);
        }
      }
    } catch (error) {
      console.error('设置套餐状态失败:', error);
      message.error('操作失败，请稍后重试');
      // 恢复之前的状态
      if (type === 'top') {
        setTopPackageEnabled(!checked);
      } else {
        setPublishPackageEnabled(!checked);
      }
    } finally {
      setSwitchLoading(false);
    }
  };

  // 显示添加套餐模态框
  const showAddModal = (type: PackageType) => {
    form.resetFields();
    form.setFieldsValue({
      type: type,
      price: 0,
      duration: 1,
      durationType: 'day',
      sortOrder: 0
    });
    setPackageType(type);
    setCurrentPackage(null);
    setModalVisible(true);
  };

  // 显示编辑套餐模态框
  const showEditModal = (record: Package, type: PackageType) => {
    setPackageType(type);
    setCurrentPackage(record);
    form.setFieldsValue({
      id: record.id,
      title: record.title,
      description: record.description,
      price: record.price,
      type: record.type,
      duration: record.duration,
      durationType: record.durationType,
      sortOrder: record.sortOrder
    });
    setModalVisible(true);
  };

  // 处理删除套餐
  const handleDelete = async (id: number) => {
    setLoading(true);
    try {
      const res = await deletePackage(id);
      if (res.code === 0) {
        message.success('套餐已成功删除');
        fetchPackages(); // 重新加载数据
      } else {
        message.error(`删除失败: ${res.message}`);
      }
    } catch (error) {
      console.error('删除套餐失败:', error);
      message.error('删除套餐失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  // 处理模态框确认
  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      
      if (currentPackage) {
        // 更新套餐
        const updateParams = {
          id: currentPackage.id,
          title: values.title,
          description: values.description,
          price: values.price,
          type: packageType,
          duration: values.duration,
          durationType: values.durationType,
          sortOrder: values.sortOrder
        };
        
        const res = await updatePackage(updateParams);
        if (res.code === 0) {
          message.success('套餐已成功更新');
          setModalVisible(false);
          fetchPackages(); // 重新加载数据
        } else {
          message.error(`更新失败: ${res.message}`);
        }
      } else {
        // 创建套餐
        const createParams = {
          title: values.title,
          description: values.description,
          price: values.price,
          type: packageType,
          duration: values.duration,
          durationType: values.durationType,
          sortOrder: values.sortOrder
        };
        
        const res = await createPackage(createParams);
        if (res.code === 0) {
          message.success('套餐已成功创建');
          setModalVisible(false);
          fetchPackages(); // 重新加载数据
        } else {
          message.error(`创建失败: ${res.message}`);
        }
      }
    } catch (error) {
      console.error('提交表单失败:', error);
      message.error('操作失败，请检查表单填写是否正确');
    } finally {
      setLoading(false);
    }
  };

  // 表格列定义
  const columns = [
    {
      title: '排序',
      dataIndex: 'sortOrder',
      key: 'sortOrder',
    },
    {
      title: '套餐名称',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: '套餐简介',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '时长',
      dataIndex: 'duration',
      key: 'duration',
      render: (duration: number, record: Package) => {
        const unitText = record.durationType === 'hour' ? '小时' : 
                         record.durationType === 'day' ? '天' : '月';
        return `${duration}${unitText}`;
      },
    },
    {
      title: '价格（元）',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => `¥${price.toFixed(2)}`,
    },
    {
      title: '操作',
      key: 'action',
      width: 160,
      render: () => null, // 移除未使用的参数
    },
  ];

  // 构建置顶套餐的表格列
  const topColumns = columns.map(column => {
    if (column.key === 'action') {
      return {
        ...column,
        render: (_: any, record: Package) => (
          <Space size="small">
            <Button 
              type="primary" 
              size="small" 
              icon={<EditOutlined />}
              onClick={() => showEditModal(record, 'top')}
            >
              编辑
            </Button>
            <Popconfirm
              title="确定要删除此套餐吗?"
              onConfirm={() => handleDelete(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button 
                danger 
                size="small" 
                icon={<DeleteOutlined />}
              >
                删除
              </Button>
            </Popconfirm>
          </Space>
        ),
      };
    }
    return column;
  });

  // 构建发布套餐的表格列
  const publishColumns = columns.map(column => {
    if (column.key === 'action') {
      return {
        ...column,
        render: (_: any, record: Package) => (
          <Space size="small">
            <Button 
              type="primary" 
              size="small" 
              icon={<EditOutlined />}
              onClick={() => showEditModal(record, 'publish')}
            >
              编辑
            </Button>
            <Popconfirm
              title="确定要删除此套餐吗?"
              onConfirm={() => handleDelete(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button 
                danger 
                size="small" 
                icon={<DeleteOutlined />}
              >
                删除
              </Button>
            </Popconfirm>
          </Space>
        ),
      };
    }
    return column;
  });

  // 定义Tab页items
  const tabItems = [
    {
      key: 'top',
      label: (
        <Space>
          <span>置顶套餐</span>
          <Switch
            checkedChildren="开启"
            unCheckedChildren="禁用"
            checked={topPackageEnabled}
            loading={activeTab === 'top' && switchLoading}
            onChange={(checked) => handlePackageTypeEnabledChange('top', checked)}
          />
        </Space>
      ),
      children: (
        <Table
          dataSource={topPackages}
          columns={topColumns}
          rowKey="id"
          loading={loading}
          pagination={false}
        />
      )
    },
    {
      key: 'publish',
      label: (
        <Space>
          <span>展示套餐</span>
          <Switch
            checkedChildren="开启"
            unCheckedChildren="禁用"
            checked={publishPackageEnabled}
            loading={activeTab === 'publish' && switchLoading}
            onChange={(checked) => handlePackageTypeEnabledChange('publish', checked)}
          />
        </Space>
      ),
      children: (
        <Table
          dataSource={publishPackages}
          columns={publishColumns}
          rowKey="id"
          loading={loading}
          pagination={false}
        />
      )
    }
  ];

  return (
    <div className="client-management-container">
      <Card variant="borderless" className="client-card">
        <div className="client-header">
          <Title level={4}>套餐设置</Title>
          <Space size="large">
            <Button 
              type="primary" 
              icon={<PlusOutlined />} 
              onClick={() => showAddModal(activeTab as PackageType)}
            >
              添加{activeTab === 'top' ? '置顶' : '展示'}套餐
            </Button>
            <Button 
              icon={<ReloadOutlined />} 
              onClick={fetchPackages}
              loading={loading}
            >
              刷新
            </Button>
          </Space>
        </div>

        <Tabs defaultActiveKey="top" onChange={(key) => setActiveTab(key)} items={tabItems} />
      </Card>

      {/* 添加/编辑套餐模态框 */}
      <Modal
        title={currentPackage ? '编辑套餐' : '添加套餐'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        okText="保存"
        cancelText="取消"
        confirmLoading={loading}
      >
        <Form
          form={form}
          layout="vertical"
        >
          <Form.Item name="id" hidden>
            <Input />
          </Form.Item>
          
          <Form.Item name="type" hidden>
            <Input />
          </Form.Item>
          
          <Form.Item
            name="title"
            label="套餐名称"
            rules={[{ required: true, message: '请输入套餐名称' }]}
          >
            <Input placeholder="请输入套餐名称" />
          </Form.Item>
          
          <Form.Item
            name="description"
            label="套餐简介"
            rules={[{ required: true, message: '请输入套餐简介' }]}
          >
            <TextArea rows={4} placeholder="请输入套餐简介" />
          </Form.Item>
          
          <Form.Item
            name="duration"
            label="时长值"
            rules={[
              { required: true, message: '请输入时长值' },
              { type: 'number', min: 1, message: '时长值必须大于0' }
            ]}
          >
            <InputNumber
              min={1}
              precision={0}
              style={{ width: '100%' }}
              placeholder="请输入时长值"
            />
          </Form.Item>

          <Form.Item
            name="durationType"
            label="时长单位"
            rules={[{ required: true, message: '请选择时长单位' }]}
          >
            <Select placeholder="请选择时长单位">
              <Select.Option value="hour">小时</Select.Option>
              <Select.Option value="day">天</Select.Option>
              <Select.Option value="month">月</Select.Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="price"
            label="价格(元)"
            rules={[
              { required: true, message: '请输入套餐价格' },
              { type: 'number', min: 0, message: '价格不能小于0' }
            ]}
          >
            <InputNumber
              min={0}
              precision={2}
              step={0.01}
              style={{ width: '100%' }}
              placeholder="请输入套餐价格"
              prefix="¥"
            />
          </Form.Item>
          
          <Form.Item
            name="sortOrder"
            label="排序"
            rules={[
              { required: true, message: '请输入排序值' },
              { type: 'number', min: 0, message: '排序值不能小于0' }
            ]}
          >
            <InputNumber
              min={0}
              precision={0}
              style={{ width: '100%' }}
              placeholder="请输入排序值，数值越小越靠前"
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default PublishSettings; 