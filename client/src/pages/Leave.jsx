import { useEffect, useState } from 'react';
import { Calendar, Plus } from 'lucide-react';
import { leaveAPI, employeeAPI } from '../api/api';

const Leave = () => {
  const [leaves, setLeaves] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    employee_id: '',
    leave_type: 'Vacation',
    start_date: '',
    end_date: '',
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [leavesRes, employeesRes] = await Promise.all([
        leaveAPI.getAll(),
        employeeAPI.getAll(),
      ]);
      setLeaves(leavesRes.data);
      setEmployees(employeesRes.data);
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await leaveAPI.create({
        employee_id: Number(formData.employee_id),
        leave_type: formData.leave_type,
        start_date: new Date(formData.start_date).toISOString(),
        end_date: new Date(formData.end_date).toISOString(),
      });
      setShowModal(false);
      fetchData();
      setFormData({ employee_id: '', leave_type: 'Vacation', start_date: '', end_date: '' });
    } catch (error) {
      console.error('Error creating leave request:', error);
    }
  };

  const getStatusBadge = (status) => {
    const badges = {
      pending: 'badge-warning',
      approved: 'badge-success',
      rejected: 'badge-danger',
    };
    return badges[status] || badges.pending;
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-secondary-900">Leave Requests</h1>
          <p className="text-secondary-600 mt-1">Manage employee leave applications and approvals</p>
        </div>
        <button
          onClick={() => setShowModal(true)}
          className="btn btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          Request Leave
        </button>
      </div>

      <div className="table-container">
        <table className="table-base">
          <thead className="table-header">
            <tr>
              <th className="table-header-cell">Employee</th>
              <th className="table-header-cell">Leave Type</th>
              <th className="table-header-cell">Start Date</th>
              <th className="table-header-cell">End Date</th>
              <th className="table-header-cell">Status</th>
            </tr>
          </thead>
          <tbody>
            {leaves.length === 0 ? (
              <tr>
                <td colSpan="5" className="text-center py-12 text-secondary-500">
                  <Calendar className="w-12 h-12 mx-auto mb-2 text-secondary-300" />
                  <p>No leave requests found</p>
                </td>
              </tr>
            ) : (
              leaves.map((leave) => (
                <tr key={leave.id} className="table-row">
                  <td className="table-cell">
                    <div className="font-medium text-secondary-900">
                      {leave.employee ? leave.employee.name : `Employee ${leave.employee_id}`}
                    </div>
                  </td>
                  <td className="table-cell">
                    <span className="badge badge-info">{leave.leave_type}</span>
                  </td>
                  <td className="table-cell text-secondary-600">
                    {leave.start_date ? new Date(leave.start_date).toLocaleDateString() : 'N/A'}
                  </td>
                  <td className="table-cell text-secondary-600">
                    {leave.end_date ? new Date(leave.end_date).toLocaleDateString() : 'N/A'}
                  </td>
                  <td className="table-cell">
                    <span className={`badge ${getStatusBadge(leave.status)} capitalize`}>
                      {leave.status}
                    </span>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl p-8 w-full max-w-md transform transition-all">
            <h2 className="text-2xl font-bold text-secondary-900 mb-6">Request Leave</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-secondary-700 mb-2">Employee</label>
                <select
                  required
                  value={formData.employee_id}
                  onChange={(e) => setFormData({ ...formData, employee_id: e.target.value })}
                  className="input-field"
                >
                  <option value="">Select Employee</option>
                  {employees.map((emp) => (
                    <option key={emp.id} value={emp.id}>
                      {emp.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-secondary-700 mb-2">Leave Type</label>
                <select
                  value={formData.leave_type}
                  onChange={(e) => setFormData({ ...formData, leave_type: e.target.value })}
                  className="input-field"
                >
                  <option value="Vacation">Vacation</option>
                  <option value="Sick Leave">Sick Leave</option>
                  <option value="Personal">Personal</option>
                  <option value="Maternity">Maternity</option>
                  <option value="Paternity">Paternity</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-secondary-700 mb-2">Start Date</label>
                <input
                  type="date"
                  required
                  value={formData.start_date}
                  onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-secondary-700 mb-2">End Date</label>
                <input
                  type="date"
                  required
                  value={formData.end_date}
                  onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                  className="input-field"
                />
              </div>
              <div className="flex justify-end gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="btn btn-secondary"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn btn-primary"
                >
                  Submit Request
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Leave;
