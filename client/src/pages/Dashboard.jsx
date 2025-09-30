import { useEffect, useState } from 'react';
import { Users, Clock, Calendar, TrendingUp, ArrowUp } from 'lucide-react';
import { employeeAPI, attendanceAPI, leaveAPI } from '../api/api';

const Dashboard = () => {
  const [stats, setStats] = useState({
    totalEmployees: 0,
    todayAttendance: 0,
    pendingLeaves: 0,
  });

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
      }
    };

    fetchStats();
  }, []);

  const cards = [
    { 
      title: 'Total Employees', 
      value: stats.totalEmployees, 
      icon: Users, 
      bgColor: 'bg-primary-50', 
      iconColor: 'bg-primary-600',
      textColor: 'text-primary-700',
      trend: '+12%'
    },
    { 
      title: "Today's Attendance", 
      value: stats.todayAttendance, 
      icon: Clock, 
      bgColor: 'bg-success-50', 
      iconColor: 'bg-success-600',
      textColor: 'text-success-700',
      trend: '+5%'
    },
    { 
      title: 'Pending Leaves', 
      value: stats.pendingLeaves, 
      icon: Calendar, 
      bgColor: 'bg-warning-50', 
      iconColor: 'bg-warning-600',
      textColor: 'text-warning-700',
      trend: '-3%'
    },
    { 
      title: 'Departments', 
      value: 3, 
      icon: TrendingUp, 
      bgColor: 'bg-purple-50', 
      iconColor: 'bg-purple-600',
      textColor: 'text-purple-700',
      trend: '0%'
    },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-secondary-900">Dashboard</h1>
        <p className="text-secondary-600 mt-1">Welcome back! Here's what's happening today.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {cards.map((card, idx) => {
          const Icon = card.icon;
          return (
            <div key={idx} className="bg-white rounded-xl shadow-sm border border-secondary-200 p-6 hover:shadow-md transition-shadow duration-200">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="text-secondary-600 text-sm font-medium">{card.title}</p>
                  <p className="text-3xl font-bold text-secondary-900 mt-3">{card.value}</p>
                  <div className="flex items-center gap-1 mt-2">
                    <ArrowUp className="w-4 h-4 text-success-600" />
                    <span className="text-sm font-medium text-success-600">{card.trend}</span>
                    <span className="text-sm text-secondary-500 ml-1">vs last month</span>
                  </div>
                </div>
                <div className={`${card.iconColor} p-3 rounded-lg`}>
                  <Icon className="w-6 h-6 text-white" />
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-white rounded-xl shadow-sm border border-secondary-200 p-6">
          <h2 className="text-xl font-semibold text-secondary-900 mb-4">Welcome to the HCM System</h2>
          <p className="text-secondary-600 leading-relaxed">
            Manage your workforce efficiently with our comprehensive Human Capital Management system.
            Navigate through the sidebar to access employees, attendance, leave requests, and salary information.
          </p>
          <div className="mt-6 flex gap-3">
            <button className="btn btn-primary">View All Employees</button>
            <button className="btn btn-secondary">Generate Report</button>
          </div>
        </div>

        <div className="bg-gradient-to-br from-primary-600 to-primary-700 rounded-xl shadow-sm p-6 text-white">
          <h3 className="text-lg font-semibold mb-2">Quick Actions</h3>
          <p className="text-primary-100 text-sm mb-4">Perform common tasks quickly</p>
          <div className="space-y-2">
            <button className="w-full text-left px-4 py-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors">
              Add New Employee
            </button>
            <button className="w-full text-left px-4 py-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors">
              Clock In/Out
            </button>
            <button className="w-full text-left px-4 py-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors">
              Request Leave
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
