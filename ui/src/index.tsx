import ReactDOM from "react-dom/client";
import { ConfigProvider } from "antd";
import dayjs from "dayjs";
import "dayjs/locale/zh-cn";
import zhCN from "antd/locale/zh_CN";
import { BrowserRouter } from "react-router-dom";
import Index from "./views/Index";
import "./utils/request";

import "./index.css";

dayjs.locale("zh-cn");

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Index />
      </BrowserRouter>
    </ConfigProvider>
  );
}

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);
root.render(<App />);
