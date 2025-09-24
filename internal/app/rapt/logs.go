package rapt

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"codeberg.org/lig/rapt/internal/k8s"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// JobInfo represents information about a job for display
type JobInfo struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Completed time.Time `json:"completed,omitempty"`
	Duration  string    `json:"duration,omitempty"`
	PodName   string    `json:"pod_name,omitempty"`
}

// ShowLogs displays logs for a tool or lists previous job runs
func ShowLogs(namespace, toolName, jobName string, follow bool, tail int) error {
	// Initialize Kubernetes client
	k8sClient, err := k8s.InitKubernetesClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to initialize kubernetes client: %w", err)
	}

	if jobName == "" {
		// List previous runs for the tool
		return listJobRuns(k8sClient, namespace, toolName)
	} else {
		// Show logs for specific job
		return showJobLogs(k8sClient, namespace, jobName, follow, tail)
	}
}

// listJobRuns lists all previous job runs for a tool
func listJobRuns(k8sClient *kubernetes.Clientset, namespace, toolName string) error {
	// Get all jobs with the tool label
	jobs, err := k8sClient.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("rapt.dev/tool=%s", toolName),
	})
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	if len(jobs.Items) == 0 {
		fmt.Printf("No previous runs found for tool '%s'\n", toolName)
		return nil
	}

	// Convert to JobInfo and sort by creation time (newest first)
	jobInfos := make([]JobInfo, len(jobs.Items))
	for i, job := range jobs.Items {
		jobInfos[i] = convertJobToInfo(job)
	}

	// Sort by creation time (newest first)
	sort.Slice(jobInfos, func(i, j int) bool {
		return jobInfos[i].Created.After(jobInfos[j].Created)
	})

	// Display in table format
	fmt.Printf("Previous runs for tool '%s':\n\n", toolName)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "JOB NAME\tSTATUS\tCREATED\tDURATION")
	
	for _, job := range jobInfos {
		created := job.Created.Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", job.Name, job.Status, created, job.Duration)
	}
	w.Flush()

	fmt.Printf("\nTo view logs for a specific job, run:\n")
	fmt.Printf("  rapt logs %s <job-name>\n", toolName)
	fmt.Printf("\nTo follow logs for the latest job, run:\n")
	fmt.Printf("  rapt logs %s %s --follow\n", toolName, jobInfos[0].Name)

	return nil
}

// showJobLogs displays logs for a specific job
func showJobLogs(k8sClient *kubernetes.Clientset, namespace, jobName string, follow bool, tail int) error {
	// Get the job
	job, err := k8sClient.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("job '%s' not found: %w", jobName, err)
	}

	// Get the pod for this job
	pods, err := k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})
	if err != nil {
		return fmt.Errorf("failed to get pods for job: %w", err)
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for job '%s'", jobName)
	}

	pod := pods.Items[0]

	// Check if pod is ready for log streaming
	if pod.Status.Phase == corev1.PodPending {
		fmt.Printf("Job '%s' is still starting up. Waiting for pod to be ready...\n", jobName)
		
		// Wait for pod to be ready
		maxRetries := 60 // Wait up to 60 seconds
		for i := 0; i < maxRetries; i++ {
			time.Sleep(1 * time.Second)
			
			updatedPod, err := k8sClient.CoreV1().Pods(namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get pod status: %w", err)
			}
			
			if updatedPod.Status.Phase != corev1.PodPending {
				pod = *updatedPod
				break
			}
		}
		
		if pod.Status.Phase == corev1.PodPending {
			return fmt.Errorf("timeout waiting for pod to be ready")
		}
	}

	// Prepare log options
	logOptions := &corev1.PodLogOptions{
		Follow: follow,
	}
	
	if tail > 0 {
		tailLines := int64(tail)
		logOptions.TailLines = &tailLines
	}

	// Get logs
	logs, err := k8sClient.CoreV1().Pods(namespace).GetLogs(pod.Name, logOptions).Stream(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to get logs: %w", err)
	}
	defer logs.Close()

	// Display job info
	fmt.Printf("Job: %s\n", jobName)
	fmt.Printf("Pod: %s\n", pod.Name)
	fmt.Printf("Status: %s\n", getJobStatus(job))
	fmt.Printf("Created: %s\n", job.CreationTimestamp.Format("2006-01-02 15:04:05"))
	if follow {
		fmt.Println("Following logs in real-time (Press Ctrl+C to stop)...")
	}
	fmt.Println(strings.Repeat("=", 50))

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

	return nil
}

// convertJobToInfo converts a Kubernetes Job to JobInfo
func convertJobToInfo(job batchv1.Job) JobInfo {
	info := JobInfo{
		Name:    job.Name,
		Status:  getJobStatus(&job),
		Created: job.CreationTimestamp.Time,
	}

	// Calculate duration if job is completed
	if job.Status.CompletionTime != nil {
		info.Completed = job.Status.CompletionTime.Time
		duration := info.Completed.Sub(info.Created)
		info.Duration = formatDuration(duration)
	}

	return info
}

// getJobStatus determines the status of a job
func getJobStatus(job *batchv1.Job) string {
	if job.Status.Succeeded > 0 {
		return "✅ Succeeded"
	}
	if job.Status.Failed > 0 {
		return "❌ Failed"
	}
	if job.Status.Active > 0 {
		return "🔄 Running"
	}
	return "⏳ Pending"
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm %.0fs", d.Minutes(), d.Seconds()-d.Truncate(time.Minute).Seconds())
	}
	return fmt.Sprintf("%.0fh %.0fm", d.Hours(), d.Minutes()-d.Truncate(time.Hour).Minutes())
}