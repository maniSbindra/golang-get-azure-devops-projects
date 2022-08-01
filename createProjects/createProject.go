package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
)

func main() {

	organizationUrl := flag.String("organizationUrl", "", "URL of the Azure DevOps organization")
	personalAccessToken := flag.String("personalAccessToken", "", "Personal Access Token")
	noOfProjects := flag.Int("noOfProjects", 100, "Number of AzDO projects to be created")
	projectNamePrefix := flag.String("projectNamePrefix", "testazdproject", "Prefix for the project name")
	projectNameStartSuffix := flag.Int("projectNameStartSuffix", 1, "Start suffix for the project name")
	noOfWorkers := flag.Int("noOfWorkers", 1, "Number of workers to be used for creating projects")

	flag.Parse()

	if *organizationUrl == "" || *personalAccessToken == "" {
		log.Fatal("organizationUrl and personalAccessToken are required to proceed further")
	}

	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(*organizationUrl, *personalAccessToken)

	ctx := context.Background()

	// Create a client to interact with the Core area
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	workQueue := make(chan *core.TeamProject, *noOfProjects)
	termChan := make(chan bool, *noOfWorkers)
	respChan := make(chan bool, *noOfProjects)

	// start workers to create projects
	for i := 0; i < *noOfWorkers; i++ {
		go createProject(ctx, coreClient, workQueue, termChan, respChan)
	}
	log.Println("All workers started")

	// initialize common project properties
	versionControl := make(map[string]string, 0)
	versionControl["sourceControlType"] = "Git"
	processTemplate := make(map[string]string, 0)
	processTemplate["templateTypeId"] = "6b724908-ef14-45cf-84f8-768b5384da45"
	capabilities := make(map[string]map[string]string, 0)
	capabilities["versioncontrol"] = versionControl
	capabilities["processTemplate"] = processTemplate

	// queue projects for creation
	for i := *projectNameStartSuffix; i < (*noOfProjects + *projectNameStartSuffix); i++ {
		projName := fmt.Sprintf("%s%d", *projectNamePrefix, i)
		teamProject := &core.TeamProject{
			Name:         &projName,
			Description:  &projName,
			Visibility:   &core.ProjectVisibilityValues.Private,
			Capabilities: &capabilities,
		}
		workQueue <- teamProject
	}
	log.Println("All projects queued for creation")

	// wait for all projects to be created
	for i := *projectNameStartSuffix; i < (*noOfProjects + *projectNameStartSuffix); i++ {
		<-respChan
		log.Printf("Main: Project creation request completed (%d) \n", i)
	}

	// terminate all workers
	for i := 0; i < *noOfWorkers; i++ {
		termChan <- true
	}
	log.Println("All workers terminated")

}

func createProject(ctx context.Context, azdoClient core.Client, workQueue chan *core.TeamProject, termChan chan bool, respChan chan<- bool) {

	var teamProject *core.TeamProject
	for {
		select {
		case <-termChan: // terminate this worker if message on termination channel
			return
		case teamProject = <-workQueue: // create projects while message on work queue
			crtPrjArgs := core.QueueCreateProjectArgs{
				ProjectToCreate: teamProject,
			}

			_, err := azdoClient.QueueCreateProject(ctx, crtPrjArgs)
			if err != nil {
				log.Printf("Error creating project: %v \n", err)
			}

			log.Printf("Project Creation request submitted successfully for project %s \n", *teamProject.Name)
			respChan <- true
		}
	}

}
