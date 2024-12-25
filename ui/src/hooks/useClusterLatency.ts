import { useAxios } from '@/utils/request'
import { useState, useEffect, useRef, useCallback } from 'react'
import axios from 'axios'

export const useClusterLatency = (cluster: string) => {
  const latencyRef = useRef<number>(0)
  const [loading, setLoading] = useState(true)
  const initialized = useRef(false)

  const { response: summaryResponse, refetch: summaryRefetch } = useAxios({
    url: `${axios.defaults.baseURL}/rest-api/v1/insight/summary`,
    method: 'GET',
    manual: true,
  })

  useEffect(() => {
    if (summaryResponse?.success) {
      latencyRef.current = summaryResponse?.data?.latency
      setLoading(false)
    } else {
      latencyRef.current = 0
      setLoading(true)
    }
  }, [summaryResponse])

  const fetchLatency = useCallback(() => {
    if (cluster) {
      summaryRefetch({
        option: {
          params: {
            cluster,
          },
        },
      })
    }
  }, [cluster, summaryRefetch])

  // Only execute once when component is first mounted
  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true
      fetchLatency()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return { latency: latencyRef, loading, refetch: fetchLatency }
}
