import { Link, Outlet, useLocation } from 'react-router-dom';
import { Users, Clock, Calendar, DollarSign, LayoutDashboard, Building2 } from 'lucide-react';
import Chatbot from './Chatbot';

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
    <div className="flex h-screen bg-secondary-50">
      <aside className="w-64 bg-white border-r border-secondary-200 flex flex-col">
        <div className="p-6 border-b border-secondary-200">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-primary-600 rounded-lg flex items-center justify-center">
              <Building2 className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-xl font-bold text-secondary-900">HCM System</h1>
              <p className="text-xs text-secondary-500">Human Capital Management</p>
            </div>
          </div>
        </div>
        <nav className="flex-1 px-3 py-4 space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = location.pathname === item.path;
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex items-center gap-3 px-3 py-2.5 rounded-lg font-medium transition-all duration-200 ${
                  isActive 
                    ? 'bg-primary-50 text-primary-700' 
                    : 'text-secondary-600 hover:bg-secondary-100 hover:text-secondary-900'
                }`}
              >
                <Icon className={`w-5 h-5 ${isActive ? 'text-primary-600' : ''}`} />
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>
      </aside>
      <main className="flex-1 overflow-y-auto bg-secondary-50">
        <div className="max-w-7xl mx-auto p-8">
          <Outlet />
        </div>
      </main>
      <Chatbot />
    </div>
  );
};

export default Layout;
