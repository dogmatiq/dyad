package dyad_test

import (
	"reflect"
	"testing"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	type tag struct{}
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, reflect.TypeOf(tag{}).PkgPath())
}
