# Copyright The Karbour Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Front-end build layer, builds front-end code and generates static files.
FROM node:20 AS ui-builder

# Create a nonroot user for final image
RUN useradd -u 10001 nonroot

# Copy front-end code
ADD . /root/ui
WORKDIR /root/ui

# Install dependencies
RUN npm install

# Build, generate static files
RUN npm run build

# Golang build layer, builds golang code and generates entrypoint binary.
FROM golang:1.19 AS go-builder

ADD . /root
WORKDIR /root
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /karbour-dashboard \
    ./main.go

# Production layer
FROM scratch AS production

# Copy the static file directory built in the previous layer to the current layer
COPY --from=ui-builder /root/ui/build /dist
COPY --from=go-builder /karbour-dashboard /karbour-dashboard

# Copy nonroot user
COPY --from=ui-builder /etc/passwd /etc/passwd
USER nonroot

# The port that the application listens on.
EXPOSE 9871 9872

ENTRYPOINT ["/karbour-dashboard", "--insecure-bind-address=0.0.0.0", "--bind-address=0.0.0.0"]
