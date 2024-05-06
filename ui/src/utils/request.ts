import axios from 'axios'
import { message } from 'antd'

export const HOST = ''
axios.defaults.baseURL = HOST

axios.interceptors.request.use(
  config => {
    return config
  },
  error => {
    return Promise.reject(error)
  },
)

axios.interceptors.response.use(
  response => {
    if (response?.status === 200) {
      return response?.data
    }
  },
  error => {
    try {
      message.error(error?.response?.data?.message)
      throw new Error(error)
    } catch (error) {}
  },
)
