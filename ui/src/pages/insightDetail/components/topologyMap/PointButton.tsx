import { Text } from "@antv/g6-react-node";

const TextCopy: any = Text;

const PointButton = (pointButtonProps: {
  width: number;
  color?: string;
  children?: string;
  onClick?: (evt) => void;
  onMouseOver?: (evt) => void;
  onMouseLeave?: (evt) => void;
  disabled?: boolean;
  marginRight?: number;
  marginLeft?: number;
}) => {
  const {
    width,
    color = "#000",
    children = "",
    onClick = () => {},
    onMouseOver = () => {},
    onMouseLeave = () => {},
    disabled = false,
  } = pointButtonProps;
  return (
    <TextCopy
      style={{
        width,
        fill: color,
        cursor: disabled ? "not-allowed" : "pointer",
        fontFamily: "PingFangSC",
        fontSize: 16,
      }}
      onClick={onClick}
      onMouseOver={onMouseOver}
      onMouseOut={onMouseLeave}
    >
      {children}
    </TextCopy>
  );
};

export default PointButton;
