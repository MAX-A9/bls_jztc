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
  Select,
} from 'antd';
import {
  SearchOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import type { UploadFile } from 'antd/es/upload/interface';
import '../ClientManagement.css'; // 引入客户管理页面的CSS
import { 
  getProductList, 
  createProduct, 
  updateProduct,
  deleteProduct,
  updateProductStatus
} from '../../api/product';
import { getShopCategoryList, type ShopCategory } from '../../api/product';
import FileSelector from '../../components/FileSelector';

const { Title, Text } = Typography;
const { Search } = Input;
const { Option } = Select;
const { confirm } = Modal;

interface ProductData {
  id: number;
  name: string;
  description: string;
  image: string;
  price: number;
  stock: number;
  category: string;
  status: 'active' | 'inactive' | 'soldout';
  timeRequired: number; // 兑换所需时长（小时）
  order: number; // 商品排序字段
  sales: number; // 新增销量字段
  tags: string; // 商品标签
}

// API响应状态转换为UI状态
const apiStatusToUiStatus = (apiStatus: number): 'active' | 'inactive' | 'soldout' => {
  switch (apiStatus) {
    case 1:
      return 'active';
    case 2:
      return 'soldout';
    case 0:
    default:
      return 'inactive';
  }
};

// UI状态转换为API状态
const uiStatusToApiStatus = (uiStatus: string): number => {
  switch (uiStatus) {
    case 'active':
      return 1;
    case 'soldout':
      return 2;
    default:
      return 0;
  }
};

const Mall: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [products, setProducts] = useState<ProductData[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [modalTitle, setModalTitle] = useState('添加商品');
  const [currentProduct, setCurrentProduct] = useState<ProductData | null>(null);
  const [imageFile, setImageFile] = useState<UploadFile[]>([]);
  const [form] = Form.useForm();
  // 保存总记录数
  const [total, setTotal] = useState(0);
  // 查询参数
  const [searchName, setSearchName] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [statusFilter, setStatusFilter] = useState<number>(-1); // 添加状态筛选条件，默认为全部(-1)
  // 文件选择器相关状态
  const [fileSelectorVisible, setFileSelectorVisible] = useState(false);
  // 商品分类相关状态
  const [categories, setCategories] = useState<ShopCategory[]>([]);
  const [categoryLoading, setCategoryLoading] = useState(false);

  useEffect(() => {
    fetchProducts();
    fetchCategories(); // 在组件初始加载时获取分类列表
  }, [searchName, currentPage, pageSize, statusFilter]);

  // 获取商品分类列表
  const fetchCategories = async () => {
    setCategoryLoading(true);
    try {
      const result = await getShopCategoryList({
        page: 1,
        size: 50,
        status: 1 // 只获取启用状态的分类
      });
      
      if (result.code === 0 && result.data) {
        const categoryList = result.data.list || [];
        if (Array.isArray(categoryList)) {
          setCategories(categoryList);
        } else {
          console.error('获取分类列表失败：返回的list不是数组', categoryList);
          setCategories([]);
        }
      } else {
        console.error('获取分类列表失败', result);
        message.error('获取分类列表失败');
        setCategories([]);
      }
    } catch (error) {
      console.error('获取分类列表失败:', error);
      message.error('获取分类列表失败');
      setCategories([]);
    } finally {
      setCategoryLoading(false);
    }
  };

  const fetchProducts = async () => {
    setLoading(true);
    try {
      // 调用API获取商品列表
      const params = {
        page: currentPage,
        size: pageSize,
        name: searchName || undefined,
        status: statusFilter as -1 | 0 | 1 | 2
      };
      
      const result = await getProductList(params);
      
      // 后端返回数据结构转换为前端使用的格式
      if (result && result.data) {
        const apiProducts = result.data.list || []; // 处理list为null的情况
        
        // 确保apiProducts是数组
        if (Array.isArray(apiProducts)) {
          const convertedProducts: ProductData[] = apiProducts.map((item: any) => ({
            id: item.id,
            name: item.name,
            description: item.description,
            image: item.thumbnail,
            price: item.price,
            stock: item.stock,
            category: item.categoryName,
            status: apiStatusToUiStatus(item.status),
            timeRequired: item.duration,
            order: item.sortOrder,
            sales: item.sales,
            tags: item.tags || '',
          }));
          
          setProducts(convertedProducts);
          setTotal(result.data.total || 0);
        } else {
          console.error('API返回的list不是数组', apiProducts);
          setProducts([]);
          setTotal(0);
        }
      } else {
        console.error('API返回数据格式错误', result);
        message.error('获取商品列表失败：数据格式错误');
        setProducts([]);
        setTotal(0);
      }
    } catch (error) {
      console.error('获取商品列表失败:', error);
      message.error('获取商品列表失败');
      setProducts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    // 搜索功能实现
    console.log('搜索:', value);
    setSearchName(value);
    setCurrentPage(1); // 搜索时重置到第一页
    setStatusFilter(-1); // 搜索时重置状态筛选条件为全部
  };

  const showAddModal = () => {
    setCurrentProduct(null);
    setModalTitle('添加商品');
    setImageFile([]);
    form.resetFields();
    // 设置默认的排序值
    const defaultOrder = products.length > 0 ? Math.max(...products.map(item => item.order)) + 1 : 1;
    form.setFieldsValue({
      status: 'inactive',
      order: defaultOrder
    });
    
    // 获取分类列表
    fetchCategories();
    
    setModalVisible(true);
  };

  const showEditModal = (product: ProductData) => {
    setCurrentProduct(product);
    setModalTitle('编辑商品');
    
    // 获取分类列表
    fetchCategories();
    
    // 从分类中找到对应的分类ID
    const categoryId = categories.find(cat => cat.name === product.category)?.id || null;
    
    form.setFieldsValue({
      ...product,
      categoryId: categoryId, // 设置分类ID
    });
    
    if (product.image) {
      setImageFile([
        {
          uid: '-1',
          name: 'image.png',
          status: 'done',
          url: product.image,
        },
      ]);
    } else {
      setImageFile([]);
    }
    
    setModalVisible(true);
  };

  const handleDelete = (id: number) => {
    confirm({
      title: '确定要删除这个商品吗?',
      icon: <ExclamationCircleOutlined />,
      content: '删除后无法恢复',
      okText: '确定',
      okType: 'danger',
      cancelText: '取消',
      onOk() {
        // 调用删除API
        deleteProduct(id)
          .then(() => {
            message.success('商品已删除');
            // 刷新商品列表
            fetchProducts();
          })
          .catch(error => {
            console.error('删除商品失败', error);
            message.error('删除商品失败');
          });
      },
    });
  };

  const handleStatusChange = (id: number, newStatus: string) => {
    // 调用API更新状态
    const status = uiStatusToApiStatus(newStatus);
    updateProductStatus(id, status)
      .then(() => {
        message.success('状态已更新');
        // 重新获取数据以反映状态更改
        fetchProducts();
      })
      .catch(error => {
        console.error('更新状态失败', error);
        message.error('更新状态失败');
      });
  };

  const handleSelectImage = (url: string) => {
    const newFile: UploadFile = {
      uid: '-1',
      name: 'image.png',
      status: 'done',
      url: url,
    };
    setImageFile([newFile]);
    form.setFieldsValue({ image: url });
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      
      // 获取图片URL
      const imageUrl = imageFile.length > 0 ? 
        (imageFile[0].url || '') : '';
      
      // 获取选中的分类名称
      const selectedCategory = categories.find(cat => cat.id === values.categoryId);
      
      // 准备API参数
      const apiParams = {
        name: values.name,
        description: values.description,
        thumbnail: imageUrl,
        images: JSON.stringify([imageUrl]),
        price: values.price || 0,
        status: uiStatusToApiStatus(values.status),
        stock: values.stock,
        duration: values.timeRequired,
        sortOrder: values.order,
        categoryId: values.categoryId || 0, // 使用选择的分类ID
        categoryName: selectedCategory?.name || '', // 添加分类名称
        tags: values.tags || '', // 添加商品标签
      };
      
      if (currentProduct) {
        // 编辑现有商品
        updateProduct({
          ...apiParams,
          id: currentProduct.id
        })
          .then(() => {
            message.success('商品已更新');
            setModalVisible(false);
            fetchProducts(); // 刷新列表
          })
          .catch(error => {
            console.error('更新商品失败', error);
            message.error('更新商品失败');
          });
      } else {
        // 添加新商品
        createProduct(apiParams)
          .then(() => {
            message.success('商品已添加');
            setModalVisible(false);
            fetchProducts(); // 刷新列表
          })
          .catch(error => {
            console.error('添加商品失败', error);
            message.error('添加商品失败');
          });
      }
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  // 表格分页变化
  const handleTableChange = (pagination: any, filters: any) => {
    setCurrentPage(pagination.current);
    setPageSize(pagination.pageSize);
    
    // 如果有筛选状态，则调用API进行筛选
    if (filters.status && filters.status.length > 0) {
      // 开始构建查询
      let statusValue;
      const statusFilterValue = filters.status[0];
      
      // 将UI状态值转换为API状态值
      if (statusFilterValue === 'all') {
        statusValue = -1;
      } else if (statusFilterValue === 'active') {
        statusValue = 1;
      } else if (statusFilterValue === 'inactive') {
        statusValue = 0;
      } else if (statusFilterValue === 'soldout') {
        statusValue = 2;
      }
      
      // 更新状态筛选条件，触发useEffect中的fetchProducts
      if (statusValue !== undefined) {
        setStatusFilter(statusValue);
      }
    } else {
      // 如果没有筛选状态，则显示全部
      setStatusFilter(-1);
    }
  };

  // 表格列定义
  const columns = [
    {
      title: '排序',
      dataIndex: 'order',
      key: 'order',
      width: 80,
      sorter: (a: ProductData, b: ProductData) => a.order - b.order,
    },
    {
      title: '商品图片',
      dataIndex: 'image',
      key: 'image',
      width: 120,
      render: (image: string, record: ProductData) => (
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          {image ? (
            <img
              src={image}
              alt={record.name}
              style={{ width: 50, height: 50, objectFit: 'cover', borderRadius: '4px' }}
            />
          ) : (
            <div style={{ width: 50, height: 50, display: 'flex', justifyContent: 'center', alignItems: 'center', background: '#f0f0f0', borderRadius: '4px' }}>
              无图片
            </div>
          )}
        </div>
      ),
    },
    {
      title: '商品名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => <Text>{text}</Text>,
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      filters: categories.map(cat => ({ text: cat.name, value: cat.name })),
      onFilter: (value: any, record: ProductData) => record.category === value,
    },
    {
      title: '标签',
      dataIndex: 'tags',
      key: 'tags',
      width: 150,
      render: (tags: string) => {
        if (!tags) return '-';
        const tagArray = tags.split(',').filter(tag => tag.trim());
        return (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
            {tagArray.map((tag, index) => (
              <span 
                key={index}
                style={{ 
                  background: '#f0f0f0', 
                  padding: '2px 6px', 
                  borderRadius: '4px',
                  fontSize: '12px' 
                }}
              >
                {tag.trim()}
              </span>
            ))}
          </div>
        );
      },
    },
    {
      title: '所需时长',
      dataIndex: 'timeRequired',
      key: 'timeRequired',
      sorter: (a: ProductData, b: ProductData) => a.timeRequired - b.timeRequired,
      render: (time: number) => `${time} 天`,
    },
    {
      title: '库存',
      dataIndex: 'stock',
      key: 'stock',
      sorter: (a: ProductData, b: ProductData) => a.stock - b.stock,
    },
    {
      title: '销量',
      dataIndex: 'sales',
      key: 'sales',
      sorter: (a: ProductData, b: ProductData) => a.sales - b.sales,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string, record: ProductData) => (
        <Select 
          value={status} 
          style={{ width: 100 }} 
          onChange={(value) => handleStatusChange(record.id, value)}
        >
          <Option value="active">已上架</Option>
          <Option value="inactive">未上架</Option>
          <Option value="soldout">已售罄</Option>
        </Select>
      ),
      filters: [
        { text: '全部', value: 'all' },
        { text: '已上架', value: 'active' },
        { text: '未上架', value: 'inactive' },
        { text: '已售罄', value: 'soldout' },
      ],
      onFilter: (value: any, record: ProductData) => {
        if (value === 'all') return true;
        return record.status === value;
      },
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_: any, record: ProductData) => (
        <Space size="small">
          <Button 
            type="primary" 
            size="small" 
            icon={<EditOutlined />}
            onClick={() => showEditModal(record)}
          >
            编辑
          </Button>
          <Button 
            danger 
            size="small" 
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="client-management-container">
      <Card className="client-card">
        <div className="client-header">
          <Title level={4}>商城管理</Title>
          <Space size="large">
            <Search
              placeholder="搜索商品名称"
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
              添加商品
            </Button>
            <Button 
              icon={<ReloadOutlined />} 
              onClick={fetchProducts}
              loading={loading}
            >
              刷新
            </Button>
          </Space>
        </div>
        
        <Table
          dataSource={products}
          columns={columns}
          rowKey="id"
          loading={loading}
          onChange={handleTableChange}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 件商品`,
          }}
        />
      </Card>
      
      <Modal
        title={modalTitle}
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
            <Col span={12}>
              <Form.Item
                name="name"
                label="商品名称"
                rules={[{ required: true, message: '请输入商品名称' }]}
              >
                <Input placeholder="请输入商品名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="categoryId"
                label="商品分类"
                rules={[{ required: true, message: '请选择商品分类' }]}
              >
                <Select 
                  placeholder="请选择商品分类" 
                  loading={categoryLoading}
                >
                  {categories.map(category => (
                    <Option key={category.id} value={category.id}>
                      {category.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>
          
          <Form.Item
            name="description"
            label="商品描述"
            rules={[{ required: true, message: '请输入商品描述' }]}
          >
            <Input.TextArea rows={4} placeholder="请输入商品描述" />
          </Form.Item>
          
          <Form.Item
            name="tags"
            label="商品标签"
            tooltip="多个标签请用逗号分隔"
          >
            <Input placeholder="请输入商品标签，多个标签用逗号分隔" />
          </Form.Item>
          
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="timeRequired"
                label="所需时长(天)"
                rules={[{ required: true, message: '请输入兑换所需时长' }]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0.5}
                  max={720} // 最多30天
                  step={0.5}
                  placeholder="请输入兑换所需时长"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="stock"
                label="库存数量"
                rules={[{ required: true, message: '请输入库存数量' }]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  placeholder="请输入库存数量"
                />
              </Form.Item>
            </Col>
          </Row>
          
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                name="status"
                label="商品状态"
                initialValue="inactive"
              >
                <Select>
                  <Option value="active">已上架</Option>
                  <Option value="inactive">未上架</Option>
                  <Option value="soldout">已售罄</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="order"
                label="排序"
                rules={[{ required: true, message: '请输入排序值' }]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={1}
                  precision={0}
                  placeholder="请输入排序值"
                />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="price"
                label="价格"
                rules={[{ required: true, message: '请输入价格' }]}
                initialValue={0}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  precision={2}
                  placeholder="请输入价格"
                />
              </Form.Item>
            </Col>
          </Row>
          
          <Form.Item
            name="image"
            label="商品图片"
            rules={[{ required: true, message: '请上传商品图片' }]}
          >
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
              {/* 已上传图片 */}
              {imageFile.length > 0 && imageFile[0].url && (
                <div style={{ width: 120, height: 120, border: '1px dashed #d9d9d9', padding: 8, boxSizing: 'border-box' }}>
                  <img 
                    src={imageFile[0].url} 
                    alt="商品图片" 
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
        </Form>
      </Modal>
      
      {/* 文件选择器组件 */}
      <FileSelector
        visible={fileSelectorVisible}
        onCancel={() => setFileSelectorVisible(false)}
        onSelect={handleSelectImage}
        title="选择商品图片"
        accept="image/*"
      />
    </div>
  );
};

export default Mall; 