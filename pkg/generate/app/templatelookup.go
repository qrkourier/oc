package app

import (
	"fmt"

	"github.com/golang/glog"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/sets"

	"github.com/openshift/origin/pkg/client"
	templateapi "github.com/openshift/origin/pkg/template/api"
)

// TemplateSearcher resolves stored template arguments into template objects
type TemplateSearcher struct {
	Client                    client.TemplatesNamespacer
	TemplateConfigsNamespacer client.TemplateConfigsNamespacer
	Namespaces                []string
	StopOnMatch               bool
}

// Search searches for a template and returns matches with the object representation
func (r TemplateSearcher) Search(terms ...string) (ComponentMatches, error) {
	matches := ComponentMatches{}

	checkedNamespaces := sets.NewString()
	for _, namespace := range r.Namespaces {
		if checkedNamespaces.Has(namespace) {
			continue
		}
		checkedNamespaces.Insert(namespace)

		templates, err := r.Client.Templates(namespace).List(kapi.ListOptions{})
		if err != nil {
			if errors.IsNotFound(err) || errors.IsForbidden(err) {
				continue
			}
			return nil, err
		}

		exact := false
		for i := range templates.Items {
			template := &templates.Items[i]
			for _, term := range terms {
				glog.V(4).Infof("checking for term %s in namespace %s", term, namespace)
				if score, scored := templateScorer(*template, term); scored {
					if score == 0.0 {
						exact = true
					}
					glog.V(4).Infof("Adding template %q in project %q with score %f", template.Name, template.Namespace, score)
					matches = append(matches, &ComponentMatch{
						Value:       term,
						Argument:    fmt.Sprintf("--template=%q", template.Name),
						Name:        template.Name,
						Description: fmt.Sprintf("Template %q in project %q", template.Name, template.Namespace),
						Score:       score,
						Template:    template,
					})
					break
				}
			}
		}
		if exact {
			break
		}
	}

	return matches, nil
}

// IsPossibleTemplateFile returns true if the argument can be a template file
func IsPossibleTemplateFile(value string) bool {
	return isFile(value)
}

// TemplateFileSearcher resolves template files into template objects
type TemplateFileSearcher struct {
	Mapper       meta.RESTMapper
	Typer        runtime.ObjectTyper
	ClientMapper resource.ClientMapper
	Namespace    string
}

// Search attemps to read template files and transform it into template objects
func (r *TemplateFileSearcher) Search(terms ...string) (ComponentMatches, error) {
	matches := ComponentMatches{}

	for _, term := range terms {
		var isSingular bool
		obj, err := resource.NewBuilder(r.Mapper, r.Typer, r.ClientMapper).
			NamespaceParam(r.Namespace).RequireNamespace().
			FilenameParam(false, term).
			Do().
			IntoSingular(&isSingular).
			Object()

		if err != nil {
			glog.V(5).Infof("Error from template search: %v", err)
			continue
		}

		if !isSingular {
			return nil, fmt.Errorf("there is more than one object in %q", term)
		}

		template, ok := obj.(*templateapi.Template)
		if !ok {
			return nil, fmt.Errorf("object in %q is not a template", term)
		}

		matches = append(matches, &ComponentMatch{
			Value:       term,
			Argument:    fmt.Sprintf("--file=%q", template.Name),
			Name:        template.Name,
			Description: fmt.Sprintf("Template file %s", term),
			Score:       0,
			Template:    template,
		})
	}

	return matches, nil
}
