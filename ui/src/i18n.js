import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import en from './locales/en.json'
import zh from './locales/zh.json'

const resources = {
  en: {
    translation: en,
  },
  zh: {
    translation: zh,
  },
}

const currentLocale = localStorage.getItem('lang') || 'en'

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: currentLocale,
    lng: currentLocale,
    debug: true,
    interpolation: {
      escapeValue: false,
    },
  })

export default i18n
