# Golang cli sample to fetch Azure DevOps projects with optional filtering

## Creating test projects in bulk to enable testing of the fetch calls

The [createProject.go](./createProjects/createProject.go) can be used to create test projects in Azure Devops. The command provides the following switches

* **organizationUrl**: Your Azure DevOps organization url
* **personalAccessToken**: your PAT
* **noOfProjects**: The number of projects to create
* **projectNamePrefix**: Project Names will start with this prefix, and then appended by a incrementing number
* **projectNameStartSuffix**: The start index which will be appended to the project name. If this value is 10, then first project name suffix will be 10, then 11 and so forth
* **noOfWorkers**: The number of go routines used to create these projects, the default value is 1. If you are creating large number of projects like 100, then you can set this to a value like 5.

### Sample commands:

*  Create 10 projects

```
export azdoPAT=YOUR_AzureDevOps_PAT
go run createProjects/createProject.go --organizationUrl="https://dev.azure.com/YOUR_ORGANIZATION/" --personalAccessToken=$azdoPAT --noOfProjects=10 --projectNamePrefix="tstazdoprj" --projectNameStartSuffix=300
```

*  create 100 projects with 10 workers

```
go run createProjects/createProject.go --organizationUrl="https://dev.azure.com/YOUR_ORGANIZATION/" --personalAccessToken=$azdoPAT --noOfProjects=100 --projectNamePrefix="tstazdoprj" --projectNameStartSuffix=400 --noOfWorkers=10
```

## Fetching Azure DevOps Projects optionally filtering by date (LastUpdateTime)
The [getProjects.go](./getProjects.go) can be used to fetch summary and details of Azure DevOps Projects. We can optionally get a filtered result of projects last modified before a specified date. This command provides the following switches:

* **organizationUrl**: Your Azure DevOps organization url
* **personalAccessToken**: your PAT
* **filterDate**: This is a date in the **"yyyy-mm-dd 15:04"** format, if specified only projects last modified prior to this date are returned
* **showProjectDetails**: As we will see from the examples below, by default only the summary regarding the number of projects is returned. If this switch is provided we also get a csv formatted list of project details.

### Sample Commands

* Get summary of all projects

```bash
> export azdoPAT=YOUR_AzureDevOps_PAT
> go run getProjects.go --organizationUrl="https://dev.azure.com/YOUR_ORGANIZATION/" --personalAccessToken=$azdoPAT

2022/08/01 17:47:32 count of all projects: 319

```

* Get details of all projects

```bash
> export azdoPAT=YOUR_AzureDevOps_PAT
> go run getProjects.go --organizationUrl="https://dev.azure.com/YOUR_ORGANIZATION/" --personalAccessToken=$azdoPAT --showProjectDetails

2022/08/01 17:47:32 count of all projects: 319


Project Name,Project ID,Last Updated
tstazdoprj47,xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx,2022-07-30 15:54
tstazdoprj100,xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx,2022-07-30 15:55
tstazdoprj76,xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx,2022-07-30 15:54
.
.
.
```
  
* Get summary of projects by filter date

```bash
> export azdoPAT=YOUR_AzureDevOps_PAT
> go run getProjects.go --organizationUrl="https://dev.azure.com/YOUR_ORGANIZATION/" --personalAccessToken=$azdoPAT --filterDate='2022-08-01 00:00'

2022/08/01 17:51:16 count of all projects: 319
2022/08/01 17:51:16 Filtering projects by date
2022/08/01 17:51:16 count of all filtered projects: 240
```
  
* Get details of projects by filter date

```bash
> export azdoPAT=YOUR_AzureDevOps_PAT
> go run getProjects.go --organizationUrl="https://dev.azure.com/YOUR_ORGANIZATION/" --personalAccessToken=$azdoPAT --filterDate='2022-08-01 00:00' --showProjectDetails

2022/08/01 17:52:22 count of all projects: 319
2022/08/01 17:52:22 Filtering projects by date
2022/08/01 17:52:22 count of all filtered projects: 240


Project Name,Project ID,Last Updated
tstazdoprj47,xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx,2022-07-30 15:54
tstazdoprj100,xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx,2022-07-30 15:55
tstazdoprj76,xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx,2022-07-30 15:54
.
.
.
```

