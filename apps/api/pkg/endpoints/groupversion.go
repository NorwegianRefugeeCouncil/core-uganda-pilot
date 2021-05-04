package endpoints

import (
  "github.com/emicklei/go-restful"
  "github.com/nrc-no/core/apps/api/pkg/registry/rest"
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  utilerrors "k8s.io/apimachinery/pkg/util/errors"
  "path"
)

type APIGroupVersion struct {
  Root             string
  Serializer       runtime.NegotiatedSerializer
  GroupVersion     schema.GroupVersion
  Storage          map[string]rest.Storage
  Typer            runtime.ObjectTyper
  Creater          runtime.ObjectCreater
  ParameterCodec   runtime.ParameterCodec
  MetaGroupVersion *schema.GroupVersion
}

func (g *APIGroupVersion) InstallREST(container *restful.Container) error {

  prefix := path.Join(g.Root, g.GroupVersion.Group, g.GroupVersion.Version)
  installer := &APIInstaller{
    group:  g,
    prefix: prefix,
  }

  _, _, errors := installer.Install(container)

  return utilerrors.NewAggregate(errors)

}
