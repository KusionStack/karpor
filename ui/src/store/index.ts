import { combineReducers, configureStore } from '@reduxjs/toolkit'
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux'
import logger from 'redux-logger'
import { persistReducer, persistStore } from 'redux-persist'
import storage from 'redux-persist/lib/storage' // 本地存储
// import storage from 'redux-persist/lib/storage/session' // 会话存储

// 多个Slice的引入；
import modules from './modules'

// 配置要存储的Slice；
const persistConfig = {
  key: 'root', // key是放入localStorage中的key
  storage,
}

const rootReducer = combineReducers(modules)
const newPersistReducer = persistReducer(persistConfig, rootReducer)
const store = configureStore({
  reducer: newPersistReducer,
  // 配置中间键
  middleware: getDefaultMiddleware =>
    // getDefaultMiddleware({ serializableCheck: false }).concat(), //不打印logger
    getDefaultMiddleware({ serializableCheck: false }).concat(logger),
  devTools: true,
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch

export const useAppDispatch = () => useDispatch<AppDispatch>()
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector

export const persistor = persistStore(store)
export default store
