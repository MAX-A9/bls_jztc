import React from 'react';
import {
  Alert,
  Form,
  Switch,
  Card,
  Row,
  Col,
  Input,
  Select,
  Button,
  Space
} from 'antd';
import { SaveOutlined } from '@ant-design/icons';

const { TextArea } = Input;
const { Option } = Select;

interface ActivityModuleProps {
  onSave?: () => void;
  activityAreaEnabled: boolean;
  toggleActivityAreaGlobalStatus: (isEnabled: boolean) => void;
}

const ActivityModule: React.FC<ActivityModuleProps> = ({ 
  onSave,
  activityAreaEnabled,
  toggleActivityAreaGlobalStatus
}) => {
  return (
    <>
      <Alert
        message="活动区域设置"
        description="配置微信小程序首页的活动区域，包含3个内容模块：左侧2个模块，右侧1个模块。"
        type="info"
        showIcon
        style={{ marginBottom: 24 }}
      />
      
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Form.Item
          name="enableActivityArea"
          label="启用活动区域"
          style={{ marginBottom: 0 }}
        >
          <Switch 
            checked={activityAreaEnabled}
            onChange={(checked) => toggleActivityAreaGlobalStatus(checked)}
            checkedChildren="已启用" 
            unCheckedChildren="已禁用" 
          />
        </Form.Item>
        
        <Button 
          type="primary" 
          icon={<SaveOutlined />} 
          onClick={onSave}
        >
          保存配置
        </Button>
      </div>

      {/* 左上模块 */}
      <Card 
        title="左上模块" 
        style={{ marginBottom: 16 }} 
        size="small"
      >
        <Row gutter={16}>
          <Col span={24}>
            <Form.Item
              name="leftTopTitle"
              label="标题"
              rules={[{ required: true, message: '请输入左上模块标题' }]}
            >
              <Input placeholder="请输入模块标题" />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="leftTopDescription"
              label="描述"
            >
              <TextArea 
                placeholder="请输入模块描述" 
                rows={2}
                maxLength={50}
                showCount
              />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="leftTopLinkType"
              label="跳转类型"
            >
              <Select placeholder="请选择跳转类型">
                <Option value="page">小程序页面</Option>
                <Option value="miniprogram">其他小程序</Option>
                <Option value="webview">网页</Option>
              </Select>
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="leftTopLinkUrl"
              label="跳转地址"
            >
              <Input placeholder="请输入跳转地址" />
            </Form.Item>
          </Col>
        </Row>
      </Card>
      
      {/* 左下模块 */}
      <Card 
        title="左下模块" 
        style={{ marginBottom: 16 }} 
        size="small"
      >
        <Row gutter={16}>
          <Col span={24}>
            <Form.Item
              name="leftBottomTitle"
              label="标题"
              rules={[{ required: true, message: '请输入左下模块标题' }]}
            >
              <Input placeholder="请输入模块标题" />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="leftBottomDescription"
              label="描述"
            >
              <TextArea 
                placeholder="请输入模块描述" 
                rows={2}
                maxLength={50}
                showCount
              />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="leftBottomLinkType"
              label="跳转类型"
            >
              <Select placeholder="请选择跳转类型">
                <Option value="page">小程序页面</Option>
                <Option value="miniprogram">其他小程序</Option>
                <Option value="webview">网页</Option>
              </Select>
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="leftBottomLinkUrl"
              label="跳转地址"
            >
              <Input placeholder="请输入跳转地址" />
            </Form.Item>
          </Col>
        </Row>
      </Card>
      
      {/* 右侧模块 */}
      <Card 
        title="右侧模块" 
        style={{ marginBottom: 16 }} 
        size="small"
      >
        <Row gutter={16}>
          <Col span={24}>
            <Form.Item
              name="rightTitle"
              label="标题"
              rules={[{ required: true, message: '请输入右侧模块标题' }]}
            >
              <Input placeholder="请输入模块标题" />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="rightDescription"
              label="描述"
            >
              <TextArea 
                placeholder="请输入模块描述" 
                rows={2}
                maxLength={50}
                showCount
              />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="rightLinkType"
              label="跳转类型"
            >
              <Select placeholder="请选择跳转类型">
                <Option value="page">小程序页面</Option>
                <Option value="miniprogram">其他小程序</Option>
                <Option value="webview">网页</Option>
              </Select>
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item
              name="rightLinkUrl"
              label="跳转地址"
            >
              <Input placeholder="请输入跳转地址" />
            </Form.Item>
          </Col>
        </Row>
      </Card>
    </>
  );
};

export default ActivityModule; 