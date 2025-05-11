import React from 'react';
import {
  Alert,
  Button,
  Form,
  Switch,
  Table,
  Space,
  Avatar,
  Popconfirm,
  Image
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';

interface NavigationModuleProps {
  miniPrograms: any[];
  handleAddMiniProgram: () => void;
  handleEditMiniProgram: (miniProgram: any) => void;
  handleDeleteMiniProgram: (id: number) => void;
  toggleMiniProgramStatus: (id: number, isEnabled: boolean) => void;
  navigationEnabled: boolean;
  toggleNavigationGlobalStatus: (isEnabled: boolean) => void;
}

const NavigationModule: React.FC<NavigationModuleProps> = ({
  miniPrograms,
  handleAddMiniProgram,
  handleEditMiniProgram,
  handleDeleteMiniProgram,
  toggleMiniProgramStatus,
  navigationEnabled,
  toggleNavigationGlobalStatus
}) => {
  return (
    <>
      <Alert
        message="导航区域设置"
        description="配置首页中跳转到其他小程序的导航区域，用户可通过点击导航直接跳转到对应的小程序。"
        type="info"
        showIcon
        style={{ marginBottom: 24 }}
      />
      
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <Form.Item
            name="enableNavigation"
            label="启用导航区域"
            style={{ marginBottom: 0 }}
          >
            <Switch 
              checked={navigationEnabled}
              onChange={(checked) => toggleNavigationGlobalStatus(checked)}
              checkedChildren="已启用" 
              unCheckedChildren="已禁用" 
            />
          </Form.Item>
        </div>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={handleAddMiniProgram}
        >
          添加小程序
        </Button>
      </div>
      
      <Table
        dataSource={miniPrograms}
        rowKey="id"
        pagination={false}
        columns={[
          {
            title: '排序',
            dataIndex: 'order',
            key: 'order',
            width: 80,
            render: (order) => (
              <div style={{ textAlign: 'center' }}>{order}</div>
            ),
          },
          {
            title: '图标',
            dataIndex: 'logo',
            key: 'logo',
            width: 100,
            render: logo => (
              <div style={{ textAlign: 'center' }}>
                {logo ? (
                  <Image 
                    src={logo} 
                    alt="小程序图标" 
                    width={40} 
                    height={40}
                    style={{ objectFit: 'cover' }}
                  />
                ) : (
                  <Avatar size={40} shape="square" style={{ backgroundColor: '#87d068' }}>
                    {miniPrograms.find(item => item.logo === logo)?.name?.slice(0, 1) || '小'}
                  </Avatar>
                )}
              </div>
            ),
          },
          {
            title: '小程序名称',
            dataIndex: 'name',
            key: 'name',
          },
          {
            title: 'AppID',
            dataIndex: 'appId',
            key: 'appId',
          },
          {
            title: '状态',
            dataIndex: 'isEnabled',
            key: 'isEnabled',
            render: (isEnabled, record) => (
              <Switch
                checked={isEnabled}
                onChange={(checked) => toggleMiniProgramStatus(record.id, checked)}
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
                  onClick={() => handleEditMiniProgram(record)}
                >
                  编辑
                </Button>
                <Popconfirm
                  title="确定要删除这个导航小程序吗？"
                  onConfirm={() => handleDeleteMiniProgram(record.id)}
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
    </>
  );
};

export default NavigationModule; 