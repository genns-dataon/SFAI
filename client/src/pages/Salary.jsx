import { useState } from 'react';
import { DollarSign, Download, FileText } from 'lucide-react';
import { salaryAPI, employeeAPI } from '../api/api';

const Salary = () => {
  const [employees, setEmployees] = useState([]);
  const [showPayslipModal, setShowPayslipModal] = useState(false);
  const [selectedEmployee, setSelectedEmployee] = useState('');
  const [payslipData, setPayslipData] = useState(null);

  useState(() => {
    employeeAPI.getAll().then((res) => setEmployees(res.data));
  }, []);

  const handleExportSalary = async () => {
    try {
      const response = await salaryAPI.export();
      alert(`Salary data exported successfully! Export ID: ${response.data.export_id}`);
    } catch (error) {
      console.error('Error exporting salary:', error);
      alert('Failed to export salary data');
    }
  };

  const handleGeneratePayslip = async () => {
    if (!selectedEmployee) return;
    try {
      const response = await salaryAPI.generatePayslip(Number(selectedEmployee));
      setPayslipData(response.data);
    } catch (error) {
      console.error('Error generating payslip:', error);
      alert('Failed to generate payslip');
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-800 mb-8">Salary & Payroll</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center mb-4">
            <Download className="w-6 h-6 text-blue-600 mr-3" />
            <h2 className="text-xl font-semibold">Export Payroll</h2>
          </div>
          <p className="text-gray-600 mb-4">
            Export complete salary data for the current period in JSON format.
          </p>
          <button
            onClick={handleExportSalary}
            className="w-full bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
          >
            Export Salary Data
          </button>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center mb-4">
            <FileText className="w-6 h-6 text-green-600 mr-3" />
            <h2 className="text-xl font-semibold">Generate Payslip</h2>
          </div>
          <p className="text-gray-600 mb-4">
            Generate a detailed payslip for a specific employee.
          </p>
          <button
            onClick={() => setShowPayslipModal(true)}
            className="w-full bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700"
          >
            Generate Payslip
          </button>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold mb-4">Salary Information</h2>
        <p className="text-gray-600">
          Manage salary components, payroll processing, and employee compensation here.
          Use the export feature to generate reports for your payroll system.
        </p>
      </div>

      {showPayslipModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-8 w-full max-w-md">
            <h2 className="text-2xl font-bold mb-4">Generate Payslip</h2>
            <div className="mb-4">
              <label className="block text-gray-700 mb-2">Select Employee</label>
              <select
                value={selectedEmployee}
                onChange={(e) => setSelectedEmployee(e.target.value)}
                className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-green-600"
              >
                <option value="">Choose an employee...</option>
                {employees.map((emp) => (
                  <option key={emp.id} value={emp.id}>
                    {emp.name}
                  </option>
                ))}
              </select>
            </div>

            {payslipData && (
              <div className="mb-4 p-4 bg-gray-50 rounded-lg">
                <h3 className="font-semibold mb-2">{payslipData.employee.name}</h3>
                <p className="text-sm text-gray-600">Period: {payslipData.period}</p>
                {payslipData.salaries.length > 0 ? (
                  <div className="mt-3">
                    {payslipData.salaries.map((salary, idx) => (
                      <div key={idx} className="flex justify-between text-sm">
                        <span>{salary.type}</span>
                        <span className="font-semibold">${salary.amount.toLocaleString()}</span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-sm text-gray-500 mt-2">No salary data available</p>
                )}
              </div>
            )}

            <div className="flex justify-end space-x-3">
              <button
                onClick={() => {
                  setShowPayslipModal(false);
                  setPayslipData(null);
                  setSelectedEmployee('');
                }}
                className="px-4 py-2 border rounded-lg hover:bg-gray-100"
              >
                Close
              </button>
              <button
                onClick={handleGeneratePayslip}
                disabled={!selectedEmployee}
                className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50"
              >
                Generate
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Salary;
