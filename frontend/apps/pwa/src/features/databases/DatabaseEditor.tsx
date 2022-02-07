import React, { FC, useState } from 'react';
import { useForm } from 'react-hook-form';
import { Database } from 'core-api-client';
import { Redirect } from 'react-router-dom';
import _ from 'lodash';

import { databaseActions } from '../../reducers/database';
import client from '../../app/client';

type FormData = {
  name: string;
};

export const DatabaseEditor: FC = (props) => {
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<FormData>();
  const [database, setDatabase] = useState<Database | undefined>(undefined);

  const onSubmit = (data: FormData) => {
    client.Database.create({ object: { name: data.name } }).then((resp) => {
      if (resp.success && resp.response) {
        databaseActions.addOne(resp.response);
        setDatabase(resp.response);
      } else if (!resp.success && resp.error) {
        _.forEach(resp.error.details.causes, (e) => {
          setError(e.field, { type: e.reason, message: e.message });
        });
      }
    });
  };
  if (database) {
    return <Redirect to={`/browse/databases/${database.id}`} />;
  }
  return (
    <div className="flex-grow-1 bg-dark text-white pt-3">
      <div className="container">
        <div className="row">
          <div className="col">
            <h3>Create New Database</h3>
            <form onSubmit={handleSubmit(onSubmit)}>
              <div className="form-group mb-2">
                <label htmlFor="name">Database Name</label>
                <input
                  {...register('name')}
                  className={`form-control ${errors.name ? 'is-invalid' : ''}`}
                  id="name"
                  aria-describedby="nameFeedback"
                />
                <div className="invalid-feedback is-invalid" id="nameFeedback">
                  {_.map(errors, (e) => {
                    return <div key={e?.message}>{e?.message}</div>;
                  })}
                </div>
              </div>
              <button className="btn btn-primary">Create New Database</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};
