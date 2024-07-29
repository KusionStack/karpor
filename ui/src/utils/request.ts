import axios from 'axios'
import { notification } from 'antd'

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
    if (
      response?.config?.url?.includes('/rest-api') &&
      !response?.data?.success
    ) {
      notification.error({
        message: `${response?.status}`,
        description: `${response?.data?.message}`,
      })
    } else {
      return response?.data
    }
  },
  error => {
    try {
      throw new Error(error)
    } catch (error) {}
  },
)
