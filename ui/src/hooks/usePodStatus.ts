import { useAxios } from '@/utils/request'
import { useState, useEffect, useRef, useCallback } from 'react'
import axios from 'axios'

export const usePodStatus = (
  cluster: string,
  namespace: string,
  name: string,
) => {
  const statusRef = useRef<string>('Loading')
  const [loading, setLoading] = useState(true)
  const initialized = useRef(false)

  const { response: summaryResponse, refetch: summaryRefetch } = useAxios({
    url: `${axios.defaults.baseURL}/rest-api/v1/insight/summary`,
    method: 'GET',
    manual: true,
    // A pod row can still be listed from the search index for a moment after
    // it's actually gone from the live cluster, so a not-found response here
    // is an expected, non-actionable state rather than a real error.
    silentError: true,
  })

  useEffect(() => {
    if (summaryResponse?.success) {
      statusRef.current = summaryResponse?.data?.resource?.status || 'unknown'
      setLoading(false)
    } else if (summaryResponse) {
      // Request completed but failed (e.g. pod not found live) -- resolve to
      // Unknown instead of spinning forever.
      statusRef.current = 'Unknown'
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
