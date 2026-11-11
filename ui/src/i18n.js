import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import en from './locales/en.json'
import zh from './locales/zh.json'
import de from './locales/de.json'
import pt from './locales/pt.json'
import ko from './locales/ko.json'
import hi from './locales/hi.json'
import ja from './locales/ja.json'

const resources = {
  en: {
    translation: en,
  },
  zh: {
    translation: zh,
  },
  de: {
    translation: de,
  },
  pt: {
    translation: pt,
  },
  ko: {
    translation: ko,
  },
  hi: {
    translation: hi,
  },
  ja: {
    translation: ja,
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
