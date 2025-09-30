import { useEffect, useState } from 'react';
import { Card, Row, Col, Statistic, Button, Typography, Space, message } from 'antd';
import { 
  UserOutlined, 
  ClockCircleOutlined, 
  CalendarOutlined, 
  RiseOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined 
} from '@ant-design/icons';
import { employeeAPI, attendanceAPI, leaveAPI } from '../api/api';

const { Title, Paragraph } = Typography;

const Dashboard = () => {
  const [stats, setStats] = useState({
    totalEmployees: 0,
    todayAttendance: 0,
    pendingLeaves: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [employeesRes, attendanceRes, leavesRes] = await Promise.all([
          employeeAPI.getAll(),
          attendanceAPI.getAll(),
          leaveAPI.getAll(),
        ]);

        const today = new Date().toISOString().split('T')[0];
        const todayAttendance = attendanceRes.data.filter(
          (a) => a.date && a.date.split('T')[0] === today
        ).length;

        const pendingLeaves = leavesRes.data.filter((l) => l.status === 'pending').length;

        setStats({
          totalEmployees: employeesRes.data.length,
          todayAttendance,
          pendingLeaves,
        });
      } catch (error) {
        console.error('Error fetching stats:', error);
        message.error('Failed to load dashboard statistics');
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  const statsCards = [
    {
      title: 'Total Employees',
      value: stats.totalEmployees,
      icon: <UserOutlined style={{ fontSize: '24px', color: '#1890ff' }} />,
      trend: 12,
      trendLabel: 'vs last month'
    },
    {
      title: "Today's Attendance",
      value: stats.todayAttendance,
      icon: <ClockCircleOutlined style={{ fontSize: '24px', color: '#52c41a' }} />,
      trend: 5,
      trendLabel: 'vs last month'
    },
    {
      title: 'Pending Leaves',
      value: stats.pendingLeaves,
      icon: <CalendarOutlined style={{ fontSize: '24px', color: '#faad14' }} />,
      trend: -3,
      trendLabel: 'vs last month'
    },
    {
      title: 'Departments',
      value: 3,
      icon: <RiseOutlined style={{ fontSize: '24px', color: '#722ed1' }} />,
      trend: 0,
      trendLabel: 'vs last month'
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Title level={2}>Dashboard</Title>
        <Paragraph type="secondary">Welcome back! Here's what's happening today.</Paragraph>
      </div>

      <Row gutter={[16, 16]}>
        {statsCards.map((stat, index) => (
          <Col xs={24} sm={12} lg={6} key={index}>
            <Card loading={loading} hoverable>
              <Space direction="vertical" size="small" style={{ width: '100%' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Paragraph type="secondary" style={{ margin: 0 }}>{stat.title}</Paragraph>
                  {stat.icon}
                </div>
                <Statistic
                  value={stat.value}
                  valueStyle={{ fontSize: '28px', fontWeight: 600 }}
                />
                <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
                  {stat.trend > 0 ? (
                    <>
                      <ArrowUpOutlined style={{ color: '#52c41a' }} />
                      <span style={{ color: '#52c41a', fontWeight: 500 }}>{stat.trend}%</span>
                    </>
                  ) : stat.trend < 0 ? (
                    <>
                      <ArrowDownOutlined style={{ color: '#ff4d4f' }} />
                      <span style={{ color: '#ff4d4f', fontWeight: 500 }}>{Math.abs(stat.trend)}%</span>
                    </>
                  ) : (
                    <span style={{ color: '#8c8c8c', fontWeight: 500 }}>0%</span>
                  )}
                  <span style={{ color: '#8c8c8c', fontSize: '14px', marginLeft: 4 }}>{stat.trendLabel}</span>
                </div>
              </Space>
            </Card>
          </Col>
        ))}
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={16}>
          <Card title="Welcome to the HCM System">
            <Paragraph>
              Manage your workforce efficiently with our comprehensive Human Capital Management system.
              Navigate through the sidebar to access employees, attendance, leave requests, and salary information.
            </Paragraph>
            <Space style={{ marginTop: 16 }}>
              <Button type="primary" size="large">View All Employees</Button>
              <Button size="large">Generate Report</Button>
            </Space>
          </Card>
        </Col>

        <Col xs={24} lg={8}>
          <Card 
            title="Quick Actions"
            style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
            styles={{ 
              header: { color: '#fff', borderBottom: '1px solid rgba(255,255,255,0.2)' },
              body: { color: '#fff' }
            }}
          >
            <Paragraph style={{ color: 'rgba(255,255,255,0.8)' }}>
              Perform common tasks quickly
            </Paragraph>
            <Space direction="vertical" style={{ width: '100%' }}>
              <Button 
                block 
                size="large" 
                style={{ 
                  background: 'rgba(255,255,255,0.1)', 
                  border: 'none', 
                  color: '#fff',
                  textAlign: 'left'
                }}
              >
                Add New Employee
              </Button>
              <Button 
                block 
                size="large" 
                style={{ 
                  background: 'rgba(255,255,255,0.1)', 
                  border: 'none', 
                  color: '#fff',
                  textAlign: 'left'
                }}
              >
                Clock In/Out
              </Button>
              <Button 
                block 
                size="large" 
                style={{ 
                  background: 'rgba(255,255,255,0.1)', 
                  border: 'none', 
                  color: '#fff',
                  textAlign: 'left'
                }}
              >
                Request Leave
              </Button>
            </Space>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
