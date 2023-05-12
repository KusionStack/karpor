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

# Copy front-end code
ADD ./ui /root/ui
WORKDIR /root/ui

# Install dependencies
RUN npm install

# Build, generate static files
RUN npm run build

FROM alpine:3.17.3 AS production

WORKDIR /

# Copy the static file directory built in the previous layer to the current layer
COPY --from=ui-builder /root/ui/build /static
COPY karbour .

# Copy nonroot user
COPY --from=ui-builder /etc/passwd /etc/passwd

USER root

ENTRYPOINT ["/karbour"]
