import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  InputNumber,
  Switch,
  Popconfirm,
  message,
  Typography,
  Image,
  Tooltip
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, ReloadOutlined, PictureOutlined } from '@ant-design/icons';
import './styles.css';
import { 
  getHomeCategoryList, 
  createHomeCategory, 
  updateHomeCategory, 
  deleteHomeCategory,
  type CategoryItem 
} from '@/api/category';
import FileSelector from '@/components/FileSelector';

const { Title } = Typography;

const HomeCategory: React.FC = () => {
  const [categories, setCategories] = useState<CategoryItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [modalVisible, setModalVisible] = useState<boolean>(false);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [form] = Form.useForm();
  const [editingId, setEditingId] = useState<number | null>(null);
  const [fileSelectorVisible, setFileSelectorVisible] = useState<boolean>(false);

  useEffect(() => {
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    setLoading(true);
    try {
      const res = await getHomeCategoryList();
      
      if (res.code === 0) {
        setCategories(res.data.list);
      } else {
        message.error(res.message || '获取首页分类失败');
      }
    } catch (error) {
      message.error('获取首页分类失败');
      console.error('获取首页分类失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingId(null);
    form.resetFields();
    form.setFieldsValue({
      sortOrder: 0,
      isActive: true
    });
    setModalVisible(true);
  };

  const handleEdit = (record: CategoryItem) => {
    setEditingId(record.id);
    form.setFieldsValue({
      name: record.name,
      sortOrder: record.sortOrder,
      isActive: record.isActive,
      icon: record.icon,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await deleteHomeCategory(id);
      
      if (res.code === 0) {
        message.success('删除分类成功');
        fetchCategories();
      } else {
        message.error(res.message || '删除分类失败');
      }
    } catch (error) {
      message.error('删除分类失败');
      console.error('删除分类失败:', error);
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      
      let res;
      if (editingId === null) {
        // 添加新分类
        res = await createHomeCategory({
          name: values.name,
          sortOrder: values.sortOrder,
          isActive: values.isActive,
          icon: values.icon
        });
      } else {
        // 更新已有分类
        res = await updateHomeCategory({
          id: editingId,
          name: values.name,
          sortOrder: values.sortOrder,
          isActive: values.isActive,
          icon: values.icon
        });
      }
      
      if (res.code === 0) {
        message.success(editingId === null ? '添加分类成功' : '更新分类成功');
        setModalVisible(false);
        fetchCategories();
      } else {
        message.error(res.message || '操作失败');
      }
    } catch (error) {
      console.error('提交表单失败:', error);
      message.error('提交表单失败，请检查表单内容');
    } finally {
      setSubmitting(false);
    }
  };

  const handleCancel = () => {
    setModalVisible(false);
  };

  const openFileSelector = () => {
    setFileSelectorVisible(true);
  };

  const handleSelectFile = (url: string) => {
    form.setFieldsValue({ icon: url });
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 70,
    },
    {
      title: '图标',
      dataIndex: 'icon',
      key: 'icon',
      width: 80,
      render: (icon: string) => (
        icon ? (
          <Image 
            src={icon} 
            alt="分类图标" 
            width={32} 
            height={32}
            style={{ objectFit: 'contain' }}
          />
        ) : (
          <span>无图标</span>
        )
      ),
    },
    {
      title: '分类名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '排序',
      dataIndex: 'sortOrder',
      key: 'sortOrder',
      width: 100,
      sorter: (a: CategoryItem, b: CategoryItem) => a.sortOrder - b.sortOrder,
    },
    {
      title: '状态',
      dataIndex: 'isActive',
      key: 'isActive',
      width: 100,
      render: (isActive: boolean) => (
        isActive ? <span style={{ color: 'green' }}>启用</span> : <span style={{ color: 'red' }}>禁用</span>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: CategoryItem) => (
        <Space size="middle">
          <Button 
            type="primary" 
            icon={<EditOutlined />} 
            onClick={() => handleEdit(record)}
            size="small"
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除此分类吗？"
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
          <Title level={4}>首页分类</Title>
          <Space>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleAdd}
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
          columns={columns}
          dataSource={categories}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10 }}
        />
        
        <Modal
          title={editingId === null ? '添加分类' : '编辑分类'}
          open={modalVisible}
          onOk={handleSubmit}
          onCancel={handleCancel}
          confirmLoading={submitting}
          maskClosable={false}
          okText="确定"
          cancelText="取消"
        >
          <Form
            form={form}
            layout="vertical"
          >
            <Form.Item
              name="name"
              label="分类名称"
              rules={[{ required: true, message: '请输入分类名称' }]}
            >
              <Input placeholder="请输入分类名称" maxLength={20} />
            </Form.Item>
            
            <Form.Item
              name="icon"
              label="分类图标"
              tooltip="请选择或上传分类图标"
            >
              <Input 
                placeholder="请选择图标" 
                readOnly
                addonAfter={
                  <Tooltip title="选择图标">
                    <PictureOutlined 
                      onClick={openFileSelector}
                      style={{ cursor: 'pointer' }}
                    />
                  </Tooltip>
                } 
              />
            </Form.Item>
            
            <Form.Item
              name="sortOrder"
              label="排序"
              rules={[{ required: true, message: '请输入排序号' }]}
              initialValue={0}
            >
              <InputNumber 
                min={0} 
                max={9999} 
                style={{ width: '100%' }} 
                placeholder="数字越小排序越靠前" 
                controls={true}
              />
            </Form.Item>
            
            <Form.Item
              name="isActive"
              label="状态"
              valuePropName="checked"
              initialValue={true}
            >
              <Switch checkedChildren="启用" unCheckedChildren="禁用" />
            </Form.Item>
          </Form>
        </Modal>

        {/* 文件选择器对话框 */}
        <FileSelector 
          visible={fileSelectorVisible}
          onCancel={() => setFileSelectorVisible(false)}
          onSelect={handleSelectFile}
          title="选择分类图标"
          accept="image/*"
        />
      </Card>
    </div>
  );
};

export default HomeCategory; 