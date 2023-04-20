import axios from 'axios'

axios.defaults.baseURL = "https://localhost:7443";

axios.interceptors.request.use(function (config) {
  return config;
}, function (error) {
  return Promise.reject(error);
});

axios.interceptors.response.use(function (response) {
  if(response?.status === 200) {
    return response?.data;
  }
}, function (error) {
  return Promise.reject(error);
});
