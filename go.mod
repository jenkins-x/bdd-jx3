module github.com/jenkins-x/bdd-jx

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/fatih/color v1.9.0
	github.com/jenkins-x/go-scm v1.5.216
	github.com/jenkins-x/jx-api/v4 v4.0.23
	github.com/jenkins-x/jx-helpers/v3 v3.0.73
	github.com/jenkins-x/jx-kube-client/v3 v3.0.2
	github.com/jenkins-x/jx-logging/v3 v3.0.3
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible

)

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)

go 1.15
