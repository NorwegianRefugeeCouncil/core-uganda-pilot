type Route = {
  title: string,
  name: string
}

const routes: { [key: string]: Route; } = {
  designSystem: {
    title: 'Design System',
    name: 'designSystem',
  },
  forms: {
    title: 'Forms',
    name: 'forms',
  },
  addRecord: {
    title: 'Add Record',
    name: 'addRecord'
  },
  viewRecord: {
    title: 'View Record',
    name: 'viewRecord'
  },
  records: {
    title: 'Records',
    name: 'records'
  }
};

export default routes;
