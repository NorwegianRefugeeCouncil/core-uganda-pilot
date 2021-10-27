import _ from "lodash";

const nestedObjectToQueryString = function (parentKey, obj) {
    let qs = _.reduce(obj, function (result, value, key) {
        if (!_.isNull(value) && !_.isUndefined(value)) {
            if (_.isArray(value)) {
                result += _.reduce(value, function (result1, value1) {
                    if (!_.isNull(value1) && !_.isUndefined(value1)) {
                        result1 += parentKey + '[' + key + ']' + '=' + value1 + '&';
                        return result1
                    } else {
                        return result1;
                    }
                }, '')
            } else if (_.isObject(value)) {
                result += nestedObjectToQueryString(key, value)
            } else {
                result += parentKey + '[' + key + ']' + '=' + value + '&';
            }
            return result;
        } else {
            return result
        }
    }, '').slice(0, -1);
    return qs;
};

export const objectToQueryString = function (obj) {
    let qs = _.reduce(obj, function (result, value, key) {
        if (!_.isNull(value) && !_.isUndefined(value)) {
            if (_.isArray(value)) {
                result += _.reduce(value, function (result1, value1) {
                    if (!_.isNull(value1) && !_.isUndefined(value1)) {
                        result1 += key + '=' + value1 + '&';
                        return result1
                    } else {
                        return result1;
                    }
                }, '')
            } else if (_.isObject(value)) {
                result += nestedObjectToQueryString(key, value) + '&';
            } else {
                result += key + '=' + value + '&';
            }
            return result;
        } else {
            return result
        }
    }, '');

    // if there is a trailing &, remove it
    if (qs[qs.length - 1] === '&') {
        qs = qs.slice(0, -1)
    }

    return qs;
};
