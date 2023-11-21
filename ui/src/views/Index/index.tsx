import { useRoutes } from "react-router-dom";
import router from "../../router";

const Index = () => {
  const element = useRoutes(router);
  return <>{element}</>;
};

export default Index;
