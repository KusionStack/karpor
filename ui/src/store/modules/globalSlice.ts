import { createSlice } from '@reduxjs/toolkit'

type InitialState = {
  isReadOnlyMode: any
  versionNumber: string
  githubBadge: boolean
}

const initialState: InitialState = {
  isReadOnlyMode: undefined,
  versionNumber: '',
  githubBadge: false,
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
    setGithubBadge: (state, action) => {
      state.githubBadge = action.payload
    },
  },
})

export const { setServerConfigMode, setVersionNumber, setGithubBadge } = globalSlice.actions

export default globalSlice.reducer
