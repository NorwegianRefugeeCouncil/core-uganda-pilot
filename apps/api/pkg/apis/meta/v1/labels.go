package v1

// Clones the given selector and returns a new selector with the given key and value added.
// Returns the given selector, if labelKey is empty.
func CloneSelectorAndAddLabel(selector *LabelSelector, labelKey, labelValue string) *LabelSelector {
  if labelKey == "" {
    // Don't need to add a label.
    return selector
  }

  // Clone.
  //newSelector := selector.DeepCopy()
  //
  //if newSelector.MatchLabels == nil {
  //  newSelector.MatchLabels = make(map[string]string)
  //}
  //
  //newSelector.MatchLabels[labelKey] = labelValue

  //return newSelector
  return nil
}

// AddLabelToSelector returns a selector with the given key and value added to the given selector's MatchLabels.
func AddLabelToSelector(selector *LabelSelector, labelKey, labelValue string) *LabelSelector {
  if labelKey == "" {
    // Don't need to add a label.
    return selector
  }
  if selector.MatchLabels == nil {
    selector.MatchLabels = make(map[string]string)
  }
  selector.MatchLabels[labelKey] = labelValue
  return selector
}

// SelectorHasLabel checks if the given selector contains the given label key in its MatchLabels
func SelectorHasLabel(selector *LabelSelector, labelKey string) bool {
  return len(selector.MatchLabels[labelKey]) > 0
}

