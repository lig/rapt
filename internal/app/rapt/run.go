package rapt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"codeberg.org/lig/rapt/internal/k8s"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// MountSpec represents a file mount specification
type MountSpec struct {
	LocalPath     string
	ContainerPath string
}

// RunTool executes a tool by creating a Kubernetes Job
func RunTool(namespace, toolName string, args map[string]string, envVars map[string]string, mounts []MountSpec, wait, follow bool, timeout int) error {
	// Initialize clients
	dynClient, err := k8s.InitDynamicClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to initialize dynamic client: %w", err)
	}

	k8sClient, err := k8s.InitKubernetesClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to initialize kubernetes client: %w", err)
	}

	// Get the tool definition
	tool, err := getToolDefinition(dynClient, namespace, toolName)
	if err != nil {
		return fmt.Errorf("failed to get tool definition: %w", err)
	}

	// Create ConfigMaps for mounted files first
	jobName := fmt.Sprintf("%s-%s", toolName, time.Now().Format("20060102-150405"))
	for i, mount := range mounts {
		// Read the local file
		fileContent, err := os.ReadFile(mount.LocalPath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", mount.LocalPath, err)
		}
		
		// Create ConfigMap name
		configMapName := fmt.Sprintf("%s-mount-%d", jobName, i)
		
		// Create ConfigMap for the file content
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: namespace,
				Labels: map[string]string{
					"rapt.dev/managed-by": "rapt",
					"rapt.dev/job":        jobName,
				},
			},
			Data: map[string]string{
				"content": string(fileContent),
			},
		}
		
		// Create the ConfigMap
		_, err = k8sClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create ConfigMap for mount %d: %w", i, err)
		}
	}

	// Create the job
	job, err := createJobFromTool(tool, toolName, args, envVars, mounts, namespace, jobName)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	// Create the job in Kubernetes
	createdJob, err := k8sClient.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create job in cluster: %w", err)
	}

	fmt.Printf("Job '%s' created successfully\n", createdJob.Name)
	fmt.Println("Streaming logs in real-time...")
	fmt.Println("Press Ctrl+C to stop following logs (job will continue running)")
	fmt.Println("=" + strings.Repeat("=", 50))

	// Always follow logs in real-time for better user experience
	return waitForJobCompletion(k8sClient, createdJob, true, timeout)
}

// getToolDefinition retrieves a tool definition from Kubernetes
func getToolDefinition(dynClient dynamic.Interface, namespace, toolName string) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "rapt.dev",
		Version:  "v1alpha1",
		Resource: "tools",
	}

	tool, err := dynClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), toolName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("tool '%s' not found in namespace '%s'", toolName, namespace)
	}

	return tool, nil
}

