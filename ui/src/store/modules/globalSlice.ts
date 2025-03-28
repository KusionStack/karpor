import { createSlice } from '@reduxjs/toolkit'

type InitialState = {
  isReadOnlyMode: any
  versionNumber: string
  isLogin: boolean
  githubBadge: boolean
  isUnsafeMode: any
  aiOptions: any
  isHighAvailability: boolean
}

const initialState: InitialState = {
  isReadOnlyMode: undefined,
  versionNumber: '',
  isLogin: false,
  githubBadge: false,
  isUnsafeMode: undefined,
  aiOptions: null,
  isHighAvailability: false,
}

export const globalSlice = createSlice({
  name: 'globalSlice',
  initialState,
  reducers: {
    setServerConfigMode: (state, action) => {
      state.isReadOnlyMode = action.payload
    },
    setVersionNumber: (state, action) => {
      state.versionNumber = action.payload
    },
    setIsLogin: (state, action) => {
      state.isLogin = action.payload
    },
    setGithubBadge: (state, action) => {
      state.githubBadge = action.payload
    },
    setIsUnsafeMode: (state, action) => {
      state.isUnsafeMode = action.payload
    },
    setAIOptions: (state, action) => {
      state.aiOptions = action.payload
    },
    setIsHighAvailability: (state, action) => {
      state.isHighAvailability = action.payload
    },
  },
})

export const {
  setServerConfigMode,
  setVersionNumber,
  setIsLogin,
  setGithubBadge,
  setIsUnsafeMode,
  setAIOptions,
  setIsHighAvailability,
} = globalSlice.actions

export default globalSlice.reducer
