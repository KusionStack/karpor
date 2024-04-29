import { Text } from '@antv/g6-react-node'
import React from 'react'

const NodeLabel = (NodeLabelProps: {
  width: number
  color?: string
  children?: string
  onClick?: (evt) => void
  onMouseOver?: (evt) => void
  onMouseLeave?: (evt) => void
  disabled?: boolean
  marginRight?: number
  marginLeft?: number
}) => {
  const {
    width,
    color = '#000',
    children = '',
    onClick,
    onMouseOver,
    onMouseLeave,
    disabled = false,
  } = NodeLabelProps
  return (
    <Text
      style={{
        width,
        fill: color,
        cursor: disabled ? 'not-allowed' : 'pointer',
        fontFamily: 'PingFangSC',
        fontSize: '16px',
      }}
      onClick={onClick}
      onMouseOver={onMouseOver}
      onMouseOut={onMouseLeave}
    >
      {children}
    </Text>
  )
}

export default NodeLabel
