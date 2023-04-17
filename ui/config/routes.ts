export const routes = [
  {
    path: '/',
    redirect: '/search',
  },
  {
    name: 'Search',
    path: '/search',
    icon: "SearchOutlined",
    component: './Home',
  },
  {
    name: 'Result',
    path: '/result',
    hideInMenu: true,
    component: './SearchList',
  },
  {
    name: 'Insight',
    path: '/insight',
    icon: "MacCommandOutlined",
    component: './SearchList',
  },
  {
    name: 'Cluster',
    path: '/cluster',
    icon: "ApartmentOutlined",
    component: './Cluster',
  },
];
