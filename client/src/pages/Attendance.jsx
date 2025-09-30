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
    return <div className="text-center py-8">Loading...</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-800">Attendance</h1>
        <button
          onClick={() => setShowModal(true)}
          className="bg-green-600 text-white px-4 py-2 rounded-lg flex items-center hover:bg-green-700"
        >
          <Clock className="w-5 h-5 mr-2" />
          Clock In
        </button>
      </div>

      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Employee</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Clock In</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Clock Out</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Location</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {attendances.map((att) => (
              <tr key={att.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  {att.employee ? att.employee.name : `Employee ${att.employee_id}`}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {att.date ? new Date(att.date).toLocaleDateString() : 'N/A'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {att.clock_in ? new Date(att.clock_in).toLocaleTimeString() : 'N/A'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {att.clock_out ? new Date(att.clock_out).toLocaleTimeString() : '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">{att.location || '-'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-8 w-full max-w-md">
            <h2 className="text-2xl font-bold mb-4">Clock In</h2>
            <form onSubmit={handleClockIn}>
              <div className="mb-4">
                <label className="block text-gray-700 mb-2">Employee</label>
                <select
                  required
                  value={formData.employee_id}
                  onChange={(e) => setFormData({ ...formData, employee_id: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-green-600"
                >
                  <option value="">Select Employee</option>
                  {employees.map((emp) => (
                    <option key={emp.id} value={emp.id}>
                      {emp.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="mb-4">
                <label className="block text-gray-700 mb-2">Location</label>
                <input
                  type="text"
                  value={formData.location}
                  onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                  placeholder="e.g., Office, Remote"
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-green-600"
                />
              </div>
              <div className="flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-4 py-2 border rounded-lg hover:bg-gray-100"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
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
