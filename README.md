# Deviceplane Controller

Deviceplane is an open source device management tool for embedded systems and edge computing. It solves various infrastructure problems related to remote device management such as:

- Network connectivity and SSH access
- Orchestration and deployment of remote updates
- Host and application monitoring
- Device organization: naming, labeling, searching, and filtering of devices
- Access and security controls

Deviceplane integrates with your device by running a lightweight static binary via your system supervisor. It can be used with nearly any Linux distro, which means you can continue using Ubuntu, Raspbian, a Yocto build, or whatever else fits your needs.

A hosted version of Deviceplane is available at [https://cloud.deviceplane.com/](https://cloud.deviceplane.com/).

## Documentation

Visit <a aria-label="next.js learn" href="https://deviceplane.com/docs">https://deviceplane.com/docs</a> to view the full documentation.

THIS REPO IS WIP

#### Build the binary

```
make build DEBUG=1
```

#### Run the binary

```
bin/controller --addr :80
	       --mysql "..."
	       --email-provider smtp
	       --email-from-address noreply@deviceplane.com
	       --smtp-server smtp.sendgrid.net
	       --smtp-port "465"
	       --smtp-username apikey
	       --smtp-password "..."
	       --auth0-audience "..."
	       --auth0-domain "..."
	       --db-max-open-conns "5"
	       --db-max-idle-conns "5"
	       --db-max-conn-lifetime 5m
	       --allowed-origin https://cloud.dev.edgeworx.io
	       --allowed-origin http://localhost:3000
	       --allowed-origin https://localhost:3000
```

#### Run the binary as intercept to Kubernetes pod

Install telepresence

```
brew install datawire/blackbird/telepresence
```

Intercept

```
telepresence connect
telepresence intercept deviceplane -n deviceplane --port 80
```

Disconnect

```
telepresence leave deviceplane-deviceplane
telepresence quit
telepresence uninstall --everything
```

## Support

For bugs, issues, and feature requests please [submit](//github.com/deviceplane/controller/issues/new) a GitHub issue.

## License

Copyright (c) Deviceplane, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.