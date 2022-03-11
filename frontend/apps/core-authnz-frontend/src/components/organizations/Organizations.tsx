import React, { FC, useEffect, useMemo, useState } from 'react';
import { Column, useTable } from 'react-table';
import { Link } from 'react-router-dom';

import { Organization } from '../../types/types';
import { useApiClient } from '../../hooks/hooks';
import { SectionTitle } from '../sectiontitle/SectionTitle';

export const Organizations: FC = () => {
  const apiClient = useApiClient();

  const [data, setData] = useState<Organization[]>([]);
  const columns = useMemo<Column<Organization>[]>(
    () => [
      {
        Header: 'Name',
        accessor: 'name',
        Cell: (props) => (
          <Link to={`/organizations/${props.row.original.id}`}>
            {props.value}
          </Link>
        ),
      },
    ],
    [],
  );

  useEffect(() => {
    apiClient.listOrganizations().then((resp) => {
      if (resp.response) {
        setData(resp.response.items);
      }
    });
  }, [apiClient]);

  const table = useTable({ columns, data });

  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } =
    table;

  return (
    <div className="container mt-3">
      <div className="row">
        <div className="col">
          <div className="card card-darkula ">
            <div className="card-body">
              <SectionTitle title="Organizations">
                <Link className="btn btn-darkula btn-sm" to="/organizations/add">
                  Add Organization
                </Link>
              </SectionTitle>
              <table
                className="table table-darkula text-light"
                {...getTableProps()}
              >
                <thead>
                  {
                    // Loop over the header rows
                    headerGroups.map((headerGroup) => (
                      // Apply the header row props
                      <tr {...headerGroup.getHeaderGroupProps()}>
                        {
                          // Loop over the headers in each row
                          headerGroup.headers.map((column) => (
                            // Apply the header cell props
                            <th {...column.getHeaderProps()}>
                              {
                                // Render the header
                                column.render('Header')
                              }
                            </th>
                          ))
                        }
                      </tr>
                    ))
                  }
                </thead>
                {/* Apply the table body props */}
                <tbody {...getTableBodyProps()}>
                  {
                    // Loop over the table rows
                    rows.map((row) => {
                      // Prepare the row for display
                      prepareRow(row);
                      return (
                        // Apply the row props
                        <tr {...row.getRowProps()}>
                          {
                            // Loop over the rows cells
                            row.cells.map((cell) => {
                              // Apply the cell props
                              return (
                                <td {...cell.getCellProps()}>
                                  {
                                    // Render the cell contents
                                    cell.render('Cell')
                                  }
                                </td>
                              );
                            })
                          }
                        </tr>
                      );
                    })
                  }
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
