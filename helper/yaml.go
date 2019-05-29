package helper

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/zdnscloud/gok8s/client"
	"k8s.io/client-go/kubernetes/scheme"
)

const YamlDelimiter = "---\n"

func CreateResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, cli.Create)
}

func DeleteResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, func(ctx context.Context, obj runtime.Object) error {
		return cli.Delete(ctx, obj, nil)
	})
}

func UpdateResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, cli.Update)
}

func mapOnYamlDocument(yaml string, fn func(context.Context, runtime.Object) error) error {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	for _, doc := range strings.Split(yaml, YamlDelimiter) {
		obj, _, err := decode([]byte(doc), nil, nil)
		if err != nil {
			return err
		}
		if err := fn(context.TODO(), obj); err != nil {
			return err
		}
	}
	return nil
}
