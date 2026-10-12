import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import enTranslation from './locales/en.json'
import zhTranslation from './locales/zh.json'
import deTranslation from './locales/de.json'
import ptTranslation from './locales/pt.json'
import koTranslation from './locales/ko.json'
import jaTranslation from './locales/ja.json'
import frTranslation from './locales/fr.json'
import esTranslation from './locales/es.json'

const currentLocale = localStorage.getItem('lang') || 'en'

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en: {
        translation: enTranslation,
      },
      zh: {
        translation: zhTranslation,
      },
      de: {
        translation: deTranslation,
      },
      pt: {
        translation: ptTranslation,
      },
      ko: {
        translation: koTranslation,
      },
      ja: {
        translation: jaTranslation,
      },
      fr: {
        translation: frTranslation,
      },
      es: {
        translation: esTranslation,
      },
    },
    fallbackLng: currentLocale,
    lng: currentLocale,
    debug: true,
    interpolation: {
      escapeValue: false,
    },
  })

export default i18n
