import React from 'react';
import Explorer from './Explorer';
import {BrowserRouter, Switch, Route} from 'react-router-dom';
import {Client} from './client';
import RecordsTable from './RecordsTable';

function App() {

  const client = new Client('https://localhost:9005');

  return (
    <BrowserRouter>
      <Switch>
        <Route path="/tables/:tableName" render={props => {
          const {tableName} = props.match.params;
          return <RecordsTable client={client} tableName={tableName}/>;
        }}/>
        <Route exact path="/">
          <Explorer client={client}/>
        </Route>
      </Switch>
    </BrowserRouter>
  );
}

export default App;
