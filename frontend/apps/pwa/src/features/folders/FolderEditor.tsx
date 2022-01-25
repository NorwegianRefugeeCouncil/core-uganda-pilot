import React, { FC, Fragment, useState } from 'react';
import { useForm } from 'react-hook-form';
import { Folder } from 'core-api-client';
import { Redirect } from 'react-router-dom';

import { databaseActions } from '../../reducers/database';
import {
  useDatabaseFromQueryParam,
  useFolderFromQueryParam,
} from '../../app/hooks';
import client from '../../app/client';
import { ApiErrorDetails } from '../../types/errors';

type FormData = {
  name: string;
};

export const FolderEditor: FC = (props) => {
  const { register, handleSubmit } = useForm<FormData>();

  const parentFolder = useFolderFromQueryParam('parentId');
  const database = useDatabaseFromQueryParam('databaseId');
  const [folder, setFolder] = useState<Folder | undefined>(undefined);
  const [errors, setErrors] = useState<ApiErrorDetails[] | undefined>();

  const onSubmit = (data: FormData) => {
    if (!database?.id) {
      return;
    }
    client
      .createFolder({
        object: {
          name: data.name,
          databaseId: database?.id,
          parentId: parentFolder?.id,
        },
      })
      .then((resp) => {
        if (resp.success && resp.response) {
          databaseActions.addOne(resp.response);
          setFolder(resp.response);
        } else if (!resp.success && resp.error) {
          setErrors(resp.error.details.causes);
        }
      });
  };

  if (!database) {
    return <>Database not found</>;
  }

  if (folder) {
    return <Redirect to={`/browse/folders/${folder.id}`} />;
  }

  return (
    <div className="flex-grow-1 bg-dark text-white pt-3">
      <div className="container">
        <div className="row">
          <div className="col">
            <h3>Create New Folder</h3>
            <form onSubmit={handleSubmit(onSubmit)}>
              <div className="form-group mb-2">
                <label htmlFor="name">Folder Name</label>
                <input
                  {...register('name')}
                  className={`form-control ${errors ? 'is-invalid' : ''}`}
                  id="name"
                  aria-describedby="nameFeedback"
                />
                {errors?.map((error) => {
                  return (
                    <div
                      className="invalid-feedback is-invalid"
                      id="nameFeedback"
                      key={`${error.field}_${error.reason}`}
                    >
                      {error.message}
                    </div>
                  );
                })}
              </div>
              <button className="btn btn-primary">Create New Folder</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};
