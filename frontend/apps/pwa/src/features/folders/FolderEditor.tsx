import React, { FC, useState } from 'react';
import { useForm } from 'react-hook-form';
import { Folder } from 'core-api-client';
import { Navigate } from 'react-router-dom';

import { databaseActions } from '../../reducers/database';
import {
  useDatabaseFromQueryParam,
  useFolderFromQueryParam,
} from '../../app/hooks';
import client from '../../app/client';

type FormData = {
  name: string;
};

export const FolderEditor: FC = () => {
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<FormData>();

  const parentFolder = useFolderFromQueryParam('parentId');
  const database = useDatabaseFromQueryParam('databaseId');
  const [folder, setFolder] = useState<Folder | undefined>(undefined);

  const onSubmit = (data: FormData) => {
    if (!database?.id) {
      return;
    }
    client.Folder.create({
      object: {
        name: data.name,
        databaseId: database?.id,
        parentId: parentFolder?.id,
      },
    }).then((resp) => {
      if (resp.success && resp.response) {
        databaseActions.addOne(resp.response);
        setFolder(resp.response);
      } else if (!resp.success && resp.error) {
        resp.error.details.causes.forEach((e: any) => {
          setError(e.field, { type: e.reason, message: e.message });
        });
      }
    });
  };

  if (!database) {
    return <>Database not found</>;
  }

  if (folder) {
    return (
      <Navigate to={`/browse/database/${database?.id}/folders/${folder.id}`} />
    );
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
                  className={`form-control ${errors.name ? 'is-invalid' : ''}`}
                  id="name"
                  aria-describedby="nameFeedback"
                />
                <div className="invalid-feedback is-invalid" id="nameFeedback">
                  {Object.values(errors).map((e) => (
                    <div key={e?.message}>{e?.message}</div>
                  ))}
                </div>
              </div>
              <button className="btn btn-primary">Create New Folder</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};
