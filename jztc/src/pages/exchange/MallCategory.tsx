import React, { useState, useEffect } from 'react';
import {
  Card,
  Typography,
  Table,
  Space,
  Button,
  Input,
  Row,
  Col,
  message,
  Modal,
  Form,
  InputNumber,
  Switch,
  Popconfirm,
  Image,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  ReloadOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import '../ClientManagement.css'; // 引入客户管理页面的CSS
import { 
  getShopCategoryList, 
  createShopCategory, 
  updateShopCategory, 
  deleteShopCategory, 
  updateShopCategoryStatus,
  type ShopCategory,
  type ShopCategoryListParams
} from '../../api/product';
import FileSelector from '../../components/FileSelector';

const { Title } = Typography;
const { Search } = Input;
const { confirm } = Modal;

interface CategoryData {
  id: number;
  name: string;
  isEnabled: boolean;
  order: number;
  productCount: number;
  image: string;
}

const MallCategory: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState<CategoryData[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [currentCategory, setCurrentCategory] = useState<CategoryData | null>(null);
  const [form] = Form.useForm();
  const [total, setTotal] = useState(0);
  const [searchName, setSearchName] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [fileSelectorVisible, setFileSelectorVisible] = useState(false);
  const [imageUrl, setImageUrl] = useState<string>('');

  useEffect(() => {
    fetchCategories();
  }, [searchName, currentPage, pageSize]);

  // 监听表单图片字段变化
  useEffect(() => {
    const unsubscribe = form.getFieldInstance('image')?.addEventListener('change', (e: any) => {
      if (e.target.value !== undefined) {
        setImageUrl(e.target.value);
      }
    });
    
    return () => {
      if (unsubscribe) {
        unsubscribe();
      }
    };
  }, [form]);

  // 将API返回的数据转换为组件使用的格式
  const convertToUIFormat = (apiData: ShopCategory[]): CategoryData[] => {
    return apiData.map(item => ({
      id: item.id,
      name: item.name,
      isEnabled: item.status === 1,
      order: item.sortOrder,
      productCount: item.productCount,
      image: item.image || '',
    }));
  };

  // 将UI数据转换为API需要的格式用于创建/更新
  const convertToAPIFormat = (uiData: Partial<CategoryData>, id?: number): any => {
    console.log('UI数据转换前:', uiData); // 添加日志，查看转换前的数据
    
    const data: any = {
      name: uiData.name || '',
      sortOrder: uiData.order || 0,
      status: uiData.isEnabled ? 1 : 0,
      image: uiData.image || imageUrl || '', // 优先使用表单中的image，如果没有则使用imageUrl
      productCount: uiData.productCount || 0,
    };
    
    if (id) {
      data.id = id;
    }
    
    console.log('UI数据转换后:', data); // 添加日志，查看转换后的数据
    return data;
  };

  const fetchCategories = async () => {
    setLoading(true);
    try {
      const params: ShopCategoryListParams = {
        page: currentPage,
        size: pageSize,
      };
      
      if (searchName) {
        params.name = searchName;
      }
      
      const result = await getShopCategoryList(params);
      
      if (result.code === 0 && result.data) {
        const listData = result.data.list || [];
        if (Array.isArray(listData)) {
          const convertedCategories = convertToUIFormat(listData);
          setCategories(convertedCategories);
          setTotal(result.data.total || 0);
        } else {
          console.error('获取分类列表失败：返回的list不是数组', result.data);
          setCategories([]);
          setTotal(0);
        }
      } else {
        message.error(result.message || '获取分类列表失败');
        setCategories([]);
        setTotal(0);
      }
    } catch (error) {
      console.error('获取分类列表失败:', error);
      message.error('获取分类列表失败');
      setCategories([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    setSearchName(value);
    setCurrentPage(1);
  };

  const showAddModal = () => {
    setCurrentCategory(null);
    form.resetFields();
    form.setFieldsValue({
      isEnabled: true,
      order: categories.length > 0 ? Math.max(...categories.map(item => item.order)) + 1 : 1,
      productCount: 0,
      image: ''
    });
    setImageUrl('');
    setModalVisible(true);
  };

  const showEditModal = (category: CategoryData) => {
    setCurrentCategory(category);
    form.setFieldsValue({
      name: category.name,
      isEnabled: category.isEnabled,
      order: category.order,
      productCount: category.productCount,
      image: category.image
    });
    setImageUrl(category.image);
    setModalVisible(true);
  };

  const handleDelete = (id: number) => {
    confirm({
      title: '确定要删除这个分类吗?',
      icon: <ExclamationCircleOutlined />,
      content: '删除后无法恢复，该分类下的商品将无法显示',
      okText: '确定',
      okType: 'danger',
      cancelText: '取消',
      onOk() {
        // 调用API删除分类
        deleteShopCategory(id)
          .then((result) => {
            if (result.code === 0) {
              message.success('分类已删除');
              fetchCategories(); // 刷新列表
            } else {
              message.error(result.message || '删除分类失败');
            }
          })
          .catch(error => {
            console.error('删除分类失败:', error);
            message.error('删除分类失败');
          });
      },
    });
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      console.log('表单值:', values); // 添加日志，查看表单值
      
      // 确保image字段存在并且使用正确的值
      if (!values.image && imageUrl) {
        values.image = imageUrl;
      }
      
      if (currentCategory) {
        // 编辑现有分类
        const apiData = convertToAPIFormat(values, currentCategory.id);
        console.log('发送到API的数据:', apiData); // 添加日志，查看API数据
        
        updateShopCategory(apiData)
          .then((result) => {
            if (result.code === 0) {
              message.success('分类已更新');
              setModalVisible(false);
              fetchCategories(); // 刷新列表
            } else {
              message.error(result.message || '更新分类失败');
            }
          })
          .catch(error => {
            console.error('更新分类失败:', error);
            message.error('更新分类失败');
          });
      } else {
        // 添加新分类
        const apiData = convertToAPIFormat(values);
        console.log('发送到API的数据:', apiData); // 添加日志，查看API数据
        
        // 确保image字段存在
        if (!apiData.image && imageUrl) {
          apiData.image = imageUrl;
        }
        
        createShopCategory(apiData)
          .then((result) => {
            if (result.code === 0) {
              message.success('分类已添加');
              setModalVisible(false);
              fetchCategories(); // 刷新列表
            } else {
              message.error(result.message || '添加分类失败');
            }
          })
          .catch(error => {
            console.error('添加分类失败:', error);
            message.error('添加分类失败');
          });
      }
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  const toggleCategoryStatus = (id: number, isEnabled: boolean) => {
    const status = isEnabled ? 1 : 0;
    
    updateShopCategoryStatus({ id, status })
      .then((result) => {
        if (result.code === 0) {
          message.success(`分类已${isEnabled ? '启用' : '禁用'}`);
          fetchCategories(); // 刷新列表
        } else {
          message.error(result.message || '更新分类状态失败');
        }
      })
      .catch(error => {
        console.error('更新分类状态失败:', error);
        message.error('更新分类状态失败');
      });
  };

  // 文件选择回调
  const handleFileSelected = (url: string) => {
    form.setFieldsValue({ 
      image: url
    });
    setImageUrl(url);
    setFileSelectorVisible(false);
  };

  const columns = [
    {
      title: '排序',
      dataIndex: 'order',
      key: 'order',
      width: 80,
      sorter: (a: CategoryData, b: CategoryData) => a.order - b.order,
    },
    {
      title: '分类图片',
      dataIndex: 'image',
      key: 'image',
      width: 100,
      render: (image: string) => (
        image ? <Image src={image} alt="分类图片" width={50} height={50} style={{ objectFit: 'cover' }} /> : '无图片'
      ),
    },
    {
      title: '分类名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '商品数量',
      dataIndex: 'productCount',
      key: 'productCount',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'isEnabled',
      key: 'isEnabled',
      width: 100,
      render: (isEnabled: boolean, record: CategoryData) => (
        <Switch
          checked={isEnabled}
          onChange={(checked) => toggleCategoryStatus(record.id, checked)}
          checkedChildren="启用"
          unCheckedChildren="禁用"
        />
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_: any, record: CategoryData) => (
        <Space size="small">
          <Button 
            type="primary" 
            size="small" 
            icon={<EditOutlined />} 
            onClick={() => showEditModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除吗?"
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
    },
  ];

  return (
    <div className="client-management-container">
      <Card className="client-card">
        <div className="client-header">
          <Title level={4}>商城分类</Title>
          <Space size="large">
            <Search
              placeholder="搜索分类名称"
              onSearch={handleSearch}
              style={{ width: 250 }}
              allowClear
              enterButton={<><SearchOutlined />搜索</>}
            />
            <Button 
              type="primary" 
              icon={<PlusOutlined />} 
              onClick={showAddModal}
            >
              添加分类
            </Button>
            <Button 
              icon={<ReloadOutlined />} 
              onClick={fetchCategories}
              loading={loading}
            >
              刷新
            </Button>
          </Space>
        </div>
        
        <Table
          dataSource={categories}
          columns={columns}
          rowKey="id"
          loading={loading}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 个分类`,
            onChange: (page, size) => {
              setCurrentPage(page);
              if (size) setPageSize(size);
            }
          }}
        />
      </Card>
      
      <Modal
        title={currentCategory ? "编辑分类" : "添加分类"}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        okText="保存"
        cancelText="取消"
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
        >
          <Row gutter={16}>
            <Col span={16}>
              <Form.Item
                name="name"
                label="分类名称"
                rules={[{ required: true, message: '请输入分类名称' }]}
              >
                <Input placeholder="请输入分类名称" maxLength={20} showCount />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="order"
                label="排序"
                rules={[{ required: true, message: '请输入排序值' }]}
                tooltip="数字越小排序越靠前"
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={1}
                  precision={0}
                  placeholder="请输入排序值"
                />
              </Form.Item>
            </Col>
          </Row>
          
          <Form.Item
            name="image"
            label="分类图片"
            hidden
            initialValue=""
          >
            <Input />
          </Form.Item>
          
          <Form.Item
            label="分类图片"
            required={false}
            help="请上传或选择分类图片"
          >
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
              {/* 已上传图片 */}
              {imageUrl && (
                <div style={{ width: 120, height: 120, border: '1px dashed #d9d9d9', padding: 8, boxSizing: 'border-box' }}>
                  <img 
                    src={imageUrl} 
                    alt="分类图片" 
                    style={{ width: '100%', height: '100%', objectFit: 'contain' }} 
                  />
                </div>
              )}
              
              {/* 上传按钮 */}
              <div 
                style={{ 
                  width: 120, 
                  height: 120, 
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
            name="productCount" 
            label="商品数量" 
            initialValue={0}
            tooltip="该值通常由系统自动管理，请谨慎修改"
          >
            <InputNumber min={0} precision={0} style={{ width: '100%' }} />
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

      {/* 文件选择器组件 */}
      <FileSelector
        visible={fileSelectorVisible}
        onCancel={() => setFileSelectorVisible(false)}
        onSelect={handleFileSelected}
        title="选择分类图片"
        accept="image/*"
      />
    </div>
  );
};

export default MallCategory; 