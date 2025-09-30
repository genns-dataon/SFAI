import { Link, Outlet, useLocation } from 'react-router-dom';
import { Users, Clock, Calendar, DollarSign, LayoutDashboard } from 'lucide-react';

const Layout = () => {
  const location = useLocation();

  const navItems = [
    { path: '/', icon: LayoutDashboard, label: 'Dashboard' },
    { path: '/employees', icon: Users, label: 'Employees' },
    { path: '/attendance', icon: Clock, label: 'Attendance' },
    { path: '/leave', icon: Calendar, label: 'Leave' },
    { path: '/salary', icon: DollarSign, label: 'Salary' },
  ];

  return (
    <div className="flex h-screen bg-gray-100">
      <aside className="w-64 bg-blue-900 text-white">
        <div className="p-6">
          <h1 className="text-2xl font-bold">HCM System</h1>
          <p className="text-xs text-blue-200 mt-1">Human Capital Management</p>
        </div>
        <nav className="mt-6">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = location.pathname === item.path;
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex items-center px-6 py-3 hover:bg-blue-800 transition-colors ${
                  isActive ? 'bg-blue-800 border-l-4 border-white' : ''
                }`}
              >
                <Icon className="w-5 h-5 mr-3" />
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>
      </aside>
      <main className="flex-1 overflow-y-auto">
        <div className="p-8">
          <Outlet />
        </div>
      </main>
    </div>
  );
};

export default Layout;
