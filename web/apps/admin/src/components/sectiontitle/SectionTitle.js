"use strict";
exports.__esModule = true;
exports.SectionTitle = void 0;
var classnames_1 = require("classnames");
var SectionTitle = function (props) {
    var title = props.title, children = props.children;
    return (<div className={(0, classnames_1["default"])("border-bottom border-secondary pb-3 my-2 d-flex flex-row justify-content-center", props.className)}>
            <span className={"flex-grow-1 fs-5"}>{title}</span>
            {children}
        </div>);
};
exports.SectionTitle = SectionTitle;
