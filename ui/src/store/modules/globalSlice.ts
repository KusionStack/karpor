import { createSlice } from '@reduxjs/toolkit'

type InitialState = {
  isReadOnlyMode: any
}

//该store分库的初始值
const initialState: InitialState = {
  isReadOnlyMode: undefined,
}

export const globalSlice = createSlice({
  // store分库名称
  name: 'globalSlice',
  // store分库初始值
  initialState,
  reducers: {
    setServerConfigMode: (state, action) => {
      console.log(state, action, 'sadsadas')
      state.isReadOnlyMode = action.payload
    },
  },
})

export const { setServerConfigMode } = globalSlice.actions

export default globalSlice.reducer
