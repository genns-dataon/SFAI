import { useState, useEffect } from 'react';
import { Card, Table, Button, Modal, Form, Input, message, Popconfirm, Space, Typography, Alert } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SettingOutlined } from '@ant-design/icons';
import { settingsAPI } from '../api/api';

const { Text } = Typography;
const { TextArea } = Input;

const Settings = () => {
  const [settings, setSettings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingKey, setEditingKey] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = async () => {
    setLoading(true);
    try {
      const response = await settingsAPI.getAll();
      setSettings(response.data);
    } catch (error) {
      message.error('Failed to fetch settings');
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    form.resetFields();
    setEditingKey(null);
    setModalVisible(true);
  };

  const handleEdit = (record) => {
    form.setFieldsValue(record);
    setEditingKey(record.key);
    setModalVisible(true);
  };

  const handleDelete = async (key) => {
    try {
      await settingsAPI.delete(key);
      message.success('Setting deleted successfully');
      fetchSettings();
    } catch (error) {
      message.error('Failed to delete setting');
    }
  };

  const handleSubmit = async (values) => {
    try {
      await settingsAPI.upsert(values);
      message.success(`Setting ${editingKey ? 'updated' : 'created'} successfully`);
      setModalVisible(false);
      form.resetFields();
      fetchSettings();
    } catch (error) {
      message.error(`Failed to ${editingKey ? 'update' : 'create'} setting`);
    }
  };

  const columns = [
    {
      title: 'Key',
      dataIndex: 'key',
      key: 'key',
      width: '25%',
    },
    {
      title: 'Value',
      dataIndex: 'value',
      key: 'value',
      width: '35%',
      render: (text) => (
        <Text style={{ 
          display: 'block',
          whiteSpace: 'pre-wrap',
          maxHeight: '100px',
          overflow: 'auto'
        }}>
          {text}
        </Text>
      ),
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
      width: '30%',
      render: (text) => <Text type="secondary">{text || 'No description'}</Text>,
    },
    {
      title: 'Actions',
      key: 'actions',
      width: '10%',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          />
          <Popconfirm
            title="Are you sure you want to delete this setting?"
            onConfirm={() => handleDelete(record.key)}
            okText="Yes"
            cancelText="No"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="p-6">
      <Card
        title={
          <Space>
            <SettingOutlined style={{ fontSize: '20px' }} />
            <span>Chatbot Settings & Configuration</span>
          </Space>
        }
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            Add Setting
          </Button>
        }
      >
        <Alert
          message="Security Warning"
          description="⚠️ This settings page is UNSECURED and allows modification of chatbot prompts. This creates security risks including prompt injection attacks and PII exposure. This feature should be secured with proper authentication or removed in production."
          type="warning"
          showIcon
          style={{ marginBottom: 16 }}
        />
        
        <Alert
          message="How it works"
          description="Settings are appended to the chatbot's system prompt as additional configuration and helpers. Use these to customize chatbot behavior, add context, or provide instructions."
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
        />

        <Table
          columns={columns}
          dataSource={settings}
          loading={loading}
          rowKey="key"
          pagination={{ pageSize: 10 }}
        />
      </Card>

      <Modal
        title={editingKey ? 'Edit Setting' : 'Add Setting'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
        }}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          style={{ marginTop: 20 }}
        >
          <Form.Item
            label="Key"
            name="key"
            rules={[
              { required: true, message: 'Please enter a key' },
              { 
                pattern: /^[a-zA-Z0-9_-]+$/, 
                message: 'Key can only contain letters, numbers, hyphens, and underscores' 
              },
            ]}
          >
            <Input 
              placeholder="e.g., company_name, working_hours, policy_url"
              disabled={editingKey !== null}
            />
          </Form.Item>

          <Form.Item
            label="Value"
            name="value"
            rules={[{ required: true, message: 'Please enter a value' }]}
          >
            <TextArea
              rows={4}
              placeholder="Enter the configuration value or instruction for the chatbot"
            />
          </Form.Item>

          <Form.Item
            label="Description (Optional)"
            name="description"
          >
            <Input placeholder="Brief description of what this setting does" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingKey ? 'Update' : 'Create'}
              </Button>
              <Button onClick={() => {
                setModalVisible(false);
                form.resetFields();
              }}>
                Cancel
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Settings;
