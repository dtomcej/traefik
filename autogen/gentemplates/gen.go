// Code generated by go-bindata.
// sources:
// templates/consul_catalog.tmpl
// templates/docker.tmpl
// templates/ecs.tmpl
// templates/eureka.tmpl
// templates/kubernetes.tmpl
// templates/kv.tmpl
// templates/marathon.tmpl
// templates/mesos.tmpl
// templates/notFound.tmpl
// templates/rancher.tmpl
// DO NOT EDIT!

package gentemplates

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesConsul_catalogTmpl = []byte(`[backends]
{{range $index, $node := .Nodes}}
  [backends."backend-{{getBackend $node}}".servers."{{getBackendName $node $index}}"]
    url = "{{getAttribute "protocol" $node.Service.Tags "http"}}://{{getBackendAddress $node}}:{{$node.Service.Port}}"
    {{$weight := getAttribute "backend.weight" $node.Service.Tags "0"}}
    {{with $weight}}
      weight = {{$weight}}
    {{end}}
{{end}}

{{range .Services}}
  {{$service := .ServiceName}}
  {{$circuitBreaker := getAttribute "backend.circuitbreaker" .Attributes ""}}
  {{with $circuitBreaker}}
  [backends."backend-{{$service}}".circuitbreaker]
    expression = "{{$circuitBreaker}}"
  {{end}}

  [backends."backend-{{$service}}".loadbalancer]
    method = "{{getAttribute "backend.loadbalancer" .Attributes "wrr"}}"
    sticky = {{getSticky .Attributes}}
    {{if hasStickinessLabel .Attributes}}
    [backends."backend-{{$service}}".loadbalancer.stickiness]
      cookieName = "{{getStickinessCookieName .Attributes}}"
    {{end}}

  {{if hasMaxconnAttributes .Attributes}}
  [backends."backend-{{$service}}".maxconn]
    amount = {{getAttribute "backend.maxconn.amount" .Attributes "" }}
    extractorfunc = "{{getAttribute "backend.maxconn.extractorfunc" .Attributes "" }}"
  {{end}}

{{end}}

[frontends]
{{range .Services}}
  [frontends."frontend-{{.ServiceName}}"]
  backend = "backend-{{.ServiceName}}"
  passHostHeader = {{getAttribute "frontend.passHostHeader" .Attributes "true"}}
  priority = {{getAttribute "frontend.priority" .Attributes "0"}}
  {{$entryPoints := getAttribute "frontend.entrypoints" .Attributes ""}}
  {{with $entryPoints}}
    entrypoints = [{{range getEntryPoints $entryPoints}}
      "{{.}}",
    {{end}}]
  {{end}}
  basicAuth = [{{range getBasicAuth .Attributes}}
  "{{.}}",
  {{end}}]
  [frontends."frontend-{{.ServiceName}}".routes."route-host-{{.ServiceName}}"]
    rule = "{{getFrontendRule .}}"
{{end}}
`)

func templatesConsul_catalogTmplBytes() ([]byte, error) {
	return _templatesConsul_catalogTmpl, nil
}

