import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const employeeAPI = {
  getAll: () => api.get('/employees'),
  getById: (id) => api.get(`/employees/${id}`),
  create: (data) => api.post('/employees', data),
  update: (id, data) => api.put(`/employees/${id}`, data),
};

export const attendanceAPI = {
  getAll: () => api.get('/attendance'),
  clockIn: (data) => api.post('/attendance/clockin', data),
  clockOut: (data) => api.post('/attendance/clockout', data),
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

export const feedbackAPI = {
  create: (data) => api.post('/feedback', data),
  getAll: () => api.get('/feedback'),
};

export const settingsAPI = {
  getAll: () => api.get('/settings'),
  get: (key) => api.get(`/settings/${key}`),
  upsert: (data) => api.post('/settings', data),
  delete: (key) => api.delete(`/settings/${key}`),
};

export default api;
