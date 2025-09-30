import { useEffect, useState } from 'react';
import { Card, Typography, Space, Tag, Avatar, Spin, Empty } from 'antd';
import { UserOutlined, TeamOutlined } from '@ant-design/icons';
import { Tree, TreeNode } from 'react-organizational-chart';
import { employeeAPI } from '../api/api';

const { Title, Text } = Typography;

const OrganizationChart = () => {
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);

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
    } finally {
      setLoading(false);
    }
  };

  const EmployeeNode = ({ employee }) => {
    const initials = employee.name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase();

    return (
      <div style={{ display: 'inline-block' }}>
        <Card
          size="small"
          style={{
            width: 240,
            textAlign: 'center',
            borderRadius: 8,
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          }}
          styles={{ body: { padding: 16 } }}
        >
          <Space direction="vertical" size="small" style={{ width: '100%' }}>
            <Avatar
              style={{
                backgroundColor: '#1890ff',
                fontSize: 18,
              }}
              size={56}
            >
              {initials}
            </Avatar>
            <Title level={5} style={{ margin: '8px 0 0 0' }}>
              {employee.name}
            </Title>
            <Tag color="blue">{employee.job_title}</Tag>
            {employee.department && (
              <Text type="secondary" style={{ fontSize: 12 }}>
                <TeamOutlined /> {employee.department.name}
              </Text>
            )}
          </Space>
        </Card>
      </div>
    );
  };

  const buildTree = (managerId = null) => {
    const directReports = employees.filter(
      (emp) => emp.manager_id === managerId
    );

    if (directReports.length === 0) return null;

    return directReports.map((employee) => (
      <TreeNode key={employee.id} label={<EmployeeNode employee={employee} />}>
        {buildTree(employee.id)}
      </TreeNode>
    ));
  };

  const renderOrganizationChart = () => {
    const topLevelEmployees = employees.filter((emp) => !emp.manager_id);

    if (topLevelEmployees.length === 0) {
      return (
        <Empty
          description="No organization structure found. Please assign managers to employees to see the hierarchy."
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      );
    }

    return (
      <div style={{ overflowX: 'auto', padding: '40px 20px' }}>
        {topLevelEmployees.map((topEmployee) => (
          <div key={topEmployee.id} style={{ marginBottom: 40 }}>
            <Tree label={<EmployeeNode employee={topEmployee} />}>
              {buildTree(topEmployee.id)}
            </Tree>
          </div>
        ))}
      </div>
    );
  };

  return (
    <div style={{ padding: '24px' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <Title level={2} style={{ margin: 0 }}>
            <TeamOutlined /> Organization Chart
          </Title>
          <Text type="secondary">
            Visual representation of company hierarchy and reporting structure
          </Text>
        </div>

        <Card style={{ minHeight: 400 }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: 100 }}>
              <Spin size="large" />
            </div>
          ) : (
            renderOrganizationChart()
          )}
        </Card>
      </Space>
    </div>
  );
};

export default OrganizationChart;
