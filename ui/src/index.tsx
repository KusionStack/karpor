import React, { useEffect, useState } from 'react'
import ReactDOM from 'react-dom/client'
import { ConfigProvider } from 'antd'
import { Provider } from 'react-redux'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import { useTranslation } from 'react-i18next'
import zhCN from 'antd/locale/zh_CN'
import enUS from 'antd/locale/en_US'
import { BrowserRouter } from 'react-router-dom'
import WrappedRoutes from '@/router'
import store from '@/store'
import './i18n'

import '@/utils/request'

import './index.css'

dayjs.locale('zh-cn')

function App() {
  const { i18n } = useTranslation()
  const currentLocale = localStorage.getItem('lang')
  const [lang, setLang] = useState(currentLocale || 'en')

  useEffect(() => {
    setLang(i18n.language)
  }, [i18n.language])

  return (
    <Provider store={store}>
      <ConfigProvider
        locale={lang === 'en' ? enUS : zhCN}
        theme={{
          token: {
            colorPrimary: '#2F54EB',
          },
        }}
      >
        <BrowserRouter>
          <WrappedRoutes />
        </BrowserRouter>
      </ConfigProvider>
    </Provider>
  )
}

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement)
root.render(<App />)
