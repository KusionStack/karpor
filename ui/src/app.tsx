
import { ConfigProvider } from 'antd';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import zhCN from 'antd/locale/zh_CN';
import SandBox from "./views/sandbox/SandBox";
import './App.css';


dayjs.locale('zh-cn');

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <SandBox />
    </ConfigProvider>
  );
}

export default App;
