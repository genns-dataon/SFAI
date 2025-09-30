import { useState, useEffect } from 'react';
import { Card, Button, Modal, Select, Typography, Space, Divider, Row, Col, message } from 'antd';
import { DownloadOutlined, FileTextOutlined, DollarOutlined } from '@ant-design/icons';
import { salaryAPI, employeeAPI } from '../api/api';

const { Title, Text, Paragraph } = Typography;

const Salary = () => {
  const [employees, setEmployees] = useState([]);
  const [showPayslipModal, setShowPayslipModal] = useState(false);
  const [selectedEmployee, setSelectedEmployee] = useState(undefined);
  const [payslipData, setPayslipData] = useState(null);

  useEffect(() => {
    employeeAPI.getAll().then((res) => setEmployees(res.data));
  }, []);

  const handleExportSalary = async () => {
    try {
      const response = await salaryAPI.export();
      message.success(`Salary data exported successfully! Export ID: ${response.data.export_id}`);
    } catch (error) {
      console.error('Error exporting salary:', error);
      message.error('Failed to export salary data');
    }
  };

  const handleGeneratePayslip = async () => {
    if (!selectedEmployee) {
      message.warning('Please select an employee first');
      return;
    }
    try {
      const response = await salaryAPI.generatePayslip(Number(selectedEmployee));
      setPayslipData(response.data);
      message.success('Payslip generated successfully');
    } catch (error) {
      console.error('Error generating payslip:', error);
      message.error('Failed to generate payslip');
    }
  };

  return (
    <div style={{ padding: '24px' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <Title level={2} style={{ margin: 0 }}>Salary & Payroll</Title>
          <Text type="secondary">Manage compensation and payroll operations</Text>
        </div>

        <Row gutter={[24, 24]}>
          <Col xs={24} lg={12}>
            <Card
              hoverable
              style={{ height: '100%' }}
            >
              <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                  <div style={{
                    width: 48,
                    height: 48,
                    backgroundColor: '#e6f7ff',
                    borderRadius: 8,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center'
                  }}>
                    <DownloadOutlined style={{ fontSize: 24, color: '#1890ff' }} />
                  </div>
                  <Title level={4} style={{ margin: 0 }}>Export Payroll</Title>
                </div>
                <Paragraph type="secondary">
                  Export complete salary data for the current period in JSON format for your records.
                </Paragraph>
                <Button
                  type="primary"
                  icon={<DownloadOutlined />}
                  onClick={handleExportSalary}
                  block
                  size="large"
                >
                  Export Salary Data
                </Button>
              </Space>
            </Card>
          </Col>

          <Col xs={24} lg={12}>
            <Card
              hoverable
              style={{ height: '100%' }}
            >
              <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                  <div style={{
                    width: 48,
                    height: 48,
                    backgroundColor: '#f6ffed',
                    borderRadius: 8,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center'
                  }}>
                    <FileTextOutlined style={{ fontSize: 24, color: '#52c41a' }} />
                  </div>
                  <Title level={4} style={{ margin: 0 }}>Generate Payslip</Title>
                </div>
                <Paragraph type="secondary">
                  Generate a detailed payslip for a specific employee with all salary components.
                </Paragraph>
                <Button
                  type="primary"
                  icon={<FileTextOutlined />}
                  onClick={() => setShowPayslipModal(true)}
                  block
                  size="large"
                  style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
                >
                  Generate Payslip
                </Button>
              </Space>
            </Card>
          </Col>
        </Row>

        <Card>
          <Space direction="vertical" size="small">
            <Title level={4}>Salary Information</Title>
            <Paragraph type="secondary">
              Manage salary components, payroll processing, and employee compensation here.
              Use the export feature to generate reports for your payroll system.
            </Paragraph>
          </Space>
        </Card>
      </Space>

      <Modal
        title={
          <Space>
            <DollarOutlined />
            <span>Generate Payslip</span>
          </Space>
        }
        open={showPayslipModal}
        onCancel={() => {
          setShowPayslipModal(false);
          setPayslipData(null);
          setSelectedEmployee(undefined);
        }}
        footer={[
          <Button
            key="close"
            onClick={() => {
              setShowPayslipModal(false);
              setPayslipData(null);
              setSelectedEmployee(undefined);
            }}
          >
            Close
          </Button>,
          <Button
            key="generate"
            type="primary"
            onClick={handleGeneratePayslip}
            disabled={!selectedEmployee}
            style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
          >
            Generate
          </Button>,
        ]}
        width={600}
      >
        <Space direction="vertical" size="large" style={{ width: '100%', marginTop: 24 }}>
          <div>
            <Text strong>Select Employee</Text>
            <Select
              placeholder="Choose an employee..."
              value={selectedEmployee}
              onChange={setSelectedEmployee}
              style={{ width: '100%', marginTop: 8 }}
              size="large"
            >
              {employees.map((emp) => (
                <Select.Option key={emp.id} value={emp.id}>
                  {emp.name}
                </Select.Option>
              ))}
            </Select>
          </div>

          {payslipData && (
            <Card style={{ backgroundColor: '#fafafa' }}>
              <Space direction="vertical" style={{ width: '100%' }}>
                <Title level={5} style={{ margin: 0 }}>{payslipData.employee.name}</Title>
                <Text type="secondary">Period: {payslipData.period}</Text>
                <Divider style={{ margin: '12px 0' }} />
                {payslipData.salaries.length > 0 ? (
                  <Space direction="vertical" style={{ width: '100%' }}>
                    {payslipData.salaries.map((salary, idx) => (
                      <div key={idx} style={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Text type="secondary">{salary.type}</Text>
                        <Text strong>${salary.amount.toLocaleString()}</Text>
                      </div>
                    ))}
                  </Space>
                ) : (
                  <Text type="secondary">No salary data available</Text>
                )}
              </Space>
            </Card>
          )}
        </Space>
      </Modal>
    </div>
  );
};

export default Salary;
