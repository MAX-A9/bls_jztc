import React from 'react';
import {
  Alert,
  Button,
  Form,
  Switch,
  Table,
  Space,
  Popconfirm,
  Image,
  Avatar,
  Tag,
  Tooltip
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';

interface BannerModuleProps {
  bannerList: any[];
  handleAddBanner: () => void;
  handleEditBanner: (banner: any) => void;
  handleDeleteBanner: (id: number) => void;
  toggleBannerStatus: (id: number, isEnabled: boolean) => void;
  bannerEnabled: boolean;
  toggleBannerGlobalStatus: (isEnabled: boolean) => void;
}

const BannerModule: React.FC<BannerModuleProps> = ({
  bannerList,
  handleAddBanner,
  handleEditBanner,
  handleDeleteBanner,
  toggleBannerStatus,
  bannerEnabled,
  toggleBannerGlobalStatus
}) => {
  return (
    <>
      <Alert
        message="轮播区域设置"
        description="配置微信小程序首页的轮播图，可以添加多张Banner图，设置跳转链接。"
        type="info"
        showIcon
        style={{ marginBottom: 24 }}
      />
      
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Form.Item
          name="enableBanner"
          label="启用轮播图"
          style={{ marginBottom: 0 }}
        >
          <Switch 
            checked={bannerEnabled}
            onChange={(checked) => toggleBannerGlobalStatus(checked)}
            checkedChildren="已启用" 
            unCheckedChildren="已禁用" 
          />
        </Form.Item>
        
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={handleAddBanner}
        >
          添加轮播图
        </Button>
      </div>
      
      <Table
        dataSource={bannerList}
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
                    alt="轮播图" 
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
                onChange={(checked) => toggleBannerStatus(record.id, checked)}
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
                  onClick={() => handleEditBanner(record)}
                >
                  编辑
                </Button>
                <Popconfirm
                  title="确定要删除这个轮播图吗？"
                  onConfirm={() => handleDeleteBanner(record.id)}
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

export default BannerModule; 