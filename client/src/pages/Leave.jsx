import { useEffect, useState } from 'react';
import { Card, Table, Button, Modal, Form, Select, DatePicker, Typography, Space, Tag, message } from 'antd';
import { CalendarOutlined, PlusOutlined } from '@ant-design/icons';
import { leaveAPI, employeeAPI } from '../api/api';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;

const Leave = () => {
  const [leaves, setLeaves] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [leavesRes, employeesRes] = await Promise.all([
        leaveAPI.getAll(),
        employeeAPI.getAll(),
      ]);
      setLeaves(leavesRes.data);
      setEmployees(employeesRes.data);
    } catch (error) {
      console.error('Error fetching data:', error);
      message.error('Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (values) => {
    try {
      await leaveAPI.create({
        employee_id: Number(values.employee_id),
        leave_type: values.leave_type,
        start_date: values.dates[0].toISOString(),
        end_date: values.dates[1].toISOString(),
      });
      message.success('Leave request submitted successfully');
      setShowModal(false);
      form.resetFields();
      fetchData();
    } catch (error) {
      console.error('Error creating leave request:', error);
      message.error('Failed to submit leave request');
    }
  };

  const getStatusColor = (status) => {
    const colors = {
      pending: 'orange',
      approved: 'green',
      rejected: 'red',
    };
    return colors[status] || 'default';
  };

  const columns = [
    {
      title: 'Employee',
      dataIndex: 'employee',
      key: 'employee',
      render: (employee, record) => (
        <Text strong>{employee ? employee.name : `Employee ${record.employee_id}`}</Text>
      ),
    },
    {
      title: 'Leave Type',
      dataIndex: 'leave_type',
      key: 'leave_type',
      render: (type) => <Tag color="blue">{type}</Tag>,
    },
    {
      title: 'Start Date',
      dataIndex: 'start_date',
      key: 'start_date',
      render: (date) => date ? new Date(date).toLocaleDateString() : 'N/A',
    },
    {
      title: 'End Date',
      dataIndex: 'end_date',
      key: 'end_date',
      render: (date) => date ? new Date(date).toLocaleDateString() : 'N/A',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={getStatusColor(status)} style={{ textTransform: 'capitalize' }}>
          {status}
        </Tag>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <Title level={2} style={{ margin: 0 }}>Leave Requests</Title>
            <Text type="secondary">Manage employee leave applications and approvals</Text>
          </div>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setShowModal(true)}
            size="large"
          >
            Request Leave
          </Button>
        </div>

        <Card>
          <Table
            columns={columns}
            dataSource={leaves}
            loading={loading}
            rowKey="id"
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showTotal: (total) => `Total ${total} requests`,
            }}
          />
        </Card>
      </Space>

      <Modal
        title={
          <Space>
            <CalendarOutlined />
            <span>Request Leave</span>
          </Space>
        }
        open={showModal}
        onCancel={() => {
          setShowModal(false);
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
            label="Employee"
            name="employee_id"
            rules={[{ required: true, message: 'Please select an employee' }]}
          >
            <Select placeholder="Select Employee" size="large">
              {employees.map((emp) => (
                <Select.Option key={emp.id} value={emp.id}>
                  {emp.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label="Leave Type"
            name="leave_type"
            rules={[{ required: true, message: 'Please select leave type' }]}
            initialValue="Vacation"
          >
            <Select size="large">
              <Select.Option value="Vacation">Vacation</Select.Option>
              <Select.Option value="Sick Leave">Sick Leave</Select.Option>
              <Select.Option value="Personal">Personal</Select.Option>
              <Select.Option value="Maternity">Maternity</Select.Option>
              <Select.Option value="Paternity">Paternity</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="Leave Duration"
            name="dates"
            rules={[{ required: true, message: 'Please select start and end dates' }]}
          >
            <RangePicker style={{ width: '100%' }} size="large" />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0, marginTop: 24 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setShowModal(false);
                form.resetFields();
              }}>
                Cancel
              </Button>
              <Button type="primary" htmlType="submit">
                Submit Request
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Leave;
