import { createSlice } from '@reduxjs/toolkit'

type InitialState = {
  isReadOnlyMode: any
  versionNumber: string
}

const initialState: InitialState = {
  isReadOnlyMode: undefined,
  versionNumber: '',
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
  },
})

export const { setServerConfigMode, setVersionNumber } = globalSlice.actions

export default globalSlice.reducer