func templatesConsul_catalogTmpl() (*asset, error) {
	bytes, err := templatesConsul_catalogTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/consul_catalog.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesDockerTmpl = []byte(`{{$backendServers := .Servers}}
[backends]{{range $backendName, $backend := .Backends}}
    {{if hasCircuitBreakerLabel $backend}}
    [backends.backend-{{$backendName}}.circuitbreaker]
      expression = "{{getCircuitBreakerExpression $backend}}"
    {{end}}

    {{if hasLoadBalancerLabel $backend}}
    [backends.backend-{{$backendName}}.loadbalancer]
      method = "{{getLoadBalancerMethod $backend}}"
      sticky = {{getSticky $backend}}
      {{if hasStickinessLabel $backend}}
      [backends.backend-{{$backendName}}.loadBalancer.stickiness]
        cookieName = "{{getStickinessCookieName $backend}}"
      {{end}}
    {{end}}

    {{if hasMaxConnLabels $backend}}
    [backends.backend-{{$backendName}}.maxconn]
      amount = {{getMaxConnAmount $backend}}
      extractorfunc = "{{getMaxConnExtractorFunc $backend}}"
    {{end}}

    {{$servers := index $backendServers $backendName}}
    {{range $serverName, $server := $servers}}
    {{if hasServices $server}}
      {{$services := getServiceNames $server}}
      {{range $serviceIndex, $serviceName := $services}}
      [backends.backend-{{getServiceBackend $server $serviceName}}.servers.service-{{$serverName}}]
      url = "{{getServiceProtocol $server $serviceName}}://{{getIPAddress $server}}:{{getServicePort $server $serviceName}}"
      weight = {{getServiceWeight $server $serviceName}}
      {{end}}
    {{else}}
      [backends.backend-{{$backendName}}.servers.server-{{$server.Name | replace "/" "" | replace "." "-"}}]
      url = "{{getProtocol $server}}://{{getIPAddress $server}}:{{getPort $server}}"
      weight = {{getWeight $server}}
    {{end}}
    {{end}}

{{end}}

[frontends]{{range $frontend, $containers := .Frontends}}
  {{$container := index $containers 0}}
  {{if hasServices $container}}
  {{$services := getServiceNames $container}}
  {{range $serviceIndex, $serviceName := $services}}
  [frontends."frontend-{{getServiceBackend $container $serviceName}}"]
  backend = "backend-{{getServiceBackend $container $serviceName}}"
  passHostHeader = {{getServicePassHostHeader $container $serviceName}}
  redirect =  "{{getServiceRedirect $container $serviceName}}"
  {{if getWhitelistSourceRange $container}}
    whitelistSourceRange = [{{range getWhitelistSourceRange $container}}
      "{{.}}",
    {{end}}]
  {{end}}
  priority = {{getServicePriority $container $serviceName}}
  entryPoints = [{{range getServiceEntryPoints $container $serviceName}}
    "{{.}}",
  {{end}}]
  basicAuth = [{{range getServiceBasicAuth $container $serviceName}}
    "{{.}}",
  {{end}}]
    [frontends."frontend-{{getServiceBackend $container $serviceName}}".routes."service-{{$serviceName | replace "/" "" | replace "." "-"}}"]
    rule = "{{getServiceFrontendRule $container $serviceName}}"
  {{end}}
  {{else}}
  [frontends."frontend-{{$frontend}}"]
  backend = "backend-{{getBackend $container}}"
  passHostHeader = {{getPassHostHeader $container}}
  redirect = "{{getRedirect $container}}"
  {{if getWhitelistSourceRange $container}}
    whitelistSourceRange = [{{range getWhitelistSourceRange $container}}
      "{{.}}",
    {{end}}]
  {{end}}
  priority = {{getPriority $container}}
  entryPoints = [{{range getEntryPoints $container}}
    "{{.}}",
  {{end}}]
  basicAuth = [{{range getBasicAuth $container}}
    "{{.}}",
  {{end}}]
  [frontends."frontend-{{$frontend}}".headers]
  {{if hasSSLRedirectHeaders $container}}
  SSLRedirect = {{getSSLRedirectHeaders $container}}
  {{end}}
  {{if hasSSLTemporaryRedirectHeaders $container}}
  SSLTemporaryRedirect = {{getSSLTemporaryRedirectHeaders $container}}
  {{end}}
  {{if hasSSLHostHeaders $container}}
  SSLHost = {{getSSLHostHeaders $container}}
  {{end}}
  {{if hasSTSSecondsHeaders $container}}
  STSSeconds = {{getSTSSecondsHeaders $container}}
  {{end}}
  {{if hasSTSIncludeSubdomainsHeaders $container}}
  STSIncludeSubdomains = {{getSTSIncludeSubdomainsHeaders $container}}
  {{end}}
  {{if hasSTSPreloadHeaders $container}}
  STSPreload = {{getSTSPreloadHeaders $container}}
  {{end}}
  {{if hasForceSTSHeaderHeaders $container}}
  ForceSTSHeader = {{getForceSTSHeaderHeaders $container}}
  {{end}}
  {{if hasFrameDenyHeaders $container}}
  FrameDeny = {{getFrameDenyHeaders $container}}
  {{end}}
  {{if hasCustomFrameOptionsValueHeaders $container}}
  CustomFrameOptionsValue = {{getCustomFrameOptionsValueHeaders $container}}
  {{end}}
  {{if hasContentTypeNosniffHeaders $container}}
  ContentTypeNosniff = {{getContentTypeNosniffHeaders $container}}
  {{end}}
  {{if hasBrowserXSSFilterHeaders $container}}
  BrowserXSSFilter = {{getBrowserXSSFilterHeaders $container}}
  {{end}}
  {{if hasContentSecurityPolicyHeaders $container}}
  ContentSecurityPolicy = {{getContentSecurityPolicyHeaders $container}}
  {{end}}
  {{if hasPublicKeyHeaders $container}}
  PublicKey = {{getPublicKeyHeaders $container}}
  {{end}}
  {{if hasReferrerPolicyHeaders $container}}
  ReferrerPolicy = {{getReferrerPolicyHeaders $container}}
  {{end}}
  {{if hasIsDevelopmentHeaders $container}}
  IsDevelopment = {{getIsDevelopmentHeaders $container}}
  {{end}}
  {{if hasRequestHeaders $container}}
    [frontends."frontend-{{$frontend}}".headers.customrequestheaders]
    {{range $k, $v := getRequestHeaders $container}}
    {{$k}} = "{{$v}}"
    {{end}}
  {{end}}
  {{if hasResponseHeaders $container}}
    [frontends."frontend-{{$frontend}}".headers.customresponseheaders]
    {{range $k, $v := getResponseHeaders $container}}
    {{$k}} = "{{$v}}"
    {{end}}
  {{end}}
  {{if hasAllowedHostsHeaders $container}}
    [frontends."frontend-{{$frontend}}".headers.AllowedHosts]
    {{range getAllowedHostsHeaders $container}}
    "{{.}}"
    {{end}}
  {{end}}
  {{if hasHostsProxyHeaders $container}}
    [frontends."frontend-{{$frontend}}".headers.HostsProxyHeaders]
    {{range getHostsProxyHeaders $container}}
    "{{.}}"
    {{end}}
  {{end}}
  {{if hasSSLProxyHeaders $container}}
    [frontends."frontend-{{$frontend}}".headers.SSLProxyHeaders]
    {{range $k, $v := getSSLProxyHeaders $container}}
    {{$k}} = "{{$v}}"
    {{end}}
  {{end}}
    [frontends."frontend-{{$frontend}}".routes."route-frontend-{{$frontend}}"]
    rule = "{{getFrontendRule $container}}"
  {{end}}
{{end}}
`)

func templatesDockerTmplBytes() ([]byte, error) {
	return _templatesDockerTmpl, nil
}

func templatesDockerTmpl() (*asset, error) {
	bytes, err := templatesDockerTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/docker.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesEcsTmpl = []byte(`[backends]{{range $serviceName, $instances := .Services}}
  [backends.backend-{{ $serviceName }}.loadbalancer]
    method = "{{ getLoadBalancerMethod $instances}}"
    sticky = {{ getLoadBalancerSticky $instances}}
    {{if hasStickinessLabel $instances}} 
    [backends.backend-{{ $serviceName }}.loadbalancer.stickiness]
      cookieName = "{{getStickinessCookieName $instances}}"
    {{end}}
    {{ if hasHealthCheckLabels $instances }}
    [backends.backend-{{ $serviceName }}.healthcheck]
      path = "{{getHealthCheckPath $instances }}"
      interval = "{{getHealthCheckInterval $instances }}"
    {{end}}

  {{range $index, $i := $instances}}
    [backends.backend-{{ $i.Name }}.servers.server-{{ $i.Name }}{{ $i.ID }}]
      url = "{{ getProtocol $i }}://{{ getHost $i }}:{{ getPort $i }}"
      weight = {{ getWeight $i}}
  {{end}}
{{end}}

[frontends]{{range $serviceName, $instances := .Services}}
  {{range filterFrontends $instances}}
    [frontends.frontend-{{ $serviceName }}]
      backend = "backend-{{ $serviceName }}"
      passHostHeader = {{ getPassHostHeader .}}
      priority = {{ getPriority .}}
      entryPoints = [{{range  getEntryPoints .}}
      "{{.}}",
    {{end}}]
      basicAuth = [{{range getBasicAuth .}}
      "{{.}}",
    {{end}}]
    [frontends.frontend-{{ $serviceName }}.routes.route-frontend-{{ $serviceName }}]
      rule = "{{getFrontendRule .}}"
  {{end}}
{{end}}`)

func templatesEcsTmplBytes() ([]byte, error) {
	return _templatesEcsTmpl, nil
}

func templatesEcsTmpl() (*asset, error) {
	bytes, err := templatesEcsTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/ecs.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesEurekaTmpl = []byte(`[backends]{{range .Applications}}
    {{ $app := .}}
    {{range .Instances}}
    [backends.backend{{$app.Name}}.servers.server-{{ getInstanceID . }}]
    url = "{{ getProtocol . }}://{{ .IpAddr }}:{{ getPort . }}"
    weight = {{ getWeight . }}
{{end}}{{end}}

[frontends]{{range .Applications}}
  [frontends.frontend{{.Name}}]
    backend = "backend{{.Name}}"
    entryPoints = ["http"]
    [frontends.frontend{{.Name }}.routes.route-host{{.Name}}]
      rule = "Host:{{ .Name | tolower }}"
{{end}}
`)

func templatesEurekaTmplBytes() ([]byte, error) {
	return _templatesEurekaTmpl, nil
}

func templatesEurekaTmpl() (*asset, error) {
	bytes, err := templatesEurekaTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/eureka.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesKubernetesTmpl = []byte(`[backends]{{range $backendName, $backend := .Backends}}
    [backends."{{$backendName}}"]
    {{if $backend.CircuitBreaker}}
    [backends."{{$backendName}}".circuitbreaker]
      expression = "{{$backend.CircuitBreaker.Expression}}"
    {{end}}
    [backends."{{$backendName}}".loadbalancer]
      method = "{{$backend.LoadBalancer.Method}}"
      {{if $backend.LoadBalancer.Sticky}}
      sticky = true
      {{end}}
      {{if $backend.LoadBalancer.Stickiness}}
      [backends."{{$backendName}}".loadbalancer.stickiness]
        cookieName = "{{$backend.LoadBalancer.Stickiness.CookieName}}"
      {{end}}
    {{range $serverName, $server := $backend.Servers}}
    [backends."{{$backendName}}".servers."{{$serverName}}"]
    url = "{{$server.URL}}"
    weight = {{$server.Weight}}
    {{end}}
{{end}}

[frontends]{{range $frontendName, $frontend := .Frontends}}
  [frontends."{{$frontendName}}"]
  backend = "{{$frontend.Backend}}"
  priority = {{$frontend.Priority}}
  passHostHeader = {{$frontend.PassHostHeader}}
  redirect =  "{{$frontend.Redirect}}"
  basicAuth = [{{range $frontend.BasicAuth}}
      "{{.}}",
  {{end}}]
  whitelistSourceRange = [{{range $frontend.WhitelistSourceRange}}
    "{{.}}",
  {{end}}]
  [frontends."{{$frontendName}}".headers]
  SSLRedirect = {{$frontend.Headers.SSLRedirect}}
  SSLTemporaryRedirect = {{$frontend.Headers.SSLTemporaryRedirect}}
  SSLHost = "{{$frontend.Headers.SSLHost}}"
  STSSeconds = {{$frontend.Headers.STSSeconds}}
  STSIncludeSubdomains = {{$frontend.Headers.STSIncludeSubdomains}}
  STSPreload = {{$frontend.Headers.STSPreload}}
  ForceSTSHeader = {{$frontend.Headers.ForceSTSHeader}}
  FrameDeny = {{$frontend.Headers.FrameDeny}}
  CustomFrameOptionsValue = "{{$frontend.Headers.CustomFrameOptionsValue}}"
  ContentTypeNosniff = {{$frontend.Headers.ContentTypeNosniff}}
  BrowserXSSFilter = {{$frontend.Headers.BrowserXSSFilter}}
  ContentSecurityPolicy = "{{$frontend.Headers.ContentSecurityPolicy}}"
  PublicKey = "{{$frontend.Headers.PublicKey}}"
  ReferrerPolicy = "{{$frontend.Headers.ReferrerPolicy}}"
  IsDevelopment = {{$frontend.Headers.IsDevelopment}}
{{if $frontend.Headers.CustomRequestHeaders}}
    [frontends."{{$frontendName}}".headers.customrequestheaders]
    {{range $k, $v := $frontend.Headers.CustomRequestHeaders}}
    {{$k}} = "{{$v}}"
    {{end}}
{{end}}
{{if $frontend.Headers.CustomResponseHeaders}}
    [frontends."{{$frontendName}}".headers.customresponseheaders]
    {{range $k, $v := $frontend.Headers.CustomResponseHeaders}}
    {{$k}} = "{{$v}}"
    {{end}}
{{end}}
{{if $frontend.Headers.AllowedHosts}}
    [frontends."{{$frontendName}}".headers.AllowedHosts]
    {{range $frontend.Headers.AllowedHosts}}
    "{{.}}"
    {{end}}
{{end}}
{{if $frontend.Headers.HostsProxyHeaders}}
    [frontends."{{$frontendName}}".headers.HostsProxyHeaders]
    {{range $frontend.Headers.HostsProxyHeaders}}
    "{{.}}"
    {{end}}
{{end}}
{{if $frontend.Headers.SSLProxyHeaders}}
    [frontends."{{$frontendName}}".headers.SSLProxyHeaders]
    {{range $k, $v := $frontend.Headers.SSLProxyHeaders}}
    {{$k}} = "{{$v}}"
    {{end}}
{{end}}
    {{range $routeName, $route := $frontend.Routes}}
    [frontends."{{$frontendName}}".routes."{{$routeName}}"]
    rule = "{{$route.Rule}}"
    {{end}}
{{end}}
`)

func templatesKubernetesTmplBytes() ([]byte, error) {
	return _templatesKubernetesTmpl, nil
}

func templatesKubernetesTmpl() (*asset, error) {
	bytes, err := templatesKubernetesTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/kubernetes.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesKvTmpl = []byte(`{{$frontends := List .Prefix "/frontends/" }}
{{$backends :=  List .Prefix "/backends/"}}
{{$tlsconfiguration := List .Prefix "/tlsconfiguration/"}}

[backends]{{range $backends}}
{{$backend := .}}
{{$backendName := Last $backend}}
{{$servers := ListServers $backend }}

{{$circuitBreaker := Get "" . "/circuitbreaker/" "expression"}}
{{with $circuitBreaker}}
[backends."{{$backendName}}".circuitBreaker]
    expression = "{{$circuitBreaker}}"
{{end}}

{{$loadBalancer := Get "" . "/loadbalancer/" "method"}}
{{with $loadBalancer}}
[backends."{{$backendName}}".loadBalancer]
    method = "{{$loadBalancer}}"
    sticky = {{ getSticky . }}
    {{if hasStickinessLabel $backend}}
    [backends."{{$backendName}}".loadBalancer.stickiness]
      cookieName = {{getStickinessCookieName $backend}}
    {{end}}
{{end}}

{{$healthCheck := Get "" . "/healthcheck/" "path"}}
{{with $healthCheck}}
[backends."{{$backendName}}".healthCheck]
    path = "{{$healthCheck}}"
    interval = "{{ Get "30s" $backend "/healthcheck/" "interval" }}"
{{end}}

{{$maxConnAmt := Get "" . "/maxconn/" "amount"}}
{{$maxConnExtractorFunc := Get "" . "/maxconn/" "extractorfunc"}}
{{with $maxConnAmt}}
{{with $maxConnExtractorFunc}}
[backends."{{$backendName}}".maxConn]
    amount = {{$maxConnAmt}}
    extractorFunc = "{{$maxConnExtractorFunc}}"
{{end}}
{{end}}

{{range $servers}}
[backends."{{$backendName}}".servers."{{Last .}}"]
    url = "{{Get "" . "/url"}}"
    weight = {{Get "0"  . "/weight"}}
{{end}}
{{end}}

[frontends]{{range $frontends}}
    {{$frontend := Last .}}
    {{$entryPoints := SplitGet . "/entrypoints"}}
    [frontends."{{$frontend}}"]
    backend = "{{Get "" . "/backend"}}"
    passHostHeader = {{Get "true" . "/passHostHeader"}}
    priority = {{Get "0" . "/priority"}}
    entryPoints = [{{range $entryPoints}}
      "{{.}}",
    {{end}}]
    {{$routes := List . "/routes/"}}
        {{range $routes}}
        [frontends."{{$frontend}}".routes."{{Last .}}"]
        rule = "{{Get "" . "/rule"}}"
        {{end}}
{{end}}

{{range $tlsconfiguration}}
{{$entryPoints := SplitGet . "/entrypoints"}}
[[tlsConfiguration]]
    entryPoints = [{{range $entryPoints}}
      "{{.}}",
    {{end}}]
    [tlsConfiguration.certificate]
        certFile = """{{Get "" . "/certificate" "/certfile"}}"""
        keyFile = """{{Get "" . "/certificate" "/keyfile"}}"""
{{end}}

`)

func templatesKvTmplBytes() ([]byte, error) {
	return _templatesKvTmpl, nil
}

func templatesKvTmpl() (*asset, error) {
	bytes, err := templatesKvTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/kv.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesMarathonTmpl = []byte(`{{$apps := .Applications}}

{{range $app := $apps}}
{{range $task := $app.Tasks}}
{{range $serviceIndex, $serviceName := getServiceNames $app}}
    [backends."backend{{getBackend $app $serviceName}}".servers."server-{{$task.ID | replace "." "-"}}{{getServiceNameSuffix $serviceName }}"]
    url = "{{getProtocol $app $serviceName}}://{{getBackendServer $task $app}}:{{getPort $task $app $serviceName}}"
    weight = {{getWeight $app $serviceName}}
{{end}}
{{end}}
{{end}}

{{range $app := $apps}}
{{range $serviceIndex, $serviceName := getServiceNames $app}}
[backends."backend{{getBackend $app $serviceName }}"]
{{ if hasMaxConnLabels $app }}
      [backends."backend{{getBackend $app $serviceName }}".maxconn]
        amount = {{getMaxConnAmount $app }}
        extractorfunc = "{{getMaxConnExtractorFunc $app }}"
{{end}}
{{ if hasLoadBalancerLabels $app }}
      [backends."backend{{getBackend $app $serviceName }}".loadbalancer]
        method = "{{getLoadBalancerMethod $app }}"
        sticky = {{getSticky $app}}
        {{if hasStickinessLabel $app}}
        [backends."backend{{getBackend $app $serviceName }}".loadbalancer.stickiness]
          cookieName = "{{getStickinessCookieName $app}}"
        {{end}}
{{end}}
{{ if hasCircuitBreakerLabels $app }}
      [backends."backend{{getBackend $app $serviceName }}".circuitbreaker]
        expression = "{{getCircuitBreakerExpression $app }}"
{{end}}
{{ if hasHealthCheckLabels $app }}
      [backends."backend{{getBackend $app $serviceName }}".healthcheck]
        path = "{{getHealthCheckPath $app }}"
        interval = "{{getHealthCheckInterval $app }}"
{{end}}
{{end}}
{{end}}

[frontends]{{range $app := $apps}}{{range $serviceIndex, $serviceName := getServiceNames .}}
  [frontends."{{ getFrontendName $app $serviceName }}"]
  backend = "backend{{getBackend $app $serviceName}}"
  passHostHeader = {{getPassHostHeader $app $serviceName}}
  priority = {{getPriority $app $serviceName}}
  entryPoints = [{{range getEntryPoints $app $serviceName}}
    "{{.}}",
  {{end}}]
  basicAuth = [{{range getBasicAuth $app $serviceName}}
    "{{.}}",
  {{end}}]
    [frontends."{{ getFrontendName $app $serviceName }}".routes."route-host{{$app.ID | replace "/" "-"}}{{getServiceNameSuffix $serviceName }}"]
    rule = "{{getFrontendRule $app $serviceName}}"
{{end}}{{end}}
`)

func templatesMarathonTmplBytes() ([]byte, error) {
	return _templatesMarathonTmpl, nil
}

func templatesMarathonTmpl() (*asset, error) {
	bytes, err := templatesMarathonTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/marathon.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesMesosTmpl = []byte(`{{$apps := .Applications}}
[backends]{{range .Tasks}}
    [backends.backend{{getBackend . $apps}}.servers.server-{{getID .}}]
    url = "{{getProtocol . $apps}}://{{getHost .}}:{{getPort . $apps}}"
    weight = {{getWeight . $apps}}
{{end}}

[frontends]{{range .Applications}}
  [frontends.frontend-{{getFrontEndName .}}]
  backend = "backend{{getFrontendBackend .}}"
  passHostHeader = {{getPassHostHeader .}}
  priority = {{getPriority .}}
  entryPoints = [{{range getEntryPoints .}}
    "{{.}}",
  {{end}}]
    [frontends.frontend-{{getFrontEndName .}}.routes.route-host{{getFrontEndName .}}]
    rule = "{{getFrontendRule .}}"
{{end}}
`)

func templatesMesosTmplBytes() ([]byte, error) {
	return _templatesMesosTmpl, nil
}

func templatesMesosTmpl() (*asset, error) {
	bytes, err := templatesMesosTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/mesos.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesNotfoundTmpl = []byte(`<!DOCTYPE html>
<html>
<head>
    <title>Traefik</title>
</head>
<body>
    Ohhhh man, this is bad...
</body>
</html>`)

func templatesNotfoundTmplBytes() ([]byte, error) {
	return _templatesNotfoundTmpl, nil
}

func templatesNotfoundTmpl() (*asset, error) {
	bytes, err := templatesNotfoundTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/notFound.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesRancherTmpl = []byte(`{{$backendServers := .Backends}}
[backends]{{range $backendName, $backend := .Backends}}
    {{if hasCircuitBreakerLabel $backend}}
    [backends.backend-{{$backendName}}.circuitbreaker]
      expression = "{{getCircuitBreakerExpression $backend}}"
    {{end}}

    {{if hasLoadBalancerLabel $backend}}
    [backends.backend-{{$backendName}}.loadbalancer]
      method = "{{getLoadBalancerMethod $backend}}"
      sticky = {{getSticky $backend}}
      {{if hasStickinessLabel $backend}}
      [backends.backend-{{$backendName}}.loadbalancer.stickiness]
        cookieName = "{{getStickinessCookieName $backend}}"
      {{end}}
    {{end}}

    {{if hasMaxConnLabels $backend}}
    [backends.backend-{{$backendName}}.maxconn]
      amount = {{getMaxConnAmount $backend}}
      extractorfunc = "{{getMaxConnExtractorFunc $backend}}"
    {{end}}

    {{range $index, $ip := $backend.Containers}}
      [backends.backend-{{$backendName}}.servers.server-{{$index}}]
      url = "{{getProtocol $backend}}://{{$ip}}:{{getPort $backend}}"
      weight = {{getWeight $backend}}
    {{end}}

{{end}}

[frontends]{{range $frontendName, $service := .Frontends}}
    [frontends."frontend-{{$frontendName}}"]
    backend = "backend-{{getBackend $service}}"
    passHostHeader = {{getPassHostHeader $service}}
    priority = {{getPriority $service}}
    redirect = "{{getRedirect $service}}"
    entryPoints = [{{range getEntryPoints $service}}
        "{{.}}",
    {{end}}]
    basicAuth = [{{range getBasicAuth $service}}
        "{{.}}",
    {{end}}]
    [frontends."frontend-{{$frontendName}}".routes."route-frontend-{{$frontendName}}"]
    rule = "{{getFrontendRule $service}}"
{{end}}
`)

func templatesRancherTmplBytes() ([]byte, error) {
	return _templatesRancherTmpl, nil
}

func templatesRancherTmpl() (*asset, error) {
	bytes, err := templatesRancherTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/rancher.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/consul_catalog.tmpl": templatesConsul_catalogTmpl,
	"templates/docker.tmpl":         templatesDockerTmpl,
	"templates/ecs.tmpl":            templatesEcsTmpl,
	"templates/eureka.tmpl":         templatesEurekaTmpl,
	"templates/kubernetes.tmpl":     templatesKubernetesTmpl,
	"templates/kv.tmpl":             templatesKvTmpl,
	"templates/marathon.tmpl":       templatesMarathonTmpl,
	"templates/mesos.tmpl":          templatesMesosTmpl,
	"templates/notFound.tmpl":       templatesNotfoundTmpl,
	"templates/rancher.tmpl":        templatesRancherTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"consul_catalog.tmpl": &bintree{templatesConsul_catalogTmpl, map[string]*bintree{}},
		"docker.tmpl":         &bintree{templatesDockerTmpl, map[string]*bintree{}},
		"ecs.tmpl":            &bintree{templatesEcsTmpl, map[string]*bintree{}},
		"eureka.tmpl":         &bintree{templatesEurekaTmpl, map[string]*bintree{}},
		"kubernetes.tmpl":     &bintree{templatesKubernetesTmpl, map[string]*bintree{}},
		"kv.tmpl":             &bintree{templatesKvTmpl, map[string]*bintree{}},
		"marathon.tmpl":       &bintree{templatesMarathonTmpl, map[string]*bintree{}},
		"mesos.tmpl":          &bintree{templatesMesosTmpl, map[string]*bintree{}},
		"notFound.tmpl":       &bintree{templatesNotfoundTmpl, map[string]*bintree{}},
		"rancher.tmpl":        &bintree{templatesRancherTmpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
