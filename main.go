package main

import (
	"cloud.google.com/go/container/apiv1"
	pb "google.golang.org/genproto/googleapis/container/v1"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	MachineType    = envString("MACHINE_TYPE", "n1-standard-1")
	Project        = os.Getenv("PROJECT")
	Location       = os.Getenv("LOCATION")
	Cluster        = os.Getenv("CLUSTER")
	NodePool       = envString("NODE_POOL", "pool-1")
	ServiceAccount = os.Getenv("SERVICE_ACCOUNT")
	OAuthScopes    = os.Getenv("OAUTH_SCOPES")
	RandomPath     = envString("RANDOM_PATH", "8YkpGd2LQN3nBuWJfXRb")
	Preemptible = envBool("PREEMPTIBLE", true)
	DiskSize = envInt("DISK_SIZE", 50)
	InitialNodeCount = envInt("INITIAL_NODE_COUNT", 1)
	MaxNodeCount = envInt("MAX_NODE_COUNT", 2)
)

func envString(e, d string) string {
	s := os.Getenv(e)
	if s == "" {
		return d
	}

	return s
}

func envBool(e string, d bool) bool {
	s := os.Getenv(e)
	if s == "" {
		return d
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatalf("%s has invalid option, it should be a boolean value", e)
	}

	return b
}

func envInt(e string, d int) int {
	s := os.Getenv(e)
	if s == "" {
		return d
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s has invalid option it should be a number value", e)
	}

	return i
}


func main() {

	http.HandleFunc(fmt.Sprintf("/%s/up", RandomPath), spinUp)
	http.HandleFunc(fmt.Sprintf("/%s/down", RandomPath), spinDown)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

func spinUp(w http.ResponseWriter, r *http.Request) {

	c, err := container.NewClusterManagerClient(r.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &pb.CreateNodePoolRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/clusters/build", Project, Location),
		NodePool: &pb.NodePool{
			Name: NodePool,
			Autoscaling: &pb.NodePoolAutoscaling{
				Enabled: true,
				MaxNodeCount: int32(MaxNodeCount),
				MinNodeCount: 0,
			},
			Management: &pb.NodeManagement{
				AutoRepair: true,
				AutoUpgrade: false,
			},
			InitialNodeCount: int32(InitialNodeCount),
			Config: &pb.NodeConfig{
				DiskSizeGb: int32(DiskSize),
				MachineType: MachineType,
				ServiceAccount: ServiceAccount,
				OauthScopes: strings.Split(OAuthScopes, ","),
				Preemptible: Preemptible,
			},
		},
	}

	_, err = c.CreateNodePool(r.Context(), req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func spinDown(w http.ResponseWriter, r *http.Request) {

	c, err := container.NewClusterManagerClient(r.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &pb.DeleteNodePoolRequest{
		Name:fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", Project, Location, Cluster, NodePool),
	}

	_, err = c.DeleteNodePool(r.Context(), req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
