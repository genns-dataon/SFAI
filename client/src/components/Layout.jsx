import { Link, Outlet, useLocation } from 'react-router-dom';
import { Layout as AntLayout, Menu, Typography } from 'antd';
import {
  DashboardOutlined,
  TeamOutlined,
  ClockCircleOutlined,
  CalendarOutlined,
  DollarOutlined,
  BankOutlined
} from '@ant-design/icons';

const { Sider, Content } = AntLayout;
const { Title, Text } = Typography;

const Layout = () => {
  const location = useLocation();

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link to="/">Dashboard</Link>,
    },
    {
      key: '/employees',
      icon: <TeamOutlined />,
      label: <Link to="/employees">Employees</Link>,
    },
    {
      key: '/attendance',
      icon: <ClockCircleOutlined />,
      label: <Link to="/attendance">Attendance</Link>,
    },
    {
      key: '/leave',
      icon: <CalendarOutlined />,
      label: <Link to="/leave">Leave</Link>,
    },
    {
      key: '/salary',
      icon: <DollarOutlined />,
      label: <Link to="/salary">Salary</Link>,
    },
  ];

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider 
        width={250} 
        style={{ 
          background: '#fff',
          borderRight: '1px solid #f0f0f0'
        }}
      >
        <div style={{ 
          padding: '24px',
          borderBottom: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          gap: 12
        }}>
          <div style={{
            width: 40,
            height: 40,
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            borderRadius: 8,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <BankOutlined style={{ fontSize: 20, color: '#fff' }} />
          </div>
          <div>
            <Title level={4} style={{ margin: 0, fontSize: 18 }}>HCM System</Title>
            <Text type="secondary" style={{ fontSize: 12 }}>Human Capital Management</Text>
          </div>
        </div>
        
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          style={{ 
            border: 'none',
            marginTop: 16,
            paddingLeft: 8,
            paddingRight: 8
          }}
        />
      </Sider>
      
      <Content style={{ background: '#f5f5f5' }}>
        <div style={{ maxWidth: 1400, margin: '0 auto', padding: 32 }}>
          <Outlet />
        </div>
      </Content>
    </AntLayout>
  );
};

export default Layout;
