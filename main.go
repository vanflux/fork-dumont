package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Request struct {
	RepoURL string `json:"repo_url"`
	Port string `json:"container_port"`
	ContainerId string `json:"container_id"`
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("./frontend/react/dist")))   
	http.HandleFunc("/cloneRepo", handleClone)
	http.HandleFunc("/getContainerList", handleContainerList)
	http.HandleFunc("/deleteContainer", handleDeleteContainer)

	fmt.Println("Server started on :3322")
	err := http.ListenAndServe(":3322", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func gitClone(repoURL, dest string) error {
	cmd := exec.Command("git", "clone", repoURL, dest)
	return cmd.Run()
}

func dockerBuild(imageName, context string) error {
	cmd := exec.Command("docker", "build", "-t", imageName, context)
	return cmd.Run()
}

func dockerRun(imageName string, port string) error {
    cmd := exec.Command("docker", "run", "-d", "-p", port+":"+port, imageName)
	return cmd.Run()
}

func handleClone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}

	rId := fmt.Sprintf("%d", rand.Intn(100))
	repoPath := "repositories/" + rId
	if err := gitClone(req.RepoURL, repoPath); err != nil {
		http.Error(w, fmt.Sprintf("Error cloning repository: %s", err), http.StatusInternalServerError)
		return
	}

	imageName := fmt.Sprintf("repo_image_%s", rId)
	if err := dockerBuild(imageName, repoPath); err != nil {
		http.Error(w, fmt.Sprintf("Error building Docker image: %s", err), http.StatusInternalServerError)
		return
	}

	if err := dockerRun(imageName, req.Port); err != nil {
		http.Error(w, fmt.Sprintf("Error running Docker container: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleContainerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.Command("docker", "ps", "-a", "--format", "{{json .}}")
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Error getting container list: %s, Details: %s", err, stderrBuf.String()), http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(&stdoutBuf)
	var containers []json.RawMessage
	for scanner.Scan() {
		containers = append(containers, json.RawMessage(scanner.Bytes()))
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error reading container list: %s", err), http.StatusInternalServerError)
		return
	}

	jsonContainers, err := json.Marshal(containers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling container list: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonContainers)
}

func handleDeleteContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}
	if err := cmdFactory("docker stop " + req.ContainerId); err != nil {
		http.Error(w, fmt.Sprintf("Error stoping container: %s", err), http.StatusInternalServerError)
		return
	}
		if err := cmdFactory("docker remove " + req.ContainerId); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting container: %s", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}


func cmdFactory(command string)error{
    args := strings.Fields(command)
    if len(args) == 0 {
        fmt.Println("No command provided")
    }
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run() 
}
