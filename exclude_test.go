package hclog

import (
	"regexp"
	"testing"

	hlog "github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestExclude(t *testing.T) {
	t.Run("excludes by message", func(t *testing.T) {
		var em ExcludeByMessage
		em.Add("foo")
		em.Add("bar")

		assert.True(t, em.Exclude(hlog.Info, "foo"))
		assert.True(t, em.Exclude(hlog.Info, "bar"))
		assert.False(t, em.Exclude(hlog.Info, "qux"))
		assert.False(t, em.Exclude(hlog.Info, "foo qux"))
		assert.False(t, em.Exclude(hlog.Info, "qux bar"))
	})

	t.Run("excludes by prefix", func(t *testing.T) {
		ebp := ExcludeByPrefix("foo: ")

		assert.True(t, ebp.Exclude(hlog.Info, "foo: rocks"))
		assert.False(t, ebp.Exclude(hlog.Info, "foo"))
		assert.False(t, ebp.Exclude(hlog.Info, "qux foo: bar"))
	})

	t.Run("exclude by regexp", func(t *testing.T) {
		ebr := &ExcludeByRegexp{
			Regexp: regexp.MustCompile("(foo|bar)"),
		}

		assert.True(t, ebr.Exclude(hlog.Info, "foo"))
		assert.True(t, ebr.Exclude(hlog.Info, "bar"))
		assert.True(t, ebr.Exclude(hlog.Info, "foo qux"))
		assert.True(t, ebr.Exclude(hlog.Info, "qux bar"))
		assert.False(t, ebr.Exclude(hlog.Info, "qux"))
	})

	t.Run("excludes many funcs", func(t *testing.T) {
		ef := ExcludeFuncs{
			ExcludeByPrefix("foo: ").Exclude,
			ExcludeByPrefix("bar: ").Exclude,
		}

		assert.True(t, ef.Exclude(hlog.Info, "foo: rocks"))
		assert.True(t, ef.Exclude(hlog.Info, "bar: rocks"))
		assert.False(t, ef.Exclude(hlog.Info, "foo"))
		assert.False(t, ef.Exclude(hlog.Info, "qux foo: bar"))

	})
}
