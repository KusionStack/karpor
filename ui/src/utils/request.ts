import axios from 'axios'
import { notification } from 'antd'
import { useState, useEffect } from 'react'
import { useDispatch } from 'react-redux'
import { setIsLogin } from '@/store/modules/globalSlice'

export const HOST = 'https://karpor-demo.kusionstack.io/'
axios.defaults.baseURL = HOST

axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (!config?.headers?.Authorization && token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  },
)

axios.interceptors.response.use(
  response => {
    return response
  },
  error => {
    try {
      const { response } = error || {}
      if (response) {
        return Promise.resolve(response)
      }
    } catch (error) {}
  },
)

export const useAxios = ({
  url = '',
  option = {},
  manual = true,
  method = 'GET',
  callbackFn = null,
  successParams = {},
}) => {
  const [options, setOptions] = useState({
    url,
    option,
    manual,
    method,
    callbackFn,
    successParams,
  })
  const [response, setResponse] = useState(null)
  const [loading, setLoading] = useState(!manual)
  const dispatch = useDispatch()

  function handleResponse(res, callbackFn, successParams) {
    if (res?.status === 403) {
      notification.error({
        message: `${res?.status}`,
        description: `${res?.data?.message}`,
      })
      return
    }
    if (res?.status === 401) {
      dispatch(setIsLogin(false))
      if (res?.config?.url?.includes('/rest-api/v1/authn')) {
        setResponse(res?.data)
      }
      return
    }
    if (res?.config?.url?.includes('/rest-api') && !res?.data?.success) {
      notification.error({
        message: `${res?.status}`,
        description: `${res?.data?.message}`,
      })
    } else {
      setResponse({
        ...res?.data,
        ...(callbackFn ? { callbackFn } : {}),
        ...(successParams ? { successParams } : {}),
      })
    }
  }

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        const res = await axios({
          url: url,
          method: method.toLowerCase(),
          ...options?.option,
        })
        handleResponse(res, callbackFn, successParams)
      } catch (err) {
        throw new Error(err)
      } finally {
        setLoading(false)
      }
    }

    if (!manual) {
      fetchData()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url, options, method, manual, dispatch])

  const refetch = async (newParams = {}) => {
    const newOptions = {
      ...options,
      ...newParams,
    }
    setOptions(newOptions)
    setLoading(true)
    try {
      const res = await axios({
        url: newOptions?.url,
        method: method.toLowerCase(),
        ...newOptions?.option,
      })
      handleResponse(res, newOptions?.callbackFn, newOptions?.successParams)
    } catch (err) {
      throw new Error(err)
    } finally {
      setLoading(false)
    }
  }

  return { response, loading, refetch }
}
