import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Chatbot from './components/Chatbot';
import Dashboard from './pages/Dashboard';
import Employees from './pages/Employees';
import Attendance from './pages/Attendance';
import Leave from './pages/Leave';
import Salary from './pages/Salary';
import OrganizationChart from './pages/OrganizationChart';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Dashboard />} />
          <Route path="employees" element={<Employees />} />
          <Route path="org-chart" element={<OrganizationChart />} />
          <Route path="attendance" element={<Attendance />} />
          <Route path="leave" element={<Leave />} />
          <Route path="salary" element={<Salary />} />
        </Route>
      </Routes>
      <Chatbot />
    </Router>
  );
}

export default App;
