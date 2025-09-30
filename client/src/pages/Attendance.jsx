import { useEffect, useState } from 'react';
import { Clock, Plus } from 'lucide-react';
import { attendanceAPI, employeeAPI } from '../api/api';

const Attendance = () => {
  const [attendances, setAttendances] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    employee_id: '',
    location: '',
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [attendanceRes, employeesRes] = await Promise.all([
        attendanceAPI.getAll(),
        employeeAPI.getAll(),
      ]);
      setAttendances(attendanceRes.data);
      setEmployees(employeesRes.data);
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleClockIn = async (e) => {
    e.preventDefault();
    try {
      await attendanceAPI.clockIn({
        employee_id: Number(formData.employee_id),
        location: formData.location,
      });
      setShowModal(false);
      fetchData();
      setFormData({ employee_id: '', location: '' });
    } catch (error) {
      console.error('Error clocking in:', error);
    }
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
          <h1 className="text-3xl font-bold text-secondary-900">Attendance</h1>
          <p className="text-secondary-600 mt-1">Track employee clock-in and clock-out times</p>
        </div>
        <button
          onClick={() => setShowModal(true)}
          className="btn btn-success flex items-center gap-2"
        >
          <Clock className="w-5 h-5" />
          Clock In
        </button>
      </div>

      <div className="table-container">
        <table className="table-base">
          <thead className="table-header">
            <tr>
              <th className="table-header-cell">Employee</th>
              <th className="table-header-cell">Date</th>
              <th className="table-header-cell">Clock In</th>
              <th className="table-header-cell">Clock Out</th>
              <th className="table-header-cell">Location</th>
            </tr>
          </thead>
          <tbody>
            {attendances.length === 0 ? (
              <tr>
                <td colSpan="5" className="text-center py-12 text-secondary-500">
                  <Clock className="w-12 h-12 mx-auto mb-2 text-secondary-300" />
                  <p>No attendance records found</p>
                </td>
              </tr>
            ) : (
              attendances.map((att) => (
                <tr key={att.id} className="table-row">
                  <td className="table-cell">
                    <div className="font-medium text-secondary-900">
                      {att.employee ? att.employee.name : `Employee ${att.employee_id}`}
                    </div>
                  </td>
                  <td className="table-cell text-secondary-600">
                    {att.date ? new Date(att.date).toLocaleDateString() : 'N/A'}
                  </td>
                  <td className="table-cell">
                    <span className="badge badge-success">
                      {att.clock_in ? new Date(att.clock_in).toLocaleTimeString() : 'N/A'}
                    </span>
                  </td>
                  <td className="table-cell">
                    {att.clock_out ? (
                      <span className="badge badge-danger">
                        {new Date(att.clock_out).toLocaleTimeString()}
                      </span>
                    ) : (
                      <span className="text-secondary-400">-</span>
                    )}
                  </td>
                  <td className="table-cell text-secondary-600">{att.location || '-'}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl p-8 w-full max-w-md transform transition-all">
            <h2 className="text-2xl font-bold text-secondary-900 mb-6">Clock In</h2>
            <form onSubmit={handleClockIn} className="space-y-4">
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
                <label className="block text-sm font-medium text-secondary-700 mb-2">Location</label>
                <input
                  type="text"
                  value={formData.location}
                  onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                  placeholder="e.g., Office, Remote"
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
                  className="btn btn-success"
                >
                  Clock In
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Attendance;
