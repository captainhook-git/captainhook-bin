package inconfig

import (
	"github.com/captainhook-go/captainhook/configuration"
	"github.com/captainhook-go/captainhook/test"
	"testing"
)

func TestCustomValueIsFalsy(t *testing.T) {
	inOut := test.CreateFakeIO()
	conf, _ := configuration.NewJsonFactory().CreateConfig("../../../test/files/config/valid_with_custom.json", &configuration.NullableAppSettings{})
	repo := test.CreateFakeRepo()

	options := configuration.NewOptions(map[string]interface{}{"value": "fiz"})
	condition := configuration.NewCondition("CaptainHook::InConfig.CustomValueIsFalsy", options, []*configuration.Condition{})

	action := NewCustomValueIsFalsy(inOut, conf, repo)
	if !action.IsTrue(condition) {
		t.Errorf("Custom value should be falsy")
	}
}

func TestCustomValueIsNotFalsy(t *testing.T) {
	inOut := test.CreateFakeIO()
	conf, _ := configuration.NewJsonFactory().CreateConfig("../../../test/files/config/valid_with_custom.json", &configuration.NullableAppSettings{})
	repo := test.CreateFakeRepo()

	options := configuration.NewOptions(map[string]interface{}{"value": "foo"})
	condition := configuration.NewCondition("CaptainHook::InConfig.CustomValueIsFalsy", options, []*configuration.Condition{})

	action := NewCustomValueIsFalsy(inOut, conf, repo)
	if action.IsTrue(condition) {
		t.Errorf("Custom value should not be falsy")
	}
}

func TestApplicableFor(t *testing.T) {
	inOut := test.CreateFakeIO()
	conf := test.CreateFakeConfig()
	repo := test.CreateFakeRepo()

	action := NewCustomValueIsFalsy(inOut, conf, repo)
	if !action.IsApplicableFor("pre-commit") {
		t.Errorf("Condition should be applicable for 'pre-commit'")
	}
	if !action.IsApplicableFor("pre-push") {
		t.Errorf("Condition should be applicable for 'pre-push'")
	}
}
