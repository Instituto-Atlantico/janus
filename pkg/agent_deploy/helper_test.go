package agent_deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSSHHostName(t *testing.T) {

	t.Run("should pass if hostname is well built", func(t *testing.T) {
		hostname := "host@hostname"

		valid := ValidateSSHHostName(hostname)

		assert.True(t, valid)
	})

	t.Run("should pass if hostname has an IP number", func(t *testing.T) {
		hostname := "host@192.168.0.3"

		valid := ValidateSSHHostName(hostname)

		assert.True(t, valid)
	})

	t.Run("should not pass if hostname is broken", func(t *testing.T) {
		hostname := "host @ hostname"

		valid := ValidateSSHHostName(hostname)

		assert.False(t, valid)
	})

	t.Run("should not pass if no username is passed", func(t *testing.T) {
		hostname := "@hostname"

		valid := ValidateSSHHostName(hostname)

		assert.False(t, valid)
	})

	t.Run("should not pass if no ip is passed", func(t *testing.T) {
		hostname := "user@"

		valid := ValidateSSHHostName(hostname)

		assert.False(t, valid)
	})

	t.Run("should not pass if no @ separator is passed", func(t *testing.T) {
		hostname := "userhostname"

		valid := ValidateSSHHostName(hostname)

		assert.False(t, valid)
	})
}
