import axios from 'axios';

const getApiUrl = () => {
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL;
  }
  
  if (typeof window !== 'undefined' && window.location.hostname.includes('replit.dev')) {
    const backendHost = window.location.hostname.replace('-5000', '-8080');
    return `https://${backendHost}/api`;
  }
  
  return 'http://localhost:8080/api';
};

const API_URL = getApiUrl();

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const employeeAPI = {
  getAll: () => api.get('/employees'),
  getById: (id) => api.get(`/employees/${id}`),
  create: (data) => api.post('/employees', data),
  update: (id, data) => api.put(`/employees/${id}`, data),
};

export const attendanceAPI = {
  getAll: () => api.get('/attendance'),
  clockIn: (data) => api.post('/attendance/clockin', data),
};

export const leaveAPI = {
  getAll: () => api.get('/leave'),
  create: (data) => api.post('/leave', data),
};

export const salaryAPI = {
  export: () => api.get('/salary/export'),
  generatePayslip: (employeeId) => api.post('/salary/payslip', { employee_id: employeeId }),
};

export const chatAPI = {
  sendMessage: (message) => api.post('/chat', { message }),
};

export default api;