// createJobFromTool creates a Kubernetes Job from a tool definition
func createJobFromTool(tool *unstructured.Unstructured, toolName string, args map[string]string, envVars map[string]string, mounts []MountSpec, namespace, jobName string) (*batchv1.Job, error) {
	// Extract tool spec
	spec, found, err := unstructured.NestedMap(tool.Object, "spec")
	if err != nil || !found {
		return nil, fmt.Errorf("invalid tool spec")
	}

	jobTemplate, found, err := unstructured.NestedMap(spec, "jobTemplate")
	if err != nil || !found {
		return nil, fmt.Errorf("tool spec missing jobTemplate")
	}

	// Extract image
	image, found, err := unstructured.NestedString(jobTemplate, "image")
	if err != nil || !found {
		return nil, fmt.Errorf("tool spec missing image")
	}

	// Extract command
	var command []string
	if cmd, found, err := unstructured.NestedStringSlice(jobTemplate, "command"); err == nil && found {
		command = cmd
	}

	// Extract existing environment variables
	var env []corev1.EnvVar
	if envList, found, err := unstructured.NestedSlice(jobTemplate, "env"); err == nil && found {
		for _, envItem := range envList {
			if envMap, ok := envItem.(map[string]interface{}); ok {
				name, _ := envMap["name"].(string)
				value, _ := envMap["value"].(string)
				if name != "" {
					env = append(env, corev1.EnvVar{
						Name:  name,
						Value: value,
					})
				}
			}
		}
	}

	// Add runtime environment variables
	for key, value := range envVars {
		env = append(env, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	// Handle file mounts
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume
	
	for i, mount := range mounts {
		// Add volume mount
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      fmt.Sprintf("mount-%d", i),
			MountPath: mount.ContainerPath,
			SubPath:   "content",
		})
		
		// Add volume (ConfigMap name will be set later)
		volumes = append(volumes, corev1.Volume{
			Name: fmt.Sprintf("mount-%d", i),
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "", // Will be set later
					},
				},
			},
		})
	}

	// Build command arguments from tool arguments
	var jobArgs []string
	if toolArgs, found, err := unstructured.NestedSlice(spec, "arguments"); err == nil && found {
		for _, argItem := range toolArgs {
			if argMap, ok := argItem.(map[string]interface{}); ok {
				argName, _ := argMap["name"].(string)
				if argName != "" {
					if value, exists := args[argName]; exists {
						jobArgs = append(jobArgs, value)
					} else if defaultValue, hasDefault := argMap["default"].(string); hasDefault && defaultValue != "" {
						jobArgs = append(jobArgs, defaultValue)
					} else if _, isRequired := argMap["required"].(bool); isRequired {
						return nil, fmt.Errorf("required argument '%s' not provided", argName)
					}
				}
			}
		}
	}


	// Update volume references with ConfigMap names
	for i := range mounts {
		configMapName := fmt.Sprintf("%s-mount-%d", jobName, i)
		volumes[i].VolumeSource.ConfigMap.LocalObjectReference.Name = configMapName
	}

	// Create the Job
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			Labels: map[string]string{
				"rapt.dev/tool": toolName,
				"rapt.dev/managed-by": "rapt",
			},
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: int32Ptr(300), // Clean up after 5 minutes
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:         "tool",
							Image:        image,
							Command:      command,
							Args:         jobArgs,
							Env:          env,
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}

	return job, nil
}

// waitForJobCompletion waits for a job to complete and optionally follows logs
func waitForJobCompletion(k8sClient *kubernetes.Clientset, job *batchv1.Job, follow bool, timeout int) error {
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	// Watch for job status changes
	watcher, err := k8sClient.BatchV1().Jobs(job.Namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", job.Name),
	})
	if err != nil {
		return fmt.Errorf("failed to watch job: %w", err)
	}
	defer watcher.Stop()

	// Start log following if requested
	if follow {
		go followJobLogs(k8sClient, job)
	}

	// Wait for job completion
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Modified:
			updatedJob, ok := event.Object.(*batchv1.Job)
			if !ok {
				continue
			}

			// Check if job is complete
			if updatedJob.Status.Succeeded > 0 {
				fmt.Printf("Job '%s' completed successfully\n", job.Name)
				return nil
			}

			if updatedJob.Status.Failed > 0 {
				fmt.Printf("Job '%s' failed\n", job.Name)
				return fmt.Errorf("job failed")
			}
		case watch.Error:
			return fmt.Errorf("error watching job: %v", event.Object)
		}
	}

	return fmt.Errorf("job watch ended unexpectedly")
}

// followJobLogs follows the logs of a job
func followJobLogs(k8sClient *kubernetes.Clientset, job *batchv1.Job) {
	// Wait a bit for the pod to be created
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		time.Sleep(1 * time.Second)
		
		// Get the pod for this job
		pods, err := k8sClient.CoreV1().Pods(job.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", job.Name),
		})
		if err != nil {
			fmt.Printf("Error getting pods for job: %v\n", err)
			return
		}

		if len(pods.Items) > 0 {
			pod := pods.Items[0]
			
			// Check if pod is ready to stream logs
			if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
				// Follow logs
				logs, err := k8sClient.CoreV1().Pods(job.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
					Follow: true,
				}).Stream(context.TODO())
				if err != nil {
					fmt.Printf("Error getting logs: %v\n", err)
					return
				}
				defer logs.Close()

				// Stream logs
				buffer := make([]byte, 1024)
				for {
					n, err := logs.Read(buffer)
					if n > 0 {
						fmt.Print(string(buffer[:n]))
					}
					if err != nil {
						break
					}
				}
				return
			}
		}
	}
	
	fmt.Printf("Timeout waiting for pod to be ready for job '%s'\n", job.Name)
}

// int32Ptr returns a pointer to an int32
func int32Ptr(i int32) *int32 { return &i }