package exec

import (
	"github.com/captainhook-go/captainhook/configuration"
	"github.com/captainhook-go/captainhook/git"
	"github.com/captainhook-go/captainhook/hooks"
	"github.com/captainhook-go/captainhook/hooks/app"
	"github.com/captainhook-go/captainhook/hooks/conditions"
	"github.com/captainhook-go/captainhook/io"
	"strings"
)

// ConditionRunner executes conditions
// If the execution is successful the condition is considered true.
type ConditionRunner struct {
	cIO  io.IO
	conf *configuration.Configuration
	repo git.Repo
}

// Run executes the ConditionRunner
func (c *ConditionRunner) Run(hook string, condition *configuration.Condition) bool {
	if isLogicCondition(condition.Run()) {
		return checkLogicCondition(app.NewContext(c.cIO, c.conf, c.repo), condition, hook)
	}

	if len(condition.Conditions()) > 0 {
		c.cIO.Write("ConditionRunner: "+condition.Run()+" not allowed to have sub conditions", true, io.NORMAL)
		return true
	}

	conditionToExecute, err := c.crateCondition(condition)
	if err != nil {
		return false
	}
	if !conditionToExecute.IsApplicableFor(hook) {
		c.cIO.Write("ConditionRunner: "+condition.Run()+" not applicable for hook "+hook, true, io.VERBOSE)
		return true
	}
	return conditionToExecute.IsTrue(condition)
}

// createCondition creates the condition to execute
// It either creates an internally available condition or an `external` one that just executes a command
func (c *ConditionRunner) crateCondition(condition *configuration.Condition) (hooks.Condition, error) {
	if isInternalFunctionality(condition.Run()) {
		return c.createInternalCondition(condition)
	}
	return c.createExternalCondition()
}

// createInternalCondition creates one of CaptainHooks own conditions
func (c *ConditionRunner) createInternalCondition(condition *configuration.Condition) (hooks.Condition, error) {
	path := splitInternalPath(condition.Run())
	conditionGenerator, err := conditions.ConditionCreationFunc(path)
	if err != nil {
		c.cIO.Write("ConditionRunner: "+condition.Run()+"\n"+err.Error(), true, io.QUIET)
		return nil, err
	}
	return conditionGenerator(c.cIO, c.conf, c.repo), nil
}

// createExternalCondition creates a condition that runs the configured command
func (c *ConditionRunner) createExternalCondition() (hooks.Condition, error) {
	return conditions.NewExternalCommand(c.cIO, c.conf, c.repo), nil
}

func NewConditionRunner(cIO io.IO, conf *configuration.Configuration, repo git.Repo) *ConditionRunner {
	c := ConditionRunner{
		cIO,
		conf,
		repo,
	}
	return &c
}

func checkLogicCondition(context *app.Context, condition *configuration.Condition, hook string) bool {
	if isAndCondition(condition.Run()) {
		return DoAllConditionsApply(context, condition.Conditions(), hook)
	}
	return DoesAnyConditionApply(context, condition.Conditions(), hook)
}

func isAndCondition(run string) bool {
	return strings.Contains(strings.ToLower(run), "logic.and")
}

func DoAllConditionsApply(context *app.Context, conditions []*configuration.Condition, hook string) bool {
	conditionRunner := NewConditionRunner(context.IO(), context.Config(), context.Repository())
	for _, condition := range conditions {
		if !conditionRunner.Run(hook, condition) {
			return false
		}
	}
	return true
}

func DoesAnyConditionApply(context *app.Context, conditions []*configuration.Condition, hook string) bool {
	conditionRunner := NewConditionRunner(context.IO(), context.Config(), context.Repository())

	if len(conditions) < 1 {
		return true
	}

	for _, condition := range conditions {
		if conditionRunner.Run(hook, condition) {
			return true
		}
	}
	return false
}
