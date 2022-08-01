package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
)

func main() {

	organizationUrl := flag.String("organizationUrl", "", "URL of the Azure DevOps organization")
	personalAccessToken := flag.String("personalAccessToken", "", "Personal Access Token")

	// This is blank by default needs to be in the format of "yyyy-mm-dd 15:04"
	filterDate := flag.String("filterDate", "", "Date to filter projects by")

	showProjectDetails := flag.Bool("showProjectDetails", false, "Show project details")

	flag.Parse()

	if *organizationUrl == "" || *personalAccessToken == "" {
		log.Fatal("organizationUrl, personalAccessToken and filterDate are required to proceed further")
	}

	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(*organizationUrl, *personalAccessToken)

	ctx := context.Background()

	// Create a client to interact with the Core area
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	filteredProjs, err := getAllProjects(ctx, coreClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("count of all projects: %v\n", len(filteredProjs))

	// if filter date is provided then filter the projects by the date
	if *filterDate != "" {
		log.Println("Filtering projects by date")
		filterDateTime, err := time.Parse("2006-01-02 15:04", *filterDate)
		if err != nil {
			log.Fatalf("Filter date is not in correct format: %v", err)
		}
		filteredProjs = getProjectsByFilter(filteredProjs, filterDateTime)
		log.Printf("count of all filtered projects: %v\n", len(filteredProjs))
	}

	if *showProjectDetails {
		printProjects(filteredProjs)
	}

}

func getAllProjects(ctx context.Context, coreClient core.Client) ([]core.TeamProjectReference, error) {

	// Get first page of the list of team projects for your organization
	responseValue, err := coreClient.GetProjects(ctx, core.GetProjectsArgs{})
	if err != nil {
		log.Fatal(err)
	}

	var allProjects []core.TeamProjectReference

	// index := 0
	for responseValue != nil {
		// Log the page of team project names
		allProjects = append(allProjects, (*responseValue).Value...)

		// if continuationToken has a value, then there is at least one more page of projects to get
		if responseValue.ContinuationToken != "" {
			// Get next page of team projects
			projectArgs := core.GetProjectsArgs{
				ContinuationToken: &responseValue.ContinuationToken,
			}
			responseValue, err = coreClient.GetProjects(ctx, projectArgs)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			// log.Println("Fetching next set of projects, and continuing")
		} else {
			responseValue = nil
		}
	}
	return allProjects, nil
}

func getProjectsByFilter(projects []core.TeamProjectReference, filterTime time.Time) []core.TeamProjectReference {
	var filteredProjects []core.TeamProjectReference

	for _, project := range projects {
		if (*(project.LastUpdateTime)).Time.Before(filterTime) {
			filteredProjects = append(filteredProjects, project)
		}
	}
	return filteredProjects
}

func printProjects(projects []core.TeamProjectReference) {
	fmt.Printf("\n\n%s,%s,%s\n", "Project Name", "Project ID", "Last Updated")
	for _, project := range projects {
		fmt.Printf("%s,%s,%s\n", *project.Name, *project.Id, (*project.LastUpdateTime).Time.Format("2006-01-02 15:04"))
	}
}
