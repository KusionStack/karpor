import axios from 'axios'
import { message } from 'antd';

axios.interceptors.request.use((config) => {
  return config;
}, (error) => {
  return Promise.reject(error);
});

axios.interceptors.response.use((response) => {
  if (response?.status === 200) {
    return response?.data;
  }
}, (error) => {
  message.error(error?.response?.data?.message || "请求失败，请重试");
  return Promise.reject(error);
});
