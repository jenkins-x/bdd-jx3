package parsers_test

import (
	"testing"

	"github.com/jenkins-x/bdd-jx3/test/utils/parsers"
	"github.com/stretchr/testify/assert"
)

func TestGetApplicationsParser(t *testing.T) {
	out := `
WARNING: could not find the current user name user: Current not implemented on linux/amd64
APPLICATION           STAGING PODS URL
bdd-spring-1561456570 0.0.1   1/1  http://bdd-spring-1561456570.bdd-ghe-jx-pr-4153-100-staging.35.205.242.160.nip.io`
	applications, err := parsers.ParseJxGetApplications(out)
	assert.NoError(t, err)
	assert.Len(t, applications, 1)
}

func TestGetRemoteApplicationsParser(t *testing.T) {
	out := `APPLICATION           PRODUCTION PODS URL
bdd-spring-1617112975 0.0.1           http://bdd-spring-1617112975-myapps.34.123.71.97.nip.io`
	applications, err := parsers.ParseJxGetApplications(out)
	assert.NoError(t, err)
	assert.Len(t, applications, 1)

	for k, v := range applications {
		assert.Equal(t, "bdd-spring-1617112975", k, "should find first application")
		assert.Equal(t, "0.0.1", v.Version, "found app.Version")
		assert.Equal(t, "http://bdd-spring-1617112975-myapps.34.123.71.97.nip.io", v.Url, "found app.Url")
		assert.Equal(t, 0, v.RunningPods, "found app.RunningPods")
		assert.Equal(t, 0, v.DesiredPods, "found app.DesiredPods")
	}
}
