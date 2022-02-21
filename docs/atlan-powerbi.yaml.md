# atlan-powerbi

> This file is generated by argodocs. Please do not Edit.

|Kind|Version|Entrypoint Template|Last Updated At|
|----|----|----|----|
|WorkflowTemplate|argoproj.io/v1alpha1|main|Tuesday, 22-Feb-22 01:18:25 IST|

S3 Output Reference - which directory will have what data once the workflow completes

``` 
<bucket_name>/argo-artifacts/<connection-qualified-name>
├── current-state                 -> tarball of `transformed_metadata/`
├── lineage-current-state         -> tarball of `lineage-output`
├── lineage-diff/                 -> output of diff calculation in the lineage step
├── lineage-output/               -> output of `generate-lineage` template
├── metadata-diff/                -> output of diff calculation in the publish step
├── metadata-extract/             -> output of `process-powerbi-results` template
├── refreshable-data/             -> output of `fetch-refreshables` step in `extract-powerbi-metadata` template
├── relationship-lineage-output/  -> output of `generate-relationship-lineage` template 
├── scan-results/                 -> output of `get-scan-results` step of `run-scan` template
└── transformed-metadata/         -> output of `transform` step in `atlan-crawler/generic-publish` template
```
## Templates

A list of all the templates present in the Workflow Template

