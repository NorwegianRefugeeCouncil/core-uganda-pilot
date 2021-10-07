type Route = {
  name: string
}

const routes: { [key: string]: Route; } = {
  home: {
    name: 'Home'
  },
  individuals: {
    name: 'Individuals List'
  },
  individual: {
    name: 'Edit Individual'
  },
  cases: {
    name: 'Cases'
  }
};

export default routes;