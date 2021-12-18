package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cschleiden/go-dt/pkg/backend"
	"github.com/cschleiden/go-dt/pkg/backend/memory"
	"github.com/cschleiden/go-dt/pkg/client"
	"github.com/cschleiden/go-dt/pkg/sync"
	"github.com/cschleiden/go-dt/pkg/worker"
	"github.com/cschleiden/go-dt/pkg/workflow"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	mb := memory.NewMemoryBackend()

	// Run worker
	go RunWorker(ctx, mb)

	// Start workflow via client
	c := client.NewTaskHubClient(mb)

	startWorkflow(ctx, c)

	c2 := make(chan os.Signal, 1)
	<-c2
}

func startWorkflow(ctx context.Context, c client.TaskHubClient) {
	wf, err := c.CreateWorkflowInstance(ctx, client.WorkflowInstanceOptions{
		InstanceID: uuid.NewString(),
	}, Workflow1, "Hello world")
	if err != nil {
		panic("could not start workflow")
	}

	log.Println("Started workflow", wf.GetInstanceID())
}

func RunWorker(ctx context.Context, mb backend.Backend) {
	w := worker.NewWorker(mb)

	w.RegisterWorkflow("wf1", Workflow1)

	w.RegisterActivity("a1", Activity1)

	if err := w.Start(ctx); err != nil {
		panic("could not start worker")
	}
}

func Workflow1(ctx workflow.Context, msg string) (string, error) {
	log.Println("Entering Workflow1")
	log.Println("\tWorkflow instance input:", msg)
	log.Println("\tIsReplaying:", ctx.Replaying())

	defer func() {
		log.Println("Leaving Workflow1")
	}()

	a1, err := ctx.ExecuteActivity("a1", 35, 12)
	if err != nil {
		panic("error executing activity 1")
	}

	t, err := ctx.ScheduleTimer(5 * time.Second)
	if err != nil {
		panic("could not schedule timer")
	}

	s := ctx.NewSelector()

	s.AddFuture(t, func(f sync.Future) {
		log.Println("Timer fired")
		log.Println("\tIsReplaying:", ctx.Replaying())
	})

	s.AddFuture(a1, func(f sync.Future) {
		var r int
		f.Get(&r)
		log.Println("Result:", r)
		log.Println("\tIsReplaying:", ctx.Replaying())
	})

	s.Select()
	s.Select()

	return "result", nil
}

func Activity1(ctx context.Context, a, b int) (int, error) {
	log.Println("Entering Activity1")

	time.Sleep(6 * time.Second)

	defer func() {
		log.Println("Leaving Activity1")
	}()

	return a + b, nil
}
