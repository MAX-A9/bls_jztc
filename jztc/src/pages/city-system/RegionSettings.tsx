import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Input,
  Modal,
  Form,
  Select,
  message,
  Tag,
  Typography,
  Cascader,
} from 'antd';
import {
  SearchOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import './RegionSettings.css';
import {
  getRegionList,
  createRegion,
  updateRegion,
  deleteRegion,
  getRegionDetail,
  type RegionListItem,
  type RegionCreateReq,
  type RegionUpdateReq,
} from '@/api/region';

// 保留china-division数据用于级联选择器
import provincesRaw from 'china-division/dist/provinces.json';
import citiesRaw from 'china-division/dist/cities.json';
import areasRaw from 'china-division/dist/areas.json';

// 类型定义
interface ProvinceRaw {
  code: string;
  name: string;
}

interface CityRaw {
  code: string;
  name: string;
  provinceCode: string;
}

interface AreaRaw {
  code: string;
  name: string;
  cityCode: string;
  provinceCode: string;
}

// 将导入的数据转换为正确的类型
const provinces = provincesRaw as ProvinceRaw[];
const cities = citiesRaw as CityRaw[];
const areas = areasRaw as AreaRaw[];

const { Title } = Typography;
const { Option } = Select;

// 扩展RegionListItem接口以适应当前组件需求
interface RegionData extends RegionListItem {
  parentName?: string; // 根据location解析
}

// 级联选择器选项接口
interface CascaderOption {
  value: number | string;
  label: string;
  children?: CascaderOption[];
  isLeaf?: boolean;
  disabled?: boolean;
}

const RegionSettings: React.FC = () => {
  const [regions, setRegions] = useState<RegionData[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  
  // 搜索条件
  const [searchName, setSearchName] = useState('');
  const [searchLevel, setSearchLevel] = useState<string | undefined>(undefined);
  const [searchStatus, setSearchStatus] = useState<string | undefined>(undefined);
  
  // 模态框状态
  const [modalVisible, setModalVisible] = useState(false);
  const [modalTitle, setModalTitle] = useState('');
  const [editingRegion, setEditingRegion] = useState<RegionData | null>(null);
  const [cascaderOptions, setCascaderOptions] = useState<CascaderOption[]>([]);
  const [form] = Form.useForm();

  // 获取地区数据
  const fetchRegions = async () => {
    setLoading(true);
    try {
      const params: any = {
        page: current,
        pageSize: pageSize,
      };
      
      if (searchName) {
        params.name = searchName;
      }
      
      if (searchLevel) {
        params.level = searchLevel;
      }
      
      if (searchStatus !== undefined) {
        params.status = searchStatus;
      }

      const res = await getRegionList(params);
      if (res.code === 0) {
        // 处理后端返回的数据，添加parentName属性
        const regionList = res.data.list.map(item => {
          // 从location解析父级名称，例如"北京市/朝阳区" -> "北京市"
          const locationParts = item.location.split('/');
          const parentName = locationParts.length > 1 
            ? locationParts[locationParts.length - 2] 
            : '无';

          return {
            ...item,
            parentName,
          } as RegionData;
        });

        setRegions(regionList);
        setTotal(res.data.total);
      } else {
        message.error(res.message || '获取地区列表失败');
      }
    } catch (error) {
      console.error('获取地区列表失败:', error);
      message.error('获取地区列表失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 构建级联选择器的选项，使用china-division数据
  const buildCascaderOptions = () => {
    // 构建省级选项
    const options: CascaderOption[] = provinces.map(province => {
      // 获取该省下的所有市
      const provinceCities = cities.filter(city => city.provinceCode === province.code);
      
      return {
        value: province.name,
        label: province.name,
        children: provinceCities.length > 0 ? provinceCities.map(city => {
          // 获取该市下的所有区县
          const cityAreas = areas.filter(area => area.cityCode === city.code);
          
          return {
            value: city.name,
            label: city.name,
            children: cityAreas.length > 0 ? cityAreas.map(area => ({
              value: area.name,
              label: area.name,
              isLeaf: true,
            })) : undefined,
          };
        }) : undefined,
      };
    });
    
    setCascaderOptions(options);
  };

  // 首次加载和数据变化时构建级联选择器选项
  useEffect(() => {
    buildCascaderOptions();
  }, []);

  // 首次加载和依赖项变化时获取数据
  useEffect(() => {
    fetchRegions();
  }, [current, pageSize]);

  // 处理搜索
  const handleSearch = () => {
    setCurrent(1); // 重置到第一页
    fetchRegions();
  };

  // 重置搜索条件
  const handleReset = () => {
    setSearchName('');
    setSearchLevel(undefined);
    setSearchStatus(undefined);
    setCurrent(1);
    fetchRegions();
  };

  // 处理添加地区
  const handleAdd = () => {
    setModalTitle('添加地区');
    setEditingRegion(null);
    form.resetFields();
    // 设置默认值
    form.setFieldsValue({
      status: '0' // 默认启用
    });
    setModalVisible(true);
  };

  // 处理编辑地区
  const handleEdit = async (id: number) => {
    setLoading(true);
    try {
      const res = await getRegionDetail(id);
      if (res.code === 0) {
        setModalTitle('编辑地区');
        
        // 处理location字段，转换为级联选择器的值
        const locationPath = res.data.location ? res.data.location.split('/') : [];
        
        const regionData: RegionData = {
          ...res.data,
          parentName: locationPath.length > 1 ? locationPath[locationPath.length - 2] : '无'
        };
        
        setEditingRegion(regionData);
        
        form.setFieldsValue({
          ...res.data,
          locationCascader: locationPath.length > 0 ? locationPath : undefined,
          // 确保状态值正确设置
          status: String(res.data.status) 
        });
        
        setModalVisible(true);
      } else {
        message.error(res.message || '获取地区详情失败');
      }
    } catch (error) {
      console.error('获取地区详情失败:', error);
      message.error('获取地区详情失败');
    } finally {
      setLoading(false);
    }
  };

  // 处理删除地区
  const handleDelete = async (id: number) => {
    setLoading(true);
    try {
      const res = await deleteRegion(id);
      if (res.code === 0) {
        message.success('删除地区成功');
        fetchRegions();
      } else {
        message.error(res.message || '删除地区失败');
      }
    } catch (error) {
      console.error('删除地区失败:', error);
      message.error('删除地区失败');
    } finally {
      setLoading(false);
    }
  };

  // 删除确认对话框
  const showDeleteConfirm = (id: number, name: string) => {
    Modal.confirm({
      title: '确定要删除这个地区吗?',
      content: `地区名称: ${name}`,
      icon: <ExclamationCircleOutlined />,
      okText: '确定',
      okType: 'danger',
      cancelText: '取消',
      onOk() {
        return handleDelete(id);
      },
    });
  };

  // 处理表单提交
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();

      // 处理级联选择器的值，转换为字符串格式
      const location = values.locationCascader ? values.locationCascader.join('/') : '';
      
      if (editingRegion) {
        // 更新地区
        const updateData: RegionUpdateReq = {
          id: editingRegion.id,
          name: values.name || '',
          location: location,
          level: values.level,
          status: values.status,
        };
        
        const res = await updateRegion(updateData);
        if (res.code === 0) {
          message.success('更新地区成功');
          setModalVisible(false);
          fetchRegions();
        } else {
          message.error(res.message || '更新地区失败');
        }
      } else {
        // 创建地区
        const createData: RegionCreateReq = {
          name: values.name || '',
          location: location,
          level: values.level,
          status: values.status,
        };
        
        const res = await createRegion(createData);
        if (res.code === 0) {
          message.success('添加地区成功');
          setModalVisible(false);
          fetchRegions();
        } else {
          message.error(res.message || '添加地区失败');
        }
      }
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  // 处理表单取消
  const handleCancel = () => {
    setModalVisible(false);
  };

  // 渲染状态标签
  const renderStatusTag = (status: any) => {
    // 确保正确比较状态值，因为API返回的可能是字符串
    const statusValue = String(status);
    return statusValue === '0' ? (
      <Tag color="green">启用</Tag>
    ) : (
      <Tag color="red">禁用</Tag>
    );
  };

  // 渲染级别标签
  const renderLevelTag = (level: string) => {
    let color = 'default';
    switch (level) {
      case '省':
        color = 'blue';
        break;
      case '县':
        color = 'purple';
        break;
      case '乡':
        color = 'orange';
        break;
    }
    return <Tag color={color}>{level}</Tag>;
  };

  // 分页变化处理
  const handleTableChange = (pagination: any) => {
    setCurrent(pagination.current);
    setPageSize(pagination.pageSize);
  };

  // 表格列定义
  const columns: ColumnsType<RegionData> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '地区名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '所在地区',
      dataIndex: 'location',
      key: 'location',
    },
    {
      title: '级别',
      dataIndex: 'level',
      key: 'level',
      render: renderLevelTag,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: renderStatusTag,
      width: 100,
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
            onClick={() => handleEdit(record.id)}
          >
            编辑
          </Button>
          <Button
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => showDeleteConfirm(record.id, record.name)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="region-settings-container">
      <Card className="region-card">
        <div className="region-header">
          <Title level={4}>地区管理</Title>
          <Space>
            <Input
              placeholder="搜索地区名称"
              value={searchName}
              onChange={(e) => setSearchName(e.target.value)}
              style={{ width: 200 }}
              allowClear
            />
            <Select 
              placeholder="级别" 
              value={searchLevel} 
              onChange={value => setSearchLevel(value)}
              style={{ width: 100 }}
              allowClear
            >
              <Option value="省">省</Option>
              <Option value="县">县</Option>
              <Option value="乡">乡</Option>
            </Select>
            <Select 
              placeholder="状态" 
              value={searchStatus} 
              onChange={value => setSearchStatus(value)}
              style={{ width: 100 }}
              allowClear
            >
              <Option value="0">启用</Option>
              <Option value="1">禁用</Option>
            </Select>
            <Button 
              type="primary" 
              icon={<SearchOutlined />} 
              onClick={handleSearch}
            >
              搜索
            </Button>
            <Button 
              icon={<ReloadOutlined />} 
              onClick={handleReset}
            >
              重置
            </Button>
            <Button 
              type="primary" 
              icon={<PlusOutlined />} 
              onClick={handleAdd}
            >
              添加地区
            </Button>
          </Space>
        </div>
        
        {/* 地区列表表格 */}
        <Table
          dataSource={regions}
          columns={columns}
          rowKey="id"
          pagination={{
            current,
            pageSize,
            total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
          }}
          loading={loading}
          onChange={handleTableChange}
          className="region-table"
        />
      </Card>
      
      {/* 添加/编辑地区模态框 */}
      <Modal
        title={modalTitle}
        open={modalVisible}
        onCancel={handleCancel}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{
            status: '0' // 默认启用
          }}
        >
          <Form.Item
            name="locationCascader"
            label="所在地区"
            rules={[{ required: true, message: '请选择所在地区' }]}
          >
            <Cascader
              options={cascaderOptions}
              placeholder="请选择所在地区"
              changeOnSelect
              style={{ width: '100%' }}
            />
          </Form.Item>
          
          <Form.Item
            name="name"
            label="地区名称"
            tooltip="可选填，不填则显示为'无'"
          >
            <Input placeholder="可选填，不填则显示为'无'" />
          </Form.Item>
          
          <Form.Item
            name="level"
            label="级别"
            rules={[{ required: true, message: '请选择级别' }]}
          >
            <Select placeholder="请选择级别">
              <Option value="省">省</Option>
              <Option value="县">县</Option>
              <Option value="乡">乡</Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="status"
            label="状态"
            rules={[{ required: true, message: '请选择状态' }]}
          >
            <Select placeholder="请选择状态">
              <Option value="0">启用</Option>
              <Option value="1">禁用</Option>
            </Select>
          </Form.Item>
          
          <Form.Item>
            <div style={{ textAlign: 'right' }}>
              <Button 
                onClick={handleCancel} 
                style={{ marginRight: 8 }}
              >
                取消
              </Button>
              <Button 
                type="primary" 
                htmlType="submit"
              >
                确定
              </Button>
            </div>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default RegionSettings; 