import { useState, useEffect } from 'react';
import { DollarSign, Download, FileText } from 'lucide-react';
import { salaryAPI, employeeAPI } from '../api/api';

const Salary = () => {
  const [employees, setEmployees] = useState([]);
  const [showPayslipModal, setShowPayslipModal] = useState(false);
  const [selectedEmployee, setSelectedEmployee] = useState('');
  const [payslipData, setPayslipData] = useState(null);

  useEffect(() => {
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
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-secondary-900">Salary & Payroll</h1>
        <p className="text-secondary-600 mt-1">Manage compensation and payroll operations</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="card hover:shadow-md transition-shadow">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center">
              <Download className="w-6 h-6 text-primary-600" />
            </div>
            <h2 className="text-xl font-semibold text-secondary-900">Export Payroll</h2>
          </div>
          <p className="text-secondary-600 mb-6">
            Export complete salary data for the current period in JSON format for your records.
          </p>
          <button
            onClick={handleExportSalary}
            className="btn btn-primary w-full"
          >
            Export Salary Data
          </button>
        </div>

        <div className="card hover:shadow-md transition-shadow">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 bg-success-100 rounded-lg flex items-center justify-center">
              <FileText className="w-6 h-6 text-success-600" />
            </div>
            <h2 className="text-xl font-semibold text-secondary-900">Generate Payslip</h2>
          </div>
          <p className="text-secondary-600 mb-6">
            Generate a detailed payslip for a specific employee with all salary components.
          </p>
          <button
            onClick={() => setShowPayslipModal(true)}
            className="btn btn-success w-full"
          >
            Generate Payslip
          </button>
        </div>
      </div>

      <div className="card">
        <h2 className="text-xl font-semibold text-secondary-900 mb-4">Salary Information</h2>
        <p className="text-secondary-600 leading-relaxed">
          Manage salary components, payroll processing, and employee compensation here.
          Use the export feature to generate reports for your payroll system.
        </p>
      </div>

      {showPayslipModal && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl p-8 w-full max-w-md transform transition-all">
            <h2 className="text-2xl font-bold text-secondary-900 mb-6">Generate Payslip</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-secondary-700 mb-2">Select Employee</label>
                <select
                  value={selectedEmployee}
                  onChange={(e) => setSelectedEmployee(e.target.value)}
                  className="input-field"
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
                <div className="card bg-secondary-50 border-secondary-200">
                  <h3 className="font-semibold text-secondary-900 mb-2">{payslipData.employee.name}</h3>
                  <p className="text-sm text-secondary-600 mb-3">Period: {payslipData.period}</p>
                  {payslipData.salaries.length > 0 ? (
                    <div className="space-y-2">
                      {payslipData.salaries.map((salary, idx) => (
                        <div key={idx} className="flex justify-between text-sm">
                          <span className="text-secondary-700">{salary.type}</span>
                          <span className="font-semibold text-secondary-900">${salary.amount.toLocaleString()}</span>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-sm text-secondary-500">No salary data available</p>
                  )}
                </div>
              )}
            </div>

            <div className="flex justify-end gap-3 pt-6">
              <button
                onClick={() => {
                  setShowPayslipModal(false);
                  setPayslipData(null);
                  setSelectedEmployee('');
                }}
                className="btn btn-secondary"
              >
                Close
              </button>
              <button
                onClick={handleGeneratePayslip}
                disabled={!selectedEmployee}
                className="btn btn-success disabled:opacity-50 disabled:cursor-not-allowed"
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
