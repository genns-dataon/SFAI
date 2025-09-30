import { useEffect, useState } from 'react';
import { Card, Table, Button, Modal, Form, Input, Select, DatePicker, Space, Tag, Typography, message } from 'antd';
import { PlusOutlined, SearchOutlined, UserOutlined, EditOutlined } from '@ant-design/icons';
import { employeeAPI } from '../api/api';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

const Employees = () => {
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editingEmployee, setEditingEmployee] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchEmployees();
  }, []);

  const fetchEmployees = async () => {
    try {
      setLoading(true);
      const response = await employeeAPI.getAll();
      setEmployees(response.data);
    } catch (error) {
      console.error('Error fetching employees:', error);
      message.error('Failed to fetch employees');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (values) => {
    try {
      const formattedData = {
        ...values,
        hire_date: values.hire_date.format('YYYY-MM-DD'),
        manager_id: values.manager_id || null,
      };
      
      if (editingEmployee) {
        await employeeAPI.update(editingEmployee.id, formattedData);
        message.success('Employee updated successfully');
      } else {
        await employeeAPI.create(formattedData);
        message.success('Employee added successfully');
      }
      
      setShowModal(false);
      setEditingEmployee(null);
      form.resetFields();
      fetchEmployees();
    } catch (error) {
      console.error('Error saving employee:', error);
      message.error(editingEmployee ? 'Failed to update employee' : 'Failed to add employee');
    }
  };

  const handleEdit = (employee) => {
    setEditingEmployee(employee);
    form.setFieldsValue({
      name: employee.name,
      email: employee.email,
      job_title: employee.job_title,
      department_id: employee.department_id,
      manager_id: employee.manager_id,
      hire_date: dayjs(employee.hire_date),
    });
    setShowModal(true);
  };

  const filteredEmployees = employees.filter(
    (emp) =>
      emp.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      emp.email.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (text) => (
        <Space>
          <div style={{
            width: 40,
            height: 40,
            backgroundColor: '#e6f7ff',
            borderRadius: '50%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <Text strong style={{ color: '#1890ff' }}>
              {text.split(' ').map(n => n[0]).join('')}
            </Text>
          </div>
          <Text strong>{text}</Text>
        </Space>
      ),
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: 'Job Title',
      dataIndex: 'job_title',
      key: 'job_title',
      render: (text) => <Tag color="blue">{text}</Tag>,
    },
    {
      title: 'Department',
      dataIndex: 'department',
      key: 'department',
      render: (dept) => <Text strong>{dept ? dept.name : 'N/A'}</Text>,
    },
    {
      title: 'Manager',
      dataIndex: 'manager',
      key: 'manager',
      render: (manager) => manager ? <Text>{manager.name}</Text> : <Text type="secondary">None</Text>,
    },
    {
      title: 'Hire Date',
      dataIndex: 'hire_date',
      key: 'hire_date',
      render: (date) => date ? new Date(date).toLocaleDateString() : 'N/A',
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Button
          type="link"
          icon={<EditOutlined />}
          onClick={() => handleEdit(record)}
        >
          Edit
        </Button>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <Title level={2} style={{ margin: 0 }}>Employees</Title>
            <Text type="secondary">Manage your workforce and team members</Text>
          </div>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setShowModal(true)}
            size="large"
          >
            Add Employee
          </Button>
        </div>

        <Card>
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            <Input
              placeholder="Search employees by name or email..."
              prefix={<SearchOutlined />}
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              size="large"
            />

            <Table
              columns={columns}
              dataSource={filteredEmployees}
              loading={loading}
              rowKey="id"
              pagination={{
                pageSize: 10,
                showSizeChanger: true,
                showTotal: (total) => `Total ${total} employees`,
              }}
            />
          </Space>
        </Card>
      </Space>

      <Modal
        title={
          <Space>
            <UserOutlined />
            <span>{editingEmployee ? 'Edit Employee' : 'Add New Employee'}</span>
          </Space>
        }
        open={showModal}
        onCancel={() => {
          setShowModal(false);
          setEditingEmployee(null);
          form.resetFields();
        }}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          style={{ marginTop: 24 }}
        >
          <Form.Item
            label="Name"
            name="name"
            rules={[{ required: true, message: 'Please enter employee name' }]}
          >
            <Input placeholder="John Doe" />
          </Form.Item>

          <Form.Item
            label="Email"
            name="email"
            rules={[
              { required: true, message: 'Please enter email' },
              { type: 'email', message: 'Please enter a valid email' },
            ]}
          >
            <Input placeholder="john.doe@company.com" />
          </Form.Item>

          <Form.Item
            label="Job Title"
            name="job_title"
            rules={[{ required: true, message: 'Please enter job title' }]}
          >
            <Input placeholder="Software Engineer" />
          </Form.Item>

          <Form.Item
            label="Department"
            name="department_id"
            rules={[{ required: true, message: 'Please select department' }]}
            initialValue={1}
          >
            <Select>
              <Select.Option value={1}>Engineering</Select.Option>
              <Select.Option value={2}>Human Resources</Select.Option>
              <Select.Option value={3}>Sales</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="Manager"
            name="manager_id"
          >
            <Select placeholder="Select a manager (optional)" allowClear>
              {employees
                .filter(emp => !editingEmployee || emp.id !== editingEmployee.id)
                .map(emp => (
                  <Select.Option key={emp.id} value={emp.id}>
                    {emp.name} ({emp.job_title})
                  </Select.Option>
                ))}
            </Select>
          </Form.Item>

          <Form.Item
            label="Hire Date"
            name="hire_date"
            rules={[{ required: true, message: 'Please select hire date' }]}
          >
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0, marginTop: 24 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setShowModal(false);
                setEditingEmployee(null);
                form.resetFields();
              }}>
                Cancel
              </Button>
              <Button type="primary" htmlType="submit">
                {editingEmployee ? 'Update Employee' : 'Add Employee'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Employees;
