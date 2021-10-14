type Route = {
  title: string,
  name: string
}

const routes: { [key: string]: Route; } = {
  home: {
    title: 'Home',
    name: 'home',
  },
  individuals: {
    title: 'Individuals List',
    name: 'individuals'
  },
  individual: {
    title: 'Edit Individual',
    name: 'individual',
  },
  cases: {
    title: 'Cases',
    name: 'cases',
  }
};

export default routes;