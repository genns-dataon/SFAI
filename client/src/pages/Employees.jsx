import { useEffect, useState } from 'react';
import { Card, Table, Button, Modal, Form, Input, Select, DatePicker, Space, Tag, Typography, message, Descriptions, Divider, Tabs, InputNumber } from 'antd';
import { PlusOutlined, SearchOutlined, UserOutlined, EditOutlined, EyeOutlined } from '@ant-design/icons';
import { employeeAPI } from '../api/api';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

const Employees = () => {
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [showDetailModal, setShowDetailModal] = useState(false);
  const [selectedEmployee, setSelectedEmployee] = useState(null);
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
        date_of_birth: values.date_of_birth ? values.date_of_birth.format('YYYY-MM-DD') : '',
        probation_end_date: values.probation_end_date ? values.probation_end_date.format('YYYY-MM-DD') : '',
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

  const handleViewDetails = (employee) => {
    setSelectedEmployee(employee);
    setShowDetailModal(true);
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
      employee_number: employee.employee_number,
      date_of_birth: employee.date_of_birth ? dayjs(employee.date_of_birth) : null,
      national_id: employee.national_id,
      tax_id: employee.tax_id,
      marital_status: employee.marital_status,
      employment_type: employee.employment_type,
      employment_status: employee.employment_status,
      job_level: employee.job_level,
      work_location: employee.work_location,
      work_arrangement: employee.work_arrangement,
      base_salary: employee.base_salary,
      pay_frequency: employee.pay_frequency,
      currency: employee.currency,
      bank_account: employee.bank_account,
      benefit_eligibility: employee.benefit_eligibility,
      probation_end_date: employee.probation_end_date ? dayjs(employee.probation_end_date) : null,
      performance_rating: employee.performance_rating,
      skills: employee.skills,
      training_completed: employee.training_completed,
      career_notes: employee.career_notes,
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
      render: (text, record) => (
        <Space style={{ cursor: 'pointer' }} onClick={() => handleViewDetails(record)}>
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
          <Text strong style={{ color: '#1890ff' }}>{text}</Text>
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
        width={800}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          style={{ marginTop: 24 }}
        >
          <Tabs
            defaultActiveKey="1"
            items={[
              {
                key: '1',
                label: 'Basic Info',
                children: (
                  <>
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
                  </>
                ),
              },
              {
                key: '2',
                label: 'Personal & ID',
                children: (
                  <>
                    <Form.Item label="Employee Number" name="employee_number">
                      <Input placeholder="EMP-001" />
                    </Form.Item>

                    <Form.Item label="Date of Birth" name="date_of_birth">
                      <DatePicker style={{ width: '100%' }} />
                    </Form.Item>

                    <Form.Item label="National ID / Passport" name="national_id">
                      <Input placeholder="ID or Passport Number" />
                    </Form.Item>

                    <Form.Item label="Tax ID" name="tax_id">
                      <Input placeholder="Tax Identification Number" />
                    </Form.Item>

                    <Form.Item label="Marital Status" name="marital_status">
                      <Select placeholder="Select marital status" allowClear>
                        <Select.Option value="single">Single</Select.Option>
                        <Select.Option value="married">Married</Select.Option>
                        <Select.Option value="divorced">Divorced</Select.Option>
                        <Select.Option value="widowed">Widowed</Select.Option>
                      </Select>
                    </Form.Item>
                  </>
                ),
              },
              {
                key: '3',
                label: 'Employment',
                children: (
                  <>
                    <Form.Item label="Employment Type" name="employment_type">
                      <Select placeholder="Select employment type" allowClear>
                        <Select.Option value="full-time">Full-time</Select.Option>
                        <Select.Option value="part-time">Part-time</Select.Option>
                        <Select.Option value="contract">Contract</Select.Option>
                        <Select.Option value="intern">Intern</Select.Option>
                      </Select>
                    </Form.Item>

                    <Form.Item label="Employment Status" name="employment_status">
                      <Select placeholder="Select employment status" allowClear>
                        <Select.Option value="active">Active</Select.Option>
                        <Select.Option value="probation">Probation</Select.Option>
                        <Select.Option value="resigned">Resigned</Select.Option>
                        <Select.Option value="terminated">Terminated</Select.Option>
                      </Select>
                    </Form.Item>

                    <Form.Item label="Job Level / Grade" name="job_level">
                      <Input placeholder="e.g., Senior, Mid-level, Junior" />
                    </Form.Item>

                    <Form.Item label="Work Location" name="work_location">
                      <Input placeholder="e.g., New York Office, Singapore" />
                    </Form.Item>

                    <Form.Item label="Work Arrangement" name="work_arrangement">
                      <Select placeholder="Select work arrangement" allowClear>
                        <Select.Option value="onsite">Onsite</Select.Option>
                        <Select.Option value="hybrid">Hybrid</Select.Option>
                        <Select.Option value="remote">Remote</Select.Option>
                      </Select>
                    </Form.Item>
                  </>
                ),
              },
              {
                key: '4',
                label: 'Compensation',
                children: (
                  <>
                    <Form.Item label="Base Salary" name="base_salary">
                      <InputNumber
                        style={{ width: '100%' }}
                        placeholder="50000"
                        min={0}
                        formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                        parser={(value) => value.replace(/\$\s?|(,*)/g, '')}
                      />
                    </Form.Item>

                    <Form.Item label="Pay Frequency" name="pay_frequency">
                      <Select placeholder="Select pay frequency" allowClear>
                        <Select.Option value="monthly">Monthly</Select.Option>
                        <Select.Option value="biweekly">Biweekly</Select.Option>
                        <Select.Option value="weekly">Weekly</Select.Option>
                        <Select.Option value="annual">Annual</Select.Option>
                      </Select>
                    </Form.Item>

                    <Form.Item label="Currency" name="currency">
                      <Select placeholder="Select currency" allowClear>
                        <Select.Option value="USD">USD</Select.Option>
                        <Select.Option value="EUR">EUR</Select.Option>
                        <Select.Option value="GBP">GBP</Select.Option>
                        <Select.Option value="IDR">IDR</Select.Option>
                        <Select.Option value="SGD">SGD</Select.Option>
                      </Select>
                    </Form.Item>

                    <Form.Item label="Bank Account" name="bank_account">
                      <Input placeholder="Bank account number" />
                    </Form.Item>

                    <Form.Item label="Benefit Eligibility" name="benefit_eligibility">
                      <Input.TextArea
                        placeholder="e.g., Medical, Housing, Pension"
                        rows={3}
                      />
                    </Form.Item>
                  </>
                ),
              },
              {
                key: '5',
                label: 'Performance',
                children: (
                  <>
                    <Form.Item label="Probation End Date" name="probation_end_date">
                      <DatePicker style={{ width: '100%' }} />
                    </Form.Item>

                    <Form.Item label="Performance Rating" name="performance_rating">
                      <Select placeholder="Select performance rating" allowClear>
                        <Select.Option value="excellent">Excellent</Select.Option>
                        <Select.Option value="good">Good</Select.Option>
                        <Select.Option value="satisfactory">Satisfactory</Select.Option>
                        <Select.Option value="needs-improvement">Needs Improvement</Select.Option>
                      </Select>
                    </Form.Item>

                    <Form.Item label="Skills / Certifications" name="skills">
                      <Input.TextArea
                        placeholder="e.g., React, Python, AWS Certified"
                        rows={3}
                      />
                    </Form.Item>

                    <Form.Item label="Training Completed" name="training_completed">
                      <Input.TextArea
                        placeholder="e.g., Leadership Training, Safety Certification"
                        rows={3}
                      />
                    </Form.Item>

                    <Form.Item label="Career / Succession Notes" name="career_notes">
                      <Input.TextArea
                        placeholder="Career development notes"
                        rows={3}
                      />
                    </Form.Item>
                  </>
                ),
              },
            ]}
          />

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

      <Modal
        title={
          <Space>
            <EyeOutlined />
            <span>Employee Details</span>
          </Space>
        }
        open={showDetailModal}
        onCancel={() => {
          setShowDetailModal(false);
          setSelectedEmployee(null);
        }}
        footer={[
          <Button key="close" onClick={() => setShowDetailModal(false)}>
            Close
          </Button>,
          <Button key="edit" type="primary" onClick={() => {
            setShowDetailModal(false);
            handleEdit(selectedEmployee);
          }}>
            Edit Employee
          </Button>,
        ]}
        width={900}
      >
        {selectedEmployee && (
          <div style={{ marginTop: 24 }}>
            <Title level={4}>Basic Information</Title>
            <Descriptions bordered column={2}>
              <Descriptions.Item label="Name">{selectedEmployee.name}</Descriptions.Item>
              <Descriptions.Item label="Email">{selectedEmployee.email}</Descriptions.Item>
              <Descriptions.Item label="Job Title">{selectedEmployee.job_title}</Descriptions.Item>
              <Descriptions.Item label="Department">{selectedEmployee.department?.name || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Manager">{selectedEmployee.manager?.name || 'None'}</Descriptions.Item>
              <Descriptions.Item label="Hire Date">
                {selectedEmployee.hire_date ? new Date(selectedEmployee.hire_date).toLocaleDateString() : 'N/A'}
              </Descriptions.Item>
            </Descriptions>

            <Divider />
            <Title level={4}>Personal & Identification</Title>
            <Descriptions bordered column={2}>
              <Descriptions.Item label="Employee Number">{selectedEmployee.employee_number || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Date of Birth">
                {selectedEmployee.date_of_birth ? new Date(selectedEmployee.date_of_birth).toLocaleDateString() : 'N/A'}
              </Descriptions.Item>
              <Descriptions.Item label="National ID">{selectedEmployee.national_id || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Tax ID">{selectedEmployee.tax_id || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Marital Status">{selectedEmployee.marital_status || 'N/A'}</Descriptions.Item>
            </Descriptions>

            <Divider />
            <Title level={4}>Employment & Job Details</Title>
            <Descriptions bordered column={2}>
              <Descriptions.Item label="Employment Type">{selectedEmployee.employment_type || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Employment Status">
                <Tag color={selectedEmployee.employment_status === 'active' ? 'green' : 'orange'}>
                  {selectedEmployee.employment_status || 'N/A'}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="Job Level">{selectedEmployee.job_level || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Work Location">{selectedEmployee.work_location || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Work Arrangement">{selectedEmployee.work_arrangement || 'N/A'}</Descriptions.Item>
            </Descriptions>

            <Divider />
            <Title level={4}>Compensation & Benefits</Title>
            <Descriptions bordered column={2}>
              <Descriptions.Item label="Base Salary">
                {selectedEmployee.base_salary ? `${selectedEmployee.currency || 'USD'} ${selectedEmployee.base_salary.toLocaleString()}` : 'N/A'}
              </Descriptions.Item>
              <Descriptions.Item label="Pay Frequency">{selectedEmployee.pay_frequency || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Currency">{selectedEmployee.currency || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Bank Account">{selectedEmployee.bank_account || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Benefit Eligibility" span={2}>{selectedEmployee.benefit_eligibility || 'N/A'}</Descriptions.Item>
            </Descriptions>

            <Divider />
            <Title level={4}>Performance & Development</Title>
            <Descriptions bordered column={2}>
              <Descriptions.Item label="Probation End Date">
                {selectedEmployee.probation_end_date ? new Date(selectedEmployee.probation_end_date).toLocaleDateString() : 'N/A'}
              </Descriptions.Item>
              <Descriptions.Item label="Performance Rating">{selectedEmployee.performance_rating || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Skills" span={2}>{selectedEmployee.skills || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Training Completed" span={2}>{selectedEmployee.training_completed || 'N/A'}</Descriptions.Item>
              <Descriptions.Item label="Career Notes" span={2}>{selectedEmployee.career_notes || 'N/A'}</Descriptions.Item>
            </Descriptions>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default Employees;
