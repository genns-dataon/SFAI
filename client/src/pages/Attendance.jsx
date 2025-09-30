import { useEffect, useState } from 'react';
import { Card, Table, Button, Modal, Form, Select, Input, Typography, Space, Tag, message } from 'antd';
import { ClockCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { attendanceAPI, employeeAPI } from '../api/api';

const { Title, Text } = Typography;

const Attendance = () => {
  const [attendances, setAttendances] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showClockInModal, setShowClockInModal] = useState(false);
  const [showClockOutModal, setShowClockOutModal] = useState(false);
  const [clockInForm] = Form.useForm();
  const [clockOutForm] = Form.useForm();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [attendanceRes, employeesRes] = await Promise.all([
        attendanceAPI.getAll(),
        employeeAPI.getAll(),
      ]);
      setAttendances(attendanceRes.data);
      setEmployees(employeesRes.data);
    } catch (error) {
      console.error('Error fetching data:', error);
      message.error('Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  const handleClockIn = async (values) => {
    try {
      await attendanceAPI.clockIn({
        employee_id: Number(values.employee_id),
        location: values.location || '',
      });
      message.success('Clocked in successfully');
      setShowClockInModal(false);
      clockInForm.resetFields();
      fetchData();
    } catch (error) {
      console.error('Error clocking in:', error);
      message.error('Failed to clock in');
    }
  };

  const handleClockOut = async (values) => {
    try {
      const response = await attendanceAPI.clockOut({
        employee_id: Number(values.employee_id),
        location: values.location || '',
      });
      
      const duration = response.data.duration;
      const durationText = duration 
        ? ` Total time: ${duration.hours}h ${duration.minutes}m`
        : '';
      
      message.success(`Clocked out successfully!${durationText}`);
      setShowClockOutModal(false);
      clockOutForm.resetFields();
      fetchData();
    } catch (error) {
      console.error('Error clocking out:', error);
      message.error(error.response?.data?.error || 'Failed to clock out');
    }
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
      title: 'Date',
      dataIndex: 'date',
      key: 'date',
      render: (date) => date ? new Date(date).toLocaleDateString() : 'N/A',
    },
    {
      title: 'Clock In',
      dataIndex: 'clock_in',
      key: 'clock_in',
      render: (time) => time ? (
        <Tag color="green">{new Date(time).toLocaleTimeString()}</Tag>
      ) : 'N/A',
    },
    {
      title: 'Clock Out',
      dataIndex: 'clock_out',
      key: 'clock_out',
      render: (time) => time ? (
        <Tag color="red">{new Date(time).toLocaleTimeString()}</Tag>
      ) : (
        <Text type="secondary">-</Text>
      ),
    },
    {
      title: 'Location',
      dataIndex: 'location',
      key: 'location',
      render: (location) => location || '-',
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <Title level={2} style={{ margin: 0 }}>Attendance</Title>
            <Text type="secondary">Track employee clock-in and clock-out times</Text>
          </div>
          <Space>
            <Button
              type="primary"
              icon={<ClockCircleOutlined />}
              onClick={() => setShowClockInModal(true)}
              size="large"
              style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
            >
              Clock In
            </Button>
            <Button
              type="primary"
              icon={<ClockCircleOutlined />}
              onClick={() => setShowClockOutModal(true)}
              size="large"
              style={{ backgroundColor: '#ff4d4f', borderColor: '#ff4d4f' }}
            >
              Clock Out
            </Button>
          </Space>
        </div>

        <Card>
          <Table
            columns={columns}
            dataSource={attendances}
            loading={loading}
            rowKey="id"
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showTotal: (total) => `Total ${total} records`,
            }}
          />
        </Card>
      </Space>

      <Modal
        title={
          <Space>
            <ClockCircleOutlined />
            <span>Clock In</span>
          </Space>
        }
        open={showClockInModal}
        onCancel={() => {
          setShowClockInModal(false);
          clockInForm.resetFields();
        }}
        footer={null}
        width={500}
      >
        <Form
          form={clockInForm}
          layout="vertical"
          onFinish={handleClockIn}
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
            label="Location"
            name="location"
          >
            <Input placeholder="e.g., Office, Remote" size="large" />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0, marginTop: 24 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setShowClockInModal(false);
                clockInForm.resetFields();
              }}>
                Cancel
              </Button>
              <Button 
                type="primary" 
                htmlType="submit"
                style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
              >
                Clock In
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={
          <Space>
            <ClockCircleOutlined />
            <span>Clock Out</span>
          </Space>
        }
        open={showClockOutModal}
        onCancel={() => {
          setShowClockOutModal(false);
          clockOutForm.resetFields();
        }}
        footer={null}
        width={500}
      >
        <Form
          form={clockOutForm}
          layout="vertical"
          onFinish={handleClockOut}
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
            label="Location"
            name="location"
          >
            <Input placeholder="e.g., Office, Remote" size="large" />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0, marginTop: 24 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setShowClockOutModal(false);
                clockOutForm.resetFields();
              }}>
                Cancel
              </Button>
              <Button 
                type="primary" 
                htmlType="submit"
                style={{ backgroundColor: '#ff4d4f', borderColor: '#ff4d4f' }}
              >
                Clock Out
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Attendance;
