import * as React from 'react';

import client from './app/client';

const uuid = '44bb1ed5-687d-4c8a-b7cc-401c0e6c1f6e';

export const Foo = () => {
  const [sendId, setSendId] = React.useState(false);
  const [sendName, setSendName] = React.useState(true);
  const [sendOtherField, setSendOtherField] = React.useState(true);
  const [sendValid, setSendValid] = React.useState(true);

  const [id, setId] = React.useState(uuid);
  const [name, setName] = React.useState('A name');
  const [otherField, setOtherField] = React.useState(10);
  const [valid, setValid] = React.useState(true);

  const send = () => {
    const body: any = {};
    if (sendId) body.id = id;
    if (sendName) body.name = name;
    if (sendOtherField) body.otherField = otherField;
    if (sendValid) body.valid = valid;
    const resposne = client.createFoo(body);
    console.log(resposne);
  };

  return (
    <div>
      <h3>Send:</h3>
      <input id="sendId" type="checkbox" checked={sendId} onChange={() => setSendId(!sendId)} />
      <label htmlFor="sendId">Id</label>
      <br />
      <input id="sendName" type="checkbox" checked={sendName} onChange={() => setSendName(!sendName)} />
      <label htmlFor="sendName">sendName</label>
      <br />
      <input id="sendOtherField" type="checkbox" checked={sendOtherField} onChange={() => setSendOtherField(!sendOtherField)} />
      <label htmlFor="sendOtherField">sendOtherField</label>
      <br />
      <input id="sendValid" type="checkbox" checked={sendValid} onChange={() => setSendValid(!sendValid)} />
      <label htmlFor="sendValid">sendValid</label>
      <br />

      <h3>Value:</h3>
      <input id="id" type="text" value={id} onChange={(e) => setId(e.target.value)} />
      <label htmlFor="id">id</label>
      <br />
      <input id="name" type="text" value={name} onChange={(e) => setName(e.target.value)} />
      <label htmlFor="name">name</label>
      <br />
      <input id="otherField" type="number" value={otherField} onChange={(e) => setOtherField(Number(e.target.value))} />
      <label htmlFor="otherField">otherField</label>
      <br />
      <input id="valid" type="checkbox" checked={valid} onChange={() => setValid(!valid)} />
      <label htmlFor="sendId">valid</label>
      <br />

      <button type="button" onClick={send}>
        Send
      </button>
    </div>
  );
};
