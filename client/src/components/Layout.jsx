import { useState, useEffect } from 'react';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { Layout as AntLayout, Menu, Typography, Button, Dropdown } from 'antd';
import {
  DashboardOutlined,
  TeamOutlined,
  ClockCircleOutlined,
  CalendarOutlined,
  DollarOutlined,
  BankOutlined,
  ApartmentOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined
} from '@ant-design/icons';

const { Sider, Content } = AntLayout;
const { Title, Text } = Typography;

const Layout = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [collapsed, setCollapsed] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const [user, setUser] = useState(null);

  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
  }, []);

  // Detect screen size and auto-collapse on mobile
  useEffect(() => {
    const handleResize = () => {
      const mobile = window.innerWidth < 768;
      setIsMobile(mobile);
      if (mobile) {
        setCollapsed(true);
      }
    };

    handleResize(); // Initial check
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // Auto-close menu on mobile when navigating
  useEffect(() => {
    if (isMobile) {
      setCollapsed(true);
    }
  }, [location.pathname, isMobile]);

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
      key: '/org-chart',
      icon: <ApartmentOutlined />,
      label: <Link to="/org-chart">Organization Chart</Link>,
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
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
        trigger={null}
        breakpoint="md"
        collapsedWidth={isMobile ? 0 : 80}
        width={250} 
        style={{ 
          background: '#fff',
          borderRight: '1px solid #f0f0f0',
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          zIndex: 100
        }}
      >
        <div style={{ 
          padding: collapsed ? '24px 12px' : '24px',
          borderBottom: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          justifyContent: collapsed ? 'center' : 'flex-start',
          gap: 12,
          minHeight: 88
        }}>
          <img 
            src="/logo.png" 
            alt="SunFish Logo" 
            style={{
              width: 40,
              height: 40,
              objectFit: 'contain',
              flexShrink: 0
            }}
          />
          {!collapsed && (
            <div style={{ overflow: 'hidden' }}>
              <Title level={4} style={{ margin: 0, fontSize: 18, whiteSpace: 'nowrap' }}>HCM System</Title>
              <Text type="secondary" style={{ fontSize: 12, whiteSpace: 'nowrap' }}>Human Capital Management</Text>
            </div>
          )}
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
      
      <AntLayout style={{ marginLeft: collapsed ? (isMobile ? 0 : 80) : 250, transition: 'margin-left 0.2s' }}>
        <div style={{
          background: '#fff',
          padding: '12px 16px',
          borderBottom: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          position: 'sticky',
          top: 0,
          zIndex: 99
        }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: '16px',
              width: 40,
              height: 40,
            }}
          />
          
          {user && (
            <Dropdown
              menu={{
                items: [
                  {
                    key: 'logout',
                    label: 'Logout',
                    icon: <LogoutOutlined />,
                    onClick: () => {
                      localStorage.removeItem('token');
                      localStorage.removeItem('user');
                      navigate('/login');
                    }
                  }
                ]
              }}
              placement="bottomRight"
            >
              <Button type="text" icon={<UserOutlined />}>
                {user.username}
              </Button>
            </Dropdown>
          )}
        </div>
        
        <Content style={{ background: '#f5f5f5', minHeight: 'calc(100vh - 64px)' }}>
          <div style={{ 
            maxWidth: 1400, 
            margin: '0 auto', 
            padding: isMobile ? '16px' : '32px'
          }}>
            <Outlet />
          </div>
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;
