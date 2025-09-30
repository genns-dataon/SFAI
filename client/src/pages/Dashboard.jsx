import { useEffect, useState } from 'react';
import { Users, Clock, Calendar, TrendingUp } from 'lucide-react';
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
    { title: 'Total Employees', value: stats.totalEmployees, icon: Users, color: 'bg-blue-500' },
    { title: "Today's Attendance", value: stats.todayAttendance, icon: Clock, color: 'bg-green-500' },
    { title: 'Pending Leaves', value: stats.pendingLeaves, icon: Calendar, color: 'bg-yellow-500' },
    { title: 'Departments', value: 3, icon: TrendingUp, color: 'bg-purple-500' },
  ];

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-800 mb-8">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {cards.map((card, idx) => {
          const Icon = card.icon;
          return (
            <div key={idx} className="bg-white rounded-lg shadow-md p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">{card.title}</p>
                  <p className="text-3xl font-bold text-gray-800 mt-2">{card.value}</p>
                </div>
                <div className={`${card.color} p-3 rounded-full`}>
                  <Icon className="w-6 h-6 text-white" />
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="mt-8 bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold mb-4">Welcome to the HCM System</h2>
        <p className="text-gray-600">
          Manage your workforce efficiently with our comprehensive Human Capital Management system.
          Navigate through the sidebar to access employees, attendance, leave requests, and salary information.
        </p>
      </div>
    </div>
  );
};

export default Dashboard;
