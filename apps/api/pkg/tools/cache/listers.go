package cache

import (
	"github.com/nrc-no/core/apps/api/pkg/apis/meta"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/labels"
	"github.com/sirupsen/logrus"
)

// AppendFunc is used to add a matching item to whatever list the caller is using
type AppendFunc func(interface{})

// ListAll calls appendFn with each value retrieved from store which matches the selector.
func ListAll(store Store, selector labels.Selector, appendFn AppendFunc) error {
	selectAll := selector.Empty()
	for _, m := range store.List() {
		if selectAll {
			// Avoid computing labels of the objects to speed up common flows
			// of listing all objects.
			appendFn(m)
			continue
		}
		metadata, err := meta.Accessor(m)
		if err != nil {
			return err
		}
		if selector.Matches(labels.Set(metadata.GetLabels())) {
			appendFn(m)
		}
	}
	return nil
}

// ListAllByNamespace used to list items belongs to namespace from Indexer.
func ListAllByNamespace(indexer Indexer, namespace string, selector labels.Selector, appendFn AppendFunc) error {
	selectAll := selector.Empty()
	if namespace == metav1.NamespaceAll {
		for _, m := range indexer.List() {
			if selectAll {
				// Avoid computing labels of the objects to speed up common flows
				// of listing all objects.
				appendFn(m)
				continue
			}
			metadata, err := meta.Accessor(m)
			if err != nil {
				return err
			}
			if selector.Matches(labels.Set(metadata.GetLabels())) {
				appendFn(m)
			}
		}
		return nil
	}

	items, err := indexer.Index(NamespaceIndex, &metav1.ObjectMeta{Namespace: namespace})
	if err != nil {
		// Ignore error; do slow search without index.
		logrus.Warningf("can not retrieve list of objects using index : %v", err)
		for _, m := range indexer.List() {
			metadata, err := meta.Accessor(m)
			if err != nil {
				return err
			}
			if metadata.GetNamespace() == namespace && selector.Matches(labels.Set(metadata.GetLabels())) {
				appendFn(m)
			}

		}
		return nil
	}
	for _, m := range items {
		if selectAll {
			// Avoid computing labels of the objects to speed up common flows
			// of listing all objects.
			appendFn(m)
			continue
		}
		metadata, err := meta.Accessor(m)
		if err != nil {
			return err
		}
		if selector.Matches(labels.Set(metadata.GetLabels())) {
			appendFn(m)
		}
	}

	return nil
}
