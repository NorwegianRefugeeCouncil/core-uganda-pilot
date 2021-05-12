package v1

import "fmt"

// LabelSelectorAsMap converts the LabelSelector api type into a map of strings, ie. the
// original structure of a label selector. Operators that cannot be converted into plain
// labels (Exists, DoesNotExist, NotIn, and In with more than one value) will result in
// an error.
func LabelSelectorAsMap(ps *LabelSelector) (map[string]string, error) {
  if ps == nil {
    return nil, nil
  }
  selector := map[string]string{}
  for k, v := range ps.MatchLabels {
    selector[k] = v
  }
  for _, expr := range ps.MatchExpressions {
    switch expr.Operator {
    case LabelSelectorOpIn:
      if len(expr.Values) != 1 {
        return selector, fmt.Errorf("operator %q without a single value cannot be converted into the old label selector format", expr.Operator)
      }
      // Should we do anything in case this will override a previous key-value pair?
      selector[expr.Key] = expr.Values[0]
    case LabelSelectorOpNotIn, LabelSelectorOpExists, LabelSelectorOpDoesNotExist:
      return selector, fmt.Errorf("operator %q cannot be converted into the old label selector format", expr.Operator)
    default:
      return selector, fmt.Errorf("%q is not a valid selector operator", expr.Operator)
    }
  }
  return selector, nil
}
