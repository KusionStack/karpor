import { useAxios } from '@/utils/request'
import { useState, useEffect, useRef, useCallback } from 'react'

export const usePodStatus = (
  cluster: string,
  namespace: string,
  name: string,
) => {
  const statusRef = useRef<string>('Loading')
  const [loading, setLoading] = useState(true)
  const initialized = useRef(false)

  const { response: summaryResponse, refetch: summaryRefetch } = useAxios({
    url: '/rest-api/v1/insight/summary',
    method: 'GET',
    manual: true,
  })

  useEffect(() => {
    if (summaryResponse?.success) {
      statusRef.current = summaryResponse?.data?.resource?.status || 'unknown'
      setLoading(false)
    } else {
      statusRef.current = 'Loading'
      setLoading(true)
    }
  }, [summaryResponse])

  const fetchStatus = useCallback(() => {
    if (cluster && namespace && name) {
      summaryRefetch({
        option: {
          params: {
            cluster,
            apiVersion: 'v1',
            namespace,
            kind: 'Pod',
            name,
          },
        },
      })
    }
  }, [cluster, namespace, name, summaryRefetch])

  // Only execute once when component is first mounted
  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true
      fetchStatus()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return { status: statusRef, loading, refetch: fetchStatus }
}
