import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Form, 
  Input, 
  Select, 
  Switch, 
  Button, 
  message, 
  Breadcrumb, 
  Space, 
  Typography,
  Spin,
  Divider,
  DatePicker
} from 'antd';
import { ArrowLeftOutlined, SaveOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import WangEditor from '../../components/editor/WangEditor';
import './styles.css';
import { 
  getContentDetail, 
  createContent, 
  updateContent,
  type ContentCreateReq,
  type ContentUpdateReq
} from '@/api/content';
import { getCategoriesForSelect } from '@/api/category';
import dayjs from 'dayjs';

const { Title } = Typography;
const { Option } = Select;

const ContentForm: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState<boolean>(false);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [categories, setCategories] = useState<{ label: string; value: string; type: string }[]>([]);
  const [categoryLoading, setCategoryLoading] = useState<boolean>(false);
  const isEditMode = !!id;

  useEffect(() => {
    // 先获取分类，加载完成后再进行其他操作
    const loadInitialData = async () => {
      try {
        // 先获取分类
        await fetchCategories();
        
        // 再加载内容详情（如果是编辑模式）
        if (isEditMode) {
          await fetchContentDetails();
        } else {
          // 设置新建内容的默认值
          form.setFieldsValue({
            status: '待审核',
            isRecommended: false,
          });
        }
      } catch (error) {
        console.error('初始化数据失败:', error);
      }
    };
    
    loadInitialData();
  }, [id]);

  // 获取分类列表
  const fetchCategories = async () => {
    setCategoryLoading(true);
    try {
      const options = await getCategoriesForSelect();
      console.log('获取到的分类选项:', options);
      setCategories(options);
    } catch (error) {
      console.error('获取分类列表失败:', error);
      message.error('获取分类列表失败');
    } finally {
      setCategoryLoading(false);
    }
  };

  const fetchContentDetails = async () => {
    setLoading(true);
    try {
      const res = await getContentDetail(Number(id));
      
      if (res.code === 0) {
        const { data } = res;
        
        // 查找该分类是首页分类还是闲置分类，并设置正确的分类值
        let categoryValue = data.category;
        if (categories.length > 0) {
          // 在分类选项中查找匹配的分类名
          const found = categories.find(item => 
            item.type !== 'group' && item.label === data.category
          );
          if (found) {
            categoryValue = found.value; // 使用完整的值，例如 "home:手机数码"
          }
        }
        
        // 设置表单数据，如果有topUntil则需要转换为DatePicker可用的格式
        form.setFieldsValue({
          title: data.title,
          category: categoryValue,
          author: data.author,
          content: data.content,
          status: data.status,
          isRecommended: data.isRecommended,
          topUntilDate: data.topUntil ? dayjs(data.topUntil) : undefined,
        });
      } else {
        message.error(res.message || '获取内容详情失败');
      }
    } catch (error) {
      console.error('获取内容详情失败:', error);
      message.error('获取内容详情失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      
      // 确保内容字段是字符串，防止富文本编辑器异常
      if (values.content && typeof values.content !== 'string') {
        values.content = String(values.content);
      }
      
      // 处理置顶时间
      let topUntil = null;
      if (values.isRecommended) {
        if (values.topUntilDate) {
          // 使用日期选择器的值
          topUntil = dayjs(values.topUntilDate).format('YYYY-MM-DD HH:mm:ss');
        } else {
          message.error('请选择置顶截止时间');
          setSubmitting(false);
          return;
        }
      }
      
      // 处理分类值，从 "type:name" 格式中提取实际的分类名称
      if (values.category && typeof values.category === 'string' && values.category.includes(':')) {
        values.category = values.category.split(':')[1];
      }
      
      // 删除表单中不需要提交给API的字段
      delete values.topUntilDate;
      
      // 添加处理后的topUntil字段
      values.topUntil = topUntil;
      
      // 准备提交的数据
      let res;
      
      if (isEditMode) {
        // 编辑内容
        const updateData: ContentUpdateReq = {
          id: Number(id),
          ...values
        };
        
        res = await updateContent(updateData);
      } else {
        // 添加内容
        const createData: ContentCreateReq = values;
        res = await createContent(createData);
      }
      
      if (res.code === 0) {
        message.success(isEditMode ? '内容更新成功' : '内容添加成功');
        navigate('/content/all-content');
      } else {
        message.error(res.message || '提交失败');
      }
    } catch (error) {
      console.error('提交表单失败:', error);
      message.error('提交表单失败，请检查表单内容');
    } finally {
      setSubmitting(false);
    }
  };

  const handleBack = () => {
    navigate('/content/all-content');
  };

  // 渲染Select的选项
  const renderCategoryOptions = () => {
    // 过滤出分组标题和实际选项
    const groups = categories.filter(item => item.type === 'group');
    const options = categories.filter(item => item.type !== 'group');

    // 按分组渲染选项
    return groups.map(group => (
      <Select.OptGroup key={group.value} label={group.label}>
        {options
          .filter(option => {
            // 根据分组类型筛选选项
            if (group.value === 'home-group') {
              return option.type === 'home';
            } else if (group.value === 'idle-group') {
              return option.type === 'idle';
            }
            return false;
          })
          .map(option => (
            <Option key={option.value} value={option.value}>
              {option.label}
            </Option>
          ))}
      </Select.OptGroup>
    ));
  };

  return (
    <div className="category-container">
      <Card className="main-dashboard-card">
        <div className="category-header">
          <Breadcrumb items={[
            { title: '内容管理' },
            { title: '所有内容', href: '/content/all-content' },
            { title: isEditMode ? '编辑内容' : '添加内容' }
          ]} />
          <Space>
            <Button 
              icon={<ArrowLeftOutlined />} 
              onClick={handleBack}
            >
              返回
            </Button>
            <Button 
              type="primary" 
              icon={<SaveOutlined />} 
              onClick={handleSubmit}
              loading={submitting}
            >
              保存
            </Button>
          </Space>
        </div>
        
        <Title level={4}>{isEditMode ? '编辑内容' : '添加内容'}</Title>
        
        <Spin spinning={loading}>
          <Form
            form={form}
            layout="vertical"
            requiredMark="optional"
          >
            <Form.Item
              name="title"
              label="标题"
              rules={[{ required: true, message: '请输入内容标题' }]}
            >
              <Input placeholder="请输入内容标题" maxLength={100} />
            </Form.Item>
            
            <Form.Item
              name="category"
              label="分类"
              rules={[{ required: true, message: '请选择分类' }]}
            >
              <Select 
                placeholder="请选择分类" 
                loading={categoryLoading}
                showSearch
                optionFilterProp="children"
              >
                {renderCategoryOptions()}
              </Select>
            </Form.Item>
            
            <Form.Item
              name="author"
              label="作者"
              rules={[{ required: true, message: '请输入作者' }]}
            >
              <Input placeholder="请输入作者" />
            </Form.Item>
            
            <Form.Item
              name="status"
              label="状态"
              rules={[{ required: true, message: '请选择状态' }]}
            >
              <Select placeholder="请选择状态">
                <Option value="待审核">待审核</Option>
                <Option value="已发布">已发布</Option>
                <Option value="已下架">已下架</Option>
              </Select>
            </Form.Item>
            
            <Form.Item
              name="isRecommended"
              label="是否置顶"
              valuePropName="checked"
              initialValue={false}
            >
              <Switch checkedChildren="是" unCheckedChildren="否" />
            </Form.Item>
            
            <Form.Item 
              shouldUpdate={(prevValues, currentValues) => 
                prevValues.isRecommended !== currentValues.isRecommended
              }
            >
              {({ getFieldValue }) => (
                <Form.Item
                  name="topUntilDate"
                  label="置顶截止时间"
                  rules={[
                    { 
                      required: getFieldValue('isRecommended'), 
                      message: '请选择置顶截止时间' 
                    }
                  ]}
                >
                  <DatePicker 
                    showTime
                    placeholder="选择置顶截止时间"
                    disabled={!getFieldValue('isRecommended')}
                    format="YYYY-MM-DD HH:mm:ss"
                    disabledDate={(current) => current && current < dayjs().startOf('day')}
                  />
                </Form.Item>
              )}
            </Form.Item>
            
            <Divider />
            
            <Form.Item
              name="content"
              label="内容详情"
              rules={[{ required: true, message: '请输入内容详情' }]}
            >
              <WangEditor height={600} />
            </Form.Item>
          </Form>
        </Spin>
      </Card>
    </div>
  );
};

export default ContentForm; 