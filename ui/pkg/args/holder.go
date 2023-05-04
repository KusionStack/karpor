// Copyright The Karbour Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package args

import (
	"net"
)

var Holder = &holder{}

// Argument holder structure. It is private to make sure that only 1 instance can be created. It holds all
// arguments values passed to Dashboard binary.
type holder struct {
	insecurePort int
	port         int

	insecureBindAddress net.IP
	bindAddress         net.IP

	defaultCertDir string
	certFile       string
	keyFile        string
	staticDir      string

	autoGenerateCertificates bool
}

// GetInsecurePort 'insecure-port' argument of Dashboard binary.
func (self *holder) GetInsecurePort() int {
	return self.insecurePort
}

// GetPort 'port' argument of Dashboard binary.
func (self *holder) GetPort() int {
	return self.port
}

// GetInsecureBindAddress 'insecure-bind-address' argument of Dashboard binary.
func (self *holder) GetInsecureBindAddress() net.IP {
	return self.insecureBindAddress
}

// GetBindAddress 'bind-address' argument of Dashboard binary.
func (self *holder) GetBindAddress() net.IP {
	return self.bindAddress
}

// GetDefaultCertDir 'default-cert-dir' argument of Dashboard binary.
func (self *holder) GetDefaultCertDir() string {
	return self.defaultCertDir
}

// GetAutoGenerateCertificates 'auto-generate-certificates' argument of Dashboard binary.
func (self *holder) GetAutoGenerateCertificates() bool {
	return self.autoGenerateCertificates
}

// GetStaticDir 'static-dir' argument of Dashboard binary.
func (self *holder) GetStaticDir() string {
	return self.staticDir
}
