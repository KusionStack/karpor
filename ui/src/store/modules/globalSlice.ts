import { createSlice } from '@reduxjs/toolkit'

type InitialState = {
  isReadOnlyMode: any
}

const initialState: InitialState = {
  isReadOnlyMode: undefined,
}

export const globalSlice = createSlice({
  name: 'globalSlice',
  initialState,
  reducers: {
    setServerConfigMode: (state, action) => {
      state.isReadOnlyMode = action.payload
    },
  },
})

export const { setServerConfigMode } = globalSlice.actions

export default globalSlice.reducer
