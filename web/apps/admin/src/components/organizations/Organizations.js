"use strict";
exports.__esModule = true;
exports.Organizations = void 0;
var react_1 = require("react");
var react_table_1 = require("react-table");
var hooks_1 = require("../../hooks/hooks");
var react_router_dom_1 = require("react-router-dom");
var SectionTitle_1 = require("../sectiontitle/SectionTitle");
var Organizations = function (props) {
    var apiClient = (0, hooks_1.useApiClient)();
    var _a = (0, react_1.useState)([]), data = _a[0], setData = _a[1];
    var columns = (0, react_1.useMemo)(function () { return [
        {
            Header: "Name",
            accessor: "name",
            Cell: function (props, ctx) { return (<react_router_dom_1.Link to={"/organizations/" + props.row.original.id}>{props.value}</react_router_dom_1.Link>); }
        }
    ]; }, []);
    (0, react_1.useEffect)(function () {
        console.log(apiClient);
        apiClient.listOrganizations().then(function (resp) {
            if (resp.response) {
                setData(resp.response.items);
            }
        });
    }, [apiClient]);
    var table = (0, react_table_1.useTable)({ columns: columns, data: data });
    var getTableProps = table.getTableProps, getTableBodyProps = table.getTableBodyProps, headerGroups = table.headerGroups, rows = table.rows, prepareRow = table.prepareRow;
    return (<div className={"container mt-3"}>
            <div className={"row"}>
                <div className={"col"}>
                    <div className={"card card-darkula "}>
                        <div className={"card-body"}>
                            <SectionTitle_1.SectionTitle title={"Organizations"}>
                                <react_router_dom_1.Link className={"btn btn-darkula btn-sm"} to={"organizations/add"}>Add Organization</react_router_dom_1.Link>
                            </SectionTitle_1.SectionTitle>
                            <table className={"table table-darkula text-light"} {...getTableProps()}>
                                <thead>
                                {// Loop over the header rows
        headerGroups.map(function (headerGroup) { return (
        // Apply the header row props
        <tr {...headerGroup.getHeaderGroupProps()}>
                                            {// Loop over the headers in each row
            headerGroup.headers.map(function (column) { return (
            // Apply the header cell props
            <th {...column.getHeaderProps()}>
                                                        {// Render the header
                column.render('Header')}
                                                    </th>); })}
                                        </tr>); })}
                                </thead>
                                {/* Apply the table body props */}
                                <tbody {...getTableBodyProps()}>
                                {// Loop over the table rows
        rows.map(function (row) {
            // Prepare the row for display
            prepareRow(row);
            return (
            // Apply the row props
            <tr {...row.getRowProps()}>
                                                {// Loop over the rows cells
                row.cells.map(function (cell) {
                    // Apply the cell props
                    return (<td {...cell.getCellProps()}>
                                                                {// Render the cell contents
                        cell.render('Cell')}
                                                            </td>);
                })}
                                            </tr>);
        })}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>);
};
exports.Organizations = Organizations;