|Name|Type|Line Number|
|----|----|----|
|[main](#main)|`DAG`|28|
|[extract-powerbi-metadata](#extract-powerbi-metadata)|`DAG`|148|
|[publish-powerbi-metadata](#publish-powerbi-metadata)|`DAG`|463|
|[run-scan](#run-scan)|`DAG`|627|
|[add-page-counts](#add-page-counts)|`Script`|751|
|[process-powerbi-results](#process-powerbi-results)|`Container`|805|
|[filter-workspaces](#filter-workspaces)|`Script`|842|
|[create-endorsement-requests](#create-endorsement-requests)|`Script`|901|
|[lineage](#lineage)|`DAG`|998|
|[generate-relationship-lineage](#generate-relationship-lineage)|`Container`|1107|
|[generate-lineage](#generate-lineage)|`Container`|1161|
|[api-request](#api-request)|`DAG`|1220|
|[powerbi-basic-auth](#powerbi-basic-auth)|`Container`|1342|
|[powerbi-service-principal-auth](#powerbi-service-principal-auth)|`Container`|1567|

---

### main

Type: `DAG`

main (DAG): primary dag of the workflow with two steps, extract and publish, 
extract -> fetches data from PowerBI APIs and writes them to S3, 
publish -> transforms and writes that data to the metastore and then generates BI<>BI lineage and BI<>SQL lineage 
- Inputs
    - Parameters
        - `credentials-fetch-strategy` - Credential
        - `credential-guid`
        - `connection` - Connection Entity (stringified JSON)
        - `include-filter` - Workspace Filters
        - `exclude-filter`
        - `atlas-auth-type` - Atlas auth type
        - `publish-mode` - Publish Mode
        - `endorsement-attach-mode` - PowerBI Endorsement Attachment
- Tasks
    - `extract`
        - snkdfjnsdjkfnsdjk
        - Template: [extract-powerbi-metadata](#extract-powerbi-metadata)
    - `publish`
        - Template: [publish-powerbi-metadata](#publish-powerbi-metadata)

---

### extract-powerbi-metadata

Type: `DAG`

extract-powerbi-metadata (DAG): responsible for extracting metadata from PowerBI using the Metadata Scanner API
fetch-credentials   -> fetch credentials from the api service and only write the `authType` key value as output
auth-type-to-param  -> convert that `authType` output from an artifact to a parameter
workspaces          -> get all the workspaces in the account, output data has `name` and `id` of the workspace 
filter-workspaces   -> filter out workspaces using the include-filter and exclude-filter regex on workspace `name`
run-scans           -> divide the filtered workspaces into chunks of 100 and start metadata scans on them (reason: API limit)
fetch-refreshables  -> fetch the freshables data, this is extra metadata about how often a PowerBI dataset is refreshed
process-results     -> run https://github.com/atlanhq/marketplace-scripts/tree/master/marketplace_scripts/powerbi/scanner_results
fetch-pages         -> fetch metadata about pages
add-page-counts     -> add counts of these pages to the reports metadata
- Inputs
    - Parameters
        - `credentials-fetch-strategy`
        - `credential-guid`
        - `connection-qualified-name`
        - `heracles-uri`
        - `include-filter`
        - `exclude-filter`
        - `endorsement-attach-mode`
        - `git-kube-secret-name`
        - `git-kube-ssh-key`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
- Tasks
    - `fetch-credentials`
        - Template: rest-api::oauth2-client-credentials
    - `auth-type-to-param`
        - Template: utils::artifact-to-key-param
    - `workspaces`
        - Template: [api-request](#api-request)
    - `filter-workspaces`
        - Template: [filter-workspaces](#filter-workspaces)
    - `run-scans`
        - Template: [run-scan](#run-scan)
    - `fetch-refreshables`
        - Template: [api-request](#api-request)
    - `process-results`
        - Template: [process-powerbi-results](#process-powerbi-results)
    - `fetch-pages`
        - Template: [api-request](#api-request)
    - `add-page-counts`
        - Template: [add-page-counts](#add-page-counts)

---

### publish-powerbi-metadata

Type: `DAG`

publish-powerbi-metadata (DAG): transform result ->  create connection entity -> upload assetss to atlas -> publish lineage
publish                       -> run atlan-crawler/generic-publish template
generate-endorsement-requests -> generate requests payload from the endorsements data (when: endorsement-attach-mode == requests)
publish-endorsement-requests  -> publish those endorsement attachment requests to heracles
lineage                       -> generate BI<>BI and BI<>SQL lineage and publish it
- Inputs
    - Parameters
        - `connection`
        - `mode`
        - `source`
        - `endorsement-attach-mode`
        - `atlas-api-uri`
        - `heracles-uri`
        - `atlan-web-kube-secret`
        - `atlas-auth-type`
        - `publish-chunk-size`
        - `git-kube-secret-name`
        - `git-kube-ssh-key`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
- Tasks
    - `publish`
        - Template: atlan-crawler::generic-publish
    - `generate-endorsement-requests`
        - Template: [create-endorsement-requests](#create-endorsement-requests)
    - `publish-endorsement-requests`
        - Template: atlan-api::create-requests-bulk
    - `lineage`
        - Template: [lineage](#lineage)

---

### run-scan

Type: `DAG`

run-scan (DAG): Starts the metadata scan for a set of workspaces (max 100)
start-scans               -> Trigger metadata scan for the set of input workspaces, this returns a `scanId`
convert-scan-id-to-param  -> Convert that `scanId` from artifact to parameter
check-scan-status         -> Long poll (wait = 10s) for the status of that scan, until status == succeeded
get-scan-results          -> As scan succeeded, get the result of that scan and write it on S3
- Inputs
    - Parameters
        - `credential-guid`
        - `index`
        - `auth-type`
        - `connection-qualified-name`
        - `workspaces`
        - `git-kube-secret-name`
        - `git-kube-ssh-key`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
- Tasks
    - `start-scans`
        - Template: [api-request](#api-request)
    - `convert-scan-id-to-param`
        - Template: utils::artifact-to-key-param
    - `check-scan-status`
        - Template: [api-request](#api-request)
    - `get-scan-results`
        - Template: [api-request](#api-request)

---

### add-page-counts

Type: `Script`

add-page-counts (Script): read pages and reports data -> add count of pages in the reports by matching on id
- Inputs
    - Parameters
        - `output-prefix`
    - Artifacts
        - `metadata-extract`
- Outputs
    - Artifacts
        - `output`

---

### process-powerbi-results

Type: `Container`

process-powerbi-results (Container): Parse the result from metadata scanning and write as individual asset files
run -> https://github.com/atlanhq/marketplace-scripts/tree/master/marketplace_scripts/powerbi/scanner_results
- Inputs
    - Parameters
        - `connection-qualified-name`
        - `endorsement-attach-mode`
    - Artifacts
        - `scan-results`
        - `refreshable-data`
        - `scripts`
- Outputs
    - Artifacts
        - `output`

---

### filter-workspaces

Type: `Script`

filter-workspaces (Script): Filter the workspaces with include-filter and exclude-filder as regex and divide into chunks of 100
- Inputs
    - Parameters
        - `include-filter`
        - `exclude-filter`
    - Artifacts
        - `workspaces`
- Outputs
    - Parameters
        - `output`

---

### create-endorsement-requests

Type: `Script`

create-endorsement-requests (Script): Read the output from the scanner_results script for endorsements 
and use that to construct static metadta update requests to send to the api service
- Inputs
    - Artifacts
        - `publish-results`
        - `metadata-extract`
- Outputs
    - Artifacts
        - `output`

---

### lineage

Type: `DAG`

linegae (DAG): Generate BI<>BI and BI<>SQL lineage
generate-relationship-lineage-assets  -> run `generate-relationship-lineage` template, generates BI<>BI lineage
generate-lineage-assets               -> run `generate-lineage template`, generate BI<>SQL lineage
publish-lineage                       -> publish BI<>SQL Lineage, run atlan-crawler/generic-lineage-publish template
publish-relationship-lineage          -> publish BI<>BI Lineage, run atlan-crawler/generic-lineage-publish template
- Inputs
    - Parameters
        - `connection`
        - `heracles-uri`
        - `atlas-api-uri`
        - `atlan-web-kube-secret`
        - `publish-chunk-size`
        - `git-kube-secret-name`
        - `git-kube-ssh-key`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
- Tasks
    - `generate-relationship-lineage-assets`
        - Template: [generate-relationship-lineage](#generate-relationship-lineage)
    - `generate-lineage-assets`
        - Template: [generate-lineage](#generate-lineage)
    - `publish-lineage`
        - Template: atlan-crawler::generic-lineage-publish
    - `publish-relationship-lineage`
        - Template: atlan-crawler::generic-lineage-publish

---

### generate-relationship-lineage

Type: `Container`

generate-relationship-lineage (Container): Uses transformed metadata to generate BI<>BI lineage
run github.com/atlanhq/marketplace-scripts/tree/master/marketplace_scripts/lineage/powerbi_main.py 
- Inputs
    - Parameters
        - `connection-name`
        - `connection-qualified-name`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
    - Artifacts
        - `transfomed-metadata`
        - `scripts`
- Outputs
    - Artifacts
        - `output`

---

### generate-lineage

Type: `Container`

generate-lineage (Container): Uses scanner result and connection cache to generate BI<>SQL lineage
run github.com/atlanhq/marketplace-scripts/tree/master/marketplace_scripts/relationship_lineage/cmd/powerbi.py
- Inputs
    - Parameters
        - `connection-name`
        - `connection-qualified-name`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
    - Artifacts
        - `scan-results`
        - `connection-cache`
        - `scripts`
- Outputs
    - Artifacts
        - `output`

---

### api-request

Type: `DAG`

api-request (DAG): Conditional DAG that for authentication types `basic` and `service_principal`
basic-auth-request        -> run `powerbi-basic-auth` template, OAuth2 request with resource owner username and password
service-principal-request -> run `powerbi-service-principal-auth` template, OAuth2 reqeust with just client credentials
- Inputs
    - Parameters
        - `auth-type`
        - `credential-guid`
        - `method`
        - `url`
        - `request-config`
        - `output-prefix`
        - `pagination-wait-time`
        - `raw-input-file-pattern`
        - `raw-input-multiline`
        - `raw-input-paginate`
        - `output-chunk-size`
        - `execution-script`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
    - Artifacts
        - `raw-input`
- Outputs
    - Parameters
        - `success-num-files`
        - `failure-num-files`
    - Artifacts
        - `success`
        - `failure`
- Tasks
    - `basic-auth-request`
        - Template: [powerbi-basic-auth](#powerbi-basic-auth)
    - `service-principal-request`
        - Template: [powerbi-service-principal-auth](#powerbi-service-principal-auth)

---

### powerbi-basic-auth

Type: `Container`

powerbi-basic-auth (Container, InitContainer): OAuth2 request with resource owner username and password
InitContainer fetches the credentials and writes them to the shared `credentials` volume which is read back by the main `rest-master` container
- Inputs
    - Parameters
        - `credential-guid`
        - `method`
        - `url`
        - `request-config`
        - `oauth2-scope`
        - `pagination-wait-time`
        - `raw-input`
        - `raw-input-file-pattern`
        - `raw-input-file-sort`
        - `raw-input-multiline`
        - `execution-script`
        - `raw-input-paginate`
        - `output-chunk-size`
        - `output-secret-name`
        - `output-secret-namespace`
        - `max-retries`
        - `output-prefix`
        - `impersonate`
        - `kube-secret-name`
        - `client-id-env`
        - `client-secret-env`
        - `token-url-env`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
        - `heracles-uri`
        - `init-container-execution-script`
    - Artifacts
        - `raw-input`
        - `raw-input-file`
- Outputs
    - Parameters
        - `success-num-files`
        - `failure-num-files`
    - Artifacts
        - `success`
        - `failure`

---

### powerbi-service-principal-auth

Type: `Container`

powerbi-service-principal-auth (Container, InitContainer): OAuth2 request with just client credentials
InitContainer fetches the credentials and writes them to the shared `credentials` volume which is read back by the main `rest-master` container
- Inputs
    - Parameters
        - `credential-guid`
        - `method`
        - `url`
        - `request-config`
        - `oauth2-scope`
        - `raw-input`
        - `raw-input-file-pattern`
        - `raw-input-file-sort`
        - `raw-input-multiline`
        - `execution-script`
        - `raw-input-paginate`
        - `output-chunk-size`
        - `output-secret-name`
        - `output-secret-namespace`
        - `pagination-wait-time`
        - `max-retries`
        - `output-prefix`
        - `kube-secret-name`
        - `client-id-env`
        - `client-secret-env`
        - `token-url-env`
        - `statsd-host`
        - `statsd-port`
        - `statsd-global-tags`
        - `heracles-uri`
        - `init-container-execution-script`
    - Artifacts
        - `raw-input`
        - `raw-input-file`
- Outputs
    - Parameters
        - `success-num-files`
        - `failure-num-files`
    - Artifacts
        - `success`
        - `failure`
